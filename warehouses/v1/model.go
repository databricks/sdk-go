package v1

type ChannelName string

const (
	ChannelName_ChannelNameUnspecified ChannelName = "CHANNEL_NAME_UNSPECIFIED"
	ChannelName_ChannelNamePreview     ChannelName = "CHANNEL_NAME_PREVIEW"
	ChannelName_ChannelNameCurrent     ChannelName = "CHANNEL_NAME_CURRENT"
	ChannelName_ChannelNamePrevious    ChannelName = "CHANNEL_NAME_PREVIOUS"
	ChannelName_ChannelNameCustom      ChannelName = "CHANNEL_NAME_CUSTOM"
)

// String representation for [fmt.Print].
func (f *ChannelName) String() string {
	return string(*f)
}

type DefaultWarehouseOverrideType string

const (
	DefaultWarehouseOverrideType_DefaultWarehouseOverrideTypeUnspecified DefaultWarehouseOverrideType = "DEFAULT_WAREHOUSE_OVERRIDE_TYPE_UNSPECIFIED"
	DefaultWarehouseOverrideType_LastSelected                            DefaultWarehouseOverrideType = "LAST_SELECTED"
	DefaultWarehouseOverrideType_Custom                                  DefaultWarehouseOverrideType = "CUSTOM"
)

// String representation for [fmt.Print].
func (f *DefaultWarehouseOverrideType) String() string {
	return string(*f)
}

type EndpointSecurityPolicy string

const (
	EndpointSecurityPolicy_None              EndpointSecurityPolicy = "NONE"
	EndpointSecurityPolicy_DataAccessControl EndpointSecurityPolicy = "DATA_ACCESS_CONTROL"
	EndpointSecurityPolicy_Passthrough       EndpointSecurityPolicy = "PASSTHROUGH"
)

// String representation for [fmt.Print].
func (f *EndpointSecurityPolicy) String() string {
	return string(*f)
}

type EndpointSpotInstancePolicy string

const (
	EndpointSpotInstancePolicy_PolicyUnspecified    EndpointSpotInstancePolicy = "POLICY_UNSPECIFIED"
	EndpointSpotInstancePolicy_CostOptimized        EndpointSpotInstancePolicy = "COST_OPTIMIZED"
	EndpointSpotInstancePolicy_ReliabilityOptimized EndpointSpotInstancePolicy = "RELIABILITY_OPTIMIZED"
)

// String representation for [fmt.Print].
func (f *EndpointSpotInstancePolicy) String() string {
	return string(*f)
}

type EndpointState string

const (
	EndpointState_Starting EndpointState = "STARTING"
	EndpointState_Running  EndpointState = "RUNNING"
	EndpointState_Stopping EndpointState = "STOPPING"
	EndpointState_Stopped  EndpointState = "STOPPED"
	EndpointState_Deleting EndpointState = "DELETING"
	EndpointState_Deleted  EndpointState = "DELETED"
)

// String representation for [fmt.Print].
func (f *EndpointState) String() string {
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

type WarehouseType string

const (
	WarehouseType_TypeUnspecified WarehouseType = "TYPE_UNSPECIFIED"
	WarehouseType_Classic         WarehouseType = "CLASSIC"
	WarehouseType_Pro             WarehouseType = "PRO"
	WarehouseType_Reyden          WarehouseType = "REYDEN"
)

// String representation for [fmt.Print].
func (f *WarehouseType) String() string {
	return string(*f)
}

type EndpointHealth_Status string

const (
	EndpointHealth_Status_StatusUnspecified EndpointHealth_Status = "STATUS_UNSPECIFIED"
	EndpointHealth_Status_Healthy           EndpointHealth_Status = "HEALTHY"
	EndpointHealth_Status_Degraded          EndpointHealth_Status = "DEGRADED"
	EndpointHealth_Status_Failed            EndpointHealth_Status = "FAILED"
)

// String representation for [fmt.Print].
func (f *EndpointHealth_Status) String() string {
	return string(*f)
}

// Configures the channel name and DBSQL version of the warehouse.
// CHANNEL_NAME_CUSTOM should be chosen only when `dbsql_version` is specified..
type Channel struct {
	Name         *ChannelName `json:"name"`
	DbsqlVersion *string      `json:"dbsql_version"`
}

// Request message for CreateDefaultWarehouseOverride..
type CreateDefaultWarehouseOverrideRequest struct {
	DefaultWarehouseOverrideId *string                   `json:"default_warehouse_override_id"`
	DefaultWarehouseOverride   *DefaultWarehouseOverride `json:"default_warehouse_override"`
}

// Creates a new SQL warehouse..
type CreateWarehouse struct {
	Name                    *string                     `json:"name"`
	ClusterSize             *string                     `json:"cluster_size"`
	MinNumClusters          *int                        `json:"min_num_clusters"`
	MaxNumClusters          *int                        `json:"max_num_clusters"`
	AutoStopMins            *int                        `json:"auto_stop_mins"`
	CreatorName             *string                     `json:"creator_name"`
	InstanceProfileArn      *string                     `json:"instance_profile_arn"`
	Tags                    *EndpointTags               `json:"tags"`
	SpotInstancePolicy      *EndpointSpotInstancePolicy `json:"spot_instance_policy"`
	EnablePhoton            *bool                       `json:"enable_photon"`
	Channel                 *Channel                    `json:"channel"`
	EnableServerlessCompute *bool                       `json:"enable_serverless_compute"`
	WarehouseType           *WarehouseType              `json:"warehouse_type"`
}

type CreateWarehouse_Response struct {
	Id *string `json:"id"`
}

// Represents a per-user default warehouse override configuration. This resource
// allows users or administrators to customize how a user's default warehouse is
// selected for SQL operations. If no override exists for a user, the workspace
// default warehouse will be used..
type DefaultWarehouseOverride struct {
	Name                       *string                       `json:"name"`
	DefaultWarehouseOverrideId *string                       `json:"default_warehouse_override_id"`
	Type                       *DefaultWarehouseOverrideType `json:"type"`
	WarehouseId                *string                       `json:"warehouse_id"`
}

// Request message for DeleteDefaultWarehouseOverride..
type DeleteDefaultWarehouseOverrideRequest struct {
	Name *string `json:"name"`
}

// This is an incremental edit functionality, so all fields except id are
// optional. If a field is set, the corresponding configuration in the SQL
// warehouse is modified. If a field is unset, the existing configuration value
// in the SQL warehouse is retained. Thus, this API is not idempotent..
type EditWarehouseRequest struct {
	Id                      *string                     `json:"id"`
	Name                    *string                     `json:"name"`
	ClusterSize             *string                     `json:"cluster_size"`
	MinNumClusters          *int                        `json:"min_num_clusters"`
	MaxNumClusters          *int                        `json:"max_num_clusters"`
	AutoStopMins            *int                        `json:"auto_stop_mins"`
	CreatorName             *string                     `json:"creator_name"`
	InstanceProfileArn      *string                     `json:"instance_profile_arn"`
	Tags                    *EndpointTags               `json:"tags"`
	SpotInstancePolicy      *EndpointSpotInstancePolicy `json:"spot_instance_policy"`
	EnablePhoton            *bool                       `json:"enable_photon"`
	Channel                 *Channel                    `json:"channel"`
	EnableServerlessCompute *bool                       `json:"enable_serverless_compute"`
	WarehouseType           *WarehouseType              `json:"warehouse_type"`
}

type EditWarehouseRequest_Response struct {
}

type EndpointConfPair struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type EndpointHealth struct {
	Status        *EndpointHealth_Status `json:"status"`
	Message       *string                `json:"message"`
	FailureReason *TerminationReason     `json:"failure_reason"`
	Summary       *string                `json:"summary"`
	Details       *string                `json:"details"`
}

type EndpointInfo struct {
	Id                      *string                     `json:"id"`
	Name                    *string                     `json:"name"`
	ClusterSize             *string                     `json:"cluster_size"`
	MinNumClusters          *int                        `json:"min_num_clusters"`
	MaxNumClusters          *int                        `json:"max_num_clusters"`
	AutoStopMins            *int                        `json:"auto_stop_mins"`
	CreatorName             *string                     `json:"creator_name"`
	InstanceProfileArn      *string                     `json:"instance_profile_arn"`
	Tags                    *EndpointTags               `json:"tags"`
	SpotInstancePolicy      *EndpointSpotInstancePolicy `json:"spot_instance_policy"`
	EnablePhoton            *bool                       `json:"enable_photon"`
	Channel                 *Channel                    `json:"channel"`
	EnableServerlessCompute *bool                       `json:"enable_serverless_compute"`
	WarehouseType           *WarehouseType              `json:"warehouse_type"`
	NumClusters             *int                        `json:"num_clusters"`
	NumActiveSessions       *int64                      `json:"num_active_sessions"`
	State                   *EndpointState              `json:"state"`
	JdbcUrl                 *string                     `json:"jdbc_url"`
	OdbcParams              *OdbcParams                 `json:"odbc_params"`
	Health                  *EndpointHealth             `json:"health"`
}

type EndpointTagPair struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type EndpointTags struct {
	CustomTags []EndpointTagPair `json:"custom_tags"`
}

// Request message for GetDefaultWarehouseOverride..
type GetDefaultWarehouseOverrideRequest struct {
	Name *string `json:"name"`
}

// Fetches the warehouse info for a single SQL warehouse..
type GetWarehouse struct {
	Id *string `json:"id"`
}

type GetWarehouse_Response struct {
	Id                      *string                     `json:"id"`
	Name                    *string                     `json:"name"`
	ClusterSize             *string                     `json:"cluster_size"`
	MinNumClusters          *int                        `json:"min_num_clusters"`
	MaxNumClusters          *int                        `json:"max_num_clusters"`
	AutoStopMins            *int                        `json:"auto_stop_mins"`
	CreatorName             *string                     `json:"creator_name"`
	InstanceProfileArn      *string                     `json:"instance_profile_arn"`
	Tags                    *EndpointTags               `json:"tags"`
	SpotInstancePolicy      *EndpointSpotInstancePolicy `json:"spot_instance_policy"`
	EnablePhoton            *bool                       `json:"enable_photon"`
	Channel                 *Channel                    `json:"channel"`
	EnableServerlessCompute *bool                       `json:"enable_serverless_compute"`
	WarehouseType           *WarehouseType              `json:"warehouse_type"`
	NumClusters             *int                        `json:"num_clusters"`
	NumActiveSessions       *int64                      `json:"num_active_sessions"`
	State                   *EndpointState              `json:"state"`
	JdbcUrl                 *string                     `json:"jdbc_url"`
	OdbcParams              *OdbcParams                 `json:"odbc_params"`
	Health                  *EndpointHealth             `json:"health"`
}

// Fetches the workspace level SQL warehouse configurations. These are the
// configurations that are set centrally and shared by all SQL warehouses in a
// workspace..
type GetWorkspaceWarehouseConfigRequest struct {
}

type GetWorkspaceWarehouseConfigRequest_Response struct {
	SecurityPolicy             *EndpointSecurityPolicy    `json:"security_policy"`
	DataAccessConfig           []EndpointConfPair         `json:"data_access_config"`
	InstanceProfileArn         *string                    `json:"instance_profile_arn"`
	Channel                    *Channel                   `json:"channel"`
	EnableServerlessCompute    *bool                      `json:"enable_serverless_compute"`
	GlobalParam                *RepeatedEndpointConfPairs `json:"global_param"`
	ConfigParam                *RepeatedEndpointConfPairs `json:"config_param"`
	SqlConfigurationParameters *RepeatedEndpointConfPairs `json:"sql_configuration_parameters"`
	GoogleServiceAccount       *string                    `json:"google_service_account"`
	EnabledWarehouseTypes      []WarehouseTypePair        `json:"enabled_warehouse_types"`
}

// Request message for ListDefaultWarehouseOverrides..
type ListDefaultWarehouseOverridesRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

// Response message for ListDefaultWarehouseOverrides..
type ListDefaultWarehouseOverridesResponse struct {
	DefaultWarehouseOverrides []DefaultWarehouseOverride `json:"default_warehouse_overrides"`
	NextPageToken             *string                    `json:"next_page_token"`
}

type OdbcParams struct {
	Hostname *string `json:"hostname"`
	Path     *string `json:"path"`
	Protocol *string `json:"protocol"`
	Port     *int    `json:"port"`
}

type RepeatedEndpointConfPairs struct {
	ConfigPair         []EndpointConfPair `json:"config_pair"`
	ConfigurationPairs []EndpointConfPair `json:"configuration_pairs"`
}

// Sets the workspace level warehouse configuration that is shared by all SQL
// warehouses in this workspace.
//
// This is idempotent..
type SetWorkspaceWarehouseConfigRequest struct {
	SecurityPolicy             *EndpointSecurityPolicy    `json:"security_policy"`
	DataAccessConfig           []EndpointConfPair         `json:"data_access_config"`
	InstanceProfileArn         *string                    `json:"instance_profile_arn"`
	Channel                    *Channel                   `json:"channel"`
	EnableServerlessCompute    *bool                      `json:"enable_serverless_compute"`
	GlobalParam                *RepeatedEndpointConfPairs `json:"global_param"`
	ConfigParam                *RepeatedEndpointConfPairs `json:"config_param"`
	SqlConfigurationParameters *RepeatedEndpointConfPairs `json:"sql_configuration_parameters"`
	GoogleServiceAccount       *string                    `json:"google_service_account"`
	EnabledWarehouseTypes      []WarehouseTypePair        `json:"enabled_warehouse_types"`
}

type SetWorkspaceWarehouseConfigRequest_Response struct {
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

// Request message for UpdateDefaultWarehouseOverride..
type UpdateDefaultWarehouseOverrideRequest struct {
	DefaultWarehouseOverride *DefaultWarehouseOverride `json:"default_warehouse_override"`
	UpdateMask               *string                   `json:"update_mask"`
	AllowMissing             *bool                     `json:"allow_missing"`
}

// * Configuration values to enable or disable the access to specific warehouse
// types in the workspace..
type WarehouseTypePair struct {
	WarehouseType *WarehouseType `json:"warehouse_type"`
	Enabled       *bool          `json:"enabled"`
}

// Deletes a warehouse. This API is idempotent..
type DeleteWarehouseRequest struct {
	Id *string `json:"id"`
}

type DeleteWarehouseRequest_Response struct {
}

// Lists all of the SQL warehouses. TODO: consider paginating to limit the
// number of warehouses returned..
type ListWarehousesRequest struct {
	RunAsUserId *int64  `json:"run_as_user_id"`
	PageSize    *int    `json:"page_size"`
	PageToken   *string `json:"page_token"`
}

type ListWarehousesRequest_Response struct {
	Warehouses    []EndpointInfo `json:"warehouses"`
	NextPageToken *string        `json:"next_page_token"`
}

// Starts a SQL warehouse. This API is idempotent..
type StartRequest struct {
	Id *string `json:"id"`
}

type StartRequest_Response struct {
}

// Stops a SQL warehouse. This API is idempotent..
type StopRequest struct {
	Id *string `json:"id"`
}

type StopRequest_Response struct {
}
