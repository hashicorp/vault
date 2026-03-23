/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import EmberRouter from '@ember/routing/router';
import config from 'vault/config/environment';
import { addDocfyRoutes } from '@docfy/ember';
export default class Router extends EmberRouter {
  location = config.locationType;
  rootURL = config.rootURL;
}

Router.map(function () {
  this.route('vault', { path: '/' }, function () {
    this.route('cluster', { path: '/:cluster_name' }, function () {
      this.route('dashboard');
      this.mount('config-ui');
      this.mount('sync');
      this.route('recovery', function () {
        this.route('snapshots', function () {
          this.route('load');
          this.route('snapshot', { path: '/:snapshot_id' }, function () {
            this.route('manage');
            this.route('details');
          });
        });
      });
      this.route('oidc-provider-ns', { path: '/*namespace/identity/oidc/provider/:provider_name/authorize' });
      this.route('oidc-provider', { path: '/identity/oidc/provider/:provider_name/authorize' });
      this.route('oidc-callback', { path: '/auth/*auth_path/oidc/callback' });
      this.route('auth');
      this.route('redirect');
      this.route('init');
      this.route('logout');
      this.route('license');
      this.route('mfa-setup');
      this.route('clients', function () {
        this.route('counts', function () {
          this.route('overview');
          this.route('client-list');
        });
        this.route('config');
        this.route('edit');
      });
      this.route('usage-reporting');
      this.route('storage', { path: '/storage/raft' });
      this.route('storage-restore', { path: '/storage/raft/restore' });
      this.route('settings', function () {
        this.route('index', { path: '/' });
        this.route('seal');
        this.route('auth', function () {
          this.route('index', { path: '/' });
          this.route('enable');
          this.route('configure', { path: '/configure/:method' }, function () {
            this.route('index', { path: '/' });
            this.route('section', { path: '/:section_name' });
          });
        });
      });
      this.route('unseal');
      this.route('tools', function () {
        this.route('tool', { path: '/:selected_action' });
        this.mount('open-api-explorer', { path: '/api-explorer' });
      });
      this.route('access', function () {
        this.route('reset-password');
        this.route('methods', { path: '/' });
        this.route('method', { path: '/:path' }, function () {
          this.route('index', { path: '/' });
          this.route('item', { path: '/item/:item_type' }, function () {
            this.route('list', { path: '/' });
            this.route('create');
            this.route('edit', { path: '/edit/:item_id' });
            this.route('show', { path: '/show/:item_id' });
          });
          this.route('section', { path: '/:section_name' });
        });
        this.route('mfa', function () {
          this.route('index', { path: '/' });
          this.route('methods', function () {
            this.route('index', { path: '/' });
            this.route('create');
            this.route('method', { path: '/:id' }, function () {
              this.route('edit');
              this.route('enforcements');
            });
          });
          this.route('enforcements', function () {
            this.route('index', { path: '/' });
            this.route('create');
            this.route('enforcement', { path: '/:name' }, function () {
              this.route('edit');
            });
          });
        });
        this.route('leases', function () {
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
        // the "identity" routes expect :item_type to be plural
        this.route('identity', { path: '/identity/:item_type' }, function () {
          this.route('index', { path: '/' });
          this.route('create');
          this.route('merge');
          this.route('edit', { path: '/edit/:item_id' });
          this.route('show', { path: '/:item_id/:section' });
          this.route('aliases', function () {
            this.route('index', { path: '/' });
            this.route('add', { path: '/add/:item_id' });
            this.route('edit', { path: '/edit/:item_alias_id' });
            this.route('show', { path: '/:item_alias_id/:section' });
          });
        });
        this.route('control-groups');
        this.route('control-groups-configure', { path: '/control-groups/configure' });
        this.route('control-group-accessor', { path: '/control-groups/:accessor' });
        this.route('namespaces', function () {
          this.route('index', { path: '/' });
          this.route('create');
        });
        this.route('oidc', function () {
          this.route('clients', function () {
            this.route('create');
            this.route('client', { path: '/:name' }, function () {
              this.route('details');
              this.route('providers');
              this.route('edit');
            });
          });
          this.route('keys', function () {
            this.route('create');
            this.route('key', { path: '/:name' }, function () {
              this.route('details');
              this.route('clients');
              this.route('edit');
            });
          });
          this.route('assignments', function () {
            this.route('create');
            this.route('assignment', { path: '/:name' }, function () {
              this.route('details');
              this.route('edit');
            });
          });
          this.route('providers', function () {
            this.route('create');
            this.route('provider', { path: '/:name' }, function () {
              this.route('details');
              this.route('clients');
              this.route('edit');
            });
          });
          this.route('scopes', function () {
            this.route('create');
            this.route('scope', { path: '/:name' }, function () {
              this.route('details');
              this.route('edit');
            });
          });
        });
      });
      this.route('secrets-redirect', { path: '/secrets' }); // legacy redirect
      this.route('secrets', { path: '/secrets-engines' }, function () {
        this.route('enable', function () {
          this.route('create', { path: '/:mount_type' });
        });
        this.route('backends', { path: '/' });
        this.route('backend', { path: '/:backend' }, function () {
          this.mount('kmip');
          this.mount('kubernetes');
          this.mount('kv');
          this.mount('ldap');
          this.mount('pki');
          this.route('index', { path: '/' });
          this.route('configuration', function () {
            this.route('general-settings');
            this.route('plugin-settings');
            // only CONFIGURABLE_SECRET_ENGINES can be configured and access the edit route
            this.route('edit');
          });
          // because globs / params can't be empty,
          // we have to special-case ids of '' with their own routes
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

          // ssh sign
          this.route('sign-root', { path: '/sign/' });
          this.route('sign', { path: '/sign/*secret' });
          // transit-specific routes
          this.route('actions-root', { path: '/actions/' });
          this.route('actions', { path: '/actions/*secret' });
          // database specific route
          this.route('overview');
        });
      });
      this.route('policies', { path: '/policies/:type' }, function () {
        this.route('index', { path: '/' });
        this.route('create');
      });
      this.route('policy', { path: '/policy/:type' }, function () {
        this.route('show', { path: '/:policy_name' });
        this.route('edit', { path: '/:policy_name/edit' });
      });
      this.route('resilience-recovery');
      this.route('replication-dr-promote', function () {
        this.route('details');
      });
      this.mount('replication');
      // Catch-all for unmatched routes within the cluster (e.g., /vault/fake-route)
      // renders template: vault/cluster/not-found.hbs
      this.route('not-found', { path: '/*path' });
    });
  });
  if (config.environment !== 'production') {
    addDocfyRoutes(this);
  }
});
