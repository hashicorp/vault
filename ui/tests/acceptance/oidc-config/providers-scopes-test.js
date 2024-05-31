/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn, findAll, currentRouteName } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import oidcConfigHandlers from 'vault/mirage/handlers/oidc-config';
import authPage from 'vault/tests/pages/auth';
import { create } from 'ember-cli-page-object';
import { selectChoose } from 'ember-power-select/test-support';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import fm from 'vault/tests/pages/components/flash-message';
import {
  OIDC_BASE_URL,
  SELECTORS,
  CLIENT_LIST_RESPONSE,
  SCOPE_LIST_RESPONSE,
  SCOPE_DATA_RESPONSE,
  PROVIDER_LIST_RESPONSE,
  PROVIDER_DATA_RESPONSE,
  clearRecord,
} from 'vault/tests/helpers/oidc-config';
import { capabilitiesStub, overrideResponse } from 'vault/tests/helpers/stubs';
const searchSelect = create(ss);
const flashMessage = create(fm);

// OIDC_BASE_URL = '/vault/access/oidc'

module('Acceptance |  oidc-config providers and scopes', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    oidcConfigHandlers(this.server);
    this.store = this.owner.lookup('service:store');
    // mock client list so OIDC BASE URL does not redirect to landing call-to-action image
    this.server.get('/identity/oidc/client', () => overrideResponse(null, { data: CLIENT_LIST_RESPONSE }));
    return authPage.login();
  });

  // LIST SCOPES EMPTY
  test('it navigates to scopes list view and renders empty state when no scopes are configured', async function (assert) {
    assert.expect(4);
    this.server.get('/identity/oidc/scope', () => overrideResponse(404));
    await visit(OIDC_BASE_URL);
    await click('[data-test-tab="scopes"]');
    assert.strictEqual(currentURL(), '/vault/access/oidc/scopes');
    assert.dom('[data-test-tab="scopes"]').hasClass('active', 'scopes tab is active');
    assert
      .dom(SELECTORS.scopeEmptyState)
      .hasText(
        `No scopes yet Use scope to define identity information about the authenticated user. OIDC provider scopes`,
        'renders empty state no scopes are configured'
      );
    assert
      .dom(SELECTORS.scopeCreateButton)
      .hasAttribute('href', '/ui/vault/access/oidc/scopes/create', 'toolbar renders create scope link');
  });

  // LIST SCOPE EXIST
  test('it renders scope list when scopes exist', async function (assert) {
    assert.expect(11);
    this.server.get('/identity/oidc/scope', () => overrideResponse(null, { data: SCOPE_LIST_RESPONSE }));
    this.server.get('/identity/oidc/scope/test-scope', () =>
      overrideResponse(null, { data: SCOPE_DATA_RESPONSE })
    );
    await visit(OIDC_BASE_URL + '/scopes');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.index',
      'redirects to scopes index route when scopes exist'
    );
    assert
      .dom('[data-test-oidc-scope-linked-block="test-scope"]')
      .exists('displays linked block for test scope');

    // navigates to/from create, edit, detail views from list view
    await click(SELECTORS.scopeCreateButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.create',
      'scope index toolbar navigates to create form'
    );
    await click(SELECTORS.scopeCancelButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.index',
      'create form navigates back to index on cancel'
    );

    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-scope-menu-link="edit"]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.edit',
      'linked block popup menu navigates to edit'
    );
    await click(SELECTORS.scopeCancelButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'scope edit form navigates back to details on cancel'
    );

    // navigate to details from index page
    await click('[data-test-breadcrumb-link="oidc-scopes"] a');
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-scope-menu-link="details"]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'popup menu navigates to details'
    );
    // check that details tab has all the information
    assert.dom(SELECTORS.scopeDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.scopeDeleteButton).exists('toolbar renders delete option');
    assert.dom(SELECTORS.scopeEditButton).exists('toolbar renders edit button');
    assert.strictEqual(findAll('[data-test-component="info-table-row"]').length, 2, 'renders all info rows');
  });

  // ERROR DELETING SCOPE
  test('it throws error when trying to delete when scope is currently being associated with any provider', async function (assert) {
    assert.expect(3);
    this.server.get('/identity/oidc/scope', () => overrideResponse(null, { data: SCOPE_LIST_RESPONSE }));
    this.server.get('/identity/oidc/scope/test-scope', () =>
      overrideResponse(null, { data: SCOPE_DATA_RESPONSE })
    );
    this.server.get('/identity/oidc/provider', () =>
      overrideResponse(null, { data: PROVIDER_LIST_RESPONSE })
    );
    this.server.get('/identity/oidc/provider/test-provider', () => {
      overrideResponse(null, { data: PROVIDER_DATA_RESPONSE });
    });
    // throw error when trying to delete test-scope since it is associated to test-provider
    this.server.delete(
      '/identity/oidc/scope/test-scope',
      () => ({
        errors: [
          'unable to delete scope "test-scope" because it is currently referenced by these providers: test-provider',
        ],
      }),
      400
    );
    await visit(OIDC_BASE_URL + '/scopes');
    await click('[data-test-oidc-scope-linked-block="test-scope"]');
    assert.dom('[data-test-oidc-scope-header]').hasText('test-scope', 'renders scope name');
    assert.dom(SELECTORS.scopeDetailsTab).hasClass('active', 'details tab is active');

    // try to delete scope
    await click(SELECTORS.scopeDeleteButton);
    await click(SELECTORS.confirmActionButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'unable to delete scope "test-scope" because it is currently referenced by these providers: test-provider',
      'renders error flash upon scope deletion'
    );
  });

  // CRUD SCOPE + CRUD PROVIDER
  test('it creates a scope, and creates a provider with that scope', async function (assert) {
    assert.expect(28);

    //* clear out test state
    await clearRecord(this.store, 'oidc/scope', 'test-scope');
    await clearRecord(this.store, 'oidc/provider', 'test-provider');

    // create a new scope
    await visit(OIDC_BASE_URL + '/scopes/create');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.create',
      'navigates to create form'
    );
    await fillIn('[data-test-input="name"]', 'test-scope');
    await fillIn('[data-test-input="description"]', 'this is a test');
    await fillIn('[data-test-component="code-mirror-modifier"] textarea', SCOPE_DATA_RESPONSE.template);
    await click(SELECTORS.scopeSaveButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Successfully created the scope test-scope.',
      'renders success flash upon scope creation'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'navigates to scope detail view after save'
    );
    assert.dom(SELECTORS.scopeDetailsTab).hasClass('active', 'scope details tab is active');
    assert.dom('[data-test-value-div="Name"]').hasText('test-scope', 'has correct created name');
    assert
      .dom('[data-test-value-div="Description"]')
      .hasText('this is a test', 'has correct created description');

    // edit scope
    await click(SELECTORS.scopeEditButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.edit',
      'navigates to edit page from details'
    );
    await fillIn('[data-test-input="description"]', 'this is an edit test');
    await click(SELECTORS.scopeSaveButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Successfully updated the scope test-scope.',
      'renders success flash upon scope updating'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'navigates back to scope details on update'
    );
    assert
      .dom('[data-test-value-div="Description"]')
      .hasText('this is an edit test', 'has correct edited description');

    // create a provider using test-scope
    await click('[data-test-breadcrumb-link="oidc-scopes"] a');
    await click('[data-test-tab="providers"]');
    assert.dom('[data-test-tab="providers"]').hasClass('active', 'providers tab is active');
    await click('[data-test-oidc-provider-create]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.providers.create',
      'navigates to provider create form'
    );
    await fillIn('[data-test-input="name"]', 'test-provider');
    await clickTrigger('#scopesSupported');
    await selectChoose('#scopesSupported', 'test-scope');
    await click(SELECTORS.providerSaveButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Successfully created the OIDC provider test-provider.',
      'renders success flash upon provider creation'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.providers.provider.details',
      'navigates to provider detail view after save'
    );

    // assert default values in details view are correct
    assert.dom('[data-test-value-div="Issuer URL"]').hasTextContaining('http://', 'issuer includes scheme');
    assert
      .dom('[data-test-value-div="Issuer URL"]')
      .hasTextContaining('identity/oidc/provider/test', 'issuer path populates correctly');
    assert
      .dom('[data-test-value-div="Scopes"] a')
      .hasAttribute('href', '/ui/vault/access/oidc/scopes/test-scope/details', 'lists scopes as links');

    // check provider's application list view
    await click(SELECTORS.providerClientsTab);
    assert.strictEqual(
      findAll('[data-test-oidc-client-linked-block]').length,
      2,
      'all applications appear in provider applications tab'
    );

    // edit and limit applications
    await click(SELECTORS.providerDetailsTab);
    await click(SELECTORS.providerEditButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.providers.provider.edit',
      'navigates to provider edit page from details'
    );
    await click('[data-test-oidc-radio="limited"]');
    await click('[data-test-component="search-select"]#allowedClientIds .ember-basic-dropdown-trigger');
    await fillIn('.ember-power-select-search input', 'test-app');
    await searchSelect.options.objectAt(0).click();
    await click(SELECTORS.providerSaveButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Successfully updated the OIDC provider test-provider.',
      'renders success flash upon provider updating'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.providers.provider.details',
      'navigates back to provider details after updating'
    );
    const providerModel = this.store.peekRecord('oidc/provider', 'test-provider');
    assert.propEqual(
      providerModel.allowedClientIds,
      ['whaT7KB0C3iBH1l3rXhd5HPf0n6vXU0s'],
      'provider saves client_id (not id or name) in allowed_client_ids param'
    );
    await click(SELECTORS.providerClientsTab);
    assert
      .dom('[data-test-oidc-client-linked-block]')
      .hasTextContaining('test-app', 'list of applications is just test-app');

    // edit back to allow all
    await click(SELECTORS.providerDetailsTab);
    await click(SELECTORS.providerEditButton);
    await click('[data-test-oidc-radio="allow-all"]');
    await click(SELECTORS.providerSaveButton);
    await click(SELECTORS.providerClientsTab);
    assert.strictEqual(
      findAll('[data-test-oidc-client-linked-block]').length,
      2,
      'all applications appear in provider applications tab'
    );

    // delete
    await click(SELECTORS.providerDetailsTab);
    await click(SELECTORS.providerDeleteButton);
    await click(SELECTORS.confirmActionButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Provider deleted successfully',
      'success flash message renders after deleting provider'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.providers.index',
      'navigates back to list view after delete'
    );

    // delete scope
    await visit(OIDC_BASE_URL + '/scopes/test-scope/details');
    await click(SELECTORS.scopeDeleteButton);
    await click(SELECTORS.confirmActionButton);
    assert.strictEqual(
      flashMessage.latestMessage,
      'Scope deleted successfully',
      'renders success flash upon deleting scope'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.index',
      'navigates back to list view after delete'
    );
  });

  // LIST PROVIDERS
  test('it lists default provider and navigates to details', async function (assert) {
    assert.expect(7);
    await visit(OIDC_BASE_URL);
    await click('[data-test-tab="providers"]');
    assert.dom('[data-test-tab="providers"]').hasClass('active', 'providers tab is active');
    assert.strictEqual(currentURL(), '/vault/access/oidc/providers');
    assert
      .dom('[data-test-oidc-provider-linked-block="default"]')
      .exists('index page lists default provider');
    await click('[data-test-popup-menu-trigger]');

    await click('[data-test-oidc-provider-menu-link="edit"]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.providers.provider.edit',
      'provider linked block popup menu navigates to edit'
    );
    await click(SELECTORS.providerCancelButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.providers.provider.details',
      'provider edit form navigates back to details on cancel'
    );

    // navigate to details from index page
    await click('[data-test-breadcrumb-link="oidc-providers"] a');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.oidc.providers.index',
      'providers breadcrumb navigates back to list view'
    );
    await click('[data-test-oidc-provider-linked-block="default"] [data-test-popup-menu-trigger]');
    await click('[data-test-oidc-provider-menu-link="details"]');
    assert.dom(SELECTORS.providerDeleteButton).doesNotExist('delete button hidden for default provider');
  });

  // PROVIDER DELETE + EDIT PERMISSIONS
  test('it hides delete and edit for a provider when no permission', async function (assert) {
    assert.expect(3);
    this.server.get('/identity/oidc/providers', () =>
      overrideResponse(null, { data: { providers: ['test-provider'] } })
    );
    this.server.get('/identity/oidc/provider/test-provider', () =>
      overrideResponse(null, {
        data: {
          allowed_client_ids: ['*'],
          issuer: 'http://127.0.0.1:8200/v1/identity/oidc/provider/test-provider',
          scopes_supported: ['test-scope'],
        },
      })
    );
    this.server.post('/sys/capabilities-self', () =>
      capabilitiesStub(OIDC_BASE_URL + '/provider/test-provider', ['read'])
    );

    await visit(OIDC_BASE_URL + '/providers');
    await click('[data-test-oidc-provider-linked-block]');
    assert.dom(SELECTORS.providerDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.providerDeleteButton).doesNotExist('delete option is hidden');
    assert.dom(SELECTORS.providerEditButton).doesNotExist('edit button is hidden');
  });
});
