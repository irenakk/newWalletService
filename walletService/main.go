package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"newWalletService/config"
	"newWalletService/internal/handler"
	"newWalletService/internal/middleware"
	"newWalletService/internal/repository"
	"newWalletService/internal/usecase"
	pb "newWalletService/proto/wallet"
)

type walletServer struct {
	pb.UnimplementedWalletServiceServer
	accounts map[string]*pb.Account
}

// Implement the GetMessage RPC
func (s *walletServer) Add(c *gin.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	message := fmt.Sprintf("Hello, %s!", req.Name)
	return &pb.MessageResponse{Message: message}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось слушать порт: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userServer{})

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
	walletRepository := repository.NewWalletRepository(db)
	accountRepository := repository.NewAccountRepository(db)

	walletUsecase := usecase.NewWalletUsecase(userRepository, walletRepository, accountRepository)

	// Initialize handlers with JWT configuration
	walletHandler := handler.NewWalletHandler(walletUsecase)

	// Protected routes with JWT middleware
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware([]byte(cfg.JWT.Secret)))
	{
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
