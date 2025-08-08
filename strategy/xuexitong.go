package strategy

import (
	"encoding/json"
	"errors"
	"github.com/thedevsaddam/gojsonq"
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
	err = json.Unmarshal([]byte(courses), &xueXiTCourse)
	if err != nil {
		log2.Print(log2.INFO, "["+user.GetAccount()+"] "+" 解析失败")
		panic(err)
	}
	log2.Print(log2.INFO, "["+user.GetAccount()+"] "+" 课程数量："+strconv.Itoa(len(xueXiTCourse.ChannelList)))
	// log2.Print(log2.INFO, "["+cache.Name+"] "+courses)

	var courseList = make([]XueXiTCourse, 0)
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
			courseImage  string
		)

		for _, v := range channel.Content.Course.Data {
			teacher = v.Teacherfactor
			courseName = v.Name
			courseDataID = v.Id
			userID := strings.Split(v.CourseSquareUrl, "userId=")[1]
			x.UserID = userID
			classId = strings.Split(strings.Split(v.CourseSquareUrl, "classId=")[1], "&userId")[0]
			courseID = strings.Split(strings.Split(v.CourseSquareUrl, "courseId=")[1], "&personId")[0]
			courseImage = v.Imageurl
		}

		course := XueXiTCourse{
			Cpi:           channel.Cpi,
			Key:           classId,
			CourseID:      courseID,
			ChatID:        channel.Content.Chatid,
			CourseTeacher: teacher,
			CourseName:    courseName,
			CourseImage:   courseImage,
			CourseDataID:  courseDataID,
			ContentID:     channel.Content.Id,
		}
		for _, course := range courseList {
			if course.CourseID == courseID {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		courseList = append(courseList, course)
	}
	var iCourses []interfaces.ICourse
	for _, course := range courseList {
		iCourses = append(iCourses, &course)
	}
	return iCourses, nil
}

func (x *XueXiTCourse) GetID() string {
	return x.Key
}

func (x *XueXiTCourse) GetName() string {
	return x.CourseName
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
