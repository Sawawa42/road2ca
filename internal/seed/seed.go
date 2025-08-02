package seed

import (
	"encoding/csv"
	"os"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"strconv"
	"fmt"
)

// Seed SettingとItemをCSVから読み込み、DBに保存する。保存前にテーブルを空にする。
func Seed(mysqlItem repository.MySQLItemRepo, mysqlSetting repository.MySQLSettingRepo, collectionRepo repository.CollectionRepo) error {
	settings, err := loadSettingsFromCSV("internal/seed/csv/settings.csv")
	if err != nil {
		return err
	}

	if err := mysqlSetting.Truncate(); err != nil {
		return err
	}

	for _, setting := range settings {
		if err := mysqlSetting.Save(setting); err != nil {
			return err
		}
	}

	setting, err := mysqlSetting.FindLatest()
	if err != nil {
		return err
	}

	items, err := loadItemsFromCSV("internal/seed/csv/items.csv")
	if err != nil {
		return err
	}

	if err := collectionRepo.Truncate(); err != nil {
		return err
	}

	if err := mysqlItem.Truncate(); err != nil {
		return err
	}

	if err := setWeightToItems(items, setting); err != nil {
		return err
	}

	if err := mysqlItem.Save(items); err != nil {
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
		Rarity3Ratio, _ := strconv.ParseFloat(record[5], 64)
		Rarity2Ratio, _ := strconv.ParseFloat(record[6], 64)
		Rarity1Ratio, _ := strconv.ParseFloat(record[7], 64)

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

func setWeightToItems(items []*entity.Item, setting *entity.Setting) error {
	if setting.Rarity3Ratio + setting.Rarity2Ratio + setting.Rarity1Ratio != 100 {
		return fmt.Errorf("total rarity ratio must be 100, got %f + %f + %f = %f",
			setting.Rarity3Ratio, setting.Rarity2Ratio, setting.Rarity1Ratio,
			setting.Rarity3Ratio + setting.Rarity2Ratio + setting.Rarity1Ratio)
	}

	totalRarityWeights := map[int]int{
		3: int(setting.Rarity3Ratio * 100000),
		2: int(setting.Rarity2Ratio * 100000),
		1: int(setting.Rarity1Ratio * 100000),
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

	return nil
}
