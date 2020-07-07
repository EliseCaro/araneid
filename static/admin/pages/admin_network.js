const admin_network = {
        update:function (message,item) {
            const me = this
            if(me.object.data.labels.length >= 11){
                me.object.data.labels.shift()
                me.object.data.datasets.forEach((dataset,index) => {
                    dataset.data.shift();
                });me.object.update(0);
            }
            me.object.data.labels.push(message);
            me.object.data.datasets.forEach((d,i) => {
                if(i === 0 ){
                    const value = (item.bytesRecv/1024).toFixed(2)
                    d.data.push(value);
                    $(".network_up").html(value)
                }else {
                    const value = (item.bytesSent/1024).toFixed(2)
                    d.data.push(value);
                    $(".network_dow").html(value)
                }
            });
            me.object.update();
        },
        initialized:function (result) {
            const me = this
            me.chartData.labels.push(result.message)
            me.object = new Chart(me.chartCon,{
                type: 'line',
                data: me.chartData,
                options:me.options,
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
                      return data.datasets[tooltipItems.datasetIndex].label + " : " + tooltipItems.yLabel + " KB";
                  }
              },
          },
          layout: {
              padding: {top: 10, left: 30, right: 40, bottom: 20}
          },
      },
      object:{},
      chartCon:$(".js-chart-dashboard-network"),
      chartData:{
        labels: [],
        datasets:[
            {
                label: '上行',
                fill: true,
                backgroundColor: 'rgba(132, 94, 247, .3)',
                pointHoverBorderColor: 'rgba(132, 94, 247, 1)',
                data: [0]
            },
            {
                label: '下行',
                fill: true,
                backgroundColor: 'rgba(233, 236, 239, 1)',
                pointHoverBorderColor: 'rgba(233, 236, 239, 1)',
                data:[0]
            }
        ]
      }
}