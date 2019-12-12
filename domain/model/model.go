package model

import (
	"encoding/json"
	"github.com/becosuke/guestbook-api/config"
	"github.com/go-redis/redis/v7"
)

const documentKey = "document"
const numPerPage = 10

var client = redis.NewClient(&redis.Options{
	Addr:     config.GetConfig().Redis.Addr,
	Password: "", // no password set
	DB:       0,  // use default DB
})

type Document struct {
	Name string `json:"name"`
	Body string `json:"body"`
	Score int64 `json:"score"`
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

func Range(page int) ([]byte, error) {
	start := (page - 1) * numPerPage
	end := page * numPerPage - 1
	documents := make([]Document, 0, numPerPage)
	for _, record := range client.ZRange(documentKey, int64(start), int64(end)).Val() {
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
