import Pretender from 'pretender';
import { moduleFor, test } from 'ember-qunit';
import { storeMVP } from './_test-cases';

moduleFor('adapter:identity/entity-merge', 'Unit | Adapter | identity/entity-merge', {
  needs: ['service:auth', 'service:flash-messages', 'service:control-group', 'service:version'],
  beforeEach() {
    this.server = new Pretender(function() {
      this.post('/v1/**', response => {
        return [response, { 'Content-Type': 'application/json' }, JSON.stringify({})];
      });
    });
  },
  afterEach() {
    this.server.shutdown();
  },
});

test(`entity-merge#createRecord`, function(assert) {
  assert.expect(2);
  let adapter = this.subject();
  adapter.createRecord(storeMVP, { modelName: 'identity/entity-merge' }, { attr: x => x });
  let { url, method } = this.server.handledRequests[0];
  assert.equal(url, `/v1/identity/entity/merge`, ` calls the correct url`);
  assert.equal(method, 'POST', `uses the correct http verb: POST`);
});
