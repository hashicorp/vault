import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('overview');
  this.route('configuration'); // ARG TODO return to "Edit Configuration" depending on convo with Ivana
  this.route('configure', function () {
    // configure the engines tidy, crl, or url
    this.route('tidy'); // ARG TODO appear only to "tidy" the action. There is automatic-tidy, but that doesn't appear to be accounted for in the designs.
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
