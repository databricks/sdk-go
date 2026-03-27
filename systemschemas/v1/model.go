package v1

type DisableSystemSchema struct {
	Schema      *string `json:"schema"`
	MetastoreId *string `json:"metastore_id"`
}

type DisableSystemSchema_Response struct {
}

type EnableSystemSchema struct {
	Schema      *string `json:"schema"`
	MetastoreId *string `json:"metastore_id"`
	CatalogName *string `json:"catalog_name"`
}

type EnableSystemSchema_Response struct {
}

type ListSystemSchemas struct {
	MetastoreId *string `json:"metastore_id"`
	MaxResults  *int    `json:"max_results"`
	PageToken   *string `json:"page_token"`
}

type ListSystemSchemas_Response struct {
	Schemas       []SystemSchemaInfo `json:"schemas"`
	NextPageToken *string            `json:"next_page_token"`
}

type SystemSchemaInfo struct {
	Schema *string `json:"schema"`
	State  *string `json:"state"`
}
