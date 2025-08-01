package repository

import (
	"database/sql"
	"road2ca/internal/entity"
)

type MySQLSettingRepo interface {
	// Save SettingをDBに保存する
	Save(setting *entity.Setting) error
	// FindLatest 最新のsettingを取得する
	FindLatest() (*entity.Setting, error)
	// Truncate テーブルを空にする
	Truncate() error
}

type mysqlSettingRepo struct {
	db  *sql.DB
}

func NewMySQLSettingRepo(db *sql.DB) MySQLSettingRepo {
	return &mysqlSettingRepo{
		db: db,
	}
}

// Save SettingをDBに保存する
func (r *mysqlSettingRepo) Save(setting *entity.Setting) error {
	query := `
	INSERT INTO settings (id, name, gachaCoinConsumption, drawGachaMaxTimes, getRankingLimit, rewardCoin)
	VALUES (?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	name = VALUES(name),
	gachaCoinConsumption = VALUES(gachaCoinConsumption),
	drawGachaMaxTimes = VALUES(drawGachaMaxTimes),
	getRankingLimit = VALUES(getRankingLimit),
	rewardCoin = VALUES(rewardCoin)
	`
	uuidBytes, err := GetUUIDv7Bytes()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(query, uuidBytes, setting.Name, setting.GachaCoinConsumption, setting.DrawGachaMaxTimes, setting.GetRankingLimit, setting.RewardCoin)
	if err != nil {
		return err
	}
	return nil
}

// FindLatest 最新のsettingを取得する
func (r *mysqlSettingRepo) FindLatest() (*entity.Setting, error) {
	var setting entity.Setting
	query := "SELECT id, name, gachaCoinConsumption, drawGachaMaxTimes, getRankingLimit, rewardCoin FROM settings ORDER BY id DESC LIMIT 1"
	row := r.db.QueryRow(query)

	err := row.Scan(&setting.ID, &setting.Name, &setting.GachaCoinConsumption, &setting.DrawGachaMaxTimes, &setting.GetRankingLimit, &setting.RewardCoin)
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
