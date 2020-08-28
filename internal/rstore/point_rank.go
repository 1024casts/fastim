package rstore

import (
	"strconv"

	redis2 "github.com/1024casts/banhui/pkg/redis"
	"github.com/go-redis/redis"
	"github.com/lexkong/log"
)

// 用户相关计数
// 使用hash table来实现

const (
	PrefixPointRank = "point:rank"
)

var PointRank = NewPointRank()

type pointRank struct {
	redisClient *redis.Client
}

func NewPointRank() *pointRank {
	return &pointRank{
		redisClient: redis2.Client,
	}
}

// 返回新值
func (s *pointRank) ZIncrBy(userId uint64, point int) (float64, error) {
	redisKey := PrefixPointRank
	userIdStr := strconv.Itoa(int(userId))

	pointStr := strconv.Itoa(point)
	pointFloat, _ := strconv.ParseFloat(pointStr, 64)
	newValue, err := s.redisClient.ZIncrBy(redisKey, pointFloat, userIdStr).Result()
	if err != nil {
		log.Warnf("[rstore] ZIncrBy err, %+v", err)
		return newValue, err
	}

	return newValue, nil
}

// 获取排行列表
func (s *pointRank) GetRankList(limit int64) ([]redis.Z, error) {
	redisKey := PrefixPointRank

	newValue, err := s.redisClient.ZRevRangeByScoreWithScores(redisKey, redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  limit,
	}).Result()
	if err != nil {
		log.Warnf("[rstore] ZRangeByScoreWithScores err, %+v", err)
		return newValue, err
	}

	return newValue, nil
}
