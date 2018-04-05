import Pretender from 'pretender';
import { moduleFor, test } from 'ember-qunit';
import testCases from './_test-cases';

const noop = response => {
  return function() {
    return [response, { 'Content-Type': 'application/json' }, JSON.stringify({})];
  };
};

moduleFor('adapter:identity/group', 'Unit | Adapter | identity/group', {
  needs: ['service:auth', 'service:flash-messages'],
  beforeEach() {
    this.server = new Pretender(function() {
      this.post('/v1/**', noop());
      this.put('/v1/**', noop());
      this.get('/v1/**', noop());
      this.delete('/v1/**', noop(204));
    });
  },
  afterEach() {
    this.server.shutdown();
  },
});

const cases = testCases('identit/entity');

cases.forEach(testCase => {
  test(`group#${testCase.adapterMethod}`, function(assert) {
    assert.expect(2);
    let adapter = this.subject();
    adapter[testCase.adapterMethod](...testCase.args);
    let { url, method } = this.server.handledRequests[0];
    assert.equal(url, testCase.url, `${testCase.adapterMethod} calls the correct url: ${testCase.url}`);
    assert.equal(
      method,
      testCase.method,
      `${testCase.adapterMethod} uses the correct http verb: ${testCase.method}`
    );
  });
});
