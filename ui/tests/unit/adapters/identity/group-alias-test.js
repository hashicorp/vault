import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import testCases from './_test-cases';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

module('Unit | Adapter | identity/group-alias', function(hooks) {
  setupTest(hooks);

  hooks.beforeEach(function() {
    this.server = apiStub();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  const cases = testCases('identity/group-alias');

  cases.forEach(testCase => {
    test(`group-alias#${testCase.adapterMethod}`, function(assert) {
      assert.expect(2);
      let adapter = this.owner.lookup('adapter:identity/group-alias');
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
});
