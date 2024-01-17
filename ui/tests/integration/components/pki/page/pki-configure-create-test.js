/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-configure-create';
import sinon from 'sinon';

module('Integration | Component | page/pki-configure-create', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.context = { owner: this.engine }; // this.engine set by setupEngine
    this.store = this.owner.lookup('service:store');
    this.cancelSpy = sinon.spy();
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: 'pki', route: 'overview' },
      { label: 'configure' },
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
    assert.dom(SELECTORS.breadcrumbContainer).exists('breadcrumbs exist');
    assert.dom(SELECTORS.title).hasText('Configure PKI');
    assert.dom(SELECTORS.option).exists({ count: 3 }, 'Three configuration options are shown');
    assert.dom(SELECTORS.cancelButton).exists('Cancel link is shown');
    assert.dom(SELECTORS.saveButton).isDisabled('Done button is disabled');

    await click(SELECTORS.optionByKey('import'));
    assert.dom(SELECTORS.optionByKey('import')).isChecked('Selected item is checked');

    await click(SELECTORS.optionByKey('generate-csr'));
    assert.dom(SELECTORS.optionByKey('generate-csr')).isChecked('Selected item is checked');

    await click(SELECTORS.optionByKey('generate-root'));
    assert.dom(SELECTORS.optionByKey('generate-root')).isChecked('Selected item is checked');

    await click(SELECTORS.generateRootCancel);
    assert.ok(this.cancelSpy.calledOnce);
  });
});
