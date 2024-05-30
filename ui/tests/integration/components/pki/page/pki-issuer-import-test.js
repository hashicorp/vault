/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { click, fillIn, render } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { Response } from 'miragejs';
import { v4 as uuidv4 } from 'uuid';
import { setupRenderingTest } from 'vault/tests/helpers';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_CONFIGURE_CREATE } from 'vault/tests/helpers/pki/pki-selectors';

/**
 * this test is for the page component only. A separate test is written for the form rendered
 */
module('Integration | Component | page/pki-issuer-import', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.breadcrumbs = [{ label: 'something' }];
    this.model = this.store.createRecord('pki/action', {
      actionType: 'generate-csr',
    });
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-component';
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
  });

  test('it renders correct title before and after submit', async function (assert) {
    assert.expect(3);
    this.server.post(`/pki-component/issuers/import/bundle`, () => {
      assert.true(true, 'Import endpoint called');
      return {
        request_id: uuidv4(),
        data: {},
      };
    });

    await render(hbs`<Page::PkiIssuerImport @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert.dom(GENERAL.title).hasText('Import a CA');
    await click('[data-test-text-toggle]');
    await fillIn('[data-test-text-file-textarea]', 'dummy-pem-bundle');
    await click(PKI_CONFIGURE_CREATE.importSubmit);
    assert.dom(GENERAL.title).hasText('View imported items');
  });

  test('it does not update title if API response is an error', async function (assert) {
    assert.expect(2);
    this.server.post(`/pki-component/issuers/import/bundle`, () => new Response(404, {}, { errors: [] }));

    await render(hbs`<Page::PkiIssuerImport @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert.dom(GENERAL.title).hasText('Import a CA');
    // Fill in
    await click('[data-test-text-toggle]');
    await fillIn('[data-test-text-file-textarea]', 'dummy-pem-bundle');
    await click(PKI_CONFIGURE_CREATE.importSubmit);
    assert.dom(GENERAL.title).hasText('Import a CA', 'title does not change if response is unsuccessful');
  });
});
