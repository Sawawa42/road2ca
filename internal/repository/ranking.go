package repository

import (
	"context"
	"road2ca/internal/entity"
	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
)

type RankingRepo interface {
	Save(user *entity.User) error
	FindInRange(start, end int) ([]*entity.Ranking, error)
}

type rankingRepo struct {
	rdb *redis.Client
}

func NewRankingRepo(rdb *redis.Client) RankingRepo {
	return &rankingRepo{
		rdb: rdb,
	}
}

// Save ランキング情報をキャッシュに保存する
func (r *rankingRepo) Save(user *entity.User) error {
	ctx := context.Background()
	// sorted setを使用してランキングを保存する
	if err := r.rdb.ZAdd(ctx, "rankings", redis.Z{
		Score: calculateScore(user),
		Member: user.ID,
	}).Err(); err != nil {
		return err
	}

	return nil
}

func calculateScore(user *entity.User) float64 {
	// 同一スコアの場合IDの昇順にするため、Scoreの小数部にUUIDv7の時間情報を埋め込む
	sec, _ := user.ID.Time().UnixTime()
	return float64(user.HighScore) + (1.0 - (float64(sec) / (1e12 + 1.0)))
}

// FindInRange キャッシュから指定範囲のランキングを取得する
func (r *rankingRepo) FindInRange(start, end int) ([]*entity.Ranking, error) {
	ctx := context.Background()
	// startは1から始まるので、Redisのインデックスに合わせて-1している
	scores, err := r.rdb.ZRevRangeWithScores(ctx, "rankings", int64(start-1), int64(end-1)).Result()
	if err != nil {
		return nil, err
	}
	if len(scores) == 0 {
		return []*entity.Ranking{}, nil
	}

	results := make([]*entity.Ranking, 0, len(scores))
	for i, score := range scores {
		userid := score.Member.(uuid.UUID)
		results = append(results, &entity.Ranking{
			UserID: userid,
			Score:  int(score.Score),
			Rank:   start + i,
		})
	}
	return results, nil
}
