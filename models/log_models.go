package db_models

//
//// DeltaStartupLogs action events
//type DeltaStartupLogs struct {
//	ID            int64     `gorm:"primaryKey"` // auto increment
//	NodeInfo      string    `json:"node_info"`
//	OSDetails     string    `json:"os_details"`
//	IPAddress     string    `json:"ip_address"`
//	DeltaNodeUuid string    `json:"delta_node_uuid"`
//	CreatedAt     time.Time `json:"created_at"` // auto set
//	UpdatedAt     time.Time `json:"updated_at"`
//}
//
//type DealEndpointRequestLog struct {
//	ID          int64  `gorm:"primaryKey"` // auto increment
//	NodeInfo    string `json:"node_info"`
//	Information string `json:"information"`
//}
//
//type DealContentRequestLog struct {
//	ID               int64     `gorm:"primaryKey"` // auto increment
//	NodeInfo         string    `json:"node_info"`
//	RequestingApiKey string    `json:"requesting_api_key"`
//	EventMessage     string    `json:"event_message"`
//	CreatedAt        time.Time `json:"created_at"` // auto set
//	UpdatedAt        time.Time `json:"updated_at"`
//}
//
//type DealPieceCommitmentRequestLog struct {
//	ID               int64     `gorm:"primaryKey"` // auto increment
//	NodeInfo         string    `json:"node_info"`
//	RequestingApiKey string    `json:"requesting_api_key"`
//	EventMessage     string    `json:"event_message"`
//	CreatedAt        time.Time `json:"created_at"` // auto set
//	UpdatedAt        time.Time `json:"updated_at"`
//}
//
//type ContentPrepareLog struct {
//	ID                  int64               `gorm:"primaryKey"` // auto increment
//	NodeInfo            string              `json:"node_info"`
//	RequestingApiKey    string              `json:"requesting_api_key"`
//	EventMessage        string              `json:"event_message"`
//	ContentDealProposal ContentDealProposal `json:"content_deal_proposal"`
//	CreatedAt           time.Time           `json:"created_at"` // auto set
//	UpdatedAt           time.Time           `json:"updated_at"`
//}
//type ContentAnnounceLog struct {
//	ID                  int64               `gorm:"primaryKey"` // auto increment
//	NodeInfo            string              `json:"node_info"`
//	RequestingApiKey    string              `json:"requesting_api_key"`
//	EventMessage        string              `json:"event_message"`
//	ContentDealProposal ContentDealProposal `json:"content_deal_proposal"`
//	CreatedAt           time.Time           `json:"created_at"` // auto set
//	UpdatedAt           time.Time           `json:"updated_at"`
//}
//
//type OpenStatsLog struct {
//	ID            int64     `gorm:"primaryKey"` // auto increment
//	NodeInfo      string    `json:"node_info"`
//	RequesterInfo string    `json:"requester_info"`
//	CreatedAt     time.Time `json:"created_at"` // auto set
//	UpdatedAt     time.Time `json:"updated_at"`
//}
//
//type RepairRequestLog struct {
//	ID               int64     `gorm:"primaryKey"` // auto increment
//	NodeInfo         string    `json:"node_info"`
//	RequesterInfo    string    `json:"requester_info"`
//	RequestingApiKey string    `json:"requesting_api_key"`
//	EventMessage     string    `json:"event_message"`
//	CreatedAt        time.Time `json:"created_at"` // auto set
//	UpdatedAt        time.Time `json:"updated_at"`
//}
//
//type NodeRequestLog struct {
//	ID               int64     `gorm:"primaryKey"` // auto increment
//	NodeInfo         string    `json:"node_info"`
//	RequesterInfo    string    `json:"requester_info"`
//	RequestingApiKey string    `json:"requesting_api_key"`
//	CreatedAt        time.Time `json:"created_at"` // auto set
//	UpdatedAt        time.Time `json:"updated_at"`
//}
//
//// job events
//type PieceCommitmentJobLog struct {
//	ID               int64     `gorm:"primaryKey"` // auto increment
//	NodeInfo         string    `json:"node_info"`
//	RequesterInfo    string    `json:"requester_info"`
//	RequestingApiKey string    `json:"requesting_api_key"`
//	CreatedAt        time.Time `json:"created_at"` // auto set
//	UpdatedAt        time.Time `json:"updated_at"`
//}
//
//type StorageDealMakeJobLog struct {
//	ID               int64     `gorm:"primaryKey"` // auto increment
//	NodeInfo         string    `json:"node_info"`
//	RequesterInfo    string    `json:"requester_info"`
//	RequestingApiKey string    `json:"requesting_api_key"`
//	CreatedAt        time.Time `json:"created_at"` // auto set
//	UpdatedAt        time.Time `json:"updated_at"`
//}
//
//type DataTransferStatusJobLog struct {
//	ID               int64     `gorm:"primaryKey"` // auto increment
//	NodeInfo         string    `json:"node_info"`
//	RequesterInfo    string    `json:"requester_info"`
//	RequestingApiKey string    `json:"requesting_api_key"`
//	CreatedAt        time.Time `json:"created_at"` // auto set
//	UpdatedAt        time.Time `json:"updated_at"`
//}
//
//type InstanceMetaJobLog struct {
//	ID               int64
//	NodeInfo         string    `json:"node_info"`
//	RequesterInfo    string    `json:"requester_info"`
//	RequestingApiKey string    `json:"requesting_api_key"`
//	DeltaNodeUuid    string    `json:"delta_node_uuid"`
//	CreatedAt        time.Time `json:"created_at"` // auto set
//	UpdatedAt        time.Time `json:"updated_at"`
//}
//
//// ContentLog time series log events
//type ContentLog struct {
//	ID                int64     `gorm:"primaryKey"` // auto increment
//	Name              string    `json:"name"`
//	Size              int64     `json:"size"`
//	Cid               string    `json:"cid"`
//	RequestingApiKey  string    `json:"requesting_api_key,omitempty"`
//	PieceCommitmentId int64     `json:"piece_commitment_id,omitempty"`
//	Status            string    `json:"status"`
//	ConnectionMode    string    `json:"connection_mode"` // offline or online
//	AutoRetry         bool      `json:"auto_retry"`
//	LastMessage       string    `json:"last_message"`
//	NodeInfo          string    `json:"node_info"`
//	RequesterInfo     string    `json:"requester_info"`
//	DeltaNodeUuid     string    `json:"delta_node_uuid"`
//	SystemContentId   int64     `json:"system_content_id"`
//	CreatedAt         time.Time `json:"created_at"` // auto set
//	UpdatedAt         time.Time `json:"updated_at"`
//}
//
//// ContentDealLog time series content deal events
//type ContentDealLog struct {
//	ID                  int64     `gorm:"primaryKey"`
//	Content             int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
//	PropCid             string    `json:"propCid"`
//	DealUUID            string    `json:"dealUuid"`
//	Miner               string    `json:"miner"`
//	DealID              int64     `json:"dealId"`
//	Failed              bool      `json:"failed"`
//	Verified            bool      `json:"verified"`
//	Slashed             bool      `json:"slashed"`
//	FailedAt            time.Time `json:"failedAt,omitempty"`
//	DTChan              string    `json:"dtChan" gorm:"index"`
//	TransferStarted     time.Time `json:"transferStarted"`
//	TransferFinished    time.Time `json:"transferFinished"`
//	OnChainAt           time.Time `json:"onChainAt"`
//	SealedAt            time.Time `json:"sealedAt"`
//	LastMessage         string    `json:"lastMessage"`
//	DealProtocolVersion string    `json:"deal_protocol_version"`
//	MinerVersion        string    `json:"miner_version,omitempty"`
//	NodeInfo            string    `json:"node_info"`
//	RequesterInfo       string    `json:"requester_info"`
//	RequestingApiKey    string    `json:"requesting_api_key"`
//	DeltaNodeUuid       string    `json:"delta_node_uuid"`
//	SystemContentDealId int64     `json:"system_content_deal_id"`
//	CreatedAt           time.Time `json:"created_at"` // auto set
//	UpdatedAt           time.Time `json:"updated_at"`
//}
//
//type ContentMinerLog struct {
//	ID                   int64     `gorm:"primaryKey"`
//	Content              int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
//	Miner                string    `json:"miner"`
//	NodeInfo             string    `json:"node_info"`
//	RequesterInfo        string    `json:"requester_info"`
//	RequestingApiKey     string    `json:"requesting_api_key"`
//	DeltaNodeUuid        string    `json:"delta_node_uuid"`
//	SystemContentMinerId int64     `json:"system_content_miner_id"`
//	CreatedAt            time.Time `json:"created_at"` // auto set
//	UpdatedAt            time.Time `json:"updated_at"`
//}
//
//type ContentWalletLog struct {
//	ID                    int64     `gorm:"primaryKey"`
//	Content               int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
//	WalletId              int64     `json:"wallet_id" gorm:"index:,option:CONCURRENTLY"`
//	NodeInfo              string    `json:"node_info"`
//	RequesterInfo         string    `json:"requester_info"`
//	RequestingApiKey      string    `json:"requesting_api_key"`
//	DeltaNodeUuid         string    `json:"delta_node_uuid"`
//	SystemContentWalletId int64     `json:"system_content_miner_id"`
//	CreatedAt             time.Time `json:"created_at"` // auto set
//	UpdatedAt             time.Time `json:"updated_at"`
//}
//
//type PieceCommitmentLog struct {
//	ID                             int64     `gorm:"primaryKey"`
//	Cid                            string    `json:"cid"`
//	Piece                          string    `json:"piece"`
//	Size                           int64     `json:"size"`
//	PaddedPieceSize                uint64    `json:"padded_piece_size"`
//	UnPaddedPieceSize              uint64    `json:"unnpadded_piece_size"`
//	Status                         string    `json:"status"` // open, in-progress, completed (closed).
//	LastMessage                    string    `json:"last_message"`
//	NodeInfo                       string    `json:"node_info"`
//	RequesterInfo                  string    `json:"requester_info"`
//	RequestingApiKey               string    `json:"requesting_api_key"`
//	DeltaNodeUuid                  string    `json:"delta_node_uuid"`
//	SystemContentPieceCommitmentId int64     `json:"system_content_piece_commitment_id"`
//	CreatedAt                      time.Time `json:"created_at"` // auto set
//	UpdatedAt                      time.Time `json:"updated_at"`
//}
//
//type ContentDealProposalLog struct {
//	ID                          int64     `gorm:"primaryKey"`
//	Content                     int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
//	Unsigned                    string    `json:"unsigned"`
//	Signed                      string    `json:"signed"`
//	Meta                        string    `json:"meta"`
//	NodeInfo                    string    `json:"node_info"`
//	RequesterInfo               string    `json:"requester_info"`
//	RequestingApiKey            string    `json:"requesting_api_key"`
//	DeltaNodeUuid               string    `json:"delta_node_uuid"`
//	SystemContentDealProposalId int64     `json:"system_content_deal_proposal_id"`
//	CreatedAt                   time.Time `json:"created_at"` // auto set
//	UpdatedAt                   time.Time `json:"updated_at"`
//}
//
//type ContentDealProposalParametersLog struct {
//	ID                                    int64     `gorm:"primaryKey"`
//	Content                               int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
//	Label                                 string    `json:"label,omitempty"`
//	Duration                              int64     `json:"duration,omitempty"`
//	StartEpoch                            int64     `json:"start_epoch,omitempty"`
//	EndEpoch                              int64     `json:"end_epoch,omitempty"`
//	TransferParams                        string    `json:"transfer_params,omitempty"`
//	RemoveUnsealedCopy                    bool      `json:"remove_unsealed_copy"`
//	SkipIPNIAnnounce                      bool      `json:"skip_ipni_announce"`
//	VerifiedDeal                          bool      `json:"verified_deal"`
//	UnverifiedDealMaxPrice                string    `json:"unverified_deal_max_price,omitempty"`
//	NodeInfo                              string    `json:"node_info"`
//	RequesterInfo                         string    `json:"requester_info"`
//	RequestingApiKey                      string    `json:"requesting_api_key"`
//	DeltaNodeUuid                         string    `json:"delta_node_uuid"`
//	SystemContentDealProposalParametersId int64     `json:"system_content_deal_proposal_id"`
//	CreatedAt                             time.Time `json:"created_at"` // auto set
//	UpdatedAt                             time.Time `json:"updated_at"`
//}
//
//type WalletLog struct {
//	ID               int64     `gorm:"primaryKey"`
//	UuId             string    `json:"uuid"`
//	Addr             string    `json:"addr"`
//	Owner            string    `json:"owner"`
//	KeyType          string    `json:"key_type"`
//	PrivateKey       string    `json:"private_key"`
//	NodeInfo         string    `json:"node_info"`
//	RequesterInfo    string    `json:"requester_info"`
//	RequestingApiKey string    `json:"requesting_api_key"`
//	DeltaNodeUuid    string    `json:"delta_node_uuid"`
//	SystemWalletId   int64     `json:"system_content_deal_proposal_id"`
//	CreatedAt        time.Time `json:"created_at"` // auto set
//	UpdatedAt        time.Time `json:"updated_at"`
//}
//
//type InstanceMetaLog struct {
//	// gorm id
//	ID                               int64     `gorm:"primary_key" json:"id"`
//	InstanceUuid                     string    `json:"instance_uuid"`
//	InstanceHostName                 string    `json:"instance_host_name"`
//	InstanceNodeName                 string    `json:"instance_node_name"`
//	OSDetails                        string    `json:"os_details"`
//	PublicIp                         string    `json:"public_ip"`
//	MemoryLimit                      uint64    `json:"memory_limit"`
//	CpuLimit                         uint64    `json:"cpu_limit"`
//	StorageLimit                     uint64    `json:"storage_limit"`
//	DisableRequest                   bool      `json:"disable_requests"`
//	DisableCommitmentPieceGeneration bool      `json:"disable_commitment_piece_generation"`
//	DisableStorageDeal               bool      `json:"disable_storage_deal"`
//	DisableOnlineDeals               bool      `json:"disable_online_deals"`
//	DisableOfflineDeals              bool      `json:"disable_offline_deals"`
//	NumberOfCpus                     uint64    `json:"number_of_cpus"`
//	StorageInBytes                   uint64    `json:"storage_in_bytes"`
//	SystemMemory                     uint64    `json:"system_memory"`
//	HeapMemory                       uint64    `json:"heap_memory"`
//	HeapInUse                        uint64    `json:"heap_in_use"`
//	StackInUse                       uint64    `json:"stack_in_use"`
//	InstanceStart                    time.Time `json:"instance_start"`
//	BytesPerCpu                      uint64    `json:"bytes_per_cpu"`
//	NodeInfo                         string    `json:"node_info"`
//	RequesterInfo                    string    `json:"requester_info"`
//	DeltaNodeUuid                    string    `json:"delta_node_uuid"`
//	SystemInstanceMetaId             int64     `json:"system_instance_meta_id"`
//	CreatedAt                        time.Time `json:"created_at"`
//	UpdatedAt                        time.Time `json:"updated_at"`
//}
//
//type BatchImportLog struct {
//	ID                  int64     `gorm:"primaryKey"`
//	Uuid                string    `json:"uuid" gorm:"index:,option:CONCURRENTLY"`
//	Status              string    `json:"status"`
//	NodeInfo            string    `json:"node_info"`
//	RequesterInfo       string    `json:"requester_info"`
//	DeltaNodeUuid       string    `json:"delta_node_uuid"`
//	SystemBatchImportId int64     `json:"system_batch_import_id"`
//	CreatedAt           time.Time `json:"created_at"`
//	UpdatedAt           time.Time `json:"updated_at"`
//}
//
//// BatchContent associate the content to a batch
//type BatchImportContentLog struct {
//	ID                   int64     `gorm:"primaryKey"`
//	BatchImportID        int64     `json:"batch_import_id" gorm:"index:,option:CONCURRENTLY"`
//	ContentID            int64     `json:"content_id" gorm:"index:,option:CONCURRENTLY"` // check status of the content
//	NodeInfo             string    `json:"node_info"`
//	RequesterInfo        string    `json:"requester_info"`
//	DeltaNodeUuid        string    `json:"delta_node_uuid"`
//	SystemBatchContentId int64     `json:"system_batch_content_id"`
//	CreatedAt            time.Time `json:"created_at"`
//	UpdatedAt            time.Time `json:"updated_at"`
//}
