package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaiWittawat/openclaw-empire/config"
	"github.com/MaiWittawat/openclaw-empire/internal/handler"
	"github.com/MaiWittawat/openclaw-empire/internal/middleware"
	"github.com/MaiWittawat/openclaw-empire/internal/repository"
	"github.com/MaiWittawat/openclaw-empire/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cnf, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err) // เปลี่ยนจาก panic เป็น Fatal ที่มี context
	}

	// --- 1. Database Security & Optimization ---
	db, err := gorm.Open(postgres.Open(cnf.DB.GetConnStr()), &gorm.Config{
		PrepareStmt: true,                                  // ป้องกัน SQL Injection และช่วยเรื่อง Performance
		Logger:      logger.Default.LogMode(logger.Silent), // ปิด Log ที่อาจหลุดข้อมูล Sensitive ใน Production
	})
	if err != nil {
		logrus.Fatalf("Unable to connect to database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// --- 2. Middleware & Security Setup ---
	r := gin.New()                       // ใช้ New แทน Default เพื่อคุม Middleware เอง
	r.Use(gin.Recovery())                // ป้องกัน Server Crash จาก Panic
	r.Use(middleware.LogrusMiddleware()) // ใช้ Logrus ร่วมกับ Gin

	// เพิ่ม CORS Policy (ตัวอย่าง)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- 3. Routing ---
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	api := r.Group("/api")
	{
		api.GET("/agents", h.GetAgents)
		api.GET("/tasks", h.GetTasks)
		api.POST("/tasks", h.CreateTask)
		api.GET("/stats", h.GetStats)
	}
	r.GET("/ws", h.HandleWebSocket)
	r.Static("/public", "./public") // จำกัด Path ให้ชัดเจน

	// --- 4. Graceful Shutdown (สำคัญมาก) ---
	srv := &http.Server{
		Addr:    ":" + cnf.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Listen: %s\n", err)
		}
	}()

	// รอรับ Signal เพื่อปิดระบบ
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}
}
