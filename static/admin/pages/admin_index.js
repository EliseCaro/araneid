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
        initializedLoop:function () {
            setInterval(function () {
                admin_index.begin.initializedAjax(function (res) {
                    admin_memory.update(res.message,res.data.Memory)
                    admin_network.update(res.message,res.data.Network)
                },0)
            },1000);
        }
    }
}
$(document).ready(function() {
    admin_index.begin.initialized()
    admin_index.begin.initializedLoop()
});