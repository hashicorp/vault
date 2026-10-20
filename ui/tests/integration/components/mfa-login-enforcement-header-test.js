/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | mfa-login-enforcement-header', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
      },
    });
  });

  test('it renders heading', async function (assert) {
    await render(hbs`<Mfa::MfaLoginEnforcementHeader @heading="New enforcement" />`);

    assert.dom(GENERAL.hdsPageHeaderTitle).includesText('New enforcement');
    assert.dom(GENERAL.icon('lock')).hasClass('hds-icon-lock', 'Lock icon renders');
    assert
      .dom(GENERAL.hdsPageHeaderDescription)
      .includesText('An enforcement will define which auth types', 'Description renders');
    assert.dom('[data-test-mleh-radio]').doesNotExist('Radio cards are hidden when not inline display mode');
    assert
      .dom('[data-test-component="search-select"]')
      .doesNotExist('Search select is hidden when not inline display mode');
  });

  test('it renders inline', async function (assert) {
    this.server.get('/identity/mfa/login-enforcement', () => {
      assert.ok(true, 'Request made to fetch enforcements');
      return {
        data: {
          key_info: {
            foo: { name: 'foo' },
          },
          keys: ['foo'],
        },
      };
    });

    await render(hbs`
      <Mfa::MfaLoginEnforcementHeader
        @isInline={{true}}
        @radioCardGroupValue={{this.value}}
        @onRadioCardSelect={{fn (mut this.value)}}
        @onEnforcementSelect={{fn (mut this.enforcement)}}
      />
    `);

    assert.dom('h2').hasText('Enforcement', 'Inline heading renders');
    assert.dom('p').includesText('An enforcement includes the authentication types', 'Description renders');
    assert.dom('fieldset legend').hasText('Enforcement', 'The radio group has an accessible legend');
    for (const option of ['new', 'existing', 'skip']) {
      await click(`[data-test-mleh-radio="${option}"]`);
      assert.strictEqual(this.value, option, 'Value is updated on radio select');
      assert.dom(`[data-test-mleh-radio="${option}"]`).isChecked('The selected radio is checked');
      assert
        .dom(`[data-test-mleh-radio="${option}"]`)
        .hasAttribute('name', 'mfa-enforcement-preference', 'Radios share the group name');
      if (option === 'existing') {
        assert.dom('[data-test-component="search-select"]').exists('Existing enforcement selector renders');
        await clickTrigger();
        await click('.ember-power-select-option');
      } else {
        assert
          .dom('[data-test-component="search-select"]')
          .doesNotExist('Existing enforcement selector is hidden for other preferences');
      }
    }

    assert.strictEqual(this.enforcement.name, 'foo', 'Existing enforcement is selected');
  });

  test('it preserves controlled selection and disables existing when no enforcements are available', async function (assert) {
    this.server.get('/identity/mfa/login-enforcement', () => ({ data: { keys: [], key_info: {} } }));
    this.value = 'skip';

    await render(hbs`
      <Mfa::MfaLoginEnforcementHeader
        @isInline={{true}}
        @radioCardGroupValue={{this.value}}
        @onRadioCardSelect={{fn (mut this.value)}}
        @onEnforcementSelect={{fn (mut this.enforcement)}}
      />
    `);

    assert.dom('[data-test-mleh-radio="skip"]').isChecked('The initial preference is selected');
    assert.dom('[data-test-mleh-radio="existing"]').isDisabled('Existing requires an available enforcement');
    assert.dom('[data-test-mleh-radio="new"]').isNotDisabled('Creating a new enforcement remains available');
    assert.dom('[data-test-mleh-radio="skip"]').isNotDisabled('Skipping enforcement remains available');

    this.set('value', 'new');
    await settled();
    assert.dom('[data-test-mleh-radio="new"]').isChecked('The group follows parent selection changes');
    assert.dom('[data-test-mleh-radio="skip"]').isNotChecked('The previous preference is unchecked');

    await click('[data-test-mleh-radio="skip"]');
    assert.strictEqual(this.value, 'skip', 'Selecting another card still invokes the parent callback');
    assert.dom('[data-test-mleh-radio="new"]').isNotChecked('Only one preference remains checked');
  });
});
