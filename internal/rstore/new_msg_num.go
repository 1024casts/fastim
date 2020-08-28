package rstore

import (
	"fmt"
	"strconv"

	"github.com/lexkong/log"

	redis2 "github.com/1024casts/banhui/pkg/redis"
	"github.com/go-redis/redis"
)

// 记录用户会话的新消息数

const (
	DefaultExpireTime = 604800 // 7天

	PrefixNewMsgNum = "im:new_msg_num:%d" // redis前缀key
)

var NewMsgNum = NewNewMsgNum()

type newMsgNum struct {
	redisClient *redis.Client
}

func NewNewMsgNum() *newMsgNum {
	return &newMsgNum{
		redisClient: redis2.Client,
	}
}

// 有序集合-增加分值
func (n *newMsgNum) Incr(userId, YUserId uint64) bool {
	redisKey := fmt.Sprintf(PrefixNewMsgNum, userId)
	// 字符串转uint64
	userIdStr := strconv.FormatUint(YUserId, 10)
	err := n.redisClient.ZIncrBy(redisKey, 1, userIdStr).Err()
	if err != nil {
		log.Warnf("[rstore] ZIncrBy err, %+v", err)
		return false
	}

	err = n.redisClient.Expire(redisKey, DefaultExpireTime).Err()
	if err != nil {
		log.Warnf("[rstore] set expire err, %+v", err)
		return false
	}

	return true
}

// 从有序集合中删除一条
func (n *newMsgNum) DelOne(userId, YUserId uint64) bool {
	redisKey := fmt.Sprintf(PrefixNewMsgNum, userId)
	userIdStr := strconv.FormatUint(YUserId, 10)
	err := n.redisClient.ZRem(redisKey, userIdStr).Err()
	if err != nil {
		log.Warnf("[rstore] ZRem err, %+v", err)
		return false
	}

	err = n.redisClient.Expire(redisKey, DefaultExpireTime).Err()
	if err != nil {
		log.Warnf("[rstore] set expire err, %+v", err)
		return false
	}

	return true
}

// 删除整个有序集合
func (n *newMsgNum) DelAll(userId uint64) (int64, error) {
	redisKey := fmt.Sprintf(PrefixNewMsgNum, userId)
	val, err := n.redisClient.Del(redisKey).Result()
	if err != nil {
		log.Warnf("[rstore] ZRem Del, %+v", err)
		return 0, err
	}

	return val, nil
}

// 获取整个有序集合
func (n *newMsgNum) GetAll(userId uint64) (map[uint64]int, error) {
	redisKey := fmt.Sprintf(PrefixNewMsgNum, userId)
	valSlice, err := n.redisClient.ZRangeWithScores(redisKey, 0, -1).Result()

	retMap := make(map[uint64]int)
	if err != nil {
		log.Warnf("[rstore] ZRem Del, %+v", err)
		return retMap, err
	}

	log.Infof("new msg num get all, result: %+v", valSlice)
	for _, val := range valSlice {
		retMap[val.Member.(uint64)] = int(val.Score)
	}
	log.Infof("new msg num get all, result map: %+v", retMap)

	return retMap, nil
}
