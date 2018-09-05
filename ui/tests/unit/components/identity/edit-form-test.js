import { moduleForComponent, test } from 'ember-qunit';
import sinon from 'sinon';
import Ember from 'ember';

moduleForComponent('identity/edit-form', 'Unit | Component | identity/edit-form', {
  unit: true,
  needs: ['service:auth', 'service:flash-messages'],
});

let testCases = [
  {
    identityType: 'entity',
    mode: 'create',
    expected: 'vault.cluster.access.identity',
  },
  {
    identityType: 'entity',
    mode: 'edit',
    expected: 'vault.cluster.access.identity.show',
  },
  {
    identityType: 'entity-merge',
    mode: 'merge',
    expected: 'vault.cluster.access.identity',
  },
  {
    identityType: 'entity-alias',
    mode: 'create',
    expected: 'vault.cluster.access.identity.aliases',
  },
  {
    identityType: 'entity-alias',
    mode: 'edit',
    expected: 'vault.cluster.access.identity.aliases.show',
  },
  {
    identityType: 'group',
    mode: 'create',
    expected: 'vault.cluster.access.identity',
  },
  {
    identityType: 'group',
    mode: 'edit',
    expected: 'vault.cluster.access.identity.show',
  },
  {
    identityType: 'group-alias',
    mode: 'create',
    expected: 'vault.cluster.access.identity.aliases',
  },
  {
    identityType: 'group-alias',
    mode: 'edit',
    expected: 'vault.cluster.access.identity.aliases.show',
  },
];
testCases.forEach(function(testCase) {
  let model = Ember.Object.create({
    identityType: testCase.identityType,
    rollbackAttributes: sinon.spy(),
  });
  test(`it computes cancelLink properly: ${testCase.identityType} ${testCase.mode}`, function(assert) {
    let component = this.subject();

    component.set('mode', testCase.mode);
    component.set('model', model);
    assert.equal(component.get('cancelLink'), testCase.expected, 'cancel link is correct');
  });
});
