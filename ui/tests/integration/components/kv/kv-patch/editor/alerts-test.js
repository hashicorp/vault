/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { FORM } from 'vault/tests/helpers/kv/kv-selectors';

module('Integration | Component | kv | kv-patch/editor/alerts', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(function () {
    this.keyError = '';
    this.keyWarning = '';
    this.valueWarning = '';
    this.renderComponent = async () => {
      return render(
        hbs`
      <KvPatch::Editor::Alerts
        @idx={{1}}
        @keyError={{this.keyError}}
        @keyWarning={{this.keyWarning}}
        @valueWarning={{this.valueWarning}}
      />`,
        { owner: this.engine }
      );
    };
  });

  test('it renders', async function (assert) {
    await this.renderComponent();
    assert.dom(FORM.patchAlert('validation', 1)).doesNotExist();
    assert.dom(FORM.patchAlert('value-warning', 1)).doesNotExist();
    assert.dom(FORM.patchAlert('key-warning', 1)).doesNotExist();
  });

  test('it renders key error', async function (assert) {
    this.keyError = "There's a problem with your key";
    await this.renderComponent();
    assert.dom(FORM.patchAlert('validation', 1)).hasClass('hds-alert--color-critical');
    assert.dom(`${FORM.patchAlert('validation', 1)} ${GENERAL.icon('alert-diamond-fill')}`).exists();

    assert.dom(FORM.patchAlert('key-warning', 1)).doesNotExist();
    assert.dom(FORM.patchAlert('value-warning', 1)).doesNotExist();
  });

  test('it renders key warning', async function (assert) {
    this.keyWarning = 'Key warning';
    await this.renderComponent();
    assert.dom(FORM.patchAlert('key-warning', 1)).hasClass('hds-alert--color-warning');
    assert.dom(`${FORM.patchAlert('key-warning', 1)} ${GENERAL.icon('alert-triangle-fill')}`).exists();
    assert.dom(FORM.patchAlert('validation', 1)).doesNotExist();
    assert.dom(FORM.patchAlert('value-warning', 1)).doesNotExist();
  });

  test('it renders value warning', async function (assert) {
    this.valueWarning = 'Value warning';
    await this.renderComponent();
    assert.dom(FORM.patchAlert('value-warning', 1)).hasClass('hds-alert--color-warning');
    assert.dom(`${FORM.patchAlert('value-warning', 1)} ${GENERAL.icon('alert-triangle-fill')}`).exists();
    assert.dom(FORM.patchAlert('validation', 1)).doesNotExist();
    assert.dom(FORM.patchAlert('key-warning', 1)).doesNotExist();
  });

  test('it renders all three alerts', async function (assert) {
    this.keyError = "There's a problem with your key";
    this.keyWarning = 'Key warning';
    this.valueWarning = 'Value warning';
    await this.renderComponent();
    assert.dom(FORM.patchAlert('validation', 1)).hasText(this.keyError);
    assert.dom(FORM.patchAlert('key-warning', 1)).hasText(this.keyWarning);
    assert.dom(FORM.patchAlert('value-warning', 1)).hasText(this.valueWarning);
  });
});
