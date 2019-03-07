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

module('Unit | Service | permissions', function(hooks) {
  setupTest(hooks);

  hooks.beforeEach(function() {
    this.server = new Pretender();
    this.server.get('/v1/sys/internal/ui/resultant-acl', () => {
      return [200, { 'Content-Type': 'application/json' }, JSON.stringify(PERMISSIONS_RESPONSE)];
    });
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('sets paths properly', async function(assert) {
    let service = this.owner.lookup('service:permissions');
    await service.getPaths.perform();
    assert.deepEqual(service.get('exactPaths'), PERMISSIONS_RESPONSE.data.exact_paths);
    assert.deepEqual(service.get('globPaths'), PERMISSIONS_RESPONSE.data.glob_paths);
  });

  test('returns true if a policy includes access to an exact path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
    assert.equal(service.hasPermission('foo'), true);
  });

  test('returns true if a paths prefix is included in the policys exact paths', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
    assert.equal(service.hasPermission('bar'), true);
  });

  test('it returns true if a policy includes access to a glob path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.equal(service.hasPermission('baz/biz/hi'), true);
  });

  test('it returns true if a policy includes access to the * glob path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    const splatPath = { '': {} };
    service.set('globPaths', splatPath);
    assert.equal(service.hasPermission('hi'), true);
  });

  test('it returns false if the matched path includes the deny capability', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.equal(service.hasPermission('boo'), false);
  });

  test('it returns true if passed path does not end in a slash but globPath does', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.equal(service.hasPermission('ends/in/slash'), true, 'matches without slash');
    assert.equal(service.hasPermission('ends/in/slash/'), true, 'matches with slash');
  });

  test('it returns false if a policy does not includes access to a path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    assert.equal(service.hasPermission('danger'), false);
  });

  test('sets the root token', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.setPaths({ data: { root: true } });
    assert.equal(service.canViewAll, true);
  });

  test('returns true with the root token', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('canViewAll', true);
    assert.equal(service.hasPermission('hi'), true);
  });

  test('it returns true if a policy has the specified capabilities on a path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.equal(service.hasPermission('bar/bee', ['create', 'list']), true);
    assert.equal(service.hasPermission('baz/biz', ['read']), true);
  });

  test('it returns false if a policy does not have the specified capabilities on a path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('exactPaths', PERMISSIONS_RESPONSE.data.exact_paths);
    service.set('globPaths', PERMISSIONS_RESPONSE.data.glob_paths);
    assert.equal(service.hasPermission('bar/bee', ['create', 'delete']), false);
    assert.equal(service.hasPermission('foo', ['create']), false);
  });

  test('defaults to show all items when policy cannot be found', async function(assert) {
    let service = this.owner.lookup('service:permissions');
    this.server.get('/v1/sys/internal/ui/resultant-acl', () => {
      return [403, { 'Content-Type': 'application/json' }];
    });
    await service.getPaths.perform();
    assert.equal(service.canViewAll, true);
  });

  test('returns the first allowed nav route for policies', function(assert) {
    let service = this.owner.lookup('service:permissions');
    const policyPaths = {
      'sys/policies/acl': {
        capabilities: ['deny'],
      },
      'sys/policies/rgp': {
        capabilities: ['read'],
      },
    };
    service.set('exactPaths', policyPaths);
    assert.equal(service.navPathParams('policies'), 'rgp');
  });

  test('returns the first allowed nav route for access', function(assert) {
    let service = this.owner.lookup('service:permissions');
    const accessPaths = {
      'sys/auth': {
        capabilities: ['deny'],
      },
      'identity/entities': {
        capabilities: ['read'],
      },
    };
    const expected = ['vault.cluster.access.identity', 'entities'];
    service.set('exactPaths', accessPaths);
    assert.deepEqual(service.navPathParams('access'), expected);
  });

  test('hasNavPermission returns true if a policy includes access to at least one path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    const accessPaths = {
      'sys/auth': {
        capabilities: ['deny'],
      },
      'sys/leases/lookup': {
        capabilities: ['read'],
      },
    };
    service.set('exactPaths', accessPaths);
    assert.equal(service.hasNavPermission('access', 'leases'), true);
  });

  test('hasNavPermission returns false if a policy does not include access to any paths', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('exactPaths', {});
    assert.equal(service.hasNavPermission('access'), false);
  });

  test('appends the namespace to the path if there is one', function(assert) {
    const namespaceService = Service.extend({
      path: 'marketing',
    });
    this.owner.register('service:namespace', namespaceService);
    let service = this.owner.lookup('service:permissions');
    assert.equal(service.pathNameWithNamespace('sys/auth'), 'marketing/sys/auth');
  });
});
