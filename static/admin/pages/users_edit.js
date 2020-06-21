const users_edit = {
    cms:{},
    begin:{
        wizard:function () {
            $('.js-wizard-simple').bootstrapWizard({
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
    }
}
$(document).ready(function() {
    users_edit.begin.wizard();
});