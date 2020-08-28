package rstore

import (
	"fmt"
	"time"

	"github.com/1024casts/banhui/util"

	redis2 "github.com/1024casts/banhui/pkg/redis"
	"github.com/go-redis/redis"
	"github.com/lexkong/log"
)

// 用户相关计数
// 使用hash table来实现

const (
	DefaultExpireTime3 = 604800 * time.Second // 7天 86400*7

	PrefixUserStat = "user:stat:%d"

	DefaultIncrStep = 1  // 默认步长
	DefaultDecrStep = -1 // 默认步长

	// 通知类
	TypeNoticeComment = "new_notice_comment" // 新评论通知数
	TypeNoticeReply   = "new_notice_reply"   // 新回复通知数
	TypeNoticeLike    = "new_notice_like"    // 新喜欢通知数
	TypeNoticeSys     = "new_notice_sys"     // 新系统消息通知数
	TypeNoticeClass   = "new_notice_class"   // 新班级管理通知数

	// 个人信息类
	TypeUserPost      = "user_post_count"      // 个人帖子数
	TypeUserComment   = "user_comment_count"   // 个人评论数
	TypeUserReply     = "user_reply_count"     // 个人回复数
	TypeUserClass     = "user_class_count"     // 个人班级数
	TypeUserClassmate = "user_classmate_count" // 同学数

	// 班级数
	TypeClassPost    = "class_post_count"    // 班级帖子数
	TypeClassComment = "class_comment_count" // 班级评论数
	TypeClassReply   = "class_reply_count"   // 班级回复数

	// 用户喜欢数
	// 个人中心里的喜欢数就可以是下面3个值的总和
	TypeUserLikePost    = "like_post_count"    // 用户喜欢贴子数
	TypeUserLikeComment = "like_comment_count" // 用户喜欢评论数
	TypeUserLikeReply   = "like_reply_count"   // 用户喜欢回复数
)

var UserStat = NewUserStat()

type newUserStat struct {
	redisClient *redis.Client
}

func NewUserStat() *newUserStat {
	return &newUserStat{
		redisClient: redis2.Client,
	}
}

// 返回新值
func (s *newUserStat) HIncrBy(userId uint64, typ string) (int64, error) {
	redisKey := fmt.Sprintf(PrefixUserStat, userId)
	newValue, err := s.redisClient.HIncrBy(redisKey, typ, DefaultIncrStep).Result()
	if err != nil {
		log.Warnf("[rstore] HIncrBy err, %+v", err)
		return newValue, err
	}

	err = s.redisClient.Expire(redisKey, DefaultExpireTime3).Err()
	if err != nil {
		log.Warnf("[rstore] set expire err, %+v", err)
		return 0, err
	}

	return newValue, nil
}

func (s *newUserStat) HDecrBy(userId uint64, typ string) (int64, error) {
	redisKey := fmt.Sprintf(PrefixUserStat, userId)
	newValue, err := s.redisClient.HIncrBy(redisKey, typ, DefaultDecrStep).Result()
	if err != nil {
		log.Warnf("[rstore] HIncrBy err, %+v", err)
		return newValue, err
	}

	err = s.redisClient.Expire(redisKey, DefaultExpireTime3).Err()
	if err != nil {
		log.Warnf("[rstore] set expire err, %+v", err)
		return 0, err
	}

	return newValue, nil
}

// 1: 添加成功返回，0:已经存在或被取代
func (s *newUserStat) HSet(userId uint64, typ string, val string) (bool, error) {
	redisKey := fmt.Sprintf(PrefixUserStat, userId)
	res, err := s.redisClient.HSet(redisKey, typ, val).Result()
	log.Infof("hset res: %+v, err: %+v", res, err)
	if err != nil {
		log.Warnf("[rstore] HSet err, %+v", err)
		return false, err
	}

	err = s.redisClient.Expire(redisKey, DefaultExpireTime3).Err()
	if err != nil {
		log.Warnf("[rstore] set expire err, %+v", err)
		return false, err
	}

	return true, nil
}

// 哈希表不存在或者hash key 不存在返回错误
func (s *newUserStat) HGet(userId uint64, typ string) (int, error) {
	redisKey := fmt.Sprintf(PrefixUserStat, userId)
	res, err := s.redisClient.HGet(redisKey, typ).Result()
	if err != nil {
		log.Warnf("[rstore] HGet err, %+v", err)
		return 0, err
	}

	val, err := util.StringToInt(res)
	if err != nil {
		log.Warnf("[rstore] string to int err, %+v", err)
		return 0, err
	}

	return val, nil
}

// 返回删除的记录数
func (s *newUserStat) HDel(userId uint64, typ string) (int64, error) {
	redisKey := fmt.Sprintf(PrefixUserStat, userId)
	res, err := s.redisClient.HDel(redisKey, typ).Result()
	if err != nil {
		log.Warnf("[rstore] HDel err, %+v", err)
		return 0, err
	}

	return res, nil
}

func (s *newUserStat) HGetAll(userId uint64, typ string) (map[string]string, error) {
	redisKey := fmt.Sprintf(PrefixUserStat, userId)
	resMap, err := s.redisClient.HGetAll(redisKey).Result()
	if err != nil {
		log.Warnf("[rstore] HGetAll err, %+v", err)
		return nil, err
	}

	return resMap, nil
}
