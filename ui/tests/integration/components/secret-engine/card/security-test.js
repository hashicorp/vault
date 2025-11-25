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
    assert.dom(GENERAL.fieldLabel('local')).hasText('Local');
    assert
      .dom(GENERAL.helpTextByAttr('local'))
      .hasText('Secrets stay in one cluster and are not replicated.');
    assert.dom(GENERAL.inputByAttr('local')).isChecked();
    assert.dom(GENERAL.fieldLabel('seal_wrap')).hasText('Seal wrap');
    assert
      .dom(GENERAL.helpTextByAttr('seal_wrap'))
      .hasText('Wrap secrets with an additional encryption layer using a seal.');
    assert.dom(GENERAL.inputByAttr('seal_wrap')).isNotChecked();
  });
});
