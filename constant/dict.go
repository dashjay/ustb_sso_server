package constant

const (
	GetRandTokenURL = "http://jwstu.ustb.edu.cn/stu_comm/login.html"
	SISGetState     = "https://sis.ustb.edu.cn/connect/state"
	SISAuth         = "https://sis.ustb.edu.cn/auth"
	SISGetQr        = "https://sis.ustb.edu.cn/connect/qrpage"
	JWSTU           = "http://jwstu.ustb.edu.cn/smvc/WeiXinLoginService/sisWebLoginReturn.json?redirectDomain=jwstu.ustb.edu.cn&appid=9673f8e6b2c74d48b0a871152ce1d383&auth_code=%s&rand_token=%s"

	GetGradeURL          = "http://jwstu.ustb.edu.cn/smvc/StuQueryInfoService/listStuScore.json"              // 获取成绩
	GetHeadImageURL      = "http://jwstu.ustb.edu.cn/smvc/PullFilesService/imgByOwn.json?imgType=stu_img"     // 获取个人图片
	GetStuInfo           = "http://jwstu.ustb.edu.cn/smvc/CommonService/obtainCurUserStu.json"                // 获取个人信息
	GetCourseTable       = "http://jwstu.ustb.edu.cn/smvc/StuQueryInfoService/viewStuCourseSchedule.json"     // 课程表
	GetExamArrangement   = "http://jwstu.ustb.edu.cn/smvc/StuQueryInfoService/viewExamArrangement.json"       // 考试安排
	GetCETScore          = "http://jwstu.ustb.edu.cn/smvc/QueryScoreService/queryScore.json"                  // 获取四六级分数
	GetStuInfoPerfection = "http://jwstu.ustb.edu.cn/smvc/StuInfoChangeService/obtain1StuInfoPerfection.json" // 获取用户采集的信息

	// {"startSection":"1","endSection":"2","queryDate":"2020-02-24"} POST 传输
	GetEmptyClassRoom = "http://jwstu.ustb.edu.cn/smvc/StuQueryInfoService/listFreeClassroom.json"
)
