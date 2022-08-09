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
    await visit(SCOPES_URL);
    assert.equal(currentURL(), '/vault/access/oidc/scopes');

    // check empty state
    assert.dom(SELECTORS.scopeEmptyState).hasText(
      `No scopes yet Use scope to define identity information about the authenticated user. Learn more. Create scope`,
      'renders empty state no scopes are configured'
    );
  });

  test('it creates a scope from empty state create scope button', async function (assert) {
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
    await visit(SCOPES_URL);
    // create a new scope
    await click(SELECTORS.scopeCreateButton);
    assert.equal(currentRouteName(), 'vault.cluster.access.oidc.scopes.create', 'navigates to create form');
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
});
