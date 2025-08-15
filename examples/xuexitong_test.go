package examples

import (
	"testing"
	"yatori-go-coreV3/common"
	"yatori-go-coreV3/global"
	"yatori-go-coreV3/yatori"
)

func setup() {
	common.InitConfig("./")
}

func TestLogin(t *testing.T) {

	setup()
	gUser := global.Config.Users[0]
	user := yatori.NewUser(gUser.Account, gUser.Password, "")
	err := user.On(gUser.AccountType)
	if err != nil {
		t.Error(err)
	}
	err = user.Login()
	courseList, err := user.CourseList()
	if err != nil {
		t.Error(err)
	}
	for _, course := range courseList {
		t.Log(course.GetID(), course.GetName(), course.GetCourseID(), course.GetUserID())
		course.GetDetail()
	}
}
