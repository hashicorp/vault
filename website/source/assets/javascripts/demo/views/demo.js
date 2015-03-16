Demo.DemoView = Ember.View.extend({
  classNames: ['demo-overlay'],

  didInsertElement: function() {
    var element  = this.$();

    element.hide().fadeIn(300);

    // Scroll to the bottom of the element
    element.scrollTop(element[0].scrollHeight);

    // Focus
    element.find('input.shell')[0].focus();

    // Hijack scrolling to only work within terminal
    //
    $(element).on('DOMMouseScroll mousewheel', function(ev) {
      var scrolledEl = $(this),
          scrollTop = this.scrollTop,
          scrollHeight = this.scrollHeight,
          height = scrolledEl.height(),
          delta = (ev.type == 'DOMMouseScroll' ?
              ev.originalEvent.detail * -40 :
              ev.originalEvent.wheelDelta),
          up = delta > 0;

      var prevent = function() {
          ev.stopPropagation();
          ev.preventDefault();
          ev.returnValue = false;
          return false;
      };

      if (!up && -delta > scrollHeight - height - scrollTop) {
          // Scrolling down, but this will take us past the bottom.
          scrolledEl.scrollTop(scrollHeight);
          return prevent();
      } else if (up && delta > scrollTop) {
          // Scrolling up, but this will take us past the top.
          scrolledEl.scrollTop(0);
          return prevent();
      }
    });
  },

  willDestroyElement: function() {
    var element  = this.$();

    element.fadeOut(400);

    // Allow scrolling
    $('body').unbind('DOMMouseScroll mousewheel');
  },

  click: function() {
    var element  = this.$();
    // Focus
    element.find('input.shell')[0].focus();
  },

  keyDown: function(ev) {
    var cursor = this.get('controller.cursor'),
        currentLength = this.get('controller.currentLog.length');

    console.log(ev);

    switch(ev.keyCode) {
      // Down arrow
      case 40:
        if (cursor === 0) {
            return;
        }

        this.incrementProperty('controller.cursor');
        break;

      // Up arrow
      case 38:
        if ((currentLength + cursor) === 0) {
            return;
        }

        this.decrementProperty('controller.cursor');
        break;

      // command + k
      case 75:
        if (ev.metaKey) {
          this.set('controller.currentLog', []);
          this.set('controller.notCleared', false);
        }
        break;

      // escape
      case 27:
        this.get('controller').transitionTo('index');
        break;
    }
  },

  submitted: function() {
    var element  = this.$();

    // Focus the input
    element.find('input.shell')[0].focus();

    // Scroll to the bottom of the element
    element.scrollTop(element[0].scrollHeight);

  }.observes('controller.currentLog')
});
