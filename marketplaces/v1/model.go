package v1

type AssetType string

const (
	AssetType_AssetTypeUnspecified        AssetType = "ASSET_TYPE_UNSPECIFIED"
	AssetType_AssetTypeGitRepo            AssetType = "ASSET_TYPE_GIT_REPO"
	AssetType_AssetTypeDataTable          AssetType = "ASSET_TYPE_DATA_TABLE"
	AssetType_AssetTypeModel              AssetType = "ASSET_TYPE_MODEL"
	AssetType_AssetTypeNotebook           AssetType = "ASSET_TYPE_NOTEBOOK"
	AssetType_AssetTypeMedia              AssetType = "ASSET_TYPE_MEDIA"
	AssetType_AssetTypePartnerIntegration AssetType = "ASSET_TYPE_PARTNER_INTEGRATION"
	AssetType_AssetTypeApp                AssetType = "ASSET_TYPE_APP"
	AssetType_AssetTypeMcp                AssetType = "ASSET_TYPE_MCP"
)

// String representation for [fmt.Print].
func (f *AssetType) String() string {
	return string(*f)
}

type Category string

const (
	Category_AdvertisingAndMarketing    Category = "ADVERTISING_AND_MARKETING"
	Category_ClimateAndEnvironment      Category = "CLIMATE_AND_ENVIRONMENT"
	Category_Commerce                   Category = "COMMERCE"
	Category_Demographics               Category = "DEMOGRAPHICS"
	Category_Economics                  Category = "ECONOMICS"
	Category_Education                  Category = "EDUCATION"
	Category_Energy                     Category = "ENERGY"
	Category_Financial                  Category = "FINANCIAL"
	Category_Gaming                     Category = "GAMING"
	Category_Geospatial                 Category = "GEOSPATIAL"
	Category_Health                     Category = "HEALTH"
	Category_LookupTables               Category = "LOOKUP_TABLES"
	Category_Manufacturing              Category = "MANUFACTURING"
	Category_Media                      Category = "MEDIA"
	Category_Other                      Category = "OTHER"
	Category_PublicSector               Category = "PUBLIC_SECTOR"
	Category_Retail                     Category = "RETAIL"
	Category_Security                   Category = "SECURITY"
	Category_ScienceAndResearch         Category = "SCIENCE_AND_RESEARCH"
	Category_Sports                     Category = "SPORTS"
	Category_TransportationAndLogistics Category = "TRANSPORTATION_AND_LOGISTICS"
	Category_TravelAndTourism           Category = "TRAVEL_AND_TOURISM"
	Category_OpenSource                 Category = "OPEN_SOURCE"
)

// String representation for [fmt.Print].
func (f *Category) String() string {
	return string(*f)
}

type Cost string

const (
	Cost_Free Cost = "FREE"
	Cost_Paid Cost = "PAID"
)

// String representation for [fmt.Print].
func (f *Cost) String() string {
	return string(*f)
}

type DataRefresh string

const (
	DataRefresh_None      DataRefresh = "NONE"
	DataRefresh_Second    DataRefresh = "SECOND"
	DataRefresh_Minute    DataRefresh = "MINUTE"
	DataRefresh_Hourly    DataRefresh = "HOURLY"
	DataRefresh_Daily     DataRefresh = "DAILY"
	DataRefresh_Weekly    DataRefresh = "WEEKLY"
	DataRefresh_Monthly   DataRefresh = "MONTHLY"
	DataRefresh_Quarterly DataRefresh = "QUARTERLY"
	DataRefresh_Yearly    DataRefresh = "YEARLY"
)

// String representation for [fmt.Print].
func (f *DataRefresh) String() string {
	return string(*f)
}

type DeltaSharingRecipientType string

const (
	DeltaSharingRecipientType_DeltaSharingRecipientTypeDatabricks DeltaSharingRecipientType = "DELTA_SHARING_RECIPIENT_TYPE_DATABRICKS"
	DeltaSharingRecipientType_DeltaSharingRecipientTypeOpen       DeltaSharingRecipientType = "DELTA_SHARING_RECIPIENT_TYPE_OPEN"
)

// String representation for [fmt.Print].
func (f *DeltaSharingRecipientType) String() string {
	return string(*f)
}

type ExchangeFilterType string

const (
	ExchangeFilterType_GlobalMetastoreId ExchangeFilterType = "GLOBAL_METASTORE_ID"
)

// String representation for [fmt.Print].
func (f *ExchangeFilterType) String() string {
	return string(*f)
}

type FileParentType string

const (
	FileParentType_Provider        FileParentType = "PROVIDER"
	FileParentType_Listing         FileParentType = "LISTING"
	FileParentType_ListingResource FileParentType = "LISTING_RESOURCE"
)

// String representation for [fmt.Print].
func (f *FileParentType) String() string {
	return string(*f)
}

type FileStatus string

const (
	FileStatus_FileStatusPublished          FileStatus = "FILE_STATUS_PUBLISHED"
	FileStatus_FileStatusStaging            FileStatus = "FILE_STATUS_STAGING"
	FileStatus_FileStatusSanitizing         FileStatus = "FILE_STATUS_SANITIZING"
	FileStatus_FileStatusSanitizationFailed FileStatus = "FILE_STATUS_SANITIZATION_FAILED"
)

// String representation for [fmt.Print].
func (f *FileStatus) String() string {
	return string(*f)
}

type FulfillmentType string

const (
	FulfillmentType_RequestAccess FulfillmentType = "REQUEST_ACCESS"
	FulfillmentType_Install       FulfillmentType = "INSTALL"
)

// String representation for [fmt.Print].
func (f *FulfillmentType) String() string {
	return string(*f)
}

type InstallationStatus string

const (
	InstallationStatus_Installed InstallationStatus = "INSTALLED"
	InstallationStatus_Failed    InstallationStatus = "FAILED"
)

// String representation for [fmt.Print].
func (f *InstallationStatus) String() string {
	return string(*f)
}

type ListingShareType string

const (
	ListingShareType_Sample ListingShareType = "SAMPLE"
	ListingShareType_Full   ListingShareType = "FULL"
)

// String representation for [fmt.Print].
func (f *ListingShareType) String() string {
	return string(*f)
}

type ListingStatus string

const (
	ListingStatus_Draft     ListingStatus = "DRAFT"
	ListingStatus_Pending   ListingStatus = "PENDING"
	ListingStatus_Published ListingStatus = "PUBLISHED"
	ListingStatus_Suspended ListingStatus = "SUSPENDED"
)

// String representation for [fmt.Print].
func (f *ListingStatus) String() string {
	return string(*f)
}

type ListingTagType string

const (
	ListingTagType_ListingTagTypeUnspecified ListingTagType = "LISTING_TAG_TYPE_UNSPECIFIED"
	ListingTagType_ListingTagTypeLanguage    ListingTagType = "LISTING_TAG_TYPE_LANGUAGE"
	ListingTagType_ListingTagTypeTask        ListingTagType = "LISTING_TAG_TYPE_TASK"
)

// String representation for [fmt.Print].
func (f *ListingTagType) String() string {
	return string(*f)
}

type ListingType string

const (
	ListingType_Standard     ListingType = "STANDARD"
	ListingType_Personalized ListingType = "PERSONALIZED"
)

// String representation for [fmt.Print].
func (f *ListingType) String() string {
	return string(*f)
}

type MarketplaceFileType string

const (
	MarketplaceFileType_ProviderIcon     MarketplaceFileType = "PROVIDER_ICON"
	MarketplaceFileType_EmbeddedNotebook MarketplaceFileType = "EMBEDDED_NOTEBOOK"
	MarketplaceFileType_App              MarketplaceFileType = "APP"
)

// String representation for [fmt.Print].
func (f *MarketplaceFileType) String() string {
	return string(*f)
}

type PersonalizationRequestStatus string

const (
	PersonalizationRequestStatus_New            PersonalizationRequestStatus = "NEW"
	PersonalizationRequestStatus_RequestPending PersonalizationRequestStatus = "REQUEST_PENDING"
	PersonalizationRequestStatus_Fulfilled      PersonalizationRequestStatus = "FULFILLED"
	PersonalizationRequestStatus_Denied         PersonalizationRequestStatus = "DENIED"
)

// String representation for [fmt.Print].
func (f *PersonalizationRequestStatus) String() string {
	return string(*f)
}

type Visibility string

const (
	Visibility_Public  Visibility = "PUBLIC"
	Visibility_Private Visibility = "PRIVATE"
)

// String representation for [fmt.Print].
func (f *Visibility) String() string {
	return string(*f)
}

type AddExchangeForListingRequest struct {
	ListingId  *string `json:"listing_id"`
	ExchangeId *string `json:"exchange_id"`
}

type AddExchangeForListingResponse struct {
	ExchangeForListing *ExchangeListing `json:"exchange_for_listing"`
}

type BatchGetListingsRequest struct {
	Ids []string `json:"ids"`
}

type BatchGetListingsResponse struct {
	Listings []Listing `json:"listings"`
}

type BatchGetProvidersRequest struct {
	Ids []string `json:"ids"`
}

type BatchGetProvidersResponse struct {
	Providers []ProviderInfo `json:"providers"`
}

type ConsumerTerms struct {
	Version *string `json:"version"`
}

// contact info for the consumer requesting data or performing a listing
// installation.
type ContactInfo struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Company   *string `json:"company"`
}

type CreateExchangeFilterRequest struct {
	Filter *ExchangeFilter `json:"filter"`
}

type CreateExchangeFilterResponse struct {
	FilterId *string `json:"filter_id"`
}

type CreateExchangeRequest struct {
	Exchange *Exchange `json:"exchange"`
}

type CreateExchangeResponse struct {
	ExchangeId *string `json:"exchange_id"`
}

type CreateFile struct {
	FileParent          *FileParent          `json:"file_parent"`
	MarketplaceFileType *MarketplaceFileType `json:"marketplace_file_type"`
	MimeType            *string              `json:"mime_type"`
	DisplayName         *string              `json:"display_name"`
}

type CreateFile_Response struct {
	UploadUrl *string   `json:"upload_url"`
	FileInfo  *FileInfo `json:"file_info"`
}

type CreateListing struct {
	Listing *Listing `json:"listing"`
}

type CreateListing_Response struct {
	ListingId *string `json:"listing_id"`
}

// Data request messages also creates a lead (maybe).
type CreatePersonalizationRequest struct {
	ListingId             *string                    `json:"listing_id"`
	Comment               *string                    `json:"comment"`
	IntendedUse           *string                    `json:"intended_use"`
	FirstName             *string                    `json:"first_name"`
	LastName              *string                    `json:"last_name"`
	Company               *string                    `json:"company"`
	IsFromLighthouse      *bool                      `json:"is_from_lighthouse"`
	RecipientType         *DeltaSharingRecipientType `json:"recipient_type"`
	AcceptedConsumerTerms *ConsumerTerms             `json:"accepted_consumer_terms"`
}

type CreatePersonalizationRequest_Response struct {
	Id *string `json:"id"`
}

type CreateProvider struct {
	Provider *ProviderInfo `json:"provider"`
}

type CreateProvider_Response struct {
	Id *string `json:"id"`
}

type CreateProviderAnalyticsDashboard struct {
}

type CreateProviderAnalyticsDashboard_Response struct {
	Id *string `json:"id"`
}

type DataRefreshInfo struct {
	Interval *int64       `json:"interval"`
	Unit     *DataRefresh `json:"unit"`
}

type DeleteExchangeFilterRequest struct {
	Id *string `json:"id"`
}

type DeleteExchangeFilterResponse struct {
}

type DeleteExchangeRequest struct {
	Id *string `json:"id"`
}

type DeleteExchangeResponse struct {
}

type DeleteFile struct {
	FileId *string `json:"file_id"`
}

type DeleteFile_Response struct {
}

type DeleteListing struct {
	Id *string `json:"id"`
}

type DeleteListing_Response struct {
}

type DeleteProvider struct {
	Id *string `json:"id"`
}

type DeleteProvider_Response struct {
}

type Exchange struct {
	Id             *string           `json:"id"`
	Name           *string           `json:"name"`
	Comment        *string           `json:"comment"`
	Filters        []ExchangeFilter  `json:"filters"`
	CreatedAt      *int64            `json:"created_at"`
	CreatedBy      *string           `json:"created_by"`
	UpdatedAt      *int64            `json:"updated_at"`
	UpdatedBy      *string           `json:"updated_by"`
	LinkedListings []ExchangeListing `json:"linked_listings"`
}

type ExchangeFilter struct {
	Id          *string             `json:"id"`
	ExchangeId  *string             `json:"exchange_id"`
	FilterValue *string             `json:"filter_value"`
	Name        *string             `json:"name"`
	CreatedAt   *int64              `json:"created_at"`
	CreatedBy   *string             `json:"created_by"`
	UpdatedAt   *int64              `json:"updated_at"`
	UpdatedBy   *string             `json:"updated_by"`
	FilterType  *ExchangeFilterType `json:"filter_type"`
}

type ExchangeListing struct {
	Id           *string `json:"id"`
	ExchangeId   *string `json:"exchange_id"`
	ExchangeName *string `json:"exchange_name"`
	ListingId    *string `json:"listing_id"`
	ListingName  *string `json:"listing_name"`
	CreatedAt    *int64  `json:"created_at"`
	CreatedBy    *string `json:"created_by"`
}

type FileInfo struct {
	Id                  *string              `json:"id"`
	MarketplaceFileType *MarketplaceFileType `json:"marketplace_file_type"`
	FileParent          *FileParent          `json:"file_parent"`
	MimeType            *string              `json:"mime_type"`
	DownloadLink        *string              `json:"download_link"`
	CreatedAt           *int64               `json:"created_at"`
	UpdatedAt           *int64               `json:"updated_at"`
	DisplayName         *string              `json:"display_name"`
	Status              *FileStatus          `json:"status"`
	StatusMessage       *string              `json:"status_message"`
}

type FileParent struct {
	ParentId       *string         `json:"parent_id"`
	FileParentType *FileParentType `json:"file_parent_type"`
}

type GetAllInstallations struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type GetAllInstallations_Response struct {
	Installations []InstallationDetail `json:"installations"`
	NextPageToken *string              `json:"next_page_token"`
}

type GetAllPersonalizationRequestsForConsumer struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type GetAllPersonalizationRequestsForConsumer_Response struct {
	PersonalizationRequests []PersonalizationRequest `json:"personalization_requests"`
	NextPageToken           *string                  `json:"next_page_token"`
}

type GetExchangeRequest struct {
	Id *string `json:"id"`
}

type GetExchangeResponse struct {
	Exchange *Exchange `json:"exchange"`
}

type GetFile struct {
	FileId *string `json:"file_id"`
}

type GetFile_Response struct {
	FileInfo *FileInfo `json:"file_info"`
}

type GetInstallationDetails struct {
	ListingId *string `json:"listing_id"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type GetInstallationDetails_Response struct {
	Installations []InstallationDetail `json:"installations"`
	NextPageToken *string              `json:"next_page_token"`
}

// this is effectively a static request for now and will return latest version
// of the dashboard template that exists on server..
type GetLatestVersionProviderAnalyticsDashboard struct {
}

type GetLatestVersionProviderAnalyticsDashboard_Response struct {
	Version *int64 `json:"version"`
}

type GetListing struct {
	Id *string `json:"id"`
}

type GetListing_Response struct {
	Listing *Listing `json:"listing"`
}

type GetListingContent struct {
	ListingId *string `json:"listing_id"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type GetListingContent_Response struct {
	SharedDataObjects []SharedDataObject `json:"shared_data_objects"`
	NextPageToken     *string            `json:"next_page_token"`
}

type GetListingFulfillments struct {
	ListingId *string `json:"listing_id"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type GetListingFulfillments_Response struct {
	Fulfillments  []ListingFulfillment `json:"fulfillments"`
	NextPageToken *string              `json:"next_page_token"`
}

type GetListings struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type GetListings_Response struct {
	Listings      []Listing `json:"listings"`
	NextPageToken *string   `json:"next_page_token"`
}

type GetPersonalizationRequestsForConsumer struct {
	ListingId *string `json:"listing_id"`
}

type GetPersonalizationRequestsForConsumer_Response struct {
	PersonalizationRequests []PersonalizationRequest `json:"personalization_requests"`
}

type GetPersonalizationRequestsForProvider struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type GetPersonalizationRequestsForProvider_Response struct {
	PersonalizationRequests []PersonalizationRequest `json:"personalization_requests"`
	NextPageToken           *string                  `json:"next_page_token"`
}

type GetProvider struct {
	Id *string `json:"id"`
}

type GetProvider_Response struct {
	Provider *ProviderInfo `json:"provider"`
}

type GetPublishedListingForConsumer struct {
	Id *string `json:"id"`
}

type GetPublishedListingForConsumer_Response struct {
	Listing *Listing `json:"listing"`
}

// Listing messages.
type GetPublishedListingsForConsumer struct {
	PageToken         *string      `json:"page_token"`
	PageSize          *int         `json:"page_size"`
	Assets            []AssetType  `json:"assets"`
	Categories        []Category   `json:"categories"`
	Tags              []ListingTag `json:"tags"`
	IsFree            *bool        `json:"is_free"`
	IsPrivateExchange *bool        `json:"is_private_exchange"`
	IsStaffPick       *bool        `json:"is_staff_pick"`
	ProviderIds       []string     `json:"provider_ids"`
}

type GetPublishedListingsForConsumer_Response struct {
	Listings      []Listing `json:"listings"`
	NextPageToken *string   `json:"next_page_token"`
}

type GetPublishedProviderForConsumer struct {
	Id *string `json:"id"`
}

type GetPublishedProviderForConsumer_Response struct {
	Provider *ProviderInfo `json:"provider"`
}

type InstallListing struct {
	ListingId             *string                    `json:"listing_id"`
	ShareName             *string                    `json:"share_name"`
	CatalogName           *string                    `json:"catalog_name"`
	RepoDetail            *RepoInstallation          `json:"repo_detail"`
	RecipientType         *DeltaSharingRecipientType `json:"recipient_type"`
	AcceptedConsumerTerms *ConsumerTerms             `json:"accepted_consumer_terms"`
}

type InstallListing_Response struct {
	Installation *InstallationDetail `json:"installation"`
}

type InstallationDetail struct {
	Id            *string                    `json:"id"`
	ListingId     *string                    `json:"listing_id"`
	ShareName     *string                    `json:"share_name"`
	CatalogName   *string                    `json:"catalog_name"`
	InstalledOn   *int64                     `json:"installed_on"`
	Status        *InstallationStatus        `json:"status"`
	ErrorMessage  *string                    `json:"error_message"`
	ListingName   *string                    `json:"listing_name"`
	RepoName      *string                    `json:"repo_name"`
	RepoPath      *string                    `json:"repo_path"`
	RecipientType *DeltaSharingRecipientType `json:"recipient_type"`
	Tokens        []TokenInfo                `json:"tokens"`
	TokenDetail   *TokenDetail               `json:"token_detail"`
}

type ListExchangeFiltersRequest struct {
	ExchangeId *string `json:"exchange_id"`
	PageToken  *string `json:"page_token"`
	PageSize   *int    `json:"page_size"`
}

type ListExchangeFiltersResponse struct {
	Filters       []ExchangeFilter `json:"filters"`
	NextPageToken *string          `json:"next_page_token"`
}

type ListExchangesForListingRequest struct {
	ListingId *string `json:"listing_id"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListExchangesForListingResponse struct {
	ExchangeListing []ExchangeListing `json:"exchange_listing"`
	NextPageToken   *string           `json:"next_page_token"`
}

type ListExchangesRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListExchangesResponse struct {
	Exchanges     []Exchange `json:"exchanges"`
	NextPageToken *string    `json:"next_page_token"`
}

type ListFiles struct {
	FileParent *FileParent `json:"file_parent"`
	PageToken  *string     `json:"page_token"`
	PageSize   *int        `json:"page_size"`
}

type ListFiles_Response struct {
	FileInfos     []FileInfo `json:"file_infos"`
	NextPageToken *string    `json:"next_page_token"`
}

type ListListingsForExchangeRequest struct {
	ExchangeId *string `json:"exchange_id"`
	PageToken  *string `json:"page_token"`
	PageSize   *int    `json:"page_size"`
}

type ListListingsForExchangeResponse struct {
	ExchangeListings []ExchangeListing `json:"exchange_listings"`
	NextPageToken    *string           `json:"next_page_token"`
}

type ListProviderAnalyticsDashboard struct {
}

type ListProviderAnalyticsDashboard_Response struct {
	Id          *string `json:"id"`
	Version     *int64  `json:"version"`
	DashboardId *string `json:"dashboard_id"`
}

type ListProviders struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListProviders_Response struct {
	Providers     []ProviderInfo `json:"providers"`
	NextPageToken *string        `json:"next_page_token"`
}

type ListPublishedProvidersForConsumer struct {
	PageToken  *string `json:"page_token"`
	PageSize   *int    `json:"page_size"`
	IsFeatured *bool   `json:"is_featured"`
}

type ListPublishedProvidersForConsumer_Response struct {
	Providers     []ProviderInfo `json:"providers"`
	NextPageToken *string        `json:"next_page_token"`
}

type Listing struct {
	Id      *string         `json:"id"`
	Summary *ListingSummary `json:"summary"`
	Detail  *ListingDetail  `json:"detail"`
}

type ListingDetail struct {
	Description               *string          `json:"description"`
	TermsOfService            *string          `json:"terms_of_service"`
	DocumentationLink         *string          `json:"documentation_link"`
	SupportLink               *string          `json:"support_link"`
	FileIds                   []string         `json:"file_ids"`
	PrivacyPolicyLink         *string          `json:"privacy_policy_link"`
	EmbeddedNotebookFileInfos []FileInfo       `json:"embedded_notebook_file_infos"`
	GeographicalCoverage      *string          `json:"geographical_coverage"`
	Cost                      *Cost            `json:"cost"`
	PricingModel              *string          `json:"pricing_model"`
	UpdateFrequency           *DataRefreshInfo `json:"update_frequency"`
	CollectionGranularity     *DataRefreshInfo `json:"collection_granularity"`
	CollectionDateStart       *int64           `json:"collection_date_start"`
	CollectionDateEnd         *int64           `json:"collection_date_end"`
	DataSource                *string          `json:"data_source"`
	Size                      *float64         `json:"size"`
	Assets                    []AssetType      `json:"assets"`
	License                   *string          `json:"license"`
	Tags                      []ListingTag     `json:"tags"`
}

type ListingFulfillment struct {
	ListingId       *string                    `json:"listing_id"`
	FulfillmentType *FulfillmentType           `json:"fulfillment_type"`
	ShareInfo       *ShareInfo                 `json:"share_info"`
	RepoInfo        *RepoInfo                  `json:"repo_info"`
	RecipientType   *DeltaSharingRecipientType `json:"recipient_type"`
}

type ListingSetting struct {
	Visibility *Visibility `json:"visibility"`
}

type ListingSummary struct {
	Name           *string         `json:"name"`
	Subtitle       *string         `json:"subtitle"`
	Status         *ListingStatus  `json:"status"`
	Share          *ShareInfo      `json:"share"`
	ProviderRegion *RegionInfo     `json:"provider_region"`
	Setting        *ListingSetting `json:"setting"`
	CreatedAt      *int64          `json:"created_at"`
	CreatedBy      *string         `json:"created_by"`
	UpdatedAt      *int64          `json:"updated_at"`
	UpdatedBy      *string         `json:"updated_by"`
	PublishedAt    *int64          `json:"published_at"`
	PublishedBy    *string         `json:"published_by"`
	Categories     []Category      `json:"categories"`
	ListingType    *ListingType    `json:"listingType"`
	CreatedById    *int64          `json:"created_by_id"`
	UpdatedById    *int64          `json:"updated_by_id"`
	ProviderId     *string         `json:"provider_id"`
	ExchangeIds    []string        `json:"exchange_ids"`
	GitRepo        *RepoInfo       `json:"git_repo"`
}

type ListingTag struct {
	TagName   *ListingTagType `json:"tag_name"`
	TagValues []string        `json:"tag_values"`
}

type PersonalizationRequest struct {
	Id               *string                       `json:"id"`
	ConsumerRegion   *RegionInfo                   `json:"consumer_region"`
	ContactInfo      *ContactInfo                  `json:"contact_info"`
	Comment          *string                       `json:"comment"`
	IntendedUse      *string                       `json:"intended_use"`
	Status           *PersonalizationRequestStatus `json:"status"`
	StatusMessage    *string                       `json:"status_message"`
	Share            *ShareInfo                    `json:"share"`
	CreatedAt        *int64                        `json:"created_at"`
	ListingId        *string                       `json:"listing_id"`
	UpdatedAt        *int64                        `json:"updated_at"`
	MetastoreId      *string                       `json:"metastore_id"`
	ListingName      *string                       `json:"listing_name"`
	IsFromLighthouse *bool                         `json:"is_from_lighthouse"`
	ProviderId       *string                       `json:"provider_id"`
	RecipientType    *DeltaSharingRecipientType    `json:"recipient_type"`
}

type ProviderInfo struct {
	Id                   *string `json:"id"`
	Name                 *string `json:"name"`
	Description          *string `json:"description"`
	IconFilePath         *string `json:"icon_file_path"`
	BusinessContactEmail *string `json:"business_contact_email"`
	SupportContactEmail  *string `json:"support_contact_email"`
	IsFeatured           *bool   `json:"is_featured"`
	PublishedBy          *string `json:"published_by"`
	CompanyWebsiteLink   *string `json:"company_website_link"`
	IconFileId           *string `json:"icon_file_id"`
	TermOfServiceLink    *string `json:"term_of_service_link"`
	PrivacyPolicyLink    *string `json:"privacy_policy_link"`
	DarkModeIconFileId   *string `json:"dark_mode_icon_file_id"`
	DarkModeIconFilePath *string `json:"dark_mode_icon_file_path"`
}

type RegionInfo struct {
	Cloud  *string `json:"cloud"`
	Region *string `json:"region"`
}

type RemoveExchangeForListingRequest struct {
	Id *string `json:"id"`
}

type RemoveExchangeForListingResponse struct {
}

type RepoInfo struct {
	GitRepoUrl *string `json:"git_repo_url"`
}

type RepoInstallation struct {
	RepoName *string `json:"repo_name"`
	RepoPath *string `json:"repo_path"`
}

type SearchPublishedListingsForConsumer struct {
	Query             *string     `json:"query"`
	IsFree            *bool       `json:"is_free"`
	IsPrivateExchange *bool       `json:"is_private_exchange"`
	ProviderIds       []string    `json:"provider_ids"`
	Categories        []Category  `json:"categories"`
	Assets            []AssetType `json:"assets"`
	PageToken         *string     `json:"page_token"`
	PageSize          *int        `json:"page_size"`
}

type SearchPublishedListingsForConsumer_Response struct {
	Listings      []Listing `json:"listings"`
	NextPageToken *string   `json:"next_page_token"`
}

type ShareInfo struct {
	Name *string           `json:"name"`
	Type *ListingShareType `json:"type"`
}

type SharedDataObject struct {
	Name           *string `json:"name"`
	DataObjectType *string `json:"data_object_type"`
}

type TokenDetail struct {
	ShareCredentialsVersion *int    `json:"shareCredentialsVersion"`
	BearerToken             *string `json:"bearerToken"`
	Endpoint                *string `json:"endpoint"`
	ExpirationTime          *string `json:"expirationTime"`
}

type TokenInfo struct {
	Id             *string `json:"id"`
	CreatedAt      *int64  `json:"created_at"`
	CreatedBy      *string `json:"created_by"`
	ActivationUrl  *string `json:"activation_url"`
	ExpirationTime *int64  `json:"expiration_time"`
	UpdatedAt      *int64  `json:"updated_at"`
	UpdatedBy      *string `json:"updated_by"`
}

type UninstallListing struct {
	ListingId      *string `json:"listing_id"`
	InstallationId *string `json:"installation_id"`
}

type UninstallListing_Response struct {
}

type UpdateExchangeFilterRequest struct {
	Id     *string         `json:"id"`
	Filter *ExchangeFilter `json:"filter"`
}

type UpdateExchangeFilterResponse struct {
	Filter *ExchangeFilter `json:"filter"`
}

type UpdateExchangeRequest struct {
	Id       *string   `json:"id"`
	Exchange *Exchange `json:"exchange"`
}

type UpdateExchangeResponse struct {
	Exchange *Exchange `json:"exchange"`
}

type UpdateInstallationDetail struct {
	ListingId      *string             `json:"listing_id"`
	InstallationId *string             `json:"installation_id"`
	Installation   *InstallationDetail `json:"installation"`
	RotateToken    *bool               `json:"rotate_token"`
}

type UpdateInstallationDetail_Response struct {
	Installation *InstallationDetail `json:"installation"`
}

type UpdateListing struct {
	Id      *string  `json:"id"`
	Listing *Listing `json:"listing"`
}

type UpdateListing_Response struct {
	Listing *Listing `json:"listing"`
}

type UpdatePersonalizationRequestStatus struct {
	ListingId *string                       `json:"listing_id"`
	RequestId *string                       `json:"request_id"`
	Status    *PersonalizationRequestStatus `json:"status"`
	Reason    *string                       `json:"reason"`
	Share     *ShareInfo                    `json:"share"`
}

type UpdatePersonalizationRequestStatus_Response struct {
	Request *PersonalizationRequest `json:"request"`
}

type UpdateProvider struct {
	Id       *string       `json:"id"`
	Provider *ProviderInfo `json:"provider"`
}

type UpdateProvider_Response struct {
	Provider *ProviderInfo `json:"provider"`
}

type UpdateProviderAnalyticsDashboard struct {
	Id      *string `json:"id"`
	Version *int64  `json:"version"`
}

type UpdateProviderAnalyticsDashboard_Response struct {
	Id          *string `json:"id"`
	Version     *int64  `json:"version"`
	DashboardId *string `json:"dashboard_id"`
}
