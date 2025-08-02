package repository

import (
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"road2ca/internal/entity"
)

func FindItems(mysqlRepo MySQLItemRepo, redisRepo RedisItemRepo) ([]*entity.Item, error) {
	items, err := redisRepo.Find()
	if err != nil || len(items) == 0 {
		// キャッシュにアイテムがない場合はMySQLから取得
		if err == redis.Nil || len(items) == 0 {
			items, err = mysqlRepo.Find()
			if err != nil {
				return nil, err
			}
			return items, nil
		} else {
			return nil, err
		}
	}
	return items, nil
}

func GetUUIDv7Bytes() ([]byte, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return uuid.MarshalBinary()
}

func FindSetting(mysqlRepo MySQLSettingRepo, redisRepo RedisSettingRepo) (*entity.Setting, error) {
	setting, err := redisRepo.FindLatest()
	if err != nil || setting == nil {
		// キャッシュに設定がない場合はMySQLから取得
		if err == redis.Nil || setting == nil {
			setting, err = mysqlRepo.FindLatest()
			if err != nil {
				return nil, err
			}
			return setting, nil
		} else {
			return nil, err
		}
	}
	return setting, nil
}
