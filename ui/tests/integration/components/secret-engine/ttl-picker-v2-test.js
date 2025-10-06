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

module('Integration | Component | SecretEngine::TtlPickerV2', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = keyMgmtMockModel;
  });

  test('it shows default ttl picker', async function (assert) {
    assert.expect(4);
    this.ttlKey = 'default_lease_ttl';
    await render(hbs`
      <SecretEngine::TtlPickerV2 @model={{this.model}} @ttlKey={{this.ttlKey}} />
    `);
    assert.dom(GENERAL.fieldLabelbyAttr(this.ttlKey)).hasText('Default time-to-live (TTL)');
    assert.dom(GENERAL.helpTextByAttr(this.ttlKey)).hasText('How long secrets in this engine stay valid.');
    await fillIn(GENERAL.inputByAttr(this.ttlKey), 5);
    await fillIn(GENERAL.selectByAttr(this.ttlKey), 'm');
    assert.dom(GENERAL.inputByAttr(this.ttlKey)).hasValue('5');
    assert.dom(GENERAL.selectByAttr(this.ttlKey)).hasValue('m');
  });

  test('it shows max ttl picker', async function (assert) {
    assert.expect(4);
    this.ttlKey = 'max_lease_duration';
    await render(hbs`
      <SecretEngine::TtlPickerV2 @model={{this.model}} @ttlKey={{this.ttlKey}} />
    `);
    assert.dom(GENERAL.fieldLabelbyAttr(this.ttlKey)).hasText('Maximum time-to-live (TTL)');
    assert
      .dom(GENERAL.helpTextByAttr(this.ttlKey))
      .hasText('Maximum extension for the secrets life beyond default.');
    await fillIn(GENERAL.inputByAttr(this.ttlKey), 10);
    await fillIn(GENERAL.selectByAttr(this.ttlKey), 'm');
    assert.dom(GENERAL.inputByAttr(this.ttlKey)).hasValue('10');
    assert.dom(GENERAL.selectByAttr(this.ttlKey)).hasValue('m');
  });

  test('it shows an error message if ttl picker time is not a number value', async function (assert) {
    assert.expect(2);
    this.ttlKey = 'max_lease_duration';
    await render(hbs`
      <SecretEngine::TtlPickerV2 @model={{this.model}} @ttlKey={{this.ttlKey}} />
    `);
    await fillIn(GENERAL.inputByAttr(this.ttlKey), 'some text');
    await fillIn(GENERAL.selectByAttr(this.ttlKey), 'm');
    assert.dom(GENERAL.messageError).exists();
    assert.dom(GENERAL.messageError).hasText('Only use numbers for this setting.');
  });
});
