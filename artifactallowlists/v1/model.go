package v1

type ArtifactType string

const (
	ArtifactType_ArtifactTypeUnspecified ArtifactType = "ARTIFACT_TYPE_UNSPECIFIED"
	ArtifactType_InitScript              ArtifactType = "INIT_SCRIPT"
	ArtifactType_LibraryJar              ArtifactType = "LIBRARY_JAR"
	ArtifactType_LibraryMaven            ArtifactType = "LIBRARY_MAVEN"
)

// String representation for [fmt.Print].
func (f *ArtifactType) String() string {
	return string(*f)
}

type ArtifactMatcher_MatchType string

const (
	ArtifactMatcher_MatchType_MatchTypeUnspecified ArtifactMatcher_MatchType = "MATCH_TYPE_UNSPECIFIED"
	ArtifactMatcher_MatchType_PrefixMatch          ArtifactMatcher_MatchType = "PREFIX_MATCH"
)

// String representation for [fmt.Print].
func (f *ArtifactMatcher_MatchType) String() string {
	return string(*f)
}

type ArtifactAllowlistInfo struct {
	ArtifactMatchers []ArtifactMatcher `json:"artifact_matchers"`
	MetastoreId      *string           `json:"metastore_id"`
	CreatedBy        *string           `json:"created_by"`
	CreatedAt        *int64            `json:"created_at"`
}

type ArtifactMatcher struct {
	Artifact  *string                    `json:"artifact"`
	MatchType *ArtifactMatcher_MatchType `json:"match_type"`
}

type GetArtifactAllowlist struct {
	ArtifactType *ArtifactType `json:"artifact_type"`
}

type SetArtifactAllowlist struct {
	ArtifactType     *ArtifactType     `json:"artifact_type"`
	ArtifactMatchers []ArtifactMatcher `json:"artifact_matchers"`
	MetastoreId      *string           `json:"metastore_id"`
	CreatedBy        *string           `json:"created_by"`
	CreatedAt        *int64            `json:"created_at"`
}
