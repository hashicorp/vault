/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { keyMgmtMockModel } from 'vault/tests/helpers/secret-engine/mocks';
import { SELECTORS } from 'vault/tests/helpers/secret-engine/general-settings-selectors';

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
    assert.dom(SELECTORS.engineType).hasAnyText(keyMgmtMockModel.secretsEngine.type);
    assert.dom(SELECTORS.currentVersion).hasAnyText(keyMgmtMockModel.secretsEngine.running_plugin_version);
    assert.dom(SELECTORS.versionCard.versionsDropdown).doesNotExist();
  });
});
