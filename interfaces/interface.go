package interfaces

type IUser interface {
	Login() error
	GetAccount() string
	GetPreUrl() string
	GetPassword() string
	GetCookie() string
	UserInfo() (map[string]any, error) // 可获取用户参数（如userID、username）
	CacheData() (map[string]any, error)
}

type ICourse interface {
	GetID() string
	GetCourseID() string
	GetName() string
	GetDetail() []IDetail
}

type IDetail interface {
	GetVideo()
	GetWork()
}

type ICourseList interface {
	CourseList(user IUser) ([]ICourse, error)
}
