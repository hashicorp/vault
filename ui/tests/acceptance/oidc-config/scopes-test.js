import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn, findAll, pauseTest, currentRouteName } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { OIDC_BASE_URL, SELECTORS } from 'vault/tests/helpers/oidc-config';
import logout from 'vault/tests/pages/logout';
import { create } from 'ember-cli-page-object';
import fm from 'vault/tests/pages/components/flash-message';
const flashMessage = create(fm);
const SCOPES_URL = OIDC_BASE_URL.concat('/scopes');

module('Acceptance | oidc-config/scopes', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'oidcConfig';
  });

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    const model = await this.store.peekRecord('oidc/scope', 'test');
    if (model) model.destroyRecord();
    return authPage.login();
  });

  hooks.afterEach(async function () {
    const model = await this.store.peekRecord('oidc/scope', 'test');
    if (model) model.destroyRecord();
    return logout.visit();
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it renders empty state when no scopes are configured', async function (assert) {
    assert.expect(2);
    await visit(SCOPES_URL);
    assert.equal(currentURL(), '/vault/access/oidc/scopes');

    // check empty state
    assert.dom(SELECTORS.scopeEmptyState).hasText(
      `No scopes yet Use scope to define identity information about the authenticated user. Learn more. Create scope`,
      'renders empty state no scopes are configured'
    );
  });

  test('it creates a scope from empty state create scope button', async function (assert) {
    assert.expect(3);
    await visit(SCOPES_URL);
    // create a new scope
    await click(SELECTORS.scopeCreateButtonEmptyState);
    assert.equal(currentRouteName(), 'vault.cluster.access.oidc.scopes.create', 'navigates to create form in empty state');
    await fillIn('[data-test-input="name"]', 'test');
    await click(SELECTORS.scopeSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully created a scope',
      'renders success flash upon scope creation'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'navigates to detail view after save'
    );
  });

  test('it creates, updates and deletes a scope', async function (assert) {
    assert.expect(11);
    await visit(SCOPES_URL);
    // create a new scope
    await click(SELECTORS.scopeCreateButton);
    assert.equal(currentRouteName(), 'vault.cluster.access.oidc.scopes.create', 'navigates to create form');
    await fillIn('[data-test-input="name"]', 'test');
    await fillIn('[data-test-input="description"]', 'this is a test');

    await click(SELECTORS.scopeSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully created a scope',
      'renders success flash upon scope creation'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.scopes.scope.details',
      'navigates to detail view after save'
    );
      
    // assert values for new scope is correct
    assert
      .dom('[data-test-value-div="Name"]')
      .hasText('test', 'has correct created name');
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
      'Successfully updated scope',
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
    await click(SELECTORS.confirmDeleteButton);
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
  });

  test('it renders scope list when scopes exist', async function (assert) {
    assert.expect(7);
    this.server.get('/identity/oidc/scope', () => {
      return {
        request_id: '8b89adf5-d086-5fe5-5876-59b0aaf5c0c3',
        lease_id: '',
        renewable: false,
        lease_duration: 0,
        data: {
          keys: ['test'],
        },
        wrap_info: null,
        warnings: null,
        auth: null,
      };
    });
    this.server.get('/identity/oidc/scope/test', () => {
      return {      
        request_id: 'test-id',
        data: {
          description:'this is a test',
          template:'{ test }',
        },
      };
    });
    await visit(OIDC_BASE_URL);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.clients.index',
      'redirects to clients index route when clients exist'
    );
    assert
      .dom('[data-test-oidc-client-linked-block]')
      .hasText('some-app', 'displays linked block for client');

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
});
