package prefectV2

import (
	"bytes"
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/value"
	requests "crm-uplift-ii24-backend/internal/infrastructure/workflow/prefectV2/requests"
	responses "crm-uplift-ii24-backend/internal/infrastructure/workflow/prefectV2/responses"
	"crm-uplift-ii24-backend/pkg/logging"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	applicationJSON                           string = "application/json"
	ErrorCheckFlowRunCompletionByDeploymentID string = "[PrefectClientV2] Error CheckFlowRunCompletionByDeploymentID"
)

type PrefectClientV2 struct {
	prefectApiUrl string
	httpClient    *http.Client
}

// NewPrefectClient initializes and returns a new instance of PrefectClient.
// It sets the Prefect API URL and configures an HTTP client with a timeout.
//
// Parameters:
//
//	prefectApiUrl: The URL of the Prefect API to connect to.
//
// Returns:
//
//	A pointer to a PrefectClient configured with the specified API URL and HTTP client.
func NewPrefectClientV2(prefectApiUrl string, insecureTLS bool) entity.StageExecutor {
	return &PrefectClientV2{
		prefectApiUrl: prefectApiUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecureTLS,
				},
			},
		},
	}
}

// CreateFlowRun initiates a new flow run for a given deployment in the Prefect API.
// It constructs a POST request with the specified deployment ID and parameters,
// sends the request to the Prefect server, and returns the response containing
// the flow run details or an error if the operation fails.
//
// Parameters:
//
//	deploymentID: The ID of the deployment for which the flow run is to be created.
//	parameters: A map of parameters to be passed to the flow run.
//
// Returns:
//
//	A pointer to FlowRunResponse containing details of the created flow run,
//	or an error if the request fails or the server returns a non-201 status code.
func (pc *PrefectClientV2) Run(ctx context.Context, deploymentID string, parameters *map[string]interface{}) (flowRunID *string, flowRunState *value.StateType, err error) {
	url := fmt.Sprintf("%s/deployments/%s/create_flow_run", pc.prefectApiUrl, deploymentID)

	reqBody := requests.FlowRunRequest{
		Parameters: parameters,
	}

	logging.Debug("[PrefectClientV2] Run", zap.Any("request", reqBody))

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("accept", applicationJSON)
	req.Header.Set("Content-Type", applicationJSON)

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, nil, fmt.Errorf("failed to start flow run: %s", resp.Status)
	}

	var response responses.FlowRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, err
	}

	logging.Debug("[PrefectClientV2] Run", zap.Any("response", response))

	if deploymentID != response.DeploymnentID {
		return nil, nil, errors.New("response deployment id doesn't match incoming deployment id")
	}

	if response.StateType != value.Scheduled {
		return nil, nil, errors.New("flow run hasn't sheduled")
	}

	return &response.FlowID, &response.StateType, nil
}

// GetFlowRunStatus retrieves the status of a flow run from the Prefect API using the provided flowRunID.
// It sends a GET request to the Prefect API and returns a FlowRunResponse object if successful.
// Returns an error if the request fails or if the response cannot be decoded.
func (pc *PrefectClientV2) Status(ctx context.Context, flowRunID string) (stateType *value.StateType, err error) {
	url := fmt.Sprintf("%s/flow_runs/%s", pc.prefectApiUrl, flowRunID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", applicationJSON)

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to read flow run status: %s", resp.Status)
	}

	var response responses.FlowRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	logging.Debug("[PrefectClientV2] Status", zap.Any("response", response))

	return &response.StateType, nil
}

// CheckFlowRunCompletionByDeploymentID checks if a flow run has completed for a given deployment ID
// within a specified time interval. It sends a POST request to the Prefect API to retrieve flow run
// history and checks the response for completed states. Returns an error if the flow run hasn't
// completed or if any request/response errors occur.
//
// Parameters:
//
//	ctx - the context for the request
//	hisoryStart - the start time for the history interval
//	historyEnd - the end time for the history interval
//	deploymentID - the ID of the deployment to check
//
// Returns:
//
//	error - an error if the flow run hasn't completed or if any request/response errors occur
func (pc *PrefectClientV2) CheckFlowRunCompletionByDeploymentID(ctx context.Context, hisoryStart time.Time, historyEnd time.Time, deploymentID string) error {
	url := fmt.Sprintf("%s/flow_runs/history", pc.prefectApiUrl)

	reqBody := requests.FlowRunCompletedRequest{
		HistoryStart:           hisoryStart,
		HistoryEnd:             historyEnd,
		HistoryIntervalSeconds: 3600,
		FlowRuns: requests.FlowRuns{
			DeploymentID: requests.DeploymentFilter{
				Any: []string{deploymentID},
			},
			State: requests.StateFilter{
				Type: requests.StateType{
					Any: []string{string(value.Completed)},
				},
			},
			StartTime: requests.TimeFilter{
				After:  hisoryStart,
				Before: historyEnd,
			},
		},
		Sort:  "START_TIME_DESC",
		Limit: 1,
	}

	logging.Debug("[PrefectClientV2] CheckFlowRunCompletionByDeploymentID", zap.Any("request", reqBody))

	body, err := json.Marshal(reqBody)
	if err != nil {
		return logging.WrapError(ErrorCheckFlowRunCompletionByDeploymentID, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return logging.WrapError(ErrorCheckFlowRunCompletionByDeploymentID, err)
	}
	req.Header.Set("accept", applicationJSON)
	req.Header.Set("Content-Type", applicationJSON)

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return logging.WrapError(ErrorCheckFlowRunCompletionByDeploymentID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return logging.WrapError(ErrorCheckFlowRunCompletionByDeploymentID, errors.New("Status not ok"))
	}

	var intervals []responses.Interval
	if err := json.NewDecoder(resp.Body).Decode(&intervals); err != nil {
		return logging.WrapError(ErrorCheckFlowRunCompletionByDeploymentID, errors.New("error intervals decoding"))
	}

	logging.Debug("[PrefectClientV2] CheckFlowRunCompletionByDeploymentID", zap.Any("intervals", intervals))

	for _, interval := range intervals {
		if len(interval.States) > 0 {
			return nil
		}
	}
	return logging.WrapError(ErrorCheckFlowRunCompletionByDeploymentID, errors.New("flow run hasn't completed over period"))
}

// GetDeploymentParameters retrieves the parameters of a specified deployment
// from the Prefect API using the provided deployment ID. It sends a GET request
// to the Prefect server and returns a map of parameters if successful, or an error
// if the request fails or the response cannot be decoded.
//
// Parameters:
//
//	ctx - The context for the HTTP request, allowing for cancellation and timeouts.
//	deploymentID - The unique identifier of the deployment whose parameters are to be fetched.
//
// Returns:
//
//	A map containing the deployment parameters if the request is successful.
//	An error if the request fails or if the response cannot be decoded.
func (pc *PrefectClientV2) GetDeploymentParameters(ctx context.Context, deploymentID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/deployments/%s", pc.prefectApiUrl, deploymentID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", applicationJSON)

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to read deployment parameters: %s", resp.Status)
	}

	var response responses.ParametersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Parameters, nil
}
