/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const { overviewCard } = GENERAL;
module('Integration | Component | kv | kv-subkeys-card', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  hooks.beforeEach(function () {
    this.isPatchAllowed = true;
    this.subkeys = {
      foo: null,
      bar: {
        baz: null,
      },
    };
    this.renderComponent = async () => {
      return render(
        hbs`<KvSubkeysCard @subkeys={{this.subkeys}} @isPatchAllowed={{this.isPatchAllowed}} />`,
        {
          owner: this.engine,
        }
      );
    };
  });

  test('it renders', async function (assert) {
    await this.renderComponent();

    assert.dom(overviewCard.title('Subkeys')).exists();
    assert
      .dom(overviewCard.description('Subkeys'))
      .hasText(
        'The table is displaying the top level subkeys. Toggle on the JSON view to see the full depth.'
      );
    assert.dom(overviewCard.content('Subkeys')).hasText('Keys foo bar');
    assert.dom(GENERAL.toggleInput('kv-subkeys')).isNotChecked('JSON toggle is not checked by default');
    assert.dom(overviewCard.actionText('Patch secret')).exists();
  });

  test('it hides patch action when isPatchAllowed is false', async function (assert) {
    this.isPatchAllowed = false;
    await this.renderComponent();
    assert.dom(overviewCard.title('Subkeys')).exists();
    assert.dom(overviewCard.actionText('Patch secret')).doesNotExist();
  });

  test('it toggles to JSON', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.toggleInput('kv-subkeys')).isNotChecked();
    await click(GENERAL.toggleInput('kv-subkeys'));
    assert.dom(GENERAL.toggleInput('kv-subkeys')).isChecked('JSON toggle is checked');
    assert.dom(overviewCard.description('Subkeys')).hasText(
      'These are the subkeys within this secret. All underlying values of leaf keys are not retrieved and are replaced with null instead. Subkey API documentation .' // space is intentional because a trailing icon renders after the inline link
    );
    assert.dom(overviewCard.content('Subkeys')).hasText(JSON.stringify(this.subkeys, null, 2));
  });
});
