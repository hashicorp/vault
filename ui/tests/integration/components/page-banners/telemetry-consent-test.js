/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

const SELECTOR = {
  banner: '[data-test-telemetry-consent-banner]',
  accept: '[data-test-telemetry-consent-accept]',
  decline: '[data-test-telemetry-consent-decline]',
};

module('Integration | Component | page-banners/telemetry-consent', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.analytics = this.owner.lookup('service:analytics');
    this.recordConsent = sinon.stub(this.analytics, 'recordConsent');
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('it renders the consent banner', async function (assert) {
    await render(hbs`<PageBanners::TelemetryConsent />`);
    assert.dom(SELECTOR.banner).exists('the banner renders');
    assert.dom(SELECTOR.accept).exists('the Accept action renders');
    assert.dom(SELECTOR.decline).exists('the Decline action renders');
  });

  test('clicking Accept records consent as accepted', async function (assert) {
    await render(hbs`<PageBanners::TelemetryConsent />`);

    await click(SELECTOR.accept);

    assert.true(this.recordConsent.calledOnceWith(true), 'Consent is recorded as true in analytics service');
  });

  test('clicking Decline records consent as declined', async function (assert) {
    await render(hbs`<PageBanners::TelemetryConsent />`);

    await click(SELECTOR.decline);

    assert.true(
      this.recordConsent.calledOnceWith(false),
      'Consent is recorded as false in analytics service'
    );
  });
});
