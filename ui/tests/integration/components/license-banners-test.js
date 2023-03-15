/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import subDays from 'date-fns/subDays';
import addDays from 'date-fns/addDays';
import formatRFC3339 from 'date-fns/formatRFC3339';
import timestamp from 'core/utils/timestamp';

module('Integration | Component | license-banners', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const mockNow = timestamp.now();
    this.now = mockNow;
    this.yesterday = subDays(mockNow, 1);
    this.nextMonth = addDays(mockNow, 30);
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
  });

  test('it does not render if no expiry', async function (assert) {
    assert.expect(1);
    await render(hbs`<LicenseBanners />`);
    assert.dom('[data-test-license-banner]').doesNotExist('License banner does not render');
  });

  test('it renders an error if expiry is before now', async function (assert) {
    assert.expect(2);
    this.set('expiry', formatRFC3339(this.yesterday));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert.dom('[data-test-license-banner-expired]').exists('Expired license banner renders');
    assert.dom('.message-title').hasText('License expired', 'Shows correct title on alert');
  });

  test('it renders a warning if expiry is within 30 days', async function (assert) {
    assert.expect(2);
    this.set('expiry', formatRFC3339(this.nextMonth));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert.dom('[data-test-license-banner-warning]').exists('Warning license banner renders');
    assert.dom('.message-title').hasText('Vault license expiring', 'Shows correct title on alert');
  });

  test('it does not render a banner if expiry is outside 30 days', async function (assert) {
    assert.expect(1);
    const outside30 = addDays(this.mockNow, 32);
    this.set('expiry', formatRFC3339(outside30));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert.dom('[data-test-license-banner]').doesNotExist('License banner does not render');
  });

  test('it does not render the expired banner if it has been dismissed', async function (assert) {
    assert.expect(3);
    this.set('expiry', formatRFC3339(this.yesterday));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    await click('[data-test-dismiss-expired]');
    assert.dom('[data-test-license-banner-expired]').doesNotExist('Expired license banner does not render');

    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    const localStorageResult = JSON.parse(localStorage.getItem(`dismiss-license-banner-1.13.1+ent`));
    assert.strictEqual(localStorageResult, 'expired');
    assert
      .dom('[data-test-license-banner-expired]')
      .doesNotExist('The expired banner still does not render after a re-render.');
    localStorage.removeItem(`dismiss-license-banner-1.13.1+ent`);
  });

  test('it does not render the warning banner if it has been dismissed', async function (assert) {
    assert.expect(3);
    this.set('expiry', formatRFC3339(this.nextMonth));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    await click('[data-test-dismiss-warning]');
    assert.dom('[data-test-license-banner-warning]').doesNotExist('Warning license banner does not render');

    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    const localStorageResult = JSON.parse(localStorage.getItem(`dismiss-license-banner-1.13.1+ent`));
    assert.strictEqual(localStorageResult, 'warning');
    assert
      .dom('[data-test-license-banner-warning]')
      .doesNotExist('The warning banner still does not render after a re-render.');
    localStorage.removeItem(`dismiss-license-banner-1.13.1+ent`);
  });

  test('it renders a banner if the vault license has changed', async function (assert) {
    assert.expect(3);
    this.version.version = '1.12.1+ent';
    this.set('expiry', formatRFC3339(this.nextMonth));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    await click('[data-test-dismiss-warning]');
    this.version.version = '1.13.1+ent';
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert
      .dom('[data-test-license-banner-warning]')
      .exists('The warning banner shows even though we have dismissed it earlier.');

    await click('[data-test-dismiss-warning]');
    const localStorageResultNewVersion = JSON.parse(
      localStorage.getItem(`dismiss-license-banner-1.13.1+ent`)
    );
    const localStorageResultOldVersion = JSON.parse(
      localStorage.getItem(`dismiss-license-banner-1.12.1+ent`)
    );
    // Check that localStorage was cleaned and no longer contains the old version storage key.
    assert.strictEqual(localStorageResultOldVersion, null);
    assert.strictEqual(localStorageResultNewVersion, 'warning');
    // If debugging this test remember to clear localStorage if the test was not run to completion.
    localStorage.removeItem(`dismiss-license-banner-1.13.1+ent`);
  });
});
