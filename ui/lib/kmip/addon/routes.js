import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function() {
  // Define your engine's route map here
  this.route('index', { path: '/' });
  this.route('scopes', { path: '/scopes' }, function() {
    this.route('scope', { path: '/:scope_name' }, function() {
      this.route('roles', { path: '/roles' }, function() {
        this.route('role', { path: '/:role_name' }, function() {
          this.route('credentials', { path: '/credentials' }, function() {
            this.route('cred', { path: '/:serial' });
          });
        });
      });
    });
  });
});
