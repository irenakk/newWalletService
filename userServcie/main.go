package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	config2 "newWalletService/config"
	"newWalletService/internal/handler"
	"newWalletService/internal/middleware"
	"newWalletService/internal/repository"
	"newWalletService/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config2.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := config2.NewDatabase(cfg.GetDSN())
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
	walletRepository := repository.NewWalletRepository(db)
	accountRepository := repository.NewAccountRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository, walletRepository, accountRepository)
	walletUsecase := usecase.NewWalletUsecase(userRepository, walletRepository, accountRepository)

	// Initialize handlers with JWT configuration
	authHandler := handler.NewAuthHandler([]byte(cfg.JWT.Secret), userUsecase)
	walletHandler := handler.NewWalletHandler(walletUsecase)

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
		protected.POST("/add", walletHandler.Add)
		protected.POST("/transfer", walletHandler.Transfer)
	}

	// Start server with configured host and port
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}
