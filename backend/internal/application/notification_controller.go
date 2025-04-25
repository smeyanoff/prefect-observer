package application

import (
	"crm-uplift-ii24-backend/internal/services"
	"crm-uplift-ii24-backend/pkg/logging"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	ErrorSendopostRunNotificator string = "[NotificationController] Error SendopostRunNotificator"
)

type NotificationController struct {
	notificationService *services.SenpostRunNotificationService
	upgrader            websocket.Upgrader
}

func NewNotificationController(allowedOrigins []string, notificationService *services.SenpostRunNotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return isOriginAllowed(allowedOrigins, r)
			},
		},
	}
}

//	@Summary		Connect to WebSocket notifications for sendpost execution
//	@Description	Establishes a WebSocket connection to receive status updates on sendpost execution.
//	@Tags			Notifications
//	@Accept			json
//	@Produce		json
//	@Param			sendpost_id	path		int					true	"Sendpost ID"
//	@Success		101			{string}	string				"Switching Protocols"
//	@Failure		400			{object}	map[string]string	"Invalid ID"
//	@Failure		500			{object}	map[string]string	"Internal Server Error"
//	@Router			/sendposts/{sendpost_id}/run/ws [get]
func (c *NotificationController) SendopostRunNotificatorAddListener(ctx *gin.Context) {
	logging.Info("[NotificationController] SendopostRunNotificatorAddListener")

	idStr := ctx.Param("sendpost_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorSendopostRunNotificator, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	logging.Debug("[NotificationController] SendopostRunNotificatorAddListener", zap.Int("sendpost_id", id))

	conn, err := c.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	if err := c.notificationService.AddListener(uint(id), conn); err != nil {
		logging.Warn(ErrorSendopostRunNotificator, zap.Error(err))
		return
	}

	logging.Debug("[NotificationController] Client connected")
}
