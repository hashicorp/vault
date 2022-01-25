import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Transform | stringarray', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    let transform = this.owner.lookup('transform:stringarray');
    assert.ok(transform);
  });
});
