/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { currentRouteName, currentURL, visit } from '@ember/test-helpers';
import sinon from 'sinon';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Acceptance | enterprise | pki | external | dns-providers route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
    this.dnsProvidersListStub = sinon.stub(api.secrets, 'pkiExternalCaListConfigDns');
    this.awsRoute53ReadStub = sinon.stub(api.secrets, 'pkiExternalCaReadConfigDnsAwsRoute53');
    this.azureReadStub = sinon.stub(api.secrets, 'pkiExternalCaReadConfigDnsAzureDns');
    this.gcpReadStub = sinon.stub(api.secrets, 'pkiExternalCaReadConfigDnsGoogleCloudDns');
    this.rfc2136ReadStub = sinon.stub(api.secrets, 'pkiExternalCaReadConfigDnsRfc2136');
    this.mountPath = `pki-external-ca-${uuidv4()}`;

    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.dnsProvidersURL = `/vault/secrets-engines/${this.mountPath}/pki/external/dns-providers`;
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it fetches DNS provider details for each provider type', async function (assert) {
    this.dnsProvidersListStub.resolves({
      keys: ['aws-prod', 'azure-staging', 'gcp-dev', 'rfc2136-internal'],
      key_info: {
        'aws-prod': { name: 'aws-prod', type: 'aws-route53' },
        'azure-staging': { name: 'azure-staging', type: 'azure' },
        'gcp-dev': { name: 'gcp-dev', type: 'google-cloud-dns' },
        'rfc2136-internal': { name: 'rfc2136-internal', type: 'rfc2136' },
      },
    });

    // Mock individual provider reads
    this.awsRoute53ReadStub.withArgs('aws-prod', this.mountPath).resolves({
      name: 'aws-prod',
      type: 'aws-route53',
      access_key_id: 'AKIAIOSFODNN7EXAMPLE',
      region: 'us-east-1',
      hosted_zone_id: 'Z3M3LMPEXAMPLE',
    });

    this.azureReadStub.withArgs('azure-staging', this.mountPath).resolves({
      name: 'azure-staging',
      type: 'azure',
      client_id: '12345678-1234-1234-1234-123456789012',
      tenant_id: '87654321-4321-4321-4321-210987654321',
      subscription_id: 'abcdef12-3456-7890-abcd-ef1234567890',
      resource_group_name: 'my-resource-group',
    });

    this.gcpReadStub.withArgs('gcp-dev', this.mountPath).resolves({
      name: 'gcp-dev',
      type: 'google-cloud-dns',
      project: 'my-gcp-project',
      service_account_key: '{"type":"service_account"}',
    });

    this.rfc2136ReadStub.withArgs('rfc2136-internal', this.mountPath).resolves({
      name: 'rfc2136-internal',
      type: 'rfc2136',
      nameserver: '192.168.1.1:53',
      tsig_key_name: 'example-key',
      tsig_algorithm: 'hmac-sha256',
    });

    await visit(this.dnsProvidersURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.dns-providers',
      'navigated to dns-providers route'
    );
    assert.true(this.dnsProvidersListStub.calledOnce, 'providers list called once in parent route');
    assert.true(this.awsRoute53ReadStub.calledOnce, 'AWS Route53 read called once');
    assert.true(this.azureReadStub.calledOnce, 'Azure read called once');
    assert.true(this.gcpReadStub.calledOnce, 'GCP read called once');
    assert.true(this.rfc2136ReadStub.calledOnce, 'RFC2136 read called once');

    // Verify providers cards display
    assert.dom(GENERAL.cardContainer()).exists({ count: 4 });
    const provider = ['aws-prod', 'azure-staging', 'gcp-dev', 'rfc2136-internal'];
    provider.forEach((p) => {
      assert.dom(`${GENERAL.cardContainer(p)} ${GENERAL.textDisplay()}`).hasText(p);
    });
  });

  test('it handles a 404', async function (assert) {
    this.dnsProvidersListStub.rejects(getErrorResponse()); // Throws 404
    await visit(this.dnsProvidersURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.dns-providers',
      'navigated to dns-providers route'
    );
    assert.dom(GENERAL.emptyStateTitle).hasText('No DNS providers exist yet');
  });

  test('it displays 403 permission denied error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.dnsProvidersListStub.rejects(getErrorResponse(error, 403));
    await visit(this.dnsProvidersURL);
    assert.true(this.dnsProvidersListStub.calledOnce, 'providers list called once');
    assert.strictEqual(currentURL(), this.dnsProvidersURL, 'it renders dns-providers URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'it redirects to error route'
    );
    assert.dom(GENERAL.pageError.title(403)).exists().hasText('ERROR 403 Not authorized');
  });

  test('it displays 500 internal server error', async function (assert) {
    const error = { errors: ['Internal server error'] };
    this.dnsProvidersListStub.rejects(getErrorResponse(error, 500));
    await visit(this.dnsProvidersURL);
    assert.true(this.dnsProvidersListStub.calledOnce, 'providers list called once');
    assert.strictEqual(currentURL(), this.dnsProvidersURL, 'it renders dns-providers URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'it redirects to error route'
    );
    assert.dom(GENERAL.pageError.title(500)).exists().hasText('ERROR 500 Error');
  });

  test('it throws error for unsupported DNS provider type', async function (assert) {
    this.dnsProvidersListStub.resolves({
      keys: ['invalid-provider'],
      key_info: {
        'invalid-provider': { type: 'unsupported-type' },
      },
    });
    await visit(this.dnsProvidersURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'it redirects to error route for unsupported type'
    );
    assert.dom(GENERAL.pageError.description).containsText('Unsupported DNS type');
  });

  test('it handles partial failures when reading individual providers', async function (assert) {
    this.dnsProvidersListStub.resolves({
      keys: ['aws-prod', 'azure-staging'],
      key_info: {
        'aws-prod': { type: 'aws-route53' },
        'azure-staging': { type: 'azure' },
      },
    });

    // First provider succeeds
    this.awsRoute53ReadStub.withArgs('aws-prod', this.mountPath).resolves({
      name: 'aws-prod',
      type: 'aws-route53',
      access_key_id: 'AKIAIOSFODNN7EXAMPLE',
      region: 'us-east-1',
    });

    // Second provider fails
    this.azureReadStub
      .withArgs('azure-staging', this.mountPath)
      .rejects(getErrorResponse({ errors: ['1 error occurred:\n\t* permission denied\n\n'] }, 403));

    await visit(this.dnsProvidersURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.dns-providers',
      'stays on dns-providers route despite partial failure'
    );
    assert.true(this.awsRoute53ReadStub.calledOnce, 'AWS Route53 read called once');
    assert.true(this.azureReadStub.calledOnce, 'Azure read called once');
    // Successful account displays
    assert
      .dom(`${GENERAL.cardContainer('aws-prod')} ${GENERAL.infoRowValue('Access key ID')}`)
      .hasText('AKIAIOSFODNN7EXAMPLE');
    // Failed provider should show error
    assert
      .dom(`${GENERAL.cardContainer('azure-staging')} ${GENERAL.infoRowValue('Error')}`)
      .hasText('You do not have permission to read configurations for this provider');
  });
});
