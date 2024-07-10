# Golang 後端模板

### 包名 string_backend_0001

```
要使用前建議使用全局替換 "string_backend_0001" 成自己的包名稱
```

## 已安裝的包

* [gin](https://github.com/gin-gonic/gin)
* [gorm](https://gorm.io/index.html)
* [swag](https://github.com/swaggo/swag)

## 特色

* 上手簡單
* 附常用方法 internal/pkg
* 內置簡易 internal/logger
* Google Oauth2 登入
* Line Oauth2 登入

## 使用

### 初始化

```shell
git clone https://github.com/Wuchieh/backend.git
go mod init
```

### 執行

```shell
swag init
go build . -o backend && ./backend
```