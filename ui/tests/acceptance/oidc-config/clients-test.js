import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn, findAll } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { BASE_URL, SELECTORS } from 'vault/tests/helpers/oidc-config';
import logout from 'vault/tests/pages/logout';

// in congruency with backend verbiage 'applications' are referred to as 'clients
// throughout the codebase ('applications' only appears in the UI)

// BASE_URL = '/vault/access/oidc'

module('Acceptance | oidc-config/clients', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'oidcConfig';
  });

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    const model = await this.store.peekRecord('oidc/client', 'some-app');
    if (model) model.destroyRecord();
    return authPage.login();
  });

  hooks.afterEach(async function () {
    const model = await this.store.peekRecord('oidc/client', 'some-app');
    if (model) model.destroyRecord();
    return logout.visit();
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it configures a new client that allows all assignments', async function (assert) {
    assert.expect(6);
    await visit(BASE_URL);
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
    assert.equal(currentURL(), '/vault/access/oidc/clients/create');
    await fillIn('[data-test-input="name"]', 'some-app');
    await click(SELECTORS.clientSaveButton);

    assert.equal(findAll('[data-test-component="info-table-row"]').length, 9);
  });
});
