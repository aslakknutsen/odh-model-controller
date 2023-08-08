package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-logr/logr"
	predictorv1 "github.com/kserve/modelmesh-serving/apis/serving/v1alpha1"
	inferenceservicev1 "github.com/kserve/modelmesh-serving/apis/serving/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewModelCheck(client client.Client, log logr.Logger) http.Handler {
	return &modelCheckHandler{Client: client, Log: log}
}

type modelCheckHandler struct {
	Client client.Client
	Log    logr.Logger
}

func (m *modelCheckHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	ctx := context.Background()

	m.Log.Info("received request", "url", req.URL)

	ns := req.URL.Query().Get("ns")
	if ns == "" {
		m.Log.Error(nil, "missing namespace")
		respond(resp, false)
		return
	}
	modelID := req.URL.Query().Get("modelid")
	if modelID == "" {
		m.Log.Error(nil, "missing modelid")
		respond(resp, false)
		return
	}

	inferenceService := &inferenceservicev1.InferenceService{}
	err := m.Client.Get(ctx, types.NamespacedName{Namespace: ns, Name: modelID}, inferenceService)
	if err != nil {
		m.Log.Error(err, "failed to get inference service")
		respond(resp, false)
		return
	}

	servingRuntimeName := *inferenceService.Spec.Predictor.Model.Runtime
	if servingRuntimeName == "" {
		m.Log.Error(nil, "missing servingruntime.spec.predicator.model.runtime")
		respond(resp, false)
		return
	}
	servingRuntime := &predictorv1.ServingRuntime{}
	err = m.Client.Get(ctx, types.NamespacedName{Namespace: ns, Name: servingRuntimeName}, servingRuntime)
	if err != nil {
		m.Log.Error(err, "failed to get serving runtime")
		respond(resp, false)
		return
	}
	if servingRuntime.Annotations[EnableAuthAnnotation] != "true" {
		respond(resp, true)
		return
	}
	m.Log.Info("serving runtime has auth-enabled", "servingruntime", servingRuntimeName, "inferenceservice", modelID)
	respond(resp, false)
}

func respond(resp http.ResponseWriter, anonymousAccessOk bool) {
	response := "{\"anonymous\": " + strconv.FormatBool(anonymousAccessOk) + "}\n"
	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(200)
	resp.Write([]byte(response))
}
