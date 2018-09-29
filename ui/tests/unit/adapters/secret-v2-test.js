import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | secret-v2', function(hooks) {
  setupTest(hooks);

  test('secret api urls', function(assert) {
    let url, method, options;
    let adapter = this.owner.factoryFor('adapter:secret-v2').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve({});
      },
    });

    adapter.query({}, 'secret', { id: '', backend: 'secret' });
    assert.equal(url, '/v1/secret/metadata/', 'query generic url OK');
    assert.equal('GET', method, 'query generic method OK');
    assert.deepEqual(options, { data: { list: true } }, 'query generic url OK');

    adapter.queryRecord({}, 'secret', { id: 'foo', backend: 'secret' });
    assert.equal(url, '/v1/secret/data/foo', 'queryRecord generic url OK');
    assert.equal('GET', method, 'queryRecord generic method OK');
  });
});
