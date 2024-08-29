package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var Conf Config

type Db struct {
	DbType string `json:"dbType"`
	Detail struct {
		Host        string `json:"host"`
		Port        int    `json:"port"`
		User        string `json:"user"`
		Password    string `json:"password"`
		Schema      string `json:"schema"`
		MaxPoolSize int    `json:"maxPoolSize"`
		MaxIdleSize int    `json:"maxIdleSize"`
	} `json:"detail"`
}
type Logging struct {
	File       string `json:"file"`
	Level      string `json:"level"`
	MaxAge     int    `json:"maxAge"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	Compress   bool   `json:"compress"`
}
type Store struct {
	Cache struct {
		Type   string `json:"type"`
		Detail struct {
			Expire          int    `json:"expire"`
			UseSSL          bool   `json:"useSSL"`
			AccessID        string `json:"accessID"`
			EndPoint        string `json:"endPoint"`
			SecretAccessKey string `json:"secretAccessKey"`
		} `json:"detail"`
	} `json:"cache"`
	Lasting struct {
		Type   string `json:"type"`
		Detail struct {
			Expire          int    `json:"expire"`
			UseSSL          bool   `json:"useSSL"`
			AccessID        string `json:"accessID"`
			EndPoint        string `json:"endPoint"`
			SecretAccessKey string `json:"secretAccessKey"`
		} `json:"detail"`
	} `json:"lasting"`
}
type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	IndexDb  int    `json:"indexDb"`
	User     string `json:"user"`
	Password string `json:"password"`
}
type Server struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	MsgEncodeKey       string `json:"msgEncodeKey"`
	MaxRecount         int    `json:"maxRecount"`
	HealthyInterval    int    `json:"healthyInterval"`
	AlgoManageProtocol string `json:"AlgoManageProtocol"`
	AlgoManageURL      string `json:"AlgoManageUrl"`
	AppDebug           bool   `json:"appDebug"`
	CasbinModel        string `json:"casbinModel"`
	EnableTLS          bool   `json:"enableTLS"`
	Cert               string `json:"cert"`
	KeyFile            string `json:"keyFile"`
	JwtSignKey         string `json:"jwtSignKey"`
	JwtOnlineUsers     int    `json:"jwtOnlineUsers"`
	JwtExpireAt        int    `json:"jwtExpireAt"`
	IsRedis            bool   `json:"isRedis"`
}

type Config struct {
	Db      Db      `json:"db"`
	Logging Logging `json:"logging"`
	Redis   Redis   `json:"redis"`
	Store   Store   `json:"store"`
	Forward Server  `json:"forward"`
}

func LoadConfig() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = os.Stat(filepath.Join(dir, "config.json"))
	if err != nil {
		//生成默认配置及文件
		configFile, err := os.Create(filepath.Join(dir, "config.json"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_, err = io.Copy(configFile, strings.NewReader(defaultConfigJsonV1))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("已经生成默认配置文件,请填写配置信息后,重新启动")
		os.Exit(1)
	}
	//读取配置文件
	configContent, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = json.Unmarshal(configContent, &Conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(Conf)
}

const defaultConfigJsonV1 = `{
"db": {
  "dbType": "sqlite",
  "detail": {
    "host": "localhost",
    "port": 1234,
    "user": "root",
    "password": "password",
    "schema": "",
    "maxPoolSize": 10,
    "maxIdleSize": 5
  }
},
  "logging": {
    "file": "logs/forward.log",
    "level": "debug",
    "maxAge": 30,
    "maxSize": 10,
    "maxBackups": 5
  },
  "store": {
    "cache": {
      "type": "minio",
       "detail": {
       "expire": 3,
         "useSSL": false,
         "accessID": "minio_admin",
         "endPoint": "10.10.1.40:9000",
         "secretAccessKey": "9ijnBHU*@123"
       }

    },
    "lasting": {
      "type": "ftp",
      "detail": {

      }
    }
  },
  "forward": {
    "msgEncodeKey":"",
    "maxRecount": 5,
    "healthyInterval": 15,
    "AlgoManageProtocol":"websocket",
    "AlgoManageUrl": "http://localhost"
  }
}`
