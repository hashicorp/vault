/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('sets paths properly', async function (assert) {
    const service = this.owner.lookup('service:permissions');
    await service.getPaths.perform();
    assert.deepEqual(service.get('exactPaths'), PERMISSIONS_RESPONSE.data.exact_paths);
    assert.deepEqual(service.get('globPaths'), PERMISSIONS_RESPONSE.data.glob_paths);
  });

  test('returns true if a policy includes access to an exact path', function (assert) {
    const service = this.owner.lookup('service:permissions');
    service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
    assert.true(service.hasPermission('foo'));
  });

  test('returns true if a paths prefix is included in the policys exact paths', function (assert) {
    const service = this.owner.lookup('service:permissions');
    service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
    assert.true(service.hasPermission('bar'));
  });

  test('it returns true if a policy includes access to a glob path', function (assert) {
    const service = this.owner.lookup('service:permissions');
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.true(service.hasPermission('baz/biz/hi'));
  });

  test('it returns true if a policy includes access to the * glob path', function (assert) {
    const service = this.owner.lookup('service:permissions');
    const splatPath = { '': {} };
    service.set('globPaths', splatPath);
    assert.true(service.hasPermission('hi'));
  });

  test('it returns false if the matched path includes the deny capability', function (assert) {
    const service = this.owner.lookup('service:permissions');
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.false(service.hasPermission('boo'));
  });

  test('it returns true if passed path does not end in a slash but globPath does', function (assert) {
    const service = this.owner.lookup('service:permissions');
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.true(service.hasPermission('ends/in/slash'), 'matches without slash');
    assert.true(service.hasPermission('ends/in/slash/'), 'matches with slash');
  });

  test('it returns false if a policy does not includes access to a path', function (assert) {
    const service = this.owner.lookup('service:permissions');
    assert.false(service.hasPermission('danger'));
  });

  test('sets the root token', function (assert) {
    const service = this.owner.lookup('service:permissions');
    service.setPaths({ data: { root: true } });
    assert.true(service.canViewAll);
  });

  test('returns true with the root token', function (assert) {
    const service = this.owner.lookup('service:permissions');
    service.set('canViewAll', true);
    assert.true(service.hasPermission('hi'));
  });

  test('it returns true if a policy has the specified capabilities on a path', function (assert) {
    const service = this.owner.lookup('service:permissions');
    service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.true(service.hasPermission('bar/bee', ['create', 'list']));
    assert.true(service.hasPermission('baz/biz', ['read']));
  });

  test('it returns false if a policy does not have the specified capabilities on a path', function (assert) {
    const service = this.owner.lookup('service:permissions');
    service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.false(service.hasPermission('bar/bee', ['create', 'delete']));
    assert.false(service.hasPermission('foo', ['create']));
  });

  test('defaults to show all items when policy cannot be found', async function (assert) {
    const service = this.owner.lookup('service:permissions');
    this.server.get('/v1/sys/internal/ui/resultant-acl', () => {
      return [403, { 'Content-Type': 'application/json' }];
    });
    await service.getPaths.perform();
    assert.true(service.canViewAll);
  });

  test('returns the first allowed nav route for policies', function (assert) {
    const service = this.owner.lookup('service:permissions');
    const policyPaths = {
      'sys/policies/acl': {
        capabilities: ['deny'],
      },
      'sys/policies/rgp': {
        capabilities: ['read'],
      },
    };
    service.set('exactPaths', policyPaths);
    assert.strictEqual(service.navPathParams('policies').models[0], 'rgp');
  });

  test('returns the first allowed nav route for access', function (assert) {
    const service = this.owner.lookup('service:permissions');
    const accessPaths = {
      'sys/auth': {
        capabilities: ['deny'],
      },
      'identity/entity/id': {
        capabilities: ['read'],
      },
    };
    const expected = { route: 'vault.cluster.access.identity', models: ['entities'] };
    service.set('exactPaths', accessPaths);
    assert.deepEqual(service.navPathParams('access'), expected);
  });

  test('hasNavPermission returns true if a policy includes the required capabilities for at least one path', function (assert) {
    const service = this.owner.lookup('service:permissions');
    const accessPaths = {
      'sys/auth': {
        capabilities: ['deny'],
      },
      'identity/group/id': {
        capabilities: ['list', 'read'],
      },
    };
    service.set('exactPaths', accessPaths);
    assert.true(service.hasNavPermission('access', 'groups'));
  });

  test('hasNavPermission returns false if a policy does not include the required capabilities for at least one path', function (assert) {
    const service = this.owner.lookup('service:permissions');
    const accessPaths = {
      'sys/auth': {
        capabilities: ['deny'],
      },
      'identity/group/id': {
        capabilities: ['read'],
      },
    };
    service.set('exactPaths', accessPaths);
    assert.false(service.hasNavPermission('access', 'groups'));
  });

  test('appends the namespace to the path if there is one', function (assert) {
    const namespaceService = Service.extend({
      path: 'marketing',
    });
    this.owner.register('service:namespace', namespaceService);
    const service = this.owner.lookup('service:permissions');
    assert.strictEqual(service.pathNameWithNamespace('sys/auth'), 'marketing/sys/auth');
  });
});
