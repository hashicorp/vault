import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Service | user-preference', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.service = this.owner.lookup('service:user-preference');
    this.service.kvDisplaySetting = 'unset';
  });
  hooks.afterEach(function () {
    this.service.kvDisplaySetting = 'unset';
  });

  test('it manages the kvDisplaySetting correctly', function (assert) {
    assert.strictEqual(this.service.kvDisplaySetting, 'unset', 'kvDisplaySetting is unset by default');
    assert.false(
      this.service.calculateInitialKvJson(false),
      'inital view is not json when secret not advanced'
    );
    assert.true(this.service.calculateInitialKvJson(true), 'inital view is json when secret is advanced');
    this.service.setKvDisplayPreference(true);
    assert.strictEqual(this.service.kvDisplaySetting, 'json', 'kvDisplaySetting is set to json');
    assert.true(this.service.calculateInitialKvJson(false), 'inital view is json due to user preference');
    assert.true(this.service.calculateInitialKvJson(true), 'inital view is json due to user preference');
    this.service.setKvDisplayPreference(false);
    assert.strictEqual(this.service.kvDisplaySetting, 'keyvalue', 'kvDisplaySetting is set to keyvalue');
    assert.false(
      this.service.calculateInitialKvJson(false),
      'inital view is key-value due to user preference'
    );
    assert.false(
      this.service.calculateInitialKvJson(true),
      'inital view is key-value due to user preference'
    );
  });
});
