package entity

type Setting struct {
	ID                   []byte
	Name                 string
	GachaCoinConsumption int
	DrawGachaMaxTimes    int
	GetRankingLimit      int
	RewardCoin           int
	Rarity3Ratio         int
	Rarity2Ratio         int
	Rarity1Ratio         int
}
