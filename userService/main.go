package main

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"newWalletService/config"
	"newWalletService/internal/handler"
	"newWalletService/internal/middleware"
	"newWalletService/internal/repository"
	"newWalletService/internal/rpctransfer"
	"newWalletService/internal/usecase"
	"newWalletService/proto"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := config.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.DB.Close()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router with middleware
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	userRepository := repository.NewUserRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository)

	authHandler := handler.NewAuthHandler([]byte(cfg.JWT.Secret), userUsecase)

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes with JWT middleware
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware([]byte(cfg.JWT.Secret)))
	{
		protected.GET("/hello", authHandler.Hello)
	}

	// ---- HTTP server ----
	httpAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", httpAddr)
	srv := &http.Server{
		Addr:         httpAddr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("HTTP server starting on %s", httpAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// ---- gRPC server ----
	grpcAddr := ":9090"
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", grpcAddr, err)
	}

	grpcServer := grpc.NewServer()
	h := &rpctransfer.Handlers{
		Usecase: usecase.NewUserUsecase(userRepository),
	}
	proto.RegisterUserServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	log.Printf("gRPC server starting on %s", grpcAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
