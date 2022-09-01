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
} from 'vault/tests/helpers/oidc-config';
const searchSelect = create(ss);
const flashMessage = create(fm);

// OIDC_BASE_URL = '/vault/access/oidc'

module('Acceptance | oidc-config/keys', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'oidcConfig';
  });

  hooks.beforeEach(async function () {
    this.store = await this.owner.lookup('service:store');
    // mock client list so OIDC url does not redirect to landing page
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
    return authPage.login();
  });

  hooks.afterEach(function () {
    return logout.visit();
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it lists default key and navigates to details', async function (assert) {
    assert.expect(7);
    await visit(OIDC_BASE_URL);
    await click('[data-test-tab="keys"]');
    assert.dom('[data-test-tab="keys"]').hasClass('active', 'keys tab is active');
    assert.equal(currentURL(), '/vault/access/oidc/keys');
    assert.dom('[data-test-oidc-key-linked-block]').hasText('default', 'index page lists default key');
    await click('[data-test-popup-menu-trigger]');

    await click('[data-test-oidc-key-menu-link="edit"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.edit',
      'key linked block popup menu navigates to edit'
    );
    await click(SELECTORS.keyCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.details',
      'key edit form navigates back to details on cancel'
    );

    // navigate to details from index page
    await click('[data-test-breadcrumb-link="oidc-keys"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.index',
      'keys breadcrumb navigates back to list view'
    );
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-key-menu-link="details"]');
    assert.dom(SELECTORS.keyDeleteButton).isDisabled('delete button is disabled for default key');
  });

  test('it creates, updates, rotates and deletes a key', async function (assert) {
    assert.expect(20);
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
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.details',
      'navigates to detail view after save'
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
    // edit and limit applications
    await click(SELECTORS.keyEditButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.keys.key.edit',
      'navigates to key edit page from details'
    );
    await click('label[for=limited]');
    await clickTrigger();
    assert.equal(searchSelect.options.length, 1, 'dropdown has only application that uses this key');
    assert
      .dom('.ember-power-select-option')
      .hasTextContaining('app-1', 'dropdown renders correct application');
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
    assert
      .dom('[data-test-oidc-client-linked-block]')
      .hasTextContaining('app-1', 'list of applications is just app-1');

    // edit back to allow all
    await click(SELECTORS.keyDetailsTab);
    await click(SELECTORS.keyEditButton);
    await click('label[for=allow-all]');
    await click(SELECTORS.keySaveButton);
    await click(SELECTORS.keyClientsTab);
    assert.equal(
      findAll('[data-test-oidc-client-linked-block]').length,
      2,
      'all applications appear in key applications tab'
    );
    // delete
    await click(SELECTORS.keyDetailsTab);
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

  test('it hides delete and edit when no permission', async function (assert) {
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
