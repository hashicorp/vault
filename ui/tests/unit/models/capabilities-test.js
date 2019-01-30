import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { SUDO_PATHS, SUDO_PATH_PREFIXES } from 'vault/models/capabilities';

import { run } from '@ember/runloop';

module('Unit | Model | capabilities', function(hooks) {
  setupTest(hooks);

  test('it exists', function(assert) {
    let model = run(() => this.owner.lookup('service:store').createRecord('capabilities'));
    assert.ok(!!model);
  });

  test('it reads capabilities', function(assert) {
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('capabilities', {
        path: 'foo',
        capabilities: ['list', 'read'],
      })
    );

    assert.ok(model.get('canRead'));
    assert.ok(model.get('canList'));
    assert.notOk(model.get('canUpdate'));
    assert.notOk(model.get('canDelete'));
  });

  test('it allows everything if root is present', function(assert) {
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('capabilities', {
        path: 'foo',
        capabilities: ['root', 'deny', 'read'],
      })
    );
    assert.ok(model.get('canRead'));
    assert.ok(model.get('canCreate'));
    assert.ok(model.get('canUpdate'));
    assert.ok(model.get('canDelete'));
    assert.ok(model.get('canList'));
  });

  test('it denies everything if deny is present', function(assert) {
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('capabilities', {
        path: 'foo',
        capabilities: ['sudo', 'deny', 'read'],
      })
    );
    assert.notOk(model.get('canRead'));
    assert.notOk(model.get('canCreate'));
    assert.notOk(model.get('canUpdate'));
    assert.notOk(model.get('canDelete'));
    assert.notOk(model.get('canList'));
  });

  test('it requires sudo on sudo paths', function(assert) {
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('capabilities', {
        path: SUDO_PATHS[0],
        capabilities: ['sudo', 'read'],
      })
    );
    assert.ok(model.get('canRead'));
    assert.notOk(model.get('canCreate'), 'sudo requires the capability to be set as well');
    assert.notOk(model.get('canUpdate'));
    assert.notOk(model.get('canDelete'));
    assert.notOk(model.get('canList'));
  });

  test('it requires sudo on sudo paths prefixes', function(assert) {
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('capabilities', {
        path: SUDO_PATH_PREFIXES[0] + '/foo',
        capabilities: ['sudo', 'read'],
      })
    );
    assert.ok(model.get('canRead'));
    assert.notOk(model.get('canCreate'), 'sudo requires the capability to be set as well');
    assert.notOk(model.get('canUpdate'));
    assert.notOk(model.get('canDelete'));
    assert.notOk(model.get('canList'));
  });

  test('it does not require sudo on sys/leases/revoke if update capability is present and path is not fully a sudo prefix', function(assert) {
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('capabilities', {
        path: 'sys/leases/revoke',
        capabilities: ['update', 'read'],
      })
    );
    assert.ok(model.get('canRead'));
    assert.notOk(model.get('canCreate'), 'sudo requires the capability to be set as well');
    assert.ok(model.get('canUpdate'), 'should not require sudo if it has update');
    assert.notOk(model.get('canDelete'));
    assert.notOk(model.get('canList'));
  });

  test('it requires sudo on prefix path even if capability is present', function(assert) {
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('capabilities', {
        path: SUDO_PATH_PREFIXES[0] + '/aws',
        capabilities: ['update', 'read'],
      })
    );
    assert.notOk(model.get('canRead'));
    assert.notOk(model.get('canCreate'));
    assert.notOk(model.get('canUpdate'), 'should still require sudo');
    assert.notOk(model.get('canDelete'));
    assert.notOk(model.get('canList'));
  });

  test('it does not require sudo on prefix path if both update and sudo capabilities are present', function(assert) {
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('capabilities', {
        path: SUDO_PATH_PREFIXES[0] + '/aws',
        capabilities: ['sudo', 'update', 'read'],
      })
    );
    assert.ok(model.get('canRead'));
    assert.notOk(model.get('canCreate'));
    assert.ok(model.get('canUpdate'), 'should not require sudo');
    assert.notOk(model.get('canDelete'));
    assert.notOk(model.get('canList'));
  });
});
