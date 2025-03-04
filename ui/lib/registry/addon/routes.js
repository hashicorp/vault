import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('plugins', function () {
    this.route('plugin', { path: '/:id' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
});
