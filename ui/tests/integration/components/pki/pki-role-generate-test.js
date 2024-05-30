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
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_ROLE_GENERATE } from 'vault/tests/helpers/pki/pki-selectors';
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
    assert.dom(PKI_ROLE_GENERATE.form).exists('shows the cert generate form');
    assert.dom(GENERAL.inputByAttr('commonName')).exists('shows the common name field');
    assert.dom(PKI_ROLE_GENERATE.optionsToggle).exists('toggle exists');
    await fillIn(GENERAL.inputByAttr('commonName'), 'example.com');
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
    assert.dom(PKI_ROLE_GENERATE.form).doesNotExist('Does not show the form');
    assert.dom(PKI_ROLE_GENERATE.downloadButton).exists('shows the download button');
    assert.dom(PKI_ROLE_GENERATE.revokeButton).exists('shows the revoke button');
    assert.dom(GENERAL.infoRowValue('Certificate')).exists({ count: 1 }, 'shows certificate info row');
    assert
      .dom(GENERAL.infoRowValue('Serial number'))
      .hasText('abcd-efgh-ijkl', 'shows serial number info row');
  });
});
