package xuexitong

// 注意Api类文件主需要写最原始的接口请求和最后的json的string形式返回，不需要用结构体序列化。
// 序列化和具体的功能实现请移步到Action代码文件中
const (
	ApiLoginWeb = "https://passport2.chaoxing.com/fanyalogin"

	ApiPullCourses = "https://mooc1-api.chaoxing.com/mycourse/backclazzdata"

	// ApiChapterPoint 接口-课程章节任务点状态
	ApiChapterPoint = "https://mooc1-api.chaoxing.com/job/myjobsnodesmap"
	ApiChapterCards = "https://mooc1-api.chaoxing.com/gas/knowledge"
	ApiPullChapter  = "https://mooc1-api.chaoxing.com/gas/clazz"

	// PageMobileChapterCard SSR页面-客户端章节任务卡片
	PageMobileChapterCard = "https://mooc1-api.chaoxing.com/knowledge/cards"

	// APIChapterCardResource 接口-课程章节卡片资源
	APIChapterCardResource = "https://mooc1-api.chaoxing.com/ananas/status"
	// APIVideoPlayReport 接口-视频播放上报
	APIVideoPlayReport  = "https://mooc1.chaoxing.com/mooc-ans/multimedia/log/a"
	APIVideoPlayReport2 = "https://mooc1-api.chaoxing.com/multimedia/log/a" // cxkitty的

	// ApiWorkCommit 接口-单元作业答题提交
	ApiWorkCommit = "https://mooc1-api.chaoxing.com/work/addStudentWorkNew"
	// ApiWorkCommitNew 接口-新的作业提交答案接口
	ApiWorkCommitNew = "https://mooc1.chaoxing.com/mooc-ans/work/addStudentWorkNew"

	// 接口-课程文档阅读上报
	ApiDocumentReadingReport = "https://mooc1.chaoxing.com/ananas/job/document"

	// PageMobileWork SSR页面-客户端单元测验答题页
	PageMobileWork  = "https://mooc1-api.chaoxing.com/android/mworkspecial"           // 这是个cxkitty中的
	PageMobileWorkY = "https://mooc1-api.chaoxing.com/mooc-ans/work/phone/doHomeWork" // 这个是自己爬的

	KEY = "u2oh6Vu^HWe4_AES" // 注意 Go 语言中字符串默认就是 UTF-8 编码
)
