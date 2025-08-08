package xuexitong

import (
	"io/ioutil"
	"net/http"
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
