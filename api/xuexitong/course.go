package xuexitong

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yatori-go-coreV3/interfaces"
)

// CourseListApi 拉取对应账号的课程数据
func CourseListApi(user interfaces.IUser) (string, error) {

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, ApiPullCourses, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("Cookie", user.GetCookie())
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func DetailApi(cookie string, cpi, key string) (string, error) {
	method := "GET"

	params := url.Values{}
	params.Add("id", key)
	params.Add("personid", cpi)
	params.Add("fields", "id,bbsid,classscore,isstart,allowdownload,chatid,name,state,isfiled,visiblescore,hideclazz,begindate,forbidintoclazz,coursesetting.fields(id,courseid,hiddencoursecover,coursefacecheck),course.fields(id,belongschoolid,name,infocontent,objectid,app,bulletformat,mappingcourseid,imageurl,teacherfactor,jobcount,knowledge.fields(id,name,indexOrder,parentnodeid,status,isReview,layer,label,jobcount,begintime,endtime,attachment.fields(id,type,objectid,extension).type(video)))")
	params.Add("view", "json")

	client := &http.Client{}
	req, err := http.NewRequest(method, ApiPullChapter+"?"+params.Encode(), nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("User-Agent", " Dalvik/2.1.0 (Linux; U; Android 12; SM-N9006 Build/70e2a6b.1) (schild:e9b05c3f9fb49fef2f516e86ac3c4ff1) (device:SM-N9006) Language/zh_CN com.chaoxing.mobile/ChaoXingStudy_3_6.3.7_android_phone_10822_249 (@Kalimdor)_4627cad9c4b6415cba5dc6cac39e6c96")
	req.Header.Add("Accept-Language", " zh_CN")
	req.Header.Add("Host", " mooc1-api.chaoxing.com")
	req.Header.Add("Connection", " Keep-Alive")
	req.Header.Add("Cookie", cookie)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil
}

func DetailPointStatusApi(cookie, key, userID, cpi, courseID string, keyList []int) (string, error) {
	method := "POST"

	strInts := make([]string, len(keyList))
	for i, v := range keyList {
		strInts[i] = fmt.Sprintf("%d", v)
	}

	ts := time.Now().UnixNano() / 1000000
	join := strings.Join(strInts, ",")
	values := url.Values{
		"view":     {"json"},
		"nodes":    {join},
		"clazzid":  {key},
		"time":     {strconv.FormatInt(ts, 10)},
		"userid":   {userID},
		"cpi":      {cpi},
		"courseid": {courseID},
	}
	// 编码请求体
	payload := strings.NewReader(values.Encode())

	client := &http.Client{}
	req, err := http.NewRequest(method, ApiChapterPoint, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "mooc1-api.chaoxing.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", cookie)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 解码响应体（假设服务器返回的内容是 ISO-8859-1 编码）
	// decodedBody, _, err := transform.Bytes(charmap.ISO8859_1.NewDecoder(), body)
	return string(body), nil
}

// FetchDetailCords 以课程序号拉取对应“章节”的任务节点卡片资源接口2
// Args:
//
//	nodes: 任务点集合 , index: 任务点索引
func FetchDetailCords(cookie, clazzid, courseid, knowledgeid, num, cpi string) (string, error) {
	//https://mooc1.chaoxing.com/mooc-ans/knowledge/cards?clazzid={_clazzid}&courseid={_courseid}&knowledgeid={_knowledgeid}&num={_possible_num}&ut=s&cpi={_cpi}&v=20160407-3&mooc2=1"
	method := "GET"
	values := url.Values{}
	values.Add("clazzid", clazzid)
	values.Add("courseid", courseid)
	values.Add("knowledgeid", knowledgeid)
	values.Add("num", num)
	values.Add("cpi", cpi)

	client := &http.Client{}
	req, err := http.NewRequest(method, "https://mooc1.chaoxing.com/mooc-ans/knowledge/cards"+"?"+values.Encode(), nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("User-Agent", " Dalvik/2.1.0 (Linux; U; Android 12; SM-N9006 Build/70e2a6b.1) (schild:e9b05c3f9fb49fef2f516e86ac3c4ff1) (device:SM-N9006) Language/zh_CN com.chaoxing.mobile/ChaoXingStudy_3_6.3.7_android_phone_10822_249 (@Kalimdor)_4627cad9c4b6415cba5dc6cac39e6c96")
	req.Header.Add("Accept-Language", " zh_CN")
	req.Header.Add("Host", " mooc1-api.chaoxing.com")
	req.Header.Add("Connection", " Keep-Alive")
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Accept", "*/*")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	//try{
	//	mArg = {"hiddenConfig":false,"isMirror":false,"attachments":[{"otherInfo":"nodeId_188905976-cpi_414911026","authEnc":"d155dab728fd219baad655566da436c5","mid":"15060881566971645856824803","liveDragEnc":"936f2f0943745b0c05680eb0cbd52411","type":"live","begins":0,"liveSetEnc":"1a95837d77f659cf0a447800ead23df0","jobid":"live-6000122221272107","isNotDrag":1,"ends":0,"property":{"vdoid":"vdoid94295833C1CkX","jobid":"live-6000122221272107","module":"insertlive","mid":"15060881566971645856824803","title":"第一章毛泽东思想及历史地位","courseid":206105634,"userId":"168718513","liveId":6000122221272107,"streamName":"NEWLIVEY6237350vdoid94295833C1CkX","live":true,"liveStatus":"直播回看","_jobid":"live-6000122221272107"},"enc":"e799afd3b75ac94c7661eb392e3657e3","job":true,"aid":1001745669,"liveSwDsEnc":"feb6a3c5d1979cdfeeab265a703df279"},{"otherInfo":"nodeId_188905976-cpi_414911026","authEnc":"aac74e13474ecb7ef0bb98d2fe0d3ea6","mid":"12143258042351645617956623","liveDragEnc":"e44ccdd68ff0d6791ee2e6506067b698","type":"live","begins":0,"liveSetEnc":"6f84a1a9a69fe6a9cda2677f97bdc6c5","jobid":"live-1000122037998115","isNotDrag":1,"ends":0,"property":{"fid":"6776","module":"insertlive","mid":"12143258042351645617956623","title":"毛泽东思想和中国特色社会主义理论体系概论","userId":"39241191","liveId":1000122037998115,"streamName":"NEWLIVE06X89W91vdoid94113705413D4","vdoid":"vdoid94113705413D4","jobid":"live-1000122037998115","courseid":206105634,"live":true,"liveStatus":"未开始","_jobid":"live-1000122037998115"},"enc":"eabe3704f60d54575be2214c3a52eb35","job":true,"aid":1001745670,"liveSwDsEnc":"59b483c39c05c8bfa66220f92b0acf70"},{"headOffset":124000,"otherInfo":"nodeId_188905976-cpi_414911026-rt_d-ds_0-ff_1-be_0_0-vt_1-v_6-enc_f6a7026e045bd57ce354d1d2336d424e&courseId=206105634","isPassed":false,"mid":"7704186497151561693751561","jumpTimePointList":[],"type":"video","begins":0,"jobid":"1558340423438519","customType":0,"attDurationEnc":"0a817830625be8e8dbb95690921599d5","videoFaceCaptureEnc":"11ea9ae486200527c3d006b4c3ef15f9","ends":0,"randomCaptureTime":658,"property":{"jobid":"1558340423438519","switchwindow":"true","size":395947998,"fastforward":"true","hsize":"377.61 MB","module":"insertvideo","name":"1.2.mp4","mid":"7704186497151561693751561","type":".mp4","doublespeed":0,"objectid":"e26ecd3d3dbcf659afa5dbcc7a2de5ca","_jobid":"1558340423438519"},"playTime":97000,"attDuration":725,"headOffsetVersion":0,"job":true,"aid":1001745671,"objectId":"e26ecd3d3dbcf659afa5dbcc7a2de5ca"}],"coursename":"毛泽东思想和中国特色社会主义理论体系概论","defaults":{"fid":"6776","ktoken":"138f8d3800775b88333e2b109080b2ef","mtEnc":"1c0ad1e39e2a74a3b026784e0c0089d7","appInfo":"","playingCapture":1,"videoAutoPlay":0,"userid":"348514942","reportTimeInterval":60,"showVideoWaterMark":0,"schooldoublespeed":0,"endCapture":0,"defenc":"60b8104b56caf3e40a199adb336f894c","cardid":169812899,"imageUrl":"https://p.ananas.chaoxing.com/star3/270_169c/2c77783bb5c4c8c4f8aeae29903d326b.png","state":0,"cpi":414911026,"captureInterval":0,"playAginCapture":0,"startCapture":1,"isFiled":0,"ignoreVideoCtrl":0,"reportUrl":"https://mooc1.chaoxing.com/mooc-ans/multimedia/log/a/414911026","chapterCapture":0,"initdataUrl":"https://mooc1.chaoxing.com/mooc-ans/richvideo/initdatawithviewer","cFid":"46175","knowledgeid":188905976,"videoTopicCloud":0,"qnenc":"a7de7c6e1ac8a0a03487d96291111433","clazzId":115946061,"chapterCollectionType":0,"lastmodifytime":1740537950000,"aiVideoInterpret":0,"courseid":206105634,"subtitleUrl":"https://mooc1.chaoxing.com/mooc-ans/richvideo/subtitle","playingLoopCapture":1},"mooc2":0,"knowledgename":"马克思主义中国化命题的提出与科学内涵","openShowHotMap":false,"control":true,"chapterVideoTranslate":0,"isErya":1};
	//}catch(e){
	//}
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	sprintf := fmt.Sprintf(`mArg = ([^;]{6,})`)
	compile := regexp.MustCompile(sprintf)
	find := compile.FindAllStringSubmatch(string(body), -1)
	for _, v := range find {
		return v[1], nil
	}
	return string(body), nil
}
