import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('overview');
  this.route('configuration');
  this.route('configure', function () {
    // configure the engines tidy, crl, or url
    this.route('create');
    this.route('tidy', function () {
      this.route('details');
      this.route('edit');
    });
    this.route('crl', function () {
      this.route('details');
      this.route('edit');
    });
    this.route('url', function () {
      this.route('details');
      this.route('edit');
    });
    // generate root, ca or intermediate
    this.route('generate');
    this.route('root', function () {
      this.route('index', { path: '/' });
      this.route('details');
      this.route('edit');
    });
    this.route('import', function () {
      this.route('index', { path: '/' });
      this.route('details');
    });
    this.route('intermediate', function () {
      this.route('index', { path: '/' });
      this.route('details');
    });
  });
  this.route('roles', function () {
    this.route('index', { path: '/' });
    this.route('create');
    this.route('role', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('issuers', function () {
    this.route('index', { path: '/' });
    this.route('create');
    this.route('issuer', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('certificates', function () {
    this.route('index', { path: '/' });
    this.route('create');
    this.route('certificate', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('keys', function () {
    this.route('index', { path: '/' });
    this.route('create');
    this.route('key', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
});
