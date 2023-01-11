import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { click, currentURL, currentRouteName, visit } from '@ember/test-helpers';
import { SELECTORS } from 'vault/tests/helpers/pki/overview';

const tokenWithPolicy = async function (name, policy) {
  await consoleComponent.runCommands([
    `write sys/policies/acl/${name} policy=${btoa(policy)}`,
    `write -field=client_token auth/token/create policies=${name}`,
  ]);
  return consoleComponent.lastLogOutput;
};

const consoleComponent = create(consoleClass);

const runCommands = async function (commands) {
  try {
    await consoleComponent.runCommands(commands);
    const res = consoleComponent.lastLogOutput;
    if (res.includes('Error')) {
      throw new Error(res);
    }
    return res;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error(
      `The following occurred when trying to run the command(s):\n ${commands.join('\n')} \n\n ${
        consoleComponent.lastLogOutput
      }`
    );
    throw error;
  }
};

module('Acceptance | pki overview', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    // Setup PKI engine
    const mountPath = `pki`;
    await enablePage.enable('pki', mountPath);
    this.mountPath = mountPath;
    await runCommands([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
    const pki_admin_policy = `
    path "${this.mountPath}/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
    const pki_issuers_list_policy = `
    path "${this.mountPath}/issuers" {
      capabilities = ["list"]
    },
    `;
    const pki_roles_list_policy = `
    path "${this.mountPath}/roles" {
      capabilities = ["list"]
    },
    `;

    this.pkiRolesList = await tokenWithPolicy('pki-roles-list', pki_roles_list_policy);
    this.pkiIssuersList = await tokenWithPolicy('pki-issuers-list', pki_issuers_list_policy);
    this.pkiAdminToken = await tokenWithPolicy('pki-admin', pki_admin_policy);
    await logout.visit();
  });

  hooks.afterEach(async function () {
    await logout.visit();
    await authPage.login();
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.mountPath}`]);
    await logout.visit();
  });

  test('shows the correct information on issuer card', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert.dom(SELECTORS.issuersCardTitle).hasText('Issuers');
    assert.dom(SELECTORS.issuersCardOverviewNum).hasText('1');
    await click(SELECTORS.issuersCardLink);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers`);
    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
  });

  test('shows the correct information on roles card', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert.dom(SELECTORS.rolesCardTitle).hasText('Roles');
    assert.dom(SELECTORS.rolesCardOverviewNum).hasText('0');
    await click(SELECTORS.rolesCardLink);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
    await runCommands([
      `write ${this.mountPath}/roles/some-role \
    issuer_ref="default" \
    allowed_domains="example.com" \
    allow_subdomains=true \
    max_ttl="720h"`,
    ]);
    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert.dom(SELECTORS.rolesCardOverviewNum).hasText('1');
  });

  test('hides roles card if user does not have permissions', async function (assert) {
    await authPage.login(this.pkiIssuersList);
    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert.dom(SELECTORS.rolesCardTitle).doesNotExist('Roles card does not exist');
    assert.dom(SELECTORS.issuersCardTitle).exists('Issuers card exists');
  });

  test('hides issuers card if user does not have permissions', async function (assert) {
    await authPage.login(this.pkiRolesList);
    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert.dom(SELECTORS.issuersCardTitle).doesNotExist('Issuers card does not exist');
    assert.dom(SELECTORS.rolesCardTitle).exists('Roles card does exist');
  });

  test('navigates to generate credentials page for Issue Certificates card', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await runCommands([
      `write ${this.mountPath}/roles/some-role \
    issuer_ref="default" \
    allowed_domains="example.com" \
    allow_subdomains=true \
    max_ttl="720h"`,
    ]);
    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert.dom(SELECTORS.issueCertificate).hasText('Issue certificate');
    await click(SELECTORS.issueCertificateInput);
    await click(SELECTORS.firstPowerSelectOption);
    await click(SELECTORS.issueCertificateButton);
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.roles.role.generate');
  });

  test('navigates to certificate details page for Issue Certificates card', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert.dom(SELECTORS.viewCertificate).hasText('View certificate');
    await click(SELECTORS.viewCertificateInput);
    await click(SELECTORS.firstPowerSelectOption);
    await click(SELECTORS.viewCertificateButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.certificates.certificate.details'
    );
  });
});
