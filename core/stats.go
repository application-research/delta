package core

type StatsService struct {
	DeltaNode DeltaNode
}

type StatsParam struct {
	RequestingApiKey string `json:"requesting_api_key"`
}
type CommpStatsParam struct {
	StatsParam
	PieceCommpId int64 `json:"piece_commp_id"`
}
type ContentStatsParam struct {
	StatsParam
	ContentId int64 `json:"content_id"`
}
type DealStatsParam struct {
	StatsParam
	DealId int64 `json:"deal_id"`
}

type StatsResult struct {
	Content          []Content         `json:"content"`
	Deals            []ContentDeal     `json:"deals"`
	PieceCommitments []PieceCommitment `json:"piece_commitments"`
}

type StatsContentResult struct {
	Content Content `json:"content"`
}

type StatsDealResult struct {
	Deals ContentDeal `json:"deals"`
}

type StatsPieceCommitmentResult struct {
	PieceCommitments PieceCommitment `json:"piece_commitments"`
}

func NewStatsStatsService() *StatsService {
	return &StatsService{}
}

func (s *StatsService) Status(param StatsParam) (StatsResult, error) {
	var content []Content
	s.DeltaNode.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ?", param.RequestingApiKey).Scan(&content)

	var contentDeal []ContentDeal
	s.DeltaNode.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ?", param.RequestingApiKey).Scan(&contentDeal)

	// select * from piece_commitments pc, content c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?;
	var pieceCommitments []PieceCommitment
	s.DeltaNode.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?", param.RequestingApiKey).Scan(&pieceCommitments)

	return StatsResult{
		Content:          content,
		Deals:            contentDeal,
		PieceCommitments: pieceCommitments}, nil
}

func (s *StatsService) CommpStatus(param CommpStatsParam) (StatsPieceCommitmentResult, error) {
	// select * from piece_commitments pc, content c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?;
	var pieceCommitment PieceCommitment
	s.DeltaNode.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.requesting_api_key = ? and pc.id = ?", param.RequestingApiKey, param.PieceCommpId).Scan(&pieceCommitment)

	return StatsPieceCommitmentResult{
		PieceCommitments: pieceCommitment}, nil
}

func (s *StatsService) ContentStatus(param ContentStatsParam) (StatsContentResult, error) {
	var content Content
	s.DeltaNode.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ? and c.id = ?", param.RequestingApiKey, param.ContentId).Scan(&content)

	return StatsContentResult{Content: content}, nil
}

func (s *StatsService) DealStatus(param DealStatsParam) (StatsDealResult, error) {
	var contentDeal ContentDeal
	s.DeltaNode.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ? and cd.id = ?", param.RequestingApiKey, param.DealId).Scan(&contentDeal)

	return StatsDealResult{Deals: contentDeal}, nil
}
