/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { find, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | pki | external-pki | ExternalPki::ConfigDetails', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.renderComponent = () =>
      render(hbs`<ExternalPki::ConfigDetails @config={{this.config}} @engineId={{this.engineId}} />`, {
        owner: this.engine,
      });
  });

  test('it handles undefined config', async function (assert) {
    this.config = undefined;

    await this.renderComponent();
    assert.dom('[data-test-row-label]').doesNotExist();
  });

  test('it handles empty configuration', async function (assert) {
    this.config = {};

    await this.renderComponent();
    assert.dom('[data-test-row-label]').doesNotExist();
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
      identification_station: 'example.com',
      testid: 'test',
    };

    await this.renderComponent();
    // Should transform "Id" to "ID" in labels
    assert.dom(GENERAL.infoRowLabel('ID start')).exists('it formats labels that start with id_');
    assert.dom(GENERAL.infoRowLabel('Client ID key')).exists('it formats labels with _id_ in the middle');
    assert.dom(GENERAL.infoRowLabel('Tenant ID')).exists('it formats labels that end with _id');
    assert
      .dom(GENERAL.infoRowLabel('Identification station'))
      .exists('it does not format when "id" begins the word');
    assert.dom(GENERAL.infoRowLabel('Testid')).exists('it does not format when "id" is part of the word');
  });

  test('it renders EncodedDataCard for trusted_ca, ca_chain, certificate, and private_key fields', async function (assert) {
    const certData =
      '-----BEGIN CERTIFICATE-----\nMIIBkTCB+wIJAKHHCgVZU1WOMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnZh\n-----END CERTIFICATE-----';

    this.config = {
      trusted_ca: certData,
      ca_chain: certData,
      certificate: certData,
      private_key: certData,
    };

    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('Trusted CA')).exists('trusted_ca renders formatted label');
    assert.dom(GENERAL.infoRowLabel('CA chain')).exists('ca_chain renders formatted label');
    assert.dom(GENERAL.infoRowLabel('Certificate')).exists('certificate renders formatted label');
    assert.dom(GENERAL.infoRowLabel('Private key')).exists('private_key renders formatted label');
    assert
      .dom('[data-test-certificate-card]')
      .exists({ count: 4 }, 'EncodedDataCard renders for all four fields');
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
      ca_chain: 'CA chain',
      directory_url: 'https://acme-v02.api.letsencrypt.org/directory',
      key_type: 'ec-256',
      nameserver: 'internal.host.server',
      not_after: '2025-01-01T00:00:00Z',
      not_before: '2024-01-01T00:00:00Z',
      trusted_ca: 'Trusted CA',
      tsig_algorithm: 'TSIG algorithm',
      tsig_key_name: 'TSIG key name',
      ttl: 2355,
    };

    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('IAM role ARN to assume')).exists();
    assert.dom(GENERAL.infoRowLabel('CA chain')).exists();
    assert.dom(GENERAL.infoRowLabel('Directory URL')).exists();
    assert.dom(GENERAL.infoRowLabel('Active key type')).exists();
    assert.dom(GENERAL.infoRowLabel('DNS server address')).exists();
    assert.dom(GENERAL.infoRowLabel('Valid until')).exists('not_after renders as "Valid until"');
    assert.dom(GENERAL.infoRowLabel('Valid after')).exists('not_before renders as "Valid after"');
    assert.dom(GENERAL.infoRowLabel('Trusted CA')).exists();
    assert.dom(GENERAL.infoRowLabel('TSIG algorithm')).exists();
    assert.dom(GENERAL.infoRowLabel('TSIG key name')).exists();
    assert.dom(GENERAL.infoRowLabel('Time to live')).exists();
    assert.dom(GENERAL.messageError).doesNotExist('error banner is not rendered without last_error');
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

  test('it renders error banner when last_error is present', async function (assert) {
    this.config = {
      directory_url: 'https://acme.example.com/directory',
      last_error: 'ACME challenge failed: DNS record not found',
    };

    await this.renderComponent();
    assert.dom(GENERAL.messageError).exists('renders error banner when last_error is set');
    assert.dom(GENERAL.messageError).containsText('ACME challenge failed: DNS record not found');
    // last_error itself is excluded from the info table
    assert
      .dom(GENERAL.infoRowLabel('Last error'))
      .doesNotExist('last_error field is excluded from the table');
  });

  test('it renders a copy button for order_id and serial_number fields', async function (assert) {
    this.config = {
      order_id: 'abc-123-order',
      serial_number: '12:34:56:78',
      directory_url: 'https://acme.example.com/directory',
    };

    await this.renderComponent();

    assert
      .dom(GENERAL.copyButton, find(GENERAL.infoRowValue('Order ID')))
      .exists('copy button renders for order_id');
    assert
      .dom(GENERAL.copyButton, find(GENERAL.infoRowValue('Serial number')))
      .exists('copy button renders for serial_number');
    assert
      .dom(GENERAL.copyButton, find(GENERAL.infoRowValue('Directory URL')))
      .doesNotExist('copy button is not rendered for non-copyable fields');
  });

  test('it formats date fields as UTC using MM/dd/yyyy, HH:mm format', async function (assert) {
    const isoDate = '2024-03-15T14:30:00Z';
    const expectedDate = '03/15/2024, 14:30 UTC';

    this.config = {
      creation_date: isoDate,
      last_update: isoDate,
      last_updated: isoDate,
      last_updated_date: isoDate,
      last_update_date: isoDate,
      next_work_date: isoDate,
      expires: isoDate,
      not_after: isoDate,
      not_before: isoDate,
    };

    await this.renderComponent();

    assert
      .dom(GENERAL.infoRowValue('Creation date'))
      .hasText(expectedDate, 'creation_date is formatted as UTC');
    assert.dom(GENERAL.infoRowValue('Last update')).hasText(expectedDate, 'last_update is formatted as UTC');
    assert
      .dom(GENERAL.infoRowValue('Last updated'))
      .hasText(expectedDate, 'last_updated is formatted as UTC');
    assert
      .dom(GENERAL.infoRowValue('Last updated date'))
      .hasText(expectedDate, 'last_updated_date is formatted as UTC');
    assert
      .dom(GENERAL.infoRowValue('Last update date'))
      .hasText(expectedDate, 'last_update_date is formatted as UTC');
    assert
      .dom(GENERAL.infoRowValue('Next work date'))
      .hasText(expectedDate, 'next_work_date is formatted as UTC');
    assert.dom(GENERAL.infoRowValue('Expires')).hasText(expectedDate, 'expires is formatted as UTC');
    assert.dom(GENERAL.infoRowValue('Valid until')).hasText(expectedDate, 'not_after is formatted as UTC');
    assert.dom(GENERAL.infoRowValue('Valid after')).hasText(expectedDate, 'not_before is formatted as UTC');
  });

  test('it renders a link to the role if @engineId is provided', async function (assert) {
    this.engineId = 'mypkimount';
    this.config = {
      order_id: 'abc-123-order',
      role_name: 'myrole',
      directory_url: 'https://acme.example.com/directory',
    };

    await this.renderComponent();

    assert.dom(GENERAL.infoRowValue('Role name')).exists().hasText('myrole');
    assert.dom(GENERAL.linkTo('myrole')).exists().hasText('myrole');
  });
});
