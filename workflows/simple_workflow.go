package workflows

import (
	"time"

	"github.com/dapr/durabletask-go/workflow"

	"github.com/dapr/kit/logger"
)

var log = logger.NewLogger("workflows.simple_workflow")

// OrderWorkflow is a sample workflow function.
func SimpleWorkflow(ctx *workflow.WorkflowContext) (any, error) {
	if err := ctx.CallActivity(Activity1).Await(nil); err != nil {
		return nil, err
	}

	if err := ctx.CallActivity(Activity2).Await(nil); err != nil {
		return nil, err
	}

	if err := ctx.WaitForExternalEvent("event", time.Minute*5).Await(nil); err != nil {
		return nil, err
	}

	if err := ctx.CallChildWorkflow(ChildWorkflow).Await(nil); err != nil {
		return nil, err
	}

	return nil, nil
}

func ChildWorkflow(ctx *workflow.WorkflowContext) (any, error) {
	if err := ctx.CallActivity(Activity3).Await(nil); err != nil {
		return nil, err
	}

	return nil, nil
}

func Activity1(workflow.ActivityContext) (any, error) {
	log.Info("Activity 1 called")
	// If you had a traceparent in input or context, you could log it here, e.g.:
	// log.Printf("Traceparent: %s", input.TraceParent)

	time.Sleep(1 * time.Second)

	log.Info("Activity 1 finished")
	return nil, nil
}

func Activity2(workflow.ActivityContext) (any, error) {
	log.Info("Activity 2 called")
	// If you had a traceparent in input or context, you could log it here, e.g.:
	// log.Printf("Traceparent: %s", input.TraceParent)

	time.Sleep(1 * time.Second)

	log.Info("Activity 1 finished")
	return nil, nil
}

func Activity3(workflow.ActivityContext) (any, error) {
	log.Info("Activity 3 called")
	// If you had a traceparent in input or context, you could log it here, e.g.:
	// log.Printf("Traceparent: %s", input.TraceParent)

	time.Sleep(1 * time.Second)

	log.Info("Activity 1 finished")
	return nil, nil
}