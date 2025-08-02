package seed

import (
	"encoding/csv"
	"os"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"strconv"
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

	setting, err := r.MySQLSetting.FindLatest()
	if err != nil {
		return err
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

	setWeightToItems(items, setting)

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
		Rarity3Ratio, _ := strconv.Atoi(record[5])
		Rarity2Ratio, _ := strconv.Atoi(record[6])
		Rarity1Ratio, _ := strconv.Atoi(record[7])

		setting := &entity.Setting{
			Name:                 record[0],
			GachaCoinConsumption: GachaCoinConsumption,
			DrawGachaMaxTimes:    DrawGachaMaxTimes,
			GetRankingLimit:      GetRankingLimit,
			RewardCoin:           RewardCoin,
			Rarity3Ratio:         Rarity3Ratio,
			Rarity2Ratio:         Rarity2Ratio,
			Rarity1Ratio:         Rarity1Ratio,
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

		item := &entity.Item{
			Name:   record[0],
			Rarity: rarity,
			Weight: 0,
		}
		items = append(items, item)
	}

	return items, nil
}

func setWeightToItems(items []*entity.Item, setting *entity.Setting) {
	totalRarityWeights := map[int]int{
		3: setting.Rarity3Ratio,
		2: setting.Rarity2Ratio,
		1: setting.Rarity1Ratio,
	}

	itemsByRarity := make(map[int][]*entity.Item)
	for _, item := range items {
		itemsByRarity[item.Rarity] = append(itemsByRarity[item.Rarity], item)
	}

	for rarity, itemsInGroup := range itemsByRarity {
		totalWeight := totalRarityWeights[rarity]
		count := len(itemsInGroup)

		if count == 0 {
			continue
		}

		baseWeight := totalWeight / count
		remainder := totalWeight % count

		for i, item := range itemsInGroup {
			item.Weight = baseWeight
			if i < remainder {
				item.Weight++
			}
		}
	}
}
