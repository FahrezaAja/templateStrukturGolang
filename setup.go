package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	fmt.Println("Masukkan nama proyek (tanpa spasi):")
	var projectName string
	fmt.Scanln(&projectName)
	if projectName == "" {
		log.Fatal("❌ Nama proyek tidak boleh kosong.")
	}

	fmt.Println("Masukkan path folder tujuan (kosongkan untuk lokasi sekarang):")
	var dest string
	fmt.Scanln(&dest)

	var basePath string
	if dest == "" {
		basePath = filepath.Join(".", projectName)
	} else {
		basePath = filepath.Join(dest, projectName)
	}

	dirs := []string{
		filepath.Join(basePath, "cmd"),
		filepath.Join(basePath, "app", "routes"),
		filepath.Join(basePath, "app", "controllers"),
		filepath.Join(basePath, "app", "models"),
		filepath.Join(basePath, "app", "services"),
		filepath.Join(basePath, "app", "utils"),
		filepath.Join(basePath, "app", "config"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Gagal membuat folder %s: %v", dir, err)
		}
	}

	fmt.Println("📦 Membuat go.mod ...")
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = basePath
	if err := cmd.Run(); err != nil {
		log.Fatalf("Gagal go mod init: %v", err)
	}

	fmt.Println("⚙️ Menginstal dependency ... (harap tunggu)")
	deps := []string{
		"github.com/gin-gonic/gin",
		"github.com/joho/godotenv",
		"gorm.io/gorm",
		"gorm.io/driver/mysql",
	}
	for _, dep := range deps {
		cmd = exec.Command("go", "get", "-u", dep)
		cmd.Dir = basePath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Gagal menginstal %s: %v", dep, err)
		}
	}

	mainContent := fmt.Sprintf(`package main

import (
	"%s/app/config"
	"%s/app/models"
	"%s/app/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	port := config.GetEnv("PORT", "8080")
	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":" + port)
}
`, projectName, projectName, projectName)

	if err := os.WriteFile(filepath.Join(basePath, "cmd/main.go"), []byte(mainContent), 0644); err != nil {
		log.Fatalf("Gagal membuat main.go: %v", err)
	}

	configGo := `package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Gagal load .env, gunakan default environment")
	}
}

func GetEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func ConnectDB() {
	user := GetEnv("DB_USER", "root")
	pass := GetEnv("DB_PASSWORD", "")
	host := GetEnv("DB_HOST", "127.0.0.1")
	port := GetEnv("DB_PORT", "3306")
	dbName := GetEnv("DB_NAME", "mydb")

	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, host, port)
	dbRoot, err := sql.Open("mysql", dsnRoot)
	if err != nil {
		log.Fatalf("❌ Gagal Tersabung  ke MySQL root: %v", err)
	}
	defer dbRoot.Close()

	_, err = dbRoot.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + " CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;")
	if err != nil {
		log.Fatalf("❌ Gagal membuat database: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal Tersambung ke database: %v", err)
	}
}
`
	if err := os.WriteFile(filepath.Join(basePath, "app/config/config.go"), []byte(configGo), 0644); err != nil {
		log.Fatalf("Gagal membuat config.go: %v", err)
	}

	userModel := `package models

type User struct {
	ID       uint   ` + "`json:\"id\" gorm:\"primaryKey;autoIncrement\"`" + `
	Name     string ` + "`json:\"name\"`" + `
	Email    string ` + "`json:\"email\" gorm:\"uniqueIndex\"`" + `
	Password string ` + "`json:\"password\"`" + `
}
`
	if err := os.WriteFile(filepath.Join(basePath, "app/models/user.go"), []byte(userModel), 0644); err != nil {
		log.Fatalf("Gagal membuat user.go: %v", err)
	}

	routesContent := `package routes

import (
	"github.com/gin-gonic/gin"
	"myapp/app/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", controllers.Home)
}
`
	routesContent = replaceModuleName(routesContent, projectName)
	if err := os.WriteFile(filepath.Join(basePath, "app/routes/routes.go"), []byte(routesContent), 0644); err != nil {
		log.Fatalf("Gagal membuat routes.go: %v", err)
	}

	homeController := `package controllers

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello from Gin + GORM!",
	})
}
`
	if err := os.WriteFile(filepath.Join(basePath, "app/controllers/home_controller.go"), []byte(homeController), 0644); err != nil {
		log.Fatalf("Gagal membuat home_controller.go: %v", err)
	}

	userService := `package services

import "%s/app/models"

func GetSampleUser() models.User {
	return models.User{
		ID: 1,
		Name: "Eja Tampan",
		Email: "ejatampan123@example.com",
		Password: "hashedpassword",
	}
}
`
	userService = fmt.Sprintf(userService, projectName)
	if err := os.WriteFile(filepath.Join(basePath, "app/services/user_service.go"), []byte(userService), 0644); err != nil {
		log.Fatalf("Gagal membuat user_service.go: %v", err)
	}

	hashUtil := `package utils

import "fmt"

func HashPassword(pwd string) string {
	return fmt.Sprintf("hashed(%s)", pwd)
}
`
	if err := os.WriteFile(filepath.Join(basePath, "app/utils/hash.go"), []byte(hashUtil), 0644); err != nil {
		log.Fatalf("Gagal membuat hash.go: %v", err)
	}

	envContent := fmt.Sprintf(`# Server
PORT=8080
APP_NAME=%s
ENV=development

# Database
DB_DRIVER=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=%s
`, projectName, projectName)
	if err := os.WriteFile(filepath.Join(basePath, ".env"), []byte(envContent), 0644); err != nil {
		log.Fatalf("Gagal membuat .env: %v", err)
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = basePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Gagal menjalankan go mod tidy: %v", err)
	}

	fmt.Println()
	fmt.Println("✅ Proyek berhasil dibuat dengan GORM + Gin + godotenv + auto-create database!")
	fmt.Println("----------------------------------")
	fmt.Printf("📁 Lokasi: %s\n", basePath)
	fmt.Println("➡ Jalankan perintah berikut:")
	fmt.Printf("   cd \"%s\"\n", basePath)
	fmt.Println("   go run cmd/main.go")
	fmt.Println("----------------------------------")
}

func replaceModuleName(content, moduleName string) string {
	
	content = strings.ReplaceAll(content, "myapp", moduleName)
	content = strings.ReplaceAll(content, "myApp", moduleName)
	return content
}
