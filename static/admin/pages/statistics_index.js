const statistics_index = {
    initialized:function () {
        Chart.defaults.scale.ticks.beginAtZero         = true;
        Chart.defaults.global.legend.labels.boxWidth   = 12;
        Chart.defaults.global.elements.point.radius    = 2;
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
        meObject.dataObject.labels = week.date
        meObject.dataObject.datasets = week.items
        meObject.options.tooltips.callbacks = {
            label: function(tooltipItems, data) {
                return data.datasets[tooltipItems.datasetIndex].label + " : " + tooltipItems.yLabel + " 个";
            }
        };
        meObject.journalObject = new Chart(meObject.journalCon,{
            type: 'line',data: meObject.dataObject,options:meObject.options,
        })
    },
}
$(document).ready(function() {
    statistics_index.initialized()
    statistics_index.initializedJournal()
});