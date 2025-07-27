package entity

type Setting struct {
	ID                   []byte
	Name                 string
	GachaCoinConsumption int
	DrawGachaMaxTimes    int
	GetRankingLimit      int
	RewardCoin           int
}
