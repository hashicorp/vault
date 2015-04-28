Demo.DemoStepView = Ember.View.extend({
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
      element.blur();
    }

  }.observes('controller.isLoading'),

  focus: function() {
    var element = this.$().find('input.shell');
    element.focus();
  }.observes('controller.cursor'),

  submitted: function() {
    var element  = this.$();

    // Focus the input
    element.find('input.shell')[0].focus();

    // guarantees that the log is scrolled when updated
    Ember.run.scheduleOnce('afterRender', this, function() {
      window.scrollTo(0, document.body.scrollHeight);
    });

  }.observes('controller.logs.length')
});
