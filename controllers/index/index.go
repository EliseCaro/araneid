package index

import (
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	service2 "github.com/beatrice950201/araneid/extend/service"
	"regexp"
)

type Index struct {
	Main
}

// 自然语言处理：
// 1：关键字重新提取
// 1: 从训练模型提取关键词随机近义词；如果没有进入流程2
// 2: 每个关键词提取1*n个近义词放入训练模型
// 3: 进入流程1重组关键词近义词
// 4： 返回该文章重组后的SEO关键词
// 2：内容描述重新提取
// 1：进行2*n次数的自然语言重组；
// 2：使用以上的关键词处理方式进行替换原文关键词；（如果不能保持语言自然则放弃）
// 3：返回该文章的SEO描述
// 3: 文章标题：
// 1：进行2*n次数的自然语言重组；
// 2： 直接返回
// 3: 文章内容：
// 1：过滤所有html标签；剩余纯文本以及标点符号；// todo 文章内容不考虑内容内嵌图等资源标签
// 2：进行2*n次数的自然语言重组；
// 3: 使用以上关键词以及训练模型进行关键词批量替换
// 3：根据标点符号重组html结构，句号为<P>
// 4：返回内容
//
// 内容模型挂载；
// 1： 文章分类；通过链接采集或手动添加->自然语言处理->入库到模型分类
// 2： 文章详情：通过挂载的采集器自动远行采集，自动发布->自然语言处理->入库->均分入库到模型文章；
// 3   挂载模板： 通过挂载模板分组以决定此模型的内容应该用那些模板；使用使用文件将响应的每个url以html文件的形式进行保存；解决链接内容不固定等问题；
// 蜘蛛池：
// 1：蜘蛛池挂载内容模型：
// 2：配置泛解析或者前缀解析；
// 3：设置内容模型的数据最大重复挂载次数；（重复挂载次数如：【在这个蜘蛛池每个文章不包含TKDB处理最多允许重复几次；】）
//    虽说重复但加入TKDB之后重复几率很小；
// 4：挂载TKDB模板；TKDB模板配置在原来文章上的随机组合方式；（通过蜘蛛池关键词库使用自定义标签组合）
// 5： 关键词库；用户导入关键词库，或者通过关键词训练模型提取

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

// @router /index/test [get]
func (c *Index) Test() {
	item := c.articleService.One(16)
	model := c.modelsService.One(item.Model)
	res, err := c.disguiseService.DisguiseHandleManage(model.Disguise, &spider.HandleModule{
		Title:       item.Title,
		Context:     item.Context,
		Keywords:    item.Keywords,
		Description: item.Description,
	})
	if err != nil {
		c.Fail(&controllers.ResultJson{Message: err.Error()})
	} else {
		c.Succeed(&controllers.ResultJson{Message: "success", Data: res})
	}
}
