package rstore

import (
	"fmt"
	"strconv"
	"time"

	redis2 "github.com/1024casts/banhui/pkg/redis"
	"github.com/go-redis/redis"
	"github.com/lexkong/log"
)

// 用户相关计数
// 使用string来实现

const (

	// 用户机会次数
	TypeChanceUsedTimes  = "user:chance:%d:used_times"  // 已经使用的机会次数
	TypeChanceAwardTimes = "user:chance:%d:award_times" // 奖励的机会次数
)

var UserChance = NewUserChance()

type userChance struct {
	redisClient *redis.Client
}

func NewUserChance() *userChance {
	return &userChance{
		redisClient: redis2.Client,
	}
}

func (s *userChance) IncrUsedTimes(userId uint64) error {
	redisKey := fmt.Sprintf(TypeChanceUsedTimes, userId)
	_, err := s.redisClient.Incr(redisKey).Result()
	if err != nil {
		log.Warnf("[rstore] Incr err, %+v", err)
		return err
	}

	dateStr := time.Now().Format("2006-01-02")
	timeStr := dateStr + " 23:59:59"
	loc, _ := time.LoadLocation("Local")
	expireTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	if err != nil {
		log.Warnf("[rstore] ParseInLocation err, %+v", err)
		return err
	}
	err = s.redisClient.ExpireAt(redisKey, expireTime).Err()
	if err != nil {
		log.Warnf("[rstore] set expire err, %+v", err)
		return err
	}

	return nil
}

func (s *userChance) GetUsedTimes(userId uint64) (int, error) {
	redisKey := fmt.Sprintf(TypeChanceUsedTimes, userId)
	res, err := s.redisClient.Get(redisKey).Result()
	if err == redis2.Nil {
		return 0, nil
	} else if err != nil {
		log.Warnf("[rstore] Get err, %+v", err)
		return 0, err
	}

	val, err := strconv.Atoi(res)
	if err != nil {
		log.Warnf("[rstore] strconv atoi err, %+v", err)
		return 0, err
	}

	return val, nil
}

func (s *userChance) IncrAwardTimes(userId uint64) error {
	redisKey := fmt.Sprintf(TypeChanceAwardTimes, userId)
	_, err := s.redisClient.Incr(redisKey).Result()
	if err != nil {
		log.Warnf("[rstore] Incr err, %+v", err)
		return err
	}

	dateStr := time.Now().Format("2006-01-02")
	timeStr := dateStr + " 23:59:59"
	loc, _ := time.LoadLocation("Local")
	expireTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	if err != nil {
		log.Warnf("[rstore] ParseInLocation err, %+v", err)
		return err
	}
	err = s.redisClient.ExpireAt(redisKey, expireTime).Err()
	if err != nil {
		log.Warnf("[rstore] set expire err, %+v", err)
		return err
	}

	return nil
}

func (s *userChance) GetAwardTimes(userId uint64) (int, error) {
	redisKey := fmt.Sprintf(TypeChanceAwardTimes, userId)
	res, err := s.redisClient.Get(redisKey).Result()
	if err == redis2.Nil {
		return 0, nil
	} else if err != nil {
		log.Warnf("[rstore] Get err, %+v", err)
		return 0, err
	}

	val, err := strconv.Atoi(res)
	if err != nil {
		log.Warnf("[rstore] strconv atoi err, %+v", err)
		return 0, err
	}

	return val, nil
}
