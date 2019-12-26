package model

import (
	"encoding/json"
	"github.com/becosuke/guestbook-api/config"
	"github.com/go-redis/redis/v7"
	"math/rand"
)

func init() {
	rand.Seed(config.NowDatetime().UnixNano())
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const documentKey = "document"
const NumPerPage = 10

var client = redis.NewClient(&redis.Options{
	Addr:     config.GetConfig().Redis.Addr,
	Password: "", // no password set
	DB:       0,  // use default DB
})

type Document struct {
	Name  string `json:"name"`
	Body  string `json:"body"`
	Score int64  `json:"score"`
}

func Ping() error {
	_, err := client.Ping().Result()
	return err
}

func Add(body []byte) ([]byte, error) {
	document := Document{}
	if err := json.Unmarshal(body, &document); err != nil {
		return []byte{}, err
	}
	document.Score = config.NowDatetime().UnixNano()

	record, err := json.Marshal(document)
	if err != nil {
		return []byte{}, err
	}
	client.ZAdd(documentKey, &redis.Z{Score: float64(document.Score), Member: record})
	return record, nil
}

func Range(start int, end int) ([]byte, error) {
	documents := make([]Document, 0, NumPerPage)
	for _, record := range client.ZRevRange(documentKey, int64(start), int64(end)).Val() {
		document := Document{}
		if err := json.Unmarshal([]byte(record), &document); err != nil {
			return []byte{}, err
		}
		documents = append(documents, document)
	}
	res, err := json.Marshal(documents)
	if err != nil {
		return []byte{}, err
	}
	return res, nil
}

func Flush() {
	client.FlushAll()
}

func Random() {
	for i := 0; i < 1000; i++ {
		document := Document{Name: randString(10), Body: randString(50), Score: config.NowDatetime().UnixNano()}
		record, err := json.Marshal(document)
		if err != nil {
			continue
		}
		client.ZAdd(documentKey, &redis.Z{Score: float64(document.Score), Member: record})
	}
}

func randString(n int) string {
	bs := make([]byte, n)
	for i := range bs {
		bs[i] = letters[rand.Intn(len(letters))]
	}
	return string(bs)
}
