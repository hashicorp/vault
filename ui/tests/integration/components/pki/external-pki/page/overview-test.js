/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { CreationMethod } from 'vault/utils/constants/snippet';
import sinon from 'sinon';

module('Integration | Component | pki | external-pki | ExternalPki::Page::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    const perms = {
      pkiExternalConfigAcmeAccount: {
        canCreate: true,
        canDelete: true,
        canList: true,
        canPatch: true,
        canRead: true,
        canSudo: true,
        canUpdate: true,
      },
      pkiExternalConfigDns: {
        canCreate: true,
        canDelete: true,
        canList: true,
        canPatch: true,
        canRead: true,
        canSudo: true,
        canUpdate: true,
      },
      pkiExternalRole: {
        canCreate: true,
        canDelete: true,
        canList: true,
        canPatch: true,
        canRead: true,
        canSudo: true,
        canUpdate: true,
      },
      pkiExternalLookupOrders: {
        canCreate: true,
        canDelete: true,
        canList: true,
        canPatch: true,
        canRead: true,
        canSudo: true,
        canUpdate: true,
      },
    };

    this.setupModel = (permOverrides) => {
      return {
        engine: new SecretsEngineResource({
          accessor: 'pki-external-ca_e158c567',
          type: 'pki-external-ca',
          path: 'my-pki-external-ca/',
        }),
        permissions: { ...perms, ...permOverrides },
        acmeAccounts: { keys: [], errorMsg: '' },
        dnsProviders: { keys: [], errorMsg: '' },
        roles: { keys: [], errorMsg: '' },
      };
    };
    this.isNotConfigured = undefined;
    this.renderComponent = () =>
      render(
        hbs`<ExternalPki::Page::Overview @model={{this.model}} @isNotConfigured={{this.isNotConfigured}} />
        `,
        { owner: this.engine }
      );
  });

  test('it renders overview cards by default (when isNotConfigured=undefined)', async function (assert) {
    this.model = this.setupModel();
    await this.renderComponent();

    assert.dom('h1').doesNotExist();
    assert.dom(GENERAL.radioCardByAttr()).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('terraform')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('cli')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('api')).doesNotExist();
  });

  module('is not configured', function (hooks) {
    hooks.beforeEach(function () {
      this.isNotConfigured = true;
      this.model = this.setupModel();
    });

    test('it renders implementation select when isNotConfigured is true', async function (assert) {
      await this.renderComponent();

      assert.dom('h1').hasText('Choose your implementation method');
      assert
        .dom(GENERAL.radioCardByAttr(CreationMethod.TERRAFORM))
        .exists()
        .isChecked('Terraform is initially selected');
      assert.dom(GENERAL.radioCardByAttr(CreationMethod.APICLI)).exists().isNotChecked();

      assert
        .dom(GENERAL.fieldByAttr('terraform'))
        .exists({ count: 2 }, 'Terraform snippet renders for each config');
      assert.dom(GENERAL.fieldByAttr('cli')).doesNotExist();
      assert.dom(GENERAL.fieldByAttr('api')).doesNotExist();
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

    test('it generates correct config snippets for Terraform', async function (assert) {
      await this.renderComponent();

      const expectedTfvpAcme = `resource "vault_pki_external_ca_secret_backend_acme_account" "<local identifier>" {
 mount = "${this.model.engine.id}"
 name = <name>
 directory_url = <directory_url>
 email_contacts = [<email_contacts>]
}`;

      const expectedTfvpRole = `resource "vault_pki_external_ca_secret_backend_role" "<local identifier>" {
 mount = "${this.model.engine.id}"
 name = <name>
 allowed_domains = [<allowed_domains>]
 allowed_domain_options = [<allowed_domain_options>]
 allowed_challenge_types = [<allowed_challenge_types>]
}`;

      const [firstSnippet, secondSnippet] = findAll(GENERAL.fieldByAttr('terraform'));
      assert.dom(firstSnippet).hasText(expectedTfvpAcme, 'first snippet has expected tfvp');
      assert.dom(secondSnippet).hasText(expectedTfvpRole, 'second snippet has expected tfvp');
    });

    test('it generates correct config snippets for API/CLI', async function (assert) {
      await this.renderComponent();
      await click(`input${GENERAL.radioCardByAttr(CreationMethod.APICLI)}`);
      const expectedCliAcme = `vault write ${this.model.engine.id}/config/acme-account/<name> \\
  directory_url="<directory_url>" \\
  email_contacts="<email_contacts>" \\
`;
      const expectedCliRole = `vault write ${this.model.engine.id}/role/<name> \\
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
  $VAULT_ADDR/v1/${this.model.engine.id}/config/acme-account/<name>
`;
      const expectedApiRole = `curl \\
  --header "X-Vault-Token: $VAULT_TOKEN" \\
  --request POST \\
  --data '{"allowed_domains":"[<allowed_domains>]","allowed_domain_options":"[<allowed_domain_options>]","allowed_challenge_types":"[<allowed_challenge_types>]"}' \\
  $VAULT_ADDR/v1/${this.model.engine.id}/role/<name>
`;

      [firstSnippet, secondSnippet] = findAll(GENERAL.fieldByAttr('api'));
      assert.dom(firstSnippet).hasText(expectedApiAcme);
      assert.dom(secondSnippet).hasText(expectedApiRole);
    });
  });

  module('is configured', function (hooks) {
    hooks.beforeEach(function () {
      this.isNotConfigured = false;
      const router = this.owner.lookup('service:router');
      this.transitionToStub = sinon.stub(router, 'transitionTo');
    });

    test('it renders overview cards with zero counts', async function (assert) {
      this.model = this.setupModel();
      await this.renderComponent();

      assert.dom('h1').doesNotExist('implementation select header does not render');
      assert.dom(GENERAL.radioCardByAttr()).doesNotExist('radio cards do not render');

      // Check lookup cards render
      assert
        .dom(GENERAL.overviewCard.container('View certificate'))
        .exists()
        .hasTextContaining('View certificate');
      assert.dom(GENERAL.overviewCard.container('Orders')).exists().hasTextContaining('Orders');

      // Check count cards
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).exists();
      assert.dom(GENERAL.overviewCard.content('ACME accounts')).hasText('0');
      assert.dom(GENERAL.overviewCard.container('DNS providers')).exists();
      assert.dom(GENERAL.overviewCard.content('DNS providers')).hasText('0');
      assert.dom(GENERAL.overviewCard.container('Roles')).exists();
      assert.dom(GENERAL.overviewCard.content('Roles')).hasText('0');
    });

    test('it renders overview cards with non-zero counts', async function (assert) {
      this.model = this.setupModel();
      this.model.acmeAccounts = { keys: ['account-1', 'account-2'], errorMsg: null };
      this.model.dnsProviders = { keys: ['provider-1'], errorMsg: null };
      this.model.roles = { keys: ['role-1', 'role-2', 'role-3'], errorMsg: null };
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.content('ACME accounts')).hasText('2');
      assert.dom(GENERAL.overviewCard.content('DNS providers')).hasText('1');
      assert.dom(GENERAL.overviewCard.content('Roles')).hasText('3');
    });

    test('it hides count cards when no read or list permissions', async function (assert) {
      this.model = this.setupModel({
        pkiExternalConfigAcmeAccount: { canList: false, canRead: false },
        pkiExternalConfigDns: { canList: false, canRead: false },
        pkiExternalRole: { canList: false, canRead: false },
      });
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).doesNotExist();
      assert.dom(GENERAL.overviewCard.container('DNS providers')).doesNotExist();
      assert.dom(GENERAL.overviewCard.container('Roles')).doesNotExist();
    });

    test('it renders count cards with read only permissions', async function (assert) {
      this.model = this.setupModel({
        pkiExternalConfigAcmeAccount: { canRead: true },
        pkiExternalConfigDns: { canRead: true },
        pkiExternalRole: { canRead: true },
      });
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).exists();
      assert.dom(GENERAL.overviewCard.container('DNS providers')).exists();
      assert.dom(GENERAL.overviewCard.container('Roles')).exists();
    });

    test('it renders count cards with list only permissions', async function (assert) {
      this.model = this.setupModel({
        pkiExternalConfigAcmeAccount: { canList: true },
        pkiExternalConfigDns: { canList: true },
        pkiExternalRole: { canList: true },
      });
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).exists();
      assert.dom(GENERAL.overviewCard.container('DNS providers')).exists();
      assert.dom(GENERAL.overviewCard.container('Roles')).exists();
    });

    test('it hides only ACME accounts card', async function (assert) {
      this.model = this.setupModel({ pkiExternalConfigAcmeAccount: { canList: false } });
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).doesNotExist();
      assert.dom(GENERAL.overviewCard.container('DNS providers')).exists();
      assert.dom(GENERAL.overviewCard.container('Roles')).exists();
    });

    test('it hides only DNS providers card', async function (assert) {
      this.model = this.setupModel({ pkiExternalConfigDns: { canList: false } });
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).exists();
      assert.dom(GENERAL.overviewCard.container('DNS providers')).doesNotExist();
      assert.dom(GENERAL.overviewCard.container('Roles')).exists();
    });

    test('it hides only Roles card', async function (assert) {
      this.model = this.setupModel({ pkiExternalRole: { canList: false } });
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).exists();
      assert.dom(GENERAL.overviewCard.container('DNS providers')).exists();
      assert.dom(GENERAL.overviewCard.container('Roles')).doesNotExist();
    });

    test('it displays error message when count card has error', async function (assert) {
      this.model = this.setupModel();
      this.model.acmeAccounts.errorMsg = 'ACME accounts server error';
      this.model.dnsProviders.errorMsg = 'DNS provider server error';

      await this.renderComponent();
      assert
        .dom(`${GENERAL.overviewCard.container('ACME accounts')} ${GENERAL.messageError}`)
        .exists()
        .hasText('Error ACME accounts server error');
      assert
        .dom(`${GENERAL.overviewCard.container('DNS providers')} ${GENERAL.messageError}`)
        .exists()
        .hasText('Error DNS provider server error');
      assert.dom(GENERAL.overviewCard.container('Roles')).hasTextContaining('0');
    });

    test('it transitions to look up certificate', async function (assert) {
      this.model = this.setupModel();
      await this.renderComponent();
      await fillIn(`${GENERAL.overviewCard.container('View certificate')} input`, '03:e7:1f:');
      await click(GENERAL.button('Lookup certificate'));
      const [route, engineId, serialNumber] = this.transitionToStub.lastCall.args;
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.pki.external.certificates.certificate',
        'it calls transition with expected route'
      );
      assert.strictEqual(engineId, 'my-pki-external-ca', 'it calls transition with expected engine ID');
      assert.strictEqual(serialNumber, '03:e7:1f:', 'it calls transition with inputted serial number');
    });

    test('it transitions to look up order', async function (assert) {
      this.model = this.setupModel();
      await this.renderComponent();
      await fillIn(`${GENERAL.overviewCard.container('Orders')} input`, '01936d8e-7c3');
      await click(GENERAL.button('Lookup order'));
      const [route, engineId, orderId] = this.transitionToStub.lastCall.args;
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.pki.external.orders.order',
        'it calls transition with expected route'
      );
      assert.strictEqual(engineId, 'my-pki-external-ca', 'it calls transition with expected engine ID');
      assert.strictEqual(orderId, '01936d8e-7c3', 'it calls transition with inputted order ID');
    });
  });
});
