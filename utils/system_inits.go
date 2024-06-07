package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var ERROR error
var Redis *redis.Client

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited")
	fmt.Println("config mysql inited")
}
func InitMySql() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	DB, ERROR = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	if ERROR != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("config mysql inited")
	}
}

func InitRedis() {
	ctx := context.Background()
	Redis = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdeConns"),
	})
	pong, err := Redis.Ping(ctx).Result()
	if err != nil {
		fmt.Println("failed to connect redis", err)
	} else {
		fmt.Println("redis connected...", pong)
	}
}

const (
	PublishKey = "websocket"
)

// publish 发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	err := Redis.Publish(ctx, channel, msg).Err()
	return err
}

func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Redis.Subscribe(ctx, channel)
	fmt.Println("subscribed...")
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("subscribed...", msg.Payload)
	return msg.Payload, err
}
