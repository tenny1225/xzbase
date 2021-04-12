package conf

import "github.com/astaxie/beego/config"

var Config config.Configer

func Init(name string) {
	c, e := config.NewConfig("ini", name)
	if e != nil {
		panic(e)
	}
	Config = c
}
