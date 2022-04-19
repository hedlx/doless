package ops

import (
	"context"
	"fmt"

	api "github.com/hedlx/doless/client"
)

type CreateLambdaM struct {
	Name     string
	Runtime  string
	Endpoint string
}

func CreateLambda(ctx context.Context, lambda CreateLambdaM, path string) (*api.Lambda, error) {
	uploadID, err := upload(ctx, path, true)
	if err != nil {
		return nil, err
	}

	createResp, _, err := client.LambdaApi.
		CreateLambda(ctx).
		CreateLambda(api.CreateLambda{
			Name:     lambda.Name,
			Runtime:  lambda.Runtime,
			Endpoint: lambda.Endpoint,
			Archive:  uploadID,
		}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("Error when calling `LambdaApi.CreateLambda``: %v", err)
	}

	return createResp, nil
}

func GetLambda(ctx context.Context, id string) (*api.Lambda, error) {
	resp, _, err := client.LambdaApi.
		GetLambda(ctx, id).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("Error when calling `LambdaApi.GetLambda``: %v", err)
	}

	return resp, nil
}

func ListLambdas(ctx context.Context) ([]api.Lambda, error) {
	resp, _, err := client.LambdaApi.
		ListLambdas(ctx).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("Error when calling `LambdaApi.ListLambdas``: %v", err)
	}

	return resp, nil
}

func StartLambda(ctx context.Context, id string) (*api.Lambda, error) {
	taskResp, _, err := client.LambdaApi.
		StartLambda(ctx, id).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("Error when calling `LambdaApi.StartLambda``: %v", err)
	}

	res, err := pollTask(ctx, taskResp.GetTask())
	if err != nil {
		return nil, err
	}

	if res.GetStatus() == "FAILED" {
		return nil, fmt.Errorf("Error when starting lambda: %v", res.GetDetails()["error"])
	}

	return GetLambda(ctx, id)
}

func DestroyLambda(ctx context.Context, id string) error {
	taskResp, _, err := client.LambdaApi.
		DestroyLambda(ctx, id).
		Execute()
	if err != nil {
		return fmt.Errorf("Error when calling `LambdaApi.DestroyLambda``: %v", err)
	}

	res, err := pollTask(ctx, taskResp.GetTask())
	if err != nil {
		return err
	}

	if res.GetStatus() == "FAILED" {
		return fmt.Errorf("Error when destroying lambda: %v", res.GetDetails()["error"])
	}

	return nil
}

func DeployLambda(ctx context.Context, input CreateLambdaM, path string) (*api.Lambda, error) {
	lambda, err := CreateLambda(ctx, input, path)
	if err != nil {
		return nil, err
	}

	return StartLambda(ctx, lambda.GetId())
}