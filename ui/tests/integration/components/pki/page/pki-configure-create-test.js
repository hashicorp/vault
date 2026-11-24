/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_CONFIGURE_CREATE } from 'vault/tests/helpers/pki/pki-selectors';
import { configCapabilities } from 'vault/tests/helpers/pki/pki-helpers';

module('Integration | Component | page/pki-configure-create', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.context = { owner: this.engine }; // this.engine set by setupEngine
    this.cancelSpy = sinon.spy();
    this.capabilities = configCapabilities;
  });

  test('it renders', async function (assert) {
    await render(
      hbs`
        <Page::PkiConfigureCreate
          @capabilities={{this.capabilities}}
          @onCancel={{this.cancelSpy}}
        />
    `,
      this.context
    );
    assert.dom(PKI_CONFIGURE_CREATE.option).exists({ count: 3 });
    assert.dom(GENERAL.cancelButton).exists('Cancel link is shown');
    assert.dom(GENERAL.submitButton).isDisabled('Done button is disabled');

    await click(PKI_CONFIGURE_CREATE.optionByKey('import'));
    assert.dom(PKI_CONFIGURE_CREATE.optionByKey('import')).isChecked();

    await click(PKI_CONFIGURE_CREATE.optionByKey('generate-csr'));
    assert.dom(PKI_CONFIGURE_CREATE.optionByKey('generate-csr')).isChecked();

    await click(PKI_CONFIGURE_CREATE.optionByKey('generate-root'));
    assert.dom(PKI_CONFIGURE_CREATE.optionByKey('generate-root')).isChecked();

    await click(GENERAL.cancelButton);
    assert.true(this.cancelSpy.calledOnce, 'cancel action is called');
  });
});
