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
    this.route('index', { path: '/' });
    this.route('create');
    this.route('role', { path: '/:id' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('issuers', function () {
    this.route('index', { path: '/' });
    this.route('issuer', { path: '/:id' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('certificates', function () {
    this.route('index', { path: '/' });
    this.route('create');
    this.route('certificate', { path: '/:id' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('keys', function () {
    this.route('index', { path: '/' });
    this.route('generate');
    this.route('import');
    this.route('key', { path: '/:id' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
});
