package file_handler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/horlakz/api.secretariat_repository/dto"
	"github.com/horlakz/api.secretariat_repository/handler"
	"github.com/horlakz/api.secretariat_repository/lib/constants"
	"github.com/horlakz/api.secretariat_repository/payload/request"
	"github.com/horlakz/api.secretariat_repository/payload/response"
	file_repository "github.com/horlakz/api.secretariat_repository/repository/file"
	file_service "github.com/horlakz/api.secretariat_repository/service/file"
)

type fileHandler struct {
	fileService file_service.FileServiceInterface
}

type FileHandlerInterface interface {
	CreateFile(c *fiber.Ctx) error
	GetAllFiles(c *fiber.Ctx) error
	GetFileById(c *fiber.Ctx) error
	TransferFile(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

func NewFileHandler(fileService file_service.FileServiceInterface) FileHandlerInterface {
	return &fileHandler{fileService: fileService}
}

func (f *fileHandler) GetFileExt(name string) string {
	return strings.Split(name, ".")[len(strings.Split(name, "."))-1]

}

func (f *fileHandler) CreateFile(c *fiber.Ctx) error {
	var resp response.Response
	var createFileReq request.CreateFileRequest
	var fileDto dto.FileDTO

	userId := handler.GetUserId(c)

	if err := c.BodyParser(&createFileReq); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	fileDto.Name = createFileReq.Name
	fileDto.Key = createFileReq.Key
	fileDto.MimeType = createFileReq.MimeType
	fileDto.Size = createFileReq.Size
	fileDto.UserId = userId
	fileDto.Ext = f.GetFileExt(createFileReq.Name)

	if _, err := f.fileService.CreateFile(fileDto); err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "File created successfully"

	return c.Status(http.StatusCreated).JSON(resp)
}

func (f *fileHandler) GetAllFiles(c *fiber.Ctx) error {
	var resp response.Response
	var files []dto.FileDTO
	var pageable file_repository.FilePageable

	pageable.Pageable = handler.GeneratePageable(c)
	pageable.UserId = handler.GetUserId(c)

	files, pagination, err := f.fileService.FindAllFiles(pageable)
	if err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Files retrieved successfully"
	resp.Data = map[string]interface{}{"result": files, "pagination": pagination}

	return c.Status(http.StatusOK).JSON(resp)
}

func (f *fileHandler) GetFileById(c *fiber.Ctx) error {
	var resp response.Response

	fileId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	file, err := f.fileService.FindFileById(fileId)
	if err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "File retrieved successfully"
	resp.Data = map[string]interface{}{"result": file}

	return c.Status(http.StatusOK).JSON(resp)
}

func (f *fileHandler) TransferFile(c *fiber.Ctx) error {
	var resp response.Response
	var transferReq request.TransferRequest
	var transferDto dto.TransferDTO

	userId := handler.GetUserId(c)

	if err := c.BodyParser(&transferReq); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	fileId, err := uuid.Parse(transferReq.FileId)
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	transferDto.FileId = fileId
	transferDto.FromUserId = userId

	if _, err := f.fileService.TransferFile(transferDto, transferReq.ToEmail); err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "File transferred successfully"

	return c.Status(http.StatusOK).JSON(resp)
}

func (f *fileHandler) DeleteFile(c *fiber.Ctx) error {
	var resp response.Response

	fileId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if err := f.fileService.DeleteFile(fileId); err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "File deleted successfully"

	return c.Status(http.StatusOK).JSON(resp)
}
