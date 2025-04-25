package main

import (
	"crm-uplift-ii24-backend/config"
	_ "crm-uplift-ii24-backend/docs"
	"crm-uplift-ii24-backend/internal/application"
	runstatus "crm-uplift-ii24-backend/internal/infrastructure/notifications/runStatus"
	"crm-uplift-ii24-backend/internal/infrastructure/persistence"
	"crm-uplift-ii24-backend/internal/infrastructure/persistence/repository"
	"crm-uplift-ii24-backend/internal/infrastructure/workflow/prefectV2"
	"crm-uplift-ii24-backend/internal/services"
	"crm-uplift-ii24-backend/internal/services/runners"
	log "crm-uplift-ii24-backend/pkg/logging"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			OBSERVER backend
// @version		1.2
// @description	This is a backend for OBSERVER app.
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host			localhost:8180
// @BasePath		/v1
// @schemes		http
func main() {
	log.Init()
	defer log.Sync()

	cfg := config.LoadConfig()

	// Infrustructure
	db, err := persistence.ConnectDB(cfg)
	if err != nil {
		log.Error("Couldn`t connect to db", zap.String("err", err.Error()))
	}
	if err := persistence.AutoMigrate(db); err != nil {
		log.Error("Failed to migrate", zap.String("err", err.Error()))
	}
	stageExecutor := prefectV2.NewPrefectClientV2(cfg.App.PrefectApiUrl, cfg.App.InsecureSkipVerify)
	sendpostRunNotificator := runstatus.NewNotificatorWS()

	// Repository
	sendpostRepo := repository.NewGormSendpostRepository(db)
	stageRepo := repository.NewGormSendpostStageRepository(db)

	// Services
	stageService := services.NewStageService(stageRepo, sendpostRepo)
	sendpostService := services.NewSendpostService(sendpostRepo, stageService)
	stageRunnerService := services.NewStageRunnerService(stageExecutor, stageService, cfg.App.StageStatusQueryTimeout)
	stageRunnerFactory := runners.NewStageRunnerFactory(stageRunnerService, stageService)
	sendpostRunNotificationService := services.NewSenpostRunNotificationService(sendpostRunNotificator)
	sendpostRunnerService := services.NewSendpostRunService(sendpostService, stageService, sendpostRunNotificationService, stageRunnerFactory)

	// Controllers
	sendpostController := application.NewSendpostController(sendpostService)
	stageController := application.NewStageController(stageService, stageExecutor)
	sendpostRunnerController := application.NewSendpostRunnerController(sendpostRunnerService)
	notificationController := application.NewNotificationController(cfg.CORS.AllowOrigins, sendpostRunNotificationService)

	// Router
	r := gin.Default()

	// Настройка CORS через данные из конфигурации
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           12 * time.Hour,
	}))
	apiV1 := r.Group("/v1")

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// sendposts
	apiV1.POST("/sendposts", sendpostController.CreateSendpost)
	apiV1.GET("/sendposts", sendpostController.GetSendposts)
	apiV1.POST("/sendposts/:sendpost_id", sendpostController.CopySendpost)
	apiV1.GET("/sendposts/:sendpost_id", sendpostController.GetSendpost)
	apiV1.DELETE("/sendposts/:sendpost_id", sendpostController.DeleteSendpost)
	apiV1.POST("/sendposts/:sendpost_id/parameters", sendpostController.AddUpdateSendpostParameters)
	apiV1.DELETE("/sendposts/:sendpost_id/parameters/:key", sendpostController.DeleteSendpostParameter)

	// sendpost run
	apiV1.POST("/sendposts/:sendpost_id/run", sendpostRunnerController.Start)

	// notifications
	apiV1.GET("/sendposts/:sendpost_id/run/ws", notificationController.SendopostRunNotificatorAddListener)

	// stages
	apiV1.POST("/sendposts/:sendpost_id/stages", stageController.AddStageToSendpost)
	apiV1.GET("/sendposts/:sendpost_id/stages", stageController.GetSendpostStages)
	apiV1.GET("/sendposts/:sendpost_id/stages/:stage_id", stageController.GetStageDetailedInfo)
	apiV1.DELETE("/sendposts/:sendpost_id/stages/:stage_id", stageController.DeleteStage)
	apiV1.PATCH("/sendposts/:sendpost_id/stages/:stage_id", stageController.BlockUnblockStage)
	apiV1.PUT("/sendposts/:sendpost_id/stages/:stage_id", stageController.UpdateParameters)

	// stage info
	apiV1.GET("/prefectV2/:deployment_id/parameters", stageController.GetStageParameters)

	// sub-stages
	apiV1.POST("/sendposts/:sendpost_id/stages/:stage_id/sub-stages", stageController.AddSubStage)
	apiV1.GET("/sendposts/:sendpost_id/stages/:stage_id/sub-stages", stageController.GetSubStages)

	// Start Server
	log.Info("Server is running on port: 8081")
	r.Run("0.0.0.0" + ":" + cfg.App.Port)
}
