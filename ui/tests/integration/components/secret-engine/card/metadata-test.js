/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { keyMgmtMockModel } from 'vault/tests/helpers/secret-engine/mocks';

module('Integration | Component | SecretEngine::Card::Metadata', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = keyMgmtMockModel;
  });

  test('it shows metadata card information', async function (assert) {
    assert.expect(4);
    await render(hbs`
      <SecretEngine::Card::Metadata @model={{this.model}} />
    `);
    assert.dom(`${GENERAL.cardContainer('metadata')} h2`).hasText('Metadata');
    assert.dom(GENERAL.inputByAttr('path')).hasValue(this.model.secretsEngine.path);
    assert.dom(GENERAL.inputByAttr('accessor')).hasValue(this.model.secretsEngine.accessor);
    await fillIn(GENERAL.textareaByAttr('description'), 'Some awesome description');
    assert.dom(GENERAL.textareaByAttr('description')).hasValue('Some awesome description');
  });
});
