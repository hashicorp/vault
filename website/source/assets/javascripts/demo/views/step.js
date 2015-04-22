Demo.DemoStepView = Ember.View.extend({
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
    // Remove overlay
    $('.sidebar-overlay').removeClass('active');

    var element  = this.$();

    element.fadeOut(400);

    // Allow scrolling
    $('body').unbind('DOMMouseScroll mousewheel');
  },

  // click: function() {
  //   var element  = this.$();

  //   // Record scoll position
  //   var x = element.scrollX, y = element.scrollY;
  //   // Focus
  //   element.find('input.shell')[0].focus();
  //   // Scroll back to where you were
  //   element.scrollTo(x, y);
  // },

  keyDown: function(ev) {
    var cursor = this.get('controller.cursor'),
        currentLength = this.get('controller.commandLog.length');

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
          this.set('controller.logs', '');
          this.set('controller.notCleared', false);
        }
        break;

      // escape
      case 27:
        this.get('controller').transitionTo('index');
        break;
    }
  },

  deFocus: function() {
    var element = this.$().find('input.shell');

    // defocus while loading
    if (this.get('controller.isLoading')) {
      element.blur()
    }

  }.observes('controller.isLoading'),

  focus: function() {
    var element = this.$().find('input.shell');
    element.focus()
  }.observes('controller.cursor'),

  submitted: function() {
    var element  = this.$();

    console.log("submitted");

    // Focus the input
    element.find('input.shell')[0].focus();

    // Scroll to the bottom of the element
    element.scrollTop(element[0].scrollHeight);

  }.observes('controller.logs.length')
});
