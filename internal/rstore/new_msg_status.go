package rstore

import (
	"fmt"

	"github.com/go-redis/redis"

	redis2 "github.com/1024casts/banhui/pkg/redis"
	"github.com/1024casts/snake/pkg/log"
)

// 记录用户新消息状态，也就是是否有新消息

const (
	DefaultExpireTime2 = 86400

	PrefixNewMsgStatus = "im:new_msg_status:%d" // redis前缀key
)

var NewMsgStatus = NewNewMsgStatus()

type newMsgStatus struct {
	redisClient *redis.Client
}

func NewNewMsgStatus() *newMsgStatus {
	return &newMsgStatus{
		redisClient: redis2.Client,
	}
}

func (n *newMsgStatus) Set(userId uint64) bool {
	redisKey := fmt.Sprintf(PrefixNewMsgStatus, userId)
	err := n.redisClient.Set(redisKey, 1, DefaultExpireTime2).Err()
	if err != nil {
		log.Warnf("[user_new_msg] set err, %v", err)
		return false
	}

	return true
}

func (n *newMsgStatus) Get(userId uint64) (string, error) {
	redisKey := fmt.Sprintf(PrefixNewMsgStatus, userId)
	val, err := n.redisClient.Get(redisKey).Result()
	if err == redis.Nil {
		return "", err
	} else if err != nil {
		log.Warnf("[user_new_msg] get err, %v", err)
		return "", err
	} else {
		return val, nil
	}
}

func (n *newMsgStatus) Del(userId uint64) int64 {
	redisKey := fmt.Sprintf(PrefixNewMsgStatus, userId)
	rows, err := n.redisClient.Del(redisKey).Result()
	if err != nil {
		log.Warnf("[user_new_msg] del err, %v", err)
		return rows
	}
	return rows
}
