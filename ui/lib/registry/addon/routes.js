import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('plugins', function () {
    this.route('plugin', { path: '/:plugin_name' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
});
