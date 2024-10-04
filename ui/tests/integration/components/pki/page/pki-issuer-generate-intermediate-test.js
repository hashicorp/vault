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

/**
 * this test is for the page component only. A separate test is written for the form rendered
 */
module('Integration | Component | page/pki-issuer-generate-intermediate', function (hooks) {
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
    this.server.post(`/pki-component/issuers/generate/intermediate/internal`, () => {
      assert.true(true, 'Issuers endpoint called');
      return {
        request_id: uuidv4(),
        data: {
          csr: '------BEGIN CERTIFICATE------',
          key_id: 'some-key-id',
        },
      };
    });

    await render(
      hbs`<Page::PkiIssuerGenerateIntermediate @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(GENERAL.title).hasText('Generate intermediate CSR');
    await fillIn(GENERAL.inputByAttr('type'), 'internal');
    await fillIn(GENERAL.inputByAttr('commonName'), 'foobar');
    await click('[data-test-save]');
    assert.dom(GENERAL.title).hasText('View Generated CSR');
  });

  test('it does not update title if API response is an error', async function (assert) {
    assert.expect(2);
    this.server.post(
      '/pki-component/issuers/generate/intermediate/internal',
      () => new Response(403, {}, { errors: ['API returns this error'] })
    );

    await render(
      hbs`<Page::PkiIssuerGenerateIntermediate @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(GENERAL.title).hasText('Generate intermediate CSR');
    // Fill in
    await fillIn(GENERAL.inputByAttr('type'), 'internal');
    await fillIn(GENERAL.inputByAttr('commonName'), 'foobar');
    await click('[data-test-save]');
    assert
      .dom(GENERAL.title)
      .hasText('Generate intermediate CSR', 'title does not change if response is unsuccessful');
  });
});
