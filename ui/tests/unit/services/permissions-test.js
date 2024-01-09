/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import Pretender from 'pretender';
import Service from '@ember/service';

const PERMISSIONS_RESPONSE = {
  data: {
    exact_paths: {
      foo: {
        capabilities: ['read'],
      },
      'bar/bee': {
        capabilities: ['create', 'list'],
      },
      boo: {
        capabilities: ['deny'],
      },
    },
    glob_paths: {
      'baz/biz': {
        capabilities: ['read'],
      },
      'ends/in/slash/': {
        capabilities: ['list'],
      },
    },
  },
};

module('Unit | Service | permissions', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.server = new Pretender();
    this.server.get('/v1/sys/internal/ui/resultant-acl', () => {
      return [200, { 'Content-Type': 'application/json' }, JSON.stringify(PERMISSIONS_RESPONSE)];
    });
    this.service = this.owner.lookup('service:permissions');
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('sets paths properly', async function (assert) {
    await this.service.getPaths.perform();
    assert.deepEqual(this.service.get('exactPaths'), PERMISSIONS_RESPONSE.data.exact_paths);
    assert.deepEqual(this.service.get('globPaths'), PERMISSIONS_RESPONSE.data.glob_paths);
  });

  test('sets the root token', function (assert) {
    this.service.setPaths({ data: { root: true } });
    assert.true(this.service.canViewAll);
  });

  test('defaults to show all items when policy cannot be found', async function (assert) {
    this.server.get('/v1/sys/internal/ui/resultant-acl', () => {
      return [403, { 'Content-Type': 'application/json' }];
    });
    await this.service.getPaths.perform();
    assert.true(this.service.canViewAll);
  });

  test('returns the first allowed nav route for policies', function (assert) {
    const policyPaths = {
      'sys/policies/acl': {
        capabilities: ['deny'],
      },
      'sys/policies/rgp': {
        capabilities: ['read'],
      },
    };
    this.service.set('exactPaths', policyPaths);
    assert.strictEqual(this.service.navPathParams('policies').models[0], 'rgp');
  });

  test('returns the first allowed nav route for access', function (assert) {
    const accessPaths = {
      'sys/auth': {
        capabilities: ['deny'],
      },
      'identity/entity/id': {
        capabilities: ['read'],
      },
    };
    const expected = { route: 'vault.cluster.access.identity', models: ['entities'] };
    this.service.set('exactPaths', accessPaths);
    assert.deepEqual(this.service.navPathParams('access'), expected);
  });

  module('hasPermission', function () {
    test('returns true if a policy includes access to an exact path', function (assert) {
      this.service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
      assert.true(this.service.hasPermission('foo'));
    });

    test('returns true if a paths prefix is included in the policys exact paths', function (assert) {
      this.service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
      assert.true(this.service.hasPermission('bar'));
    });

    test('it returns true if a policy includes access to a glob path', function (assert) {
      this.service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
      assert.true(this.service.hasPermission('baz/biz/hi'));
    });

    test('it returns true if a policy includes access to the * glob path', function (assert) {
      const splatPath = { '': {} };
      this.service.set('globPaths', splatPath);
      assert.true(this.service.hasPermission('hi'));
    });

    test('it returns false if the matched path includes the deny capability', function (assert) {
      this.service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
      assert.false(this.service.hasPermission('boo'));
    });

    test('it returns true if passed path does not end in a slash but globPath does', function (assert) {
      this.service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
      assert.true(this.service.hasPermission('ends/in/slash'), 'matches without slash');
      assert.true(this.service.hasPermission('ends/in/slash/'), 'matches with slash');
    });

    test('it returns false if a policy does not includes access to a path', function (assert) {
      assert.false(this.service.hasPermission('danger'));
    });
    test('returns true with the root token', function (assert) {
      this.service.set('canViewAll', true);
      assert.true(this.service.hasPermission('hi'));
    });
    test('it returns true if a policy has the specified capabilities on a path', function (assert) {
      this.service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
      this.service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
      assert.true(this.service.hasPermission('bar/bee', ['create', 'list']));
      assert.true(this.service.hasPermission('baz/biz', ['read']));
    });

    test('it returns false if a policy does not have the specified capabilities on a path', function (assert) {
      this.service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
      this.service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
      assert.false(this.service.hasPermission('bar/bee', ['create', 'delete']));
      assert.false(this.service.hasPermission('foo', ['create']));
    });
  });

  module('hasNavPermission', function () {
    test('returns true if a policy includes the required capabilities for at least one path', function (assert) {
      const accessPaths = {
        'sys/auth': {
          capabilities: ['deny'],
        },
        'identity/group/id': {
          capabilities: ['list', 'read'],
        },
      };
      this.service.set('exactPaths', accessPaths);
      assert.true(this.service.hasNavPermission('access', 'groups'));
    });

    test('returns false if a policy does not include the required capabilities for at least one path', function (assert) {
      const accessPaths = {
        'sys/auth': {
          capabilities: ['deny'],
        },
        'identity/group/id': {
          capabilities: ['read'],
        },
      };
      this.service.set('exactPaths', accessPaths);
      assert.false(this.service.hasNavPermission('access', 'groups'));
    });

    test('should handle routeParams as array', function (assert) {
      const getPaths = (override) => ({
        'sys/auth': {
          capabilities: [override || 'read'],
        },
        'identity/mfa/method': {
          capabilities: [override || 'read'],
        },
        'identity/oidc/client': {
          capabilities: [override || 'deny'],
        },
      });

      this.service.set('exactPaths', getPaths());
      assert.true(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc']),
        'hasNavPermission returns true for array of route params when any route is permitted'
      );
      assert.false(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc'], true),
        'hasNavPermission returns false for array of route params when any route is not permitted and requireAll is passed'
      );

      this.service.set('exactPaths', getPaths('read'));
      assert.true(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc'], true),
        'hasNavPermission returns true for array of route params when all routes are permitted and requireAll is passed'
      );

      this.service.set('exactPaths', getPaths('deny'));
      assert.false(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc']),
        'hasNavPermission returns false for array of route params when no routes are permitted'
      );
      assert.false(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc'], true),
        'hasNavPermission returns false for array of route params when no routes are permitted and requireAll is passed'
      );
    });
  });

  module('pathWithNamespace', function () {
    test('appends the namespace to the path if there is one', function (assert) {
      const namespaceService = Service.extend({
        path: 'marketing',
      });
      this.owner.register('service:namespace', namespaceService);
      assert.strictEqual(this.service.pathNameWithNamespace('sys/auth'), 'marketing/sys/auth');
    });

    test('appends the chroot and namespace when both present', function (assert) {
      const namespaceService = Service.extend({
        path: 'marketing',
      });
      this.owner.register('service:namespace', namespaceService);
      this.service.set('chrootNamespace', 'admin/');
      assert.strictEqual(this.service.pathNameWithNamespace('sys/auth'), 'admin/marketing/sys/auth');
    });
    test('appends the chroot when no namespace', function (assert) {
      this.service.set('chrootNamespace', 'admin');
      assert.strictEqual(this.service.pathNameWithNamespace('sys/auth'), 'admin/sys/auth');
    });
    test('handles superfluous slashes', function (assert) {
      const namespaceService = Service.extend({
        path: '/marketing',
      });
      this.owner.register('service:namespace', namespaceService);
      this.service.set('chrootNamespace', '/admin/');
      assert.strictEqual(this.service.pathNameWithNamespace('/sys/auth'), 'admin/marketing/sys/auth');
      assert.strictEqual(
        this.service.pathNameWithNamespace('/sys/policies/'),
        'admin/marketing/sys/policies/'
      );
    });
  });
});
