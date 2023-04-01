// A package that is used to define the constants used in the project.
package utils

import "github.com/application-research/delta-db/messaging"

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

var GlobalDeltaDataReporter = messaging.NewDeltaMetricsTracer()

const (
	DELTA_LABEL               string = "seal-the-delta-deal"
	CONTENT_PINNED            string = "pinned"
	CONTENT_FAILED_TO_PIN     string = "failed-to-pin"
	CONTENT_FAILED_TO_PROCESS string = "failed-to-process"

	CONTENT_PIECE_COMPUTING        = "piece-computing"
	CONTENT_PIECE_COMPUTED         = "piece-computed"
	CONTENT_PIECE_COMPUTING_FAILED = "piece-computing-failed"
	CONTENT_PIECE_ASSIGNED         = "piece-assigned"

	CONTENT_DEAL_MAKING_PROPOSAL  = "making-deal-proposal"
	CONTENT_DEAL_SENDING_PROPOSAL = "sending-deal-proposal"
	CONTENT_DEAL_PROPOSAL_SENT    = "deal-proposal-sent"
	CONTENT_DEAL_PROPOSAL_FAILED  = "deal-proposal-failed"

	DEAL_STATUS_TRANSFER_STARTED  = "transfer-started"
	DEAL_STATUS_TRANSFER_FINISHED = "transfer-finished"
	DEAL_STATUS_TRANSFER_FAILED   = "transfer-failed"

	COMMP_STATUS_OPEN     = "open"
	COMMP_STATUS_COMITTED = "committed"

	CONNECTION_MODE_E2E    = "e2e"
	CONNECTION_MODE_IMPORT = "import"

	DEAL_VERIFIED   = "verified"
	DEAL_UNVERIFIED = "unverified"

	LOTUS_API        = "http://api.chain.love"
	API_AUTH         = "https://auth.estuary.tech/check-api-key"
	SP_SELECTION_API = "https://sp-select.delta.store/api/providers"

	EPOCH_540_DAYS              = 1555200
	EPOCH_PER_DAY               = 2880
	EPOCH_PER_HOUR              = 60 * 2
	FILECOIN_GENESIS_UNIX_EPOCH = 1598306400
	DEFAULT_DURATION            = EPOCH_540_DAYS - (EPOCH_PER_DAY * 21)

	COMMP_MODE_FAST     = "fast"
	COMMP_MODE_STREAM   = "stream"
	COMPP_MODE_FILBOOST = "filboost"
)
