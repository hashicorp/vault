import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('overview');
  this.route('configuration', function () {
    this.route('index', { path: '/' });
    this.route('tidy');
    this.route('create', function () {
      this.route('index', { path: '/' });
      this.route('import-ca');
      this.route('generate-root');
      this.route('generate-csr');
    });
    this.route('edit');
    this.route('details');
  });
  this.route('roles', function () {
    this.route('index', { path: '/' }); // ARG TODO remove
    this.route('create');
    this.route('role', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('issuers', function () {
    this.route('create');
    this.route('issuer', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('certificates', function () {
    this.route('create');
    this.route('certificate', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('keys', function () {
    this.route('create');
    this.route('key', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
});
