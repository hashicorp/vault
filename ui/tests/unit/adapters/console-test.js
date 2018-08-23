import { moduleFor, test } from 'ember-qunit';
import needs from 'vault/tests/unit/adapters/_adapter-needs';

moduleFor('adapter:console', 'Unit | Adapter | console', {
  needs,
});

test('it builds the correct URL', function(assert) {
  let adapter = this.subject();
  let sysPath = 'sys/health';
  let awsPath = 'aws/roles/my-other-role';
  assert.equal(adapter.buildURL(sysPath), '/v1/sys/health');
  assert.equal(adapter.buildURL(awsPath), '/v1/aws/roles/my-other-role');
});
