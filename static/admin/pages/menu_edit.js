const menu_edit = {
    cms:{
       child_menus:function (module) {
          const object = $("#module");
          const datas  = { active : parseInt(object.data("active")), module : module };
          const func   = function (res) {
              application.cms.loader("hide");
              if(res && res.status === false){
                  application.cms.alert({
                      title : application.cms.requestStatus.errorTitle,
                      text  : res.message,
                      type  : "error",
                      timer : application.cms.requestStatus.layerTime
                  });
                  return false;
              }
              res &&  menu_edit.cms.child_html(module,res.data || [],datas.active);
          }
          application.ajax.post(object.data("action"), datas, func)
       },
        child_html:function (module,res = [],active) {
          let html  = "<option value='"+ module +"'>顶级栏目</option>";
              html += menu_edit.cms.createHtmlNode(res,0,active);
          $("#pid").html(html);
        },
        createHtmlNode:function (res,level,pid) {
            let html = '';
            for (let i in res) {
                let prefix_html = "";
                if ( level > 0 ) {
                    for (let k = 0; k < level * 3; k++) {
                        prefix_html += "&nbsp;";
                    }
                    prefix_html += "┝";
                }
                let item = res[i];
                let selected = item.id === pid ? "selected" : "";
                html += '<option '+selected+' value="'+ item.id +'">'+prefix_html + item.title+'</option>';
                if(item.child && (level + 1) < 2 ){
                    html += menu_edit.cms.createHtmlNode(item.child,level + 1,pid);
                }
            }
            return html;
        }
    },
    begin:{
        wizard:function () {
            $('.js-wizard-simple').bootstrapWizard({
                nextSelector     : '[data-wizard="next"]',
                previousSelector : '[data-wizard="prev"]',
                finishSelector   : '[data-wizard="finish"]',
                onTabShow: (tab, nav, index) => {
                    let percent = ((index + 1) / nav.find('li').length) * 100;
                    let progress = nav.parents('.block').find('[data-wizard="progress"] > .progress-bar');
                    if (progress.length) {
                        progress.css({ width: percent + 1 + '%' });
                    }
                }
            });
        },
        render_child:function () {
          const module = parseInt($("#module").val());
          if (module > 0) {
              menu_edit.cms.child_menus(module)
          }else {
              $("#pid").html("<option value='0' selected>顶级模块</option>");
          }
        },
        change_module:function () {
            $("#module").change(function () {
               menu_edit.begin.render_child();
            })
        }

    }
}
$(document).ready(function() {
    menu_edit.begin.wizard();
    menu_edit.begin.render_child();
    menu_edit.begin.change_module();
});