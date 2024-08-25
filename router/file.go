package router

import (
	"github.com/gofiber/fiber/v2"

	file_handler "github.com/horlakz/api.secretariat_repository/handler/file"
	"github.com/horlakz/api.secretariat_repository/lib/constants"
	"github.com/horlakz/api.secretariat_repository/lib/database"
	"github.com/horlakz/api.secretariat_repository/middleware"
	file_repository "github.com/horlakz/api.secretariat_repository/repository/file"
	user_repository "github.com/horlakz/api.secretariat_repository/repository/user"
	file_service "github.com/horlakz/api.secretariat_repository/service/file"
)

func InitializeFileRouter(router fiber.Router, db database.DatabaseInterface, env constants.Env) {
	userRepository := user_repository.NewUserRepository(db)
	transferRepository := file_repository.NewTransferRepository(db)
	fileRepository := file_repository.NewFileRepository(db)

	fileService := file_service.NewFileService(fileRepository, transferRepository, userRepository)

	authMiddleware := middleware.Protected()

	fileHandler := file_handler.NewFileHandler(fileService)

	fileRouter := router.Group("file", authMiddleware)

	fileRouter.Post("/", fileHandler.CreateFile)
	fileRouter.Get("/", fileHandler.GetAllFiles)
	fileRouter.Get("/:id", fileHandler.GetFileById)
	fileRouter.Delete("/:id", fileHandler.DeleteFile)
	fileRouter.Post("/share", fileHandler.TransferFile)
}
