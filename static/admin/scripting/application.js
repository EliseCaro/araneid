const application = {
    cms:{
        tablePageSize       :     20,
        socketHost          :     'ws://' + window.location.host + '/admin/admin/socket',
        uploadImageMultiple :     "/admin/attachment/upload_images",
        uploadImageSingle   :     "/admin/attachment/upload_image",
        uploadFileSingle    :     "/admin/attachment/upload_file",
        requestStatus:{
             errorTitle    : "失败提示",
             successTitle  : "成功提示",
             errorMessage  : "请求失败;网络断开或服务器已停止;请稍后再试！",
             defaultStatus : false,
             defaultUrls   : "",
             layerTime     : 3000
         },
        authorizationMessage:{
            title : "删除警告",
            text  : "您正在删除条目；删除后将无法恢复！您确定删除吗？",
            confirmButtonText  : "确定",
            cancelButtonText   : "取消",
        },
        authorizationStatus:{
            title : "禁用警告",
            text  : "禁用条目后将在启用之前无法使用！您确定该操作吗？",
        },
         loader:function (type) {
             if(typeof(One) !== "undefined" && typeof(One) !== undefined){
                 One.loader(type)
             }
         },
         alert:function (option,func) {
             parent === self ? swal(option) : parent.swal(option);
             setTimeout(function(){
                 func && func();
             }, option.timer + 500);
         },
        authorization:function (option,func) {
            swal($.extend({
                type: "info",
                confirmButtonText: application.cms.authorizationMessage.confirmButtonText,
                cancelButtonText: application.cms.authorizationMessage.cancelButtonText,
                showCancelButton: true,
                closeOnConfirm: false,
                showLoaderOnConfirm: true,
            },option),func)
        },
        layout_theme:function (t) {
            const o = t || "dark";
            One.layout("header_style_"  + o)
            One.layout("sidebar_style_" + o)
        },
        layer_iframe:function (option = {}) {
            const options = {type: 2,scrollbar: false,offset: 'auto',zIndex:100,resize:false, shadeClose : true, anim: 1, fixed: false};
            return layer.open(
                $.extend({
                    title: option.title || false,
                    content: option.urls,
                    skin: option.skin || "layui-layer-molv",
                    area: option.area,
            },options));
        }
    },
    binds:{
        init:function () {
            $(".ajax_from").submit(function () {
                application.ajax.post(
                    $(this).attr("action"),
                    $(this).serialize(),
                );
                return false;
            });
            $("body").on("click",".open_iframe",function () {
                let option = {
                    urls : $(this).attr("href"),
                    area : $(this).data("area").split(","),
                    title: $(this).data("title") || false,
                    skin:  $(this).data("skin") || "",
                };
                if (option.area.length === 1){
                    option.area = option.area[0]
                }
                application.cms.layer_iframe(option)
                return false;
            });
            $("body").on("click",".ids_delete",function () {
                const object = $(this);
                application.cms.authorization({
                    title: application.cms.authorizationMessage.title,
                    text: application.cms.authorizationMessage.text,
                }, function () {
                    application.ajax.post(object.data("action"), {
                        ":ids": object.data("ids"),
                        hideLoader: true
                    });
                });
            })
            $("body").on("click",".ids_enable",function () {
                const object = $(this);
                application.ajax.post(object.data("action"),{
                    ":ids":object.data("ids"),
                    ":status":1,
                    ":field":object.data("field")
                });
            });
            $("body").on("click",".ids_disable",function () {
                const object = $(this);
                application.cms.authorization({
                    title : application.cms.authorizationStatus.title,
                    text  : application.cms.authorizationStatus.text,
                },function () {
                    application.ajax.post(object.data("action"),{
                        ":ids":object.data("ids"),
                        ":status":0,
                        ":field":object.data("field"),
                        hideLoader : true,
                    });
                });
            });
            $("body").on("click",".jump_urls",function(){
                window.location.href = $(this).data("action");
            });
            // 更多绑定项
        }
    },
    ajax:{
        post:function (url,data,func) {
            !data.hideLoader && application.cms.loader("show");
            $.ajax({
                type        : "POST",
                dataType    : "json",
                url         : url,
                data        : data,
                traditional : true,
                success     : function (result) {
                    func ? func(result) : application.ajax.requestBack(result.message||application.cms.requestStatus.errorMessage,result.status||application.cms.requestStatus.defaultStatus,result.url||application.cms.requestStatus.defaultUrls);
                },
                error       : function(data) {
                    console.log(data);
                    application.ajax.requestBack(application.cms.requestStatus.errorMessage,application.cms.requestStatus.defaultStatus,application.cms.requestStatus.defaultUrls);
                }
            })
        },
        requestBack:function (message,status,urls) {
            application.cms.loader("hide");
            const requestStatus = application.cms.requestStatus;
            const option = {
                 title : (status === true) ? requestStatus.successTitle : requestStatus.errorTitle,
                 text  : message,
                 type  : (status === true) ? "success" : "error",
                 timer : requestStatus.layerTime
             }
            application.cms.alert(option,function () {
                 if (urls) {
                     parent === self ? window.location.href = urls : parent.window.location.href = urls;
                 }
            });
        },
    }
};

$(document).ready(function(){
    application.binds.init();
    if(typeof(layout_theme_type) !== "undefined"){
        application.cms.layout_theme(layout_theme_type)
    }
});