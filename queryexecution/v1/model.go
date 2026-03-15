package v1

// cancel query request for published Dashboards. Since published dashboards
// have the option of running as the publisher (as opposed to the viewer), the
// cancel request require additional parameters (dashboardName and
// dashboardRevisionId) to retrieve the corresponding user context for rpc calls
// to sql-exec-api.
type CancelPublishedQueryExecutionRequest struct {
	Tokens              []string `json:"tokens"`
	DashboardName       *string  `json:"dashboard_name"`
	DashboardRevisionId *string  `json:"dashboard_revision_id"`
}

type CancelQueryExecutionResponse struct {
	Status []CancelQueryExecutionResponseStatus `json:"status"`
}

type CancelQueryExecutionResponseStatus struct {
	DataToken *string `json:"data_token"`
	Success   *Empty  `json:"success"`
	Pending   *Empty  `json:"pending"`
}

// Represents an empty message, similar to google.protobuf.Empty, which is not
// available in the firm right now..
type Empty struct {
}

// Execute query request for published Dashboards. Since published dashboards
// have the option of running as the publisher, the datasets, warehouse_id are
// excluded from the request and instead read from the source (lakeview-config)
// via the additional parameters (dashboardName and dashboardRevisionId).
type ExecutePublishedDashboardQueryRequest struct {
	DashboardName       *string `json:"dashboard_name"`
	DashboardRevisionId *string `json:"dashboard_revision_id"`
	OverrideWarehouseId *string `json:"override_warehouse_id"`
}

type ExecuteQueryResponse struct {
}

type PendingStatus struct {
	DataToken *string `json:"data_token"`
}

// poll query request for published Dashboards. Since published dashboards have
// the option of running as the publisher (as opposed to the viewer), the
// polling request require additional parameters (dashboardName and
// dashboardRevisionId) to retrieve the corresponding user context for rpc calls
// to sql-exec-api.
type PollPublishedQueryStatusRequest struct {
	Tokens              []string `json:"tokens"`
	DashboardName       *string  `json:"dashboard_name"`
	DashboardRevisionId *string  `json:"dashboard_revision_id"`
}

type PollQueryStatusResponse struct {
	Data []PollQueryStatusResponseData `json:"data"`
}

type PollQueryStatusResponseData struct {
	Status *QueryResponseStatus `json:"status"`
}

type QueryResponseStatus struct {
	Success     *SuccessStatus `json:"success"`
	Pending     *PendingStatus `json:"pending"`
	Canceled    *Empty         `json:"canceled"`
	Closed      *Empty         `json:"closed"`
	StatementId *string        `json:"statement_id"`
}

type SuccessStatus struct {
	DataToken *string `json:"data_token"`
	Truncated *bool   `json:"truncated"`
}
