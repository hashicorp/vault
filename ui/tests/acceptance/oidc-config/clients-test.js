import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn, findAll, currentRouteName, find } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { OIDC_BASE_URL, SELECTORS } from 'vault/tests/helpers/oidc-config';
import logout from 'vault/tests/pages/logout';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import fm from 'vault/tests/pages/components/flash-message';
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

  hooks.afterEach(async function () {
    this.store.findRecord('oidc/client', 'some-app').then((model) => {
      if (model) model.destroyRecord();
      this.store.findRecord('oidc/assignment', 'assignment-1').then((record) => {
        if (record) record.destroyRecord().then(() => {
          console.log(`destroyed ${record}!`)
        })
      });
      this.store.findRecord('oidc/assignment', 'assignment-2').then((record) => {
        if (record) record.destroyRecord().then( () => {
        console.log(`destroyed ${record}!`)
      });
    });
    return logout.visit();
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it saves a new client with defaults and edits to limit assignments', async function (assert) {
    assert.expect(28);
    await visit(OIDC_BASE_URL);
    if (find('[data-test-oidc-client-linked-block]')) {
      await click('[data-test-oidc-client-linked-block]');
      await click(SELECTORS.clientDeleteButton);
      await click(SELECTORS.confirmDeleteButton);
    }
    assert.equal(currentURL(), '/vault/access/oidc');
    // check empty state
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

    // create a new application
    await click(SELECTORS.oidcClientCreateButton);
    assert.equal(currentRouteName(), 'vault.cluster.access.oidc.clients.create', 'navigates to create form');
    assert
      .dom(SELECTORS.clientFormBreadcrumb + ' a')
      .hasText('Applications', 'create form breadcrumb links back to applications');
    await fillIn('[data-test-input="name"]', 'some-app');
    await click('[data-test-toggle-group="More options"]');
    // toggle ttls to false, testing it sets correct default duration
    await click('[data-test-input="idTokenTtl"]');
    await click('[data-test-input="accessTokenTtl"]');
    await click(SELECTORS.clientSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully created an application',
      'renders success flash upon client creation'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'navigates to detail view after save'
    );
    assert.dom('[data-test-oidc-client-header]').hasText('some-app', 'renders application name as title');
    assert.dom(SELECTORS.clientDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.clientDeleteButton).exists('toolbar renders delete option');
    assert.dom(SELECTORS.clientEditButton).exists('toolbar renders edit button');

    // assert default values in details view are as expected
    assert.equal(findAll('[data-test-component="info-table-row"]').length, 9, 'renders all info rows');
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

    // since we have a separate form component test, just test navigation
    await click(SELECTORS.clientDetailsTab);
    await click(SELECTORS.clientEditButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.edit',
      'navigates to edit page from details'
    );

    // edit client
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
    await click(SELECTORS.assignSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully created an assignment',
      'renders success flash upon assignment creating'
    );
    await click(SELECTORS.clientSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully updated application',
      'renders success flash upon client updating'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'navigates back to details on update'
    );
    assert.dom('[data-test-value-div="Redirect URI"]').hasText('some-url.com', 'shows updated attribute');
    assert.dom('[data-test-value-div="Assignment"]').hasText('assignment-1', 'updated to limited assignment');

    // reset state
    const assign1 = this.store.peekRecord('oidc/assignment', 'assignment-1');
    assign1.destroyRecord();
    await click(SELECTORS.clientDeleteButton);
    await click(SELECTORS.confirmDeleteButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully deleted client',
      'renders success flash upon deleting client'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients',
      'navigates back to list view after delete'
    );
  });

  test('it saves a new client with limited assignments and edits to allow_all', async function (assert) {
    assert.expect(6);
    await visit(OIDC_BASE_URL);
    if (find('[data-test-oidc-client-linked-block]')) {
      await click('[data-test-oidc-client-linked-block]');
      await click(SELECTORS.clientDeleteButton);
      await click(SELECTORS.confirmDeleteButton);
    }
    assert.equal(currentURL(), '/vault/access/oidc');

    // create a new client
    await click(SELECTORS.oidcClientCreateButton);
    assert.equal(currentURL(), '/vault/access/oidc/clients/create');
    await fillIn('[data-test-input="name"]', 'some-app');
    await click('label[for=limited]');

    // create a new assignment
    await clickTrigger();
    await fillIn('.ember-power-select-search input', 'assignment-1');
    await searchSelect.options.objectAt(0).click();
    await click('[data-test-search-select="entities"] .ember-basic-dropdown-trigger');
    await searchSelect.options.objectAt(0).click();
    await click('[data-test-search-select="groups"] .ember-basic-dropdown-trigger');
    await searchSelect.options.objectAt(0).click();
    await click(SELECTORS.assignSaveButton);
    // create second assignment
    await clickTrigger();
    await fillIn('.ember-power-select-search input', 'assignment-2');
    await searchSelect.options.objectAt(0).click();
    await click('[data-test-search-select="entities"] .ember-basic-dropdown-trigger');
    await searchSelect.options.objectAt(0).click();
    await click('[data-test-search-select="groups"] .ember-basic-dropdown-trigger');
    await searchSelect.options.objectAt(0).click();
    await click(SELECTORS.assignSaveButton);
    await click(SELECTORS.clientSaveButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'navigates back to details on save'
    );
    assert
      .dom('[data-test-row-value="Assignment"]')
      .hasText('assignment-1,assignment-2', 'lists both assignments');

    // edit client to allow all assignments
    await click(SELECTORS.clientEditButton);
    assert.dom(SELECTORS.clientSaveButton).hasText('Update', 'form button renders correct text');
    await click('label[for=allow-all]');
    await click(SELECTORS.clientSaveButton);
    assert
      .dom('[data-test-value-div="Assignment"]')
      .hasText('allow_all', 'client updated to allow all assignments');

    // reset state
    await click(SELECTORS.clientDeleteButton);
    await click(SELECTORS.confirmDeleteButton);
  });

  test('it renders client list when clients exist', async function (assert) {
    assert.expect(3);
    await visit(OIDC_BASE_URL);
    assert.equal(currentURL(), '/vault/access/oidc');

    // create a new client
    await click(SELECTORS.oidcClientCreateButton);
    await fillIn('[data-test-input="name"]', 'some-app');
    await click(SELECTORS.clientSaveButton);

    // navigate to list view
    await click(SELECTORS.clientHeaderBreadcrumb);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.index',
      'navigates to clients index route'
    );
    assert
      .dom('[data-test-oidc-client-linked-block]')
      .hasText('some-app', 'displays linked block for client');
    // navigate to create from index page
    await click('[data-test-oidc-client-create]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.create',
      'clients index toolbar navigates to create form'
    );
    assert
      .dom(SELECTORS.clientFormBreadcrumb + ' a')
      .hasText('Applications', 'create form breadcrumb has correct text');
    await click(SELECTORS.clientCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.index',
      'create form navigates back to index on cancel'
    );

    // navigate to edit from index page
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-client-menu-link="edit"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.edit',
      'linked block popup menu navigates to edit'
    );
    assert
      .dom(SELECTORS.clientFormBreadcrumb + ' a')
      .hasText('Details', 'edit form breadcrumb has correct text');
    await click(SELECTORS.clientCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'edit form navigates back to details on cancel'
    );
    assert
      .dom(SELECTORS.clientHeaderBreadcrumb + ' a')
      .hasText('Applications', 'details breadcrumb has correct text');
    await click(SELECTORS.clientHeaderBreadcrumb);

    // navigate to details from index page
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-client-menu-link="details"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.client.details',
      'popup menu navigates to details'
    );

    // reset state
    await click(SELECTORS.clientDeleteButton);
    await click(SELECTORS.confirmDeleteButton);
  });
});
