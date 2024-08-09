/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { createConfig } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import sinon from 'sinon';

module('Integration | Component | SecretEngine/configure-ssh', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = createConfig(this.store, 'ssh-test', 'ssh');
    this.saveConfig = sinon.stub();
  });

  test('it shows create fields if not configured', async function (assert) {
    await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.model}}
    @configured={{false}}
    @saveConfig={{this.saveConfig}}
    @loading={{false}}
  />
    `);
    assert.dom(GENERAL.maskedInput('privateKey')).hasNoText('Private key is empty and reset');
    assert.dom(GENERAL.inputByAttr('publicKey')).hasNoText('Public key is empty and reset');
    assert
      .dom(GENERAL.inputByAttr('generate-signing-key-checkbox'))
      .isChecked('Generate signing key is checked by default');
  });

  test('it calls save with correct arg', async function (assert) {
    await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.model}}
    @configured={{false}}
    @saveConfig={{this.saveConfig}}
    @loading={{false}}
  />
    `);
    await click(SES.ssh.save);
    assert.ok(
      this.saveConfig.withArgs({ delete: false }).calledOnce,
      'calls the saveConfig action with args delete:false'
    );
  });

  test('it shows masked key if model is not new', async function (assert) {
    // replace model with model that has public_key
    this.model = {
      publicKey:
        'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC3lCZ7W2eJZ9W9qzv7K9GJ5qJYQ2cY6C+5Kv8Jtjz8h6wqZJ9U9K1lJ9Z6zq4sX0f7Q5X2l8L4gTt2+2ZKpVv6g1KQ6JG5H4QbVrQq2r4FzZQ2B0Y8q5c7q3Y5X6q4Q6',
    };
    await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.model}}
    @configured={{true}}
    @saveConfig={{this.saveConfig}}
    @loading={{false}}
  />
    `);
    assert
      .dom(SES.ssh.editConfigSection)
      .exists('renders the edit configuration section of the form and not the create part');
    assert.dom(GENERAL.inputByAttr('public-key')).hasText('***********', 'public key is masked');
    await click('[data-test-button="toggle-masked"]');
    assert
      .dom(GENERAL.inputByAttr('public-key'))
      .hasText(this.model.publicKey, 'public key is unmasked and shows the actual value');
  });

  test('it calls delete correctly', async function (assert) {
    await render(hbs`
      <SecretEngine::ConfigureSsh
    @model={{this.model}}
    @configured={{true}}
    @saveConfig={{this.saveConfig}}
    @loading={{false}}
  />
    `);
    // delete Public key
    await click(SES.ssh.deletePublicKey);
    assert.dom(GENERAL.confirmMessage).hasText('This will remove the CA certificate information.');
    await click(GENERAL.confirmButton);
    assert.ok(
      this.saveConfig.withArgs({ delete: true }).calledOnce,
      'calls the saveConfig action with args delete:true'
    );
  });
});
