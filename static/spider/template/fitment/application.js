$(document).ready(function() {
    new Swiper('.swiper-container', {slidesPerView: 1,spaceBetween: 0,loop: true,pagination: {el: '.swiper-pagination', clickable: true},
        navigation: {nextEl: '.swiper-button-next',prevEl: '.swiper-button-prev'},
    });
    $(function(){
        var search=$('.search');
        var menu=$('.menu');
        $('.navbtn').click(function(){
            if(search.is(':hidden') && menu.is(':hidden')){
                search.addClass('showEl');
                menu.addClass('showEl');
            }
        })
        $('.close').click(function(){
            search.removeClass('showEl');
            menu.removeClass('showEl');
        })
    })
})