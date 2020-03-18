package structs

import "fmt"

type Reply struct {
	Msg  string  `json:"msg" bson:"msg"`
	Body []Grade `json:"body" bson:"body"`
}

type Grade struct {
	SchoolYear    int     `json:"schoolYear" bson:"schoolYear"`
	Semester      int     `json:"semester" bson:"semester"`
	CourseNo      string  `json:"courseNo" bson:"courseNo"`
	CourseName    string  `json:"courseName" bson:"courseName"`
	Period        int     `json:"period" bson:"period"`
	Credit        float32 `json:"credit" bson:"credit"`
	ExamResults   float32 `json:"examResults" bson:"examResults"`
	CourseResults string  `json:"courseResults" bson:"courseResults"`
	OriginResults float32 `json:"originResults" bson:"originResults"`
	Username      string  `json:"username" bson:"username"`
}

func (g *Grade) Print() string {
	return fmt.Sprintf("课程名称=%s,分数=%.1f,学分=%.1f", g.CourseName, g.ExamResults, g.Credit)
}
