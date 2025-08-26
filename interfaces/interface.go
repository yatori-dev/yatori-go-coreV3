package interfaces

type Str string

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
	GetUserID() string
	GetCourseID() string
	GetName() string
	GetDetail() []IDetail
	IStatusStruct
}

// IDetail 课程详情 这里可能各种课程软件就会有很大的分歧
type IDetail interface {
	GetVideo()
	GetWork()
	IStatusStruct
}

type IStatusStruct interface {
	StatusStruct() IGet
}

type IGet interface {
	Get() Str
}

type ICourseList interface {
	CourseList(user IUser) ([]ICourse, error)
}
