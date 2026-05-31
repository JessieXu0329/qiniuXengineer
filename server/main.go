package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aipr/ai-pr-reviewer/controller"
	"github.com/aipr/ai-pr-reviewer/model"
	"github.com/aipr/ai-pr-reviewer/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Global database and cache variables
var (
	db  *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
)

// CORS Middleware to enable communication with the cyber-theme frontend
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func main() {
	// Initialize Gin engine
	r := gin.Default()
	r.Use(CORSMiddleware())

	// Print Cyber Banner
	fmt.Println(`
    ===================================================================
      ▲   I    P   R      R   E   V   I   E   W      A   S   S   I   S   T
      🤖 AI PR Review Assistant (Gin Backend Core) - Cyber Tech v2
    ===================================================================
    `)

	// Initialize MySQL & Redis
	initInfrastructure()

	// Instantiate multi-layered service and controller modules
	reviewService := service.NewReviewService(db, rdb)
	reviewController := controller.NewReviewController(reviewService)

	// Base API V1 Group
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "online",
				"timestamp": time.Now().Format(time.RFC3339),
				"engine":    "Gin Web Framework Core v2",
			})
		})

		// Analysis triggering and operations bindings
		v1.POST("/review/analyze", reviewController.Analyze)
		v1.GET("/review/:id", reviewController.GetDetail)
		v1.GET("/review/dashboard", reviewController.GetDashboard)
	}

	log.Println("[INFO] Cyber-engine Backend listening on port :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("[FATAL] Server startup failed: %v", err)
	}
}

func initInfrastructure() {
	// Connect to local native MySQL at 127.0.0.1:3306
	dsn := "aipr_user:aipr_password@tcp(127.0.0.1:3306)/aipr_db?charset=utf8mb4&parseTime=True&loc=Local"
	log.Println("[CONNECT] Connecting to MySQL database...")
	
	// Create context with short timeout
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	
	var err error
	done := make(chan bool)
	go func() {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// Configure connection pool to prevent connection drop errors
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				sqlDB.SetMaxIdleConns(10)
				sqlDB.SetMaxOpenConns(100)
				sqlDB.SetConnMaxLifetime(time.Hour)
			}
			// Auto migrate GORM database schemas dynamically
			log.Println("[MIGRATE] Bootstrapping GORM database schemas auto-migrations...")
			migrationErr := db.AutoMigrate(&model.PRReviewTask{}, &model.ReviewComment{})
			if migrationErr != nil {
				log.Printf("[WARNING] Schema migration failed: %v", migrationErr)
			} else {
				log.Println("[MIGRATE] GORM tables successfully migrated.")
			}
		}
		done <- true
	}()

	select {
	case <-done:
		if err != nil {
			log.Printf("[WARNING] MySQL not reachable: %v. Review persistence and dashboard metrics are disabled.", err)
		} else {
			log.Println("[SUCCESS] MySQL connected successfully.")
		}
	case <-dbCtx.Done():
		log.Println("[WARNING] MySQL connection timeout. Review persistence and dashboard metrics are disabled.")
	}

	// Connect to local Redis at 127.0.0.1:6379
	log.Println("[CONNECT] Connecting to Redis...")
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // No password by default
		DB:       0,  // default DB
	})

	redisStatus := rdb.Ping(ctx)
	if redisStatus.Err() != nil {
		log.Printf("[WARNING] Redis not reachable: %v. Running without caching.", redisStatus.Err())
	} else {
		log.Println("[SUCCESS] Redis connected successfully.")
	}
}
