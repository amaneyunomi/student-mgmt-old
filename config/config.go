// config/config.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	StudentDBHost     string
	StudentDBPort     string
	StudentDBUser     string
	StudentDBPassword string
	StudentDBName     string

	GradeDBHost     string
	GradeDBPort     string
	GradeDBUser     string
	GradeDBPassword string
	GradeDBName     string
}

func LoadConfig() *Config {
	// 仅在开发环境加载 .env（生产环境用系统环境变量）
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: No .env file found")
		}
	}

	return &Config{
		StudentDBHost:     getEnv("STUDENT_DB_HOST", "localhost"),
		StudentDBPort:     getEnv("STUDENT_DB_PORT", "5432"),
		StudentDBUser:     getEnv("STUDENT_DB_USER", "postgres"),
		StudentDBPassword: getEnv("STUDENT_DB_PASSWORD", ""),
		StudentDBName:     getEnv("STUDENT_DB_NAME", "student_info_db"),

		GradeDBHost:     getEnv("GRADE_DB_HOST", "localhost"),
		GradeDBPort:     getEnv("GRADE_DB_PORT", "5432"),
		GradeDBUser:     getEnv("GRADE_DB_USER", "postgres"),
		GradeDBPassword: getEnv("GRADE_DB_PASSWORD", ""),
		GradeDBName:     getEnv("GRADE_DB_NAME", "grade_db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetStudentDBURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.StudentDBHost, c.StudentDBPort, c.StudentDBUser, c.StudentDBPassword, c.StudentDBName)
}

func (c *Config) GetGradeDBURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.GradeDBHost, c.GradeDBPort, c.GradeDBUser, c.GradeDBPassword, c.GradeDBName)
}
