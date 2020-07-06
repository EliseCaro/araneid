const admin_index = {
    option :{
        chartOptions:{
            tooltips: {
                callbacks: {
                    label: function(tooltipItems, data) {
                        return data.datasets[tooltipItems.datasetIndex].label + " : " + tooltipItems.yLabel + " KB";
                    }
                },
                intersect: false,mode:'index',titleAlign:"center",
                titleFontSize:15,bodyAlign:"left",bodySpacing:7,xPadding:10,yPadding:10,
            },
            layout: {
                padding: {left: 30,right: 40,top: 10,bottom: 20}
            }
        },
        networkObject:{},
        chartNetworkCon:$(".js-chart-dashboard-network"),
        chartNetworkData:{labels: [],datasets: []},
    },
    begin  :{
        initialized:function () {
            Chart.defaults.scale.ticks.beginAtZero = true;
            Chart.defaults.global.legend.labels.boxWidth = 12;
        },
        initializedNetwork:function(){
           admin_index.begin.initializedAjax(function (result) {
               const network = new Chart(admin_index.option.chartNetworkCon,{
                   type: 'line',
                   data: admin_index.option.chartNetworkData,
                   options: admin_index.option.chartOptions,
               })
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
                }else {
                    application.ajax.requestBack(result.message,result.status,result.url);
                }
            })
        }
    }
}
$(document).ready(function() {
    admin_index.begin.initialized()
    admin_index.begin.initializedNetwork()
    setInterval(function () {
        admin_index.begin.initializedAjax(function (res) {
            console.log(res)
        },0)
    },2000);
});