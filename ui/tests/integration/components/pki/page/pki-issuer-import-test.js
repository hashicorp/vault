/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_CONFIGURE_CREATE } from 'vault/tests/helpers/pki/pki-selectors';
import sinon from 'sinon';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

/**
 * this test is for the page component only. A separate test is written for the form rendered
 */
module('Integration | Component | page/pki-issuer-import', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.breadcrumbs = [{ label: 'something' }];
    this.backend = 'pki-component';
    this.secretMountPath = this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.importStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'pkiIssuersImportBundle')
      .resolves({});

    this.renderComponent = () =>
      render(hbs`<Page::PkiIssuerImport @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
  });

  test('it renders correct title before and after submit', async function (assert) {
    assert.expect(3);

    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Import a CA');

    const pem_bundle = 'dummy-pem-bundle';
    await click(GENERAL.textToggle);
    await fillIn(GENERAL.maskedInput, pem_bundle);
    await click(PKI_CONFIGURE_CREATE.importSubmit);

    assert.true(this.importStub.calledWith(this.backend, { pem_bundle }), 'API called with correct params');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('View imported items');
  });

  test('it does not update title if API response is an error', async function (assert) {
    assert.expect(2);

    this.importStub.rejects(getErrorResponse());

    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Import a CA');
    // Fill in
    await click(GENERAL.textToggle);
    await fillIn(GENERAL.maskedInput, 'dummy-pem-bundle');
    await click(PKI_CONFIGURE_CREATE.importSubmit);
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText('Import a CA', 'title does not change if response is unsuccessful');
  });
});
