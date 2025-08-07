package interfaces

type IUser interface {
	Login() error
	GetAccount() string
	GetPreUrl() string
	GetPassword() string
	UserInfo() (map[string]any, error) // 可获取用户参数（如userID、username）
	CacheData() (map[string]any, error)
}

type ICourse interface {
	GetID() string
	GetName() string
}

type ICourseList interface {
	CourseList() ([]ICourse, error)
}
