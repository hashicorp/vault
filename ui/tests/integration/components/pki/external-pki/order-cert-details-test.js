/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { findAll, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | ExternalPki::OrderCertDetails', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.orderId = 'test-order-123';
    this.engineId = 'pki-external';
    this.certificate = { details: undefined };
    this.order = { details: undefined };

    this.renderComponent = () =>
      render(
        hbs`<ExternalPki::OrderCertDetails
          @order={{this.order}}
          @certificate={{this.certificate}}
          @orderId={{this.orderId}}
          @engineId={{this.engineId}}
        />`,
        { owner: this.engine }
      );
  });

  // Rendering order challenges table

  test('it renders empty state when no challenges provided', async function (assert) {
    this.order = { details: { challenges: null } };
    await this.renderComponent();
    assert.dom(GENERAL.tableRow()).doesNotExist('no table rows rendered without order data');
    assert.dom(GENERAL.cardContainer('Order information')).doesNotExist();
  });

  test('it renders table with single identifier and single challenge', async function (assert) {
    this.order = {
      details: {
        challenges: {
          'example.com': [
            {
              challenge_status: 'pending',
              challenge_type: 'dns-01',
              expires: '2026-07-24T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
          ],
        },
      },
    };

    await this.renderComponent();
    assert.dom(GENERAL.cardContainer('Order information')).exists();
    assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'renders one table row');
    assert
      .dom(GENERAL.tableData(0, 'identifier'))
      .hasText('Toggle example.com', 'displays identifier as expandable row');
    assert.dom(GENERAL.tableData(0, 'challenge_status')).hasText('Pending', 'displays pending status badge');
    assert.dom(GENERAL.tableData(0, 'challenge_type')).hasText('', 'no challenge type shown for pending');
    assert
      .dom(`${GENERAL.tableData(0, 'identifier')} button`)
      .hasAttribute('aria-expanded', 'true', 'nested rows are open by default');

    assert
      .dom(GENERAL.badge('challenge_status'))
      .exists({ count: 2 }, 'it renders a badge for parent and nested row');

    // Nested rows conveniently don't get their own row index number and are not nested within the parent's index
    // so find all of an element and get the second item
    const [, childStatus] = findAll('[data-test-table-data="challenge_status"]');
    const [, childType] = findAll('[data-test-table-data="challenge_type"]');
    const [, childExpires] = findAll('[data-test-table-data="expires"]');
    assert.dom(childStatus).hasText('Pending', 'challenge status is pending');
    assert.dom(childType).hasText('DNS-01', 'challenge type is uppercase');
    assert.dom(childExpires).hasTextContaining('07/24/2026');
  });

  test('it renders table with single identifier and multiple challenges', async function (assert) {
    this.order = {
      details: {
        challenges: {
          'example.com': [
            {
              challenge_status: 'valid',
              challenge_type: 'dns-01',
              expires: '2026-07-24T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
            {
              challenge_status: 'pending',
              challenge_type: 'http-01',
              expires: '2026-07-25T21:34:36Z',
              requires_manual_fulfillment: 'true',
            },
          ],
        },
      },
    };

    await this.renderComponent();
    assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'renders one table row for identifier');
    assert.dom(GENERAL.tableData(0, 'identifier')).hasText('Toggle example.com', 'displays identifier');
    assert
      .dom(GENERAL.tableData(0, 'challenge_status'))
      .hasText('Valid', 'displays valid status when at least one challenge is valid');
    assert
      .dom(GENERAL.tableData(0, 'challenge_type'))
      .hasText('DNS-01', 'displays only valid challenge types in uppercase');

    const [, firstChildStatus, secondChildStatus] = findAll('[data-test-table-data="challenge_status"]');
    const [, firstChildType, secondChildType] = findAll('[data-test-table-data="challenge_type"]');
    const [, firstChildExpires, secondChildExpires] = findAll('[data-test-table-data="expires"]');

    assert.dom(firstChildStatus).hasText('Valid', 'first child has "Valid" status');
    assert.dom(firstChildType).hasText('DNS-01', 'first child displays type');
    assert.dom(firstChildExpires).hasTextContaining('07/24/2026', 'first child displays expiry');

    assert.dom(secondChildStatus).hasText('Pending', 'second child has "Pending" status');
    assert.dom(secondChildType).hasText('HTTP-01', 'second child displays type');
    assert.dom(secondChildExpires).hasTextContaining('07/25/2026', 'second child displays expiry');
  });

  test('it renders table with multiple identifiers', async function (assert) {
    this.order = {
      details: {
        challenges: {
          'example.com': [
            {
              challenge_status: 'valid',
              challenge_type: 'dns-01',
              expires: '2026-07-24T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
          ],
          'test.example.com': [
            {
              challenge_status: 'pending',
              challenge_type: 'http-01',
              expires: '2026-07-25T21:34:36Z',
              requires_manual_fulfillment: 'true',
            },
          ],
        },
      },
    };

    await this.renderComponent();
    assert.dom(GENERAL.tableRow()).exists({ count: 2 }, 'renders two table rows for two identifiers');
    assert.dom(GENERAL.tableData(0, 'identifier')).hasText('Toggle example.com');
    assert.dom(GENERAL.tableData(1, 'identifier')).hasText('Toggle test.example.com');
  });

  test('it shows valid status when multiple challenges are valid', async function (assert) {
    this.order = {
      details: {
        challenges: {
          'example.com': [
            {
              challenge_status: 'valid',
              challenge_type: 'dns-01',
              expires: '2026-07-24T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
            {
              challenge_status: 'valid',
              challenge_type: 'http-01',
              expires: '2026-07-25T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
          ],
        },
      },
    };

    await this.renderComponent();

    assert.dom(GENERAL.tableData(0, 'challenge_status')).hasText('Valid', 'displays valid status');
    assert
      .dom(GENERAL.tableData(0, 'challenge_type'))
      .hasText('DNS-01, HTTP-01', 'displays all valid challenge types comma-separated and uppercase');
  });

  test('it shows pending status when no challenges are valid', async function (assert) {
    this.order = {
      details: {
        challenges: {
          'example.com': [
            {
              challenge_status: 'pending',
              challenge_type: 'dns-01',
              expires: '2026-07-24T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
            {
              challenge_status: 'invalid',
              challenge_type: 'http-01',
              expires: '2026-07-25T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
          ],
        },
      },
    };

    await this.renderComponent();

    assert
      .dom(GENERAL.tableData(0, 'challenge_status'))
      .containsText('Pending', 'displays pending status when no valid challenges');
    assert
      .dom(GENERAL.tableData(0, 'challenge_type'))
      .hasText('', 'displays empty challenge type when no valid challenges');
  });

  test('it formats challenge types to uppercase', async function (assert) {
    this.order = {
      details: {
        challenges: {
          'example.com': [
            {
              challenge_status: 'valid',
              challenge_type: 'dns-01',
              expires: '2026-07-24T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
          ],
          'test.com': [
            {
              challenge_status: 'valid',
              challenge_type: 'http-01',
              expires: '2026-07-24T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
          ],
          'another.com': [
            {
              challenge_status: 'valid',
              challenge_type: 'tls-alpn-01',
              expires: '2026-07-24T21:34:36Z',
              requires_manual_fulfillment: 'false',
            },
          ],
        },
      },
    };

    await this.renderComponent();

    assert.dom(GENERAL.tableData(0, 'challenge_type')).hasText('DNS-01', 'dns-01 formatted to DNS-01');
    assert.dom(GENERAL.tableData(1, 'challenge_type')).hasText('HTTP-01', 'http-01 formatted to HTTP-01');
    assert
      .dom(GENERAL.tableData(2, 'challenge_type'))
      .hasText('TLS-ALPN-01', 'tls-alpn-01 formatted to TLS-ALPN-01');
  });

  // Rendering order details

  test('completed order: only shows order_status, last_order_update, role_name in order section', async function (assert) {
    this.order = {
      details: {
        order_status: 'completed',
        last_update: '2026-07-20T10:00:00Z',
        role_name: 'my-role',
        creation_date: '2026-07-01T10:00:00Z',
        expires: '2026-08-01T10:00:00Z',
        challenges: {},
      },
    };

    await this.renderComponent();

    assert.dom(GENERAL.infoRowLabel('Order status')).exists();
    assert.dom(GENERAL.infoRowLabel('Last order update')).exists();
    assert.dom(GENERAL.infoRowLabel('Role name')).exists();
    assert
      .dom(GENERAL.infoRowLabel('Creation date'))
      .doesNotExist('creation_date not shown for completed orders');
    assert.dom(GENERAL.infoRowLabel('Expires')).doesNotExist('expires not shown for completed orders');
  });

  test('expired order: only shows order_status, last_order_update, creation_date, expires, role_name', async function (assert) {
    this.order = {
      details: {
        order_status: 'expired',
        last_update: '2026-07-20T10:00:00Z',
        creation_date: '2026-07-01T10:00:00Z',
        expires: '2026-08-01T10:00:00Z',
        next_work_date: '0001-01-01T00:00:00Z',
        role_name: 'my-role',
        challenges: {},
      },
    };

    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('Order ID')).exists();
    assert.dom(GENERAL.infoRowLabel('Order status')).exists();
    assert.dom(GENERAL.infoRowLabel('Last order update')).exists();
    assert.dom(GENERAL.infoRowLabel('Creation date')).exists();
    assert.dom(GENERAL.infoRowLabel('Expires')).exists();
    assert.dom(GENERAL.infoRowLabel('Role name')).exists();
    assert.dom(GENERAL.infoRowLabel('Next work date')).doesNotExist();
    assert.dom(GENERAL.messageError).doesNotExist();
  });

  test('failed order: only shows order_status, last_order_update, creation_date, expires, role_name', async function (assert) {
    const error =
      'failed to get authorization for order 019f8c25-fc43-75bc-8cf5-ca1e0328a7bc: 1 error occurred:\n\t* authorization status is Invalid for identifier dns:host.docker.internal within order 019f8c25-fc43-75bc-8cf5-ca1e0328a7bc: 403 urn:ietf:params:acme:error:unauthorized: Non-200 status code from HTTP: http://host.docker.internal:5002/.well-known/acme-challenge/Adwp5HWy-ck82EvMthHvvq_Fxsgs5IK3Firu0GYFuA4 returned 404\n\n';
    this.order = {
      details: {
        challenges: {
          'host.docker.internal': [
            {
              challenge_status: 'invalid',
              challenge_type: 'http-01',
              error:
                '403 urn:ietf:params:acme:error:unauthorized: Non-200 status code from HTTP: http://host.docker.internal:5002/.well-known/acme-challenge/Adwp5HWy-ck82EvMthHvvq_Fxsgs5IK3Firu0GYFuA4 returned 404',
              expires: '2026-07-23T00:25:28Z',
              requires_manual_fulfillment: 'true',
            },
          ],
        },
        creation_date: '2026-07-22T16:25:27-07:00',
        csr: '',
        expires: '2026-07-23T23:25:28Z',
        identifiers: ['host.docker.internal'],
        last_error: error,
        last_update: '2026-07-22T16:30:49-07:00',
        next_work_date: '0001-01-01T00:00:00Z',
        order_status: 'error',
        role_name: 'pebble-july-3',
        serial_number: '',
      },
    };

    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('Order ID')).exists();
    assert.dom(GENERAL.infoRowLabel('Order status')).exists();
    assert.dom(GENERAL.infoRowLabel('Last order update')).exists();
    assert.dom(GENERAL.infoRowLabel('Creation date')).exists();
    assert.dom(GENERAL.infoRowLabel('Expires')).exists();
    assert.dom(GENERAL.infoRowLabel('Role name')).exists();
    assert.dom(GENERAL.infoRowLabel('Next work date')).doesNotExist();
    assert.dom(GENERAL.messageError).exists().hasText(`Error ${error}`);
  });

  test('pending order: shows all order fields including order_id', async function (assert) {
    this.order = {
      details: {
        order_status: 'awaiting-challenge-fulfillment',
        role_name: 'my-role',
        challenges: {
          'host.docker.internal': [
            {
              challenge_status: 'pending',
              challenge_type: 'tls-alpn-01',
              expires: '2026-07-23T00:40:18Z',
              requires_manual_fulfillment: 'true',
            },
          ],
        },
        creation_date: '2026-07-22T16:40:17-07:00',
        csr: '',
        expires: '2026-07-23T23:40:18Z',
        identifiers: ['host.docker.internal'],
        last_error: '',
        last_update: '2026-07-23T10:44:07-07:00',
        next_work_date: '2026-07-23T11:44:07-07:00',
        serial_number: '',
      },
    };

    await this.renderComponent();

    assert.dom(GENERAL.infoRowLabel('Order ID')).exists();
    assert.dom(GENERAL.infoRowLabel('Order status')).exists();
    assert.dom(GENERAL.infoRowLabel('Next work date')).exists();
  });

  // Rendering certificate details

  test('certificate card renders when certificate details are present', async function (assert) {
    this.order = { details: { order_status: 'completed', challenges: {} } };
    this.certificate = {
      details: {
        serial_number: 'ab:cd:ef',
        certificate: '-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----',
      },
    };

    await this.renderComponent();

    assert.dom(GENERAL.cardContainer('Certificate details')).exists('certificate card is rendered');
    assert.dom(GENERAL.infoRowLabel('Serial number')).exists();
  });

  test('certificate card does not without certificate details', async function (assert) {
    this.order = { details: { order_status: 'pending', challenges: {} } };
    this.certificate = { details: undefined };

    await this.renderComponent();
    assert.dom(GENERAL.cardContainer('Certificate details')).doesNotExist();
  });

  // Error states

  test('order 404 error: renders alert title with no message body', async function (assert) {
    this.order = {
      error: { status: 404 },
    };

    await this.renderComponent();

    assert.dom(GENERAL.messageError).exists().hasText('Order status is unavailable');
  });

  test('order 500 error: renders API error message', async function (assert) {
    this.order = {
      error: { status: 500, message: 'internal server error' },
    };
    this.certificate = {
      details: {
        serial_number: 'ab:cd:ef',
        certificate: '-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----',
      },
    };
    await this.renderComponent();
    assert.dom(GENERAL.cardContainer('Certificate details')).exists();
    assert.dom(GENERAL.messageError).exists().hasText('Order status is unavailable internal server error');
  });

  test('order 403 error: renders permission error', async function (assert) {
    this.order = {
      error: { status: 403, message: 'permission denied', path: '/v1/pki/roles/my-role/order/123' },
    };
    this.certificate = {
      details: {
        serial_number: 'ab:cd:ef',
        certificate: '-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----',
      },
    };
    await this.renderComponent();
    assert
      .dom(GENERAL.messageError)
      .exists()
      .hasText(
        'Order status is unavailable You do not have "read" permissions for the path: /v1/pki/roles/my-role/order/123'
      );
    assert.dom(GENERAL.cardContainer('Certificate details')).exists('certificate card still renders');
    assert.dom(GENERAL.infoRowLabel('Serial number')).exists('cert details are visible');
  });

  test('certificate 400 error: renders API error message', async function (assert) {
    this.order = { details: { order_status: 'expired', challenges: {} } };
    this.certificate = {
      error: { status: 400, message: 'order has status expired, must be completed to fetch cert' },
    };

    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('Order status')).exists('order section still renders');
    assert.dom(GENERAL.cardContainer('Certificate details')).doesNotExist();
    assert
      .dom(GENERAL.messageError)
      .exists()
      .hasText('Certificate data is unavailable order has status expired, must be completed to fetch cert');
  });

  test('certificate 404 error: renders alert title with no message body', async function (assert) {
    this.order = { details: { order_status: 'completed', challenges: {} } };
    this.certificate = {
      error: { status: 404 },
    };

    await this.renderComponent();

    assert.dom(GENERAL.messageError).exists().hasText('Certificate data is unavailable');
  });

  test('certificate 403 error: renders permission error', async function (assert) {
    this.order = { details: { order_status: 'expired', challenges: {} } };
    this.certificate = {
      error: {
        status: 403,
        message: 'permission denied',
        path: '/v1/pki/roles/my-role/myorderid/fetch-cert',
      },
    };

    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('Order status')).exists('order section still renders');
    assert.dom(GENERAL.cardContainer('Certificate details')).doesNotExist();
    assert
      .dom(GENERAL.messageError)
      .exists()
      .hasText(
        'Certificate data is unavailable You do not have "read" permissions for the path: /v1/pki/roles/my-role/myorderid/fetch-cert'
      );
  });
});
