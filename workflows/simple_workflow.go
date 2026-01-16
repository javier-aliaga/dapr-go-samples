package workflows

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/dapr/durabletask-go/workflow"

	"github.com/dapr/kit/logger"
)

var log = logger.NewLogger("workflows.simple_workflow")
var tracer = otel.Tracer("workflows.simple_workflow")

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

func Activity1(ctx workflow.ActivityContext) (any, error) {
	log.Infof("Activity 1 called with traceparent: %s", ctx.GetTraceContext().TraceParent)
	_, childSpan := tracer.Start(ctx.Context(), "Custom||Activity1")
	defer childSpan.End()

	time.Sleep(1 * time.Second)

	log.Info("Activity 1 finished")
	return nil, nil
}

func Activity2(ctx workflow.ActivityContext) (any, error) {
	log.Infof("Activity 2 called with traceparent: %s", ctx.GetTraceContext().TraceParent)
	// If you had a traceparent in input or context, you could log it here, e.g.:
	// log.Printf("Traceparent: %s", input.TraceParent)
	_, childSpan := tracer.Start(ctx.Context(), "Custom||Activity2")
	defer childSpan.End()
	time.Sleep(1 * time.Second)

	log.Info("Activity 2 finished")
	return nil, nil
}

func Activity3(ctx workflow.ActivityContext) (any, error) {
	log.Infof("Activity 3 called with traceparent: %s", ctx.GetTraceContext().TraceParent)
	// If you had a traceparent in input or context, you could log it here, e.g.:
	// log.Printf("Traceparent: %s", input.TraceParent)

	// Create a new context
	newCtx := context.Background()

	// Create a TextMapCarrier with the traceparent
	carrier := propagation.MapCarrier{}
	carrier.Set("traceparent", ctx.GetTraceContext().TraceParent)

	// Use the TraceContext propagator to extract the trace context
	propagator := propagation.TraceContext{}
	newCtx = propagator.Extract(newCtx, carrier)

	_, childSpan := tracer.Start(newCtx, "Custom||Activity3")
	defer childSpan.End()

	time.Sleep(1 * time.Second)

	log.Info("Activity 3 finished")
	return nil, nil
}