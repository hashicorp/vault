Demo.Router.map(function() {
  this.route('demo', function() {
    this.route('step', { path: '/:id' });
  });
});
