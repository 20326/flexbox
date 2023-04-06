package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// https://github.com/spark-golang/spark-url/database/redis/redis_command.go

type RedisStorage Storage

func (rc *RedisStorage) Get(key string) string {
	var mes string
	var cmd *redis.StringCmd

	cmd = rc.db.Get(context.Background(), key)

	if err := cmd.Err(); err != nil {
		mes = ""
	} else {
		mes = cmd.Val()
	}
	return mes
}

func (rc *RedisStorage) Set(key string, value interface{}, expire time.Duration) bool {
	var err error
	err = rc.db.Set(context.Background(), key, value, expire).Err()

	if err != nil {
		return false
	}
	return true
}

func (rc *RedisStorage) TTL(key string) time.Duration {
	return rc.db.TTL(context.Background(), key).Val()
}

func (rc *RedisStorage) GetRaw(key string) (bts []byte, err error) {
	bts, err = rc.db.Get(context.Background(), key).Bytes()

	if err != nil && err != redis.Nil {
		return []byte{}, err
	}
	return bts, nil
}

func (rc *RedisStorage) MGet(keys ...string) ([]string, error) {
	var sliceCmd *redis.SliceCmd
	sliceCmd = rc.db.MGet(context.Background(), keys...)

	if err := sliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	tmp := sliceCmd.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice, nil
}

func (rc *RedisStorage) MGets(keys ...string) (ret []interface{}, err error) {
	ret, err = rc.db.MGet(context.Background(), keys...).Result()

	if err != nil && err != redis.Nil {
		return []interface{}{}, err
	}
	return ret, nil
}

// HGetAll 从redis获取hash的所有键值对
func (rc *RedisStorage) HGetAll(key string) map[string]string {
	var hash map[string]string
	var stringMapCmd *redis.MapStringStringCmd
	stringMapCmd = rc.db.HGetAll(context.Background(), key)

	if err := stringMapCmd.Err(); err != nil && err != redis.Nil {
		hash = make(map[string]string)
	} else {
		hash = stringMapCmd.Val()
	}

	return hash
}

// HGet 从redis获取hash单个值
func (rc *RedisStorage) HGet(key string, fields string) (string, error) {
	var stringCmd *redis.StringCmd
	stringCmd = rc.db.HGet(context.Background(), key, fields)

	err := stringCmd.Err()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == redis.Nil {
		return "", nil
	}
	return stringCmd.Val(), nil
}

// HMGet 批量获取hash值
func (rc *RedisStorage) HMGet(key string, fileds ...string) []string {
	var sliceCmd *redis.SliceCmd
	sliceCmd = rc.db.HMGet(context.Background(), key, fileds...)

	if err := sliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	tmp := sliceCmd.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice
}

// HMGetMap 批量获取hash值，返回map
func (rc *RedisStorage) HMGetMap(key string, fields ...string) map[string]string {
	if len(fields) == 0 {
		return make(map[string]string)
	}

	var sliceCmd *redis.SliceCmd
	sliceCmd = rc.db.HMGet(context.Background(), key, fields...)

	if err := sliceCmd.Err(); err != nil && err != redis.Nil {
		return make(map[string]string)
	}

	tmp := sliceCmd.Val()
	hashRet := make(map[string]string, len(tmp))

	var tmpTagID string

	for k, v := range tmp {
		tmpTagID = fields[k]
		if v != nil {
			hashRet[tmpTagID] = v.(string)
		} else {
			hashRet[tmpTagID] = ""
		}
	}
	return hashRet
}

// HMSet 设置redis的hash
func (rc *RedisStorage) HMSet(key string, hash map[string]interface{}, expire time.Duration) bool {
	if len(hash) > 0 {
		var err error
		err = rc.db.HMSet(context.Background(), key, hash).Err()

		if err != nil {
			return false
		}
		rc.db.Expire(context.Background(), key, expire)
		return true
	}
	return false
}

// HSet hset
func (rc *RedisStorage) HSet(key string, field string, value interface{}) bool {
	var err error
	err = rc.db.HSet(context.Background(), key, field, value).Err()
	if err != nil {
		return false
	}
	return true
}

// HDel ...
func (rc *RedisStorage) HDel(key string, field ...string) bool {
	var intCmd *redis.IntCmd
	intCmd = rc.db.HDel(context.Background(), key, field...)

	if err := intCmd.Err(); err != nil {
		return false
	}

	return true
}

// SetWithErr ...
func (rc *RedisStorage) SetWithErr(key string, value interface{}, expire time.Duration) error {

	return rc.db.Set(context.Background(), key, value, expire).Err()
}

// SetNx 设置redis的string 如果键已存在
func (rc *RedisStorage) SetNx(key string, value interface{}, expiration time.Duration) bool {
	var res bool
	var err error
	res, err = rc.db.SetNX(context.Background(), key, value, expiration).Result()

	if err != nil {
		return false
	}

	return res
}

// SetNxWithErr 设置redis的string 如果键已存在
func (rc *RedisStorage) SetNxWithErr(key string, value interface{}, expiration time.Duration) (bool, error) {
	return rc.db.SetNX(context.Background(), key, value, expiration).Result()
}

// Incr redis自增
func (rc *RedisStorage) Incr(key string) bool {
	var err error
	err = rc.db.Incr(context.Background(), key).Err()

	if err != nil {
		return false
	}
	return true
}

// IncrWithErr ...
func (rc *RedisStorage) IncrWithErr(key string) (int64, error) {
	return rc.db.Incr(context.Background(), key).Result()
}

// IncrBy 将 key 所储存的值加上增量 increment 。
func (rc *RedisStorage) IncrBy(key string, increment int64) (int64, error) {
	var intCmd *redis.IntCmd
	intCmd = rc.db.IncrBy(context.Background(), key, increment)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// Decr redis自减
func (rc *RedisStorage) Decr(key string) bool {
	var err error
	err = rc.db.Decr(context.Background(), key).Err()

	if err != nil {
		return false
	}
	return true
}

// Type ...
func (rc *RedisStorage) Type(key string) (string, error) {
	var statusCmd *redis.StatusCmd
	statusCmd = rc.db.Type(context.Background(), key)

	if err := statusCmd.Err(); err != nil {
		return "", err
	}
	return statusCmd.Val(), nil
}

// ZRevRange 倒序获取有序集合的部分数据
func (rc *RedisStorage) ZRevRange(key string, start, stop int64) ([]string, error) {
	var stringSliceCmd *redis.StringSliceCmd

	stringSliceCmd = rc.db.ZRevRange(context.Background(), key, start, stop)
	if err := stringSliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return stringSliceCmd.Val(), nil
}

// ZRevRangeWithScores ...
func (rc *RedisStorage) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	var zSliceCmd *redis.ZSliceCmd
	zSliceCmd = rc.db.ZRevRangeWithScores(context.Background(), key, start, stop)

	if err := zSliceCmd.Err(); err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}
	return zSliceCmd.Val(), nil
}

// ZRange ...
func (rc *RedisStorage) ZRange(key string, start, stop int64) ([]string, error) {
	var stringSliceCmd *redis.StringSliceCmd
	stringSliceCmd = rc.db.ZRange(context.Background(), key, start, stop)

	if err := stringSliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return stringSliceCmd.Val(), nil
}

// ZRevRank ...
func (rc *RedisStorage) ZRevRank(key string, member string) (int64, error) {
	var intCmd *redis.IntCmd

	intCmd = rc.db.ZRevRank(context.Background(), key, member)

	if err := intCmd.Err(); err != nil && err != redis.Nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// ZRevRangeByScore ...
func (rc *RedisStorage) ZRevRangeByScore(key string, opt *redis.ZRangeBy) (res []string, err error) {

	res, err = rc.db.ZRevRangeByScore(context.Background(), key, opt).Result()

	if err != nil && err != redis.Nil {
		return []string{}, err
	}

	return res, nil
}

// ZRevRangeByScoreWithScores ...
// ZRevRange 倒序获取有序集合的部分数据
func (rc *RedisStorage) ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) (res []redis.Z, err error) {
	res, err = rc.db.ZRevRangeByScoreWithScores(context.Background(), key, opt).Result()

	if err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}

	return res, nil
}

// ZCard 获取有序集合的基数
func (rc *RedisStorage) nZCard(key string) (int64, error) {
	var intCmd *redis.IntCmd
	intCmd = rc.db.ZCard(context.Background(), key)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// ZScore 获取有序集合成员 member 的 score 值
func (rc *RedisStorage) ZScore(key string, member string) (float64, error) {
	var floatCmd *redis.FloatCmd
	floatCmd = rc.db.ZScore(context.Background(), key, member)

	err := floatCmd.Err()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	return floatCmd.Val(), err
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
func (rc *RedisStorage) ZAdd(key string, members ...redis.Z) (int64, error) {
	var intCmd *redis.IntCmd
	intCmd = rc.db.ZAdd(context.Background(), key, members...)

	if err := intCmd.Err(); err != nil && err != redis.Nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (rc *RedisStorage) ZCount(key string, min, max string) (int64, error) {
	var intCmd *redis.IntCmd
	intCmd = rc.db.ZCount(context.Background(), key, min, max)

	if err := intCmd.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return intCmd.Val(), nil
}

// ZIncrBy 有序集合中对指定成员的分数加上增量 increment
func (rc *RedisStorage) ZIncrBy(key, member string, increment float64) (float64, error) {
	var floatCmd *redis.FloatCmd
	floatCmd = rc.db.ZIncrBy(context.Background(), key, increment, member)
	if err := floatCmd.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return floatCmd.Val(), nil
}

// Del redis删除
func (rc *RedisStorage) Del(key string) int64 {
	var res int64
	var err error

	res, err = rc.db.Del(context.Background(), key).Result()

	if err != nil {
		return 0
	}
	return res
}

// DelWithErr ...
func (rc *RedisStorage) DelWithErr(key string) (int64, error) {
	return rc.db.Del(context.Background(), key).Result()
}

// HIncrBy 哈希field自增
func (rc *RedisStorage) HIncrBy(key string, field string, incr int) {
	rc.db.HIncrBy(context.Background(), key, field, int64(incr))
}

// Exists 键是否存在
func (rc *RedisStorage) Exists(key string) bool {
	var res int64
	var err error
	res, err = rc.db.Exists(context.Background(), key).Result()

	if err != nil {
		return false
	}
	return res == 1
}

// ExistsWithErr ...
func (rc *RedisStorage) ExistsWithErr(key string) (bool, error) {
	var res int64
	var err error
	res, err = rc.db.Exists(context.Background(), key).Result()

	if err != nil {
		return false, nil
	}
	return res == 1, nil
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
func (rc *RedisStorage) LPush(key string, values ...interface{}) (int64, error) {
	var intCmd *redis.IntCmd

	intCmd = rc.db.LPush(context.Background(), key, values...)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}

	return intCmd.Val(), nil
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func (rc *RedisStorage) RPush(key string, values ...interface{}) (int64, error) {
	var intCmd *redis.IntCmd
	intCmd = rc.db.RPush(context.Background(), key, values...)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}

	return intCmd.Val(), nil
}

// RPop 移除并返回列表 key 的尾元素。
func (rc *RedisStorage) RPop(key string) (string, error) {
	var stringCmd *redis.StringCmd
	stringCmd = rc.db.RPop(context.Background(), key)

	if err := stringCmd.Err(); err != nil {
		return "", err
	}

	return stringCmd.Val(), nil
}

// LRange 获取列表指定范围内的元素
func (rc *RedisStorage) LRange(key string, start, stop int64) (res []string, err error) {
	res, err = rc.db.LRange(context.Background(), key, start, stop).Result()

	if err != nil {
		return []string{}, err
	}

	return res, nil
}

// LLen ...
func (rc *RedisStorage) LLen(key string) int64 {
	intCmd := rc.db.LLen(context.Background(), key)

	if err := intCmd.Err(); err != nil {
		return 0
	}

	return intCmd.Val()
}

// LLenWithErr ...
func (rc *RedisStorage) LLenWithErr(key string) (int64, error) {

	return rc.db.LLen(context.Background(), key).Result()
}

// LRem ...
func (rc *RedisStorage) LRem(key string, count int64, value interface{}) int64 {
	var intCmd *redis.IntCmd
	intCmd = rc.db.LRem(context.Background(), key, count, value)

	if err := intCmd.Err(); err != nil {
		return 0
	}

	return intCmd.Val()
}

// LIndex ...
func (rc *RedisStorage) LIndex(key string, idx int64) (string, error) {

	return rc.db.LIndex(context.Background(), key, idx).Result()
}

// LTrim ...
func (rc *RedisStorage) LTrim(key string, start, stop int64) (string, error) {

	return rc.db.LTrim(context.Background(), key, start, stop).Result()
}

// ZRemRangeByRank 移除有序集合中给定的排名区间的所有成员
func (rc *RedisStorage) ZRemRangeByRank(key string, start, stop int64) (res int64, err error) {
	res, err = rc.db.ZRemRangeByRank(context.Background(), key, start, stop).Result()

	if err != nil {
		return 0, err
	}

	return res, nil
}

// Expire 设置过期时间
func (rc *RedisStorage) Expire(key string, expiration time.Duration) (res bool, err error) {
	res, err = rc.db.Expire(context.Background(), key, expiration).Result()

	if err != nil {
		return false, err
	}

	return res, err
}

// ZRem 从zset中移除变量
func (rc *RedisStorage) ZRem(key string, members ...interface{}) (res int64, err error) {
	res, err = rc.db.ZRem(context.Background(), key, members...).Result()

	if err != nil {
		return 0, err
	}
	return res, nil
}

// SAdd 向set中添加成员
func (rc *RedisStorage) SAdd(key string, member ...interface{}) (int64, error) {
	var intCmd *redis.IntCmd
	intCmd = rc.db.SAdd(context.Background(), key, member...)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// SMembers 返回set的全部成员
func (rc *RedisStorage) SMembers(key string) ([]string, error) {
	var stringSliceCmd *redis.StringSliceCmd
	stringSliceCmd = rc.db.SMembers(context.Background(), key)

	if err := stringSliceCmd.Err(); err != nil {
		return []string{}, err
	}
	return stringSliceCmd.Val(), nil
}

// SIsMember ...
func (rc *RedisStorage) SIsMember(key string, member interface{}) (bool, error) {
	var boolCmd *redis.BoolCmd
	boolCmd = rc.db.SIsMember(context.Background(), key, member)

	if err := boolCmd.Err(); err != nil {
		return false, err
	}
	return boolCmd.Val(), nil
}

func (rc *RedisStorage) SCard(key string) (int64, error) {
	var intCmd *redis.IntCmd
	intCmd = rc.db.SCard(context.Background(), key)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// HKeys 获取hash的所有域
func (rc *RedisStorage) HKeys(key string) []string {
	var stringSliceCmd *redis.StringSliceCmd
	stringSliceCmd = rc.db.HKeys(context.Background(), key)

	if err := stringSliceCmd.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	return stringSliceCmd.Val()
}

// HLen 获取hash的长度
func (rc *RedisStorage) HLen(key string) int64 {
	var intCmd *redis.IntCmd
	intCmd = rc.db.HLen(context.Background(), key)

	if err := intCmd.Err(); err != nil && err != redis.Nil {
		return 0
	}
	return intCmd.Val()
}

// GeoAdd 写入地理位置
func (rc *RedisStorage) GeoAdd(key string, location *redis.GeoLocation) (res int64, err error) {
	res, err = rc.db.GeoAdd(context.Background(), key, location).Result()

	if err != nil {
		return 0, err
	}

	return res, nil
}

// GeoRadius 根据经纬度查询列表
func (rc *RedisStorage) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) (res []redis.GeoLocation, err error) {
	res, err = rc.db.GeoRadius(context.Background(), key, longitude, latitude, query).Result()

	if err != nil {
		return []redis.GeoLocation{}, err
	}

	return res, nil
}
