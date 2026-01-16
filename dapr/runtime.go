package dapr

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/dapr/durabletask-go/workflow"
	"github.com/dapr/go-sdk/client"

	"github.com/dapr/kit/logger"

	"github.com/javier-aliaga/dapr-go-samples/workflows"
)

var log = logger.NewLogger("dapr.runtime")

type WorkflowRuntime struct {
	client  *workflow.Client
	runtime *workflow.Registry
}

func (w *WorkflowRuntime) Client() *workflow.Client {
	return w.client
}

// StartWorkflowRuntime bootstraps the Dapr Workflow runtime and registers workflows.
func StartWorkflowRuntime(ctx context.Context) (*WorkflowRuntime, error) {
	r := workflow.NewRegistry()

	// Register your workflows and activities
	err := r.AddWorkflowN("SimpleWorkflow", workflows.SimpleWorkflow)
	if err != nil {
		return nil, fmt.Errorf("register workflow: %w", err)
	}
	err = r.AddWorkflowN("ChildWorkflow", workflows.ChildWorkflow)
	if err != nil {
		return nil, fmt.Errorf("register workflow: %w", err)
	}
	err = r.AddActivity(workflows.Activity1)
	if err != nil {
		return nil, fmt.Errorf("register activity: %w", err)
	}
	err = r.AddActivity(workflows.Activity2)
	if err != nil {
		return nil, fmt.Errorf("register activity: %w", err)
	}
	err = r.AddActivity(workflows.Activity3)
	if err != nil {
		return nil, fmt.Errorf("register activity: %w", err)
	}

	wClient, err := client.NewWorkflowClient(grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	//wClient, err := client.NewWorkflowClient()
	if err != nil {
		return nil, fmt.Errorf("create workflow client: %w", err)
	}

	// Start runtime in background
	go func() {
		if err := wClient.StartWorker(ctx, r); err != nil {
			// In a real app, propagate or log more robustly
			log.Info("workflow runtime stopped: %v\n", err)
		}
	}()

	return &WorkflowRuntime{
		client:  wClient,
		runtime: r,
	}, nil
}