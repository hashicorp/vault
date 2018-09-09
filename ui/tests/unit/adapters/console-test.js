import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | console', function(hooks) {
  setupTest(hooks);

  test('it builds the correct URL', function(assert) {
    let adapter = this.owner.lookup('adapter:console');
    let sysPath = 'sys/health';
    let awsPath = 'aws/roles/my-other-role';
    assert.equal(adapter.buildURL(sysPath), '/v1/sys/health');
    assert.equal(adapter.buildURL(awsPath), '/v1/aws/roles/my-other-role');
  });
});
