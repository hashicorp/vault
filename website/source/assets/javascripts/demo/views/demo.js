Demo.DemoView = Ember.View.extend({
  classNames: ['demo-overlay'],

  mouseUp: function(ev) {
    var selection = window.getSelection().toString();

    if (selection.length > 0) {
      // Ignore clicks when they are trying to select something
      return;
    }

    var element = this.$();

    // Record scroll position
    var x = element.scrollX, y = element.scrollY;
    // Focus
    element.find('input.shell')[0].focus();
    // Scroll back  to where you were
    element.scrollTop(x, y);
  },

  didInsertElement: function() {
    var controller = this.get('controller'),
        overlay    = $('.sidebar-overlay'),
        element    = this.$();

    $('body').addClass('demo-active');

    overlay.addClass('active');

    overlay.on('click', function() {
      controller.transitionTo('index');
    });

    // Scroll to the bottom of the element
    element.scrollTop(element[0].scrollHeight);

    // Focus
    element.find('input.shell')[0].focus();
  },

  willDestroyElement: function() {
    // Remove overlay
    $('.sidebar-overlay').removeClass('active');

    var element  = this.$();

    element.fadeOut(400);

    $('body').removeClass('demo-active');

    // reset scroll to top after closing demo
    window.scrollTo(0, 0);
  },

});
