import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function() {
  this.route('configuration');
  this.route('configure');
  this.route('scopes', { path: '/scopes' }, function() {
    this.route('index', { path: '/' });
    this.route('create');
  });
  this.route('scope', { path: '/scopes/:scope_name/roles' }, function() {
    this.route('roles', { path: '/' });
    this.route('roles.create', { path: '/create' });
  });
  this.route('role', { path: '/scopes/:scope_name/roles/:role_name' });
  this.route('role.edit', { path: '/scopes/:scope_name/roles/:role_name/edit' });
  this.route('credentials', { path: '/scopes/:scope_name/roles/:role_name/credentials' }, function() {
    this.route('index', { path: '/' });
    this.route('generate');
    this.route('show', { path: '/:serial' });
  });
});
