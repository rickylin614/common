package credis

import "github.com/go-redis/redismock/v8"

var Mock redismock.ClientMock

func NewRedisMock() {
	client, Mock = redismock.NewClientMock()
	isCluster = false
}

func GetMock() redismock.ClientMock {
	return Mock
}
