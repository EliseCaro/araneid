const domain_links = {
    option:{
        initializedValue:function () {
            let result = []
            $(".link_item").each(function(){
                const id = "#item_" + $(this).data("key")
                const title = $(id + ' input[name="title"]').val();
                const urls  = $(id + ' input[name="urls"]').val();
                if (title && urls){
                    result.push({title : title,urls  : urls})
                }
            })
           return result
        }
    },
    begin:{
        initialized:function () {
            $("body").on("click",".prefix_delete",function () {
                $(this).parent().remove();
            })
        },
        submitClick:function () {
            $(".submit_links").click(function () {
                const value = domain_links.option.initializedValue();
                application.ajax.post($(this).data("action"),{
                    inputs : JSON.stringify(value),
                })
            })
        }
    },
}
$(document).ready(function() {
    domain_links.begin.initialized();
    domain_links.begin.submitClick();
});