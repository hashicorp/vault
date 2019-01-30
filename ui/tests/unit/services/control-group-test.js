import { set } from '@ember/object';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

import { storageKey, CONTROL_GROUP_PREFIX, TOKEN_SEPARATOR } from 'vault/services/control-group';

let versionStub = Service.extend();
let routerStub = Service.extend({
  transitionTo: sinon.stub(),
  urlFor: sinon.stub().returns('/ui/vault/foo'),
});

function storage() {
  return {
    items: {},
    getItem(key) {
      var item = this.items[key];
      return item && JSON.parse(item);
    },

    setItem(key, val) {
      return (this.items[key] = JSON.stringify(val));
    },

    removeItem(key) {
      delete this.items[key];
    },

    keys() {
      return Object.keys(this.items);
    },
  };
}

module('Unit | Service | control group', function(hooks) {
  setupTest(hooks);

  hooks.beforeEach(function() {
    this.owner.register('service:version', versionStub);
    this.version = this.owner.lookup('service:version');
    this.owner.register('service:router', routerStub);
    this.router = this.owner.lookup('service:router');
  });

  hooks.afterEach(function() {});

  let isOSS = context => set(context, 'version.isOSS', true);
  let isEnt = context => set(context, 'version.isOSS', false);
  let resolvesArgs = (assert, result, expectedArgs) => {
    return result.then((...args) => {
      return assert.deepEqual(args, expectedArgs, 'resolves with the passed args');
    });
  };

  [
    [
      'it resolves isOSS:true, wrapTTL: true, response: has wrap_info',
      isOSS,
      [[{ one: 'two', three: 'four' }], { wrap_info: { token: 'foo', accessor: 'bar' } }, true],
      (assert, result) => resolvesArgs(assert, result, [{ one: 'two', three: 'four' }]),
    ],
    [
      'it resolves isOSS:true, wrapTTL: false, response: has no wrap_info',
      isOSS,
      [[{ one: 'two', three: 'four' }], { wrap_info: null }, false],
      (assert, result) => resolvesArgs(assert, result, [{ one: 'two', three: 'four' }]),
    ],
    [
      'it resolves isOSS: false and wrapTTL:true response: has wrap_info',
      isEnt,
      [[{ one: 'two', three: 'four' }], { wrap_info: { token: 'foo', accessor: 'bar' } }, true],
      (assert, result) => resolvesArgs(assert, result, [{ one: 'two', three: 'four' }]),
    ],
    [
      'it resolves isOSS: false and wrapTTL:false response: has no wrap_info',
      isEnt,
      [[{ one: 'two', three: 'four' }], { wrap_info: null }, false],
      (assert, result) => resolvesArgs(assert, result, [{ one: 'two', three: 'four' }]),
    ],
    [
      'it rejects isOSS: false, wrapTTL:false, response: has wrap_info',
      isEnt,
      [
        [{ one: 'two', three: 'four' }],
        { foo: 'bar', wrap_info: { token: 'secret', accessor: 'lookup' } },
        false,
      ],
      (assert, result) => {
        // ensure failure if we ever don't reject
        assert.expect(2);

        return result.then(
          () => {},
          err => {
            assert.equal(err.token, 'secret');
            assert.equal(err.accessor, 'lookup');
          }
        );
      },
    ],
  ].forEach(function([name, setup, args, expectation]) {
    test(`checkForControlGroup: ${name}`, function(assert) {
      if (setup) {
        setup(this);
      }
      let service = this.owner.lookup('service:control-group');
      let result = service.checkForControlGroup(...args);
      return expectation(assert, result);
    });
  });

  test(`handleError: transitions to accessor when there is no transition passed in`, function(assert) {
    let error = {
      accessor: '12345',
      token: 'token',
      creation_path: 'kv/',
      creation_time: new Date().toISOString(),
      ttl: 400,
    };
    let url;
    let expected = { ...error, uiParams: { url } };
    let transition = {
      targetName: 'vault.cluster.foo',
    };
    let service = this.owner.factoryFor('service:control-group').create({
      urlFromTransition: sinon.spy(),
      storeControlGroupToken: sinon.spy(),
    });
    service.handleError(error, transition);
    assert.ok(service.urlFromTransition.calledWith(transition), 'calls urlFromTransition');
    assert.ok(service.storeControlGroupToken.calledWith(expected), 'calls storeControlGroupToken');
    assert.ok(
      this.router.transitionTo.calledWith('vault.cluster.access.control-group-accessor', '12345'),
      'calls router transitionTo'
    );
  });

  test(`logFromError: returns correct content string`, function(assert) {
    let error = {
      accessor: '12345',
      token: 'token',
      creation_path: 'kv/',
      creation_time: new Date().toISOString(),
      ttl: 400,
    };
    let service = this.owner.factoryFor('service:control-group').create({
      storeControlGroupToken: sinon.spy(),
    });
    let contentString = service.logFromError(error);
    assert.ok(
      this.router.urlFor.calledWith('vault.cluster.access.control-group-accessor', '12345'),
      'calls urlFor with accessor'
    );
    assert.ok(service.storeControlGroupToken.calledWith(error), 'calls storeControlGroupToken');
    assert.ok(contentString.content.includes('12345'), 'contains accessor');
    assert.ok(contentString.content.includes('kv/'), 'contains creation path');
    assert.ok(contentString.content.includes('token'), 'contains token');
  });

  test('urlFromTransition', function(assert) {
    let transition = {
      targetName: 'vault.cluster.foo',
      params: {
        vault: {},
        cluster: { cluster_name: 'vault' },
        foo: { bar: '1' },
      },
      queryParams: {},
    };
    let expected = [transition.targetName, { cluster_name: 'vault' }, { bar: '1' }, { queryParams: {} }];
    let service = this.owner.lookup('service:control-group');
    service.urlFromTransition(transition);
    assert.ok(this.router.urlFor.calledWith(...expected), 'calls urlFor with expected args');
  });

  test('storageKey', function(assert) {
    let accessor = '12345';
    let path = 'kv/foo/bar';
    let expectedKey = `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
    assert.equal(storageKey(accessor, path), expectedKey, 'uses expected key');
  });

  test('keyFromAccessor', function(assert) {
    let store = storage();
    let accessor = '12345';
    let path = 'kv/foo/bar';
    let data = { foo: 'bar' };
    let expectedKey = `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
    let subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });

    store.setItem(expectedKey, data);
    store.setItem(`${CONTROL_GROUP_PREFIX}2345${TOKEN_SEPARATOR}${path}`, 'ok');

    assert.equal(subject.keyFromAccessor(accessor), expectedKey, 'finds key given the accessor');
    assert.equal(subject.keyFromAccessor('foo'), null, 'returns null if no key was found');
  });

  test('storeControlGroupToken', function(assert) {
    let store = storage();
    let subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });
    let info = {
      accessor: '12345',
      creation_path: 'foo/',
      creation_time: new Date().toISOString(),
      ttl: 300,
    };
    let key = `${CONTROL_GROUP_PREFIX}${info.accessor}${TOKEN_SEPARATOR}${info.creation_path}`;

    subject.storeControlGroupToken(info);
    assert.deepEqual(store.items[key], JSON.stringify(info), 'stores the whole info object');
  });

  test('deleteControlGroupToken', function(assert) {
    let store = storage();
    let subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });
    let accessor = 'foo';
    let path = 'kv/one';

    let expectedKey = `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
    store.setItem(expectedKey, { one: '2' });
    subject.deleteControlGroupToken(accessor);
    assert.equal(Object.keys(store.items).length, 0, 'there are no keys stored in storage');
  });

  test('deleteTokens', function(assert) {
    let store = storage();
    let subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });

    let keyOne = `${CONTROL_GROUP_PREFIX}foo`;
    let keyTwo = `${CONTROL_GROUP_PREFIX}bar`;
    store.setItem(keyOne, { one: '2' });
    store.setItem(keyTwo, { two: '2' });
    store.setItem('value', 'one');
    assert.equal(Object.keys(store.items).length, 3, 'stores 3 values');
    subject.deleteTokens();
    assert.equal(Object.keys(store.items).length, 1, 'removes tokens with control group prefix');
    assert.equal(store.getItem('value'), 'one', 'keeps the non-prefixed value');
  });

  test('wrapInfoForAccessor', function(assert) {
    let store = storage();
    let subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });

    let keyOne = `${CONTROL_GROUP_PREFIX}foo`;
    store.setItem(keyOne, { one: '2' });
    assert.deepEqual(subject.wrapInfoForAccessor('foo'), { one: '2' });
  });
});
