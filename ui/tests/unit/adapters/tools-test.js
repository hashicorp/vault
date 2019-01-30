import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | tools', function(hooks) {
  setupTest(hooks);

  test('wrapping api urls', function(assert) {
    let url, method, options;
    let adapter = this.owner.factoryFor('adapter:tools').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    let clientToken;
    let data = { foo: 'bar' };
    adapter.toolAction('wrap', data, { wrapTTL: '30m' });
    assert.equal('/v1/sys/wrapping/wrap', url, 'wrapping:wrap url OK');
    assert.equal('POST', method, 'wrapping:wrap method OK');
    assert.deepEqual({ data: data, wrapTTL: '30m', clientToken }, options, 'wrapping:wrap options OK');

    data = { token: 'token' };
    adapter.toolAction('lookup', data);
    assert.equal('/v1/sys/wrapping/lookup', url, 'wrapping:lookup url OK');
    assert.equal('POST', method, 'wrapping:lookup method OK');
    assert.deepEqual({ data, clientToken }, options, 'wrapping:lookup options OK');

    adapter.toolAction('unwrap', data);
    assert.equal('/v1/sys/wrapping/unwrap', url, 'wrapping:unwrap url OK');
    assert.equal('POST', method, 'wrapping:unwrap method OK');
    assert.deepEqual({ data, clientToken }, options, 'wrapping:unwrap options OK');

    adapter.toolAction('rewrap', data);
    assert.equal('/v1/sys/wrapping/rewrap', url, 'wrapping:rewrap url OK');
    assert.equal('POST', method, 'wrapping:rewrap method OK');
    assert.deepEqual({ data, clientToken }, options, 'wrapping:rewrap options OK');
  });

  test('tools api urls', function(assert) {
    let url, method;
    let adapter = this.owner.factoryFor('adapter:tools').create({
      ajax: (...args) => {
        [url, method] = args;
        return resolve();
      },
    });

    adapter.toolAction('hash', { input: 'someBase64' });
    assert.equal(url, '/v1/sys/tools/hash', 'sys tools hash: url OK');
    assert.equal('POST', method, 'sys tools hash: method OK');

    adapter.toolAction('random', { bytes: '32' });
    assert.equal(url, '/v1/sys/tools/random', 'sys tools random: url OK');
    assert.equal('POST', method, 'sys tools random: method OK');
  });
});
