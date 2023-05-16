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

// GlobalDeltaDataReporter the global metrics tracer
// helps us improve the product.
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

	BATCH_IMPORT_STATUS_COMPLETED = "completed"
	BATCH_IMPORT_STATUS_FAILED    = "failed"
	BATCH_IMPORT_STATUS_STARTED   = "started"

	COMMP_STATUS_OPEN     = "open"
	COMMP_STATUS_COMITTED = "committed"

	CONNECTION_MODE_E2E    = "e2e"
	CONNECTION_MODE_IMPORT = "import"

	DEAL_VERIFIED   = "verified"
	DEAL_UNVERIFIED = "unverified"

	EPOCH_540_DAYS              = 1555200
	EPOCH_PER_DAY               = 2880
	EPOCH_PER_HOUR              = 60 * 2
	FILECOIN_GENESIS_UNIX_EPOCH = 1598306400
	DEFAULT_DURATION            = EPOCH_540_DAYS - (EPOCH_PER_DAY * 21)

	COMMP_MODE_FAST     = "fast"
	COMMP_MODE_STREAM   = "stream"
	COMPP_MODE_FILBOOST = "filboost"

	MAX_DEAL_RETRY = 10
)
