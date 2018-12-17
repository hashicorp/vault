import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import Pretender from 'pretender';

const PERMISSIONS_RESPONSE = {
  data: {
    exact_paths: {
      foo: {
        capabilities: ['read'],
      },
      bar: {
        capabilities: ['create'],
      },
      boo: {
        capabilities: ['deny'],
      },
    },
    glob_paths: {
      'baz/biz': {
        capabilities: ['read'],
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

  test('it returns false if a policy does not includes access to a path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    assert.equal(service.hasPermission('danger'), false);
  });

  test('returns true with the root token', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('isRootToken', true);
    assert.equal(service.hasPermission('hi'), true);
  });
});
