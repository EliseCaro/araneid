const adapter_index = {
    option :{},
    begin  :{
        handleForm:function () {
            $(".handleForm").submit(function () {
                application.ajax.post($(this).attr("action"), $(this).serialize(),function (result) {
                    application.ajax.requestBack(result.message,result.status,result.url);
                    if(result && result.data){
                        $(".download").attr("data-href",result.data).show();
                    }
                    if(result && result.status === true && !result.data){
                        $(".download").remove();
                    }
                });
                return false;
            });
        },
        download:function(){
            $(".download").click(function () {
                const url = $(this).data("href");
                if(url){ window.location.href = "/" + url }
            })
        },
        uploadCustom:function () {
            $('.upload_custom').uploadFile({
                url:application.cms.uploadFileSingle,
                uploadStr:'<div class="text-center py-4"><p><i class="fa fa-4x fa-cloud-upload-alt text-gray"></i></p><h4 class="mb-0 text-gray success_title">选择读取文件对象</h4></div>',
                dragDropStr:false, maxFileCount:1, showCancel:false,
                showAbort:false, showDone:false, showDelete:true,
                showStatusAfterSuccess:false, showError:false,
                showFileSize:false, showPreview:false, fileName:"files",
                onSuccess:function(files,data,xhr,pd){
                    if(data && data.status === true){
                        $('.upload_custom').parent().children("input[type='hidden']").val(data.data.id);
                        $(".success_title").html(data.data.name)
                    }else {
                        module_upload.option.upload_success(data,pd,"single");
                    }
                },
            });
        }
    }
}
$(document).ready(function() {
    adapter_index.begin.handleForm()
    adapter_index.begin.uploadCustom()
    adapter_index.begin.download()
});