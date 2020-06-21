const module_upload = {
    option:{
        multiple_object: $(".upload_multiple"),
        single_object  : $(".upload_single"),
        onLoad_class   : ".upload_onLoad",
        previewHeight  : "90px",
        previewWidth   : "90px",
        dragdropWidth  : "90px",
        single_button  :'<div class="upload-dragdrop-plus"><i class="far fa-2x fa-images"></i></div>',
        uploadObject   : {},
        upload_onLoad:function (object) {
            const items_object = module_upload.option.single_object.parent().children(this.onLoad_class);
            items_object.children().each(function(i,n){
                object.createProgress($(n).data("name"),$(n).data("path"),$(n).data("size"));
            });
        },
        upload_success:function (data,pd,type) {
            if(data && data.status === false){
                application.cms.alert({
                    title : application.cms.requestStatus.errorTitle,
                    text  : data.message,
                    type  : "error",
                    timer : application.cms.layerTime
                },function () {
                    pd.statusbar.remove();
                    module_upload.option.uploadObject.reset(false);
                })
            }else {
                const items_inputs = module_upload.option.single_object.parent().children("input[type='hidden']");
                if(type === "single"){
                    items_inputs.val(data.data.id)
                }
            }
        },
        delete_value:function (data,pd,type){
            const items_inputs = module_upload.option.single_object.parent().children("input[type='hidden']");
            if(type === "single"){
                items_inputs.val(0)
            }
        }
    },
    begin :{
        single:function () {
            module_upload.option.uploadObject = module_upload.option.single_object.uploadFile({
                url:application.cms.uploadImageSingle,
                uploadStr:module_upload.option.single_button,
                dragDropStr:false,
                maxFileCount:1,
                showPreview:true,
                fileName:module_upload.option.single_object.data("name"),
                previewHeight: module_upload.option.previewHeight,
                previewWidth: module_upload.option.previewWidth,
                dragdropWidth : module_upload.option.dragdropWidth,
                showDelete: true,
                deleteStr:'<i class="si si-close"></i>',
                showAbort:false,
                showError:false,
                showFileSize:false,
                showProgress:true,
                onLoad:function(obj){
                    module_upload.option.upload_onLoad(obj)
                },
                onSuccess:function(files,data,xhr,pd){
                    module_upload.option.upload_success(data,pd,"single")
                },
                deleteCallback:function (data,pd) {
                    module_upload.option.delete_value(data,pd,"single")
                }
            })
        },
        multiple:function () {}
    }
}

$(document).ready(function() {
    module_upload.begin.single(); // 初始化单文件上传
    module_upload.begin.multiple(); // 初始化多文件上传
});
