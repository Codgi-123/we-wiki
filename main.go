package main

import (
	_ "github.com/Codgi-123/we-wiki/app"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/memcache"
	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/astaxie/beego/session/redis_cluster"
)

func main() {
	beego.Run()
}
