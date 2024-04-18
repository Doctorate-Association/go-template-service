# Go template microservice
This is a template for a microservice written in Go. It uses the [Hertz](https://github.com/cloudwego/hertz) framework underneath. and go-jwt to authenticate the requests was from authorised user.

## Pre requirements
- Go 1.22
- Docker
- Docker-compose
- Go mod

## Directory structure

|  catalog   | introduce  |
|  ----  | ----  |
| conf  | Configuration files |
| main.go  | Startup file |
| hertz_gen  | Hertz generated model |
| biz/handler  | Used for request processing, validation and return of response. |
| biz/service  | The actual business logic. |
| biz/dal  | Logic for operating the storage layer |
| biz/route  | Routing and middleware registration |
| biz/utils  | Wrapped some common methods |

## How to run

```shell
sh build.sh
sh output/bootstrap.sh
```
