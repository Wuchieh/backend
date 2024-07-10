# Golang Backend Template

[中文](README_TW.md) | English

### Package name: string_backend_0001

It is recommended to use global replacement of "string_backend_0001" with your own package name before use.


## Installed Packages

* [gin](https://github.com/gin-gonic/gin)
* [gorm](https://gorm.io/index.html)
* [swag](https://github.com/swaggo/swag)

## Features

* Easy to get started
* Includes common methods in internal/pkg
* Built-in simple Logger in internal/logger
* Google OAuth2 login
* Line OAuth2 login

## Usage

### Initialization

```shell
git clone https://github.com/Wuchieh/backend.git
cd backend
go mod init
```

### Execution

```shell
swag init
go build . -o backend && ./backend
```