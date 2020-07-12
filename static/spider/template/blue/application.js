$(document).ready(function() {
    jQuery(".focusBox").slide({
        titCell: ".num li",
        mainCell: ".pic",
        effect: "fold",
        autoPlay: true,
        trigger: "click",
        startFun: function (i) {
            jQuery(".focusBox .txt li").eq(i).animate({"bottom": 0}).siblings().animate({"bottom": -34});
        }
    });
})