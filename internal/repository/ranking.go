package repository

import (
	"github.com/redis/go-redis/v9"
	"road2ca/internal/entity"
	"context"
)

type RankingRepository interface {
	SaveToCache(user *entity.User) error
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
// 保存時(/game/finish)、contextから取得したuserのhighscoreを更新したuserを引数に取る
func (r *rankingRepository) SaveToCache(user *entity.User) error {
	ctx := context.Background()
	// sorted setを使用してランキングを保存する
	if err := r.rdb.ZAdd(ctx, "rankings", redis.Z{
		// 同一スコアの場合IDの昇順にするため、Scoreの少数部でIDを表現する
		Score: float64(user.HighScore) + (1.0 - (float64(user.ID) / (1e12 + 1.0))),
		Member: user.ID,
	}).Err(); err != nil {
		return err
	}

	return nil
}

// FindInRangeFromCache キャッシュから指定範囲のランキングを取得する
func (r *rankingRepository) FindInRangeFromCache(start, end int) ([]*entity.Ranking, error) {
	ctx := context.Background()
	scores, err := r.rdb.ZRevRangeWithScores(ctx, "rankings", int64(start), int64(end)).Result()
	if err != nil {
		return nil, err
	}
	if len(scores) == 0 {
		return []*entity.Ranking{}, nil
	}

	results := make([]*entity.Ranking, 0, len(scores))
	for i, score := range scores {
		results = append(results, &entity.Ranking{
			UserID: int(score.Member.(int)),
			Score:  int(score.Score),
			Rank:   start + i + 1,
		})
	}

	return results, nil
}
