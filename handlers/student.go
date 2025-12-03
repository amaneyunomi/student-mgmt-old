package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"student-management/database"
	"student-management/models"
)

// GetStudents 获取所有学生信息
func GetStudents(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT s.student_id, s.student_number, s.name, s.class_id, 
               c.class_name, c.grade_id, g.grade_name
        FROM students s
        JOIN classes c ON s.class_id = c.class_id
        JOIN grades g ON c.grade_id = g.grade_id
    `

	rows, err := database.StudentDB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.StudentWithClass
	for rows.Next() {
		var student models.StudentWithClass
		err := rows.Scan(
			&student.StudentID,
			&student.StudentNumber,
			&student.Name,
			&student.ClassID,
			&student.ClassName,
			&student.GradeID,
			&student.GradeName,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GetStudentByNumber 根据学号获取学生信息
func GetStudentByNumber(w http.ResponseWriter, r *http.Request) {
	studentNumber := r.URL.Query().Get("student_number")
	if studentNumber == "" {
		http.Error(w, "student_number parameter is required", http.StatusBadRequest)
		return
	}

	query := `
        SELECT s.student_id, s.student_number, s.name, s.class_id, 
               c.class_name, c.grade_id, g.grade_name
        FROM students s
        JOIN classes c ON s.class_id = c.class_id
        JOIN grades g ON c.grade_id = g.grade_id
        WHERE s.student_number = $1
    `

	var student models.StudentWithClass
	err := database.StudentDB.QueryRow(query, studentNumber).Scan(
		&student.StudentID,
		&student.StudentNumber,
		&student.Name,
		&student.ClassID,
		&student.ClassName,
		&student.GradeID,
		&student.GradeName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// GetStudentsByClass 根据班级获取学生信息
func GetStudentsByClass(w http.ResponseWriter, r *http.Request) {
	classIDStr := r.URL.Query().Get("class_id")
	if classIDStr == "" {
		http.Error(w, "class_id parameter is required", http.StatusBadRequest)
		return
	}

	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		http.Error(w, "Invalid class_id", http.StatusBadRequest)
		return
	}

	query := `
        SELECT s.student_id, s.student_number, s.name, s.class_id, 
               c.class_name, c.grade_id, g.grade_name
        FROM students s
        JOIN classes c ON s.class_id = c.class_id
        JOIN grades g ON c.grade_id = g.grade_id
        WHERE s.class_id = $1
    `

	rows, err := database.StudentDB.Query(query, classID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.StudentWithClass
	for rows.Next() {
		var student models.StudentWithClass
		err := rows.Scan(
			&student.StudentID,
			&student.StudentNumber,
			&student.Name,
			&student.ClassID,
			&student.ClassName,
			&student.GradeID,
			&student.GradeName,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
