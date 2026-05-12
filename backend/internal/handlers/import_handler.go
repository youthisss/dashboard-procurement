package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rygell-dashboard/internal/config"
	"rygell-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// ImportHandler handles Excel file upload and parsing.
type ImportHandler struct {
	parserService *services.ParserService
	importService *services.ImportService
	cfg           *config.Config
}

// NewImportHandler creates a new ImportHandler.
func NewImportHandler(parserService *services.ParserService, importService *services.ImportService, cfg *config.Config) *ImportHandler {
	return &ImportHandler{
		parserService: parserService,
		importService: importService,
		cfg:           cfg,
	}
}

// UploadAndParse handles Excel file upload, saves the file, and returns parsed preview data.
func (h *ImportHandler) UploadAndParse(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	safeFileName, err := sanitizeUploadedExcelFilename(file.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate file size (max 50MB)
	if file.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 50MB limit"})
		return
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(h.cfg.UploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	// Save file with timestamp
	timestamp := time.Now().Format("20060102_150405")
	savedName := fmt.Sprintf("%s_%s", timestamp, safeFileName)
	savedPath := filepath.Join(h.cfg.UploadDir, savedName)

	if err := c.SaveUploadedFile(file, savedPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Parse the uploaded file
	sheets, err := h.parserService.ParseExcelFile(savedPath)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "failed to parse excel file",
			"details": err.Error(),
		})
		return
	}

	// Validate rows and collect errors
	var warnings []string
	for _, sheet := range sheets {
		for i, row := range sheet.Rows {
			rowErrors := h.parserService.ValidateRow(row, sheet.SheetType)
			for _, e := range rowErrors {
				warnings = append(warnings, fmt.Sprintf("Sheet '%s', row %d: %s", sheet.SheetName, i+2, e))
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"file_name": file.Filename,
		"saved_as":  savedName,
		"sheets":    sheets,
		"warnings":  warnings,
	})
}

func sanitizeUploadedExcelFilename(name string) (string, error) {
	cleaned := strings.TrimSpace(name)
	if cleaned == "" || cleaned == "." {
		return "", fmt.Errorf("invalid file name")
	}
	if strings.ContainsAny(cleaned, `/\`) || filepath.IsAbs(cleaned) || filepath.Base(cleaned) != cleaned {
		return "", fmt.Errorf("invalid file name")
	}
	ext := strings.ToLower(filepath.Ext(cleaned))
	if ext != ".xlsx" && ext != ".xls" {
		return "", fmt.Errorf("only .xlsx and .xls files are allowed")
	}
	return cleaned, nil
}

// ConfirmImport receives a previously saved filename and bulk-inserts the data into the database.
func (h *ImportHandler) ConfirmImport(c *gin.Context) {
	var req struct {
		SavedAs      string                 `json:"saved_as" binding:"required"`
		ContractMode string                 `json:"contract_mode"`
		Sheets       []services.ParsedSheet `json:"sheets"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "saved_as field is required"})
		return
	}
	contractMode := strings.TrimSpace(req.ContractMode)
	if contractMode == "" {
		contractMode = "plan"
	}
	if contractMode != "plan" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "actual contract import is not available yet"})
		return
	}

	// Validate filename and resolve inside upload directory only
	cleanName := filepath.Clean(strings.TrimSpace(req.SavedAs))
	if cleanName == "" || cleanName == "." || filepath.IsAbs(cleanName) || cleanName != filepath.Base(cleanName) || strings.Contains(cleanName, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid saved_as value"})
		return
	}

	uploadDirAbs, err := filepath.Abs(h.cfg.UploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve upload directory"})
		return
	}

	filePathAbs, err := filepath.Abs(filepath.Join(uploadDirAbs, cleanName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve file path"})
		return
	}

	uploadPrefix := uploadDirAbs + string(filepath.Separator)
	if !strings.HasPrefix(filePathAbs, uploadPrefix) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid saved_as value"})
		return
	}

	hasParsedSheets := len(req.Sheets) > 0
	if !hasParsedSheets {
		if _, err := os.Stat(filePathAbs); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				c.JSON(http.StatusNotFound, gin.H{"error": "file not found on server"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to access uploaded file"})
			return
		}
	}

	var result *services.ImportResult
	if hasParsedSheets {
		result, err = h.importService.ConfirmParsedSheets(req.Sheets)
	} else {
		result, err = h.importService.ConfirmImport(filePathAbs)
	}
	if err != nil {
		if services.IsImportValidationError(err) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":  err.Error(),
				"result": result,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clean up the uploaded file after successful import (unless IMPORT_KEEP_FILES is true)
	if !h.cfg.ImportKeepFiles {
		if removeErr := os.Remove(filePathAbs); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			log.Printf("[IMPORT] Warning: failed to remove file after import: %v", removeErr)
		} else {
			log.Printf("[IMPORT] Cleaned up imported file: %s", cleanName)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Import confirmed and data saved to database.",
		"result":  result,
	})
}
