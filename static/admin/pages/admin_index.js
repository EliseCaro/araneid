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
            },
        },
        chartMemoryOptions:{
            tooltips: {
                callbacks: {
                    label: function(tooltipItems, data) {
                        return data.datasets[tooltipItems.datasetIndex].label + " : " + tooltipItems.yLabel + " %";
                    }
                },
                intersect: false,mode:'index',titleAlign:"center",
                titleFontSize:15,bodyAlign:"left",bodySpacing:7,xPadding:10,yPadding:10,
            },
            layout: {
                padding: {left: 30,right: 40,top: 10,bottom: 20}
            },
            scales: {
                yAxes: [{
                    ticks: {
                        suggestedMax: 100
                    }
                }]
            },
        },
        networkObject:{},
        chartNetworkCon:$(".js-chart-dashboard-network"),
        chartNetworkData:{
            labels: [],
            datasets:[
                {label: '上行',fill: true,backgroundColor: 'rgba(132, 94, 247, .3)',pointHoverBorderColor: 'rgba(132, 94, 247, 1)',data: [0]},
                {label: '下行',fill: true,backgroundColor: 'rgba(233, 236, 239, 1)',pointHoverBorderColor: 'rgba(233, 236, 239, 1)',data:[0]}
            ]
        },

        memoryObject:{},
        chartMemoryCon:$(".js-chart-dashboard-memory"),
        chartMemoryData:{
            labels: [],
            datasets:[
                {
                    label: '使用百分比',
                    fill: true,
                    backgroundColor: 'rgba(34, 184, 207, .3)',
                    pointHoverBorderColor: 'rgba(34, 184, 207, 1)',
                    data: [0]
                },
            ]
        },
    },
    begin  :{
        initialized:function () {
            Chart.defaults.scale.ticks.beginAtZero = true;
            Chart.defaults.global.legend.labels.boxWidth = 12;
            Chart.defaults.scale.gridLines.color                = 'transparent'; // 网格线的颜色
            Chart.defaults.scale.gridLines.zeroLineColor        = 'transparent'; // 第一条网格线的颜色
            Chart.defaults.global.elements.point.radius         = 0; // 每个顶点上的圆角
        },
        initializedNetwork:function(){
           admin_index.begin.initializedAjax(function (result) {
               admin_index.option.chartNetworkData.labels.push(result.message)
               admin_index.option.networkObject = new Chart(admin_index.option.chartNetworkCon,{
                   type: 'line',
                   data: admin_index.option.chartNetworkData,
                   options: admin_index.option.chartOptions,
               })
               admin_index.option.chartMemoryData.labels.push(result.message)
               admin_index.option.memoryObject = new Chart(admin_index.option.chartMemoryCon,{
                   type: 'line',
                   data: admin_index.option.chartMemoryData,
                   options: admin_index.option.chartMemoryOptions,
               })
           },1)
        },
        updateNetwork:function(message,item){
            if(admin_index.option.networkObject.data.labels.length >= 11){
                admin_index.option.networkObject.data.labels.shift()
                admin_index.option.networkObject.data.datasets.forEach((dataset,index) => {
                    dataset.data.shift();
                });
                admin_index.option.networkObject.update(0);
            }
            admin_index.option.networkObject.data.labels.push(message);
            admin_index.option.networkObject.data.datasets.forEach((dataset,index) => {
                if(index === 0 ){
                    const value = (item.bytesRecv/1024).toFixed(2)
                    $(".network_up").html(value)
                    dataset.data.push(value);
                }else {
                    const value = (item.bytesSent/1024).toFixed(2)
                    dataset.data.push(value);
                    $(".network_dow").html(value)
                }
            });
            admin_index.option.networkObject.update();
        },
        updateMemory:function(message,item){
            const object = admin_index.option.memoryObject;
            if(object.data.labels.length >= 11){
                object.data.labels.shift()
                object.data.datasets.forEach((dataset) => {
                    dataset.data.shift();
                });
                object.update(0);
            }
            object.data.labels.push(message);
            object.data.datasets.forEach((d,i) => {
                $(".memory_total").html((item.total/1024/1024/1024).toFixed(2))
                $(".memory_used").html((item.used/1024/1024/1024).toFixed(2))
                d.data.push(item.usedPercent.toFixed(2));
            });
            object.update();
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
                    const data = res.data
                    if (data && data.Network) {
                        admin_index.begin.updateNetwork(res.message,data.Network)
                    }
                    if (data && data.Memory) {
                        admin_index.begin.updateMemory(res.message,data.Memory)
                    }
                },0)
            },2000);
        }
    }
}
$(document).ready(function() {
    admin_index.begin.initialized()
    admin_index.begin.initializedNetwork()
    admin_index.begin.initializedLoop()
});