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

module('Integration | Component | ExternalPki::OrderInfoCard', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.order = undefined;
    this.hasBorder = undefined;
    this.renderComponent = () =>
      render(hbs`<ExternalPki::OrderInfoCard @order={{this.order}} @hasBorder={{this.hasBorder}} />`, {
        owner: this.engine,
      });
  });

  test('it renders empty state when no order provided', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No order challenges to display');
    assert.dom(GENERAL.tableRow()).doesNotExist('no table rows rendered without order data');
  });

  test('it renders empty state when order has no challenges', async function (assert) {
    this.order = { challenges: null };
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No order challenges to display');
    assert.dom(GENERAL.tableRow()).doesNotExist('no table rows rendered without order data');
  });

  test('it renders table with single identifier and single challenge', async function (assert) {
    this.order = {
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
    };

    await this.renderComponent();
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
    };

    await this.renderComponent();
    assert.dom(GENERAL.tableRow()).exists({ count: 2 }, 'renders two table rows for two identifiers');
    assert.dom(GENERAL.tableData(0, 'identifier')).hasText('Toggle example.com');
    assert.dom(GENERAL.tableData(1, 'identifier')).hasText('Toggle test.example.com');
  });

  test('it shows valid status when multiple challenges are valid', async function (assert) {
    this.order = {
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
    };

    await this.renderComponent();

    assert.dom(GENERAL.tableData(0, 'challenge_status')).hasText('Valid', 'displays valid status');
    assert
      .dom(GENERAL.tableData(0, 'challenge_type'))
      .hasText('DNS-01, HTTP-01', 'displays all valid challenge types comma-separated and uppercase');
  });

  test('it shows pending status when no challenges are valid', async function (assert) {
    this.order = {
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
    };

    await this.renderComponent();

    assert.dom(GENERAL.tableData(0, 'challenge_type')).hasText('DNS-01', 'dns-01 formatted to DNS-01');
    assert.dom(GENERAL.tableData(1, 'challenge_type')).hasText('HTTP-01', 'http-01 formatted to HTTP-01');
    assert
      .dom(GENERAL.tableData(2, 'challenge_type'))
      .hasText('TLS-ALPN-01', 'tls-alpn-01 formatted to TLS-ALPN-01');
  });

  test('it yields orderDetails block when provided', async function (assert) {
    this.order = {
      challenges: {
        'example.com': [
          {
            challenge_status: 'valid',
            challenge_type: 'dns-01',
            expires: '2026-07-24T21:34:36Z',
            requires_manual_fulfillment: 'false',
          },
        ],
      },
    };

    await render(
      hbs`
      <ExternalPki::OrderInfoCard @order={{this.order}}>
        <:orderDetails>
          <div data-test-custom-details>Custom order details</div>
        </:orderDetails>
      </ExternalPki::OrderInfoCard>
    `,
      { owner: this.engine }
    );

    assert.dom('[data-test-custom-details]').hasText('Custom order details', 'yields orderDetails block');
  });
});
