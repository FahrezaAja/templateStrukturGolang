package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 1️⃣ Tanya nama proyek
	fmt.Println("Masukkan nama proyek:")
	var projectName string
	fmt.Scanln(&projectName)
	if projectName == "" {
		log.Fatal("❌ Nama proyek tidak boleh kosong.")
	}

	// 2️⃣ Tanya lokasi tujuan
	fmt.Println("Masukkan path folder tujuan (kosongkan untuk lokasi sekarang):")
	var dest string
	fmt.Scanln(&dest)

	var basePath string
	if dest == "" {
		basePath = filepath.Join(".", projectName)
	} else {
		basePath = filepath.Join(dest, projectName)
	}

	// 3️⃣ Buat folder utama dan subfolder
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
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Gagal membuat folder %s: %v", dir, err)
		}
	}

	// 4️⃣ Inisialisasi go.mod
	fmt.Println("📦 Membuat file go.mod ...")
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = basePath
	if err := cmd.Run(); err != nil {
		log.Fatalf("Gagal menjalankan go mod init: %v", err)
	}

	// 5️⃣ Install Gin
	fmt.Println("⚙️ Menginstal framework Gin ... (harap tunggu)")
	cmd = exec.Command("go", "get", "-u", "github.com/gin-gonic/gin")
	cmd.Dir = basePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Gagal menginstal Gin: %v", err)
	}

	// 6️⃣ Buat file main.go
	mainContent := fmt.Sprintf(`package main

import (
	"%s/app/config"
	"%s/app/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	port := config.GetEnv("PORT", "8080")

	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":" + port)
}
`, projectName, projectName)

	err := os.WriteFile(filepath.Join(basePath, "cmd", "main.go"), []byte(mainContent), 0644)
	if err != nil {
		log.Fatalf("Gagal membuat cmd/main.go: %v", err)
	}

	// 7️⃣ Buat routes.go
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
	err = os.WriteFile(filepath.Join(basePath, "app", "routes", "routes.go"), []byte(routesContent), 0644)
	if err != nil {
		log.Fatalf("Gagal membuat routes.go: %v", err)
	}

	// 8️⃣ Buat home_controller.go
	homeController := `package controllers

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello from Gin!",
	})
}
`
	err = os.WriteFile(filepath.Join(basePath, "app", "controllers", "home_controller.go"), []byte(homeController), 0644)
	if err != nil {
		log.Fatalf("Gagal membuat controllers/home_controller.go: %v", err)
	}

	// 9️⃣ Buat user model
	userModel := `package models

type User struct {
	ID       uint   ` + "`json:\"id\"`" + `
	Name     string ` + "`json:\"name\"`" + `
	Email    string ` + "`json:\"email\"`" + `
	Password string ` + "`json:\"password\"`" + `
}
`
	err = os.WriteFile(filepath.Join(basePath, "app", "models", "user.go"), []byte(userModel), 0644)
	if err != nil {
		log.Fatalf("Gagal membuat models/user.go: %v", err)
	}

	// 10️⃣ Buat sample service
	userService := `package services

import "myapp/app/models"

func GetSampleUser() models.User {
	return models.User{
		ID: 1,
		Name: "John Doe",
		Email: "john@example.com",
		Password: "hashedpassword",
	}
}
`
	userService = replaceModuleName(userService, projectName)
	err = os.WriteFile(filepath.Join(basePath, "app", "services", "user_service.go"), []byte(userService), 0644)
	if err != nil {
		log.Fatalf("Gagal membuat services/user_service.go: %v", err)
	}

	// 11️⃣ Buat utils/hash.go
	hashUtil := `package utils

import "fmt"

func HashPassword(pwd string) string {
	return fmt.Sprintf("hashed(%s)", pwd)
}
`
	err = os.WriteFile(filepath.Join(basePath, "app", "utils", "hash.go"), []byte(hashUtil), 0644)
	if err != nil {
		log.Fatalf("Gagal membuat utils/hash.go: %v", err)
	}

	// 12️⃣ Buat config/config.go
	configGo := `package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Gagal load .env, gunakan default environment")
	}
}

func GetEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
`
	err = os.WriteFile(filepath.Join(basePath, "app", "config", "config.go"), []byte(configGo), 0644)
	if err != nil {
		log.Fatalf("Gagal membuat config/config.go: %v", err)
	}

	// 13️⃣ Buat .env
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

# JWT
JWT_SECRET=mysecretkey
JWT_EXPIRE=24h

# Logging
LOG_LEVEL=debug
`, projectName, projectName)
	err = os.WriteFile(filepath.Join(basePath, ".env"), []byte(envContent), 0644)
	if err != nil {
		log.Fatalf("Gagal membuat .env: %v", err)
	}

	fmt.Println()
	fmt.Println("✅ Proyek berhasil dibuat!")
	fmt.Println("----------------------------------")
	fmt.Printf("📁 Lokasi: %s\n", basePath)
	fmt.Println("➡ Jalankan perintah berikut:")
	fmt.Printf("   cd \"%s\"\n", basePath)
	fmt.Println("   go run cmd/main.go")
	fmt.Println("----------------------------------")
}

// replaceModuleName replaces placeholder 'myapp' with actual module name
func replaceModuleName(content, moduleName string) string {
	return string([]byte(fmt.Sprintf(content, moduleName)))
}
