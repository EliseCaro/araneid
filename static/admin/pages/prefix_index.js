const prefix_index = {
     option :{
          model  : 0,
          action : "",
          edit   : "",
          delete : "",
          html:`<div class="item_prefix form-group col-sm-5 p-1 badge-primary mb-3" data-key="[ID]">
                 <div class="input-daterange input-group">
                     <input type="text" class="form-control form-control-sm" name="title" placeholder="名字">
                     <div class="input-group-prepend input-group-append">
                        <span class="input-group-text font-w600"><i class="fa fa-fw fa-link"></i></span>
                     </div>
                     <input type="text" class="form-control form-control-sm col-sm-3" name="tags" placeholder="前缀">
                 </div>
                 <div class="prefix_delete"><i class="si si-close"></i></div>
           </div>`,
     },
     begin:{
          initialized_model:function () {
               prefix_index.option.model  = $("input[name='model']").val();
               prefix_index.option.action =  $(".item_create").data("action");
               prefix_index.option.edit   =  $(".items_prefix").data("edit");
               prefix_index.option.delete =  $(".items_prefix").data("delete");
          },
          initialized_create:function () {
               $(".item_create").click(function () {
                    application.ajax.post(prefix_index.option.action,{
                         model : prefix_index.option.model,
                         hideLoader:true
                    },function (result) {
                         if(result && result.status === false) {
                              application.ajax.requestBack(result.message,result.status,result.url);
                         }else {
                              let html = prefix_index.option.html.replace(/\[ID\]/g,result.data);
                              $(".items_prefix").prepend(html)
                         }
                    })
               });
          },
          initialized_input:function () {
               $("body").on("blur",".form-control",function () {
                    const value = $(this).val(),name = $(this).attr("name");
                    const key = $(this).parent().parent().data("key");
                    if (value) {
                         let option = {id : key,field:name,hideLoader:true};name === "title" ? option["title"] = value : option["tags"] = value;
                         application.ajax.post(prefix_index.option.edit,option,function (result) {
                              if(result && result.status === false) {
                                   application.ajax.requestBack(result.message,result.status,result.url);
                              }
                         })
                    }
             });
          },
          initialized_delete:function () {
               $("body").on("click",".prefix_delete",function () {
                  const object = $(this).parent();
                    application.ajax.post(prefix_index.option.delete, {id:object.data("key"),hideLoader:true},function (result) {
                         if(result && result.status === false) {
                              application.ajax.requestBack(result.message,result.status,result.url);
                         }else {
                              object.remove();
                         }
                    })
               })
          }
     },
}
$(document).ready(function() {
     prefix_index.begin.initialized_model();
     prefix_index.begin.initialized_create();
     prefix_index.begin.initialized_input();
     prefix_index.begin.initialized_delete();
});