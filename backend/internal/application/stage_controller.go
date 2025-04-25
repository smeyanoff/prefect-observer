package application

import (
	"crm-uplift-ii24-backend/internal/application/requests"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/services"
	"crm-uplift-ii24-backend/pkg/logging"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ErrorAddStageToSendpost   string = "[Stage controller] Error AddStageToSendpost"
	ErrorAddSubStage          string = "[Stage controller] Error AddSubStage"
	ErrorDeleteStage          string = "[Stage controller] Error DeleteStage"
	ErrorGetStageDetailedInfo string = "[Stage controller] Error GetStageDetailedInfo"
	ErrorGetSendpostStages    string = "[Stage controller] Error GetSendpostStages"
	ErrorBlockUnblock         string = "[Stage controller] Error BlockUnblock"
	ErrorGetSubStages         string = "[Stage controller] Error GetSubStages"
	ErrorGetStageParameters   string = "[Stage controller] Error GetStageParameters"
	ErrorUpdateParameters     string = "[Stage controller] Error UpdateParameters"
)

type StageController struct {
	stageService *services.StageService
	executor     entity.StageExecutor
}

func NewStageController(stageService *services.StageService, executor entity.StageExecutor) *StageController {
	return &StageController{stageService: stageService, executor: executor}
}

//	@Summary		Add a stage to a sendpost
//	@Description	Adds a new stage to the specified sendpost.
//	@Description	If `previous_stage_id` is provided adds stage after.
//	@Description	If field `next_stage_id` in the previous_stage is not null changes `next_stage_id` in previous_stage on the new provided stage id.
//	@Description	At the same time writes the new provided stage `next_stage_id` with previous_stage `next_stage_id` a.k.a this method allows insert stage between two stages.
//	@Description	Field `type` could be `PARALLEL|SEQUENTIAL|OBSERVER`.
//	@ID				AddStageToSendpost
//	@Tags			Stage
//	@Param			sendpost_id	path	int	true	"Sendpost ID"
//	@Accept			json
//	@Produce		json
//	@Param			request	body		requests.Stage	true	"Stage creation data"
//	@Success		201		{object}	responses.Stage	"Successfully added stage"
//	@Failure		400		{string}	string			InvalidRequestBodyErr
//	@Failure		500		{string}	string			"Internal server error"
//	@Router			/sendposts/{sendpost_id}/stages [post]
func (sc *StageController) AddStageToSendpost(ctx *gin.Context) {
	logging.Info("[Stage controller] AddStageToSendpost request")

	idStr := ctx.Param("sendpost_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorAddStageToSendpost, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidRequestBodyErr)
		return
	}

	var request requests.Stage
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logging.Warn(ErrorAddStageToSendpost, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidRequestBodyErr)
		return
	}

	logging.Debug("[Stage Controller] AddStageToSendpost", zap.Int("sendpost_id", id), zap.Any("request", request))

	stage := unmarshalStage(uint(id), &request)

	if err := sc.stageService.AddStage(ctx, stage, request.PreviousStageID); err != nil {
		logging.Warn(ErrorAddStageToSendpost, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusCreated, mapStage(stage))
}

//	@Summary		Get parallel stage sub-stages
//	@Description	Retrieves the sub-stages of the specified stage.
//
//	@ID				GetSubStages
//	@Tags			Stage
//	@Param			sendpost_id	path	int	true	"Sendpost ID"
//	@Param			stage_id	path	int	true	"Stage ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.SendpostStages	"Successfully retrieved stages"
//	@Failure		400	{string}	string						InvalidRequestBodyErr
//	@Failure		500	{string}	string						"Internal server error"
//	@Router			/sendposts/{sendpost_id}/stages/{stage_id}/sub-stages [get]
func (sc *StageController) GetSubStages(ctx *gin.Context) {
	logging.Info("[Stage controller] GetSubStages request")

	idStr := ctx.Param("stage_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorAddStageToSendpost, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidRequestBodyErr)
		return
	}

	logging.Debug("[Stage Controller] GetSubStages", zap.Int("stage_id", id))

	subStages, err := sc.stageService.GetSubStages(ctx, uint(id))
	if err != nil {
		logging.Warn(ErrorGetSubStages, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, mapStages(subStages))
}

//	@Summary		Add a sub-stage to a parent stage
//	@Description	Adds a sub-stage to an existing parent stage.
//	@Description	The sub-stage will be linked to the parent and can have deployment parameters.
//	@Description	Could only add sub-stage to PARALLEL stage type.
//
//	@ID				AddSubStage
//
//	@Param			sendpost_id	path	int	true	"Sendpost ID"
//	@Param			stage_id	path	int	true	"Stage ID"
//
//	@Tags			Stage
//	@Accept			json
//	@Produce		json
//	@Param			request	body		requests.Stage	true	"Sub-stage creation data"
//	@Success		201		{object}	responses.Stage	"Successfully added sub-stage"
//	@Failure		400		{string}	string			InvalidRequestBodyErr
//	@Failure		500		{string}	string			"Internal server error"
//	@Router			/sendposts/{sendpost_id}/stages/{stage_id}/sub-stages [post]
func (sc *StageController) AddSubStage(ctx *gin.Context) {
	logging.Info("[Stage controller] AddSubStage request")

	stageIdStr := ctx.Param("stage_id")
	stageId, err := strconv.Atoi(stageIdStr)
	if err != nil {
		logging.Warn(ErrorAddSubStage, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidRequestBodyErr)
		return
	}

	sendpostIdStr := ctx.Param("sendpost_id")
	sendpostId, err := strconv.Atoi(sendpostIdStr)
	if err != nil {
		logging.Warn(ErrorAddSubStage, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidRequestBodyErr)
		return
	}

	var request *requests.Stage
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logging.Warn(ErrorAddSubStage, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidRequestBodyErr)
		return
	}

	logging.Debug("[Stage controller] AddSubStage", zap.Int("stage_id", stageId), zap.Any("request", request))

	stage := unmarshalStage(uint(sendpostId), request)
	if err := sc.stageService.AddSubStage(ctx, uint(stageId), stage); err != nil {
		logging.Warn(ErrorAddSubStage, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusCreated, mapStage(stage))
}

//	@Summary		Delete a stage
//	@Description	Deletes a stage by its ID. If the stage is linked to other stages, they will be updated accordingly.
//	@ID				DeleteStage
//	@Tags			Stage
//	@Accept			json
//	@Produce		json
//	@Param			sendpost_id	path		int		true	"Sendpost ID"
//	@Param			stage_id	path		int		true	"Stage ID"
//	@Success		200			{string}	string	"Successfully deleted"
//	@Failure		400			{string}	string	"Invalid ID format"
//	@Failure		500			{string}	string	"Internal server error"
//	@Router			/sendposts/{sendpost_id}/stages/{stage_id} [delete]
func (sc *StageController) DeleteStage(ctx *gin.Context) {
	logging.Info("[Stage controller] DeleteStage request")

	idStr := ctx.Param("stage_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorDeleteStage, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	logging.Debug("[Stage controller] DeleteStage", zap.Int("stage_id", id))

	if err := sc.stageService.DeleteStage(ctx, uint(id)); err != nil {
		logging.Warn(ErrorDeleteStage, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, "Succesfully deleted")
}

//	@Summary		Get stage detailed info
//	@Description	Get detailed information about a stage by its ID.
//	@ID				GetStageDetailedInfo
//	@Tags			Stage
//	@Accept			json
//	@Produce		json
//	@Param			sendpost_id	path		int						true	"Sendpost ID"
//	@Param			stage_id	path		int						true	"Stage ID"
//	@Success		200			{object}	responses.StageDetailed	"Successfully retrieved stage"
//	@Failure		400			{string}	string					"Invalid ID format"
//	@Failure		500			{string}	string					"Internal server error"
//	@Router			/sendposts/{sendpost_id}/stages/{stage_id} [get]
func (sc *StageController) GetStageDetailedInfo(ctx *gin.Context) {
	logging.Info("[Stage controller] GetStageDetailedInfo request")

	idStr := ctx.Param("stage_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorGetStageDetailedInfo, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	logging.Debug("[Stage controller] GetStageDetailedInfo", zap.Int("stage_id", id))

	stage, err := sc.stageService.GetStage(ctx, uint(id))
	if err != nil {
		logging.Warn(ErrorGetStageDetailedInfo, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, mapStageDetailed(stage))
}

//	@Summary		Get stages of a sendpost
//	@Description	Get all stages of a sendpost by its ID.
//	@ID				GetSendpostStages
//	@Tags			Stage
//	@Accept			json
//	@Produce		json
//	@Param			sendpost_id	path		int							true	"Sendpost ID"
//	@Success		200			{object}	responses.SendpostStages	"Successfully retrieved stages"
//	@Failure		400			{string}	string						"Invalid ID format"
//	@Failure		500			{string}	string						"Internal server error"
//	@Router			/sendposts/{sendpost_id}/stages [get]
func (sc *StageController) GetSendpostStages(ctx *gin.Context) {
	logging.Info("[Stage controller] GetSendpostStages request")

	idStr := ctx.Param("sendpost_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorGetSendpostStages, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	logging.Debug("[Stage controller] GetSendpostStages", zap.Int("sendpost_id", id))

	stages, err := sc.stageService.GetSendpostStages(ctx, uint(id))
	if err != nil {
		logging.Warn(ErrorGetSendpostStages, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, mapStages(stages))
}

//	@Summary		Block/Unblock a stage
//	@Description	Block or unblock a stage by its ID.
//	@ID				BlockUnblockStage
//	@Tags			Stage
//	@Accept			json
//	@Produce		json
//	@Param			sendpost_id	path		int		true	"Sendpost ID"
//	@Param			stage_id	path		int		true	"Stage ID"
//	@Success		200			{string}	string	"Successfully blocked/unblocked"
//	@Failure		400			{string}	string	"Invalid ID format"
//	@Router			/sendposts/{sendpost_id}/stages/{stage_id} [patch]
func (sc *StageController) BlockUnblockStage(ctx *gin.Context) {
	logging.Info("[Stage controller] BlockUnblockStage request")

	idStr := ctx.Param("stage_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorBlockUnblock, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	logging.Debug("[Stage controller] BlockUnblockStage", zap.Int("stage_id", id))

	if err := sc.stageService.BlockUnblockStage(ctx, uint(id)); err != nil {
		logging.Warn(ErrorBlockUnblock, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, "Succesfully blocked/unblocked")
}

//	@Summary		Get stage parameters
//	@Description	Get parameters of a stage by its ID.
//	@ID				GetStageParameters
//	@Tags			Stage Info
//	@Accept			json
//	@Produce		json
//	@Param			deployment_id	path		string					true	"Deployment ID"
//	@Success		200				{object}	map[string]interface{}	"Successfully retrieved parameters"
//	@Failure		400				{string}	string					"Invalid ID format"
//	@Failure		500				{string}	string					"Internal server error"
//	@Router			/prefectV2/{deployment_id}/parameters [get]
func (sc *StageController) GetStageParameters(ctx *gin.Context) {
	logging.Info("[Stage controller] GetStageParameters request")

	idStr := ctx.Param("deployment_id")

	parameters, err := sc.executor.GetDeploymentParameters(ctx, idStr)
	if err != nil {
		logging.Warn(ErrorGetStageParameters, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, parameters)
}

//	@Summary		Update stage parameters
//	@Description	Update parameters of a stage by its ID.
//	@ID				UpdateStageParameters
//	@Tags			Stage
//	@Accept			json
//	@Produce		json
//	@Param			sendpost_id	path		int					true	"Sendpost ID"
//	@Param			stage_id	path		int					true	"Stage ID"
//	@Param			parameters	body		requests.Parameters	true	"Parameters"
//	@Success		200			{string}	string				"Successfully updated parameters"
//	@Failure		400			{string}	string				"Invalid ID format"
//	@Failure		500			{string}	string				"Internal server error"
//	@Router			/sendposts/{sendpost_id}/stages/{stage_id} [put]
func (sc *StageController) UpdateParameters(ctx *gin.Context) {
	logging.Info("[Stage controller] UpdateParameters request")

	idStr := ctx.Param("stage_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.Warn(ErrorUpdateParameters, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidIDErr)
		return
	}

	var request *requests.Parameters
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logging.Warn(ErrorUpdateParameters, zap.Error(err))
		ctx.JSON(http.StatusBadRequest, InvalidRequestBodyErr)
		return
	}

	if err := sc.stageService.UpdateParameters(ctx, uint(id), request.Parameters); err != nil {
		logging.Warn(ErrorUpdateParameters, zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	ctx.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
