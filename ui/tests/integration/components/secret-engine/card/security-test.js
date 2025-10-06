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

module('Integration | Component | SecretEngine::Card::Security', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = keyMgmtMockModel;
  });

  test('it shows security card information', async function (assert) {
    assert.expect(7);
    await render(hbs`
      <SecretEngine::Card::Security @model={{this.model}} />
    `);
    assert.dom(`${GENERAL.cardContainer('security')} h2`).hasText('Security');
    assert.dom(SELECTORS.label('local')).hasText('Local');
    assert.dom(SELECTORS.helperText('local')).hasText('Secrets stay in one cluster and are not replicated.');
    assert.dom(GENERAL.inputByAttr('local')).isChecked();
    assert.dom(SELECTORS.label('seal_wrap')).hasText('Seal wrap');
    assert
      .dom(SELECTORS.helperText('seal_wrap'))
      .hasText('Wrap secrets with an additional encryption layer using a seal.');
    assert.dom(GENERAL.inputByAttr('seal_wrap')).isNotChecked();
  });
});
