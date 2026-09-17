/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { getPreference } from 'vault/utils/preferences';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | user-preferences', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.flags = this.owner.lookup('service:flags');
    window.localStorage.clear();
  });

  module('data-privacy HVD visibility', function () {
    test('it renders the Data & Privacy section for non-HVD clusters', async function (assert) {
      // Verifies the section is visible when the cluster is not HVD-managed.
      this.set('model', { isHvdManaged: false });
      await render(hbs`
        {{#unless this.model.isHvdManaged}}
          <UserPreferences::DataPrivacy />
        {{/unless}}
      `);

      assert
        .dom('[data-test-data-privacy-section]')
        .exists('Data & Privacy section renders for non-HVD clusters');
    });

    test('it hides the Data & Privacy section for HVD-managed clusters', async function (assert) {
      // Verifies the section is suppressed when the cluster is HVD-managed, because
      // HVD telemetry uses PostHog and the Segment toggle is meaningless for those users.
      this.set('model', { isHvdManaged: true });
      await render(hbs`
        {{#unless this.model.isHvdManaged}}
          <UserPreferences::DataPrivacy />
        {{/unless}}
      `);

      assert
        .dom('[data-test-data-privacy-section]')
        .doesNotExist('Data & Privacy section is hidden for HVD-managed clusters');
    });
  });

  module('data-privacy', function () {
    test('it renders the Share usage section card from HDS components', async function (assert) {
      await render(hbs`<UserPreferences::DataPrivacy />`);

      assert.dom('[data-test-data-privacy-section]').exists('the Data & Privacy section renders');
      assert.dom(GENERAL.toggleInput('telemetry-consent')).exists('the consent toggle renders');
      assert.dom('[data-test-data-privacy-included] li').exists({ count: 3 }, 'included items render');
      assert.dom('[data-test-data-privacy-excluded] li').exists({ count: 3 }, 'excluded items render');
    });

    test('the toggle defaults off when no consent value is stored', async function (assert) {
      await render(hbs`<UserPreferences::DataPrivacy />`);

      assert.dom(GENERAL.toggleInput('telemetry-consent')).isNotChecked('toggle is off by default (opt-in)');
    });

    test('toggling on persists consent through the registry', async function (assert) {
      await render(hbs`<UserPreferences::DataPrivacy />`);

      await click(GENERAL.toggleInput('telemetry-consent'));

      assert.dom(GENERAL.toggleInput('telemetry-consent')).isChecked('toggle reflects the on state');
      assert.true(getPreference('telemetryConsent'), 'consent is persisted via the registry');
      assert.strictEqual(
        window.localStorage.getItem('vault:prefs:telemetryConsent'),
        'true',
        'value is written under the namespaced registry key'
      );
    });

    test('toggling off writes the off value through the registry', async function (assert) {
      await render(hbs`<UserPreferences::DataPrivacy />`);

      await click(GENERAL.toggleInput('telemetry-consent'));
      await click(GENERAL.toggleInput('telemetry-consent'));

      assert.dom(GENERAL.toggleInput('telemetry-consent')).isNotChecked('toggle is off again');
      assert.false(getPreference('telemetryConsent'), 'consent persists as off');
    });
  });
});
