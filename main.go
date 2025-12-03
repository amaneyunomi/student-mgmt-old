package StudentMgmt

import (
	"log"
	"net/http"
	"student-management/config"
	"student-management/database"
	"student-management/handlers"

	//"github.com/gorilla/mux"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	// 初始化数据库连接
	if err := database.InitStudentDB(cfg); err != nil {
		log.Fatal("Failed to connect to student database:", err)
	}

	if err := database.InitGradeDB(cfg); err != nil {
		log.Fatal("Failed to connect to grade database:", err)
	}

	// 创建路由器
	r := mux.NewRouter()

	// 学生相关路由
	r.HandleFunc("/api/students", handlers.GetStudents).Methods("GET")
	r.HandleFunc("/api/students/search", handlers.GetStudentByNumber).Methods("GET")
	r.HandleFunc("/api/students/class", handlers.GetStudentsByClass).Methods("GET")

	// 成绩相关路由
	r.HandleFunc("/api/grades", handlers.GetGradeRecords).Methods("GET")
	r.HandleFunc("/api/grades", handlers.CreateGradeRecord).Methods("POST")
	r.HandleFunc("/api/grades", handlers.UpdateGradeRecord).Methods("PUT")
	r.HandleFunc("/api/grades", handlers.DeleteGradeRecord).Methods("DELETE")

	// 基础数据路由
	r.HandleFunc("/api/subjects", handlers.GetSubjects).Methods("GET")
	r.HandleFunc("/api/exam-types", handlers.GetExamTypes).Methods("GET")

	// 启动服务器
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
