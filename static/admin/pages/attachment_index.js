const attachment_index = {
    option:{},
    begin:{
        spaceCount:function () {
            $(".space_count").click(function () {
                application.ajax.post($(this).data("action"),{},function (e) {
                    application.cms.loader("hide");
                    const requestStatus = application.cms.requestStatus;
                    const option = {
                        title : (e.status === true) ? "空间计算结果" : requestStatus.errorTitle,
                        text  : e.message,
                        type  : (e.status === true) ? "success" : "error",
                        timer : requestStatus.layerTime * 5
                    }
                    application.cms.alert(option,function () {
                      swal.close();
                    })
                });
            });
        },
    }
}
$(document).ready(function() {
    attachment_index.begin.spaceCount();
});