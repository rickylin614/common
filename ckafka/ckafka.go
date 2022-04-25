package ckafka

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rickylin614/common/zlog"
	"github.com/segmentio/kafka-go"
)

type IManager interface {
	NewReader(topic, groupId string) (<-chan kafka.Message, error)
	SetBrokers(broker []string)
	SetLeaderAddr(lead string)
	WriteMultiTopic(key, value []byte, topics []string) error
	Write(key, value []byte, topic string) error
}

type Manager struct {
	brokers    []string
	leaderAddr string
}

var Manage IManager = &Manager{}

func (this *Manager) SetBrokers(broker []string) {
	this.brokers = broker
}

func (this *Manager) SetLeaderAddr(lead string) {
	this.leaderAddr = lead
}

// 閱讀器 groupId為空可所有程序皆可接收指定訊息
func (this *Manager) NewReader(topic, groupId string) (<-chan kafka.Message, error) {
	// 給brokers值
	if this.brokers == nil || len(this.brokers) == 0 {
		return nil, errors.New("not setting kafka.brokers")
	}

	// 設定kafka連線數據
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   this.brokers,
		Topic:     topic,
		GroupID:   groupId,
		Partition: 0,
		MinBytes:  1,                      // 1B
		MaxBytes:  10e6,                   // 10MB
		MaxWait:   time.Millisecond * 500, // 500ms
	})
	zlog.Info("ready to read", topic)

	// 創建訊息接收通道
	msgChan := make(chan kafka.Message)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// 建立協程接收訊息
	go func() {
	loop:
		for {
			select {
			case <-sigchan: //收到關閉訊號 直接退出不再接收
				break loop
			default:
				m, err := r.ReadMessage(context.Background())
				if err != nil {
					zlog.Error("read msg fail err:", err)
					continue loop
				}
				msgChan <- m
			}
		}
		if err := r.Close(); err != nil {
			zlog.Error("kafka read close fail", err)
			return
		}
	}()
	return msgChan, nil
}

/* 寫入多個topic */
func (this *Manager) WriteMultiTopic(key, value []byte, topics []string) error {
	if this.leaderAddr == "" {
		return errors.New("not setting kafka.leader")
	}

	// 寫入器基礎設定
	w := &kafka.Writer{
		Addr:         kafka.TCP(this.leaderAddr),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: time.Millisecond * 100,
	}

	// 組合要送的訊息
	kmsg := make([]kafka.Message, 0)
	for _, v := range topics {
		k := kafka.Message{
			Topic: v,
			Key:   key,
			Value: value,
		}
		kmsg = append(kmsg, k)
	}

	// 送出訊息
	err := w.WriteMessages(context.Background(), kmsg...)
	if err != nil {
		zlog.Error("failed to write messages:", err)
		return err
	}

	if err := w.Close(); err != nil {
		zlog.Error("failed to close writer:", err)
		return err
	}
	return nil
}

// 寫速單個topic裡面
func (this *Manager) Write(key, value []byte, topic string) error {
	if this.leaderAddr == "" {
		return errors.New("not setting kafka.leader")
	}

	// 寫入器基礎設定
	w := &kafka.Writer{
		Topic:        topic,
		Addr:         kafka.TCP(this.leaderAddr),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: time.Millisecond * 100,
	}

	// 送出訊息
	err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		zlog.Error("failed to write messages:", err)
		return err
	}
	if err := w.Close(); err != nil {
		zlog.Error("failed to close writer:", err)
		return err
	}
	return nil
}
