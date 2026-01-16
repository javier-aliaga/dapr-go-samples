package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/javier-aliaga/dapr-go-samples/dapr"

	"github.com/dapr/durabletask-go/workflow"
	"github.com/dapr/kit/logger"
)

var log = logger.NewLogger("api.handlers")
var lastInstanceID string

func RegisterRoutes(mux *http.ServeMux, runtime *dapr.WorkflowRuntime) {
	mux.HandleFunc("/healthz", healthHandler)

	mux.HandleFunc("/workflow/event", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			raiseEvent(w, r, runtime)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/workflow", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			startSimpleWorkflow(w, r, runtime)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func startSimpleWorkflow(w http.ResponseWriter, r *http.Request, runtime *dapr.WorkflowRuntime) {
	client := runtime.Client()
	ctx := context.Background()

	log.Infof("Starting workflow")

	instanceID, err := client.ScheduleWorkflow(ctx, "SimpleWorkflow", workflow.WithStartTime(time.Now()))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to start workflow: %v", err), http.StatusInternalServerError)
		return
	}
	lastInstanceID = instanceID
	resp := "New Workflow Instance created " + instanceID
	writeJSON(w, http.StatusAccepted, resp)
}

func raiseEvent(w http.ResponseWriter, r *http.Request, runtime *dapr.WorkflowRuntime) {
	client := runtime.Client()
	ctx := context.Background()

	err := client.RaiseEvent(ctx, lastInstanceID, "event", workflow.WithEventPayload("testData"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to raise event for workflow %s: %v", lastInstanceID, err), http.StatusInternalServerError)
		return
	}

	resp := "Event raised for " + lastInstanceID
	writeJSON(w, http.StatusAccepted, resp)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}