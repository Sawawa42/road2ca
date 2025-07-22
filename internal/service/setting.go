package service

import (
	"road2ca/internal/repository"
)

type GetSettingResponseDTO struct {
	GachaCoinConsumption int `json:"gachaCoinConsumption"`
	DrawGachaMaxTimes    int `json:"drawGachaMaxTimes"`
	GetRankingLimit      int `json:"getRankingLimit"`
	RewardCoin           int `json:"rewardCoin"`
}

type SettingService interface {
	SetSettingToCache() error
	GetSettings() (*GetSettingResponseDTO, error)
}

type settingService struct {
	settingRepo repository.SettingRepo
}

func NewSettingService(settingRepo repository.SettingRepo) SettingService {
	return &settingService{
		settingRepo: settingRepo,
	}
}

func (s *settingService) SetSettingToCache() error {
	setting, err := s.settingRepo.FindLatest()
	if err != nil {
		return err
	}

	if err := s.settingRepo.Save(setting); err != nil {
		return err
	}

	return nil
}

func (s *settingService) GetSettings() (*GetSettingResponseDTO, error) {
	setting, err := s.settingRepo.FindLatest()
	if err != nil {
		return nil, err
	}

	return &GetSettingResponseDTO{
		GachaCoinConsumption: setting.GachaCoinConsumption,
		DrawGachaMaxTimes:    setting.DrawGachaMaxTimes,
		GetRankingLimit:      setting.GetRankingLimit,
		RewardCoin:           setting.RewardCoin,
	}, nil
}
