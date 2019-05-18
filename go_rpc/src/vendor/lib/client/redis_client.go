package client

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"lib/utils"

	"github.com/garyburd/redigo/redis"

	"lib/config"
)

type RedisClient struct {
	config          *config.ClientInfo
	redisClientPool *redis.Pool
}

func NewRedisClient(cinfo *config.ClientInfo) *RedisClient {
	rc := &RedisClient{}
	rc.config = cinfo
	rc.redisClientPool = &redis.Pool{
		MaxIdle:     cinfo.MaxIdleConn,
		MaxActive:   cinfo.MaxActiveConn,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", cinfo.Addrs[0], time.Duration(cinfo.Timeout)*time.Millisecond,
				time.Duration(cinfo.Timeout)*time.Millisecond, time.Duration(cinfo.Timeout)*time.Millisecond)
			if err != nil {
				return nil, err
			}

			if cinfo.Password != "" {
				if _, err := c.Do("AUTH", cinfo.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if cinfo.DbName != "" {
				if _, err := c.Do("SELECT", cinfo.DbName); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
	}
	return rc
}

func (rc *RedisClient) Exist(key string) (interface{}, error) {

	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":Exists")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		var data interface{}
		data, err = conn.Do("EXISTS", key)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}

	return nil, err
}

func (rc *RedisClient) Get(key string) (interface{}, error) {

	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":Get")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		var data interface{}
		data, err = conn.Do("GET", key)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}

	return nil, err
}

func (rc *RedisClient) Set(key string, value interface{}) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":Set")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("SET", key, value)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return nil
		}
	}

	return err
}

func (rc *RedisClient) SetWithttl(key string, value interface{}, ttl int) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":Set")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("SET", key, value, "EX", ttl)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return nil
		}
	}

	return err
}

func (rc *RedisClient) ZSet(param []interface{}) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":ZSet")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("ZADD", param...)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return nil
		}
	}

	return err

}

func (rc *RedisClient) ZRem(param []interface{}) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":ZRem")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("ZREM",param...)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return nil
		}
	}
	return err
}

func (rc *RedisClient) ZGet(key string, start int, end int) (interface{}, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":ZGet")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		var data interface{}
		data, err = conn.Do("ZREVRANGE", key, start, end, "WITHSCORES")
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}
	return nil, err
}

func (rc *RedisClient) ZGetWithoutScore(key string, start int, end int) (interface{}, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":ZGetWithoutScore")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		var data interface{}
		data, err = conn.Do("ZREVRANGE", key, start, end)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}
	return nil, err
}

func (rc *RedisClient) ZRange(key string, start int, end int) (interface{}, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":ZRange")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		var data interface{}
		data, err = conn.Do("ZRANGE", key, start, end, "WITHSCORES")
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}

	return nil, err
}

//set结构
func (rc *RedisClient) SAdd(key string, param interface{}) (int, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":SAdd")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		// doParam := append([]interface{}{key}, param...)
		var data int
		data, err = redis.Int(conn.Do("SADD", redis.Args{}.Add(key).AddFlat(param)...))
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}
	return 0, err
}
func (rc *RedisClient) SMembers(key string) (interface{}, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":SMembers")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		var data interface{}
		data, err = conn.Do("SMEMBERS", key)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}
	return nil, err
}

func (rc *RedisClient) Scan(cursor int, match string, count int) (next int, keys []string, err error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":Scan")
	countScope.SetErr()
	defer countScope.End()

	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		var values []interface{}
		scanCommand := redis.Args{}.Add(cursor)
		if "" != match {
			scanCommand.Add("MATCH", match)
		}
		scanCommand = scanCommand.Add("COUNT", count)
		values, err = redis.Values(conn.Do("SCAN", scanCommand...))
		conn.Close()
		if err == nil {
			_, err = redis.Scan(values, &next, &keys)
			if err == nil {
				countScope.SetOk()
				return
			}
		}
	}
	return
}

func (rc *RedisClient) IScan(node, cursor int, match string, count int) (next int, keys []string, err error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":IScan")
	countScope.SetErr()
	defer countScope.End()

	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		var values []interface{}
		scanCommand := redis.Args{}.Add(node, cursor)
		if "" != match {
			scanCommand.Add("MATCH", match)
		}
		scanCommand = scanCommand.Add("COUNT", count)
		values, err = redis.Values(conn.Do("ISCAN", scanCommand...))
		conn.Close()
		if err == nil {
			_, err = redis.Scan(values, &next, &keys)
			if err == nil {
				countScope.SetOk()
				return
			}
		}
	}
	return
}

func (rc *RedisClient) Keys(pattern string) (interface{}, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":Keys")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		var data interface{}
		data, err = conn.Do("KEYS", pattern)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}
	return nil, err
}

func (rc *RedisClient) Expire(key string, livetime int) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":Expire")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("expire", key, livetime)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return nil
		}
	}

	return err
}

func (rc *RedisClient) MExpire(keys []interface{}, livetime int64) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":MExpire")
	countScope.SetErr()
	defer countScope.End()
	var err error
	klen := len(keys)

	if klen == 0 {
		return ErrFunctionParamErr
	}

	conn := rc.redisClientPool.Get()

	defer conn.Close()

	for i := 0; i < klen; i++ {
		err = conn.Send("expire", keys[i], livetime)
		if err != nil {
			return err
		}
	}
	conn.Flush()

	result := make([]interface{}, klen)
	for i := 0; i < klen; i++ {
		result[i], err = conn.Receive()
		if err != nil {
			return err
		}
	}

	countScope.SetOk()
	return nil
}

func (rc *RedisClient) Delete(key string) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":Delete")
	countScope.SetErr()
	defer countScope.End()

	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("del", key)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return nil
		}
	}

	return err
}

func (rc *RedisClient) MGet(keys []interface{}) (interface{}, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":MGet")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		var data interface{}
		conn := rc.redisClientPool.Get()
		data, err = conn.Do("mget", keys...)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}

	return nil, err
}

func (rc *RedisClient) MIncFloat(keys []interface{}, values []float64) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":MIncFloat")
	countScope.SetErr()
	defer countScope.End()
	var err error
	klen := len(keys)

	if klen == 0 || klen != len(values) {
		return ErrFunctionParamErr
	}

	conn := rc.redisClientPool.Get()

	defer conn.Close()

	for i := 0; i < klen; i++ {
		err = conn.Send("INCRBYFLOAT", keys[i], values[i])
		if err != nil {
			return err
		}
	}
	conn.Flush()
	countScope.SetOk()
	return nil
}

func (rc *RedisClient) MSet(keys_values []interface{}) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":MSet")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("MSET", keys_values...)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return nil
		}
	}

	return err
}

// TupleSSList is a helper that converts an array of strings
// into a TupleSSList. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func ToTupleSSList(result interface{}, err error) ([]utils.TupleSS, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	m := make([]utils.TupleSS, len(values)/2)
	j := 0
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		value, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: ScanMap key not a bulk string value")
		}
		m[j].First = string(key)
		m[j].Second = string(value)
		j++
	}
	return m, nil
}

// TupleSFList is a helper that converts an array of strings
// into a TupleSFList. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func ToTupleSFList(result interface{}, err error) ([]utils.TupleSF, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: ToTupleSFList expects even number of values result")
	}
	ret := make([]utils.TupleSF, len(values)/2)
	j := 0
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		value, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: ScanMap key not a bulk string value")
		}
		ret[j].First = string(key)
		ret[j].Second, _ = strconv.ParseFloat(string(value), 64)
		j++
	}
	return ret, nil
}

func (rc *RedisClient) BatchSetListFloatWithTTL(tupleList []utils.TupleSF, ttl int) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":BatchSetFloatWithTTL")
	countScope.SetErr()
	defer countScope.End()

	klen := len(tupleList)
	if klen == 0 || klen != len(tupleList) {
		return ErrFunctionParamErr
	}
	conn := rc.redisClientPool.Get()
	defer conn.Close()
	for i := 0; i < klen; i++ {
		err := conn.Send("SETEX", tupleList[i].First, ttl, strconv.FormatFloat(tupleList[i].Second, 'f', -1, 64))
		if err != nil {
			return err
		}
	}

	err := conn.Flush()
	if err == nil {
		countScope.SetOk()
	}

	return err
}

func (rc *RedisClient) BatchInCrByFloat(keys []string, values []float64) ([]interface{}, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":BatchInCrByFloat")
	countScope.SetErr()
	defer countScope.End()
	var err error

	klen := len(keys)

	if klen == 0 || klen != len(values) {
		return nil, ErrFunctionParamErr
	}
	/*
		trynum := 1
			if rc.config.TryNum > 0 {
				trynum = rc.config.TryNum
			} not try */
	//var data interface{}
	conn := rc.redisClientPool.Get()

	defer conn.Close()

	for i := 0; i < klen; i++ {
		err = conn.Send("INCRBYFLOAT", keys[i], strconv.FormatFloat(values[i], 'f', -1, 64))
		if err != nil {
			return nil, err
		}
	}

	conn.Flush()

	result := make([]interface{}, klen)
	for i := 0; i < klen; i++ {
		result[i], err = conn.Receive()
		if err != nil {
			return nil, err
		}
	}

	countScope.SetOk()
	return result, nil

}

//HGetKVList  以列表的形式存储key-value
func (rc *RedisClient) HGetKVList(key interface{}) ([]utils.TupleSS, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":HGetKVList")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		data, err := ToTupleSSList(conn.Do("HGETALL", key))
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}

	return nil, err
}

func (rc *RedisClient) HGetAll(key interface{}) (map[string]string, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":HGetAll")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		data, err := redis.StringMap(conn.Do("HGETALL", key))
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}

	return nil, err
}

func (rc *RedisClient) HMGet(key_values []interface{}) ([]string, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":HMGet")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		data, err := Strings(conn.Do("HMGet", key_values...))
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}

	return nil, err
}

// Strings is a helper that converts an array command reply to a []string. If
// err is not equal to nil, then Strings returns nil, err. Nil array items are
// converted to "" in the output slice. Strings returns an error if an array
// item is not a bulk string or nil.
func Strings(reply interface{}, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []interface{}:
		result := make([]string, len(reply))
		for i := range reply {
			if reply[i] == nil {
				result[i] = ""
				continue
			}
			p, ok := reply[i].([]byte)
			if !ok {
				return nil, fmt.Errorf("redigo: unexpected element type for Strings, got type %T", reply[i])
			}
			result[i] = string(p)
		}
		return result, nil
	case nil:
		return nil, redis.ErrNil
	case redis.Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: unexpected type for Strings, got type %T", reply)
}

func (rc *RedisClient) RPush(key_values []interface{}) (int, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":RPush")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("RPUSH", key_values...)
		conn.Close()
		if err == nil {
			//if ret, ok := data.(int); ok {
			countScope.SetOk()
			return 0, nil
			//} else {
			//	err = ErrRedisDataFormat
			//}
		}
	}

	return 0, err
}

func (rc *RedisClient) LPush(key_values []interface{}) (int, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":LPush")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("LPUSH", key_values...)
		conn.Close()
		if err == nil {
			//if ret, ok := data.(int); ok {
			countScope.SetOk()
			return 0, nil
			//} else {
			//	err = ErrRedisDataFormat
			//}
		}
	}

	return 0, err
}
func (rc *RedisClient) LRange(key_values []interface{}) ([]string, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":LRange")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		data, err := redis.Strings(conn.Do("LRANGE", key_values...))
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}

	return nil, err
}

func (rc *RedisClient) Cmd(cmd string, key_values []interface{}) (interface{}, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":" + cmd)
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		var data interface{}
		conn := rc.redisClientPool.Get()
		data, err = conn.Do(cmd, key_values...)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return data, nil
		}
	}

	return nil, err
}

func (rc *RedisClient) HMSet(key_values []interface{}) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":HMSet")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("HMSet", key_values...)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return nil
		}
	}

	return err
}

func (rc *RedisClient) GetConf() *config.ClientInfo {
	return rc.config
}

func (rc *RedisClient) GetConn() redis.Conn {
	return rc.redisClientPool.Get()
}

func (rc * RedisClient) BatchHGetAll(keys []string) ([]map[string]string, error){
	values := make([]map[string]string,len(keys),len(keys))
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":BatchHGetAll")
	countScope.SetErr()
	defer countScope.End()
	conn := rc.redisClientPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	for _, key := range keys{
		err := conn.Send("HGETALL", key)
		if err != nil {
			return nil, err
		}
	}
	resp, err := conn.Do("EXEC")
	if err != nil{
		return nil, err
	}
	reply := resp.([]interface{})
	for i:=0;i<len(keys);i++{
		s, _ := redis.StringMap(reply[i], nil)
		values[i] = s
	}
	countScope.SetOk()
	return values, nil
}

func (rc * RedisClient) BatchHMSet(values [][]interface{}) (int, error){
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":BatchHMSet")
	countScope.SetErr()
	defer countScope.End()
	conn := rc.redisClientPool.Get()
	defer conn.Close()
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	success := 0
	conn.Send("MULTI")
	for _, value := range values{
		for t := 0; t < trynum; t++ {
			err := conn.Send("HMSet", value...)
			if err == nil {
				success+=1
				countScope.SetOk()
				break
			}
		}
	}
	_, err := conn.Do("EXEC")
	if err != nil{
		return 0, err
	}
	return success, nil
}

func (rc *RedisClient) ExpireAt(key string, livetime int64) error {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":ExpireAt")
	countScope.SetErr()
	defer countScope.End()
	var err error
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	for t := 0; t < trynum; t++ {
		conn := rc.redisClientPool.Get()
		_, err = conn.Do("EXPIREAT", key, livetime)
		conn.Close()
		if err == nil {
			countScope.SetOk()
			return nil
		}
	}
	return err
}

func (rc *RedisClient) MExpireAt(keys []interface{}, expire_at int64) (int,error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":MExpireAt")
	countScope.SetErr()
	defer countScope.End()
	klen := len(keys)

	if klen == 0 {
		return 0, ErrFunctionParamErr
	}

	conn := rc.redisClientPool.Get()
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	defer conn.Close()
	conn.Send("MULTI")
	success := 0
	for i := 0; i < klen; i++ {
		for t := 0; t < trynum; t++ {
			err := conn.Send("EXPIREAT", keys[i], expire_at)
			if err == nil {
				success+=1
				break
			}
		}
	}
	_, err := conn.Do("EXEC")
	if err != nil{
		return 0, err
	}
	countScope.SetOk()
	return success, nil
}


func (rc *RedisClient) BatchZGet(keys []interface{}, start int, end int) ([]interface{}, error) {
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":BatchZGet")
	countScope.SetErr()
	defer countScope.End()
	klen := len(keys)
	if klen == 0 {
		return nil, ErrFunctionParamErr
	}

	conn := rc.redisClientPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	for i := 0; i < klen; i++ {
		err := conn.Send("ZREVRANGE", keys[i], start, end, "WITHSCORES")
		if err != nil {
			return nil, err
		}
	}
	resp, err := conn.Do("EXEC")
	if err != nil{
		return nil, err
	}
	data := resp.([]interface{})
	countScope.SetOk()
	return data, nil
}

func (rc * RedisClient) BatchZSet(values [][]interface{}) (int, error){
	countScope := countstat.GetCountGlobal().NewCountScope("Redisclient:" + rc.config.Name + ":BatchZSet")
	countScope.SetErr()
	defer countScope.End()
	conn := rc.redisClientPool.Get()
	defer conn.Close()
	trynum := 1
	if rc.config.TryNum > 0 {
		trynum = rc.config.TryNum
	}
	success := 0
	conn.Send("MULTI")
	for _, value := range values{
		for t := 0; t < trynum; t++ {
			err := conn.Send("ZADD", value...)
			if err == nil {
				success+=1
				break
			}
		}
	}
	_, err := conn.Do("EXEC")
	if err != nil{
		return 0, err
	}
	countScope.SetOk()
	return success, nil
}
