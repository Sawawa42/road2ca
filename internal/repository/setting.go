package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"road2ca/internal/entity"
)

type SettingRepo interface {
	Save(setting *entity.Setting) error
	FindLatest() (*entity.Setting, error)
}

type settingRepo struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewSettingRepo(db *sql.DB, rdb *redis.Client) SettingRepo {
	return &settingRepo{
		db:  db,
		rdb: rdb,
	}
}

// Save 設定情報をキャッシュする
func (r *settingRepo) Save(setting *entity.Setting) error {
	ctx := context.Background()
	key := "setting:latest"
	jsonData, err := json.Marshal(setting)
	if err != nil {
		return err
	}
	err = r.rdb.Set(ctx, key, jsonData, 0).Err()
	return err
}

// FindLatest 最新の設定情報を取得する。キャッシュに存在しない場合はDBから取得する
func (r *settingRepo) FindLatest() (*entity.Setting, error) {
	ctx := context.Background()
	key := "setting:latest"

	// キャッシュから取得
	jsonData, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// キャッシュに存在しない場合はDBから取得
		return r.findLatestFromDB()
	} else if err != nil {
		return nil, err
	}

	var setting entity.Setting
	err = json.Unmarshal([]byte(jsonData), &setting)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingRepo) findLatestFromDB() (*entity.Setting, error) {
	var setting entity.Setting
	query := "SELECT id, name, gachaCoinConsumption, drawGachaMaxTimes, getRankingLimit, rewardCoin FROM settings ORDER BY id DESC LIMIT 1"
	row := r.db.QueryRow(query)

	err := row.Scan(&setting.ID, &setting.Name, &setting.GachaCoinConsumption, &setting.DrawGachaMaxTimes, &setting.GetRankingLimit, &setting.RewardCoin)
	if err != nil {
		return nil, err
	}

	return &setting, nil
}
