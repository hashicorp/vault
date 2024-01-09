/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupEngine } from 'ember-engines/test-support';
import Sinon from 'sinon';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-role-generate';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module('Integration | Component | pki-role-generate', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.model = this.store.createRecord('pki/certificate/generate', {
      role: 'my-role',
    });
    this.onSuccess = Sinon.spy();
  });

  test('it should render the component with the form by default', async function (assert) {
    assert.expect(4);
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiRoleGenerate
          @model={{this.model}}
          @onSuccess={{this.onSuccess}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.form).exists('shows the cert generate form');
    assert.dom(SELECTORS.commonNameField).exists('shows the common name field');
    assert.dom(SELECTORS.optionsToggle).exists('toggle exists');
    await fillIn(SELECTORS.commonNameField, 'example.com');
    assert.strictEqual(this.model.commonName, 'example.com', 'Filling in the form updates the model');
  });

  test('it should render the component displaying the cert', async function (assert) {
    assert.expect(5);
    const record = this.store.createRecord('pki/certificate/generate', {
      role: 'my-role',
      serialNumber: 'abcd-efgh-ijkl',
      certificate: 'my-very-cool-certificate',
    });
    this.set('model', record);
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiRoleGenerate
          @model={{this.model}}
          @onSuccess={{this.onSuccess}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.form).doesNotExist('Does not show the form');
    assert.dom(SELECTORS.downloadButton).exists('shows the download button');
    assert.dom(SELECTORS.revokeButton).exists('shows the revoke button');
    assert.dom(SELECTORS.certificate).exists({ count: 1 }, 'shows certificate info row');
    assert.dom(SELECTORS.serialNumber).hasText('abcd-efgh-ijkl', 'shows serial number info row');
  });
});
