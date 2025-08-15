package strategy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"log"
	"sort"
	"strconv"
	"strings"
	"yatori-go-coreV3/api/xuexitong"
	"yatori-go-coreV3/interfaces"
	log2 "yatori-go-coreV3/utils/log"
)

type XueXiTUserStrategy struct {
	UserID string
	cookie string
	User   interfaces.IUser
}

type XueXiTCourse struct {
	UserID        string `json:"userId"`
	Cpi           int    `json:"cpi"`      // 用户唯一标识
	Key           string `json:"key"`      // classID 在课程API中为key
	CourseID      string `json:"courseId"` // 课程ID
	ChatID        string `json:"chatId"`
	CourseTeacher string `json:"courseTeacher"` // 课程老师
	CourseName    string `json:"courseName"`    //课程名
	CourseImage   string `json:"courseImage"`
	// 两个标识 暂时不知道有什么用
	CourseDataID int `json:"courseDataId"`
	ContentID    int `json:"ContentID"`

	Detail DetailList

	cookie string
}

// XueXiTCourseJson 课程所有信息
type XueXiTCourseJson struct {
	Result           int           `json:"result"`
	Msg              string        `json:"msg"`
	ChannelList      []ChannelItem `json:"channelList"`
	Mcode            string        `json:"mcode"`
	Createcourse     int           `json:"createcoursed"`
	TeacherEndCourse int           `json:"teacherEndCourse"`
	ShowEndCourse    int           `json:"showEndCourse"`
	HasMore          bool          `json:"hasMore"`
	StuEndCourse     int           `json:"stuEndCourse"`
}

// ChannelItem 课程列表
type ChannelItem struct {
	Cfid     int    `json:"cfid"`
	Norder   int    `json:"norder"`
	CataName string `json:"cataName"`
	Cataid   string `json:"cataid"`
	Id       int    `json:"id"`
	Cpi      int    `json:"cpi"`
	Key      any    `json:"key"`
	Content  struct {
		Studentcount int    `json:"studentcount"`
		Chatid       string `json:"chatid"`
		IsFiled      int    `json:"isFiled"`
		Isthirdaq    int    `json:"isthirdaq"`
		Isstart      bool   `json:"isstart"`
		Isretire     int    `json:"isretire"`
		Name         string `json:"name"`
		Course       struct {
			Data []struct {
				BelongSchoolId     string `json:"belongSchoolId"`
				Coursestate        int    `json:"coursestate"`
				Teacherfactor      string `json:"teacherfactor"`
				IsCourseSquare     int    `json:"isCourseSquare"`
				Schools            string `json:"schools"`
				CourseSquareUrl    string `json:"courseSquareUrl"`
				Imageurl           string `json:"imageurl"`
				AppInfo            string `json:"appInfo"`
				Name               string `json:"name"`
				DefaultShowCatalog int    `json:"defaultShowCatalog"`
				Id                 int    `json:"id"`
				AppData            int    `json:"appData"`
			} `json:"data"`
		} `json:"course"`
		Roletype int    `json:"roletype"`
		Id       int    `json:"id"`
		State    int    `json:"state"`
		Cpi      int    `json:"cpi"`
		Bbsid    string `json:"bbsid"`
		IsSquare int    `json:"isSquare"`
	} `json:"content"`
	Topsign int `json:"topsign"`
}

type DetailList struct {
	ChatID    string          `json:"chatid"`
	Knowledge []KnowledgeItem `json:"knowledge"`
}

// KnowledgeItem 结构体用于存储 knowledge 中的每个项目
type KnowledgeItem struct {
	JobCount      int           `json:"jobcount"` // 作业数量
	IsReview      int           `json:"isreview"` // 是否为复习
	Attachment    []interface{} `json:"attachment"`
	IndexOrder    int           `json:"indexorder"` // 节点顺序
	Name          string        `json:"name"`       // 章节名称
	ID            int           `json:"id"`
	Label         string        `json:"label"`        // 节点标签
	Layer         int           `json:"layer"`        // 节点层级
	ParentNodeID  int           `json:"parentnodeid"` // 父节点 ID
	Status        string        `json:"status"`       // 节点状态
	PointTotal    int
	PointFinished int
}

type ChapterPointDTO map[string]struct {
	ClickCount    int `json:"clickcount"`    // 是否还有节点
	FinishCount   int `json:"finishcount"`   // 已完成节点
	TotalCount    int `json:"totalcount"`    // 总节点
	OpenLock      int `json:"openlock"`      // 是否有锁
	UnFinishCount int `json:"unfinishcount"` // 未完成节点
}

// NewXueXiTUserStrategy 创建策略实例（接收接口类型）
func NewXueXiTUserStrategy(user interfaces.IUser) *XueXiTUserStrategy {
	return &XueXiTUserStrategy{User: user}
}

// NewXueXiTCourse 创建课程实例
func NewXueXiTCourse() interfaces.ICourseList {
	return &XueXiTCourse{}
}

func (x *XueXiTCourse) CourseList(user interfaces.IUser) ([]interfaces.ICourse, error) {
	courses, err := xuexitong.CourseListApi(user)
	if err != nil {
		log2.Print(log2.INFO, "["+user.GetAccount()+"] "+" 拉取失败")
	}
	var xueXiTCourse XueXiTCourseJson
	var iCourses []interfaces.ICourse
	err = json.Unmarshal([]byte(courses), &xueXiTCourse)
	if err != nil {
		log2.Print(log2.INFO, "["+user.GetAccount()+"] "+" 解析失败")
		panic(err)
	}
	log2.Print(log2.INFO, "["+user.GetAccount()+"] "+" 课程数量："+strconv.Itoa(len(xueXiTCourse.ChannelList)))

	for i, channel := range xueXiTCourse.ChannelList {
		var flag = false
		if channel.Content.Course.Data == nil && i >= 0 && i < len(xueXiTCourse.ChannelList) {
			xueXiTCourse.ChannelList = append(xueXiTCourse.ChannelList[:i], xueXiTCourse.ChannelList[i+1:]...)
			continue
		}
		var (
			teacher      string
			courseName   string
			courseDataID int
			classId      string
			courseID     string
			userID       string
			courseImage  string
		)

		for _, v := range channel.Content.Course.Data {
			teacher = v.Teacherfactor
			courseName = v.Name
			courseDataID = v.Id
			userID = strings.Split(v.CourseSquareUrl, "userId=")[1]
			classId = strings.Split(strings.Split(v.CourseSquareUrl, "classId=")[1], "&userId")[0]
			courseID = strings.Split(strings.Split(v.CourseSquareUrl, "courseId=")[1], "&personId")[0]
			courseImage = v.Imageurl
		}

		course := XueXiTCourse{
			Cpi:           channel.Cpi,
			Key:           classId,
			CourseID:      courseID,
			UserID:        userID,
			ChatID:        channel.Content.Chatid,
			CourseTeacher: teacher,
			CourseName:    courseName,
			CourseImage:   courseImage,
			CourseDataID:  courseDataID,
			ContentID:     channel.Content.Id,

			cookie: user.GetCookie(),
		}

		for _, c := range iCourses {
			if c.GetCourseID() == courseID {
				flag = true
				break
			}
		}

		if flag {
			continue
		}
		iCourses = append(iCourses, &course)
	}

	return iCourses, nil
}

func (x *XueXiTCourse) SetCookie(cookie string) {
	x.cookie = cookie
}

// GetID 在学习通中为Key
func (x *XueXiTCourse) GetID() string {
	return x.Key
}

func (x *XueXiTCourse) GetUserID() string {
	return x.UserID
}

func (x *XueXiTCourse) GetCpi() string {
	return strconv.Itoa(x.Cpi)
}

func (x *XueXiTCourse) GetName() string {
	return x.CourseName
}

func (x *XueXiTCourse) GetCourseID() string {
	return x.CourseID
}

func (x *XueXiTCourse) GetDetail() []interfaces.IDetail {
	var keyList []int
	loglevel := log2.INFO
	detail, err := xuexitong.DetailApi(x.cookie, strconv.Itoa(x.Cpi), x.Key)
	if err != nil {
		log2.Print(loglevel, "["+x.GetName()+"] "+" 拉取失败")
		return nil
	}

	var detailMap map[string]interface{}
	err = json.Unmarshal([]byte(detail), &detailMap)
	if err != nil {
		log2.Print(loglevel, "["+x.GetName()+"] "+"["+x.GetID()+"] "+" 解析失败:"+err.Error())
		return nil
	}

	chapterMapJson, err := json.Marshal(detailMap["data"])
	if len(chapterMapJson) == 2 {
		log2.Print(loglevel, "["+x.GetName()+"] "+"["+x.GetID()+"] "+" 课程获取失败")
		return nil
	}
	// 解析 JSON 数据为 map 切片
	var chapterData []map[string]interface{}
	if err := json.Unmarshal(chapterMapJson, &chapterData); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	chatid := chapterData[0]["chatid"].(string)
	// 提取 knowledge
	var knowledgeData []map[string]interface{}
	course, ok := chapterData[0]["course"].(map[string]interface{})
	if !ok {
		log2.Print(loglevel, "["+x.GetName()+"] "+"["+x.GetID()+"] "+" 无法提取 course")
		return nil
	}
	data, ok := course["data"].([]interface{})
	if !ok {
		log2.Print(loglevel, "["+x.GetName()+"] "+"["+x.GetID()+"] "+" 无法提取 course data")
		return nil
	}
	if len(data) > 0 {
		knowledge, ok := data[0].(map[string]interface{})["knowledge"].(map[string]interface{})["data"].([]interface{})
		if !ok {
			log2.Print(loglevel, "["+x.GetName()+"] "+"["+x.GetID()+"] "+" 无法提取 knowledge data")
			return nil
		}
		for _, item := range knowledge {
			knowledgeMap := item.(map[string]interface{})
			knowledgeData = append(knowledgeData, knowledgeMap)
		}
	} else {
		log2.Print(loglevel, "["+x.GetName()+"] "+"["+x.GetID()+"] "+" course data 为空")
		return nil
	}

	// 将提取的数据封装到 CourseInfo 结构体中
	var knowledgeItems []KnowledgeItem
	for _, item := range knowledgeData {
		knowledgeItem := KnowledgeItem{
			JobCount:     int(item["jobcount"].(float64)),
			IsReview:     int(item["isreview"].(float64)),
			Attachment:   item["attachment"].(map[string]interface{})["data"].([]interface{}),
			IndexOrder:   int(item["indexorder"].(float64)),
			Name:         item["name"].(string),
			ID:           int(item["id"].(float64)),
			Label:        item["label"].(string),
			Layer:        int(item["layer"].(float64)),
			ParentNodeID: int(item["parentnodeid"].(float64)),
			Status:       item["status"].(string),
		}
		knowledgeItems = append(knowledgeItems, knowledgeItem)
	}
	x.Detail = DetailList{
		ChatID:    chatid,
		Knowledge: knowledgeItems,
	}
	if len(x.Detail.Knowledge) == 0 {
		log2.Print(loglevel, "["+x.GetName()+"] "+"["+x.GetID()+"] "+" 课程章节为空")
		return nil
	}
	// 按照任务点节点重排顺序
	sort.Slice(x.Detail.Knowledge, func(i, j int) bool {
		iLabelParts := strings.Split(x.Detail.Knowledge[i].Label, ".")
		jLabelParts := strings.Split(x.Detail.Knowledge[j].Label, ".")
		for k := range iLabelParts {
			if k >= len(jLabelParts) {
				return false // i has more parts, so it should come after j
			}
			iv, _ := strconv.Atoi(iLabelParts[k])
			jv, _ := strconv.Atoi(jLabelParts[k])
			if iv != jv {
				return iv < jv
			}
		}
		return len(iLabelParts) < len(jLabelParts)
	})
	log2.Print(loglevel, "["+x.GetName()+"] "+"获取课程章节成功 (共 ", log2.Yellow, strconv.Itoa(len(x.Detail.Knowledge)), log2.Default, " 个) ")

	for _, item := range x.Detail.Knowledge {
		keyList = append(keyList, item.ID)
	}
	status, err := xuexitong.DetailPointStatusApi(x.cookie, x.Key, x.UserID, strconv.Itoa(x.Cpi), x.CourseID, keyList)
	if err != nil || gojsonq.New().JSONString(status).Find("msg") == "用户不存在" {
		log2.Print(loglevel, "["+x.GetName()+"] "+"["+x.GetID()+"] "+" PointStatus拉取失败"+status)
		return nil
	}
	var cp ChapterPointDTO
	if err := json.NewDecoder(bytes.NewReader([]byte(status))).Decode(&cp); err != nil {
		log2.Print(loglevel, "failed to decode JSON response: %v", err)
		return nil
	}

	for i := range x.Detail.Knowledge {
		x.Detail.Knowledge[i].updatePointStatus(cp)
	}
	log2.Print(loglevel, "["+x.GetName()+"] "+"["+x.GetID()+"] "+" PointStatus更新成功")
	var details []interfaces.IDetail
	for _, item := range x.Detail.Knowledge {
		details = append(details, &item)
	}
	return details
}

func (k *KnowledgeItem) GetWork() {
	// TODO 获取作业
}

func (k *KnowledgeItem) GetVideo() {
	// TODO 获取视频
}

// Status 课程具体结构
func (x *XueXiTCourse) Status() any {
	return x
}

func (k *KnowledgeItem) StatusStruct() any {
	return k
}

// updatePointStatus 更新节点状态 单独对应ChaptersList每个KnowledgeItem
func (c *KnowledgeItem) updatePointStatus(chapterPoint ChapterPointDTO) {
	pointData, exists := chapterPoint[fmt.Sprintf("%d", c.ID)]
	if !exists {
		fmt.Printf("Chapter ID %d not found in API response\n", c.ID)
		return
	}
	// 当存在未完成节点 Item 中Total 记录数为未完成数数量
	// TotalCount == 0 没有节点 或者 属于顶级标签
	// 两种条件都不符合 则 记录此章节总结点数量
	if pointData.UnFinishCount != 0 && pointData.TotalCount == 0 {
		c.PointTotal = pointData.UnFinishCount
	} else {
		c.PointTotal = pointData.TotalCount
	}
	c.PointFinished = pointData.FinishCount
}

func (n *XueXiTUserStrategy) Login() error {
	loginStr, cookie, err := xuexitong.Login(n.User.GetAccount(), n.User.GetPassword())
	log2.Print(log2.DEBUG, "["+n.User.GetAccount()+"] "+" 登录成功", loginStr, err)
	log2.Print(log2.INFO, "["+n.User.GetAccount()+"] "+" 登录成功", loginStr, err)
	if err != nil {
		if gojsonq.New().JSONString(loginStr).Find("msg2") != nil {
			return errors.New(gojsonq.New().JSONString(loginStr).Find("msg2").(string))
		} else {
			return err
		}

	}

	n.cookie = cookie

	return nil
}

func (n *XueXiTUserStrategy) UserInfo() (map[string]any, error) {
	return map[string]any{
		"account": n.User.GetAccount(),
		"type":    "XUEXITONG",
	}, nil
}

func (n *XueXiTUserStrategy) CacheData() (map[string]any, error) {
	cacheData, err := n.User.CacheData()
	if err != nil {
		return nil, err
	}
	return cacheData, nil
}

func (n *XueXiTUserStrategy) GetAccount() string {
	return n.User.GetAccount()
}

func (n *XueXiTUserStrategy) GetPassword() string {
	return n.User.GetPassword()
}

func (n *XueXiTUserStrategy) GetCookie() string {
	return n.cookie
}

func (n *XueXiTUserStrategy) GetPreUrl() string {
	return n.User.GetPreUrl()
}
