const admin_index = {
    begin  :{
        initialized:function () {
            Chart.defaults.scale.ticks.beginAtZero = true;
            Chart.defaults.global.legend.labels.boxWidth = 12;
            Chart.defaults.scale.gridLines.color                = 'transparent'; // 网格线的颜色
            Chart.defaults.scale.gridLines.zeroLineColor        = 'transparent'; // 第一条网格线的颜色
            Chart.defaults.global.elements.point.radius         = 0; // 每个顶点上的圆角
            admin_index.begin.initializedAjax(function (result) {
                admin_network.initialized(result)
                admin_memory.initialized(result)
            },1)
            One.helpers(['easy-pie-chart']);
        },
        initializedAjax:function (callback,init) {
            application.ajax.post($(".chart_box").data("action"),{
                hideLoader:true,
                init:init
            },function (result) {
                if (result.status === true) {
                    application.cms.loader("hide");
                    callback && callback(result)
                }
            })
        },
        updateCpu:function(message,item){
         const object = $(".cpu_chart_pie");
         const footer = object.parent().parent().children(".chart_footer").children(".text-muted")
         object.data('easyPieChart').update(item.used_percent);
         object.children(".chart_content").html(item.used_percent.toFixed(2) + " %")
         footer.children("span").html(item.cores)
        },
        updateLoad:function(message,item){
            const object = $(".load_chart_pie");
            const footer = object.parent().parent().children(".chart_footer").children(".text-muted")
            object.data('easyPieChart').update(item.Load1);
            object.children(".chart_content").html(item.load1.toFixed(2) + " %")
            footer.children("span").html(item.load1 >= 80 ? "资源耗尽" : "运行流畅")
            let html = "";
                html+= "最近1分钟平均负载 : "  + item.load1.toFixed(2) + "<br/>"
                html+= "最近5分钟平均负载 : "  + item.load5.toFixed(2) + "<br/>"
                html+= "最近15分钟平均负载 : " + item.load15.toFixed(2)
            footer.children("i").attr("data-content",html)
        },
        initializedLoop:function () {
            setInterval(function () {
                admin_index.begin.initializedAjax(function (res) {
                    admin_memory.update(res.message,res.data.Memory)
                    admin_network.update(res.message,res.data.Network)
                    admin_index.begin.updateCpu(res.message,res.data.CPU)
                    admin_index.begin.updateLoad(res.message,res.data.Load)
                },0)
            },1000);
        }
    }
}
$(document).ready(function() {
    admin_index.begin.initialized()
    admin_index.begin.initializedLoop()
});