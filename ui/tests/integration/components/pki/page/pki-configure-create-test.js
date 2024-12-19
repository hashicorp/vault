/**
 * Copyright (c) HashiCorp, Inc.
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

module('Integration | Component | page/pki-configure-create', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.context = { owner: this.engine }; // this.engine set by setupEngine
    this.store = this.owner.lookup('service:store');
    this.cancelSpy = sinon.spy();
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: 'pki', route: 'overview', model: 'pki' },
      { label: 'Configure' },
    ];
    this.config = this.store.createRecord('pki/action');
    this.urls = this.store.createRecord('pki/config/urls');
  });

  test('it renders', async function (assert) {
    await render(
      hbs`
      <Page::PkiConfigureCreate
        @breadcrumbs={{this.breadcrumbs}}
        @config={{this.config}}
        @urls={{this.urls}}
        @onCancel={{this.cancelSpy}}
      />
    `,
      this.context
    );
    assert.dom(GENERAL.breadcrumbs).exists();
    assert.dom(GENERAL.title).hasText('Configure PKI');
    assert.dom(PKI_CONFIGURE_CREATE.option).exists({ count: 3 });
    assert.dom(GENERAL.cancelButton).exists('Cancel link is shown');
    assert.dom(GENERAL.saveButton).isDisabled('Done button is disabled');

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
