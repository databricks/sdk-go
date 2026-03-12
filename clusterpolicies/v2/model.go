package v2

type ListOrder string

const (
	ListOrder_Desc ListOrder = "DESC"
	ListOrder_Asc  ListOrder = "ASC"
)

// String representation for [fmt.Print].
func (f *ListOrder) String() string {
	return string(*f)
}

type PolicySortColumn string

const (
	PolicySortColumn_PolicyCreationTime PolicySortColumn = "POLICY_CREATION_TIME"
	PolicySortColumn_PolicyName         PolicySortColumn = "POLICY_NAME"
)

// String representation for [fmt.Print].
func (f *PolicySortColumn) String() string {
	return string(*f)
}

type CreatePolicy struct {
	Name                            *string   `json:"name"`
	Definition                      *string   `json:"definition"`
	Description                     *string   `json:"description"`
	PolicyFamilyId                  *string   `json:"policy_family_id"`
	PolicyFamilyDefinitionOverrides *string   `json:"policy_family_definition_overrides"`
	MaxClustersPerUser              *int64    `json:"max_clusters_per_user"`
	Libraries                       []Library `json:"libraries"`
}

type CreatePolicy_Response struct {
	PolicyId *string `json:"policy_id"`
}

type DeletePolicy struct {
	PolicyId *string `json:"policy_id"`
}

type DeletePolicy_Response struct {
}

type EditPolicy struct {
	PolicyId                        *string   `json:"policy_id"`
	Name                            *string   `json:"name"`
	Definition                      *string   `json:"definition"`
	Description                     *string   `json:"description"`
	PolicyFamilyId                  *string   `json:"policy_family_id"`
	PolicyFamilyDefinitionOverrides *string   `json:"policy_family_definition_overrides"`
	MaxClustersPerUser              *int64    `json:"max_clusters_per_user"`
	Libraries                       []Library `json:"libraries"`
}

type EditPolicy_Response struct {
}

type GetPolicy struct {
	PolicyId *string `json:"policy_id"`
}

type Library struct {
	Jar          *string            `json:"jar"`
	Egg          *string            `json:"egg"`
	Pypi         *PythonPyPiLibrary `json:"pypi"`
	Maven        *MavenLibrary      `json:"maven"`
	Cran         *RCranLibrary      `json:"cran"`
	Whl          *string            `json:"whl"`
	Requirements *string            `json:"requirements"`
}

type ListPolicies struct {
	SortOrder  *ListOrder        `json:"sort_order"`
	SortColumn *PolicySortColumn `json:"sort_column"`
}

type ListPolicies_Response struct {
	Policies []Policy `json:"policies"`
}

type MavenLibrary struct {
	Coordinates *string  `json:"coordinates"`
	Repo        *string  `json:"repo"`
	Exclusions  []string `json:"exclusions"`
}

// Describes a Cluster Policy entity..
type Policy struct {
	PolicyId                        *string   `json:"policy_id"`
	CreatorUserName                 *string   `json:"creator_user_name"`
	CreatedAtTimestamp              *int64    `json:"created_at_timestamp"`
	IsDefault                       *bool     `json:"is_default"`
	Name                            *string   `json:"name"`
	Definition                      *string   `json:"definition"`
	Description                     *string   `json:"description"`
	PolicyFamilyId                  *string   `json:"policy_family_id"`
	PolicyFamilyDefinitionOverrides *string   `json:"policy_family_definition_overrides"`
	MaxClustersPerUser              *int64    `json:"max_clusters_per_user"`
	Libraries                       []Library `json:"libraries"`
}

type PythonPyPiLibrary struct {
	Package *string `json:"package"`
	Repo    *string `json:"repo"`
}

type RCranLibrary struct {
	Package *string `json:"package"`
	Repo    *string `json:"repo"`
}
