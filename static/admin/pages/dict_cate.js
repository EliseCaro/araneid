const dict_cate = {
    option :{},
    begin  :{
        handle_collect:function () {
          $(".handle_collect").click(function () {
            application.ajax.post($(this).data("action"),{
                status : $(this).data("status"),
            },function (result) {
                application.ajax.requestBack(result.message,result.status,result.url);
            })
          });
        }
    }
}
$(document).ready(function() {
     dict_cate.begin.handle_collect()
});