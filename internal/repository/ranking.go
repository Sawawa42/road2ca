package repository

import (
	"github.com/redis/go-redis/v9"
	"road2ca/internal/entity"
	"context"
	"encoding/json"
	"fmt"
)

type RankingRepository interface {
	SaveToCache(ranking *entity.Ranking) error
	FindInRangeFromCache(start, end int) ([]*entity.Ranking, error)
}

type rankingRepository struct {
	rdb *redis.Client
}

func NewRankingRepository(rdb *redis.Client) RankingRepository {
	return &rankingRepository{
		rdb: rdb,
	}
}

// SaveToCache ランキング情報をキャッシュに保存する
func (r *rankingRepository) SaveToCache(ranking *entity.Ranking) error {
	ctx := context.Background()
	key := fmt.Sprintf("ranking:%d", ranking.UserID)
	jsonData, err := json.Marshal(ranking)
	if err != nil {
		return fmt.Errorf("failed to marshal ranking: %w", err)
	}

	if err := r.rdb.Set(ctx, key, jsonData, 0).Err(); err != nil {
		return fmt.Errorf("failed to save ranking to cache: %w", err)
	}
	return nil
}

// FindInRangeFromCache キャッシュから指定範囲のランキングを取得する
func (r *rankingRepository) FindInRangeFromCache(start, end int) ([]*entity.Ranking, error) {
	ctx := context.Background()
	keys, err := r.rdb.Keys(ctx, "ranking:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get ranking keys: %w", err)
	}

	var rankings []*entity.Ranking
	for _, key := range keys {
		val, err := r.rdb.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue // Key does not exist
			}
			return nil, fmt.Errorf("failed to get ranking from cache: %w", err)
		}

		var ranking entity.Ranking
		if err := json.Unmarshal([]byte(val), &ranking); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ranking: %w", err)
		}

		if ranking.Rank >= start && ranking.Rank <= end {
			rankings = append(rankings, &ranking)
		}
	}

	return rankings, nil
}
