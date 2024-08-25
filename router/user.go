package router

import (
	"github.com/gofiber/fiber/v2"

	userHandler "github.com/horlakz/api.secretariat_repository/handler/user"
	"github.com/horlakz/api.secretariat_repository/lib/config"
	"github.com/horlakz/api.secretariat_repository/lib/constants"
	"github.com/horlakz/api.secretariat_repository/lib/database"
	"github.com/horlakz/api.secretariat_repository/middleware"
	user_repository "github.com/horlakz/api.secretariat_repository/repository/user"
	"github.com/horlakz/api.secretariat_repository/service"
	user_service "github.com/horlakz/api.secretariat_repository/service/user"
)

func InitializeUserRouter(router fiber.Router, db database.DatabaseInterface, env constants.Env) {
	// Repositories
	userRepository := user_repository.NewUserRepository(db)
	verificationCodeRepository := user_repository.NewVerificationCodeRepository(db)

	// config
	mailConfig := config.NewEmail(env)

	// Services
	emailService := service.NewEmailService(mailConfig, db.Cache())
	userService := user_service.NewUserService(userRepository)
	verificationCodeService := user_service.NewVerficationCodeService(userRepository, verificationCodeRepository)
	authService := user_service.NewAuthService(userService, verificationCodeService, emailService)

	// Handler
	authHandler := userHandler.NewAuthHandler(authService)
	baseUserHandler := userHandler.NewUserHandler(userService)

	// Middlewares
	authMiddleware := middleware.Protected()
	roleMiddleware := middleware.NewRoleMiddleware(userRepository)

	// Routers
	authRoute := router.Group("/auth")
	userRoute := router.Group("/user", authMiddleware)

	// Routes
	authRoute.Post("/check-email", authHandler.CheckEmail)
	authRoute.Post("/login", authHandler.Login)
	authRoute.Post("/register", authHandler.Register)
	authRoute.Post("/refresh-token", authHandler.RefreshAccessToken)
	authRoute.Post("/resend-email", authHandler.ResendEmailVerification)
	authRoute.Post("/verify-email", authHandler.VerifyEmail)
	authRoute.Post("/verify-email-code", authHandler.VerifyEmailAndCode)
	authRoute.Post("/forgot-password", authHandler.ForgotPassword)
	authRoute.Post("/reset-password", authHandler.ResetPassword)

	userRoute.Get("/", baseUserHandler.UserDetails)
	userRoute.Get("/all", roleMiddleware.ValidateRole(user_service.UserRoleAdmin), baseUserHandler.FindAllUsers)

}
