import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn, findAll, currentRouteName } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import fm from 'vault/tests/pages/components/flash-message';
import {
  OIDC_BASE_URL,
  SELECTORS,
  clearRecord,
  overrideCapabilities,
  overrideMirageResponse,
  CLIENT_LIST_RESPONSE,
  CLIENT_DATA_RESPONSE,
} from 'vault/tests/helpers/oidc-config';
const searchSelect = create(ss);
const flashMessage = create(fm);
// in congruency with backend verbiage 'applications' are referred to as 'clients
// throughout the codebase ('applications' only appears in the UI)

// OIDC_BASE_URL = '/vault/access/oidc'

module('Acceptance | oidc-config/clients', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'oidcConfig';
  });

  hooks.beforeEach(async function () {
    this.store = await this.owner.lookup('service:store');
    return authPage.login();
  });

  hooks.afterEach(function () {
    return logout.visit();
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it renders empty state when no clients are configured', async function (assert) {
    assert.expect(5);
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(404));

    await visit(OIDC_BASE_URL);
    assert.equal(currentURL(), '/vault/access/oidc');
    assert.dom('h1.title.is-3').hasText('OIDC Provider');
    assert.dom(SELECTORS.oidcHeader).hasText(
      `Configure Vault to act as an OIDC identity provider, and offer Vault’s various authentication
    methods and source of identity to any client applications. Learn more Create your first app`,
      'renders call to action header when no clients are configured'
    );
    assert.dom('[data-test-oidc-landing]').exists('landing page renders when no clients are configured');
    assert
      .dom(SELECTORS.oidcLandingImg)
      .hasAttribute('src', '/ui/images/oidc-landing.png', 'image renders image when no clients configured');
  });

  test('it creates, updates and deletes a client', async function (assert) {
    assert.expect(21);

    //* clear out test state
    await clearRecord(this.store, 'oidc/client', 'some-app');
    await clearRecord(this.store, 'oidc/client', 'test-app');
    await clearRecord(this.store, 'oidc/assignment', 'assignment-1');

    await visit(OIDC_BASE_URL + '/clients/create');
    // create a new application
    assert.equal(currentRouteName(), 'vault.cluster.access.oidc.clients.create', 'navigates to create form');
    await fillIn('[data-test-input="name"]', 'some-app');
    await click('[data-test-toggle-group="More options"]');
    // toggle ttls to false, testing it sets correct default duration
    await click('[data-test-input="idTokenTtl"]');
    await click('[data-test-input="accessTokenTtl"]');
    await click(SELECTORS.clientSaveButton);

    assert.equal(
      flashMessage.latestMessage,
      'Successfully created the application some-app.',
      'renders success flash upon client creation'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'navigates to detail view after save'
    );
    // assert default values in details view are correct
    assert.dom('[data-test-value-div="Assignment"]').hasText('allow_all', 'client allows all assignments');
    assert.dom('[data-test-value-div="Type"]').hasText('confidential', 'type defaults to confidential');
    assert
      .dom('[data-test-value-div="Key"] a')
      .hasText('default', 'client uses default key and renders a link');
    assert
      .dom('[data-test-value-div="Client ID"] [data-test-copy-button]')
      .exists('client ID exists and has copy button');
    assert
      .dom('[data-test-value-div="Client Secret"] [data-test-copy-button]')
      .exists('client secret exists and has copy button');
    assert
      .dom('[data-test-value-div="ID Token TTL"]')
      .hasText('1 day', 'ID token ttl toggled off sets default of 24h');
    assert
      .dom('[data-test-value-div="Access Token TTL"]')
      .hasText('1 day', 'access token ttl toggled off sets default of 24h');

    // edit client
    await click(SELECTORS.clientDetailsTab);
    await click(SELECTORS.clientEditButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.edit',
      'navigates to edit page from details'
    );
    await fillIn('[data-test-input="redirectUris"] [data-test-string-list-input="0"]', 'some-url.com');
    // limit & create new assignment
    await click('label[for=limited]');
    await clickTrigger();
    await fillIn('.ember-power-select-search input', 'assignment-1');
    await searchSelect.options.objectAt(0).click();
    await click('[data-test-search-select="entities"] .ember-basic-dropdown-trigger');
    await searchSelect.options.objectAt(0).click();
    await click('[data-test-search-select="groups"] .ember-basic-dropdown-trigger');
    await searchSelect.options.objectAt(0).click();
    await click(SELECTORS.assignmentSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully created the assignment assignment-1.',
      'renders success flash upon assignment creating'
    );
    await click(SELECTORS.clientSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully updated the application some-app.',
      'renders success flash upon client updating'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'navigates back to details on update'
    );
    assert.dom('[data-test-value-div="Redirect URI"]').hasText('some-url.com', 'shows updated attribute');
    assert.dom('[data-test-value-div="Assignment"]').hasText('assignment-1', 'updated to limited assignment');

    // edit back to allow_all
    await click(SELECTORS.clientEditButton);
    assert.dom(SELECTORS.clientSaveButton).hasText('Update', 'form button renders correct text');
    await click('label[for=allow-all]');
    await click(SELECTORS.clientSaveButton);
    assert
      .dom('[data-test-value-div="Assignment"]')
      .hasText('allow_all', 'client updated to allow all assignments');

    // create another client
    await visit(OIDC_BASE_URL + '/clients/create');
    await fillIn('[data-test-input="name"]', 'test-app');

    await click(SELECTORS.clientSaveButton);
    // immediately delete client, test transition
    await click(SELECTORS.clientDeleteButton);
    await click(SELECTORS.confirmDeleteButton);
    assert.equal(
      flashMessage.latestMessage,
      'Application deleted successfully',
      'renders success flash upon deleting client'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.index',
      'navigates back to list view after delete'
    );

    // delete some-app client (last client)
    await click('[data-test-oidc-client-linked-block]');
    await click(SELECTORS.clientDeleteButton);
    await click(SELECTORS.confirmDeleteButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.index',
      'redirects to call to action if only existing client is deleted'
    );
  });

  test('it renders client list when clients exist', async function (assert) {
    assert.expect(7);
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
    this.server.get('/identity/oidc/client/some-app', () =>
      overrideMirageResponse(null, CLIENT_DATA_RESPONSE)
    );
    await visit(OIDC_BASE_URL);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.index',
      'redirects to clients index route when clients exist'
    );
    assert
      .dom('[data-test-oidc-client-linked-block]')
      .hasText('some-app Client ID: whaT7KB0C3iBH1l3rXhd5HPf0n6vXU0s', 'displays linked block for client');

    // navigates to/from create, edit, detail views from list view
    await click('[data-test-oidc-client-create]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.create',
      'clients index toolbar navigates to create form'
    );
    await click(SELECTORS.clientCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.index',
      'create form navigates back to index on cancel'
    );

    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-client-menu-link="edit"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.edit',
      'linked block popup menu navigates to edit'
    );
    await click(SELECTORS.clientCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'edit form navigates back to details on cancel'
    );

    // navigate to details from index page
    await click('[data-test-link="oidc"]');
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-client-menu-link="details"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'popup menu navigates to details'
    );
  });

  test('it renders client details and providers', async function (assert) {
    assert.expect(10);
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
    this.server.get('/identity/oidc/client/some-app', () =>
      overrideMirageResponse(null, CLIENT_DATA_RESPONSE)
    );
    await visit(OIDC_BASE_URL);
    await click('[data-test-oidc-client-linked-block]');
    assert.dom('[data-test-oidc-client-header]').hasText('some-app', 'renders application name as title');
    assert.dom(SELECTORS.clientDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.clientDeleteButton).exists('toolbar renders delete option');
    assert.dom(SELECTORS.clientEditButton).exists('toolbar renders edit button');
    assert.equal(findAll('[data-test-component="info-table-row"]').length, 9, 'renders all info rows');

    await click(SELECTORS.clientProvidersTab);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.providers',
      'navigates to client providers route'
    );
    assert.dom(SELECTORS.clientProvidersTab).hasClass('active', 'providers tab is active');
    assert.dom('[data-test-oidc-provider="default"]').exists('default provider shows');
    assert.dom(SELECTORS.clientDeleteButton).doesNotExist('provider tab does not have delete option');
    assert.dom(SELECTORS.clientEditButton).doesNotExist('provider tab does not have edit button');
  });

  test('it hides delete and edit when no permission', async function (assert) {
    assert.expect(5);
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
    this.server.get('/identity/oidc/client/some-app', () =>
      overrideMirageResponse(null, CLIENT_DATA_RESPONSE)
    );
    this.server.post('/sys/capabilities-self', () =>
      overrideCapabilities(OIDC_BASE_URL + '/client/some-app', ['read'])
    );

    await visit(OIDC_BASE_URL);
    await click('[data-test-oidc-client-linked-block]');
    assert.dom('[data-test-oidc-client-header]').hasText('some-app', 'renders application name as title');
    assert.dom(SELECTORS.clientDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.clientDeleteButton).doesNotExist('delete option is hidden');
    assert.dom(SELECTORS.clientEditButton).doesNotExist('edit button is hidden');
    assert.equal(findAll('[data-test-component="info-table-row"]').length, 9, 'renders all info rows');
  });
});
