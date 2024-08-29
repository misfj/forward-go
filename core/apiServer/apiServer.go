package apiServer

import (
	"context"
	"errors"
	"fmt"
	"forward-go/global"
	"forward-go/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiServer struct {
	eng       *gin.Engine
	httpServe *http.Server
	ctx       context.Context
}

func NewApiServer() *ApiServer {
	if global.GlobalConfig.Forward.AppDebug {
		//GIN开启Debug
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	return &ApiServer{
		eng: gin.New(),
	}
}
func (s *ApiServer) Init() error {
	//初始化路由
	return nil
}

func (s *ApiServer) Name() string {
	return "APIServer"
}

func (s *ApiServer) Startup(ctx context.Context) error {
	if !global.GlobalConfig.Forward.EnableTLS {
		s.httpServe = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", global.GlobalConfig.Forward.Host, global.GlobalConfig.Forward.Port),
			Handler: s.eng,
		}
	} else {
		//开启https,读取证书
	}
	if err := s.httpServe.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
	return nil
}

func (s *ApiServer) Close() error {
	log.Info("APIServer close..")
	return s.httpServe.Shutdown(s.ctx)
}
