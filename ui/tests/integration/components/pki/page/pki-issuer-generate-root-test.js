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
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

/**
 * this test is for the page component only. A separate test is written for the form rendered
 */
module('Integration | Component | page/pki-issuer-generate-root', function (hooks) {
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
    setRunOptions({
      rules: {
        // something strange happening here
        'link-name': { enabled: false },
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it renders correct title before and after submit', async function (assert) {
    assert.expect(3);
    this.server.post(`/pki-component/root/generate/internal`, () => {
      assert.true(true, 'Root endpoint called');
      return {
        request_id: uuidv4(),
        data: {
          certificate: '------BEGIN CERTIFICATE------',
          key_id: 'some-key-id',
        },
      };
    });

    await render(
      hbs`<Page::PkiIssuerGenerateRoot @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(GENERAL.title).hasText('Generate root');
    await fillIn(GENERAL.inputByAttr('type'), 'internal');
    await fillIn(GENERAL.inputByAttr('commonName'), 'foobar');
    await click(GENERAL.saveButton);
    assert.dom(GENERAL.title).hasText('View generated root');
  });

  test('it does not update title if API response is an error', async function (assert) {
    assert.expect(2);
    this.server.post(`/pki-component/root/generate/internal`, () => new Response(404, {}, { errors: [] }));

    await render(
      hbs`<Page::PkiIssuerGenerateRoot @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(GENERAL.title).hasText('Generate root');
    // Fill in
    await fillIn(GENERAL.inputByAttr('type'), 'internal');
    await fillIn(GENERAL.inputByAttr('commonName'), 'foobar');
    await click(GENERAL.saveButton);
    assert.dom(GENERAL.title).hasText('Generate root', 'title does not change if response is unsuccessful');
  });
});
