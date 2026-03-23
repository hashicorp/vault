/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { click, currentURL, currentRouteName, visit } from '@ember/test-helpers';
import { runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { PKI_OVERVIEW } from 'vault/tests/helpers/pki/pki-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
const { overviewCard } = GENERAL;

module('Acceptance | pki overview', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
    // Setup PKI engine
    const mountPath = `pki-${uuidv4()}`;
    await enablePage.enable('pki', mountPath);
    this.mountPath = mountPath;
    await runCmd([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
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

    this.pkiRolesList = await runCmd(tokenWithPolicyCmd('pki-roles-list', pki_roles_list_policy));
    this.pkiIssuersList = await runCmd(tokenWithPolicyCmd('pki-issuers-list', pki_issuers_list_policy));
    this.pkiAdminToken = await runCmd(tokenWithPolicyCmd('pki-admin', pki_admin_policy));
  });

  hooks.afterEach(async function () {
    await login();
    // Cleanup engine
    await runCmd([`delete sys/mounts/${this.mountPath}`]);
  });

  test('navigates to view issuers when link is clicked on issuer card', async function (assert) {
    await login(this.pkiAdminToken);
    await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
    assert.dom(overviewCard.title('Issuers')).hasText('Issuers');
    assert.dom(`${overviewCard.container('Issuers')} p`).hasText('1');
    await click(overviewCard.actionLink('Issuers'));
    assert.strictEqual(currentURL(), `/vault/secrets-engines/${this.mountPath}/pki/issuers`);
    await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
  });

  test('navigates to view roles when link is clicked on roles card', async function (assert) {
    await login(this.pkiAdminToken);
    await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
    assert.dom(overviewCard.title('Roles')).hasText('Roles');
    assert.dom(`${overviewCard.container('Roles')} p`).hasText('0');
    await click(overviewCard.actionLink('Roles'));
    assert.strictEqual(currentURL(), `/vault/secrets-engines/${this.mountPath}/pki/roles`);
    await runCmd([
      `write ${this.mountPath}/roles/some-role \
    issuer_ref="default" \
    allowed_domains="example.com" \
    allow_subdomains=true \
    max_ttl="720h"`,
    ]);
    await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
    assert.dom(`${overviewCard.container('Roles')} p`).hasText('1');
  });

  test('hides roles and certificates card if user does not have permissions', async function (assert) {
    await login(this.pkiIssuersList);
    await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
    assert.dom(overviewCard.title('Roles')).doesNotExist('Roles card does not exist');
    assert.dom(overviewCard.title('Certificates')).doesNotExist('Certificates card does not exist');
    assert.dom(overviewCard.title('Issuers')).hasText('Issuers');
  });

  test('navigates to generate certificate page for Issue Certificates card', async function (assert) {
    await login(this.pkiAdminToken);
    await runCmd([
      `write ${this.mountPath}/roles/some-role \
    issuer_ref="default" \
    allowed_domains="example.com" \
    allow_subdomains=true \
    max_ttl="720h"`,
    ]);
    await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
    await click(PKI_OVERVIEW.issueCertificatePowerSearch);
    await click(PKI_OVERVIEW.firstPowerSelectOption);
    await click(PKI_OVERVIEW.issueCertificateButton);
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.roles.role.generate');
  });

  test('navigates to certificate details page for View Certificates card', async function (assert) {
    await login(this.pkiAdminToken);
    await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
    await click(PKI_OVERVIEW.viewCertificatePowerSearch);
    await click(PKI_OVERVIEW.firstPowerSelectOption);
    await click(PKI_OVERVIEW.viewCertificateButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.certificates.certificate.details'
    );
  });

  test('navigates to issuer details page for View Issuer card', async function (assert) {
    await login(this.pkiAdminToken);
    await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
    await click(PKI_OVERVIEW.viewIssuerPowerSearch);
    await click(PKI_OVERVIEW.firstPowerSelectOption);
    await click(PKI_OVERVIEW.viewIssuerButton);
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.issuers.issuer.details');
  });
});
