package _func

import (
	"encoding/binary"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/extend/cache"
	"github.com/beatrice950201/araneid/extend/model/attachment"
	"github.com/beatrice950201/araneid/extend/model/config"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/** Ip2long 将 IPv4 字符串形式转为 uint32 **/
func Ip2long(ipString string) uint32 {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

/** uint32 字符串形式转为 IPV4 **/
func Long2IP(intIP uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

/** 解析debug为bool格式 **/
func AnalysisDebug() bool {
	t := beego.AppConfig.String("runmode")
	if t == "dev" {
		return true
	} else {
		return false
	}
}

/**  文件是否存在 **/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/** 获取加载的扩展模板 **/
func LayoutSections(cname, aname, module string) map[string]string {
	root := beego.AppConfig.String("admin_root")
	cname = strings.Replace(snakeString(cname), "_", "/", -1)
	aname = strings.ToLower(aname)
	res := map[string]string{
		"header": module + "/layout/header/" + cname + root + aname + ".html",
		"footer": module + "/layout/footer/" + cname + root + aname + ".html",
	}
	for k, v := range res {
		if b, _ := PathExists("./views/" + v); b == false {
			res[k] = ""
		}
	}
	return res
}

/** 写入缓存 **/
func SetCache(key string, val interface{}) error {
	return cache.Bm.Put(key, val, 86400*time.Second)
}

/** 读取缓存 **/
func GetCache(key string) interface{} {
	return cache.Bm.Get(key)
}

/** 判断数组是否存在某个值 **/
func InArray(need int, needArr []int) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

/** 获取配置每页显示多少条 **/
func WebPageSize() int {
	o, _ := beego.AppConfig.Int("system_page_size")
	return o
}

/** 解析字符串为map **/
func ParseAttrConfigMap(s string) map[string]string {
	maps := make(map[string]string)
	spaceRe, _ := regexp.Compile("[,;\\r\\n]+")
	res := spaceRe.Split(strings.Trim(s, ",;\r\n"), -1)
	if strings.Index(s, ":") > 0 {
		for _, v := range res {
			if countSplit := strings.Split(v, ":"); len(countSplit) == 2 {
				maps[countSplit[0]] = countSplit[1]
			}
		}
	}
	return maps
}

/** 解析字符串为map **/
func ParseAttrConfigArray(s string) []string {
	spaceRe, _ := regexp.Compile("[,;\\r\\n]+")
	return spaceRe.Split(strings.Trim(s, ",;\r\n"), -1)
}

/** todo 对接云端路径 获取资源前缀 **/
func DomainStatic(driver string) string {
	if driver == "local" {
		return "/"
	} else {
		return "http://www.eibk.com/"
	}
}

/** 获取缓存配置 **/
func CacheConfig() (c map[string]string) {
	return resolverConfig()
}

/** 获取一个随机数 **/
func RandomString() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

/** 获取一个文件地址 **/
func FilePath(id int) string {
	var item attachment.Attachment
	_ = orm.NewOrm().QueryTable(new(attachment.Attachment)).Filter("id", id).One(&item)
	return DomainStatic(item.Driver) + item.Path
}

/*********************************以下为私有方法******************************/

/** 获取全部配置并解析 **/
func resolverConfig() map[string]string {
	var (
		list         []*config.Config
		resolverList = make(map[string]string)
	)
	_, _ = orm.NewOrm().QueryTable(new(config.Config)).Filter("status", 1).All(&list)
	for _, v := range list {
		if v.Form == "image" && v.Value != "" {
			i, _ := strconv.Atoi(v.Value)
			v.Value = FilePath(i)
		}
		resolverList[v.Class+"_"+v.Name] = v.Value
	}
	return resolverList
}

/** 驼峰转下划线 **/
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
