const table_begin = {
    cms:{
        object       : $('.js-dataTable_init'),
        tableObject  : {},
        language     : {
            infoEmpty        : "",
            infoFiltered     : " 从_MAX_条数据中筛选",
            info             : "<span class='badge btn-alt-primary p-2'>当前第<strong>_PAGE_</strong>页 &nbsp; 总共<strong>_PAGES_</strong>页</span>",
            search           : "_INPUT_",
            zeroRecords      : "没有匹配到相关内容",
            searchPlaceholder: "输入搜索关键字",
            processing       : '<div class="spinner-grow text-primary" role="status"><span class="sr-only">加载中...</span></div>',
            paginate: {
                first    : '<i class="fa fa-angle-double-left"></i>',
                previous : '<i class="fa fa-angle-left"></i>',
                next     : '<i class="fa fa-angle-right"></i>',
                last     : '<i class="fa fa-angle-double-right"></i>'
            },
        },
        drawReset:function(res){
            application.cms.loader("hide");
            if(res.status === true){
                swal.close();
                table_begin.cms.tableObject.api().draw(true);
            }else {
                application.cms.alert({
                    title : application.cms.requestStatus.errorTitle,
                    text  : res.message,
                    type  : "error",
                    timer : application.cms.requestStatus.layerTime
                });
            }
        },
        checkbox_select:function(){
            let checkbox_vale = [];
            $(".js-dataTable_init input:checkbox[name='check[]']:checked").each(function(index){
                checkbox_vale[index] = parseInt($(this).attr("data-id"));
            });
            return checkbox_vale;
        },
        checkAllHtml:function (){
            let html =  '<div class="custom-control custom-checkbox d-inline-block">';
                html += '<input type="checkbox" class="custom-control-input" id="check-all" name="check-all">';
                html += '<label class="custom-control-label" for="check-all"></label>';
                html += '</div>';
            return html;
        },
        bindDom: "<'row'<'col-sm-12 col-md-6'B><'col-sm-12 col-md-6'f>><'row'<'col-sm-12'tr>><'row'<'col-sm-12 col-md-5'i><'col-sm-12 col-md-7'p>>"
    },
    drawOn:function(){
        const htmlDom = table_begin.cms.checkAllHtml();
        table_begin.cms.tableObject.api().on('draw', function () {
            $('.js-dataTable_init tr:eq(0) th:eq(0)').html(htmlDom);
            table_begin.cms.object.removeClass('js-table-checkable-enabled');
            One.helpers(['table-tools-checkable','core-bootstrap-tooltip']);
        });
    },
    check_delete:function(){
        $("body").on("click",".ids_deletes",function () {
            const object = $(this);
            application.cms.authorization({
                title: application.cms.authorizationMessage.title,
                text: application.cms.authorizationMessage.text,
            }, function () {
                application.ajax.post(object.data("action"),{":ids[]"  : table_begin.cms.checkbox_select(), hideLoader: true},function (r) {table_begin.cms.drawReset(r)});
            });
        });
    },
    check_enables:function(){
        $("body").on("click",".ids_enables",function () {
            const object = $(this);
            application.ajax.post(object.data("action"),{
                    ":status" : 1,
                    ":field"  : object.data("field"),
                    ":ids[]"  : table_begin.cms.checkbox_select(),
            },function (r) {table_begin.cms.drawReset(r)});
        });
    },
    check_disables:function(){
        $("body").on("click",".ids_disables",function () {
            const object = $(this);
            application.cms.authorization({
                title : application.cms.authorizationStatus.title,
                text  : application.cms.authorizationStatus.text,
            },function () {
                application.ajax.post(object.data("action"),{
                    ":status"    : 0,
                    ":field"     : object.data("field"),
                    ":ids[]"     : table_begin.cms.checkbox_select(),
                    "hideLoader" : true,
                },function (r) {table_begin.cms.drawReset(r)});
            });
        });
    },
    begin:function () {
        table_begin.cms.tableObject = table_begin.cms.object.dataTable({
            lengthChange : false,
            autoWidth    : false,
            checkbox     : true,
            serverSide   : true,
            ordering     : false,
            processing   : true,
            pageLength   : builder_table.pageSize || application.cms.tablePageSize,
            language     : table_begin.cms.language,
            dom          : table_begin.cms.bindDom,
            ajax         : {
                url      : table_begin.cms.object.data("action"),
                type     : "post",
                dataSrc  : function(json){
                    json.draw            = json.data.draw;
                    json.recordsTotal    = json.data.recordsTotal;
                    json.recordsFiltered = json.data.recordsFiltered;
                    json.data            = json.data.data;
                    return json.data || [];
                }
            },
            order        : builder_table.orderTable,
            columns      : builder_table.columnsTable,
            buttons      : builder_table.buttonsTable
        });
    },
}
$(document).ready(function() {
    table_begin.begin();
    table_begin.drawOn();
    table_begin.check_enables();
    table_begin.check_disables();
    table_begin.check_delete();
})