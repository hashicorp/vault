/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { keyMgmtMockModel } from 'vault/tests/helpers/secret-engine/mocks';

module('Integration | Component | SecretEngine::Card::Version', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = keyMgmtMockModel;
  });

  test('it shows version card information', async function (assert) {
    assert.expect(4);
    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);
    assert.dom(`${GENERAL.cardContainer('version')} h2`).hasText('Version');
    assert.dom(GENERAL.infoRowValue('type')).hasAnyText(keyMgmtMockModel.secretsEngine.type);
    assert
      .dom(GENERAL.infoRowValue('running_plugin_version'))
      .hasAnyText(keyMgmtMockModel.secretsEngine.running_plugin_version);
    assert.dom(GENERAL.inputByAttr('plugin_version')).doesNotExist();
  });
});
