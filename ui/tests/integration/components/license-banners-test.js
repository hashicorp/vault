/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
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
    sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2018-04-03T14:15:30')));
    const mockNow = timestamp.now();
    this.yesterday = subDays(mockNow, 1);
    this.nextMonth = addDays(mockNow, 30);
    this.outside30 = addDays(mockNow, 32);
    this.tomorrow = addDays(mockNow, 1);
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
    this.version.type = 'enterprise';
  });

  test('it does not render if no expiry', async function (assert) {
    assert.expect(2);
    await render(hbs`<LicenseBanners />`);
    assert.dom('[data-test-license-banner-expired]').doesNotExist();
    assert.dom('[data-test-license-banner-warning]').doesNotExist();
  });

  test('it renders an error if expiry is before now', async function (assert) {
    assert.expect(2);
    this.set('expiry', formatRFC3339(this.yesterday));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert.dom('[data-test-license-banner-expired]').exists('Expired license banner renders');
    assert
      .dom('[data-test-license-banner-expired] .hds-alert__title')
      .hasText('License expired', 'Shows correct title on alert');
  });

  test('it renders a warning if expiry is within 30 days', async function (assert) {
    assert.expect(2);
    this.set('expiry', formatRFC3339(this.nextMonth));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert.dom('[data-test-license-banner-warning]').exists('Warning license banner renders');
    assert
      .dom('[data-test-license-banner-warning] .hds-alert__title')
      .hasText('Vault license expiring', 'Shows correct title on alert');
  });

  test('it does not render a banner if expiry is outside 30 days', async function (assert) {
    assert.expect(2);
    this.set('expiry', formatRFC3339(this.outside30));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert.dom('[data-test-license-banner-expired]').doesNotExist();
    assert.dom('[data-test-license-banner-warning]').doesNotExist();
  });

  test('it does not render the expired banner if it has been dismissed', async function (assert) {
    assert.expect(3);
    this.set('expiry', formatRFC3339(this.yesterday));
    const key = `dismiss-license-banner-${this.version.version}-${this.expiry}`;
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    await click('[data-test-license-banner-expired] [data-test-icon="x"]');
    assert.dom('[data-test-license-banner-expired]').doesNotExist('Expired license banner does not render');

    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    const localStorageResult = JSON.parse(localStorage.getItem(key));
    assert.strictEqual(localStorageResult, 'expired');
    assert
      .dom('[data-test-license-banner-expired]')
      .doesNotExist('The expired banner still does not render after a re-render.');
    localStorage.removeItem(key);
  });

  test('it does not render the warning banner if it has been dismissed', async function (assert) {
    assert.expect(3);
    this.set('expiry', formatRFC3339(this.nextMonth));
    const key = `dismiss-license-banner-${this.version.version}-${this.expiry}`;
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    await click('[data-test-license-banner-warning] [data-test-icon="x"]');
    assert.dom('[data-test-license-banner-warning]').doesNotExist('Warning license banner does not render');

    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    const localStorageResult = JSON.parse(localStorage.getItem(key));
    assert.strictEqual(localStorageResult, 'warning');
    assert
      .dom('[data-test-license-banner-warning]')
      .doesNotExist('The warning banner still does not render after a re-render.');
    localStorage.removeItem(key);
  });

  test('it renders a banner if the vault license has changed', async function (assert) {
    assert.expect(3);
    this.version.version = '1.12.1+ent';
    this.set('expiry', formatRFC3339(this.nextMonth));
    const keyOldVersion = `dismiss-license-banner-${this.version.version}-${this.expiry}`;
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    await click('[data-test-license-banner-warning] [data-test-icon="x"]');
    this.version.version = '1.13.1+ent';
    const keyNewVersion = `dismiss-license-banner-${this.version.version}-${this.expiry}`;
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert
      .dom('[data-test-license-banner-warning]')
      .exists('The warning banner shows even though we have dismissed it earlier.');

    await click('[data-test-license-banner-warning] [data-test-icon="x"]');
    const localStorageResultNewVersion = JSON.parse(localStorage.getItem(keyNewVersion));
    const localStorageResultOldVersion = JSON.parse(localStorage.getItem(keyOldVersion));
    // Check that localStorage was cleaned and no longer contains the old version storage key.
    assert.strictEqual(localStorageResultOldVersion, null, 'local storage was cleared for the old version');
    assert.strictEqual(
      localStorageResultNewVersion,
      'warning',
      'local storage holds the new version with a warning'
    );
    // If debugging this test remember to clear localStorage if the test was not run to completion.
    localStorage.removeItem(keyNewVersion);
  });

  test('it renders a banner if the vault expiry has changed', async function (assert) {
    assert.expect(3);
    this.set('expiry', formatRFC3339(this.tomorrow));
    const keyOldExpiry = `dismiss-license-banner-${this.version.version}-${this.expiry}`;
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    await click('[data-test-license-banner-warning] [data-test-icon="x"]');
    this.set('expiry', formatRFC3339(this.nextMonth));
    const keyNewExpiry = `dismiss-license-banner-${this.version.version}-${this.expiry}`;
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert
      .dom('[data-test-license-banner-warning]')
      .exists('The warning banner shows even though we have dismissed it earlier.');

    await click('[data-test-license-banner-warning] [data-test-icon="x"]');
    const localStorageResultNewExpiry = JSON.parse(localStorage.getItem(keyNewExpiry));
    const localStorageResultOldExpiry = JSON.parse(localStorage.getItem(keyOldExpiry));
    // Check that localStorage was cleaned and no longer contains the old version storage key.
    assert.strictEqual(localStorageResultOldExpiry, null, 'local storage was cleared for the old expiry');
    assert.strictEqual(
      localStorageResultNewExpiry,
      'warning',
      'local storage holds the new expiry with a warning'
    );
    // If debugging this test remember to clear localStorage if the test was not run to completion.
    localStorage.removeItem(keyNewExpiry);
  });
});
