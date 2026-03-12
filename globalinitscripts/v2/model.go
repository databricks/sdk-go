package v2

import (
	"io"
)

type CreateGlobalInitScript struct {
	Name     *string        `json:"name"`
	Script   *io.ReadCloser `json:"script"`
	Position *int           `json:"position"`
	Enabled  *bool          `json:"enabled"`
}

type CreateGlobalInitScript_Response struct {
	ScriptId *string `json:"script_id"`
}

type DeleteGlobalInitScript struct {
	ScriptId *string `json:"script_id"`
}

type DeleteGlobalInitScript_Response struct {
}

type GetGlobalInitScript struct {
	ScriptId *string `json:"script_id"`
}

type GlobalInitScriptDetails struct {
	ScriptId  *string `json:"script_id"`
	Name      *string `json:"name"`
	Position  *int    `json:"position"`
	Enabled   *bool   `json:"enabled"`
	CreatedBy *string `json:"created_by"`
	CreatedAt *int64  `json:"created_at"`
	UpdatedBy *string `json:"updated_by"`
	UpdatedAt *int64  `json:"updated_at"`
}

type ListGlobalInitScripts struct {
}

type ListGlobalInitScripts_Response struct {
	Scripts []GlobalInitScriptDetails `json:"scripts"`
}

type UpdateGlobalInitScript struct {
	ScriptId *string        `json:"script_id"`
	Name     *string        `json:"name"`
	Script   *io.ReadCloser `json:"script"`
	Position *int           `json:"position"`
	Enabled  *bool          `json:"enabled"`
}

type UpdateGlobalInitScript_Response struct {
}
