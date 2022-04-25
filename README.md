## how to import this project

- please entry edit gitconfig

```shell 
vi ~/.gitconfig
```

- add below:

```config
[url "https://{git account}:FdPhTgC6ydx5y9fYRaGH@gftygitlab.awoplay.com/pegcornteam/serverteam/golang/common"]
  insteadOf = https://gl.ccc/common
```

- if you want to go get gl.ccc/common , please input it with .git 
  - Example: 
```shell
go get -insecure -u github.com/rickylin614/common # go 1.16 version

go get -u github.com/rickylin614/common # go 1.17 version
```

- if you use 1.17 or upper version, you must add GOINSECURE=gl.ccc* and GOPRIVATE=gl.ccc* in your system environment variables.

> if login fail , you can entry: "git clone http://gl.ccc/common" master
than type your account and password , you pc would save the login info.
after do it, try go get again.

## apollo 設定規範

### support

```
mysql,redis
```

- support設定的為初始化要自動執行設定的方法。

### mysql格式範例

```yml
mysql:
  -
    host: 127.0.0.1:10037
    schema: schema1
    user: root
    pwd: abcdefg
  -
    host: 127.0.0.1:10038
    schema: schema12
    user: root
    pwd: 1234567
    source: source2
```

### redis格式範例

```yml
host:
  - 127.0.0.1:6379
```

### redis cluster格式範例

```yml
host:
  - 127.0.0.1:7001
  - 127.0.0.1:7002
  - 127.0.0.1:7003
  - 127.0.0.1:7004
  - 127.0.0.1:7005
  - 127.0.0.1:7006
```

### log格式範例

```yml
infopath: "/logs/a/info.log"
errorpath: "/logs/a/error.log"
```

### etcdRegister

```yml
host: 127.0.0.1:1234
prefix: "service1/"
endpoints: "127.0.0.1:1001,127.0.0.1:1002,127.0.0.1:1003"
```

### etcdClient

```yml
endpoints: "127.0.0.1:1001,127.0.0.1:1002,127.0.0.1:1003"
```



