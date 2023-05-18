package models

import model "github.com/application-research/delta-db/db_models"

type CidRequest struct {
	Cids []string `json:"cids"`
}

type WalletRequest struct {
	Id         uint64 `json:"id,omitempty"`
	Address    string `json:"address,omitempty"`
	Uuid       string `json:"uuid,omitempty"`
	KeyType    string `json:"key_type,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

type PieceCommitmentRequest struct {
	Piece             string `json:"piece_cid,omitempty"`
	PaddedPieceSize   uint64 `json:"padded_piece_size,omitempty"`
	UnPaddedPieceSize uint64 `json:"unpadded_piece_size,omitempty"`
}

type TransferParameters struct {
	URL     string      `json:"url,omitempty"`
	Headers interface{} `json:"headers,omitempty"`
}

type DealRequest struct {
	Cid                    string                 `json:"cid,omitempty"`
	Miner                  string                 `json:"miner,omitempty"`
	Duration               int64                  `json:"duration,omitempty"`
	DurationInDays         int64                  `json:"duration_in_days,omitempty"`
	Wallet                 WalletRequest          `json:"wallet,omitempty"`
	PieceCommitment        PieceCommitmentRequest `json:"piece_commitment,omitempty"`
	TransferParameters     TransferParameters     `json:"transfer_parameters,omitempty"`
	ConnectionMode         string                 `json:"connection_mode,omitempty"`
	Size                   int64                  `json:"size,omitempty"`
	StartEpoch             int64                  `json:"start_epoch,omitempty"`
	StartEpochInDays       int64                  `json:"start_epoch_in_days,omitempty"`
	Replication            int                    `json:"replication,omitempty"`
	RemoveUnsealedCopy     bool                   `json:"remove_unsealed_copy"`
	SkipIPNIAnnounce       bool                   `json:"skip_ipni_announce"`
	AutoRetry              bool                   `json:"auto_retry"`
	Label                  string                 `json:"label,omitempty"`
	DealVerifyState        string                 `json:"deal_verify_state,omitempty"`
	UnverifiedDealMaxPrice string                 `json:"unverified_deal_max_price,omitempty"`
}

// DealResponse Creating a new struct called DealResponse and then returning it.
type DealResponse struct {
	Status                       string         `json:"status"`
	Message                      string         `json:"message"`
	ContentId                    int64          `json:"content_id,omitempty"`
	DealRequest                  interface{}    `json:"deal_request_meta,omitempty"`
	DealProposalParameterRequest interface{}    `json:"deal_proposal_parameter_request_meta,omitempty"`
	ReplicatedContents           []DealResponse `json:"replicated_contents,omitempty"`
}

type DealReplication struct {
	Content                      model.Content                       `json:"content"`
	ContentDealProposalParameter model.ContentDealProposalParameters `json:"deal_proposal_parameter"`
	DealRequest                  DealRequest                         `json:"deal_request"`
}
