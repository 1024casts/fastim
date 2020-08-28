package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/api.v7/storage"

	"github.com/lexkong/log"

	"github.com/1024casts/banhui/pkg/constvar"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestId, ok := v.(string); ok {
		return requestId
	}
	return ""
}

func GetDate() string {
	return time.Now().Format("2006/01/02")
}

// 获取整形的日期
func GetTodayDateInt() int {
	dateStr := time.Now().Format("200601")
	date, err := strconv.Atoi(dateStr)
	if err != nil {
		log.Warnf("[util] get date int err, %+v", err)
		return 0
	}
	return date
}

func TimeLayout() string {
	return "2006-01-02 15:04:05"
}

func TimeToString(ts time.Time) string {
	return time.Unix(ts.Unix(), 00).Format(TimeLayout())
}

func TimeToShortString(ts time.Time) string {
	return time.Unix(ts.Unix(), 00).Format("2006.01.02")
}

func GetDefaultAvatarUrl() string {
	return GetQiNiuPublicAccessUrl(constvar.DefaultAvatar)
}

// user's avatar, if empty, use default avatar
func GetAvatarUrl(key string) string {
	if key == "" {
		return GetDefaultAvatarUrl()
	}
	if strings.HasPrefix(key, "https://") {
		return key
	}
	return GetQiNiuPublicAccessUrl(key)
}

func GetStaticImageUrl(key string) string {
	if key == "" {
		return GetDefaultAvatarUrl()
	}
	return GetQiNiuPublicAccessUrl(key)
}

// 格式化时间
func GetShowTime(ts time.Time) string {
	duration := time.Now().Unix() - ts.Unix()
	if duration < 60 {
		//return fmt.Sprintf("%d妙前", duration)
		return fmt.Sprintf("刚刚发布")
	} else if duration < 3600 {
		return fmt.Sprintf("%d分钟前更新", duration/60)
	} else if duration < 86400 {
		return fmt.Sprintf("%d小时前更新", duration/3600)
	} else if duration < 86400*2 {
		return fmt.Sprint("昨天更新")
	} else {
		return TimeToShortString(ts) + "前更新"
	}
}

// 字符串转md5
func Md5(str string) (string, error) {
	h := md5.New()

	_, err := io.WriteString(h, str)
	if err != nil {
		return "", err
	}

	// 注意：这里不能使用string将[]byte转为字符串，否则会显示乱码
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// 获取七牛资源的公有链接
// 无需配置bucket, 域名会自动到域名所绑定的bucket去查找
func GetQiNiuPublicAccessUrl(path string) string {
	domain := viper.GetString("qiniu.cdn_url")
	key := strings.TrimPrefix(path, "/")

	publicAccessURL := storage.MakePublicURL(domain, key)

	return publicAccessURL
}
