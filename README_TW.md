# Golang 後端模板

中文 | [English](README.md)

### 包名 string_backend_0001

要使用前建議使用全局替換 "string_backend_0001" 成自己的包名稱

## 已安裝的包

* [gin](https://github.com/gin-gonic/gin)
* [gorm](https://gorm.io/index.html)
* [swag](https://github.com/swaggo/swag)

## 特色

* 上手簡單
* 附常用方法 internal/pkg
* 內置簡易 logger internal/logger
* 內建 OAuth2

## OAuth2 支持

* Line
* Google
* Discord

## 使用

### 初始化

```shell
go install github.com/swaggo/swag/cmd/swag@latest
git clone https://github.com/Wuchieh/backend.git
cd backend
go mod init
```

### 執行

```shell
swag init --parseDependency --parseInternal
go build . -o backend && ./backend
```

## 注意!
internal/web/oauth/oauth.go 僅測試用 請勿直接使用