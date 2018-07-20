package util

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"net/http"
	"strconv"
	"time"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// pong, err := redisClient.Ping().Result()
	// fmt.Println(pong, err)
}

type Movie struct {
	ID              int     `json:"id"`
	HaspromotionTag bool    `json:"haspromotionTag"`
	Img             string  `json:"img"`
	Version         string  `json:"version"`
	Nm              string  `json:"nm"`
	PreShow         bool    `json:"preShow"`
	Sc              float64 `json:"sc"`
	GlobalReleased  bool    `json:"globalReleased"`
	Wish            int     `json:"wish"`
	Star            string  `json:"star"`
	Rt              string  `json:"rt"`
	ShowInfo        string  `json:"showInfo"`
	Showst          int     `json:"showst"`
	Wishst          int     `json:"wishst"`
	NmEn            string  `json:"nmEn"`
}

type AutoGeneratedMovie struct {
	Coming   []interface{} `json:"coming"`
	MovieIds []int         `json:"movieIds"`
	Stid     string        `json:"stid"`
	Stids    []struct {
		MovieID int    `json:"movieId"`
		Stid    string `json:"stid"`
	} `json:"stids"`
	Total     int     `json:"total"`
	MovieList []Movie `json:"movieList"`
}

func FetchMaoyanApi() error {
	resp, err := http.Get("http://m.maoyan.com/ajax/movieOnInfoList")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	obj := AutoGeneratedMovie{}
	json.NewDecoder(resp.Body).Decode(&obj)

	_, err = redisClient.Pipelined(func(pipe redis.Pipeliner) error {
		skey := time.Now().Format("movie-2006-01-02")
		_, eerr := pipe.Expire(skey, time.Hour*24).Result()
		key := "maoyan_movie"
		for _, vv := range obj.MovieList {
			field := strconv.Itoa(vv.ID)
			pipe.SAdd(skey, vv.ID)
			if redisClient.HGet(key, field).Val() == "" {
				vv.NmEn = TranslateCh2En(vv.Nm)
				if b, eerr := json.Marshal(vv); eerr == nil {
					pipe.HSet(key, field, b)
				}
				time.Sleep(time.Microsecond * 5000)
			}
		}
		return eerr
	})

	return err
}

func FetchMusic163Api() error {
	resp, err := http.Get("http://m.maoyan.com/ajax/movieOnInfoList")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	obj := AutoGeneratedMovie{}
	json.NewDecoder(resp.Body).Decode(&obj)

	_, err = redisClient.Pipelined(func(pipe redis.Pipeliner) error {
		skey := time.Now().Format("2006-01-02")
		_, eerr := pipe.Expire(skey, time.Hour*24).Result()
		key := "maoyan_movie"
		for _, vv := range obj.MovieList {
			field := strconv.Itoa(vv.ID)
			pipe.SAdd(skey, vv.ID)
			if redisClient.HGet(key, field).Val() == "" {
				vv.NmEn = TranslateCh2En(vv.Nm)
				if b, eerr := json.Marshal(vv); eerr == nil {
					pipe.HSet(key, field, b)
				}
				time.Sleep(time.Microsecond * 5000)
			}
		}
		return eerr
	})

	return err
}
