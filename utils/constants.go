package utils

// status
const (
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

	CONNECTION_MODE_ONLINE  = "online"
	CONNECTION_MODE_OFFLINE = "offline"

	LOTUS_API = "http://api.chain.love"

	DEFAULT_DURATION = 1555200 - (2880 * 21)
)
