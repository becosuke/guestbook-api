package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	nowFixation bool
	nowDatetime time.Time
	SystemEnv   string
	Redis       Redis
	Rest        Rest
}

type Redis struct {
	Host string
	Port int
	Addr string
}

type Rest struct {
	Host string
	Port int
	Addr string
}

var conf *Config

func init() {
	location := "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	var now time.Time
	var nowFixation bool
	if tmp := os.Getenv("NOW_DATETIME"); tmp == "" {
		now = time.Now()
		nowFixation = false
	} else {
		if parsed, err := time.Parse("2006-01-02 15:04:05", tmp); err != nil {
			now = time.Now()
			nowFixation = false
		} else {
			now = parsed
			nowFixation = true
		}
	}

	var systemEnv string
	if systemEnv = os.Getenv("SYSTEM_ENV"); systemEnv == "" {
		systemEnv = "local"
	}

	var redisHost string
	if redisHost = os.Getenv("SYSTEM_REDIS_HOST"); redisHost == "" {
		redisHost = "localhost"
	}

	var redisPort int
	if redisPort, err = strconv.Atoi(os.Getenv("SYSTEM_REDIS_PORT")); err != nil {
		redisPort = 6379
	}

	var restHost string
	if restHost = os.Getenv("SYSTEM_REST_HOST"); restHost == "" {
		restHost = "0.0.0.0"
	}

	var restPort int
	if restPort, err = strconv.Atoi(os.Getenv("SYSTEM_REST_PORT")); err != nil {
		restPort = 80
	}

	conf = &Config{
		nowFixation: nowFixation,
		nowDatetime: now,
		SystemEnv:   systemEnv,
		Redis: Redis{
			Host: redisHost,
			Port: redisPort,
			Addr: redisHost + ":" + strconv.Itoa(redisPort),
		},
		Rest: Rest{
			Host: restHost,
			Port: restPort,
			Addr: restHost + ":" + strconv.Itoa(restPort),
		},
	}
}

func GetConfig() *Config {
	return conf
}

func NowDatetime() time.Time {
	if GetConfig().nowFixation {
		return GetConfig().nowDatetime
	}

	return time.Now()
}

func IsLocal() bool {
	return GetConfig().SystemEnv == "local"
}
