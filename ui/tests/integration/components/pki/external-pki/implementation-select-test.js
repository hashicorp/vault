/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CreationMethod } from 'vault/utils/constants/snippet';
import { SetupSteps } from 'pki/components/external-pki/implementation-select';

module('Integration | Component | pki | external-pki | implementation-select', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.engineId = 'my-pki-external-ca';
    this.title = 'Choose your implementation method';
    this.steps = undefined;
    this.renderComponent = () =>
      render(
        hbs`<ExternalPki::ImplementationSelect @engineId={{this.engineId}} @title={{this.title}} @steps={{this.steps}} />`,
        { owner: this.engine }
      );
  });

  test('it renders with title and method options', async function (assert) {
    await this.renderComponent();
    assert.dom('h1').hasText(this.title);
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.TERRAFORM)).exists().isChecked();
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.APICLI)).exists().isNotChecked();
  });

  test('it does not render title if no title is passed', async function (assert) {
    this.title = undefined;
    await this.renderComponent();
    assert.dom('h1').doesNotExist();
  });

  test('it switches between Terraform and API/CLI methods', async function (assert) {
    await this.renderComponent();

    // Initially Terraform is selected
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.TERRAFORM)).isChecked();
    // Switch to API/CLI
    await click(`input${GENERAL.radioCardByAttr(CreationMethod.APICLI)}`);
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.APICLI)).isChecked('API/CLI is now selected');
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.TERRAFORM)).isNotChecked();
    assert.dom(GENERAL.fieldByAttr('cli')).exists({ count: 2 }, 'cli snippet renders for each config');
    assert.dom(GENERAL.fieldByAttr('api')).exists({ count: 2 }, 'api snippet renders for each config');
    assert.dom(GENERAL.fieldByAttr('terraform')).doesNotExist();
    assert.dom(GENERAL.hdsTab('cli')).exists();
    assert.dom(GENERAL.hdsTab('api')).exists();

    // Switch back to Terraform
    await click(`input${GENERAL.radioCardByAttr(CreationMethod.TERRAFORM)}`);
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.TERRAFORM)).isChecked('Terraform is selected again');
    assert
      .dom(GENERAL.fieldByAttr('terraform'))
      .exists({ count: 2 }, 'Terraform snippet renders for each config');
    assert.dom(GENERAL.fieldByAttr('cli')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('api')).doesNotExist();
  });

  test('it generates config snippets for Terraform', async function (assert) {
    await this.renderComponent();

    const expectedTfvpAcme = `resource "vault_pki_external_ca_secret_backend_acme_account" "<local identifier>" {
   mount = "${this.engineId}"
   name = <name>
   directory_url = <directory_url>
   email_contacts = [<email_contacts>]
  }`;

    const expectedTfvpRole = `resource "vault_pki_external_ca_secret_backend_role" "<local identifier>" {
   mount = "${this.engineId}"
   name = <name>
   allowed_domains = [<allowed_domains>]
   allowed_domain_options = [<allowed_domain_options>]
   allowed_challenge_types = [<allowed_challenge_types>]
  }`;

    const [firstSnippet, secondSnippet] = findAll(GENERAL.fieldByAttr('terraform'));
    assert.dom(firstSnippet).hasText(expectedTfvpAcme, 'first snippet has expected tfvp');
    assert.dom(secondSnippet).hasText(expectedTfvpRole, 'second snippet has expected tfvp');
  });

  test('it generates config snippets for API/CLI', async function (assert) {
    await this.renderComponent();
    await click(`input${GENERAL.radioCardByAttr(CreationMethod.APICLI)}`);
    const expectedCliAcme = `vault write ${this.engineId}/config/acme-account/<name> \\
    directory_url="<directory_url>" \\
    email_contacts="<email_contacts>" \\
  `;
    const expectedCliRole = `vault write ${this.engineId}/role/<name> \\
    allowed_domains="<allowed_domains>" \\
    allowed_domain_options="<allowed_domain_options>" \\
    allowed_challenge_types="<allowed_challenge_types>" \\
  `;

    let [firstSnippet, secondSnippet] = findAll(GENERAL.fieldByAttr('cli'));
    assert.dom(firstSnippet).hasText(expectedCliAcme);
    assert.dom(secondSnippet).hasText(expectedCliRole);
    const expectedApiAcme = `curl \\
    --header "X-Vault-Token: $VAULT_TOKEN" \\
    --request POST \\
    --data '{"directory_url":"<directory_url>","email_contacts":"[<email_contacts>]"}' \\
    $VAULT_ADDR/v1/${this.engineId}/config/acme-account/<name>
  `;
    const expectedApiRole = `curl \\
    --header "X-Vault-Token: $VAULT_TOKEN" \\
    --request POST \\
    --data '{"allowed_domains":"[<allowed_domains>]","allowed_domain_options":"[<allowed_domain_options>]","allowed_challenge_types":"[<allowed_challenge_types>]"}' \\
    $VAULT_ADDR/v1/${this.engineId}/role/<name>
  `;

    [firstSnippet, secondSnippet] = findAll(GENERAL.fieldByAttr('api'));
    assert.dom(firstSnippet).hasText(expectedApiAcme);
    assert.dom(secondSnippet).hasText(expectedApiRole);
  });

  test('it displays both numbered steps by default', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.textDisplay()).exists({ count: 2 });
    assert.dom(GENERAL.textDisplay('0')).exists().hasText('1. Configure an ACME account');
    assert.dom(GENERAL.textDisplay('1')).exists().hasText('2. Create a role');
  });

  test('it displays only ACME config step when specified', async function (assert) {
    this.steps = [SetupSteps.ACME_CONFIG];
    await this.renderComponent();
    assert.dom(GENERAL.textDisplay()).exists({ count: 1 }).hasText('Configure an ACME account');
  });

  test('it displays only role config step when specified', async function (assert) {
    this.steps = [SetupSteps.ROLE_CONFIG];
    await this.renderComponent();
    assert.dom(GENERAL.textDisplay()).exists({ count: 1 }).hasText('Create a role');
  });
});
