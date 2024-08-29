package global

import (
	"fmt"
	"forward-go/config"
	"forward-go/log"
	"forward-go/utils"
	"github.com/glebarez/sqlite"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	stdlogger "log"
	"os"
	"strings"
	"time"
)

var (
	WebsocketMessageID int64
	GlobalConfig       *config.Config
	GlobalDB           *gorm.DB
	GlobalRedisClient  *redis.Client
	//GlobalLogger      *log.Logger
	GlobalCache       *minio.Client
	GlobalLasting     interface{}
	GlobalIDGenerator *utils.IDGenerator
	GlobalLocalCache  *utils.GCache
)

func LoadGlobal() {
	WebsocketMessageID = 0
	GlobalIDGenerator = utils.NewIDGenerator(WebsocketMessageID)
	GlobalConfig = &config.Conf
	//数据库
	if strings.EqualFold(GlobalConfig.Db.DbType, "sqlite") {
		openSqlite()
	} else {
		//支持mysql请添加Mysql链接
		fmt.Println("不支持其他数据库版本")
		os.Exit(1)
	}
	//存储
	if strings.EqualFold(GlobalConfig.Store.Cache.Type, "minio") {
		openMinio()
	}
	openLogger()
	//连接Redis,
	if GlobalConfig.Redis.Host == "" && GlobalConfig.Redis.Port == 0 {
		//説明沒有配置Redis實例,選取本地Redis,使用本地Cache
		GlobalLocalCache = utils.New(100)
	} else {
		GlobalRedisClient = redis.NewClient(&redis.Options{
			Username: GlobalConfig.Redis.User,
			Password: GlobalConfig.Redis.Password,
			Addr:     fmt.Sprintf("%s:%d", GlobalConfig.Redis.Host, GlobalConfig.Redis.Port),
		})
	}
}
func openSqlite() {
	var err error
	if GlobalDB, err = gorm.Open(sqlite.Open("forward.db"), &gorm.Config{
		PrepareStmt: true,
		Logger: logger.New(
			stdlogger.New(os.Stdout, "\r\n", stdlogger.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,          // Don't include params in the SQL log
				Colorful:                  true,          // Disable color
			},
		),
	}); err != nil {
		fmt.Println("连接数据库失败:", err)
		os.Exit(1)
	}
	//设置参数
	rawDB, _ := GlobalDB.DB()
	rawDB.SetMaxOpenConns(GlobalConfig.Db.Detail.MaxPoolSize)
	rawDB.SetMaxIdleConns(GlobalConfig.Db.Detail.MaxIdleSize)
}
func openMinio() {
	var err error
	GlobalCache, err = minio.New(GlobalConfig.Store.Cache.Detail.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(GlobalConfig.Store.Cache.Detail.AccessID, GlobalConfig.Store.Cache.Detail.SecretAccessKey, ""),
		Secure: GlobalConfig.Store.Cache.Detail.UseSSL,
	})
	if err != nil {
		fmt.Println("连接minio存储失败:", err)
		os.Exit(1)
	}
	//
}

func openLogger() {
	err := log.Init(&GlobalConfig.Logging)
	if err != nil {
		fmt.Println("初始化日志失败:", err)
		os.Exit(1)
	}
	//GlobalLogger = log.DefaultLogger
}
