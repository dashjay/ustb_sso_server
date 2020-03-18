package structs

type AuthCookie struct {
	UID        string `json:"uid"`
	JSESSIONID string `json:"jsessionid"`
	TS         string `json:"ts"`
	TKN        string `json:"tkn"`
}

// 返回结构定义
type AuthStruct struct {
	// Code 状态参数：如果 != 1 说明有数据传出，
	// Code 1说明有错读msg即可
	// Code 101已认证 URL 为空 UID 为学号
	// Code 102未认证 URL 为认证URL UID 为空
	Code int `json:"code" bson:"code"`

	//  "https://sis.ustb.edu.cn/auth" 开头的，传给学生，点开之后可以完成认证
	URL string `json:"url" bson:"url"`

	// 用户学号
	UID string `json:"uid" bson:"uid"`

	// 信息
	Msg string `json:"msg" bson:"msg"`
}
