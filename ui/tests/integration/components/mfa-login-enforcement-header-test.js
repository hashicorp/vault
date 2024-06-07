/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | mfa-login-enforcement-header', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    setRunOptions({
      rules: {
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
      },
    });
  });

  test('it renders heading', async function (assert) {
    await render(hbs`<Mfa::MfaLoginEnforcementHeader @heading="New enforcement" />`);

    assert.dom('[data-test-mleh-title]').includesText('New enforcement');
    assert.dom('[data-test-mleh-title] svg').hasClass('flight-icon-lock', 'Lock icon renders');
    assert
      .dom('[data-test-mleh-description]')
      .includesText('An enforcement will define which auth types', 'Description renders');
    assert.dom('[data-test-mleh-radio]').doesNotExist('Radio cards are hidden when not inline display mode');
    assert
      .dom('[data-test-component="search-select"]')
      .doesNotExist('Search select is hidden when not inline display mode');
  });

  test('it renders inline', async function (assert) {
    assert.expect(7);

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

    assert.dom('[data-test-mleh-title]').includesText('Enforcement');
    assert
      .dom('[data-test-mleh-description]')
      .includesText('An enforcement includes the authentication types', 'Description renders');
    for (const option of ['new', 'existing', 'skip']) {
      await click(`[data-test-mleh-radio="${option}"] input`);
      assert.strictEqual(this.value, option, 'Value is updated on radio select');
      if (option === 'existing') {
        await clickTrigger();
        await click('.ember-power-select-option');
      }
    }

    assert.strictEqual(this.enforcement.name, 'foo', 'Existing enforcement is selected');
  });
});
