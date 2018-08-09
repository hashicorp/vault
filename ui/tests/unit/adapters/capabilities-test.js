import { moduleFor, test } from 'ember-qunit';
import Ember from 'ember';

moduleFor('adapter:capabilities', 'Unit | Adapter | capabilities', {
  needs: ['service:auth', 'service:flash-messages', 'service:control-group', 'service:version'],
});

test('calls the correct url', function(assert) {
  let url, method, options;
  let adapter = this.subject({
    ajax: (...args) => {
      [url, method, options] = args;
      return Ember.RSVP.resolve();
    },
  });

  adapter.findRecord(null, 'capabilities', 'foo');
  assert.equal('/v1/sys/capabilities-self', url, 'calls the correct URL');
  assert.deepEqual({ paths: ['foo'] }, options.data, 'data params OK');
  assert.equal('POST', method, 'method OK');
});
