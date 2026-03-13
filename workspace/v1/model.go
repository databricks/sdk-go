package v1

import (
	"io"
)

type ExportFormat string

const (
	ExportFormat_Source    ExportFormat = "SOURCE"
	ExportFormat_Html      ExportFormat = "HTML"
	ExportFormat_Jupyter   ExportFormat = "JUPYTER"
	ExportFormat_Dbc       ExportFormat = "DBC"
	ExportFormat_RMarkdown ExportFormat = "R_MARKDOWN"
	ExportFormat_Auto      ExportFormat = "AUTO"
	ExportFormat_Raw       ExportFormat = "RAW"
)

// String representation for [fmt.Print].
func (f *ExportFormat) String() string {
	return string(*f)
}

type ExportOutputs string

const (
	ExportOutputs_All  ExportOutputs = "ALL"
	ExportOutputs_None ExportOutputs = "NONE"
)

// String representation for [fmt.Print].
func (f *ExportOutputs) String() string {
	return string(*f)
}

type Language string

const (
	Language_Scala  Language = "SCALA"
	Language_Python Language = "PYTHON"
	Language_Sql    Language = "SQL"
	Language_R      Language = "R"
)

// String representation for [fmt.Print].
func (f *Language) String() string {
	return string(*f)
}

type ObjectType string

const (
	ObjectType_ObjectTypeUnspecified ObjectType = "OBJECT_TYPE_UNSPECIFIED"
	ObjectType_Notebook              ObjectType = "NOTEBOOK"
	ObjectType_Directory             ObjectType = "DIRECTORY"
	ObjectType_Library               ObjectType = "LIBRARY"
	ObjectType_File                  ObjectType = "FILE"
	ObjectType_Repo                  ObjectType = "REPO"
	ObjectType_Dashboard             ObjectType = "DASHBOARD"
)

// String representation for [fmt.Print].
func (f *ObjectType) String() string {
	return string(*f)
}

type Delete struct {
	Path      *string `json:"path"`
	Recursive *bool   `json:"recursive"`
}

type Delete_Response struct {
}

type Export struct {
	Path           *string        `json:"path"`
	Format         *ExportFormat  `json:"format"`
	DirectDownload *bool          `json:"direct_download"`
	Outputs        *ExportOutputs `json:"outputs"`
}

// The request field `direct_download` determines whether a JSON response or
// binary contents are returned by this endpoint..
type Export_Response struct {
	Content  *io.ReadCloser `json:"content"`
	FileType *string        `json:"file_type"`
}

type GetStatus struct {
	Path *string `json:"path"`
}

type Import struct {
	Path      *string        `json:"path"`
	Format    *ExportFormat  `json:"format"`
	Language  *Language      `json:"language"`
	Content   *io.ReadCloser `json:"content"`
	Overwrite *bool          `json:"overwrite"`
}

type Import_Response struct {
}

type List struct {
	Path                   *string `json:"path"`
	NotebooksModifiedAfter *int64  `json:"notebooks_modified_after"`
}

type List_Response struct {
	Objects []ObjectInfo `json:"objects"`
}

type Mkdirs struct {
	Path *string `json:"path"`
}

type Mkdirs_Response struct {
}

// The information of the object in workspace. It will be returned by “list“
// and “get-status“..
type ObjectInfo struct {
	ObjectType *ObjectType `json:"object_type"`
	Path       *string     `json:"path"`
	Language   *Language   `json:"language"`
	CreatedAt  *int64      `json:"created_at"`
	ModifiedAt *int64      `json:"modified_at"`
	ObjectId   *int64      `json:"object_id"`
	Size       *int64      `json:"size"`
	ResourceId *string     `json:"resource_id"`
}
