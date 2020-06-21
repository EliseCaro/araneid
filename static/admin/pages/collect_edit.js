const collect_edit = {
    option:{
        filedCreateObject:$(".filed_create"),
        filedCreateBox   :$(".field_create_child_box"),
        filedCreateChild :$(".field_create_child_string").html(),
        postSubmitForm   :$(".ajax_from_custom"),
        toolBatchUrls    :$(".batch_urls_btn"),
        toolTestRetrieve :$(".modal_source_rule_test_btn"),
        toolTestMatching :$(".matching_test"),
        getMaxKey:function (type) {
          const k = collect_edit.option.filedCreateBox.attr("data-key");
          const n = type === "inc" ? parseInt(k) +1 : parseInt(k) -1;
          collect_edit.option.filedCreateBox.attr("data-key",n)
          return n
        },
        serializeArray:function () {
            const a = collect_edit.option.postSubmitForm.serializeArray()
            let   form = {};
            for (let i in a){form[a[i].name] = a[i].value}
            return form
        },
        matchingArray:function () {
            let   formArray = [];
            collect_edit.option.filedCreateBox.children().each(function () {
                let single = {}
                $(this).find("input").each(function () {
                    single[$(this).data("name")] = $(this).val()
                })
                $(this).find("select").each(function () {
                    single[$(this).data("name")] = parseInt($(this).val())
                })
                formArray.push(single)
            })
            return formArray.length > 0 ? formArray : "";
        },
        modalSourceRule:function (result) {
          let html = "";
          for (let i in result){
              html += "<tr>";
              html += '<td class="table-success text-center p-1"><i class="fa fa-fw fa-link text-success"></i></td>';
              html += '<td class="p-1 pl-2"><a target="_blank" href="'+ result[i].href +'">'+ result[i].title +'</a></td>';
              html += '<td class="text-right p-1"><a class="btn btn-sm btn-alt-primary" style="width:100%" target="_blank" href="'+ result[i].href +'">预览</a></td>';
              html += "</tr>";
          }
          return html
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
        filedCreate:function () {
          collect_edit.option.filedCreateObject.click(function () {
             const k = collect_edit.option.getMaxKey("inc");
             const h = collect_edit.option.filedCreateChild.replace(/__KEY__/g,k)
              collect_edit.option.filedCreateBox.prepend(h)
          });
        },
        filedDelete:function () {
            $("body").on("click",".field_delete",function(){
                $(this).parent().remove();
            });
        },
        ajaxFromCustom:function () {
           collect_edit.option.postSubmitForm.submit(function () {
               const serialize = $.extend(collect_edit.option.serializeArray(),{matching:collect_edit.option.matchingArray()})
               if(serialize["matching"]){
                   serialize["matching"] = JSON.stringify(serialize["matching"])
               }
               application.ajax.post($(this).attr("action"),serialize,function (result) {
                   application.ajax.requestBack(result.message,result.status,result.url);
               })
               return false;
           })
        },
        toolBatchUrls:function () {
            collect_edit.option.toolBatchUrls.click(function () {
                const role_init  = $("input[name='role_init']").val(),
                      role_start = $("input[name='start_number']").val(),
                      role_end   = $("input[name='end_number']").val();
                if(!role_init || !role_start || !role_end){
                    swal("生成错误"," 参数未正确填写；请完整填写后再试～","info")
                }
                let valuesString = "";
                for (let i = role_start; i <= role_end; i++) {
                    valuesString += role_init.replace("(*)",i) + ((i < role_end) ?  "\n" : "");
                }
                $("#batch_urls_values").val(valuesString);
            });
        },
        toolTestRetrieve:function () {
            collect_edit.option.toolTestRetrieve.click(function () {
                application.ajax.post($(this).data("action"),{
                    source      : $("textarea[name='source']").val(),
                    source_rule : $("input[name='source_rule']").val(),
                },function (result) {
                    application.cms.loader("hide");
                    if(result.status === true && result.data && result.data.length > 0){
                        $(".modal_source_rule_box").html(collect_edit.option.modalSourceRule(result.data))
                        $('#modal_source_rule_test').modal('show');
                    }else {
                        application.cms.alert({
                            title : "检索测试结果",
                            text  : "未检索到任何匹配结果；请检查检索条件！",
                            type  : "error",
                            timer : application.cms.requestStatus.layerTime
                        });
                    }
                })
            })
        },
        toolTestMatching:function () {
            collect_edit.option.toolTestMatching.click(function () {
                application.ajax.post($(this).data("action"),{
                    matching      : JSON.stringify(collect_edit.option.matchingArray()),
                    source        : $("textarea[name='source']").val(),
                    source_rule   : $("input[name='source_rule']").val(),
                },function (result) {
                    application.cms.loader("hide");
                    if(result.status === true && result.data){
                        $("#response_headers").JSONView(result.data);
                        $('#modal_test_field').modal('show');
                    }else {
                        application.cms.alert({
                            title : "字段测试结果",
                            text  : result.message ? result.message : "未检索到任何匹配字段；请检查检索条件！",
                            type  : "error",
                            timer : application.cms.requestStatus.layerTime
                        });
                    }
                })
            })
        }
    },
}
$(document).ready(function() {
    collect_edit.begin.wizard(); /*进度条*/
    collect_edit.begin.filedCreate(); /*添加字段解析*/
    collect_edit.begin.filedDelete(); /*删除字段解析*/
    collect_edit.begin.ajaxFromCustom(); /*提交创建更新*/
    collect_edit.begin.toolBatchUrls(); /*批量Url生成器*/
    collect_edit.begin.toolTestRetrieve();  /*采集连接验证器*/
    collect_edit.begin.toolTestMatching();  /*采集字段验证器*/
});