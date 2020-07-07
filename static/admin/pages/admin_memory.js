const admin_memory = {
    update:function(message,item){
        const me = this
        if(me.object.data.labels.length >= 11){
            me.object.data.labels.shift()
            me.object.data.datasets.forEach((dataset) => {
                dataset.data.shift();
            });
            me.object.update(0);
        }
        me.object.data.labels.push(message);
        me.object.data.datasets.forEach((d,i) => {
            $(".memory_total").html((item.total/1024/1024/1024).toFixed(2))
            $(".memory_used").html((item.used/1024/1024/1024).toFixed(2))
            d.data.push(item.usedPercent.toFixed(2));
        });
        me.object.update();
    },
    initialized:function (result) {
        const me = this
        me.chartData.labels.push(result.message)
        me.object = new Chart(me.chartCon,{
            type: 'line',
            data: me.chartData,
            options: me.options,
        })
    },
    options:{
        tooltips: {
            mode:'index',
            xPadding:10,
            yPadding:10,
            bodySpacing:7,
            bodyAlign:"left",
            intersect: false,
            titleFontSize:15,
            titleAlign:"center",
            callbacks: {
                label: function(tooltipItems, data) {
                    return data.datasets[tooltipItems.datasetIndex].label + " : " + tooltipItems.yLabel + " %";
                }
            },
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
    object:{},
    chartCon:$(".js-chart-dashboard-memory"),
    chartData:{
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
}