/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit, click, fillIn, findAll, currentRouteName } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import oidcConfigHandlers from 'vault/mirage/handlers/oidc-config';
import authPage from 'vault/tests/pages/auth';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { selectChoose } from 'ember-power-select/test-support';
import ss from 'vault/tests/pages/components/search-select';
import fm from 'vault/tests/pages/components/flash-message';
import {
  OIDC_BASE_URL, // -> '/vault/access/oidc'
  SELECTORS,
  clearRecord,
  CLIENT_LIST_RESPONSE,
  CLIENT_DATA_RESPONSE,
} from 'vault/tests/helpers/oidc-config';
import { capabilitiesStub, overrideResponse } from 'vault/tests/helpers/stubs';

const searchSelect = create(ss);
const flashMessage = create(fm);

// in congruency with backend verbiage 'applications' are referred to as 'clients'
// throughout the codebase and the term 'applications' only appears in the UI

module('Acceptance | oidc-config clients and keys', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    oidcConfigHandlers(this.server);
    this.store = this.owner.lookup('service:store');
    return authPage.login();
  });

  test('it creates a key, signs a client and edits key access to only that client', async function (assert) {
    assert.expect(21);

    //* start with clean test state
    await clearRecord(this.store, 'oidc/client', 'client-with-test-key');
    await clearRecord(this.store, 'oidc/client', 'client-with-default-key');
    await clearRecord(this.store, 'oidc/key', 'test-key');

    // create client with default key
    await visit(OIDC_BASE_URL + '/clients/create');
    await fillIn('[data-test-input="name"]', 'client-with-default-key');
    await click(SELECTORS.clientSaveButton);

    // check reroutes from oidc index to clients index when client exists
    await visit(OIDC_BASE_URL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.index',
      'redirects to clients index route when clients exist'
    );
    assert.dom('[data-test-tab="clients"]').hasClass('active', 'clients tab is active');
    assert
      .dom('[data-test-oidc-client-linked-block]')
      .hasTextContaining('client-with-default-key', 'displays linked block for client');

    // navigate to keys
    await click('[data-test-tab="keys"]');
    assert.dom('[data-test-tab="keys"]').hasClass('active', 'keys tab is active');
    assert.strictEqual(currentRouteName(), 'vault.cluster.access.oidc.keys.index');
    assert
      .dom('[data-test-oidc-key-linked-block="default"]')
      .hasText('default', 'index page lists default key');

    // navigate to default key details from pop-up menu
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-key-menu-link="details"]');
    assert.dom(SELECTORS.keyDeleteButton).doesNotExist('delete button is hidden for default key');
    await click(SELECTORS.keyEditButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.edit',
      'navigates to edit from key details'
    );
    await click(SELECTORS.keyCancelButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.details',
      'key edit form navigates back to details on cancel'
    );
    await click(SELECTORS.keyClientsTab);
    assert
      .dom('[data-test-oidc-client-linked-block="client-with-default-key"]')
      .exists('lists correct app with default');

    // create a new key
    await click('[data-test-breadcrumb-link="oidc-keys"] a');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.index',
      'keys breadcrumb navigates back to list view'
    );
    await click('[data-test-oidc-key-create]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.create',
      'navigates to key create form'
    );
    await fillIn('[data-test-input="name"]', 'test-key');
    await click(SELECTORS.keySaveButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.details',
      'navigates to key details after save'
    );

    // create client with test-key
    await visit(OIDC_BASE_URL + '/clients');
    await click('[data-test-oidc-client-create]');
    await fillIn('[data-test-input="name"]', 'client-with-test-key');
    await click('[data-test-toggle-group="More options"]');
    await click('[data-test-component="search-select"] [data-test-icon="trash"]');
    await clickTrigger('#key');
    await selectChoose('[data-test-component="search-select"]#key', 'test-key');
    await click(SELECTORS.clientSaveButton);

    // edit key and limit applications
    await visit(OIDC_BASE_URL + '/keys');
    await click('[data-test-oidc-key-linked-block="test-key"] [data-test-popup-menu-trigger]');
    await click('[data-test-oidc-key-menu-link="edit"]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.edit',
      'key linked block popup menu navigates to edit'
    );
    await click('[data-test-oidc-radio="limited"]');
    await clickTrigger();
    assert.strictEqual(searchSelect.options.length, 1, 'dropdown has only application that uses this key');
    assert
      .dom('.ember-power-select-option')
      .hasTextContaining('client-with-test-key', 'dropdown renders correct application');
    await searchSelect.options.objectAt(0).click();
    await click(SELECTORS.keySaveButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Successfully updated the key test-key.',
      'renders success flash upon key updating'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.details',
      'navigates back to details on update'
    );
    await click(SELECTORS.keyClientsTab);
    assert
      .dom('[data-test-oidc-client-linked-block="client-with-test-key"]')
      .exists('lists client-with-test-key');
    assert.strictEqual(findAll('[data-test-oidc-client-linked-block]').length, 1, 'it lists only one client');

    // edit back to allow all
    await click(SELECTORS.keyDetailsTab);
    await click(SELECTORS.keyEditButton);
    await click('[data-test-oidc-radio="allow-all"]');
    await click(SELECTORS.keySaveButton);
    await click(SELECTORS.keyClientsTab);
    assert.notEqual(
      findAll('[data-test-oidc-client-linked-block]').length,
      1,
      'more than one client appears in key applications tab'
    );

    //* clean up test state
    await clearRecord(this.store, 'oidc/client', 'client-with-test-key');
    await clearRecord(this.store, 'oidc/client', 'client-with-default-key');
    await clearRecord(this.store, 'oidc/key', 'test-key');
  });

  test('it creates, rotates and deletes a key', async function (assert) {
    assert.expect(10);
    // mock client list so OIDC url does not redirect to landing page
    this.server.get('/identity/oidc/client', () => overrideResponse(null, { data: CLIENT_LIST_RESPONSE }));
    this.server.post('/identity/oidc/key/test-key/rotate', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      assert.strictEqual(json.verification_ttl, 86400, 'request made with correct args to accurate endpoint');
    });

    //* clear out test state
    await clearRecord(this.store, 'oidc/key', 'test-key');

    // create a new key
    await visit(OIDC_BASE_URL + '/keys/create');
    await fillIn('[data-test-input="name"]', 'test-key');
    // toggle ttls to false, testing it sets correct default duration
    await click('[data-test-input="rotationPeriod"]');
    await click('[data-test-input="verificationTtl"]');
    assert
      .dom('[data-test-oidc-radio="limited"] input')
      .isDisabled('limiting access radio button is disabled on create');
    assert
      .dom('[data-test-oidc-radio="limited"]')
      .hasClass('is-disabled', 'limited radio button label has disabled class');
    await click(SELECTORS.keySaveButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Successfully created the key test-key.',
      'renders success flash upon key creation'
    );

    // assert default values in details view are correct
    assert.dom('[data-test-value-div="Algorithm"]').hasText('RS256', 'defaults to RS526 algorithm');
    assert
      .dom('[data-test-value-div="Rotation period"]')
      .hasText('1 day', 'when toggled off rotation period defaults to 1 day');
    assert
      .dom('[data-test-value-div="Verification TTL"]')
      .hasText('1 day', 'when toggled off verification ttl defaults to 1 day');

    // rotate key
    await click(SELECTORS.keyDetailsTab);
    await click(SELECTORS.keyRotateButton);
    await click(SELECTORS.confirmActionButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Success: test-key connection was rotated.',
      'renders success flash upon key rotation'
    );
    // delete
    await click(SELECTORS.keyDeleteButton);
    await click(SELECTORS.confirmActionButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Key deleted successfully',
      'success flash message renders after deleting key'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.index',
      'navigates back to list view after delete'
    );
  });

  test('it renders client details and providers', async function (assert) {
    assert.expect(8);
    this.server.get('/identity/oidc/client', () => overrideResponse(null, { data: CLIENT_LIST_RESPONSE }));
    this.server.get('/identity/oidc/client/test-app', () =>
      overrideResponse(null, { data: CLIENT_DATA_RESPONSE })
    );
    await visit(OIDC_BASE_URL);
    await click('[data-test-oidc-client-linked-block]');
    assert.dom('[data-test-oidc-client-header]').hasText('test-app', 'renders application name as title');
    assert.dom(SELECTORS.clientDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.clientDeleteButton).exists('toolbar renders delete option');
    assert.dom(SELECTORS.clientEditButton).exists('toolbar renders edit button');
    assert.strictEqual(findAll('[data-test-component="info-table-row"]').length, 9, 'renders all info rows');

    await click(SELECTORS.clientProvidersTab);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.providers',
      'navigates to client providers route'
    );
    assert.dom(SELECTORS.clientProvidersTab).hasClass('active', 'providers tab is active');
    assert.dom('[data-test-oidc-provider-linked-block="default"]').exists('lists default provider');
  });

  test('it hides delete and edit client when no permission', async function (assert) {
    assert.expect(5);
    this.server.get('/identity/oidc/client', () => overrideResponse(null, { data: CLIENT_LIST_RESPONSE }));
    this.server.get('/identity/oidc/client/test-app', () =>
      overrideResponse(null, { data: CLIENT_DATA_RESPONSE })
    );
    this.server.post('/sys/capabilities-self', () =>
      capabilitiesStub(OIDC_BASE_URL + '/client/test-app', ['read'])
    );

    await visit(OIDC_BASE_URL);
    await click('[data-test-oidc-client-linked-block]');
    assert.dom('[data-test-oidc-client-header]').hasText('test-app', 'renders application name as title');
    assert.dom(SELECTORS.clientDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.clientDeleteButton).doesNotExist('delete option is hidden');
    assert.dom(SELECTORS.clientEditButton).doesNotExist('edit button is hidden');
    assert.strictEqual(findAll('[data-test-component="info-table-row"]').length, 9, 'renders all info rows');
  });

  test('it hides delete and edit key when no permission', async function (assert) {
    assert.expect(4);
    this.server.get('/identity/oidc/keys', () => overrideResponse(null, { data: { keys: ['test-key'] } }));
    this.server.get('/identity/oidc/key/test-key', () =>
      overrideResponse(null, {
        data: {
          algorithm: 'RS256',
          allowed_client_ids: ['*'],
          rotation_period: 86400,
          verification_ttl: 86400,
        },
      })
    );
    this.server.post('/sys/capabilities-self', () =>
      capabilitiesStub(OIDC_BASE_URL + '/key/test-key', ['read'])
    );

    await visit(OIDC_BASE_URL + '/keys');
    await click('[data-test-oidc-key-linked-block]');
    assert.dom(SELECTORS.keyDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.keyDeleteButton).doesNotExist('delete option is hidden');
    assert.dom(SELECTORS.keyEditButton).doesNotExist('edit button is hidden');
    assert.strictEqual(findAll('[data-test-component="info-table-row"]').length, 4, 'renders all info rows');
  });
});
