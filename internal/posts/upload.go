package posts

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"starehian-society-platform/internal/middleware"
	"starehian-society-platform/pkg/storage"
)

type UploadHandler struct {
	storage *storage.R2Storage
}

func NewUploadHandler(storage *storage.R2Storage) *UploadHandler {
	return &UploadHandler{
		storage: storage,
	}
}

// UploadMedia handles media file uploads
func (h *UploadHandler) UploadMedia(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File too large (max 10MB)",
		})
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type (only images allowed)",
		})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer src.Close()

	// Generate unique key
	fileID := uuid.New().String()
	key := h.storage.GenerateKey(userID, fileID+ext)

	// Determine content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// Upload to R2
	url, err := h.storage.UploadFile(c.Context(), key, src, contentType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	return c.JSON(fiber.Map{
		"url": url,
		"key": key,
	})
}

// UploadMultipleMedia handles multiple file uploads
func (h *UploadHandler) UploadMultipleMedia(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	// Limit to 10 files
	if len(files) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Too many files (max 10)",
		})
	}

	var uploadedFiles []fiber.Map
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	for _, fileHeader := range files {
		// Validate file size
		if fileHeader.Size > 10*1024*1024 {
			continue
		}

		// Validate file type
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedExts[ext] {
			continue
		}

		// Open file
		src, err := fileHeader.Open()
		if err != nil {
			continue
		}

		// Generate unique key
		fileID := uuid.New().String()
		key := h.storage.GenerateKey(userID, fileID+ext)

		// Determine content type
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}

		// Upload to R2
		url, err := h.storage.UploadFile(c.Context(), key, src, contentType)
		src.Close()

		if err != nil {
			continue
		}

		uploadedFiles = append(uploadedFiles, fiber.Map{
			"url": url,
			"key": key,
		})
	}

	return c.JSON(fiber.Map{
		"files": uploadedFiles,
		"count": len(uploadedFiles),
	})
}

// ProcessImage processes an image (resize, compress)
// This is a placeholder - in production, you'd use an image processing library
func (h *UploadHandler) ProcessImage(file multipart.File) (io.Reader, error) {
	// For now, just return the original file
	// In production, use imaging library to resize/compress
	return file, nil
}
