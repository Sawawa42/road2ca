package service

import (
	"road2ca/internal/repository"
)

type GetSettingResponseDTO struct {
	GachaCoinConsumption int `json:"gachaCoinConsumption"`
	DrawGachaMaxTimes    int `json:"drawGachaMaxTimes"`
	GetRankingLimit      int `json:"getRankingLimit"`
	RewardCoin           int `json:"rewardCoin"`
	Rarity3Ratio         float64 `json:"rarity3Ratio"`
	Rarity2Ratio         float64 `json:"rarity2Ratio"`
	Rarity1Ratio         float64 `json:"rarity1Ratio"`
}

type SettingService interface {
	SetSettingToCache() error
	GetSettings() (*GetSettingResponseDTO, error)
}

type settingService struct {
	mysqlSettingRepo repository.MySQLSettingRepo
	redisSettingRepo repository.RedisSettingRepo
}

func NewSettingService(mysqlSettingRepo repository.MySQLSettingRepo, redisSettingRepo repository.RedisSettingRepo) SettingService {
	return &settingService{
		mysqlSettingRepo: mysqlSettingRepo,
		redisSettingRepo: redisSettingRepo,
	}
}

func (s *settingService) SetSettingToCache() error {
	setting, err := repository.FindSetting(s.mysqlSettingRepo, s.redisSettingRepo)
	if err != nil {
		return err
	}

	if err := s.redisSettingRepo.Save(setting); err != nil {
		return err
	}

	return nil
}

func (s *settingService) GetSettings() (*GetSettingResponseDTO, error) {
	setting, err := repository.FindSetting(s.mysqlSettingRepo, s.redisSettingRepo)
	if err != nil {
		return nil, err
	}

	return &GetSettingResponseDTO{
		GachaCoinConsumption: setting.GachaCoinConsumption,
		DrawGachaMaxTimes:    setting.DrawGachaMaxTimes,
		GetRankingLimit:      setting.GetRankingLimit,
		RewardCoin:           setting.RewardCoin,
		Rarity3Ratio:         setting.Rarity3Ratio,
		Rarity2Ratio:         setting.Rarity2Ratio,
		Rarity1Ratio:         setting.Rarity1Ratio,
	}, nil
}
