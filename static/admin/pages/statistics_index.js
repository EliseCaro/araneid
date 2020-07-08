const statistics_index = {
    initialized:function () {
        Chart.defaults.scale.ticks.beginAtZero         = true;
        Chart.defaults.global.legend.labels.boxWidth   = 12;
        Chart.defaults.scale.gridLines.color           = 'transparent'; // 网格线的颜色
        Chart.defaults.scale.gridLines.zeroLineColor   = 'transparent'; // 第一条网格线的颜色
        Chart.defaults.global.elements.point.radius    = 0;             // 每个顶点上的圆角
    },
    options:{
        layout:{padding:{top:10,left:10,right:28,bottom:10}},
        tooltips: {mode:'index',xPadding:10,yPadding:10,bodySpacing:7,bodyAlign:"left",intersect: false,titleFontSize:15,titleAlign:"center"},
    },
    dataObject:{labels: [],datasets:[{fill:true}]},

    journalObject:{},
    journalCon:$(".js-chart-dashboard-journal"),
    initializedJournal : function () {
        const meObject = this
        const labelsData = meObject.journalLabelsData();
        meObject.dataObject.labels = labelsData.label
        meObject.dataObject.datasets[0].label = "周蜘蛛量"
        meObject.dataObject.datasets[0].backgroundColor = 'rgba(27, 158, 183, 0.8)'
        meObject.dataObject.datasets[0].data  = labelsData.data
        meObject.options.tooltips.callbacks = {
            label: function(tooltipItems, data) {
                return data.datasets[tooltipItems.datasetIndex].label + " : " + tooltipItems.yLabel + " 个";
            }
        };
        meObject.journalObject = new Chart(meObject.journalCon,{
            type: 'line',data: meObject.dataObject,options:meObject.options,
        })
    },
    journalLabelsData:function () {
        let label = [],data = [];
        for (let i in week){
            label.push(week[i].date)
            data.push(week[i].count)
        }
        return {label:label,data:data}
    }
}
$(document).ready(function() {
    statistics_index.initialized()
    statistics_index.initializedJournal()
});