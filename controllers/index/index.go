package index

import (
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	service2 "github.com/beatrice950201/araneid/extend/service"
	"regexp"
)

type Index struct {
	Main
}

// @router / [get]
func (c *Index) Index() {
	beego.Info(_func.WebPageSize())
	beego.Info(regexp.MustCompile(`https://www.chazidian.com/([a-z]+)_([a-z]+)_([a-zA-Z0-9]{32})/$`).MatchString("https://www.chazidian.com/r_ci_a4d0a723bbf7bdb8e473c2737fc62255/"))
	service := service2.DefaultDictionariesService{}
	res := service.Translate(map[string]string{
		"context": `<h2 id="1">详细解释</h2><div class="tab-pane" id="c69364"><script type="text/javascript">tabPane = new WebFXTabPane( document.getElementById( "c69364" ), true );</script><div class="tab-page" id="c693643"><h2 class="tab">词语解释</h2><script type="text/javascript">tabPane.addTabPage( document.getElementById( "c693643" ) );</script><span class="dicpy">áng shǒu tǐng xiōng </span> <span class="diczy">ㄤˊ ㄕㄡˇ ㄊㄧㄥˇ ㄒㄩㄥ </span><p class="zdct7"><strong>昂首挺胸</strong>　</p><hr class="dichr"/><p class="zdct7">仰起头，挺起胸膛。形容无所畏惧或坚决的样子。 欧阳予倩 <span class="diczx1">《小英姑娘》</span>：“她伸开两手昂首挺胸，狂了似的往外跑。” 张书绅 <span class="diczx1">《正气歌》</span>：“她昂首挺胸，向吉普车走去。”亦作“ 昂头挺胸 ”。 沙汀 <span class="diczx1">《范老老师》</span>：“‘你看！提到内战，就连妇人女子都反对啦……’昂头挺胸，老老师欣喜地叫出来。”</p></div></div><div class="notice"><p>成语词典已有该词条：昂首挺胸</p></div>`,
	})
	beego.Info(res["context"])
	c.Data["Title"] = "home"
}

// @router /index/test [post]
func (c *Index) Test() {
	maps := make(map[string]interface{})
	if e := c.ParseForm(&maps); e != nil {
		beego.Warn(e)
	} else {
		beego.Info(maps)
	}
	//c.Succeed(&controllers.ResultJson{Message: " 发布数据成功拉！"})
	c.Fail(&controllers.ResultJson{Message: " 发布数据失败拉！"})
}
