import { moduleForModel, test } from 'ember-qunit';
import { SUDO_PATHS, SUDO_PATH_PREFIXES } from 'vault/models/capabilities';

moduleForModel('capabilities', 'Unit | Model | capabilities', {
  needs: ['transform:array'],
});

test('it exists', function(assert) {
  let model = this.subject();
  assert.ok(!!model);
});

test('it reads capabilities', function(assert) {
  let model = this.subject({
    path: 'foo',
    capabilities: ['list', 'read'],
  });

  assert.ok(model.get('canRead'));
  assert.ok(model.get('canList'));
  assert.notOk(model.get('canUpdate'));
  assert.notOk(model.get('canDelete'));
});

test('it allows everything if root is present', function(assert) {
  let model = this.subject({
    path: 'foo',
    capabilities: ['root', 'deny', 'read'],
  });
  assert.ok(model.get('canRead'));
  assert.ok(model.get('canCreate'));
  assert.ok(model.get('canUpdate'));
  assert.ok(model.get('canDelete'));
  assert.ok(model.get('canList'));
});

test('it denies everything if deny is present', function(assert) {
  let model = this.subject({
    path: 'foo',
    capabilities: ['sudo', 'deny', 'read'],
  });
  assert.notOk(model.get('canRead'));
  assert.notOk(model.get('canCreate'));
  assert.notOk(model.get('canUpdate'));
  assert.notOk(model.get('canDelete'));
  assert.notOk(model.get('canList'));
});

test('it requires sudo on sudo paths', function(assert) {
  let model = this.subject({
    path: SUDO_PATHS[0],
    capabilities: ['sudo', 'read'],
  });
  assert.ok(model.get('canRead'));
  assert.notOk(model.get('canCreate'), 'sudo requires the capability to be set as well');
  assert.notOk(model.get('canUpdate'));
  assert.notOk(model.get('canDelete'));
  assert.notOk(model.get('canList'));
});

test('it requires sudo on sudo paths prefixes', function(assert) {
  let model = this.subject({
    path: SUDO_PATH_PREFIXES[0] + '/foo',
    capabilities: ['sudo', 'read'],
  });
  assert.ok(model.get('canRead'));
  assert.notOk(model.get('canCreate'), 'sudo requires the capability to be set as well');
  assert.notOk(model.get('canUpdate'));
  assert.notOk(model.get('canDelete'));
  assert.notOk(model.get('canList'));
});
