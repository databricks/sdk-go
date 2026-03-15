package v2

type AwsAvailability string

const (
	AwsAvailability_Spot             AwsAvailability = "SPOT"
	AwsAvailability_OnDemand         AwsAvailability = "ON_DEMAND"
	AwsAvailability_SpotWithFallback AwsAvailability = "SPOT_WITH_FALLBACK"
)

// String representation for [fmt.Print].
func (f *AwsAvailability) String() string {
	return string(*f)
}

type AzureAvailability string

const (
	AzureAvailability_SpotAzure             AzureAvailability = "SPOT_AZURE"
	AzureAvailability_OnDemandAzure         AzureAvailability = "ON_DEMAND_AZURE"
	AzureAvailability_SpotWithFallbackAzure AzureAvailability = "SPOT_WITH_FALLBACK_AZURE"
)

// String representation for [fmt.Print].
func (f *AzureAvailability) String() string {
	return string(*f)
}

type CloudProviderNodeStatus string

const (
	CloudProviderNodeStatus_NotEnabledOnSubscription CloudProviderNodeStatus = "NotEnabledOnSubscription"
	CloudProviderNodeStatus_NotAvailableInRegion     CloudProviderNodeStatus = "NotAvailableInRegion"
)

// String representation for [fmt.Print].
func (f *CloudProviderNodeStatus) String() string {
	return string(*f)
}

type ComputeKind string

const (
	ComputeKind_ComputeKindUnspecified ComputeKind = "COMPUTE_KIND_UNSPECIFIED"
	ComputeKind_ClassicPreview         ComputeKind = "CLASSIC_PREVIEW"
)

// String representation for [fmt.Print].
func (f *ComputeKind) String() string {
	return string(*f)
}

type ConfidentialComputeType string

const (
	ConfidentialComputeType_SevSnp ConfidentialComputeType = "SEV_SNP"
)

// String representation for [fmt.Print].
func (f *ConfidentialComputeType) String() string {
	return string(*f)
}

type DataSecurityMode string

const (
	DataSecurityMode_None                      DataSecurityMode = "NONE"
	DataSecurityMode_SingleUser                DataSecurityMode = "SINGLE_USER"
	DataSecurityMode_UserIsolation             DataSecurityMode = "USER_ISOLATION"
	DataSecurityMode_LegacyTableAcl            DataSecurityMode = "LEGACY_TABLE_ACL"
	DataSecurityMode_LegacyPassthrough         DataSecurityMode = "LEGACY_PASSTHROUGH"
	DataSecurityMode_LegacySingleUser          DataSecurityMode = "LEGACY_SINGLE_USER"
	DataSecurityMode_LegacySingleUserStandard  DataSecurityMode = "LEGACY_SINGLE_USER_STANDARD"
	DataSecurityMode_DataSecurityModeStandard  DataSecurityMode = "DATA_SECURITY_MODE_STANDARD"
	DataSecurityMode_DataSecurityModeDedicated DataSecurityMode = "DATA_SECURITY_MODE_DEDICATED"
	DataSecurityMode_DataSecurityModeAuto      DataSecurityMode = "DATA_SECURITY_MODE_AUTO"
)

// String representation for [fmt.Print].
func (f *DataSecurityMode) String() string {
	return string(*f)
}

type EbsVolumeType string

const (
	EbsVolumeType_GeneralPurposeSsd      EbsVolumeType = "GENERAL_PURPOSE_SSD"
	EbsVolumeType_ThroughputOptimizedHdd EbsVolumeType = "THROUGHPUT_OPTIMIZED_HDD"
)

// String representation for [fmt.Print].
func (f *EbsVolumeType) String() string {
	return string(*f)
}

type GcpAvailability string

const (
	GcpAvailability_PreemptibleGcp             GcpAvailability = "PREEMPTIBLE_GCP"
	GcpAvailability_OnDemandGcp                GcpAvailability = "ON_DEMAND_GCP"
	GcpAvailability_PreemptibleWithFallbackGcp GcpAvailability = "PREEMPTIBLE_WITH_FALLBACK_GCP"
)

// String representation for [fmt.Print].
func (f *GcpAvailability) String() string {
	return string(*f)
}

type RuntimeEngine string

const (
	RuntimeEngine_Null     RuntimeEngine = "NULL"
	RuntimeEngine_Standard RuntimeEngine = "STANDARD"
	RuntimeEngine_Photon   RuntimeEngine = "PHOTON"
)

// String representation for [fmt.Print].
func (f *RuntimeEngine) String() string {
	return string(*f)
}

type TerminationCode string

const (
	TerminationCode_Unknown                                             TerminationCode = "UNKNOWN"
	TerminationCode_UserRequest                                         TerminationCode = "USER_REQUEST"
	TerminationCode_JobFinished                                         TerminationCode = "JOB_FINISHED"
	TerminationCode_Inactivity                                          TerminationCode = "INACTIVITY"
	TerminationCode_CloudProviderShutdown                               TerminationCode = "CLOUD_PROVIDER_SHUTDOWN"
	TerminationCode_CommunicationLost                                   TerminationCode = "COMMUNICATION_LOST"
	TerminationCode_CloudProviderLaunchFailure                          TerminationCode = "CLOUD_PROVIDER_LAUNCH_FAILURE"
	TerminationCode_InitScriptFailure                                   TerminationCode = "INIT_SCRIPT_FAILURE"
	TerminationCode_SparkStartupFailure                                 TerminationCode = "SPARK_STARTUP_FAILURE"
	TerminationCode_InvalidArgument                                     TerminationCode = "INVALID_ARGUMENT"
	TerminationCode_UnexpectedLaunchFailure                             TerminationCode = "UNEXPECTED_LAUNCH_FAILURE"
	TerminationCode_InternalError                                       TerminationCode = "INTERNAL_ERROR"
	TerminationCode_InstanceUnreachable                                 TerminationCode = "INSTANCE_UNREACHABLE"
	TerminationCode_RequestRejected                                     TerminationCode = "REQUEST_REJECTED"
	TerminationCode_TrialExpired                                        TerminationCode = "TRIAL_EXPIRED"
	TerminationCode_DriverUnreachable                                   TerminationCode = "DRIVER_UNREACHABLE"
	TerminationCode_SparkError                                          TerminationCode = "SPARK_ERROR"
	TerminationCode_DriverUnresponsive                                  TerminationCode = "DRIVER_UNRESPONSIVE"
	TerminationCode_MetastoreComponentUnhealthy                         TerminationCode = "METASTORE_COMPONENT_UNHEALTHY"
	TerminationCode_DbfsComponentUnhealthy                              TerminationCode = "DBFS_COMPONENT_UNHEALTHY"
	TerminationCode_ExecutionComponentUnhealthy                         TerminationCode = "EXECUTION_COMPONENT_UNHEALTHY"
	TerminationCode_AzureResourceManagerThrottling                      TerminationCode = "AZURE_RESOURCE_MANAGER_THROTTLING"
	TerminationCode_AzureResourceProviderThrottling                     TerminationCode = "AZURE_RESOURCE_PROVIDER_THROTTLING"
	TerminationCode_NetworkConfigurationFailure                         TerminationCode = "NETWORK_CONFIGURATION_FAILURE"
	TerminationCode_ContainerLaunchFailure                              TerminationCode = "CONTAINER_LAUNCH_FAILURE"
	TerminationCode_InstancePoolClusterFailure                          TerminationCode = "INSTANCE_POOL_CLUSTER_FAILURE"
	TerminationCode_SkippedSlowNodes                                    TerminationCode = "SKIPPED_SLOW_NODES"
	TerminationCode_AttachProjectFailure                                TerminationCode = "ATTACH_PROJECT_FAILURE"
	TerminationCode_UpdateInstanceProfileFailure                        TerminationCode = "UPDATE_INSTANCE_PROFILE_FAILURE"
	TerminationCode_DatabaseConnectionFailure                           TerminationCode = "DATABASE_CONNECTION_FAILURE"
	TerminationCode_RequestThrottled                                    TerminationCode = "REQUEST_THROTTLED"
	TerminationCode_SelfBootstrapFailure                                TerminationCode = "SELF_BOOTSTRAP_FAILURE"
	TerminationCode_GlobalInitScriptFailure                             TerminationCode = "GLOBAL_INIT_SCRIPT_FAILURE"
	TerminationCode_SlowImageDownload                                   TerminationCode = "SLOW_IMAGE_DOWNLOAD"
	TerminationCode_InvalidSparkImage                                   TerminationCode = "INVALID_SPARK_IMAGE"
	TerminationCode_NpipTunnelTokenFailure                              TerminationCode = "NPIP_TUNNEL_TOKEN_FAILURE"
	TerminationCode_HiveMetastoreProvisioningFailure                    TerminationCode = "HIVE_METASTORE_PROVISIONING_FAILURE"
	TerminationCode_AzureInvalidDeploymentTemplate                      TerminationCode = "AZURE_INVALID_DEPLOYMENT_TEMPLATE"
	TerminationCode_AzureUnexpectedDeploymentTemplateFailure            TerminationCode = "AZURE_UNEXPECTED_DEPLOYMENT_TEMPLATE_FAILURE"
	TerminationCode_SubnetExhaustedFailure                              TerminationCode = "SUBNET_EXHAUSTED_FAILURE"
	TerminationCode_BootstrapTimeout                                    TerminationCode = "BOOTSTRAP_TIMEOUT"
	TerminationCode_StorageDownloadFailure                              TerminationCode = "STORAGE_DOWNLOAD_FAILURE"
	TerminationCode_ControlPlaneRequestFailure                          TerminationCode = "CONTROL_PLANE_REQUEST_FAILURE"
	TerminationCode_BootstrapTimeoutCloudProviderException              TerminationCode = "BOOTSTRAP_TIMEOUT_CLOUD_PROVIDER_EXCEPTION"
	TerminationCode_AwsInsufficientInstanceCapacityFailure              TerminationCode = "AWS_INSUFFICIENT_INSTANCE_CAPACITY_FAILURE"
	TerminationCode_DockerImagePullFailure                              TerminationCode = "DOCKER_IMAGE_PULL_FAILURE"
	TerminationCode_AzureVnetConfigurationFailure                       TerminationCode = "AZURE_VNET_CONFIGURATION_FAILURE"
	TerminationCode_NpipTunnelSetupFailure                              TerminationCode = "NPIP_TUNNEL_SETUP_FAILURE"
	TerminationCode_AwsAuthorizationFailure                             TerminationCode = "AWS_AUTHORIZATION_FAILURE"
	TerminationCode_NephosResourceManagement                            TerminationCode = "NEPHOS_RESOURCE_MANAGEMENT"
	TerminationCode_StsClientSetupFailure                               TerminationCode = "STS_CLIENT_SETUP_FAILURE"
	TerminationCode_SecurityDaemonRegistrationException                 TerminationCode = "SECURITY_DAEMON_REGISTRATION_EXCEPTION"
	TerminationCode_AwsRequestLimitExceeded                             TerminationCode = "AWS_REQUEST_LIMIT_EXCEEDED"
	TerminationCode_AwsInsufficientFreeAddressesInSubnetFailure         TerminationCode = "AWS_INSUFFICIENT_FREE_ADDRESSES_IN_SUBNET_FAILURE"
	TerminationCode_AwsUnsupportedFailure                               TerminationCode = "AWS_UNSUPPORTED_FAILURE"
	TerminationCode_AzureQuotaExceededException                         TerminationCode = "AZURE_QUOTA_EXCEEDED_EXCEPTION"
	TerminationCode_AzureOperationNotAllowedException                   TerminationCode = "AZURE_OPERATION_NOT_ALLOWED_EXCEPTION"
	TerminationCode_NfsMountFailure                                     TerminationCode = "NFS_MOUNT_FAILURE"
	TerminationCode_K8sAutoscalingFailure                               TerminationCode = "K8S_AUTOSCALING_FAILURE"
	TerminationCode_K8sDbrClusterLaunchTimeout                          TerminationCode = "K8S_DBR_CLUSTER_LAUNCH_TIMEOUT"
	TerminationCode_SparkImageDownloadFailure                           TerminationCode = "SPARK_IMAGE_DOWNLOAD_FAILURE"
	TerminationCode_AzureVmExtensionFailure                             TerminationCode = "AZURE_VM_EXTENSION_FAILURE"
	TerminationCode_WorkspaceCancelledError                             TerminationCode = "WORKSPACE_CANCELLED_ERROR"
	TerminationCode_AwsMaxSpotInstanceCountExceededFailure              TerminationCode = "AWS_MAX_SPOT_INSTANCE_COUNT_EXCEEDED_FAILURE"
	TerminationCode_TemporarilyUnavailable                              TerminationCode = "TEMPORARILY_UNAVAILABLE"
	TerminationCode_WorkerSetupFailure                                  TerminationCode = "WORKER_SETUP_FAILURE"
	TerminationCode_IpExhaustionFailure                                 TerminationCode = "IP_EXHAUSTION_FAILURE"
	TerminationCode_GcpQuotaExceeded                                    TerminationCode = "GCP_QUOTA_EXCEEDED"
	TerminationCode_CloudProviderResourceStockout                       TerminationCode = "CLOUD_PROVIDER_RESOURCE_STOCKOUT"
	TerminationCode_GcpServiceAccountDeleted                            TerminationCode = "GCP_SERVICE_ACCOUNT_DELETED"
	TerminationCode_AzureByokKeyPermissionFailure                       TerminationCode = "AZURE_BYOK_KEY_PERMISSION_FAILURE"
	TerminationCode_SpotInstanceTermination                             TerminationCode = "SPOT_INSTANCE_TERMINATION"
	TerminationCode_AzureEphemeralDiskFailure                           TerminationCode = "AZURE_EPHEMERAL_DISK_FAILURE"
	TerminationCode_AbuseDetected                                       TerminationCode = "ABUSE_DETECTED"
	TerminationCode_ImagePullPermissionDenied                           TerminationCode = "IMAGE_PULL_PERMISSION_DENIED"
	TerminationCode_WorkspaceConfigurationError                         TerminationCode = "WORKSPACE_CONFIGURATION_ERROR"
	TerminationCode_SecretResolutionError                               TerminationCode = "SECRET_RESOLUTION_ERROR"
	TerminationCode_UnsupportedInstanceType                             TerminationCode = "UNSUPPORTED_INSTANCE_TYPE"
	TerminationCode_CloudProviderDiskSetupFailure                       TerminationCode = "CLOUD_PROVIDER_DISK_SETUP_FAILURE"
	TerminationCode_SshBootstrapFailure                                 TerminationCode = "SSH_BOOTSTRAP_FAILURE"
	TerminationCode_AwsInaccessibleKmsKeyFailure                        TerminationCode = "AWS_INACCESSIBLE_KMS_KEY_FAILURE"
	TerminationCode_InitContainerNotFinished                            TerminationCode = "INIT_CONTAINER_NOT_FINISHED"
	TerminationCode_SparkImageDownloadThrottled                         TerminationCode = "SPARK_IMAGE_DOWNLOAD_THROTTLED"
	TerminationCode_SparkImageNotFound                                  TerminationCode = "SPARK_IMAGE_NOT_FOUND"
	TerminationCode_ClusterOperationThrottled                           TerminationCode = "CLUSTER_OPERATION_THROTTLED"
	TerminationCode_ClusterOperationTimeout                             TerminationCode = "CLUSTER_OPERATION_TIMEOUT"
	TerminationCode_ServerlessLongRunningTerminated                     TerminationCode = "SERVERLESS_LONG_RUNNING_TERMINATED"
	TerminationCode_AzurePackedDeploymentPartialFailure                 TerminationCode = "AZURE_PACKED_DEPLOYMENT_PARTIAL_FAILURE"
	TerminationCode_InvalidWorkerImageFailure                           TerminationCode = "INVALID_WORKER_IMAGE_FAILURE"
	TerminationCode_WorkspaceUpdate                                     TerminationCode = "WORKSPACE_UPDATE"
	TerminationCode_InvalidAwsParameter                                 TerminationCode = "INVALID_AWS_PARAMETER"
	TerminationCode_DriverOutOfDisk                                     TerminationCode = "DRIVER_OUT_OF_DISK"
	TerminationCode_DriverOutOfMemory                                   TerminationCode = "DRIVER_OUT_OF_MEMORY"
	TerminationCode_DriverLaunchTimeout                                 TerminationCode = "DRIVER_LAUNCH_TIMEOUT"
	TerminationCode_DriverUnexpectedFailure                             TerminationCode = "DRIVER_UNEXPECTED_FAILURE"
	TerminationCode_UnexpectedPodRecreation                             TerminationCode = "UNEXPECTED_POD_RECREATION"
	TerminationCode_GcpInaccessibleKmsKeyFailure                        TerminationCode = "GCP_INACCESSIBLE_KMS_KEY_FAILURE"
	TerminationCode_GcpKmsKeyPermissionDenied                           TerminationCode = "GCP_KMS_KEY_PERMISSION_DENIED"
	TerminationCode_DriverEviction                                      TerminationCode = "DRIVER_EVICTION"
	TerminationCode_UserInitiatedVmTermination                          TerminationCode = "USER_INITIATED_VM_TERMINATION"
	TerminationCode_GcpIamTimeout                                       TerminationCode = "GCP_IAM_TIMEOUT"
	TerminationCode_AwsResourceQuotaExceeded                            TerminationCode = "AWS_RESOURCE_QUOTA_EXCEEDED"
	TerminationCode_CloudAccountSetupFailure                            TerminationCode = "CLOUD_ACCOUNT_SETUP_FAILURE"
	TerminationCode_AwsInvalidKeyPair                                   TerminationCode = "AWS_INVALID_KEY_PAIR"
	TerminationCode_DriverPodCreationFailure                            TerminationCode = "DRIVER_POD_CREATION_FAILURE"
	TerminationCode_MaintenanceMode                                     TerminationCode = "MAINTENANCE_MODE"
	TerminationCode_InternalCapacityFailure                             TerminationCode = "INTERNAL_CAPACITY_FAILURE"
	TerminationCode_ExecutorPodUnscheduled                              TerminationCode = "EXECUTOR_POD_UNSCHEDULED"
	TerminationCode_StorageDownloadFailureSlow                          TerminationCode = "STORAGE_DOWNLOAD_FAILURE_SLOW"
	TerminationCode_StorageDownloadFailureThrottled                     TerminationCode = "STORAGE_DOWNLOAD_FAILURE_THROTTLED"
	TerminationCode_DynamicSparkConfSizeExceeded                        TerminationCode = "DYNAMIC_SPARK_CONF_SIZE_EXCEEDED"
	TerminationCode_AwsInstanceProfileUpdateFailure                     TerminationCode = "AWS_INSTANCE_PROFILE_UPDATE_FAILURE"
	TerminationCode_InstancePoolNotFound                                TerminationCode = "INSTANCE_POOL_NOT_FOUND"
	TerminationCode_InstancePoolMaxCapacityReached                      TerminationCode = "INSTANCE_POOL_MAX_CAPACITY_REACHED"
	TerminationCode_AwsInvalidKmsKeyState                               TerminationCode = "AWS_INVALID_KMS_KEY_STATE"
	TerminationCode_GcpInsufficientCapacity                             TerminationCode = "GCP_INSUFFICIENT_CAPACITY"
	TerminationCode_GcpApiRateQuotaExceeded                             TerminationCode = "GCP_API_RATE_QUOTA_EXCEEDED"
	TerminationCode_GcpResourceQuotaExceeded                            TerminationCode = "GCP_RESOURCE_QUOTA_EXCEEDED"
	TerminationCode_GcpIpSpaceExhausted                                 TerminationCode = "GCP_IP_SPACE_EXHAUSTED"
	TerminationCode_GcpServiceAccountAccessDenied                       TerminationCode = "GCP_SERVICE_ACCOUNT_ACCESS_DENIED"
	TerminationCode_GcpServiceAccountNotFound                           TerminationCode = "GCP_SERVICE_ACCOUNT_NOT_FOUND"
	TerminationCode_GcpForbidden                                        TerminationCode = "GCP_FORBIDDEN"
	TerminationCode_GcpNotFound                                         TerminationCode = "GCP_NOT_FOUND"
	TerminationCode_ResourceUsageBlocked                                TerminationCode = "RESOURCE_USAGE_BLOCKED"
	TerminationCode_DataAccessConfigChanged                             TerminationCode = "DATA_ACCESS_CONFIG_CHANGED"
	TerminationCode_AccessTokenFailure                                  TerminationCode = "ACCESS_TOKEN_FAILURE"
	TerminationCode_InvalidInstancePlacementProtocol                    TerminationCode = "INVALID_INSTANCE_PLACEMENT_PROTOCOL"
	TerminationCode_BudgetPolicyResolutionFailure                       TerminationCode = "BUDGET_POLICY_RESOLUTION_FAILURE"
	TerminationCode_InPenaltyBox                                        TerminationCode = "IN_PENALTY_BOX"
	TerminationCode_DisasterRecoveryReplication                         TerminationCode = "DISASTER_RECOVERY_REPLICATION"
	TerminationCode_BootstrapTimeoutDueToMisconfig                      TerminationCode = "BOOTSTRAP_TIMEOUT_DUE_TO_MISCONFIG"
	TerminationCode_InstanceUnreachableDueToMisconfig                   TerminationCode = "INSTANCE_UNREACHABLE_DUE_TO_MISCONFIG"
	TerminationCode_StorageDownloadFailureDueToMisconfig                TerminationCode = "STORAGE_DOWNLOAD_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_ControlPlaneRequestFailureDueToMisconfig            TerminationCode = "CONTROL_PLANE_REQUEST_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_CloudProviderLaunchFailureDueToMisconfig            TerminationCode = "CLOUD_PROVIDER_LAUNCH_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_GcpSubnetNotReady                                   TerminationCode = "GCP_SUBNET_NOT_READY"
	TerminationCode_CloudOperationCancelled                             TerminationCode = "CLOUD_OPERATION_CANCELLED"
	TerminationCode_CloudProviderInstanceNotLaunched                    TerminationCode = "CLOUD_PROVIDER_INSTANCE_NOT_LAUNCHED"
	TerminationCode_GcpTrustedImageProjectsViolated                     TerminationCode = "GCP_TRUSTED_IMAGE_PROJECTS_VIOLATED"
	TerminationCode_BudgetPolicyLimitEnforcementActivated               TerminationCode = "BUDGET_POLICY_LIMIT_ENFORCEMENT_ACTIVATED"
	TerminationCode_EosSparkImage                                       TerminationCode = "EOS_SPARK_IMAGE"
	TerminationCode_NoMatchedK8s                                        TerminationCode = "NO_MATCHED_K8S"
	TerminationCode_LazyAllocationTimeout                               TerminationCode = "LAZY_ALLOCATION_TIMEOUT"
	TerminationCode_DriverNodeUnreachable                               TerminationCode = "DRIVER_NODE_UNREACHABLE"
	TerminationCode_SecretCreationFailure                               TerminationCode = "SECRET_CREATION_FAILURE"
	TerminationCode_PodSchedulingFailure                                TerminationCode = "POD_SCHEDULING_FAILURE"
	TerminationCode_PodAssignmentFailure                                TerminationCode = "POD_ASSIGNMENT_FAILURE"
	TerminationCode_AllocationTimeout                                   TerminationCode = "ALLOCATION_TIMEOUT"
	TerminationCode_AllocationTimeoutNoUnallocatedClusters              TerminationCode = "ALLOCATION_TIMEOUT_NO_UNALLOCATED_CLUSTERS"
	TerminationCode_AllocationTimeoutNoMatchedClusters                  TerminationCode = "ALLOCATION_TIMEOUT_NO_MATCHED_CLUSTERS"
	TerminationCode_AllocationTimeoutNoReadyClusters                    TerminationCode = "ALLOCATION_TIMEOUT_NO_READY_CLUSTERS"
	TerminationCode_AllocationTimeoutNoWarmedUpClusters                 TerminationCode = "ALLOCATION_TIMEOUT_NO_WARMED_UP_CLUSTERS"
	TerminationCode_AllocationTimeoutNodeDaemonNotReady                 TerminationCode = "ALLOCATION_TIMEOUT_NODE_DAEMON_NOT_READY"
	TerminationCode_AllocationTimeoutNoHealthyClusters                  TerminationCode = "ALLOCATION_TIMEOUT_NO_HEALTHY_CLUSTERS"
	TerminationCode_NetvisorSetupTimeout                                TerminationCode = "NETVISOR_SETUP_TIMEOUT"
	TerminationCode_NoMatchedK8sTestingTag                              TerminationCode = "NO_MATCHED_K8S_TESTING_TAG"
	TerminationCode_CloudProviderResourceStockoutDueToMisconfig         TerminationCode = "CLOUD_PROVIDER_RESOURCE_STOCKOUT_DUE_TO_MISCONFIG"
	TerminationCode_GkeBasedClusterTermination                          TerminationCode = "GKE_BASED_CLUSTER_TERMINATION"
	TerminationCode_AllocationTimeoutNoHealthyAndWarmedUpClusters       TerminationCode = "ALLOCATION_TIMEOUT_NO_HEALTHY_AND_WARMED_UP_CLUSTERS"
	TerminationCode_DockerInvalidOsException                            TerminationCode = "DOCKER_INVALID_OS_EXCEPTION"
	TerminationCode_DockerContainerCreationException                    TerminationCode = "DOCKER_CONTAINER_CREATION_EXCEPTION"
	TerminationCode_DockerImageTooLargeForInstanceException             TerminationCode = "DOCKER_IMAGE_TOO_LARGE_FOR_INSTANCE_EXCEPTION"
	TerminationCode_DnsResolutionError                                  TerminationCode = "DNS_RESOLUTION_ERROR"
	TerminationCode_GcpDeniedByOrgPolicy                                TerminationCode = "GCP_DENIED_BY_ORG_POLICY"
	TerminationCode_SecretPermissionDenied                              TerminationCode = "SECRET_PERMISSION_DENIED"
	TerminationCode_NetworkCheckNicFailure                              TerminationCode = "NETWORK_CHECK_NIC_FAILURE"
	TerminationCode_NetworkCheckDnsServerFailure                        TerminationCode = "NETWORK_CHECK_DNS_SERVER_FAILURE"
	TerminationCode_NetworkCheckStorageFailure                          TerminationCode = "NETWORK_CHECK_STORAGE_FAILURE"
	TerminationCode_NetworkCheckMetadataEndpointFailure                 TerminationCode = "NETWORK_CHECK_METADATA_ENDPOINT_FAILURE"
	TerminationCode_NetworkCheckControlPlaneFailure                     TerminationCode = "NETWORK_CHECK_CONTROL_PLANE_FAILURE"
	TerminationCode_NetworkCheckMultipleComponentsFailure               TerminationCode = "NETWORK_CHECK_MULTIPLE_COMPONENTS_FAILURE"
	TerminationCode_DriverUnhealthy                                     TerminationCode = "DRIVER_UNHEALTHY"
	TerminationCode_SecurityAgentsFailedInitialVerification             TerminationCode = "SECURITY_AGENTS_FAILED_INITIAL_VERIFICATION"
	TerminationCode_DriverDnsResolutionFailure                          TerminationCode = "DRIVER_DNS_RESOLUTION_FAILURE"
	TerminationCode_NoActivatedK8s                                      TerminationCode = "NO_ACTIVATED_K8S"
	TerminationCode_UsagePolicyEntitlementDenied                        TerminationCode = "USAGE_POLICY_ENTITLEMENT_DENIED"
	TerminationCode_NoActivatedK8sTestingTag                            TerminationCode = "NO_ACTIVATED_K8S_TESTING_TAG"
	TerminationCode_K8sActivePodQuotaExceeded                           TerminationCode = "K8S_ACTIVE_POD_QUOTA_EXCEEDED"
	TerminationCode_CloudAccountPodQuotaExceeded                        TerminationCode = "CLOUD_ACCOUNT_POD_QUOTA_EXCEEDED"
	TerminationCode_NetworkCheckNicFailureDueToMisconfig                TerminationCode = "NETWORK_CHECK_NIC_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_NetworkCheckDnsServerFailureDueToMisconfig          TerminationCode = "NETWORK_CHECK_DNS_SERVER_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_NetworkCheckStorageFailureDueToMisconfig            TerminationCode = "NETWORK_CHECK_STORAGE_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_NetworkCheckMetadataEndpointFailureDueToMisconfig   TerminationCode = "NETWORK_CHECK_METADATA_ENDPOINT_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_NetworkCheckControlPlaneFailureDueToMisconfig       TerminationCode = "NETWORK_CHECK_CONTROL_PLANE_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_NetworkCheckMultipleComponentsFailureDueToMisconfig TerminationCode = "NETWORK_CHECK_MULTIPLE_COMPONENTS_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_DbrImageResolutionFailure                           TerminationCode = "DBR_IMAGE_RESOLUTION_FAILURE"
	TerminationCode_ControlPlaneConnectionFailure                       TerminationCode = "CONTROL_PLANE_CONNECTION_FAILURE"
	TerminationCode_ControlPlaneConnectionFailureDueToMisconfig         TerminationCode = "CONTROL_PLANE_CONNECTION_FAILURE_DUE_TO_MISCONFIG"
	TerminationCode_RateLimited                                         TerminationCode = "RATE_LIMITED"
	TerminationCode_MtlsPortConnectivityFailure                         TerminationCode = "MTLS_PORT_CONNECTIVITY_FAILURE"
	TerminationCode_HivemetastoreConnectivityFailure                    TerminationCode = "HIVEMETASTORE_CONNECTIVITY_FAILURE"
	TerminationCode_SecretCreationAccessDenied                          TerminationCode = "SECRET_CREATION_ACCESS_DENIED"
)

// String representation for [fmt.Print].
func (f *TerminationCode) String() string {
	return string(*f)
}

type TerminationType string

const (
	TerminationType_Success      TerminationType = "SUCCESS"
	TerminationType_ClientError  TerminationType = "CLIENT_ERROR"
	TerminationType_ServiceFault TerminationType = "SERVICE_FAULT"
	TerminationType_CloudFailure TerminationType = "CLOUD_FAILURE"
)

// String representation for [fmt.Print].
func (f *TerminationType) String() string {
	return string(*f)
}

type ClusterState_ClusterState string

const (
	ClusterState_ClusterState_Pending     ClusterState_ClusterState = "PENDING"
	ClusterState_ClusterState_Running     ClusterState_ClusterState = "RUNNING"
	ClusterState_ClusterState_Restarting  ClusterState_ClusterState = "RESTARTING"
	ClusterState_ClusterState_Resizing    ClusterState_ClusterState = "RESIZING"
	ClusterState_ClusterState_Terminating ClusterState_ClusterState = "TERMINATING"
	ClusterState_ClusterState_Terminated  ClusterState_ClusterState = "TERMINATED"
	ClusterState_ClusterState_Error       ClusterState_ClusterState = "ERROR"
	ClusterState_ClusterState_Unknown     ClusterState_ClusterState = "UNKNOWN"
)

// String representation for [fmt.Print].
func (f *ClusterState_ClusterState) String() string {
	return string(*f)
}

// A storage location in Adls Gen2.
type Adlsgen2Info struct {
	Destination *string `json:"destination"`
}

type AutoScale struct {
	MinWorkers *int `json:"min_workers"`
	MaxWorkers *int `json:"max_workers"`
}

// Attributes set during cluster creation which are related to Amazon Web
// Services..
type AwsAttributes struct {
	FirstOnDemand       *int             `json:"first_on_demand"`
	Availability        *AwsAvailability `json:"availability"`
	ZoneId              *string          `json:"zone_id"`
	InstanceProfileArn  *string          `json:"instance_profile_arn"`
	SpotBidPricePercent *int             `json:"spot_bid_price_percent"`
	EbsVolumeType       *EbsVolumeType   `json:"ebs_volume_type"`
	EbsVolumeCount      *int             `json:"ebs_volume_count"`
	EbsVolumeSize       *int             `json:"ebs_volume_size"`
	EbsVolumeIops       *int             `json:"ebs_volume_iops"`
	EbsVolumeThroughput *int             `json:"ebs_volume_throughput"`
}

// Attributes set during cluster creation which are related to Microsoft Azure..
type AzureAttributes struct {
	LogAnalyticsInfo *LogAnalyticsInfo  `json:"log_analytics_info"`
	FirstOnDemand    *int               `json:"first_on_demand"`
	Availability     *AzureAvailability `json:"availability"`
	SpotBidMaxPrice  *float64           `json:"spot_bid_max_price"`
}

type ChangeClusterOwner struct {
	ClusterId     *string `json:"cluster_id"`
	OwnerUsername *string `json:"owner_username"`
}

type ChangeClusterOwner_Response struct {
}

type CloneCluster struct {
	SourceClusterId *string `json:"source_cluster_id"`
}

type CloudProviderNodeInfo struct {
	Status []CloudProviderNodeStatus `json:"status"`
}

// Describes all of the metadata about a single Spark cluster in <Databricks>..
type ClusterInfo struct {
	ClusterId                  *string                    `json:"cluster_id"`
	CreatorUserName            *string                    `json:"creator_user_name"`
	State                      *ClusterState_ClusterState `json:"state"`
	StateMessage               *string                    `json:"state_message"`
	ClusterMemoryMb            *int64                     `json:"cluster_memory_mb"`
	ClusterCores               *float64                   `json:"cluster_cores"`
	DefaultTags                map[string]string          `json:"default_tags"`
	ClusterLogStatus           *LogSyncStatus             `json:"cluster_log_status"`
	TerminationReason          *TerminationReason         `json:"termination_reason"`
	Spec                       *ClusterInfo_ComputeSpec   `json:"spec"`
	Driver                     *SparkInfo_SparkNode       `json:"driver"`
	Executors                  []SparkInfo_SparkNode      `json:"executors"`
	SparkContextId             *int64                     `json:"spark_context_id"`
	JdbcPort                   *int                       `json:"jdbc_port"`
	ClusterName                *string                    `json:"cluster_name"`
	SparkVersion               *string                    `json:"spark_version"`
	SparkConf                  map[string]string          `json:"spark_conf"`
	AwsAttributes              *AwsAttributes             `json:"aws_attributes"`
	AzureAttributes            *AzureAttributes           `json:"azure_attributes"`
	GcpAttributes              *GcpAttributes             `json:"gcp_attributes"`
	NodeTypeId                 *string                    `json:"node_type_id"`
	DriverNodeTypeId           *string                    `json:"driver_node_type_id"`
	WorkerNodeTypeFlexibility  *NodeTypeFlexibility       `json:"worker_node_type_flexibility"`
	DriverNodeTypeFlexibility  *NodeTypeFlexibility       `json:"driver_node_type_flexibility"`
	SshPublicKeys              []string                   `json:"ssh_public_keys"`
	CustomTags                 map[string]string          `json:"custom_tags"`
	ClusterLogConf             *ClusterLogConf            `json:"cluster_log_conf"`
	SparkEnvVars               map[string]string          `json:"spark_env_vars"`
	AutoterminationMinutes     *int                       `json:"autotermination_minutes"`
	EnableElasticDisk          *bool                      `json:"enable_elastic_disk"`
	InitScripts                []InitScriptInfo           `json:"init_scripts"`
	DockerImage                *DockerImage               `json:"docker_image"`
	InstancePoolId             *string                    `json:"instance_pool_id"`
	SingleUserName             *string                    `json:"single_user_name"`
	PolicyId                   *string                    `json:"policy_id"`
	EnableLocalDiskEncryption  *bool                      `json:"enable_local_disk_encryption"`
	DriverInstancePoolId       *string                    `json:"driver_instance_pool_id"`
	WorkloadType               *WorkloadType              `json:"workload_type"`
	DataSecurityMode           *DataSecurityMode          `json:"data_security_mode"`
	RuntimeEngine              *RuntimeEngine             `json:"runtime_engine"`
	Kind                       *ComputeKind               `json:"kind"`
	UseMlRuntime               *bool                      `json:"use_ml_runtime"`
	IsSingleNode               *bool                      `json:"is_single_node"`
	RemoteDiskThroughput       *int                       `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize *int                       `json:"total_initial_remote_disk_size"`
	StartTime                  *int64                     `json:"start_time"`
	TerminatedTime             *int64                     `json:"terminated_time"`
	LastStateLossTime          *int64                     `json:"last_state_loss_time"`
	LastRestartedTime          *int64                     `json:"last_restarted_time"`
	NumWorkers                 *int                       `json:"num_workers"`
	Autoscale                  *AutoScale                 `json:"autoscale"`
}

// Contains a snapshot of the latest user specified settings that were used to
// create/edit the cluster..
type ClusterInfo_ComputeSpec struct {
	ApplyPolicyDefaultValues   *bool                `json:"apply_policy_default_values"`
	ClusterName                *string              `json:"cluster_name"`
	SparkVersion               *string              `json:"spark_version"`
	SparkConf                  map[string]string    `json:"spark_conf"`
	AwsAttributes              *AwsAttributes       `json:"aws_attributes"`
	AzureAttributes            *AzureAttributes     `json:"azure_attributes"`
	GcpAttributes              *GcpAttributes       `json:"gcp_attributes"`
	NodeTypeId                 *string              `json:"node_type_id"`
	DriverNodeTypeId           *string              `json:"driver_node_type_id"`
	WorkerNodeTypeFlexibility  *NodeTypeFlexibility `json:"worker_node_type_flexibility"`
	DriverNodeTypeFlexibility  *NodeTypeFlexibility `json:"driver_node_type_flexibility"`
	SshPublicKeys              []string             `json:"ssh_public_keys"`
	CustomTags                 map[string]string    `json:"custom_tags"`
	ClusterLogConf             *ClusterLogConf      `json:"cluster_log_conf"`
	SparkEnvVars               map[string]string    `json:"spark_env_vars"`
	AutoterminationMinutes     *int                 `json:"autotermination_minutes"`
	EnableElasticDisk          *bool                `json:"enable_elastic_disk"`
	InitScripts                []InitScriptInfo     `json:"init_scripts"`
	DockerImage                *DockerImage         `json:"docker_image"`
	InstancePoolId             *string              `json:"instance_pool_id"`
	SingleUserName             *string              `json:"single_user_name"`
	PolicyId                   *string              `json:"policy_id"`
	EnableLocalDiskEncryption  *bool                `json:"enable_local_disk_encryption"`
	DriverInstancePoolId       *string              `json:"driver_instance_pool_id"`
	WorkloadType               *WorkloadType        `json:"workload_type"`
	DataSecurityMode           *DataSecurityMode    `json:"data_security_mode"`
	RuntimeEngine              *RuntimeEngine       `json:"runtime_engine"`
	Kind                       *ComputeKind         `json:"kind"`
	UseMlRuntime               *bool                `json:"use_ml_runtime"`
	IsSingleNode               *bool                `json:"is_single_node"`
	RemoteDiskThroughput       *int                 `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize *int                 `json:"total_initial_remote_disk_size"`
	NumWorkers                 *int                 `json:"num_workers"`
	Autoscale                  *AutoScale           `json:"autoscale"`
}

type ClusterInfo_ComputeSpec_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark configuration key-value pairs.
type ClusterInfo_ComputeSpec_SparkConfEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark environment variable key-value pairs.
type ClusterInfo_ComputeSpec_SparkEnvVarsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ClusterInfo_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ClusterInfo_DefaultTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark configuration key-value pairs.
type ClusterInfo_SparkConfEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark environment variable key-value pairs.
type ClusterInfo_SparkEnvVarsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Cluster log delivery config.
type ClusterLogConf struct {
	Dbfs    *DbfsStorageInfo    `json:"dbfs"`
	S3      *S3StorageInfo      `json:"s3"`
	Volumes *VolumesStorageInfo `json:"volumes"`
}

type ClusterState struct {
}

type CreateCluster struct {
	ApplyPolicyDefaultValues   *bool                `json:"apply_policy_default_values"`
	CloneFrom                  *CloneCluster        `json:"clone_from"`
	NumWorkers                 *int                 `json:"num_workers"`
	Autoscale                  *AutoScale           `json:"autoscale"`
	ClusterName                *string              `json:"cluster_name"`
	SparkVersion               *string              `json:"spark_version"`
	SparkConf                  map[string]string    `json:"spark_conf"`
	AwsAttributes              *AwsAttributes       `json:"aws_attributes"`
	AzureAttributes            *AzureAttributes     `json:"azure_attributes"`
	GcpAttributes              *GcpAttributes       `json:"gcp_attributes"`
	NodeTypeId                 *string              `json:"node_type_id"`
	DriverNodeTypeId           *string              `json:"driver_node_type_id"`
	WorkerNodeTypeFlexibility  *NodeTypeFlexibility `json:"worker_node_type_flexibility"`
	DriverNodeTypeFlexibility  *NodeTypeFlexibility `json:"driver_node_type_flexibility"`
	SshPublicKeys              []string             `json:"ssh_public_keys"`
	CustomTags                 map[string]string    `json:"custom_tags"`
	ClusterLogConf             *ClusterLogConf      `json:"cluster_log_conf"`
	SparkEnvVars               map[string]string    `json:"spark_env_vars"`
	AutoterminationMinutes     *int                 `json:"autotermination_minutes"`
	EnableElasticDisk          *bool                `json:"enable_elastic_disk"`
	InitScripts                []InitScriptInfo     `json:"init_scripts"`
	DockerImage                *DockerImage         `json:"docker_image"`
	InstancePoolId             *string              `json:"instance_pool_id"`
	SingleUserName             *string              `json:"single_user_name"`
	PolicyId                   *string              `json:"policy_id"`
	EnableLocalDiskEncryption  *bool                `json:"enable_local_disk_encryption"`
	DriverInstancePoolId       *string              `json:"driver_instance_pool_id"`
	WorkloadType               *WorkloadType        `json:"workload_type"`
	DataSecurityMode           *DataSecurityMode    `json:"data_security_mode"`
	RuntimeEngine              *RuntimeEngine       `json:"runtime_engine"`
	Kind                       *ComputeKind         `json:"kind"`
	UseMlRuntime               *bool                `json:"use_ml_runtime"`
	IsSingleNode               *bool                `json:"is_single_node"`
	RemoteDiskThroughput       *int                 `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize *int                 `json:"total_initial_remote_disk_size"`
}

type CreateCluster_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CreateCluster_Response struct {
	ClusterId *string `json:"cluster_id"`
}

// Spark configuration key-value pairs.
type CreateCluster_SparkConfEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark environment variable key-value pairs.
type CreateCluster_SparkEnvVarsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// A storage location in DBFS.
type DbfsStorageInfo struct {
	Destination *string `json:"destination"`
}

type DeleteCluster struct {
	ClusterId *string `json:"cluster_id"`
}

type DeleteCluster_Response struct {
}

type DockerBasicAuth struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type DockerImage struct {
	Url       *string          `json:"url"`
	BasicAuth *DockerBasicAuth `json:"basic_auth"`
}

type EditCluster struct {
	ClusterId                  *string              `json:"cluster_id"`
	ApplyPolicyDefaultValues   *bool                `json:"apply_policy_default_values"`
	NumWorkers                 *int                 `json:"num_workers"`
	Autoscale                  *AutoScale           `json:"autoscale"`
	ClusterName                *string              `json:"cluster_name"`
	SparkVersion               *string              `json:"spark_version"`
	SparkConf                  map[string]string    `json:"spark_conf"`
	AwsAttributes              *AwsAttributes       `json:"aws_attributes"`
	AzureAttributes            *AzureAttributes     `json:"azure_attributes"`
	GcpAttributes              *GcpAttributes       `json:"gcp_attributes"`
	NodeTypeId                 *string              `json:"node_type_id"`
	DriverNodeTypeId           *string              `json:"driver_node_type_id"`
	WorkerNodeTypeFlexibility  *NodeTypeFlexibility `json:"worker_node_type_flexibility"`
	DriverNodeTypeFlexibility  *NodeTypeFlexibility `json:"driver_node_type_flexibility"`
	SshPublicKeys              []string             `json:"ssh_public_keys"`
	CustomTags                 map[string]string    `json:"custom_tags"`
	ClusterLogConf             *ClusterLogConf      `json:"cluster_log_conf"`
	SparkEnvVars               map[string]string    `json:"spark_env_vars"`
	AutoterminationMinutes     *int                 `json:"autotermination_minutes"`
	EnableElasticDisk          *bool                `json:"enable_elastic_disk"`
	InitScripts                []InitScriptInfo     `json:"init_scripts"`
	DockerImage                *DockerImage         `json:"docker_image"`
	InstancePoolId             *string              `json:"instance_pool_id"`
	SingleUserName             *string              `json:"single_user_name"`
	PolicyId                   *string              `json:"policy_id"`
	EnableLocalDiskEncryption  *bool                `json:"enable_local_disk_encryption"`
	DriverInstancePoolId       *string              `json:"driver_instance_pool_id"`
	WorkloadType               *WorkloadType        `json:"workload_type"`
	DataSecurityMode           *DataSecurityMode    `json:"data_security_mode"`
	RuntimeEngine              *RuntimeEngine       `json:"runtime_engine"`
	Kind                       *ComputeKind         `json:"kind"`
	UseMlRuntime               *bool                `json:"use_ml_runtime"`
	IsSingleNode               *bool                `json:"is_single_node"`
	RemoteDiskThroughput       *int                 `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize *int                 `json:"total_initial_remote_disk_size"`
}

type EditCluster_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type EditCluster_Response struct {
}

// Spark configuration key-value pairs.
type EditCluster_SparkConfEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark environment variable key-value pairs.
type EditCluster_SparkEnvVarsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Attributes set during cluster creation which are related to GCP..
type GcpAttributes struct {
	UsePreemptibleExecutors *bool                    `json:"use_preemptible_executors"`
	GoogleServiceAccount    *string                  `json:"google_service_account"`
	BootDiskSize            *int                     `json:"boot_disk_size"`
	Availability            *GcpAvailability         `json:"availability"`
	ZoneId                  *string                  `json:"zone_id"`
	LocalSsdCount           *int                     `json:"local_ssd_count"`
	FirstOnDemand           *int                     `json:"first_on_demand"`
	ConfidentialComputeType *ConfidentialComputeType `json:"confidential_compute_type"`
}

// A storage location in Google Cloud Platform's GCS.
type GcsStorageInfo struct {
	Destination *string `json:"destination"`
}

type GetCluster struct {
	ClusterId *string `json:"cluster_id"`
}

// Returns the list of all Spark versions that can be used to create clusters..
type GetSparkVersions struct {
}

type GetSparkVersions_Response struct {
	Versions []SparkVersion `json:"versions"`
}

// Config for an individual init script Next ID: 11.
type InitScriptInfo struct {
	Dbfs      *DbfsStorageInfo      `json:"dbfs"`
	S3        *S3StorageInfo        `json:"s3"`
	File      *LocalFileInfo        `json:"file"`
	Gcs       *GcsStorageInfo       `json:"gcs"`
	Abfss     *Adlsgen2Info         `json:"abfss"`
	Workspace *WorkspaceStorageInfo `json:"workspace"`
	Volumes   *VolumesStorageInfo   `json:"volumes"`
}

type ListAvailableZones struct {
}

type ListAvailableZones_Response struct {
	Zones       []string `json:"zones"`
	DefaultZone *string  `json:"default_zone"`
}

type ListClusters struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListClusters_Response struct {
	Clusters      []ClusterInfo `json:"clusters"`
	NextPageToken *string       `json:"next_page_token"`
	PrevPageToken *string       `json:"prev_page_token"`
}

type ListNodeTypes struct {
}

type ListNodeTypes_Response struct {
	NodeTypes []NodeType `json:"node_types"`
}

type LocalFileInfo struct {
	Destination *string `json:"destination"`
}

type LogAnalyticsInfo struct {
	LogAnalyticsWorkspaceId *string `json:"log_analytics_workspace_id"`
	LogAnalyticsPrimaryKey  *string `json:"log_analytics_primary_key"`
}

// The log delivery status.
type LogSyncStatus struct {
	LastAttempted *int64  `json:"last_attempted"`
	LastException *string `json:"last_exception"`
}

// This structure embodies the machine type that hosts spark containers Note:
// this should be an internal data structure for now It is defined in proto in
// case we want to send it over the wire in the future (which is likely).
type NodeInstanceType struct {
	InstanceTypeId      *string `json:"instance_type_id"`
	LocalDisks          *int    `json:"local_disks"`
	LocalDiskSizeGb     *int    `json:"local_disk_size_gb"`
	LocalNvmeDiskSizeGb *int    `json:"local_nvme_disk_size_gb"`
	LocalNvmeDisks      *int    `json:"local_nvme_disks"`
}

// A description of a Spark node type including both the dimensions of the node
// and the instance type on which it will be hosted..
type NodeType struct {
	NodeTypeId            *string                `json:"node_type_id"`
	MemoryMb              *int                   `json:"memory_mb"`
	NumCores              *float64               `json:"num_cores"`
	Description           *string                `json:"description"`
	InstanceTypeId        *string                `json:"instance_type_id"`
	IsDeprecated          *bool                  `json:"is_deprecated"`
	Category              *string                `json:"category"`
	SupportEbsVolumes     *bool                  `json:"support_ebs_volumes"`
	SupportClusterTags    *bool                  `json:"support_cluster_tags"`
	NumGpus               *int                   `json:"num_gpus"`
	NodeInstanceType      *NodeInstanceType      `json:"node_instance_type"`
	IsHidden              *bool                  `json:"is_hidden"`
	SupportPortForwarding *bool                  `json:"support_port_forwarding"`
	DisplayOrder          *int                   `json:"display_order"`
	IsIoCacheEnabled      *bool                  `json:"is_io_cache_enabled"`
	NodeInfo              *CloudProviderNodeInfo `json:"node_info"`
	PhotonWorkerCapable   *bool                  `json:"photon_worker_capable"`
	PhotonDriverCapable   *bool                  `json:"photon_driver_capable"`
	IsEncryptedInTransit  *bool                  `json:"is_encrypted_in_transit"`
	IsGraviton            *bool                  `json:"is_graviton"`
}

// Configuration for flexible node types, allowing fallback to alternate node
// types during cluster launch and upscale..
type NodeTypeFlexibility struct {
	AlternateNodeTypeIds []string `json:"alternate_node_type_ids"`
}

type PermanentDeleteCluster struct {
	ClusterId *string `json:"cluster_id"`
}

type PermanentDeleteCluster_Response struct {
}

type PinCluster struct {
	ClusterId *string `json:"cluster_id"`
}

type PinCluster_Response struct {
}

type ResizeCluster struct {
	ClusterId  *string    `json:"cluster_id"`
	NumWorkers *int       `json:"num_workers"`
	Autoscale  *AutoScale `json:"autoscale"`
}

type ResizeCluster_Response struct {
}

type RestartCluster struct {
	ClusterId   *string `json:"cluster_id"`
	RestartUser *string `json:"restart_user"`
}

type RestartCluster_Response struct {
}

// A storage location in Amazon S3.
type S3StorageInfo struct {
	Destination      *string `json:"destination"`
	Region           *string `json:"region"`
	Endpoint         *string `json:"endpoint"`
	EnableEncryption *bool   `json:"enable_encryption"`
	EncryptionType   *string `json:"encryption_type"`
	KmsKey           *string `json:"kms_key"`
	CannedAcl        *string `json:"canned_acl"`
}

// Provides information about Spark running inside a cluster. This is used in
// both the [[ClusterInfo]] for Cluster APIs and persisted cluster proto..
type SparkInfo struct {
}

// Describes a specific Spark driver or executor..
type SparkInfo_SparkNode struct {
	PrivateIp         *string                                     `json:"private_ip"`
	PublicDns         *string                                     `json:"public_dns"`
	NodeId            *string                                     `json:"node_id"`
	InstanceId        *string                                     `json:"instance_id"`
	StartTimestamp    *int64                                      `json:"start_timestamp"`
	NodeAwsAttributes *SparkInfo_SparkNode_SparkNodeAwsAttributes `json:"node_aws_attributes"`
	HostPrivateIp     *string                                     `json:"host_private_ip"`
}

// Attributes specific to AWS for a Spark node..
type SparkInfo_SparkNode_SparkNodeAwsAttributes struct {
	IsSpot *bool `json:"is_spot"`
}

type SparkVersion struct {
	Key  *string `json:"key"`
	Name *string `json:"name"`
}

type StartCluster struct {
	ClusterId *string `json:"cluster_id"`
}

type StartCluster_Response struct {
}

type TerminationReason struct {
	Code       *TerminationCode  `json:"code"`
	Type       *TerminationType  `json:"type"`
	Parameters map[string]string `json:"parameters"`
}

type TerminationReason_ParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type UnpinCluster struct {
	ClusterId *string `json:"cluster_id"`
}

type UnpinCluster_Response struct {
}

type UpdateCluster struct {
	ClusterId  *string                              `json:"cluster_id"`
	Cluster    *UpdateCluster_UpdateClusterResource `json:"cluster"`
	UpdateMask *string                              `json:"update_mask"`
}

type UpdateCluster_Response struct {
}

type UpdateCluster_UpdateClusterResource struct {
	NumWorkers                 *int                 `json:"num_workers"`
	Autoscale                  *AutoScale           `json:"autoscale"`
	ClusterName                *string              `json:"cluster_name"`
	SparkVersion               *string              `json:"spark_version"`
	SparkConf                  map[string]string    `json:"spark_conf"`
	AwsAttributes              *AwsAttributes       `json:"aws_attributes"`
	AzureAttributes            *AzureAttributes     `json:"azure_attributes"`
	GcpAttributes              *GcpAttributes       `json:"gcp_attributes"`
	NodeTypeId                 *string              `json:"node_type_id"`
	DriverNodeTypeId           *string              `json:"driver_node_type_id"`
	WorkerNodeTypeFlexibility  *NodeTypeFlexibility `json:"worker_node_type_flexibility"`
	DriverNodeTypeFlexibility  *NodeTypeFlexibility `json:"driver_node_type_flexibility"`
	SshPublicKeys              []string             `json:"ssh_public_keys"`
	CustomTags                 map[string]string    `json:"custom_tags"`
	ClusterLogConf             *ClusterLogConf      `json:"cluster_log_conf"`
	SparkEnvVars               map[string]string    `json:"spark_env_vars"`
	AutoterminationMinutes     *int                 `json:"autotermination_minutes"`
	EnableElasticDisk          *bool                `json:"enable_elastic_disk"`
	InitScripts                []InitScriptInfo     `json:"init_scripts"`
	DockerImage                *DockerImage         `json:"docker_image"`
	InstancePoolId             *string              `json:"instance_pool_id"`
	SingleUserName             *string              `json:"single_user_name"`
	PolicyId                   *string              `json:"policy_id"`
	EnableLocalDiskEncryption  *bool                `json:"enable_local_disk_encryption"`
	DriverInstancePoolId       *string              `json:"driver_instance_pool_id"`
	WorkloadType               *WorkloadType        `json:"workload_type"`
	DataSecurityMode           *DataSecurityMode    `json:"data_security_mode"`
	RuntimeEngine              *RuntimeEngine       `json:"runtime_engine"`
	Kind                       *ComputeKind         `json:"kind"`
	UseMlRuntime               *bool                `json:"use_ml_runtime"`
	IsSingleNode               *bool                `json:"is_single_node"`
	RemoteDiskThroughput       *int                 `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize *int                 `json:"total_initial_remote_disk_size"`
}

type UpdateCluster_UpdateClusterResource_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark configuration key-value pairs.
type UpdateCluster_UpdateClusterResource_SparkConfEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark environment variable key-value pairs.
type UpdateCluster_UpdateClusterResource_SparkEnvVarsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// A storage location back by UC Volumes..
type VolumesStorageInfo struct {
	Destination *string `json:"destination"`
}

// Cluster Attributes showing for clusters workload types..
type WorkloadType struct {
	Clients *WorkloadType_ClientsTypes `json:"clients"`
}

type WorkloadType_ClientsTypes struct {
	Notebooks *bool `json:"notebooks"`
	Jobs      *bool `json:"jobs"`
}

// A storage location in Workspace Filesystem (WSFS).
type WorkspaceStorageInfo struct {
	Destination *string `json:"destination"`
}
