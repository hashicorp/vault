Demo.DemoController = Ember.ObjectController.extend({
  isLoading: false,
  logs: "",

  appendLog: function(data, prefix) {
    var newline;

    if (prefix) {
      data = '$ ' + data;
    } else {
      newline = '';
    }

    newline = '\n';

    this.set('logs', this.get('logs')+data+newline);

    Ember.run.later(function() {
      var element = $('.demo-terminal');
      // Scroll to the bottom of the element
      element.scrollTop(element[0].scrollHeight);

      element.find('input.shell')[0].focus();
    }, 5);
  },
});
