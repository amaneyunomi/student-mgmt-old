package models

import "time"

type Subject struct {
	SubjectID   int    `json:"subject_id"`
	SubjectName string `json:"subject_name"`
}

type ExamType struct {
	ExamTypeID   int    `json:"exam_type_id"`
	ExamTypeName string `json:"exam_type_name"`
}

type GradeRecord struct {
	RecordID      int       `json:"record_id"`
	StudentNumber string    `json:"student_number"`
	SubjectID     int       `json:"subject_id"`
	ExamTypeID    int       `json:"exam_type_id"`
	Score         float64   `json:"score"`
	Semester      string    `json:"semester"`
	CreatedAt     time.Time `json:"created_at"`
}

type GradeRecordWithDetails struct {
	RecordID      int       `json:"record_id"`
	StudentNumber string    `json:"student_number"`
	StudentName   string    `json:"student_name"`
	SubjectName   string    `json:"subject_name"`
	SubjectID     int       `json:"subject_id"`
	ExamTypeName  string    `json:"exam_type_name"`
	ExamTypeID    int       `json:"exam_type_id"`
	Score         float64   `json:"score"`
	Semester      string    `json:"semester"`
	CreatedAt     time.Time `json:"created_at"`
}
