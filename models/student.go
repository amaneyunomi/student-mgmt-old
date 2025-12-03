package models

type Grade struct {
	GradeID   int    `json:"grade_id"`
	GradeName string `json:"grade_name"`
}

type Class struct {
	ClassID   int    `json:"class_id"`
	GradeID   int    `json:"grade_id"`
	ClassName string `json:"class_name"`
}

type Student struct {
	StudentID     int    `json:"student_id"`
	StudentNumber string `json:"student_number"`
	Name          string `json:"name"`
	ClassID       int    `json:"class_id"`
}

type StudentWithClass struct {
	StudentID     int    `json:"student_id"`
	StudentNumber string `json:"student_number"`
	Name          string `json:"name"`
	ClassID       int    `json:"class_id"`
	ClassName     string `json:"class_name"`
	GradeID       int    `json:"grade_id"`
	GradeName     string `json:"grade_name"`
}
