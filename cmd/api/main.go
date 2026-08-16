package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/thitipa-palm/7solutions-assignment/internal/config"
	"github.com/thitipa-palm/7solutions-assignment/internal/handler"
	"github.com/thitipa-palm/7solutions-assignment/internal/middleware"
	"github.com/thitipa-palm/7solutions-assignment/internal/repository"
	"github.com/thitipa-palm/7solutions-assignment/internal/router"
	"github.com/thitipa-palm/7solutions-assignment/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Printf("application stopped with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	mongoClient, err := connectMongoDB(cfg.MongoURI)
	if err != nil {
		return err
	}

	//NOTE - ปิดการเชื่อมต่อ MongoDB เมื่อ run() จบ
	defer disconnectMongoDB(mongoClient)

	//SECTION - สร้าง repository สำหรับจัดการ user collection
	userRepository := repository.NewMongoUserRepository(
		mongoClient.Database(cfg.MongoDatabase),
	)

	indexCtx, indexCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	err = userRepository.EnsureIndexes(indexCtx)
	indexCancel()

	if err != nil {
		return fmt.Errorf("ensure MongoDB indexes: %w", err)
	}

	log.Println("MongoDB indexes ensured")
	//!SECTION

	//SECTION - สร้าง service
	userService := service.NewUserService(
		userRepository,
	)
	tokenService := service.NewTokenService(
		cfg.JWTSecret,
		cfg.JWTExpiration,
	)
	authService := service.NewAuthService(
		userRepository,
		userService,
		tokenService,
	)
	//!SECTION

	validate := validator.New()

	//SECTION - สร้าง handler
	authHandler := handler.NewAuthHandler(
		authService,
		validate,
	)

	userHandler := handler.NewUserHandler(
		userService,
		validate,
	)

	//NOTE - สร้าง middleware สำหรับตรวจสอบ JWT token
	authMiddleware := middleware.Authenticate(
		tokenService,
	)
	//!SECTION

	//SECTION - สร้าง goroutine สำหรับ monitor จำนวน user
	monitorCtx, stopMonitor := context.WithCancel(
		context.Background(),
	)

	monitorDone := make(chan struct{})

	go func() {
		defer close(monitorDone)

		service.RunUserCountMonitor(
			monitorCtx,
			userService,
			10*time.Second,
		)
	}()

	defer func() {
		stopMonitor()
		<-monitorDone
	}()
	//!SECTION

	//NOTE - สร้าง Fiber app และ route สำหรับ health check
	app := fiber.New()

	//NOTE - middleware logging
	app.Use(middleware.Logging)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	//NOTE - ตั้งค่า route สำหรับ auth handler
	router.Setup(app, authHandler, userHandler, authMiddleware)

	//NOTE - รับ Ctrl+C และ SIGTERM
	signalCtx, stopSignal := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopSignal()

	//NOTE ใช้รับผลลัพธ์จาก app.Listen()
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("server is running on port %s", cfg.Port)
		serverErr <- app.Listen(":" + cfg.Port)
	}()

	select {
	case err := <-serverErr:
		//NOTE - Server เปิดไม่สำเร็จหรือหยุดเอง
		if err != nil {
			return fmt.Errorf("run HTTP server: %w", err)
		}

		return nil

	case <-signalCtx.Done():
		//NOTE - ได้รับ Ctrl+C หรือ SIGTERM
		log.Println("shutdown signal received")
	}

	//NOTE - หยุดรับ request ใหม่ และรอ request เดิมสูงสุด 5 วินาที (graceful shutdown)
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	err = app.ShutdownWithContext(shutdownCtx)
	shutdownCancel()

	if err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	//NOTE - รอให้ app.Listen() หยุดจริง
	if err := <-serverErr; err != nil {
		log.Printf("HTTP server returned during shutdown: %v", err)
	}

	log.Println("HTTP server stopped")

	return nil
}

func connectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	client, err := repository.NewMongoClient(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("initialize MongoDB: %w", err)
	}

	log.Println("MongoDB connected")

	return client, nil
}

func disconnectMongoDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("disconnect MongoDB: %v", err)
		return
	}

	log.Println("MongoDB disconnected")
}
