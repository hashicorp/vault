import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import subDays from 'date-fns/subDays';
import addDays from 'date-fns/addDays';
import formatRFC3339 from 'date-fns/formatRFC3339';

module('Integration | Component | license-banners', function (hooks) {
  setupRenderingTest(hooks);

  test('it does not render if no expiry', async function (assert) {
    await render(hbs`<LicenseBanners />`);
    assert.dom('[data-test-license-banner]').doesNotExist('License banner does not render');
  });

  test('it renders an error if expiry is before now', async function (assert) {
    const yesterday = subDays(new Date(), 1);
    this.set('expiry', formatRFC3339(yesterday));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert.dom('[data-test-license-banner-expired]').exists('Expired license banner renders');
    assert.dom('.message-title').hasText('License expired', 'Shows correct title on alert');
  });

  test('it renders a warning if expiry is within 30 days', async function (assert) {
    const nextMonth = addDays(new Date(), 30);
    this.set('expiry', formatRFC3339(nextMonth));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert.dom('[data-test-license-banner-warning]').exists('Warning license banner renders');
    assert.dom('.message-title').hasText('Vault license expiring', 'Shows correct title on alert');
  });

  test('it does not render a banner if expiry is outside 30 days', async function (assert) {
    const outside30 = addDays(new Date(), 32);
    this.set('expiry', formatRFC3339(outside30));
    await render(hbs`<LicenseBanners @expiry={{this.expiry}} />`);
    assert.dom('[data-test-license-banner]').doesNotExist('License banner does not render');
  });
});
