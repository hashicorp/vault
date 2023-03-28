import timestamp from 'core/utils/timestamp';
import sinon from 'sinon';
import { module, test } from 'qunit';

/*
  This test coverage is more for an example than actually covering the utility
*/
module('Unit | Utility | timestamp', function () {
  test('it can be overridden', function (assert) {
    const stub = sinon.stub(timestamp, 'now').callsFake(() => new Date('2030-03-03T03:30:03'));
    const result = timestamp.now();
    assert.strictEqual(result.toISOString(), new Date('2030-03-03T03:30:03').toISOString());
    // Always make sure to restore the stub
    stub.restore(); // timestamp.now.restore(); also works
  });
});
