package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"student-management/database"
	"student-management/models"
	"time"
)

// GetGradeRecords 获取成绩记录
func GetGradeRecords(w http.ResponseWriter, r *http.Request) {
	baseQuery := `
        SELECT gr.record_id, gr.student_number, gr.subject_id, gr.exam_type_id, 
               gr.score, gr.semester, gr.created_at,
               s.subject_name, et.exam_type_name
        FROM grades_records gr
        JOIN subjects s ON gr.subject_id = s.subject_id
        JOIN exam_types et ON gr.exam_type_id = et.exam_type_id
    `

	var conditions []string
	var queryArgs []interface{}
	argIndex := 1

	// 构建动态查询条件
	if studentNumber := r.URL.Query().Get("student_number"); studentNumber != "" {
		conditions = append(conditions, fmt.Sprintf("gr.student_number = $%d", argIndex))
		queryArgs = append(queryArgs, studentNumber)
		argIndex++
	}

	if subjectIDStr := r.URL.Query().Get("subject_id"); subjectIDStr != "" {
		subjectID, err := strconv.Atoi(subjectIDStr)
		if err != nil {
			http.Error(w, "Invalid subject_id", http.StatusBadRequest)
			return
		}
		conditions = append(conditions, fmt.Sprintf("gr.subject_id = $%d", argIndex))
		queryArgs = append(queryArgs, subjectID)
		argIndex++
	}

	if semester := r.URL.Query().Get("semester"); semester != "" {
		conditions = append(conditions, fmt.Sprintf("gr.semester = $%d", argIndex))
		queryArgs = append(queryArgs, semester)
		argIndex++
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY gr.created_at DESC"

	rows, err := database.GradeDB.Query(query, queryArgs...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var records []models.GradeRecordWithDetails
	for rows.Next() {
		var record models.GradeRecordWithDetails
		err := rows.Scan(
			&record.RecordID,
			&record.StudentNumber,
			&record.SubjectID,
			&record.ExamTypeID,
			&record.Score,
			&record.Semester,
			&record.CreatedAt,
			&record.SubjectName,
			&record.ExamTypeName,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 获取学生姓名
		var studentName string
		err = database.StudentDB.QueryRow(
			"SELECT name FROM students WHERE student_number = $1",
			record.StudentNumber,
		).Scan(&studentName)
		if err != nil {
			if err == sql.ErrNoRows {
				studentName = "Unknown"
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		record.StudentName = studentName

		records = append(records, record)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// CreateGradeRecord 创建成绩记录
func CreateGradeRecord(w http.ResponseWriter, r *http.Request) {
	var record models.GradeRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if record.StudentNumber == "" || record.SubjectID == 0 || record.ExamTypeID == 0 || record.Semester == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 验证学生是否存在
	var exists bool
	err := database.StudentDB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM students WHERE student_number = $1)",
		record.StudentNumber,
	).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Student not found", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO grades_records (student_number, subject_id, exam_type_id, score, semester, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING record_id
    `

	var recordID int
	err = database.GradeDB.QueryRow(
		query,
		record.StudentNumber,
		record.SubjectID,
		record.ExamTypeID,
		record.Score,
		record.Semester,
		time.Now(),
	).Scan(&recordID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	record.RecordID = recordID
	record.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(record)
}

// UpdateGradeRecord 更新成绩记录
func UpdateGradeRecord(w http.ResponseWriter, r *http.Request) {
	recordIDStr := r.URL.Query().Get("record_id")
	if recordIDStr == "" {
		http.Error(w, "record_id parameter is required", http.StatusBadRequest)
		return
	}

	recordID, err := strconv.Atoi(recordIDStr)
	if err != nil {
		http.Error(w, "Invalid record_id", http.StatusBadRequest)
		return
	}

	var record models.GradeRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
        UPDATE grades_records 
        SET score = $1, semester = $2, subject_id = $3, exam_type_id = $4
        WHERE record_id = $5
    `

	result, err := database.GradeDB.Exec(
		query,
		record.Score,
		record.Semester,
		record.SubjectID,
		record.ExamTypeID,
		recordID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteGradeRecord 删除成绩记录
func DeleteGradeRecord(w http.ResponseWriter, r *http.Request) {
	recordIDStr := r.URL.Query().Get("record_id")
	if recordIDStr == "" {
		http.Error(w, "record_id parameter is required", http.StatusBadRequest)
		return
	}

	recordID, err := strconv.Atoi(recordIDStr)
	if err != nil {
		http.Error(w, "Invalid record_id", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM grades_records WHERE record_id = $1"
	result, err := database.GradeDB.Exec(query, recordID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSubjects 获取所有科目
func GetSubjects(w http.ResponseWriter, r *http.Request) {
	rows, err := database.GradeDB.Query("SELECT subject_id, subject_name FROM subjects ORDER BY subject_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		err := rows.Scan(&subject.SubjectID, &subject.SubjectName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		subjects = append(subjects, subject)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}

// GetExamTypes 获取所有考试类型
func GetExamTypes(w http.ResponseWriter, r *http.Request) {
	rows, err := database.GradeDB.Query("SELECT exam_type_id, exam_type_name FROM exam_types ORDER BY exam_type_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var examTypes []models.ExamType
	for rows.Next() {
		var examType models.ExamType
		err := rows.Scan(&examType.ExamTypeID, &examType.ExamTypeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		examTypes = append(examTypes, examType)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(examTypes)
}
