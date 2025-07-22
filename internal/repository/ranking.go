package repository

import (
	"context"
	"road2ca/internal/entity"
	"strconv"

	"github.com/redis/go-redis/v9"
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
		// 同一スコアの場合IDの昇順にするため、Scoreの少数部でIDを表現する
		Score:  float64(user.HighScore) + (1.0 - (float64(user.ID) / (1e12 + 1.0))),
		Member: user.ID,
	}).Err(); err != nil {
		return err
	}

	return nil
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
		userid, err := strconv.Atoi(score.Member.(string))
		if err != nil {
			return nil, err
		}
		results = append(results, &entity.Ranking{
			UserID: userid,
			Score:  int(score.Score),
			Rank:   start + i,
		})
	}
	return results, nil
}
