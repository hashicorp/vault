/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { createConfig } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';

module('Integration | Component | SecretEngine/configure-ssh', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    const router = this.owner.lookup('service:router');
    this.id = 'ssh-test';
    this.model = this.store.createRecord('ssh/ca-config', { backend: this.id });
    this.transitionStub = sinon.stub(router, 'transitionTo');
    this.refreshStub = sinon.stub(router, 'refresh');
  });

  test('it shows create fields if not configured', async function (assert) {
    await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.model}}
    @id={{this.id}}
  />
    `);
    assert.dom(GENERAL.inputByAttr('privateKey')).hasNoText('Private key is empty and reset');
    assert.dom(GENERAL.inputByAttr('publicKey')).hasNoText('Public key is empty and reset');
    assert
      .dom(GENERAL.inputByAttr('generateSigningKey'))
      .isChecked('Generate signing key is checked by default');
  });

  test('it should go back to parent route on cancel', async function (assert) {
    await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.model}}
    @id={{this.id}}
  />
    `);

    await click(SES.ssh.cancel);

    assert.true(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', 'ssh-test'),
      'On cancel the router transitions to the parent configuration index route.'
    );
  });

  test('it should validate form fields', async function (assert) {
    await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.model}}
    @id={{this.id}}
  />
    `);
    await fillIn(GENERAL.inputByAttr('publicKey'), 'hello');
    await click(SES.ssh.save);
    assert
      .dom(GENERAL.inlineError)
      .hasText(
        'You must provide a Public and Private keys or leave both unset.',
        'Public key validation error renders.'
      );

    await click(GENERAL.inputByAttr('generateSigningKey'));
    await click(SES.ssh.save);
    assert
      .dom(GENERAL.inlineError)
      .hasText(
        'You must provide a Public and Private keys or leave both unset.',
        'Generate signing key validation message shows.'
      );
  });

  test('it should generate signing key', async function (assert) {
    assert.expect(2);
    this.server.post('/ssh-test/config/ca', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = {
        backend: this.id,
        generate_signing_key: true,
      };
      assert.deepEqual(expected, data, 'POST request made to save ca-config with correct properties');
    });
    await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.model}}
    @id={{this.id}}
  />
    `);

    await click(SES.ssh.save);
    assert.dom(SES.ssh.editConfigSection).exists('renders the edit configuration section of the form');
  });

  module('editing', function (hooks) {
    hooks.beforeEach(function () {
      this.editId = 'ssh-edit-me';
      this.editModel = createConfig(this.store, 'ssh-edit-me', 'ssh');
    });
    test('it populates fields when editing', async function (assert) {
      await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.editModel}}
    @id={{this.editId}}
  />
    `);
      assert
        .dom(SES.ssh.editConfigSection)
        .exists('renders the edit configuration section of the form and not the create part');
      assert.dom(GENERAL.inputByAttr('public-key')).hasText('***********', 'public key is masked');
      await click('[data-test-button="toggle-masked"]');
      assert
        .dom(GENERAL.inputByAttr('public-key'))
        .hasText(this.editModel.publicKey, 'public key is unmasked and shows the actual value');
    });

    test('it allows you to delete a public key', async function (assert) {
      assert.expect(3);
      this.server.delete('/ssh-edit-me/config/ca', () => {
        assert.true(true, 'DELETE request made to ca-config with correct properties');
      });
      await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.editModel}}
    @id={{this.editId}}
  />
    `);
      // delete Public key
      await click(SES.ssh.delete);
      assert.dom(GENERAL.confirmMessage).hasText('Confirming will remove the CA certificate information.');
      await click(GENERAL.confirmButton);
      assert.true(
        this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration.edit', 'ssh-edit-me'),
        'On delete the router transitions to the current route.'
      );
    });
  });
});
