package entity

type Setting struct {
	ID                   int
	Name                 string
	GachaCoinConsumption int
	DrawGachaMaxTimes    int
	GetRankingLimit      int
	RewardCoin           int
	Rarity3Ratio         float64
	Rarity2Ratio         float64
	Rarity1Ratio         float64
}
