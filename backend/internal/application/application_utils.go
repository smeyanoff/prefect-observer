package application

import (
	"crm-uplift-ii24-backend/internal/application/requests"
	"crm-uplift-ii24-backend/internal/application/responses"
	"crm-uplift-ii24-backend/internal/domain/entity"
)

func mapSendpost(sendpost *entity.Sendpost) *responses.Sendpost {
	return &responses.Sendpost{
		ID:               sendpost.ID,
		Name:             sendpost.SendpostName,
		Description:      sendpost.Description,
		State:            string(sendpost.State),
		GlobalParameters: sendpost.GlobalParameters,
	}
}

func mapSendposts(sendposts []*entity.Sendpost) responses.Sendposts {
	var result []*responses.Sendpost
	for _, sendpost := range sendposts {
		result = append(result, mapSendpost(sendpost))
	}
	return result
}

func mapStageDetailed(stage *entity.Stage) *responses.StageDetailed {
	return &responses.StageDetailed{
		ID:              stage.ID,
		State:           stage.State,
		Type:            stage.Type,
		ParentStageID:   stage.ParentStageID,
		DeploymentID:    stage.DeploymnentID,
		StageParameters: stage.StageParameters,
	}
}

func mapStage(stage *entity.Stage) *responses.Stage {
	return &responses.Stage{
		ID:        stage.ID,
		Type:      stage.Type,
		State:     stage.State,
		IsBlocked: stage.IsBlocked,
	}
}

func mapStages(stages []*entity.Stage) responses.SendpostStages {
	var result []*responses.Stage
	for _, stage := range stages {
		result = append(result, mapStage(stage))
	}
	return result
}

func unmarshalStage(sendpostID uint, stageRequest *requests.Stage) *entity.Stage {
	return &entity.Stage{
		SendpostID:      sendpostID,
		Type:            stageRequest.StageType,
		DeploymnentID:   stageRequest.DeploymentID,
		StageParameters: stageRequest.StageParameters,
	}
}
