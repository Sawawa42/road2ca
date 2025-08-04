package repository

import (
	"database/sql"
	"road2ca/internal/entity"
)

type MySQLSettingRepo interface {
	// FindLatest 最新のsettingを取得する
	FindLatest() (*entity.Setting, error)
	// Truncate テーブルを空にする
	Truncate() error
}

type mysqlSettingRepo struct {
	db *sql.DB
}

func NewMySQLSettingRepo(db *sql.DB) MySQLSettingRepo {
	return &mysqlSettingRepo{
		db: db,
	}
}

// FindLatest 最新のsettingを取得する
func (r *mysqlSettingRepo) FindLatest() (*entity.Setting, error) {
	var setting entity.Setting
	query := "SELECT id, name, gachaCoinConsumption, drawGachaMaxTimes, getRankingLimit, rewardCoin, rarity3Ratio, rarity2Ratio, rarity1Ratio FROM settings ORDER BY id DESC LIMIT 1"
	row := r.db.QueryRow(query)

	err := row.Scan(&setting.ID, &setting.Name, &setting.GachaCoinConsumption, &setting.DrawGachaMaxTimes, &setting.GetRankingLimit, &setting.RewardCoin, &setting.Rarity3Ratio, &setting.Rarity2Ratio, &setting.Rarity1Ratio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No settings found
		}
		return nil, err
	}
	return &setting, nil
}

// Truncate テーブルを空にする
func (r *mysqlSettingRepo) Truncate() error {
	query := "TRUNCATE TABLE settings"
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
