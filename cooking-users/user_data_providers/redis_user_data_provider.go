package user_data_providers

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"log"
	"strconv"
)


const DbKey = "GAME:"

type UserProvider struct {
	Client *redis.Client
}

func (s *UserProvider) GetMany(userId string, names []string) map[string]string {
	result:=make(map[string]string)
	data, _ := s.Client.HGetAll(DbKey + userId).Result()
	for _,k := range names {
		result[k]=data[k]
	}
	return result
}

func (s *UserProvider) DecrementOne(userId string, resourceName string, value interface{}) (string, error) {
	var val int
	switch value.(type) {
	case int:
		val = value.(int)
		break
	case string:
		val, _ = strconv.Atoi(value.(string))
	}

	current, _:= strconv.Atoi(s.GetOne(userId, resourceName))


	if current < val {
		return string(current), errors.New("not enough resources")
	} else {
		newVal:=string(current-val)
		e:= s.SetOne(userId, resourceName, newVal)
		return newVal, e
	}


}

func (s *UserProvider) SetOne(userId string, key string, value interface{})   error    {
		_, err:= s.Client.HSet(DbKey+userId, key, value).Result()
		return err
}

func NewRedisUserProvider(client *redis.Client) UserDataProvider {
	provider:=UserProvider{Client:client}
	return &provider
}


func (s *UserProvider) DeleteAll(userId string) error {
	return s.Client.Del(DbKey + userId).Err()
}

func (s *UserProvider) GetOne(userId string, name string) string {
	result, _ := s.Client.HGet(DbKey + userId, name).Result()
	return result
}

func (s *UserProvider) IncrementOne(userId string, k string, v interface{}) ( string, error ) {
	result := s.Client.HIncrByFloat(DbKey+userId, k, v.(float64))
	return fmt.Sprintf("{%f}",result.Val()), result.Err()
}

func (s *UserProvider) IncrementMany(userId string, data map[string]interface{}) error {
	if len(data) != 0 {
		for k, v := range data {
			if _, err := s.Client.HIncrByFloat(DbKey+userId, k, v.(float64)).Result(); err!=nil {
				log.Fatalf("Resourcese update error: {%s} ", err.Error())
				return err
			}
		}
	}
	return nil
}

func (s *UserProvider) SetMany(userId string, data map[string]interface{}) error {
	if len(data) != 0 {
		if  _, err := s.Client.HMSet(DbKey+userId, data).Result(); err!=nil { return err }
	}
	return nil
}

func (s *UserProvider) GetAll(userId string) map[string]string {
	result, _ := s.Client.HGetAll(DbKey + userId).Result()
	return result
}


func (s *UserProvider) GetOrCreate(userId string, params map[string]interface{}) (User, error) {
	user := User {
		UserId: userId,
	}

	if values, err := s.Client.HGetAll(DbKey + userId).Result(); err == nil {
		if len(values) == 0 {
			user.SetDefault()
			s.Client.HMSet(DbKey+userId, user.Params)
		} else {
			user.Fill(values)
		}

		if params != nil && len(params) != 0 {
			user.Update(params)
			s.Client.HMSet(DbKey+userId, params)
		}
	} else {
		log.Printf("HGetAllError: {%s} ", err)
	}
	return user, nil
}




