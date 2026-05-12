package router

import (
	"strings"

	"rygell-dashboard/internal/config"
	"rygell-dashboard/internal/handlers"
	"rygell-dashboard/internal/middleware"
	"rygell-dashboard/internal/repositories"
	"rygell-dashboard/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup initializes all dependencies and registers routes.
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	if cfg != nil && cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}
	r := gin.Default()
	if cfg != nil {
		if cfg.TrustedProxies == "" {
			_ = r.SetTrustedProxies(nil)
		} else {
			_ = r.SetTrustedProxies(splitCSV(cfg.TrustedProxies))
		}
	}

	// Middleware
	r.Use(middleware.CORSMiddleware(cfg))
	r.Use(middleware.RequestBodyLimit(defaultMaxBodyBytes(cfg), "/api/v1/import/excel", "/api/v1/import/confirm"))

	// --- Initialize layers ---

	// Repositories
	masterRepo := repositories.NewMasterRepository(db)
	contractRepo := repositories.NewContractRepository(db)
	auditRepo := repositories.NewAuditRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Services
	masterService := services.NewMasterService(masterRepo)
	contractService := services.NewContractService(contractRepo, auditRepo, masterRepo)
	parserService := services.NewParserService()
	exportService := services.NewExportService(contractRepo, masterRepo)
	userService := services.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)

	importService := services.NewImportService(parserService, masterRepo, contractRepo)

	// Handlers
	masterHandler := handlers.NewMasterHandler(masterService)
	contractHandler := handlers.NewContractHandler(contractService)
	searchHandler := handlers.NewSearchHandler(masterService, contractService)
	importHandler := handlers.NewImportHandler(parserService, importService, cfg)
	exportHandler := handlers.NewExportHandler(exportService)
	authHandler := handlers.NewAuthHandler(userService, cfg)

	// --- Register Routes ---
	api := r.Group("/api/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/logout", authHandler.Logout)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret, cfg.CookieName))
	protected.Use(middleware.StrictOriginForUnsafeMethods(cfg))
	protected.GET("/auth/me", authHandler.Me)
	protected.GET("/auth/users", authHandler.ListUsers)
	protected.POST("/auth/users", middleware.AdminOnly(cfg), authHandler.CreateUser)
	protected.DELETE("/auth/users/:id", middleware.AdminOnly(cfg), authHandler.DeleteUser)
	protected.PATCH("/auth/users/:id/password", middleware.AdminOnly(cfg), authHandler.UpdateUserPassword)
	protected.GET("/search", searchHandler.Search)

	// Master Data
	mills := protected.Group("/mills")
	{
		mills.GET("", masterHandler.GetAllMills)
		mills.GET("/:id", masterHandler.GetMillByID)
		mills.POST("", middleware.AdminOnly(cfg), masterHandler.CreateMill)
		mills.PUT("/:id", middleware.AdminOnly(cfg), masterHandler.UpdateMill)
		mills.DELETE("/:id", middleware.AdminOnly(cfg), masterHandler.DeleteMill)
	}

	vendors := protected.Group("/vendors")
	{
		vendors.GET("", masterHandler.GetAllVendors)
		vendors.GET("/:id", masterHandler.GetVendorByID)
		vendors.POST("", middleware.AdminOnly(cfg), masterHandler.CreateVendor)
		vendors.PUT("/:id", middleware.AdminOnly(cfg), masterHandler.UpdateVendor)
		vendors.DELETE("/:id", middleware.AdminOnly(cfg), masterHandler.DeleteVendor)
	}

	products := protected.Group("/products")
	{
		products.GET("", masterHandler.GetAllProducts)
		products.GET("/:id", masterHandler.GetProductByID)
		products.POST("", middleware.AdminOnly(cfg), masterHandler.CreateProduct)
		products.PUT("/:id", middleware.AdminOnly(cfg), masterHandler.UpdateProduct)
		products.DELETE("/:id", middleware.AdminOnly(cfg), masterHandler.DeleteProduct)
	}

	zones := protected.Group("/zones")
	{
		zones.GET("", masterHandler.GetAllZones)
		zones.GET("/:id", masterHandler.GetZoneByID)
		zones.POST("", middleware.AdminOnly(cfg), masterHandler.CreateZone)
		zones.PUT("/:id", middleware.AdminOnly(cfg), masterHandler.UpdateZone)
		zones.DELETE("/:id", middleware.AdminOnly(cfg), masterHandler.DeleteZone)
	}

	mots := protected.Group("/mots")
	{
		mots.GET("", masterHandler.GetAllMots)
		mots.GET("/:id", masterHandler.GetMotByID)
		mots.POST("", middleware.AdminOnly(cfg), masterHandler.CreateMot)
		mots.PUT("/:id", middleware.AdminOnly(cfg), masterHandler.UpdateMot)
		mots.DELETE("/:id", middleware.AdminOnly(cfg), masterHandler.DeleteMot)
	}

	uoms := protected.Group("/uoms")
	{
		uoms.GET("", masterHandler.GetAllUoms)
		uoms.GET("/:id", masterHandler.GetUomByID)
		uoms.POST("", middleware.AdminOnly(cfg), masterHandler.CreateUom)
		uoms.PUT("/:id", middleware.AdminOnly(cfg), masterHandler.UpdateUom)
		uoms.DELETE("/:id", middleware.AdminOnly(cfg), masterHandler.DeleteUom)
	}

	// Contracts
	dedicatedFix := protected.Group("/contracts/dedicated-fix")
	{
		dedicatedFix.GET("", contractHandler.GetAllDedicatedFix)
		dedicatedFix.GET("/:id", contractHandler.GetDedicatedFixByID)
		dedicatedFix.POST("", middleware.AdminOnly(cfg), contractHandler.CreateDedicatedFix)
		dedicatedFix.PUT("/:id", middleware.AdminOnly(cfg), contractHandler.UpdateDedicatedFix)
		dedicatedFix.PATCH("/:id", middleware.AdminOnly(cfg), contractHandler.UpdateDedicatedFix)
		dedicatedFix.DELETE("/:id", middleware.AdminOnly(cfg), contractHandler.DeleteDedicatedFix)
		dedicatedFix.PATCH("/:id/agreement", middleware.AdminOnly(cfg), contractHandler.UpdateDedicatedFixAgreement)
	}

	dedicatedVar := protected.Group("/contracts/dedicated-var")
	{
		dedicatedVar.GET("", contractHandler.GetAllDedicatedVar)
		dedicatedVar.GET("/:id", contractHandler.GetDedicatedVarByID)
		dedicatedVar.POST("", middleware.AdminOnly(cfg), contractHandler.CreateDedicatedVar)
		dedicatedVar.PUT("/:id", middleware.AdminOnly(cfg), contractHandler.UpdateDedicatedVar)
		dedicatedVar.PATCH("/:id", middleware.AdminOnly(cfg), contractHandler.UpdateDedicatedVar)
		dedicatedVar.DELETE("/:id", middleware.AdminOnly(cfg), contractHandler.DeleteDedicatedVar)
		dedicatedVar.PATCH("/:id/agreement", middleware.AdminOnly(cfg), contractHandler.UpdateDedicatedVarAgreement)
	}

	oncall := protected.Group("/contracts/oncall")
	{
		oncall.GET("", contractHandler.GetAllOncall)
		oncall.GET("/:id", contractHandler.GetOncallByID)
		oncall.POST("", middleware.AdminOnly(cfg), contractHandler.CreateOncall)
		oncall.PUT("/:id", middleware.AdminOnly(cfg), contractHandler.UpdateOncall)
		oncall.PATCH("/:id", middleware.AdminOnly(cfg), contractHandler.UpdateOncall)
		oncall.DELETE("/:id", middleware.AdminOnly(cfg), contractHandler.DeleteOncall)
		oncall.PATCH("/:id/agreement", middleware.AdminOnly(cfg), contractHandler.UpdateOncallAgreement)
	}

	// Audit
	protected.GET("/audit/:entity_type/:entity_id", contractHandler.GetAuditHistory)
	protected.GET("/audit/vendor/:id", contractHandler.GetVendorAuditHistory)
	protected.PATCH("/vendors/:id/agreement", middleware.AdminOnly(cfg), contractHandler.UpdateVendorAgreement)
	protected.PATCH("/mills/:id/agreement", middleware.AdminOnly(cfg), contractHandler.UpdateMillAgreement)

	// Import / Export
	protected.POST("/import/excel", middleware.AdminOnly(cfg), middleware.RequestBodyLimit(defaultImportMaxBytes(cfg)), importHandler.UploadAndParse)
	protected.POST("/import/confirm", middleware.AdminOnly(cfg), middleware.RequestBodyLimit(defaultImportJSONBytes(cfg)), importHandler.ConfirmImport)
	protected.GET("/export/dedicated-fix", exportHandler.ExportDedicatedFix)
	protected.GET("/export/dedicated-var", exportHandler.ExportDedicatedVar)
	protected.GET("/export/oncall", exportHandler.ExportOncall)

	return r
}

func defaultMaxBodyBytes(cfg *config.Config) int64 {
	if cfg == nil || cfg.MaxBodyBytes <= 0 {
		return 10 * 1024 * 1024
	}
	return cfg.MaxBodyBytes
}

func defaultImportMaxBytes(cfg *config.Config) int64 {
	if cfg == nil || cfg.ImportMaxBytes <= 0 {
		return 55 * 1024 * 1024
	}
	return cfg.ImportMaxBytes
}

func defaultImportJSONBytes(cfg *config.Config) int64 {
	if cfg == nil || cfg.ImportJSONBytes <= 0 {
		return 20 * 1024 * 1024
	}
	return cfg.ImportJSONBytes
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
