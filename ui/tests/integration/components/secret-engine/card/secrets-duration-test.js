/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { keyMgmtMockModel } from 'vault/tests/helpers/secret-engine/mocks';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | SecretEngine::Card::LeaseDuration', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = keyMgmtMockModel;
  });

  test('it shows default and max ttl pickers', async function (assert) {
    await render(hbs`
      <SecretEngine::Card::SecretsDuration @model={{this.model}} />
    `);
    assert.dom(GENERAL.inputByAttr('default_lease_ttl')).exists();
    assert.dom(GENERAL.inputByAttr('max_lease_ttl')).exists();
    assert.dom(GENERAL.fieldLabel('default_lease_ttl')).hasText('Default time-to-live (TTL)');
    assert.dom(GENERAL.fieldLabel('max_lease_ttl')).hasText('Maximum time-to-live (TTL)');
    assert
      .dom(GENERAL.helpTextByAttr('default_lease_ttl'))
      .hasText('How long secrets in this engine stay valid.');
    assert
      .dom(GENERAL.helpTextByAttr('max_lease_ttl'))
      .hasText('Maximum extension for the secrets life beyond default.');
  });
});
