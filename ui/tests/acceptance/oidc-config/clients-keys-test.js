import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn, findAll, currentRouteName } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import { create } from 'ember-cli-page-object';
import { clickTrigger, selectChoose } from 'ember-power-select/test-support/helpers';
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

module('Acceptance | oidc-config clients and keys', function (hooks) {
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

  test('it creates a key, signs a client and edits key access to only that client', async function (assert) {
    //* clear out test state
    await clearRecord(this.store, 'oidc/client', 'app-with-test-key');
    await clearRecord(this.store, 'oidc/client', 'app-with-default-key');
    await clearRecord(this.store, 'oidc/key', 'test-key');

    // create client with default key
    await visit(OIDC_BASE_URL);
    await click(SELECTORS.oidcClientCreateButton);
    await fillIn('[data-test-input="name"]', 'app-with-default-key');
    await click(SELECTORS.clientSaveButton);

    // navigate to keys
    await visit(OIDC_BASE_URL + '/keys');
    assert.equal(currentURL(), '/vault/access/oidc/keys');
    assert
      .dom('[data-test-oidc-key-linked-block="default"]')
      .hasText('default', 'index page lists default key');

    // navigate to details from pop-up menu
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-key-menu-link="details"]');
    assert.dom(SELECTORS.keyDeleteButton).isDisabled('delete button is disabled for default key');
    await click(SELECTORS.keyEditButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.edit',
      'navigates to edit from key details'
    );
    await click(SELECTORS.keyCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.details',
      'key edit form navigates back to details on cancel'
    );
    await click(SELECTORS.keyClientsTab);
    assert
      .dom('[data-test-oidc-client-linked-block="app-with-default-key"]')
      .exists('lists correct app with default');

    // create a new key
    await click('[data-test-breadcrumb-link="oidc-keys"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.index',
      'keys breadcrumb navigates back to list view'
    );
    await click('[data-test-oidc-key-create]');
    assert.equal(currentRouteName(), 'vault.cluster.access.oidc.keys.create', 'navigates to key create form');
    await fillIn('[data-test-input="name"]', 'test-key');
    await click(SELECTORS.keySaveButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.details',
      'navigates to key details after save'
    );

    // create client with test-key
    await visit(OIDC_BASE_URL);
    await click('[data-test-oidc-client-create]');
    await fillIn('[data-test-input="name"]', 'app-with-test-key');
    await click('[data-test-toggle-group="More options"]');
    await click('[data-test-component="search-select"] [data-test-icon="trash"]');
    await clickTrigger('#key');
    await selectChoose('[data-test-component="search-select"]#key', 'test-key');
    await click(SELECTORS.clientSaveButton);

    // edit key and limit applications
    await click('[data-test-breadcrumb-link="oidc-clients"]');
    await click('[data-test-tab="keys"]');
    assert.dom('[data-test-tab="keys"]').hasClass('active', 'keys tab is active');
    await visit(OIDC_BASE_URL + '/keys');
    await click('[data-test-oidc-key-linked-block="test-key"] [data-test-popup-menu-trigger]');
    await click('[data-test-oidc-key-menu-link="edit"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.edit',
      'key linked block popup menu navigates to edit'
    );
    await click('label[for=limited]');
    await clickTrigger();
    assert.equal(searchSelect.options.length, 1, 'dropdown has only application that uses this key');
    assert
      .dom('.ember-power-select-option')
      .hasTextContaining('app-with-test-key', 'dropdown renders correct application');
    await searchSelect.options.objectAt(0).click();
    await click(SELECTORS.keySaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully updated the key test-key.',
      'renders success flash upon key updating'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.details',
      'navigates back to details on update'
    );
    await click(SELECTORS.keyClientsTab);
    assert.dom('[data-test-oidc-client-linked-block="app-with-test-key"]').exists('lists app-with-test-key');
    assert.equal(findAll('[data-test-oidc-client-linked-block]').length, 1, 'only lists one client');

    // edit back to allow all
    await click(SELECTORS.keyDetailsTab);
    await click(SELECTORS.keyEditButton);
    await click('label[for=allow-all]');
    await click(SELECTORS.keySaveButton);
    await click(SELECTORS.keyClientsTab);
    assert.equal(
      findAll('[data-test-oidc-client-linked-block]').length,
      2,
      'all clients appears in key applications tab'
    );

    //* clear out test state
    await clearRecord(this.store, 'oidc/client', 'app-with-test-key');
    await clearRecord(this.store, 'oidc/client', 'app-with-default-key');
    await clearRecord(this.store, 'oidc/key', 'test-key');
  });

  test('it renders client list when clients exist', async function (assert) {
    assert.expect(8);
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
    this.server.get('/identity/oidc/client/test-app', () =>
      overrideMirageResponse(null, CLIENT_DATA_RESPONSE)
    );
    await visit(OIDC_BASE_URL);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.index',
      'redirects to clients index route when clients exist'
    );
    assert.dom('[data-test-tab="clients"]').hasClass('active', 'clients tab is active');
    assert
      .dom('[data-test-oidc-client-linked-block]')
      .hasText('test-app Client ID: whaT7KB0C3iBH1l3rXhd5HPf0n6vXU0s', 'displays linked block for client');

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

  test('it creates, updates and deletes a client', async function (assert) {
    assert.expect(21);

    //* clear out test state
    await clearRecord(this.store, 'oidc/client', 'test-app');
    await clearRecord(this.store, 'oidc/assignment', 'assignment-1');

    await visit(OIDC_BASE_URL + '/clients/create');
    // create a new application
    assert.equal(currentRouteName(), 'vault.cluster.access.oidc.clients.create', 'navigates to create form');
    await fillIn('[data-test-input="name"]', 'test-app');
    await click('[data-test-toggle-group="More options"]');
    // toggle ttls to false, testing it sets correct default duration
    await click('[data-test-input="idTokenTtl"]');
    await click('[data-test-input="accessTokenTtl"]');
    await click(SELECTORS.clientSaveButton);

    assert.equal(
      flashMessage.latestMessage,
      'Successfully created the application test-app.',
      'renders success flash upon client creation'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'navigates to client details view after save'
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
      'Successfully updated the application test-app.',
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
    await fillIn('[data-test-input="name"]', 'app-to-delete');
    await click(SELECTORS.clientSaveButton);
    // immediately delete client, test transition
    await click(SELECTORS.clientDeleteButton);
    await click(SELECTORS.confirmActionButton);
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

    // delete last client
    await click('[data-test-oidc-client-linked-block]');
    await click(SELECTORS.clientDeleteButton);
    await click(SELECTORS.confirmActionButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.index',
      'redirects to call to action if only existing client is deleted'
    );
  });

  test('it creates, rotates and deletes a key', async function (assert) {
    assert.expect(12);
    // mock client list so OIDC url does not redirect to landing page
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
    this.server.post('/identity/oidc/key/test-key/rotate', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      assert.equal(json.verification_ttl, 86400, 'request made with correct args to accurate endpoint');
    });

    //* clear out test state
    await clearRecord(this.store, 'oidc/key', 'test-key');

    await visit(OIDC_BASE_URL + '/keys');
    // create a new key
    await click('[data-test-oidc-key-create]');
    assert.equal(currentRouteName(), 'vault.cluster.access.oidc.keys.create', 'navigates to key create form');
    await fillIn('[data-test-input="name"]', 'test-key');
    // toggle ttls to false, testing it sets correct default duration
    await click('[data-test-input="rotationPeriod"]');
    await click('[data-test-input="verificationTtl"]');
    assert.dom('input#limited').isDisabled('limiting access radio button is disabled on create');
    assert.dom('label[for=limited]').hasClass('is-disabled', 'limited radio button label has disabled class');
    await click(SELECTORS.keySaveButton);

    assert.equal(
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
    // check key's application list view
    await click(SELECTORS.keyClientsTab);
    assert.equal(
      findAll('[data-test-oidc-client-linked-block]').length,
      2,
      'all applications appear in key applications tab'
    );
    // rotate key
    await click(SELECTORS.keyDetailsTab);
    await click(SELECTORS.keyRotateButton);
    await click(SELECTORS.confirmActionButton);
    assert.equal(
      flashMessage.latestMessage,
      'Success: test-key connection was rotated.',
      'renders success flash upon key rotation'
    );
    // delete
    await click(SELECTORS.keyDeleteButton);
    await click(SELECTORS.confirmActionButton);
    assert.equal(
      flashMessage.latestMessage,
      'Key deleted successfully',
      'success flash message renders after deleting key'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.index',
      'navigates back to list view after delete'
    );
  });

  test('it renders client details and providers', async function (assert) {
    assert.expect(10);
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
    this.server.get('/identity/oidc/client/test-app', () =>
      overrideMirageResponse(null, CLIENT_DATA_RESPONSE)
    );
    await visit(OIDC_BASE_URL);
    await click('[data-test-oidc-client-linked-block]');
    assert.dom('[data-test-oidc-client-header]').hasText('test-app', 'renders application name as title');
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
    assert
      .dom('[data-test-oidc-provider-linked-block]')
      .hasTextContaining('default', 'shows default provider');
    assert.dom(SELECTORS.clientDeleteButton).doesNotExist('provider tab does not have delete option');
    assert.dom(SELECTORS.clientEditButton).doesNotExist('provider tab does not have edit button');
  });

  test('it hides delete and edit client when no permission', async function (assert) {
    assert.expect(5);
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
    this.server.get('/identity/oidc/client/test-app', () =>
      overrideMirageResponse(null, CLIENT_DATA_RESPONSE)
    );
    this.server.post('/sys/capabilities-self', () =>
      overrideCapabilities(OIDC_BASE_URL + '/client/test-app', ['read'])
    );

    await visit(OIDC_BASE_URL);
    await click('[data-test-oidc-client-linked-block]');
    assert.dom('[data-test-oidc-client-header]').hasText('test-app', 'renders application name as title');
    assert.dom(SELECTORS.clientDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.clientDeleteButton).doesNotExist('delete option is hidden');
    assert.dom(SELECTORS.clientEditButton).doesNotExist('edit button is hidden');
    assert.equal(findAll('[data-test-component="info-table-row"]').length, 9, 'renders all info rows');
  });

  test('it hides delete and edit key when no permission', async function (assert) {
    assert.expect(4);
    this.server.get('/identity/oidc/keys', () => overrideMirageResponse(null, { keys: ['test-key'] }));
    this.server.get('/identity/oidc/key/test-key', () =>
      overrideMirageResponse(null, {
        algorithm: 'RS256',
        allowed_client_ids: ['*'],
        rotation_period: 86400,
        verification_ttl: 86400,
      })
    );
    this.server.post('/sys/capabilities-self', () =>
      overrideCapabilities(OIDC_BASE_URL + '/key/test-key', ['read'])
    );

    await visit(OIDC_BASE_URL + '/keys');
    await click('[data-test-oidc-key-linked-block]');
    assert.dom(SELECTORS.keyDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.keyDeleteButton).doesNotExist('delete option is hidden');
    assert.dom(SELECTORS.keyEditButton).doesNotExist('edit button is hidden');
    assert.equal(findAll('[data-test-component="info-table-row"]').length, 4, 'renders all info rows');
  });
});
