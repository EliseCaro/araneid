package cache

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
)

var Bm cache.Cache

func Run() {
	cachePath := beego.AppConfig.String("cache_path")
	fileSuffix := beego.AppConfig.String("file_suffix")
	option := fmt.Sprintf(`{"CachePath":"%s","FileSuffix":"%s","DirectoryLevel":"2","EmbedExpiry":"120"}`, cachePath, fileSuffix)
	Bm, _ = cache.NewCache("file", option)
}
