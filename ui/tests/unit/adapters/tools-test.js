import { moduleFor, test } from 'ember-qunit';
import Ember from 'ember';

moduleFor('adapter:tools', 'Unit | Adapter | tools', {
  needs: ['service:auth', 'service:flash-messages', 'service:control-group', 'service:version'],
});

test('wrapping api urls', function(assert) {
  let url, method, options;
  let adapter = this.subject({
    ajax: (...args) => {
      [url, method, options] = args;
      return Ember.RSVP.resolve();
    },
  });

  let data = { foo: 'bar' };
  adapter.toolAction('wrap', data, { wrapTTL: '30m' });
  assert.equal('/v1/sys/wrapping/wrap', url, 'wrapping:wrap url OK');
  assert.equal('POST', method, 'wrapping:wrap method OK');
  assert.deepEqual({ data: data, wrapTTL: '30m' }, options, 'wrapping:wrap options OK');

  data = { token: 'token' };
  adapter.toolAction('lookup', data);
  assert.equal('/v1/sys/wrapping/lookup', url, 'wrapping:lookup url OK');
  assert.equal('POST', method, 'wrapping:lookup method OK');
  assert.deepEqual({ data }, options, 'wrapping:lookup options OK');

  adapter.toolAction('unwrap', data);
  assert.equal('/v1/sys/wrapping/unwrap', url, 'wrapping:unwrap url OK');
  assert.equal('POST', method, 'wrapping:unwrap method OK');
  assert.deepEqual({ data }, options, 'wrapping:unwrap options OK');

  adapter.toolAction('rewrap', data);
  assert.equal('/v1/sys/wrapping/rewrap', url, 'wrapping:rewrap url OK');
  assert.equal('POST', method, 'wrapping:rewrap method OK');
  assert.deepEqual({ data }, options, 'wrapping:rewrap options OK');
});

test('tools api urls', function(assert) {
  let url, method;
  let adapter = this.subject({
    ajax: (...args) => {
      [url, method] = args;
      return Ember.RSVP.resolve();
    },
  });

  adapter.toolAction('hash', { input: 'someBase64' });
  assert.equal(url, '/v1/sys/tools/hash', 'sys tools hash: url OK');
  assert.equal('POST', method, 'sys tools hash: method OK');

  adapter.toolAction('random', { bytes: '32' });
  assert.equal(url, '/v1/sys/tools/random', 'sys tools random: url OK');
  assert.equal('POST', method, 'sys tools random: method OK');
});
