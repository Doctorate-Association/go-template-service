package dal

import (
	"go-template-service/biz/dal/mysql"
	"go-template-service/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
