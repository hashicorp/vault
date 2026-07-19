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

module('Integration | Component | pki | external-pki | ExternalPki::ConfigDetails', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.renderComponent = () =>
      render(hbs`<ExternalPki::ConfigDetails @config={{this.config}} />`, { owner: this.engine });
  });

  test('it renders configuration parameters and excludes specified fields', async function (assert) {
    this.config = {
      name: 'test-config',
      email_contacts: ['admin@example.com', 'security@example.com'],
      active_key_version: 1,
      account_keys: {
        1: { key_version: 1, key_type: 'EC256' },
      },
    };

    await this.renderComponent();
    assert.dom(GENERAL.infoRowValue('Email contacts')).containsText('admin@example.com,security@example.com');
    assert.dom(GENERAL.infoRowValue('Active key version')).hasText('1');
    // Excluded
    assert.dom(GENERAL.infoRowLabel('Name')).doesNotExist();
    assert.dom(GENERAL.infoRowLabel('Account keys')).doesNotExist();
  });

  test('it transforms "Id" to "ID" in labels', async function (assert) {
    this.config = {
      id_start: '123',
      client_id_key: '12345678-1234-1234-1234-123456789012',
      tenant_id: '87654321-4321-4321-4321-210987654321',
      identifiers: 'example.com',
      testid: 'test',
    };

    await this.renderComponent();
    // Should transform "Id" to "ID" in labels
    assert.dom(GENERAL.infoRowLabel('ID start')).exists('it formats labels that start with id_');
    assert.dom(GENERAL.infoRowLabel('Client ID key')).exists('it formats labels with _id_ in the middle');
    assert.dom(GENERAL.infoRowLabel('Tenant ID')).exists('it formats labels that end with _id');
    assert.dom(GENERAL.infoRowLabel('Identifiers')).exists('it does not format when "id" begins the word');
    assert.dom(GENERAL.infoRowLabel('Testid')).exists('it does not format when "id" is part of the word');
  });

  test('it renders trusted_ca with EncodedDataCard', async function (assert) {
    this.config = {
      directory_url: 'https://acme-v02.api.letsencrypt.org/directory',
      trusted_ca:
        '-----BEGIN CERTIFICATE-----\nMIIBkTCB+wIJAKHHCgVZU1WOMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnZh\n-----END CERTIFICATE-----',
    };
    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('Trusted CA')).exists('it renders formatted label');
    assert.dom('[data-test-certificate-card]').exists();
  });

  test('it formats TTL values', async function (assert) {
    this.config = {
      directory_url: 'https://acme-v02.api.letsencrypt.org/directory',
      ttl: 2592000, // 30 days in seconds
    };

    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('Time to live')).exists('it renders custom label');
    assert.dom(GENERAL.infoRowValue('Time to live')).hasText('30 days');
  });

  test('it renders custom labels', async function (assert) {
    this.config = {
      assume_role_arn: 'my-role',
      directory_url: 'https://acme-v02.api.letsencrypt.org/directory',
      key_type: 'ec-256',
      nameserver: 'internal.host.server',
      trusted_ca: 'Trusted CA',
      tsig_algorithm: 'TSIG algorithm',
      tsig_key_name: 'TSIG key name',
      ttl: 2355,
    };

    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('IAM role ARN to assume')).exists();
    assert.dom(GENERAL.infoRowLabel('Directory URL')).exists();
    assert.dom(GENERAL.infoRowLabel('Active key type')).exists();
    assert.dom(GENERAL.infoRowLabel('DNS server address')).exists();
    assert.dom(GENERAL.infoRowLabel('Trusted CA')).exists();
    assert.dom(GENERAL.infoRowLabel('TSIG algorithm')).exists();
    assert.dom(GENERAL.infoRowLabel('TSIG key name')).exists();
    assert.dom(GENERAL.infoRowLabel('Time to live')).exists();
  });

  test('it handles empty configuration', async function (assert) {
    this.config = {};

    await this.renderComponent();
    assert.dom('[data-test-row-label]').doesNotExist();
  });

  test('it handles configuration with only excluded fields', async function (assert) {
    this.config = {
      name: 'test-config',
      account_keys: {
        1: { key_version: 1 },
      },
    };
    await this.renderComponent();
    assert.dom('[data-test-row-label]').doesNotExist();
  });
});
