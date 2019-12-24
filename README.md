# Yishan api-template

<img src="./assets/logo.svg" alt="drawing" width="50"/>

This is a api template repo used to start a new api project fast like a flash

It uses the following stack:
- go module, package management, https://github.com/golang/go/wiki/Modules
- gin, light weight http framework, https://github.com/gin-gonic/gin/
- air, live-reload https://github.com/cosmtrek/air
- resty, restful client, https://github.com/go-resty/resty
- redis, redis client, https://github.com/go-redis/redis
- gorm, ORM framework, http://gorm.io/docs/index.html
- jwt-go, Auth lib, https://github.com/dgrijalva/jwt-go
- logrus, logging lib, https://github.com/sirupsen/logrus
- viper, config lib, https://github.com/spf13/viper
- ginkgo&gomega, BDD test framework, https://github.com/onsi/ginkgo, https://github.com/onsi/gomega

## Source Layout

Follow https://github.com/golang-standards/project-layout standards

## Architecture

We use Hexagonal Architecture, check /examples/hex for more information.

## Get Start

1. go mod download

2. get air
```bash
curl -fLo ./bin/air \
    https://raw.githubusercontent.com/cosmtrek/air/master/bin/{linux, darwin, windows}/air
chmod +x ~/air
```

3. make dev
###