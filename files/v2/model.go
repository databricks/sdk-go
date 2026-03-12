package v2

import (
	"io"
)

type AddBlock struct {
	Handle *int64         `json:"handle"`
	Data   *io.ReadCloser `json:"data"`
}

type AddBlock_Response struct {
}

type Close struct {
	Handle *int64 `json:"handle"`
}

type Close_Response struct {
}

type Create struct {
	Path      *string `json:"path"`
	Overwrite *bool   `json:"overwrite"`
}

type Create_Response struct {
	Handle *int64 `json:"handle"`
}

type Delete struct {
	Path      *string `json:"path"`
	Recursive *bool   `json:"recursive"`
}

type Delete_Response struct {
}

// Stores the attributes of a file or directory..
type FileInfo struct {
	Path             *string `json:"path"`
	IsDir            *bool   `json:"is_dir"`
	FileSize         *int64  `json:"file_size"`
	ModificationTime *int64  `json:"modification_time"`
}

type GetStatus struct {
	Path *string `json:"path"`
}

type GetStatus_Response struct {
	Path             *string `json:"path"`
	IsDir            *bool   `json:"is_dir"`
	FileSize         *int64  `json:"file_size"`
	ModificationTime *int64  `json:"modification_time"`
}

type ListStatus struct {
	Path *string `json:"path"`
}

type ListStatus_Response struct {
	Files []FileInfo `json:"files"`
}

type MkDirs struct {
	Path *string `json:"path"`
}

type MkDirs_Response struct {
}

type Move struct {
	SourcePath      *string `json:"source_path"`
	DestinationPath *string `json:"destination_path"`
}

type Move_Response struct {
}

type Put struct {
	Path      *string        `json:"path"`
	Contents  *io.ReadCloser `json:"contents"`
	Overwrite *bool          `json:"overwrite"`
}

type Put_Response struct {
}

type Read struct {
	Path   *string `json:"path"`
	Offset *int64  `json:"offset"`
	Length *int64  `json:"length"`
}

type Read_Response struct {
	BytesRead *int64         `json:"bytes_read"`
	Data      *io.ReadCloser `json:"data"`
}
