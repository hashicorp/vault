import EmberRouter from '@ember/routing/router';
import config from './config/environment';

const Router = EmberRouter.extend({
  location: config.locationType,
  rootURL: config.rootURL,
});

Router.map(function() {
  this.route('vault', { path: '/' }, function() {
    this.route('cluster', { path: '/:cluster_name' }, function() {
      this.route('oidc-callback', { path: '/auth/*auth_path/oidc/callback' });
      this.route('auth');
      this.route('init');
      this.route('logout');
      this.route('license');
      this.route('settings', function() {
        this.route('index', { path: '/' });
        this.route('seal');
        this.route('auth', function() {
          this.route('index', { path: '/' });
          this.route('enable');
          this.route('configure', { path: '/configure/:method' }, function() {
            this.route('index', { path: '/' });
            this.route('section', { path: '/:section_name' });
          });
        });
        this.route('mount-secret-backend');
        this.route('configure-secret-backend', { path: '/secrets/configure/:backend' }, function() {
          this.route('index', { path: '/' });
          this.route('section', { path: '/:section_name' });
        });
      });
      this.route('unseal');
      this.route('tools', function() {
        this.route('tool', { path: '/:selected_action' });
      });
      this.route('access', function() {
        this.route('methods', { path: '/' });
        this.route('method', { path: '/:path' }, function() {
          this.route('index', { path: '/' });
          this.route('section', { path: '/:section_name' });
        });
        this.route('leases', function() {
          // lookup
          this.route('index', { path: '/' });
          // lookup prefix
          // revoke prefix + revoke force
          this.route('list-root', { path: '/list/' });
          this.route('list', { path: '/list/*prefix' });
          //renew + revoke
          this.route('show', { path: '/show/*lease_id' });
        });
        // the outer identity route handles group and entity items
        this.route('identity', { path: '/identity/:item_type' }, function() {
          this.route('index', { path: '/' });
          this.route('create');
          this.route('merge');
          this.route('edit', { path: '/edit/:item_id' });
          this.route('show', { path: '/:item_id/:section' });
          this.route('aliases', function() {
            this.route('index', { path: '/' });
            this.route('add', { path: '/add/:item_id' });
            this.route('edit', { path: '/edit/:item_alias_id' });
            this.route('show', { path: '/:item_alias_id/:section' });
          });
        });
        this.route('control-groups');
        this.route('control-groups-configure', { path: '/control-groups/configure' });
        this.route('control-group-accessor', { path: '/control-groups/:accessor' });
        this.route('namespaces', function() {
          this.route('index', { path: '/' });
          this.route('create');
        });
      });
      this.route('secrets', function() {
        this.route('backends', { path: '/' });
        this.route('backend', { path: '/:backend' }, function() {
          this.route('index', { path: '/' });
          this.route('configuration');
          // because globs / params can't be empty,
          // we have to special-case ids of '' with thier own routes
          this.route('list-root', { path: '/list/' });
          this.route('create-root', { path: '/create/' });
          this.route('show-root', { path: '/show/' });
          this.route('edit-root', { path: '/edit/' });

          this.route('list', { path: '/list/*secret' });
          this.route('show', { path: '/show/*secret' });
          this.route('create', { path: '/create/*secret' });
          this.route('edit', { path: '/edit/*secret' });

          this.route('credentials-root', { path: '/credentials/' });
          this.route('credentials', { path: '/credentials/*secret' });

          // kv v2 versions
          this.route('versions-root', { path: '/versions/' });
          this.route('versions', { path: '/versions/*secret' });

          // ssh sign
          this.route('sign-root', { path: '/sign/' });
          this.route('sign', { path: '/sign/*secret' });
          // transit-specific routes
          this.route('actions-root', { path: '/actions/' });
          this.route('actions', { path: '/actions/*secret' });
        });
      });
      this.route('policies', { path: '/policies/:type' }, function() {
        this.route('index', { path: '/' });
        this.route('create', { path: '/create' });
      });
      this.route('policy', { path: '/policy/:type' }, function() {
        this.route('show', { path: '/:policy_name' });
        this.route('edit', { path: '/:policy_name/edit' });
      });
      this.route('replication-dr-promote');
      this.route('replication', function() {
        this.route('index', { path: '/' });
        this.route('mode', { path: '/:replication_mode' }, function() {
          //details
          this.route('index', { path: '/' });
          this.route('manage');
          this.route('secondaries', function() {
            this.route('add', { path: '/add' });
            this.route('revoke', { path: '/revoke' });
            this.route('config-show', { path: '/config/show/:secondary_id' });
            this.route('config-edit', { path: '/config/edit/:secondary_id' });
            this.route('config-create', { path: '/config/create/:secondary_id' });
          });
        });
      });

      this.route('not-found', { path: '/*path' });
    });
    this.route('not-found', { path: '/*path' });
  });
});

export default Router;
