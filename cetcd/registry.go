package cetcd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rickylin614/common/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 服务信息
type ServiceInfo struct {
	Name string
	IP   string
}

type Service struct {
	ServiceInfo ServiceInfo
	stop        chan error
	leaseId     clientv3.LeaseID
	client      *clientv3.Client
}

// NewService 创建一个注册服务
func NewService(name, ip string, endpoints []string) (service *Service, err error) {
	info := ServiceInfo{
		Name: name,
		IP:   ip,
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 10,
	})

	if err != nil {
		zlog.Fatal(err)
		return nil, err
	}

	service = &Service{
		ServiceInfo: info,
		client:      client,
	}
	return
}

// Start 注册服务启动
func (service *Service) Start() (err error) {
	ch, err := service.keepAlive()
	if err != nil {
		zlog.Fatal(err)
		return
	}
	//創建接收強制關閉的訊號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case err := <-service.stop:
			return err
		case <-service.client.Ctx().Done():
			return errors.New("service closed")
		case resp, ok := <-ch:
			// 监听租约
			if !ok {
				zlog.Info("keep alive channel closed")
				return service.revoke()
			}
			zlog.Debugf("Recv reply from service: %s, ttl:%d", service.getKey(), resp.TTL)
		case <-quit: //接收到就結束keepAlive避免新請求取得該連結依然健康
			return
		}
	}

}

func (service *Service) Stop() {
	service.stop <- nil
}

func (service *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	info := &service.ServiceInfo
	key := info.Name + "/" + info.IP
	// 创建一个租约
	resp, err := service.client.Grant(context.TODO(), 5)
	if err != nil {
		zlog.Fatal(err)
		return nil, err
	}
	_, err = service.client.Put(context.TODO(), key, info.IP, clientv3.WithLease(resp.ID))
	fmt.Printf("put key: %s and value :%s", key, info.IP)
	if err != nil {
		zlog.Fatal(err)
		return nil, err
	}
	service.leaseId = resp.ID
	return service.client.KeepAlive(context.TODO(), resp.ID)
}

func (service *Service) revoke() error {
	_, err := service.client.Revoke(context.TODO(), service.leaseId)
	if err != nil {
		zlog.Fatal(err)
	}
	zlog.Infof("servide:%s stop\n", service.getKey())
	return err
}

func (service *Service) getKey() string {
	return service.ServiceInfo.Name + "/" + service.ServiceInfo.IP
}
