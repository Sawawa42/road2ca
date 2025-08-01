package seed

import (
	"encoding/csv"
	"os"
	"road2ca/internal/entity"
	"strconv"
	"road2ca/internal/repository"
)

// Seed SettingとItemをCSVから読み込み、DBに保存する。保存前にテーブルを空にする。
func Seed(r *repository.Repositories) error {
	settings, err := loadSettingsFromCSV("internal/seed/csv/settings.csv")
	if err != nil {
		return err
	}

	if err := r.MySQLSetting.Truncate(); err != nil {
		return err
	}

	for _, setting := range settings {
		if err := r.MySQLSetting.Save(setting); err != nil {
			return err
		}
	}

	items, err := loadItemsFromCSV("internal/seed/csv/items.csv")
	if err != nil {
		return err
	}

	if err := r.Collection.Truncate(); err != nil {
		return err
	}

	if err := r.MySQLItem.Truncate(); err != nil {
		return err
	}

	if err := r.MySQLItem.Save(items); err != nil {
		return err
	}

	return nil
}

func loadSettingsFromCSV(filePath string) ([]*entity.Setting, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // ヘッダーを読み飛ばす
	if err != nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var settings []*entity.Setting
	for _, record := range records {
		GachaCoinConsumption, _ := strconv.Atoi(record[1])
		DrawGachaMaxTimes, _ := strconv.Atoi(record[2])
		GetRankingLimit, _ := strconv.Atoi(record[3])
		RewardCoin, _ := strconv.Atoi(record[4])

		setting := &entity.Setting{
			Name:                   record[0],
			GachaCoinConsumption:   GachaCoinConsumption,
			DrawGachaMaxTimes:      DrawGachaMaxTimes,
			GetRankingLimit:        GetRankingLimit,
			RewardCoin:             RewardCoin,
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

func loadItemsFromCSV(filePath string) ([]*entity.Item, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // ヘッダーを読み飛ばす
	if err != nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var items []*entity.Item
	for _, record := range records {
		rarity, _ := strconv.Atoi(record[1])
		weight, _ := strconv.Atoi(record[2])

		item := &entity.Item{
			Name:   record[0],
			Rarity: rarity,
			Weight: weight,
		}
		items = append(items, item)
	}

	return items, nil
}
