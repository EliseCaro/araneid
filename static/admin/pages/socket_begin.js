const socket_begin = {
    option :{
        socketObject : {},
        notifications:$(".js_dynamic_notifications"),
        count_box :$(".js_count_notifications")
    },
    begin:{
        init:function () {
            socket_begin.option.socketObject = new WebSocket(application.cms.socketHost)
            socket_begin.option.socketObject.onmessage = function (e) {
                socket_begin.begin.assemble(JSON.parse(e.data))
                if(window.location.pathname  === "/admin/collect/index"){
                    table_begin.cms.tableObject.api().draw(true);
                }
            }
        },
        assemble:function (data) {
            let count  = parseInt(socket_begin.option.count_box.html());
            let length = socket_begin.option.notifications.children("li").length;
            let className = data.id ? "open_iframe" : "";
            let html  = "<li class='animated bounce'>";
                html += '<a class="text-dark media py-2 '+ className +'" data-area="380px,240px" href="/admin/inform/check?:id='+ data.id +'&:popup=1">'
                html += '<div class="mr-2 ml-3"><i class="fa fa-fw fa-eye text-amethyst-dark"></i></div>'
                html += '<div class="media-body pr-2">'
                html += '<div class="font-weight-normal font-size-sm">'+ data.context +'...</div>'
                html += '<small class="text-muted">'+ data.string_time +'</small>';
                html += '</div></a></li>';
            if (length >= 5) {
              socket_begin.option.notifications.children("li")[socket_begin.option.notifications.children("li").length-1].remove();
            }
            socket_begin.option.notifications.prepend(html)
            socket_begin.option.count_box.html(count+1)
        }
    }
}
$(document).ready(function() {
    socket_begin.begin.init()
})