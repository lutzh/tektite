package kafkaprotocol

const (
	// Standard Kafka API keys

	APIKeyProduce            = 0
	APIKeyFetch              = 1
	APIKeyListOffsets        = 2
	APIKeyMetadata           = 3
	APIKeyOffsetCommit       = 8
	APIKeyOffsetFetch        = 9
	APIKeyFindCoordinator    = 10
	ApiKeyJoinGroup          = 11
	ApiKeyHeartbeat          = 12
	ApiKeyLeaveGroup         = 13
	ApiKeySyncGroup          = 14
	ApiKeyDescribeGroups     = 15
	ApiKeyListGroups         = 16
	APIKeySaslHandshake      = 17
	APIKeyAPIVersions        = 18
	APIKeyCreateTopics       = 19
	APIKeyDeleteTopics       = 20
	APIKeyInitProducerId     = 22
	APIKeyAddPartitionsToTxn = 24
	APIKeyAddOffsetsToTxn    = 25
	APIKeyEndTxn             = 26
	APIKeyTxnOffsetCommit    = 28
	ApiKeyDescribeAcls       = 29
	ApiKeyCreateAcls         = 30
	ApiKeyDeleteAcls         = 31
	ApiKeyDescribeConfigs    = 32
	ApiKeyAlterConfigs       = 33
	APIKeySaslAuthenticate   = 36
	ApiKeyCreatePartitions   = 37
	ApiKeyDeleteGroups       = 42
	APIKeyOffsetDelete       = 47
	ApiKeyDescribeCluster    = 60

	// Custom API keys

	ApiKeyPutUserCredentialsRequest = 1000
	ApiKeyDeleteUserRequest         = 1001
)

const (
	ErrorCodeUnknownServerError                 = -1
	ErrorCodeNone                               = 0
	ErrorCodeOffsetOutOfRange                   = 1
	ErrorCodeCorruptMessage                     = 2
	ErrorCodeUnknownTopicOrPartition            = 3
	ErrorCodeInvalidFetchSize                   = 4
	ErrorCodeLeaderNotAvailable                 = 5
	ErrorCodeNotLeaderOrFollower                = 6
	ErrorCodeRequestTimedOut                    = 7
	ErrorCodeBrokerNotAvailable                 = 8
	ErrorCodeReplicaNotAvailable                = 9
	ErrorCodeMessageTooLarge                    = 10
	ErrorCodeStaleControllerEpoch               = 11
	ErrorCodeOffsetMetadataTooLarge             = 12
	ErrorCodeNetworkException                   = 13
	ErrorCodeCoordinatorLoadInProgress          = 14
	ErrorCodeCoordinatorNotAvailable            = 15
	ErrorCodeNotCoordinator                     = 16
	ErrorCodeInvalidTopicException              = 17
	ErrorCodeRecordListTooLarge                 = 18
	ErrorCodeNotEnoughReplicas                  = 19
	ErrorCodeNotEnoughReplicasAfterAppend       = 20
	ErrorCodeInvalidRequiredAcks                = 21
	ErrorCodeIllegalGeneration                  = 22
	ErrorCodeInconsistentGroupProtocol          = 23
	ErrorCodeInvalidGroupID                     = 24
	ErrorCodeUnknownMemberID                    = 25
	ErrorCodeInvalidSessionTimeout              = 26
	ErrorCodeRebalanceInProgress                = 27
	ErrorCodeInvalidCommitOffsetSize            = 28
	ErrorCodeTopicAuthorizationFailed           = 29
	ErrorCodeGroupAuthorizationFailed           = 30
	ErrorCodeClusterAuthorizationFailed         = 31
	ErrorCodeInvalidTimestamp                   = 32
	ErrorCodeUnsupportedSaslMechanism           = 33
	ErrorCodeIllegalSaslState                   = 34
	ErrorCodeUnsupportedVersion                 = 35
	ErrorCodeTopicAlreadyExists                 = 36
	ErrorCodeInvalidPartitions                  = 37
	ErrorCodeInvalidReplicationFactor           = 38
	ErrorCodeInvalidReplicaAssignment           = 39
	ErrorCodeInvalidConfig                      = 40
	ErrorCodeNotController                      = 41
	ErrorCodeInvalidRequest                     = 42
	ErrorCodeUnsupportedForMessageFormat        = 43
	ErrorCodePolicyViolation                    = 44
	ErrorCodeOutOfOrderSequenceNumber           = 45
	ErrorCodeDuplicateSequenceNumber            = 46
	ErrorCodeInvalidProducerEpoch               = 47
	ErrorCodeInvalidTxnState                    = 48
	ErrorCodeInvalidProducerIDMapping           = 49
	ErrorCodeInvalidTransactionTimeout          = 50
	ErrorCodeConcurrentTransactions             = 51
	ErrorCodeTransactionCoordinatorFenced       = 52
	ErrorCodeTransactionalIDAuthorizationFailed = 53
	ErrorCodeSecurityDisabled                   = 54
	ErrorCodeOperationNotAttempted              = 55
	ErrorCodeKafkaStorageErrorCode              = 56
	ErrorCodeLogDirNotFound                     = 57
	ErrorCodeSaslAuthenticationFailed           = 58
	ErrorCodeUnknownProducerID                  = 59
	ErrorCodeReassignmentInProgress             = 60
	ErrorCodeDelegationTokenAuthDisabled        = 61
	ErrorCodeDelegationTokenNotFounc            = 62
	ErrorCodeDelegationTokenOwnerMismatch       = 63
	ErrorCodeDelegationTokenRequestNotAllowed   = 64
	ErrorCodeDelegationTokenAuthorizationFailed = 65
	ErrorCodeDelegationTokenExpired             = 66
	ErrorCodeInvalidPrincipalType               = 67
	ErrorCodeNonEmptyGroup                      = 68
	ErrorCodeGroupIDNotFound                    = 69
	ErrorCodeFetchSessionIDNotFound             = 70
	ErrorCodeInvalidFetchSessionEpoch           = 71

	ErrorCodeNoSuchUser = 1000
)

var SupportedAPIVersions = []ApiVersionsResponseApiVersion{
	{ApiKey: APIKeyProduce, MinVersion: 3, MaxVersion: 3},
	{ApiKey: APIKeyFetch, MinVersion: 2, MaxVersion: 4},
	{ApiKey: APIKeyAPIVersions, MinVersion: 0, MaxVersion: 4},
	{ApiKey: APIKeyMetadata, MinVersion: 1, MaxVersion: 4},
	{ApiKey: APIKeyFindCoordinator, MinVersion: 0, MaxVersion: 1},
	{ApiKey: ApiKeyJoinGroup, MinVersion: 0, MaxVersion: 1},
	{ApiKey: ApiKeySyncGroup, MinVersion: 0, MaxVersion: 0},
	{ApiKey: ApiKeyHeartbeat, MinVersion: 0, MaxVersion: 0},
	{ApiKey: APIKeyListOffsets, MinVersion: 1, MaxVersion: 1},
	{ApiKey: APIKeyOffsetCommit, MinVersion: 2, MaxVersion: 2},
	{ApiKey: APIKeyOffsetFetch, MinVersion: 1, MaxVersion: 1},
	{ApiKey: ApiKeyLeaveGroup, MinVersion: 0, MaxVersion: 0},
	{ApiKey: APIKeySaslHandshake, MinVersion: 0, MaxVersion: 1},
	{ApiKey: APIKeyInitProducerId, MinVersion: 0, MaxVersion: 0},
	{ApiKey: APIKeySaslAuthenticate, MinVersion: 0, MaxVersion: 1},
	{ApiKey: APIKeyCreateTopics, MinVersion: 0, MaxVersion: 5},
	{ApiKey: APIKeyDeleteTopics, MinVersion: 0, MaxVersion: 5},
	{ApiKey: APIKeyOffsetDelete, MinVersion: 0, MaxVersion: 0},
	{ApiKey: ApiKeyListGroups, MinVersion: 0, MaxVersion: 5},
	{ApiKey: ApiKeyDescribeGroups, MinVersion: 0, MaxVersion: 0},
	{ApiKey: ApiKeyDeleteGroups, MinVersion: 0, MaxVersion: 0},
	{ApiKey: ApiKeyCreatePartitions, MinVersion: 0, MaxVersion: 0},
	/*
		Transactions are currently incomplete
		{ApiKey: APIKeyAddPartitionsToTxn, MinVersion: 3, MaxVersion: 3},
		{ApiKey: APIKeyAddOffsetsToTxn, MinVersion: 3, MaxVersion: 3},
		{ApiKey: APIKeyEndTxn, MinVersion: 3, MaxVersion: 3},
		{ApiKey: APIKeyTxnOffsetCommit, MinVersion: 3, MaxVersion: 3},
	*/
	{ApiKey: ApiKeyPutUserCredentialsRequest, MinVersion: 0, MaxVersion: 0},
	{ApiKey: ApiKeyDeleteUserRequest, MinVersion: 0, MaxVersion: 0},
	{ApiKey: ApiKeyDescribeConfigs, MinVersion: 1, MaxVersion: 1},
	{ApiKey: ApiKeyAlterConfigs, MinVersion: 0, MaxVersion: 0},
	{ApiKey: ApiKeyDescribeCluster, MinVersion: 0, MaxVersion: 0},
	{ApiKey: ApiKeyCreateAcls, MinVersion: 3, MaxVersion: 3},
	{ApiKey: ApiKeyDeleteAcls, MinVersion: 3, MaxVersion: 3},
	{ApiKey: ApiKeyDescribeAcls, MinVersion: 3, MaxVersion: 3},
}

var SupportedCustomAPIVersions = []ApiVersionsResponseApiVersion{
	{ApiKey: ApiKeyPutUserCredentialsRequest, MinVersion: 0, MaxVersion: 0},
	{ApiKey: ApiKeyDeleteUserRequest, MinVersion: 0, MaxVersion: 0},
}

type Records struct {
	Data [][]byte
}
