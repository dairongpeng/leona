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

package gstash

import (
	"context"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"

	"sync"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"github.com/dairongpeng/leona/internal/gstash/analytics"
	"github.com/dairongpeng/leona/internal/gstash/config"
	"github.com/dairongpeng/leona/internal/gstash/gstashs"
	"github.com/dairongpeng/leona/internal/gstash/options"
	"github.com/dairongpeng/leona/internal/gstash/storage"
	"github.com/dairongpeng/leona/internal/gstash/storage/redis"
	"github.com/dairongpeng/leona/pkg/log"
)

var pmps []gstashs.Gstash

type gstashServer struct {
	secInterval int
	// 是否开启脱敏
	omitDetails    bool
	mutex          *redsync.Mutex
	analyticsStore storage.AnalyticsStorage
	gstash         map[string]options.GStashConfig
}

// preparedGenericAPIServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedGstashServer struct {
	*gstashServer
}

func createGstashServer(cfg *config.Config) (*gstashServer, error) {
	// use the same redis database with authorization log history
	// 创建redis客户端
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisOptions.Host, cfg.RedisOptions.Port),
		Username: cfg.RedisOptions.Username,
		Password: cfg.RedisOptions.Password,
	})

	// 基于客户端建立连接池，基于连接池建立并发安全的分布式锁redsync
	rs := redsync.New(goredis.NewPool(client))

	server := &gstashServer{
		secInterval: cfg.PurgeDelay,
		omitDetails: cfg.OmitDetailedRecording,
		// 分布式锁
		mutex: rs.NewMutex("leona-gstash", redsync.WithExpiry(10*time.Minute)),
		// input的数据源
		analyticsStore: &redis.RedisClusterStorageManager{},
		// 输出源配置列表
		gstash: cfg.Gstashs,
	}

	// 初始化redis的配置
	if err := server.analyticsStore.Init(cfg.RedisOptions); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *gstashServer) PrepareRun() preparedGstashServer {
	// 初始化配置的插件信息
	s.initialize()

	return preparedGstashServer{s}
}

func (s preparedGstashServer) Run(stopCh <-chan struct{}) error {
	ticker := time.NewTicker(time.Duration(s.secInterval) * time.Second)
	defer ticker.Stop()

	// 启动消费redis的数据
	log.Info("Now run loop to clean data from redis")
	for {
		select {
		// 按ticker计时器来定时消费。当收到操作系统中断信号后，执行case <-stopCh后直接跳出for循环。当前进程数据停止消费，数据不会在当前进程的内存中丢失
		case <-ticker.C:
			s.gStash()
		// exit consumption cycle when receive SIGINT and SIGTERM signal
		case <-stopCh:
			log.Info("stop purge loop")

			return nil
		}
	}
}

// gStash get authorization log from redis and write to gstash.
func (s *gstashServer) gStash() {
	// 分布式锁控制，保证gstash可以横向扩展
	if err := s.mutex.Lock(); err != nil {
		log.Info("there is already an leona-gstash instance running.")

		return
	}
	defer func() {
		if _, err := s.mutex.Unlock(); err != nil {
			log.Errorf("could not release iam-gStash lock. err: %v", err)
		}
	}()

	// 消费后删除，和kafka不同的是, kafka可以基于一批数据多次消费
	analyticsValues := s.analyticsStore.GetAndDeleteSet(storage.AnalyticsKeyName)
	if len(analyticsValues) == 0 {
		return
	}

	// Convert to something clean
	keys := make([]interface{}, len(analyticsValues))

	// 拿到队列中的每一条记录，细粒度的处理。压缩，脱敏等
	for i, v := range analyticsValues {
		decoded := analytics.AnalyticsRecord{}
		err := msgpack.Unmarshal([]byte(v.(string)), &decoded)
		log.Debugf("Decoded Record: %v", decoded)
		if err != nil {
			log.Errorf("Couldn't unmarshal analytics data: %s", err.Error())
		} else {
			if s.omitDetails { // 简单的脱敏
				decoded.Policies = ""
				decoded.Deciders = ""
			}
			keys[i] = interface{}(decoded)
		}
	}

	// Send to gstash
	writeToGStashs(keys, s.secInterval)
}

func (s *gstashServer) initialize() {
	// 该server配置了多少个收集器插件，初始化这些收集器
	pmps = make([]gstashs.Gstash, len(s.gstash))
	i := 0
	for key, pmp := range s.gstash {
		gstashTypeName := pmp.Type
		if gstashTypeName == "" {
			gstashTypeName = key
		}

		// 拿到当前配置的写入插件
		pmpType, err := gstashs.GetGstashByName(gstashTypeName)
		if err != nil {
			log.Errorf("Gstash load error (skipping): %s", err.Error())
		} else {
			// New一个插件实例
			pmpIns := pmpType.New()
			// 根据插件元数据配置，初始化
			initErr := pmpIns.Init(pmp.Meta)
			if initErr != nil {
				log.Errorf("Gstash init error (skipping): %s", initErr.Error())
			} else {
				log.Infof("Init Gstash: %s", pmpIns.GetName())
				// 设置过滤规则
				pmpIns.SetFilters(pmp.Filters)
				pmpIns.SetTimeout(pmp.Timeout)
				pmpIns.SetOmitDetailedRecording(pmp.OmitDetailedRecording)
				pmps[i] = pmpIns
			}
		}
		i++
	}
}

func writeToGStashs(keys []interface{}, purgeDelay int) {
	// Send to gstash
	if pmps != nil {
		var wg sync.WaitGroup
		wg.Add(len(pmps))
		// 按照写入插件做并发写入
		for _, pmp := range pmps {
			go execGStashWriting(&wg, pmp, &keys, purgeDelay)
		}
		// 等待一轮写入结束
		wg.Wait()
	} else {
		log.Warn("No gstash defined!")
	}
}

func filterData(gstash gstashs.Gstash, keys []interface{}) []interface{} {
	filters := gstash.GetFilters()
	if !filters.HasFilter() && !gstash.GetOmitDetailedRecording() {
		return keys
	}
	filteredKeys := keys[:] // nolint: gocritic
	newLenght := 0

	for _, key := range filteredKeys {
		decoded, _ := key.(analytics.AnalyticsRecord)
		if gstash.GetOmitDetailedRecording() {
			decoded.Policies = ""
			decoded.Deciders = ""
		}
		// 记录需要被过滤掉
		if filters.ShouldFilter(decoded) {
			continue
		}
		filteredKeys[newLenght] = decoded
		newLenght++
	}
	filteredKeys = filteredKeys[:newLenght]

	return filteredKeys
}

// execGStashWriting 每个插件的具体写入规则
func execGStashWriting(wg *sync.WaitGroup, pmp gstashs.Gstash, keys *[]interface{}, purgeDelay int) {
	timer := time.AfterFunc(time.Duration(purgeDelay)*time.Second, func() {
		if pmp.GetTimeout() == 0 {
			log.Warnf(
				"Gstash %s is taking more time than the value configured of purge_delay. You should try to set a timeout for this gStash.",
				pmp.GetName(),
			)
		} else if pmp.GetTimeout() > purgeDelay {
			log.Warnf("Gstash %s is taking more time than the value configured of purge_delay. You should try lowering the timeout configured for this gStash.", pmp.GetName())
		}
	})
	defer timer.Stop()
	defer wg.Done()

	log.Debugf("Writing to: %s", pmp.GetName())

	ch := make(chan error, 1)
	var ctx context.Context
	var cancel context.CancelFunc
	// Initialize context depending if the gStash has a configured timeout
	if tm := pmp.GetTimeout(); tm > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(tm)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	defer cancel()

	go func(ch chan error, ctx context.Context, pmp gstashs.Gstash, keys *[]interface{}) {
		filteredKeys := filterData(pmp, *keys)

		ch <- pmp.WriteData(ctx, filteredKeys)
	}(ch, ctx, pmp, keys)

	select {
	// err 不管是否为nil，都会留存在当前ch中
	case err := <-ch:
		if err != nil {
			log.Warnf("Error Writing to: %s - Error: %s", pmp.GetName(), err.Error())
		}
	case <-ctx.Done():
		//nolint: errorlint
		switch ctx.Err() {
		case context.Canceled:
			log.Warnf("The writing to %s have got canceled.", pmp.GetName())
		case context.DeadlineExceeded:
			log.Warnf("Timeout Writing to: %s", pmp.GetName())
		}
	}
}
