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
    autoDataObject:{labels:[],datasets:[]},
    journalDataObject:{labels:[],datasets:[]},
    autoCon:$(".js-chart-dashboard-auto"),
    initializedAuto : function () {
        const meObject = this
        meObject.autoDataObject.labels = autoWeek.date
        meObject.autoDataObject.datasets = autoWeek.items
        meObject.options.tooltips.callbacks = {
            label: function(tooltipItems, data) {
                return data.datasets[tooltipItems.datasetIndex].label + " : " + tooltipItems.yLabel + " 条";
            }
        };
        new Chart(meObject.autoCon,{
            type: 'line',
            data: meObject.autoDataObject,
            options:meObject.options,
        })
    },
    journalCon:$(".js-chart-dashboard-journal"),
    initializedJournal : function () {
        const meObject = this
        meObject.journalDataObject.labels = week.date
        meObject.journalDataObject.datasets = week.items
        meObject.options.tooltips.callbacks = {
            label: function(tooltipItems, data) {
                return data.datasets[tooltipItems.datasetIndex].label + " : " + tooltipItems.yLabel + " 个";
            }
        };
        new Chart(meObject.journalCon,{
            type: 'line',
            data: meObject.journalDataObject,
            options:meObject.options,
        })
    },
}
$(document).ready(function() {
    statistics_index.initialized()
    statistics_index.initializedJournal()
    statistics_index.initializedAuto()
});