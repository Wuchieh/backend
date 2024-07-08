# Go後端模板

### 包名 string_backend_0001
```
要使用前建議使用全局替換 "string_backend_0001" 成自己的包名稱
```

# 已安裝的包
* [gin](https://github.com/gin-gonic/gin)
* [gorm](https://gorm.io/index.html)
* [swag](https://github.com/swaggo/swag)

# 使用
## 初始化
```shell
go mod tidy
go install github.com/swaggo/swag/cmd/swag@latest
```
## 執行
```shell
swag init
go build . -o backend && ./backend
```