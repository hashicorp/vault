import { moduleFor, test } from 'ember-qunit';
import sinon from 'sinon';

import Ember from 'ember';

let versionStub = Ember.Service.extend();
let routerStub = Ember.Service.extend({
  transitionTo: sinon.stub(),
  urlFor: sinon.stub(),
});

moduleFor('service:control-group', 'Unit | Service | control group', {
  beforeEach() {
    this.register('service:version', versionStub);
    this.inject.service('version', { as: 'version' });
    this.register('service:router', routerStub);
    this.inject.service('router', { as: 'router' });
  },
  afterEach() {},
});

let isOSS = (context) => Ember.set(context, 'version.isOSS', true);
let isEnt = (context) => Ember.set(context, 'version.isOSS', false);
let resolvesArgs = (assert, result, expectedArgs) => {
  return result.then((...args) => {
    return assert.deepEqual(args, expectedArgs, 'resolves with the passed args');
  });
};

[
  [
    'it resolves isOSS:true, wrapTTL: true, response: has wrap_info',
    isOSS,
    [[{'one': 'two', 'three': 'four'}], {wrap_info: {token: 'foo', accessor: 'bar'}}, true],
    (assert, result) => resolvesArgs(assert, result, [{'one': 'two', 'three': 'four'}])
  ],
  [
    'it resolves isOSS:true, wrapTTL: false, response: has no wrap_info',
    isOSS,
    [[{'one': 'two', 'three': 'four'}], {wrap_info: null}, false],
    (assert, result) => resolvesArgs(assert, result, [{'one': 'two', 'three': 'four'}])
  ],
  [
    'it resolves isOSS: false and wrapTTL:true response: has wrap_info',
    isEnt,
    [[{'one': 'two', 'three': 'four'}], {wrap_info: {token: 'foo', accessor: 'bar'}}, true],
    (assert, result) => resolvesArgs(assert, result, [{'one': 'two', 'three': 'four'}])
  ],
  [
    'it resolves isOSS: false and wrapTTL:false response: has no wrap_info',
    isEnt,
    [[{'one': 'two', 'three': 'four'}], {wrap_info: null}, false],
    (assert, result) => resolvesArgs(assert, result, [{'one': 'two', 'three': 'four'}])
  ],
 [
    'it rejects isOSS: false, wrapTTL:false, response: has wrap_info',
    isEnt,
    [[{'one': 'two', 'three': 'four'}], {foo: 'bar', wrap_info: {token: 'secret', accessor: 'lookup'}}, false],
    (assert, result) => {
      // ensure failure if we ever don't reject
      assert.expect(2);

      return result.then(
        () => {},
        (err) => {
          assert.equal(err.token, 'secret')
          assert.equal(err.accessor, 'lookup')
        });
    }
  ],

].forEach(function([name, setup, args, expectation]) {
  test(`checkForControlGroup: ${name}`, function(assert) {
    if (setup) {
      setup(this);
    }
    let service = this.subject();
    let result = service.checkForControlGroup(...args);
    return expectation(assert, result);
  });
});

test(`handleError: transitions to accessor when there is no transition passed in`, function(assert) {
  let service = this.subject();
  service.handleError({ accessor: '12345'});
  assert.ok(this.router.transitionTo.calledWith('vault.cluster.access.control-group-accessor', '12345'));
});

test(`logFromError: returns correct content string`, function(assert) {
  let service = this.subject();
  let contentString = service.logFromError({ accessor: '12345', creation_path: '/this/path/', token: 'asdf'});
  assert.ok(this.router.urlFor.calledWith('vault.cluster.access.control-group-accessor', '12345'), 'calls urlFor with accessor');
  assert.ok(contentString.content.includes('12345'), 'contains accessor');
  assert.ok(contentString.content.includes('/this/path/'), 'contains creation path');
  assert.ok(contentString.content.includes('asdf'), 'contains token');
});

