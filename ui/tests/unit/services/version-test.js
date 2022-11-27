import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | version', function (hooks) {
  setupTest(hooks);

  test('setting version computes isOSS properly', function (assert) {
    const service = this.owner.lookup('service:version');
    service.set('version', '0.9.5');
    assert.true(service.get('isOSS'));
    assert.false(service.get('isEnterprise'));
  });

  test('setting version computes isEnterprise properly', function (assert) {
    const service = this.owner.lookup('service:version');
    service.set('version', '0.9.5+prem');
    assert.false(service.get('isOSS'));
    assert.true(service.get('isEnterprise'));
  });

  test('setting version with hsm ending computes isEnterprise properly', function (assert) {
    const service = this.owner.lookup('service:version');
    service.set('version', '0.9.5+prem.hsm');
    assert.false(service.get('isOSS'));
    assert.true(service.get('isEnterprise'));
  });

  test('hasPerfReplication', function (assert) {
    const service = this.owner.lookup('service:version');
    assert.false(service.get('hasPerfReplication'));
    service.set('_features', ['Performance Replication']);
    assert.true(service.get('hasPerfReplication'));
  });

  test('hasDRReplication', function (assert) {
    const service = this.owner.lookup('service:version');
    assert.false(service.get('hasDRReplication'));
    service.set('_features', ['DR Replication']);
    assert.true(service.get('hasDRReplication'));
  });
});
