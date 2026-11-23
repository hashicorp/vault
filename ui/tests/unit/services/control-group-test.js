/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { set } from '@ember/object';
import { equal } from '@ember/object/computed';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

import { storageKey, CONTROL_GROUP_PREFIX, TOKEN_SEPARATOR } from 'vault/services/control-group';

const versionStub = Service.extend({
  isCommunity: equal('type', 'community'),
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

module('Unit | Service | control group', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.owner.register('service:version', versionStub);
    this.version = this.owner.lookup('service:version');
    this.router = this.owner.lookup('service:router');
    this.router.reopen({
      transitionTo: sinon.stub(),
      urlFor: sinon.stub().returns('/ui/vault/foo'),
      currentURL: '/vault/secrets-engines/kv/show/foo',
    });
  });

  hooks.afterEach(function () {});

  const isCommunity = (context) => set(context, 'version.type', 'community');
  const isEnt = (context) => set(context, 'version.type', 'enterprise');

  [
    [
      'it is not locked when isCommunity:true, wrapTTL: true, response: has wrap_info',
      isCommunity,
      [{ wrap_info: { token: 'foo', accessor: 'bar' } }, true],
      false,
    ],
    [
      'it is not locked when isCommunity:true, wrapTTL: false, response: has no wrap_info',
      isCommunity,
      [{ wrap_info: null }, false],
      false,
    ],
    [
      'it is not locked when isCommunity:true, wrapTTL: false, response: has wrap_info',
      isCommunity,
      [{ wrap_info: { token: 'foo', accessor: 'bar' } }, false],
      false,
    ],
    [
      'it is not locked when isCommunity: false and wrapTTL:true response: has wrap_info',
      isEnt,
      [{ wrap_info: { token: 'foo', accessor: 'bar' } }, true],
      false,
    ],
    [
      'it is not locked when isCommunity: false and wrapTTL:false response: has no wrap_info',
      isEnt,
      [{ wrap_info: null }, false],
      false,
    ],
    [
      'it is not locked when the response is a wrapping endpoint response',
      isEnt,
      [{ wrap_info: { token: 'foo', accessor: 'bar', creation_path: 'sys/wrapping/wrap' } }, false],
      false,
    ],
    [
      'it is locked when isCommunity: false, wrapTTL:false, response: has wrap_info',
      isEnt,
      [{ foo: 'bar', wrap_info: { token: 'secret', accessor: 'lookup' } }, false],
      true,
    ],
  ].forEach(function ([name, setup, args, expected]) {
    test(`isRequestedPathLocked: ${name}`, function (assert) {
      setup(this);
      const service = this.owner.lookup('service:control-group');
      assert.strictEqual(service.isRequestedPathLocked(...args), expected);
    });
  });

  test(`handleError: transitions to accessor and stores control group token`, function (assert) {
    const error = {
      accessor: '12345',
      token: 'token',
      creation_path: 'kv/',
      creation_time: '2022-03-17T20:00:25.594Z',
      ttl: 400,
    };
    const expected = { ...error, uiParams: { url: '/vault/secrets-engines/kv/show/foo' } };
    const service = this.owner.factoryFor('service:control-group').create({
      storeControlGroupToken: sinon.spy(),
    });
    service.handleError(error);
    assert.ok(service.storeControlGroupToken.calledWith(expected), 'calls storeControlGroupToken');
    assert.ok(
      this.router.transitionTo.calledWith('vault.cluster.access.control-group-accessor', '12345'),
      'calls router transitionTo'
    );
  });

  test(`logFromError: returns correct content string`, function (assert) {
    const error = {
      accessor: '12345',
      token: 'token',
      creation_path: 'kv/',
      creation_time: '2022-03-17T20:00:25.594Z',
      ttl: 400,
    };
    const service = this.owner.factoryFor('service:control-group').create({
      storeControlGroupToken: sinon.spy(),
    });
    const contentString = service.logFromError(error);
    assert.ok(
      this.router.urlFor.calledWith('vault.cluster.access.control-group-accessor', '12345'),
      'calls urlFor with accessor'
    );
    assert.ok(service.storeControlGroupToken.calledWith(error), 'calls storeControlGroupToken');
    assert.ok(contentString.content.includes('12345'), 'contains accessor');
    assert.ok(contentString.content.includes('kv/'), 'contains creation path');
    assert.ok(contentString.content.includes('token'), 'contains token');
  });

  test('storageKey', function (assert) {
    const accessor = '12345';
    const path = 'kv/foo/bar';
    const expectedKey = `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
    assert.strictEqual(storageKey(accessor, path), expectedKey, 'uses expected key');
  });

  test('keyFromAccessor', function (assert) {
    const store = storage();
    const accessor = '12345';
    const path = 'kv/foo/bar';
    const data = { foo: 'bar' };
    const expectedKey = `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
    const subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });

    store.setItem(expectedKey, data);
    store.setItem(`${CONTROL_GROUP_PREFIX}2345${TOKEN_SEPARATOR}${path}`, 'ok');

    assert.strictEqual(subject.keyFromAccessor(accessor), expectedKey, 'finds key given the accessor');
    assert.strictEqual(subject.keyFromAccessor('foo'), null, 'returns null if no key was found');
  });

  test('storeControlGroupToken', function (assert) {
    const store = storage();
    const subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });
    const info = {
      accessor: '12345',
      creation_path: 'foo/',
      creation_time: '2022-03-17T20:00:25.594Z',
      ttl: 300,
    };
    const key = `${CONTROL_GROUP_PREFIX}${info.accessor}${TOKEN_SEPARATOR}${info.creation_path}`;

    subject.storeControlGroupToken(info);
    assert.deepEqual(store.items[key], JSON.stringify(info), 'stores the whole info object');
  });

  test('deleteControlGroupToken', function (assert) {
    const store = storage();
    const subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });
    const accessor = 'foo';
    const path = 'kv/one';

    const expectedKey = `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
    store.setItem(expectedKey, { one: '2' });
    subject.deleteControlGroupToken(accessor);
    assert.strictEqual(Object.keys(store.items).length, 0, 'there are no keys stored in storage');
  });

  test('deleteTokens', function (assert) {
    const store = storage();
    const subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });

    const keyOne = `${CONTROL_GROUP_PREFIX}foo`;
    const keyTwo = `${CONTROL_GROUP_PREFIX}bar`;
    store.setItem(keyOne, { one: '2' });
    store.setItem(keyTwo, { two: '2' });
    store.setItem('value', 'one');
    assert.strictEqual(Object.keys(store.items).length, 3, 'stores 3 values');
    subject.deleteTokens();
    assert.strictEqual(Object.keys(store.items).length, 1, 'removes tokens with control group prefix');
    assert.strictEqual(store.getItem('value'), 'one', 'keeps the non-prefixed value');
  });

  test('wrapInfoForAccessor', function (assert) {
    const store = storage();
    const subject = this.owner.factoryFor('service:control-group').create({
      storage() {
        return store;
      },
    });

    const keyOne = `${CONTROL_GROUP_PREFIX}foo`;
    store.setItem(keyOne, { one: '2' });
    assert.deepEqual(subject.wrapInfoForAccessor('foo'), { one: '2' });
  });
});
