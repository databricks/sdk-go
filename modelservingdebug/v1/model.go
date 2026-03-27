package v1

import (
	"io"
)

// * Proto version of com.databricks.rpc.HttpOverRpcResponse.
//
// This message can be specially handled in UnaryRpcService with JettyRPC when
// the advanced feature CustomHandlingForHttpOverRpcProtoResponse is enabled -
// bypass the RPC serializer and populate HTTP status, response headers and
// response body from the proto message directly.
//
// Don't add/modify the fields before being aware of the implications..
type ExportMetricsResponse struct {
	Status     *int           `json:"status"`
	RawContent *io.ReadCloser `json:"raw_content"`
	Headers    []Header       `json:"headers"`
}

type ExportServedModelServiceLogFile struct {
	Name         *string `json:"name"`
	DeploymentId *string `json:"deployment_id"`
	Date         *string `json:"date"`
	InstanceId   *string `json:"instance_id"`
	FileName     *string `json:"file_name"`
}

type GetExportEndpointMetrics struct {
	Name *string `json:"name"`
}

type GetServedModelBuildLogs struct {
	Name            *string `json:"name"`
	ServedModelName *string `json:"served_model_name"`
}

type GetServedModelBuildLogs_Response struct {
	Logs *string `json:"logs"`
}

type GetServedModelLogs struct {
	Name            *string `json:"name"`
	ServedModelName *string `json:"served_model_name"`
}

type GetServedModelLogs_Response struct {
	Logs *string `json:"logs"`
}

type Header struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

// * Proto version of com.databricks.rpc.HttpOverRpcResponse.
//
// This message can be specially handled in UnaryRpcService with JettyRPC when
// the advanced feature CustomHandlingForHttpOverRpcProtoResponse is enabled -
// bypass the RPC serializer and populate HTTP status, response headers and
// response body from the proto message directly.
//
// Don't add/modify the fields before being aware of the implications..
type HttpOverRpcProtoResponse struct {
	Status     *int           `json:"status"`
	RawContent *io.ReadCloser `json:"raw_content"`
	Headers    []Header       `json:"headers"`
}

type ListEndpointServiceLogFiles struct {
	Name *string `json:"name"`
}

type ListEndpointServiceLogFiles_Response struct {
	Name  *string                                                          `json:"name"`
	Files []ListEndpointServiceLogFiles_Response_ServedModelServiceLogFile `json:"files"`
}

type ListEndpointServiceLogFiles_Response_ServedModelServiceLogFile struct {
	ServedEntityName    *string `json:"served_entity_name"`
	ServedEntityVersion *string `json:"served_entity_version"`
	DeploymentId        *string `json:"deployment_id"`
	FileName            *string `json:"file_name"`
	Size                *int64  `json:"size"`
	Date                *string `json:"date"`
	InstanceId          *string `json:"instance_id"`
}
