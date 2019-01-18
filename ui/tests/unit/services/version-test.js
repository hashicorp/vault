import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | version', function(hooks) {
  setupTest(hooks);

  test('setting version computes isOSS properly', function(assert) {
    let service = this.owner.lookup('service:version');
    service.set('version', '0.9.5');
    assert.equal(service.get('isOSS'), true);
    assert.equal(service.get('isEnterprise'), false);
  });

  test('setting version computes isEnterprise properly', function(assert) {
    let service = this.owner.lookup('service:version');
    service.set('version', '0.9.5+prem');
    assert.equal(service.get('isOSS'), false);
    assert.equal(service.get('isEnterprise'), true);
  });

  test('setting version with hsm ending computes isEnterprise properly', function(assert) {
    let service = this.owner.lookup('service:version');
    service.set('version', '0.9.5+prem.hsm');
    assert.equal(service.get('isOSS'), false);
    assert.equal(service.get('isEnterprise'), true);
  });

  test('hasPerfReplication', function(assert) {
    let service = this.owner.lookup('service:version');
    assert.equal(service.get('hasPerfReplication'), false);
    service.set('_features', ['Performance Replication']);
    assert.equal(service.get('hasPerfReplication'), true);
  });

  test('hasDRReplication', function(assert) {
    let service = this.owner.lookup('service:version');
    assert.equal(service.get('hasDRReplication'), false);
    service.set('_features', ['DR Replication']);
    assert.equal(service.get('hasDRReplication'), true);
  });
});
