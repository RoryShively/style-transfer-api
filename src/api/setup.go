package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/go-redis/redis"
)

var (
	dbmgr *DBManager
	dbmgrOnce  sync.Once
	pubsub *RedisManager
	pubsubOnce sync.Once
)

func getDBmgr() *DBManager {
	dbmgrOnce.Do(func() {
		dbmgr = NewDBManager()
	})
	return dbmgr
}

func getPubsub() *RedisManager {
	pubsubOnce.Do(func() {
		pubsub = NewRedisManager()
	})
	return pubsub
}


func setupDataDirs() {
	dataDirs := []string{
		"/data/uploaded",
		"/data/stylized",
	}

	for _, d := range dataDirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			os.Mkdir(d, 0777)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

// DBManagerCfg
type DBManagerCfg struct {
	Host string
	Port string
	User string
	DBName string
	Password string
}

func (cfg DBManagerCfg) ConnectionString() string {
	return fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v  sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)
}

type DBManager struct {
	DB *gorm.DB
}

func NewDBManager() *DBManager {
	cfg := DBManagerCfg{
		Host: os.Getenv("POSTGRES_HOST"),
		Port: os.Getenv("POSTGRES_PORT"),
		User: os.Getenv("POSTGRES_USER"),
		DBName: os.Getenv("POSTGRES_DB"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	var (
		connected bool
		db *gorm.DB
	)
	for !connected {
		var err error
		db, err = gorm.Open("postgres", cfg.ConnectionString())
		if err == nil {
			connected = true
		} else if err != nil && strings.Contains(err.Error(), "connection refused") {
			fmt.Println("Cant connect to DB. Waiting 5 seconds")
			time.Sleep(5 * time.Second)
		} else if err != nil {
			log.Fatal(err)
		}
	}

	return &DBManager{
		DB: db,
	}
}

func (mgr *DBManager) Migrate() {
	dbmgr.DB.AutoMigrate(&Image{})
}

func (mgr *DBManager) Close() {
	if mgr != nil {
		mgr.DB.Close()
	}
}



// RedisManagerCfg
type RedisManagerCfg struct {
	Host string
	Port string
}

func (cfg RedisManagerCfg) ConnectionString() string {
	return fmt.Sprintf("%v:%v",
		cfg.Host,
		cfg.Port,
	)
}

// RedisManager
type RedisManager struct {
	Client *redis.Client
}

func NewRedisManager() *RedisManager {
	cfg := RedisManagerCfg{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
	}

	var (
		connected bool
		rc *redis.Client
	)
	for !connected {
		rc = redis.NewClient(&redis.Options{
			Addr: cfg.ConnectionString(),
			DB: 0,
		})

		pong, err := rc.Ping().Result()
		if err != nil {
			log.Fatal(err)
		}
		if pong == "PONG" {
			connected = true
		}
	}

	return &RedisManager{
		Client: rc,
	}
}

func (mgr *RedisManager) Close() {
	if mgr.Client != nil {
		mgr.Client.Close()
	}
}