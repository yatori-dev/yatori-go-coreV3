package yatori

import (
	"errors"
	"yatori-go-coreV3/interfaces"
	"yatori-go-coreV3/strategy"
)

type YatoriUser struct {
	Account  string
	CacheMap map[string]any
	Password string
	PreUrl   string

	cookie string

	strategy  interfaces.IUser // 持有当前策略实例（核心！）
	courseSvc interfaces.ICourseList
}

func NewUser(account, password, url string) *YatoriUser {
	y := &YatoriUser{
		Account:  account,
		Password: password,
		PreUrl:   url,
	}

	return y
}

func (y *YatoriUser) On(accountType string) error {
	switch accountType {
	case "XUEXITONG":
		y.strategy = strategy.NewXueXiTUserStrategy(y)
		y.courseSvc = strategy.NewXueXiTCourse()
	default:
		return errors.New("不支持的账号类型")
	}
	return nil
}

func (y *YatoriUser) Login() error {
	if y.strategy == nil {
		return errors.New("请先通过On方法指定账号类型")
	}
	return y.strategy.Login()
}

func (y *YatoriUser) GetAccount() string {
	return y.Account
}

func (y *YatoriUser) GetPassword() string {
	return y.Password
}

func (y *YatoriUser) GetPreUrl() string {
	return y.PreUrl
}

func (y *YatoriUser) GetCookie() string {
	return y.strategy.GetCookie()
}

// UserInfo 统一用户信息入口
func (y *YatoriUser) UserInfo() (map[string]any, error) {
	if y.strategy == nil {
		return nil, errors.New("请先通过On方法指定账号类型")
	}
	return y.strategy.UserInfo()
}

func (y *YatoriUser) CacheData() (map[string]any, error) {
	if y.strategy == nil {
		return nil, errors.New("请先通过On方法指定账号类型")
	}
	return y.strategy.CacheData()
}

func (y *YatoriUser) CourseList() ([]interfaces.ICourse, error) {
	if y.courseSvc == nil {
		return nil, errors.New("课程服务未初始化")
	}
	return y.courseSvc.CourseList(y)
}
