/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | pki | external-pki | ExternalPki::Page::DnsProviders', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.setupModel = (dnsProviders = []) => {
      return {
        engine: new SecretsEngineResource({
          accessor: 'pki-external-ca_e158c567',
          type: 'pki-external-ca',
          path: 'my-pki-external-ca/',
        }),
        dnsProviders,
      };
    };
    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'scrollable-region-focusable': { enabled: false },
      },
    });

    this.renderComponent = () =>
      render(hbs`<ExternalPki::Page::DnsProviders @model={{this.model}} />`, { owner: this.engine });
  });

  test('it renders empty state when no DNS providers exist', async function (assert) {
    this.model = this.setupModel([]);
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No DNS providers exist yet');
  });

  module('with DNS providers', function (hooks) {
    hooks.beforeEach(function () {
      this.dnsProviders = [
        {
          name: 'aws-route53-prod',
          type: 'aws-route53',
          access_key_id: 'AKIAIOSFODNN7EXAMPLE',
          region: 'us-east-1',
          hosted_zone_id: 'Z3M3LMPEXAMPLE',
        },
        {
          name: 'azure-dns-staging',
          type: 'azure',
          client_id: '12345678-1234',
          tenant_id: '87654321-4321-4321-4321-210987654321',
          subscription_id: 'abcdef12-3456',
          resource_group_name: 'my-resource-group',
        },
        {
          name: 'gcp-dns-dev',
          type: 'google-cloud-dns',
          project: 'my-gcp-project',
          service_account_key: '{"type":"service_account","project_id":"my-project"}',
        },
        {
          name: 'rfc2136-internal',
          type: 'rfc2136',
          nameserver: '192.168.1.1:53',
          tsig_key_name: 'example-key',
          tsig_algorithm: 'hmac-sha256',
        },
      ];
      this.model = this.setupModel(this.dnsProviders);
    });

    test('it renders list of DNS providers', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.emptyStateTitle).doesNotExist();
      assert.dom(GENERAL.cardContainer()).exists({ count: 4 });
      const provider = ['aws-route53-prod', 'azure-dns-staging', 'gcp-dns-dev', 'rfc2136-internal'];
      provider.forEach((p) => {
        assert.dom(`${GENERAL.cardContainer(p)} ${GENERAL.textDisplay()}`).hasText(p);
      });
    });

    test('it displays configuration details', async function (assert) {
      await this.renderComponent();
      const cardConfig = (provider, label) =>
        `${GENERAL.cardContainer(provider)} ${GENERAL.infoRowValue(label)}`;

      assert.dom(cardConfig('aws-route53-prod', 'Access key ID')).containsText('AKIAIOSFODNN7EXAMPLE');
      assert.dom(cardConfig('aws-route53-prod', 'Region')).containsText('us-east-1');
      assert.dom(cardConfig('aws-route53-prod', 'Hosted zone ID')).containsText('Z3M3LMPEXAMPLE');

      assert.dom(cardConfig('azure-dns-staging', 'Client ID')).containsText('12345678-1234');
      assert.dom(cardConfig('azure-dns-staging', 'Subscription ID')).containsText('abcdef12-3456');
      assert.dom(cardConfig('azure-dns-staging', 'Resource group name')).containsText('my-resource-group');

      assert.dom(cardConfig('gcp-dns-dev', 'Project')).containsText('my-gcp-project');
      assert
        .dom(cardConfig('gcp-dns-dev', 'Service account key'))
        .containsText('{"type":"service_account","project_id":"my-project"}');

      assert.dom(cardConfig('rfc2136-internal', 'TSIG key name')).containsText('example-key');
      assert.dom(cardConfig('rfc2136-internal', 'TSIG algorithm')).containsText('hmac-sha256');
    });
  });
});
