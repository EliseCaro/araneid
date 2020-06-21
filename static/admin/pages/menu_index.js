const menu_index = {
    resultsSortNumber: 1,
    resultsSortMenus : [],
    begin:function () {
        $('#menu_list').nestable().nestable('collapseAll');
        $('.expand_all').on('click', function() {
            $('.dd').nestable('expandAll');
        });
        $('.collapse_all').on('click', function() {
            $('.dd').nestable('collapseAll');
        });
    },
    analysis:function(serialize,pid){
      for (let i in serialize){
          let resultsSortMenu = serialize[i];
          menu_index.resultsSortMenus.push({
              id    : resultsSortMenu.id,
              pid   : pid,
              sort  : menu_index.resultsSortNumber
          })
          ++ menu_index.resultsSortNumber;
          if(resultsSortMenu.children && resultsSortMenu.children.length > 0){
              menu_index.analysis(resultsSortMenu.children,resultsSortMenu.id)
          }
      }
      return menu_index.resultsSortMenus;
    },
    change:function () {
        $('#menu_list').nestable().on('change', function(){
            const object_dd    = $('.dd'),object_list  = $(this);
            menu_index.resultsSortNumber = 0;
            menu_index.resultsSortMenus = [];
            application.ajax.post(object_list.data("action"),{
                ":menus" : JSON.stringify(menu_index.analysis(object_dd.nestable('serialize'),object_list.data("root"))),
                ":root"  : object_list.data("root")
            },function (res) {
                application.cms.loader("hide");
                if(res && res.status === false){
                    application.cms.alert({
                        title : application.cms.requestStatus.errorTitle,
                        text  : res.message,
                        type  : "error",
                        timer : application.cms.requestStatus.layerTime
                    });
                    return false;
                }
            })
        });
    }
}
$(document).ready(function() {
    menu_index.begin();
    menu_index.change();
});