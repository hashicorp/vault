Demo.DemoView = Ember.View.extend({
  classNames: ['demo-overlay'],

  mouseUp: function(ev) {
    var selection = window.getSelection().toString();

    if (selection.length > 0) {
      // Ignore clicks when they are trying to select something
      return;
    }

    var element = this.$();

    // Record scoll position
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

    overlay.addClass('active');

    overlay.on('click', function() {
      controller.transitionTo('index');
    });

    element.hide().fadeIn(300);

    // Scroll to the bottom of the element
    element.scrollTop(element[0].scrollHeight);

    // Focus
    element.find('input.shell')[0].focus();

    // Hijack scrolling to only work within terminal
    //
    element.on('DOMMouseScroll mousewheel', function(e) {
        e.preventDefault();
    });

    $('.demo-terminal').on('DOMMouseScroll mousewheel', function(e) {
      var scrollTo = null;

      if (e.type == 'mousewheel') {
        scrollTo = (e.originalEvent.wheelDelta * -1);
      } else if (e.type == 'DOMMouseScroll') {
        scrollTo = 40 * e.originalEvent.detail;
      }

      if (scrollTo) {
        e.preventDefault();
        $(this).scrollTop(scrollTo + $(this).scrollTop());
      }
    });
  },

  willDestroyElement: function() {
    // Remove overlay
    $('.sidebar-overlay').removeClass('active');

    var element  = this.$();

    element.fadeOut(400);

    // Allow scrolling
    $('body').unbind('DOMMouseScroll mousewheel');
  },

});
