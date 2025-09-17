/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { keyMgmtMockModel } from 'vault/tests/helpers/secret-engine/mocks';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SELECTORS } from 'vault/tests/helpers/secret-engine/general-settings-selectors';

module('Integration | Component | SecretEngine::Card::LeaseDuration', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = keyMgmtMockModel;
  });

  test('it shows default and max ttl pickers', async function (assert) {
    assert.expect(5);
    await render(hbs`
      <SecretEngine::Card::LeaseDuration @model={{this.model}} />
    `);
    assert.dom(SELECTORS.ttlPickerV2).exists({ count: 2 });
    assert.dom(GENERAL.fieldLabelbyAttr('default_lease_ttl')).hasText('Time-to-live (TTL)');
    assert.dom(GENERAL.fieldLabelbyAttr('max_lease_ttl')).hasText('Maximum Time-to-live (TTL)');
    assert.dom(GENERAL.helpTextByAttr('default_lease_ttl')).hasText('Standard expiry deadline.');
    assert.dom(GENERAL.helpTextByAttr('max_lease_ttl')).hasText('Maximum possible extension for expiry.');
  });
});
