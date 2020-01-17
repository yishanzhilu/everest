# Yishan Everest

<img src="./assets/logo.svg" alt="drawing" width="50"/>

Yishan Workspace API Server

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

## Config
Everest use config in ./configs folder, it use viper to read file name with "viper"
in any acceptable file type.
You can use ENV `EVEREST_CONFIG_FILE_NAME` to change the filename, but it has be in
./configs folder.

## Get Start

### DEV

1. go mod download

2. get air
```bash
curl -fLo ./bin/air \
    https://raw.githubusercontent.com/cosmtrek/air/master/bin/{linux, darwin, windows}/air
chmod +x ~/air
```

3. update config file
```bash
# copy a local config file which is ignored in git
# so no one upload any sensitive data to pulbic
cp ./configs/viper.yaml ./configs/viper.local.yaml
# edit viper.local.yaml with your own credentials
vi ./configs/viper.local.yaml
```
4. make dev


### PROD
