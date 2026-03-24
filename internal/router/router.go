package router

import (
	"database/sql"
	"strings"

	"alioth-hrc/internal/config"
	"alioth-hrc/internal/handler"
	"alioth-hrc/internal/middleware"
	"alioth-hrc/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func New(cfg *config.Config, gdb *gorm.DB, sqlDB *sql.DB, rdb *redis.Client) *gin.Engine {
	setGinMode(cfg.AppEnv)

	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOriginList(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	r.Use(gin.Logger(), gin.Recovery())

	h := handler.NewHealthHandler(sqlDB, rdb)
	r.GET("/healthz", h.Healthz)
	r.GET("/readyz", h.Readyz)
	r.GET("/demo/ping", h.DemoPing)

	authSvc := service.NewAuthService(gdb, rdb, cfg)
	contactSvc := service.NewContactService(gdb)
	ledgerSvc := service.NewLedgerService(gdb)
	catSvc := service.NewGiftCategoryService(gdb)
	giftSvc := service.NewGiftRecordService(gdb, contactSvc, ledgerSvc, catSvc)
	statsSvc := service.NewStatsService(gdb, rdb)

	authH := handler.NewAuthHandler(authSvc)
	contactH := handler.NewContactHandler(contactSvc)
	ledgerH := handler.NewLedgerHandler(ledgerSvc)
	catH := handler.NewGiftCategoryHandler(catSvc)
	giftH := handler.NewGiftRecordHandler(giftSvc)
	statsH := handler.NewStatsHandler(statsSvc)
	remH := handler.NewReminderHandler()

	v1 := r.Group("/api/v1")
	v1.POST("/auth/register", authH.Register)
	v1.POST("/auth/login", authH.Login)
	v1.POST("/auth/refresh", authH.Refresh)
	v1.POST("/auth/logout", authH.Logout)

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/v1/swagger/doc.json")))

	authed := v1.Group("")
	authed.Use(middleware.BearerAuth(cfg.JWTSecret))

	authed.GET("/contacts", contactH.List)
	authed.POST("/contacts", contactH.Create)
	authed.GET("/contacts/:id", contactH.Get)
	authed.PUT("/contacts/:id", contactH.Update)
	authed.DELETE("/contacts/:id", contactH.Delete)

	authed.GET("/ledgers", ledgerH.List)
	authed.POST("/ledgers", ledgerH.Create)
	authed.GET("/ledgers/:id", ledgerH.Get)
	authed.PUT("/ledgers/:id", ledgerH.Update)
	authed.DELETE("/ledgers/:id", ledgerH.Delete)

	authed.GET("/gift-categories", catH.List)
	authed.POST("/gift-categories", catH.Create)
	authed.DELETE("/gift-categories/:id", catH.Delete)

	authed.GET("/gift-records/export.csv", giftH.ExportCSV)
	authed.GET("/gift-records", giftH.List)
	authed.POST("/gift-records", giftH.Create)
	authed.GET("/gift-records/:id", giftH.Get)
	authed.PUT("/gift-records/:id", giftH.Update)
	authed.DELETE("/gift-records/:id", giftH.Delete)

	authed.GET("/stats/contacts", statsH.ContactSummaries)
	authed.GET("/stats/summary", statsH.Summary)

	authed.GET("/reminders", remH.List)

	return r
}

func setGinMode(appEnv string) {
	switch strings.ToLower(strings.TrimSpace(appEnv)) {
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
