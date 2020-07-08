package main

import (
	"github.com/astaxie/beego"
	_ "github.com/beatrice950201/araneid/extend/begin"
	"github.com/beatrice950201/araneid/extend/cache"
	_ "github.com/beatrice950201/araneid/routers"
)

func main() {
	cache.Run()
	beego.Run()
}
