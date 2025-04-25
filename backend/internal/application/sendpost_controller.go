package application

import (
	"crm-uplift-ii24-backend/internal/application/requests"
	"crm-uplift-ii24-backend/internal/services"
	"crm-uplift-ii24-backend/pkg/logging"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ErrorCopySendpost                string = "[SendpostController] Error CopySendpost"
	ErrorAddUpdateSendpostParameters string = "[SendpostController] Error AddUpdateSendpostParameters"
	ErrorDeleteSendpostParameter     string = "[SendpostController] Error DeleteSendpostParameter"
)

type SendpostController struct {
	sendpostService *services.SendpostService
}

func NewSendpostController(sendpostService *services.SendpostService) *SendpostController {
	return &SendpostController{sendpostService: sendpostService}
}

// @Summary		Create a sendpost
// @Description	Creates a new sendpost with the specified parameters
//
// @ID				CreateSendpost
//
// @Tags			Sendpost
// @Accept			json
// @Produce		json
// @Param			request	body		requests.Senpost	true	"Sendpost creation data"
// @Success		201		{object}	responses.Sendpost	"Successfully created"
// @Failure		400		{string}	string				"Validation error"
// @Failure		500		{string}	string				"Internal server error"
// @Router			/sendposts [post]
func (s *SendpostController) CreateSendpost(ctx *gin.Context) {
	logging.Info("[Sendpost controller] CreateSendpost request")

	var request requests.Senpost
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logging.Warn("[Sendpost controller] Error CreateSendpost", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	logging.Debug("[Sendpost controller] CreateSendpost", zap.Any("request", request))

	sendpost, err := s.sendpostService.CreateSendpost(ctx, request.SendpostName, request.Description, request.GlobalParameters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusCreated, mapSendpost(sendpost))
}

// @Summary		Get a sendpost
// @Description	Get a sendpost with provided sendpost id
//
// @ID				GetSendpost
//
// @Tags			Sendpost
// @Accept			json
// @Produce		json
// @Param			sendpost_id	path		int					true	"Sendpost ID"
// @Success		200			{object}	responses.Sendpost	"Successfully get"
// @Failure		400			{string}	string				"Bad Request"
// @Failure		500			{string}	string				"Internal Server Error"
// @Router			/sendposts/{sendpost_id} [get]
func (s *SendpostController) GetSendpost(ctx *gin.Context) {
	logging.Info("[Sendpost controller] GetSendpost request")

	idStr := ctx.Param("sendpost_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn("[Sendpost controller] Error GetSendpost", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	logging.Debug("[Sendpost controller] GetSendpost", zap.Int("sendpost_id", id))

	sendpost, err := s.sendpostService.GetSendpost(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, mapSendpost(sendpost))
}

// @Summary		Get a sendposts
// @Description	Get sendposts
//
// @ID				GetSendposts
//
// @Tags			Sendpost
// @Accept			json
// @Produce		json
// @Success		200	{object}	responses.Sendposts	"Successfully get"
// @Failure		500	{string}	string				"Internal Server Error"
// @Router			/sendposts [get]
func (s *SendpostController) GetSendposts(ctx *gin.Context) {
	logging.Info("[Sendpost controller] GetSendposts request")

	sendposts, err := s.sendpostService.GetSendposts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, mapSendposts(sendposts))
}

// @Summary		Delete a sendpost
// @Description	Deletes a sendpost by its ID along with all associated stages
//
// @ID				DeleteSendpost
//
// @Tags			Sendpost
// @Param			sendpost_id	path		int		true	"Sendpost ID"
// @Success		200			{string}	string	"Successfully deleted"
// @Failure		400			{string}	string	"Invalid ID"
// @Failure		500			{string}	string	"Internal server error"
// @Router			/sendposts/{sendpost_id} [delete]
func (s *SendpostController) DeleteSendpost(ctx *gin.Context) {
	logging.Info("[Sendpost controller] DeleteSendpost request")

	idStr := ctx.Param("sendpost_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn("[Sendpost controller] Error DeleteSendpost", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	logging.Debug("[Sendpost controller] DeleteSendpost", zap.Int("sendpost_id", id))

	if err := s.sendpostService.DeleteSendpost(ctx, uint(id)); err != nil {
		logging.Warn("[Sendpost controller] Error DeleteStage", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, "Succesfully deleted")
}

// @Summary		Copy a sendpost
// @Description	Copies a sendpost by its ID
//
// @ID				CopySendpost
//
// @Tags			Sendpost
// @Param			sendpost_id	path		int					true	"Sendpost ID"
// @Param			request		body		requests.Senpost	true	"Sendpost creation data"
// @Success		201			{object}	responses.Sendpost	"Successfully copied"
// @Failure		400			{string}	string				"Invalid ID"
// @Failure		500			{string}	string				"Internal server error"
// @Router			/sendposts/{sendpost_id} [post]
func (s *SendpostController) CopySendpost(ctx *gin.Context) {
	logging.Info("[Sendpost controller] CopySendpost request")

	idStr := ctx.Param("sendpost_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorCopySendpost, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	logging.Debug(ErrorCopySendpost, zap.Int("sendpost_id", id))

	var request requests.Senpost
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logging.Warn(ErrorCopySendpost, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	sendpost, err := s.sendpostService.CopySendpost(ctx, uint(id), request.SendpostName, request.Description, request.GlobalParameters)
	if err != nil {
		logging.Warn(ErrorCopySendpost, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusCreated, mapSendpost(sendpost))
}

// @Summary		Add or update sendpost parameters
// @Description	Add or update sendpost parameters by its ID
//
// @ID				AddUpdateSendpostParameters
//
// @Tags			Sendpost
// @Param			sendpost_id	path		int					true	"Sendpost ID"
// @Param			request		body		requests.Parameters	true	"Sendpost parameters"
// @Success		200			{string}	string				"Successfully added or updated"
// @Failure		400			{string}	string				"Invalid ID"
// @Failure		500			{string}	string				"Internal server error"
// @Router			/sendposts/{sendpost_id}/parameters [post]
func (s *SendpostController) AddUpdateSendpostParameters(ctx *gin.Context) {
	logging.Info("[Sendpost controller] AddUpdateSendpostParameters request")

	idStr := ctx.Param("sendpost_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorAddUpdateSendpostParameters, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	var request requests.Parameters
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logging.Warn(ErrorCopySendpost, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	logging.Debug("[Sendpost controller] AddUpdateSendpostParameters", zap.Any("sendpost_id", id), zap.Any("parameters", request.Parameters))
	for k, v := range request.Parameters {
		if err := s.sendpostService.AddUpdateSendpostParameter(ctx, uint(id), k, v); err != nil {
			logging.Warn(ErrorAddUpdateSendpostParameters, zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	ctx.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

// @Summary		Delete sendpost parameter
// @Description	Delete sendpost parameter by its ID and key
//
// @ID				DeleteSendpostParameter
//
// @Tags			Sendpost
// @Param			sendpost_id	path		int		true	"Sendpost ID"
// @Param			key			path		string	true	"Parameter key"
// @Success		200			{string}	string	"Successfully deleted"
// @Failure		400			{string}	string	"Invalid ID or key"
// @Failure		500			{string}	string	"Internal server error"
// @Router			/sendposts/{sendpost_id}/parameters/{key} [delete]
func (s *SendpostController) DeleteSendpostParameter(ctx *gin.Context) {
	logging.Info("[Sendpost controller] DeleteSendpostParameter request")

	idStr := ctx.Param("sendpost_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorDeleteSendpostParameter, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	key := ctx.Param("key")
	if key == "" {
		logging.Warn(ErrorDeleteSendpostParameter, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidKeyErr)
		return
	}

	logging.Debug(ErrorDeleteSendpostParameter, zap.Int("sendpost_id", id), zap.String("key", key))

	if err := s.sendpostService.DeleteSendpostParameter(ctx, uint(id), key); err != nil {
		logging.Warn(ErrorDeleteSendpostParameter, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
