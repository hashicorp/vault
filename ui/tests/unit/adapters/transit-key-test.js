import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | transit key', function(hooks) {
  setupTest(hooks);

  test('transit api urls', function(assert) {
    let url, method, options;
    let adapter = this.owner.factoryFor('adapter:transit-key').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve({});
      },
    });

    adapter.query({}, 'transit-key', { id: '', backend: 'transit' });
    assert.equal(url, '/v1/transit/keys/', 'query list url OK');
    assert.equal('GET', method, 'query list method OK');
    assert.deepEqual(options, { data: { list: true } }, 'query generic url OK');

    adapter.queryRecord({}, 'transit-key', { id: 'foo', backend: 'transit' });
    assert.equal(url, '/v1/transit/keys/foo', 'queryRecord generic url OK');
    assert.equal('GET', method, 'queryRecord generic method OK');

    adapter.keyAction('rotate', { backend: 'transit', id: 'foo', payload: {} });
    assert.equal(url, '/v1/transit/keys/foo/rotate', 'keyAction:rotate url OK');

    adapter.keyAction('encrypt', { backend: 'transit', id: 'foo', payload: {} });
    assert.equal(url, '/v1/transit/encrypt/foo', 'keyAction:encrypt url OK');

    adapter.keyAction('datakey', { backend: 'transit', id: 'foo', payload: { param: 'plaintext' } });
    assert.equal(url, '/v1/transit/datakey/plaintext/foo', 'keyAction:datakey url OK');

    adapter.keyAction('export', { backend: 'transit', id: 'foo', payload: { param: ['hmac'] } });
    assert.equal(url, '/v1/transit/export/hmac-key/foo', 'transitAction:export, no version url OK');

    adapter.keyAction('export', { backend: 'transit', id: 'foo', payload: { param: ['hmac', 10] } });
    assert.equal(url, '/v1/transit/export/hmac-key/foo/10', 'transitAction:export, with version url OK');
  });
});
