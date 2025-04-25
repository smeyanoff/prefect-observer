package application

import (
	"crm-uplift-ii24-backend/internal/services"
	"crm-uplift-ii24-backend/pkg/logging"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ErrorRunningSendpost = "[Sendpost Runner Controller] Error running sendpost"
)

type SendpostRunnerController struct {
	sendpostRunnerService *services.SendpostRunnerService
}

func NewSendpostRunnerController(sendpostRunnerService *services.SendpostRunnerService) *SendpostRunnerController {
	return &SendpostRunnerController{sendpostRunnerService: sendpostRunnerService}
}

//	@Summary		Start the sendpost
//	@Description	Start the sendpost
//	@Tags			Sendpost Runner
//	@Accept			json
//	@Produce		json
//	@Param			sendpost_id	path		int		true	"Sendpost ID"
//	@Success		202			{object}	string	"Accepted"
//	@Failure		400			{object}	string	"Invalid ID"
//	@Failure		500			{object}	string	"Internal server error"
//	@Router			/sendposts/{sendpost_id}/run [post]
func (c *SendpostRunnerController) Start(ctx *gin.Context) {
	logging.Info("[Sendpost Runner Controller] Start request")

	idStr := ctx.Param("sendpost_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorRunningSendpost, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	c.sendpostRunnerService.Start(ctx, uint(id))

	ctx.JSON(http.StatusAccepted, "Accepted")
}
