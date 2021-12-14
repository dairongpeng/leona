// Copyright 2021 dairongpeng <dairongpeng@foxmail.com>. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiserver

import (
	"fmt"

	pb "github.com/dairongpeng/leona/api/proto/apiserver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/dairongpeng/leona/internal/apiserver/config"
	cachev1 "github.com/dairongpeng/leona/internal/apiserver/controller/v1/cache"
	"github.com/dairongpeng/leona/internal/apiserver/store"
	"github.com/dairongpeng/leona/internal/apiserver/store/mysql"
	genericoptions "github.com/dairongpeng/leona/internal/pkg/options"
	genericapiserver "github.com/dairongpeng/leona/internal/pkg/server"
	"github.com/dairongpeng/leona/pkg/log"
	"github.com/dairongpeng/leona/pkg/shutdown"
	"github.com/dairongpeng/leona/pkg/shutdown/shutdownmanagers/posixsignal"
)

type apiServer struct {
	gs *shutdown.GracefulShutdown
	// redisOptions     *genericoptions.RedisOptions
	gRPCAPIServer    *grpcAPIServer
	genericAPIServer *genericapiserver.GenericAPIServer
}

type preparedAPIServer struct {
	*apiServer
}

// ExtraConfig defines extra configuration for the leona-apiserver.
type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	ServerCert   genericoptions.GeneratableKeyCert
	mysqlOptions *genericoptions.MySQLOptions
	// etcdOptions      *genericoptions.EtcdOptions
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	// 创建优雅关停的实例
	gs := shutdown.New()
	// 对优雅关停的实例添加监听信号
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	// 构建通用的配置
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	// 构建额外的配置
	extraConfig, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}

	// 补全通用配置。再new实例（HTTP）
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}
	// 补全额外配置。再new实例（GRPC）
	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}

	// HTTP/GRPC服务的实例
	server := &apiServer{
		gs: gs,
		// redisOptions:     cfg.RedisOptions,
		genericAPIServer: genericServer,
		gRPCAPIServer:    extraServer,
	}

	return server, nil
}

func (s *apiServer) PrepareRun() preparedAPIServer {
	// 初始化路由配置
	initRouter(s.genericAPIServer.Engine)

	// 初始化redis数据库
	//s.initRedisStore()

	// 监听到信号后，执行回调，做一些收尾清理工作，优雅关停
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		mysqlStore, _ := mysql.GetMySQLFactoryOr(nil)
		if mysqlStore != nil {
			return mysqlStore.Close()
		}

		s.gRPCAPIServer.Close()
		s.genericAPIServer.Close()

		return nil
	}))

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	// 启动GRPC server
	go s.gRPCAPIServer.Run()

	// start shutdown managers
	// 启动优雅关停的实例
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// 启动HTTP server
	return s.genericAPIServer.Run()
}

type completedExtraConfig struct {
	*ExtraConfig
}

// Complete fills in any fields not set that are required to have valid data and can be derived from other fields.
func (c *ExtraConfig) complete() *completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{c}
}

// New create a grpcAPIServer instance.
func (c *completedExtraConfig) New() (*grpcAPIServer, error) {
	// 注释证书相关内容
	//creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	//if err != nil {
	//	log.Fatalf("Failed to generate credentials %s", err.Error())
	//}
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize)}
	grpcServer := grpc.NewServer(opts...)

	storeIns, _ := mysql.GetMySQLFactoryOr(c.mysqlOptions)
	// storeIns, _ := etcd.GetEtcdFactoryOr(c.etcdOptions, nil)
	store.SetClient(storeIns)
	cacheIns, err := cachev1.GetCacheInsOr(storeIns)
	if err != nil {
		log.Fatalf("Failed to get cache instance: %s", err.Error())
	}

	pb.RegisterCacheServer(grpcServer, cacheIns)

	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, c.Addr}, nil
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	//if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
	//	return
	//}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

//nolint: unparam
func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:       fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMsgSize: cfg.GRPCOptions.MaxMsgSize,
		// ServerCert:   cfg.SecureServing.ServerCert,
		mysqlOptions: cfg.MySQLOptions,
		// etcdOptions:      cfg.EtcdOptions,
	}, nil
}

//func (s *apiServer) initRedisStore() {
//	ctx, cancel := context.WithCancel(context.Background())
//	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
//		cancel()
//
//		return nil
//	}))
//
//	config := &storage.Config{
//		Host:                  s.redisOptions.Host,
//		Port:                  s.redisOptions.Port,
//		Addrs:                 s.redisOptions.Addrs,
//		MasterName:            s.redisOptions.MasterName,
//		Username:              s.redisOptions.Username,
//		Password:              s.redisOptions.Password,
//		Database:              s.redisOptions.Database,
//		MaxIdle:               s.redisOptions.MaxIdle,
//		MaxActive:             s.redisOptions.MaxActive,
//		Timeout:               s.redisOptions.Timeout,
//		EnableCluster:         s.redisOptions.EnableCluster,
//		UseSSL:                s.redisOptions.UseSSL,
//		SSLInsecureSkipVerify: s.redisOptions.SSLInsecureSkipVerify,
//	}
//
//	// try to connect to redis
//	go storage.ConnectToRedis(ctx, config)
//}
