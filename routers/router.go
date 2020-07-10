package routers

import (
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/controllers/admin"
	"github.com/beatrice950201/araneid/controllers/index"
)

func init() {

	beego.ErrorController(&controllers.Error{})

	// 前台
	beego.Router("/", &index.Index{}, "get:Index")
	beego.AddNamespace(
		beego.NewNamespace("/index", beego.NSInclude(
			&index.Index{},
			&index.Spider{},
			&index.Lists{},
			&index.Detail{},
		)),
	)

	// 后台
	beego.AddNamespace(
		beego.NewNamespace("/admin", beego.NSInclude(
			&admin.Admin{}, &admin.Sign{}, &admin.Menu{},
			&admin.Roles{}, &admin.Users{}, &admin.Attachment{},
			&admin.Config{}, &admin.Inform{},
			&admin.Collect{}, &admin.Collector{},
			&admin.Dictionaries{}, &admin.Lexicon{},
			&admin.Models{}, &admin.Disguise{}, &admin.Template{},
			&admin.Class{}, &admin.Article{}, &admin.Prefix{}, &admin.Match{},
			&admin.Arachnid{}, &admin.Keyword{}, &admin.Indexes{}, &admin.Statistics{},
			&admin.Adapter{}, &admin.Domain{}, &admin.Category{}, &admin.Detail{},
			&admin.Automatic{}, &admin.Journal{}, &admin.Advantage{}, &admin.News{},
			&admin.Team{},
		)),
	)
}
