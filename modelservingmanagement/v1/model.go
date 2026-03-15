package v1

type ExternalFunctionRequest_HttpMethod string

const (
	ExternalFunctionRequest_HttpMethod_HttpMethodUnspecified ExternalFunctionRequest_HttpMethod = "HTTP_METHOD_UNSPECIFIED"
	ExternalFunctionRequest_HttpMethod_Get                   ExternalFunctionRequest_HttpMethod = "GET"
	ExternalFunctionRequest_HttpMethod_Post                  ExternalFunctionRequest_HttpMethod = "POST"
	ExternalFunctionRequest_HttpMethod_Put                   ExternalFunctionRequest_HttpMethod = "PUT"
	ExternalFunctionRequest_HttpMethod_Delete                ExternalFunctionRequest_HttpMethod = "DELETE"
	ExternalFunctionRequest_HttpMethod_Patch                 ExternalFunctionRequest_HttpMethod = "PATCH"
)

// String representation for [fmt.Print].
func (f *ExternalFunctionRequest_HttpMethod) String() string {
	return string(*f)
}

// Simple Proto message for testing.
type ExternalFunctionRequest struct {
	ConnectionName *string                             `json:"connection_name"`
	Method         *ExternalFunctionRequest_HttpMethod `json:"method"`
	Path           *string                             `json:"path"`
	Json           *string                             `json:"json"`
	Headers        *string                             `json:"headers"`
	Params         *string                             `json:"params"`
	SubDomain      *string                             `json:"sub_domain"`
}

type ExternalFunctionResponse struct {
}
