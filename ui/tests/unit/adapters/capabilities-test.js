import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | capabilities', function(hooks) {
  setupTest(hooks);

  test('calls the correct url', function(assert) {
    let url, method, options;
    let adapter = this.owner.factoryFor('adapter:capabilities').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    adapter.findRecord(null, 'capabilities', 'foo');
    assert.equal('/v1/sys/capabilities-self', url, 'calls the correct URL');
    assert.deepEqual({ paths: ['foo'] }, options.data, 'data params OK');
    assert.equal('POST', method, 'method OK');
  });
});
