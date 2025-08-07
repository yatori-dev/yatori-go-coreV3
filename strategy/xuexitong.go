package strategy

import (
	"errors"
	"github.com/thedevsaddam/gojsonq"
	"yatori-go-coreV3/api/xuexitong"
	"yatori-go-coreV3/interfaces"
	log2 "yatori-go-coreV3/utils/log"
)

type XueXiTUserStrategy struct {
	cookie string
	User   interfaces.IUser
}

// NewXueXiTUserStrategy 创建策略实例（接收接口类型）
func NewXueXiTUserStrategy(user interfaces.IUser) *XueXiTUserStrategy {
	return &XueXiTUserStrategy{User: user}
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

func (n *XueXiTUserStrategy) GetPreUrl() string {
	return n.User.GetPreUrl()
}
