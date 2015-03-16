Demo.Router.map(function() {
  this.route('demo', { path: '/demo' }, function() {
    this.route('crud', { path: '/crud' });
  });
});
