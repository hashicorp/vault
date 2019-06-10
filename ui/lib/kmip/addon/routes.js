import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function() {
  this.route('configuration');
  this.route('configure');
  this.route('scopes', { path: '/scopes' }, function() {
    this.route('index', { path: '/' });
    this.route('create');
    this.route('scope', { path: '/:scope_name' }, function() {
      this.route('roles', { path: '/roles' }, function() {
        this.route('index', { path: '/' });
        this.route('create');
        this.route('role', { path: '/:role_name' }, function() {
          this.route('credentials', { path: '/credentials' }, function() {
            this.route('index', { path: '/' });
            this.route('cred', { path: '/:serial' });
          });
        });
      });
    });
  });
});
