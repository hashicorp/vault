Demo.DemoController = Ember.ObjectController.extend({
  isLoading: false,
  logs: "",

  appendLog: function(data, prefix) {
    var newline;

    if (prefix) {
      data = '$ ' + data;
    }

    if (this.get('logs.length') === 0) {
      newline = '';
    } else {
      newline = '\n';
    }

    this.set('logs', this.get('logs')+newline+data);

    Ember.run.later(function() {
      var element = $('.demo-overlay');
      // Scroll to the bottom of the element
      element.scrollTop(element[0].scrollHeight);

      element.find('input.shell')[0].focus();
    }, 5);
  },
});
