const roles_edit = {
     cms:{
         tree_object   : $('.js-tree'),
         wizard_object : $('.js-wizard-simple'),
         form_object   : $('.post_action'),
         tab_array     : {}, // 储存所有TABS
         auth_node     : [], // 储存选中结果
         parse_int     : function (index,data) {
             let resource = [];data.push(index);
             for (let i in data){
                 resource.push(parseInt(data[i]));
             }
             return resource;
         }
     },
     begin:{
         wizard:function () {
             roles_edit.cms.wizard_object.bootstrapWizard({
                 nextSelector     : '[data-wizard="next"]',
                 previousSelector : '[data-wizard="prev"]',
                 finishSelector   : '[data-wizard="finish"]',
                 onTabShow: (tab, nav, index) => {
                     let percent = ((index + 1) / nav.find('li').length) * 100;
                     let progress = nav.parents('.block').find('[data-wizard="progress"] > .progress-bar');
                     if (progress.length) {
                         progress.css({ width: percent + 1 + '%' });
                     }
                 }
             });
         },
         js_tree:function () {
             roles_edit.cms.tree_object.each(function () {
                 $(this).jstree({
                     plugins  : ["checkbox"],
                     checkbox : {
                         "keep_selected_style" : false,
                         "three_state" : false,
                         "cascade" : 'down+up'
                     }
                 });
                 roles_edit.cms.tab_array[$(this).data('tab')] = $(this);
             });
         },
         form_action(){
             $(".post_action").submit(function () {
                 let form_post = $(this).serialize();
                 $.each(roles_edit.cms.tab_array, function (index) {
                     let structure = $(this).jstree(true).get_checked();
                     if(structure.length){
                         let res = roles_edit.cms.parse_int(index,structure);
                         roles_edit.cms.auth_node.push.apply(roles_edit.cms.auth_node, res);
                     }
                 })
                 if (roles_edit.cms.auth_node.length) {
                     form_post += '&menu_auth=' + JSON.stringify(roles_edit.cms.auth_node);
                 }
                 application.ajax.post($(this).attr("action"),form_post);
                 return false;
             });
         }
     }
}
$(document).ready(function() {
     roles_edit.begin.wizard();
     roles_edit.begin.js_tree();
     roles_edit.begin.form_action();
});