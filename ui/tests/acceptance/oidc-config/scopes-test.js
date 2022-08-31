import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn, findAll, currentRouteName } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { OIDC_BASE_URL, SELECTORS } from 'vault/tests/helpers/oidc-config';
import logout from 'vault/tests/pages/logout';
import { create } from 'ember-cli-page-object';
import fm from 'vault/tests/pages/components/flash-message';
const flashMessage = create(fm);
import {
  clearRecord,
  overrideMirageResponse,
  SCOPE_LIST_RESPONSE,
  SCOPE_DATA_RESPONSE,
  PROVIDER_LIST_RESPONSE,
  PROVIDER_DATA_RESPONSE,
} from 'vault/tests/helpers/oidc-config';
const SCOPES_URL = OIDC_BASE_URL.concat('/scopes');

// OIDC_BASE_URL = '/vault/access/oidc'

module('Acceptance | oidc-config/scopes', function (hooks) {
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

  test('it renders empty state when no scopes are configured', async function (assert) {
    assert.expect(3);
    this.server.get('/identity/oidc/scope', () => overrideMirageResponse(404));

    await visit(SCOPES_URL);
    assert.equal(currentURL(), '/vault/access/oidc/scopes');
    assert.dom('[data-test-tab="scopes"]').hasClass('active', 'scopes tab is active');
    // check empty state
    assert
      .dom(SELECTORS.scopeEmptyState)
      .hasText(
        `No scopes yet Use scope to define identity information about the authenticated user. Learn more. Create scope`,
        'renders empty state no scopes are configured'
      );
  });

  test('it creates a scope from empty state create scope button', async function (assert) {
    assert.expect(3);

    this.server.get('/identity/oidc/scope', () => overrideMirageResponse(404));
    await visit(SCOPES_URL);
    // create a new scope
    await click(SELECTORS.scopeCreateButtonEmptyState);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.create',
      'navigates to create form in empty state'
    );
    await fillIn('[data-test-input="name"]', 'test-scope');
    await click(SELECTORS.scopeSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully created the scope test-scope.',
      'renders success flash upon scope creation'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'navigates to scopes detail view after save'
    );
  });

  test('it creates, updates and deletes a scope', async function (assert) {
    assert.expect(11);

    //* clear out test state
    await clearRecord(this.store, 'oidc/scope', 'test-scope');

    // create a new scope
    await visit(SCOPES_URL + '/create');
    assert.equal(currentRouteName(), 'vault.cluster.access.oidc.scopes.create', 'navigates to create form');
    await fillIn('[data-test-input="name"]', 'test-scope');
    await fillIn('[data-test-input="description"]', 'this is a test');

    await click(SELECTORS.scopeSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully created the scope test-scope.',
      'renders success flash upon scope creation'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'navigates to detail view after save'
    );

    // assert values for new scope is correct
    assert.dom('[data-test-value-div="Name"]').hasText('test-scope', 'has correct created name');
    assert
      .dom('[data-test-value-div="Description"]')
      .hasText('this is a test', 'has correct created description');

    // edit scope
    await click(SELECTORS.scopeEditButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.edit',
      'navigates to edit page from details'
    );
    await fillIn('[data-test-input="description"]', 'this is an edit test');
    await click(SELECTORS.scopeSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully updated the scope test-scope.',
      'renders success flash upon scope updating'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'navigates back to details on update'
    );
    assert
      .dom('[data-test-value-div="Description"]')
      .hasText('this is an edit test', 'has correct edited description');

    // delete scope
    await click(SELECTORS.scopeDeleteButton);
    await click(SELECTORS.confirmActionButton);
    assert.equal(
      flashMessage.latestMessage,
      'Scope deleted successfully',
      'renders success flash upon deleting scope'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.index',
      'navigates back to list view after delete'
    );

    //* clear out test state
    await clearRecord(this.store, 'oidc/scope', 'test-scope');
  });

  test('it renders scope list when scopes exist', async function (assert) {
    this.server.get('/identity/oidc/scope', () => overrideMirageResponse(null, SCOPE_LIST_RESPONSE));
    this.server.get('/identity/oidc/scope/test-scope', () =>
      overrideMirageResponse(null, SCOPE_DATA_RESPONSE)
    );
    assert.expect(12);
    await visit(SCOPES_URL);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.index',
      'redirects to scopes index route when scopes exist'
    );
    assert.dom('[data-test-oidc-scope-linked-block]').hasText('test-scope', 'displays linked block for test');

    // navigates to/from create, edit, detail views from list view
    await click(SELECTORS.scopeCreateButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.create',
      'scope index toolbar navigates to create form'
    );
    await click(SELECTORS.scopeCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.index',
      'create form navigates back to index on cancel'
    );

    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-scope-menu-link="edit"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.edit',
      'linked block popup menu navigates to edit'
    );
    await click(SELECTORS.scopeCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'edit form navigates back to details on cancel'
    );

    // navigate to details from index page
    await click('[data-test-oidc-scope-return-to-index]');
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-scope-menu-link="details"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'popup menu navigates to details'
    );
    // check that details tab has all the information
    assert.dom(SELECTORS.scopeDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.scopeDeleteButton).exists('toolbar renders delete option');
    assert.dom(SELECTORS.scopeEditButton).exists('toolbar renders edit button');
    assert.equal(findAll('[data-test-component="info-table-row"]').length, 2, 'renders all info rows');
    // check JSON code renders
    assert.dom('[data-test-component="code-mirror-modifier"]').containsText('test', 'Code mirror renders');
  });

  test('it throws error when trying to delete when scope is currently being associated with any provider', async function (assert) {
    assert.expect(5);
    this.server.get('/identity/oidc/scope', () => overrideMirageResponse(null, SCOPE_LIST_RESPONSE));
    this.server.get('/identity/oidc/scope/test-scope', () =>
      overrideMirageResponse(null, SCOPE_DATA_RESPONSE)
    );
    this.server.get('/identity/oidc/provider', () => overrideMirageResponse(null, PROVIDER_LIST_RESPONSE));
    this.server.get('/identity/oidc/provider/test-provider', () => {
      overrideMirageResponse(null, PROVIDER_DATA_RESPONSE);
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
    await visit(SCOPES_URL);
    await click('[data-test-oidc-scope-linked-block]');
    assert.dom('[data-test-oidc-scope-header]').hasText('test-scope', 'renders scope name');
    assert.dom(SELECTORS.scopeDetailsTab).hasClass('active', 'details tab is active');

    // try to delete scope
    await click(SELECTORS.scopeDeleteButton);
    await click(SELECTORS.confirmActionButton);
    assert.equal(
      flashMessage.latestMessage,
      'unable to delete scope "test-scope" because it is currently referenced by these providers: test-provider',
      'renders error flash upon scope deletion'
    );
    assert.equal(findAll('[data-test-component="info-table-row"]').length, 2, 'renders all info rows');
    assert.dom('[data-test-component="code-mirror-modifier"]').containsText('test', 'Code mirror renders');
  });
});
