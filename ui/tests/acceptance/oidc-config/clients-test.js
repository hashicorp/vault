/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit, click, fillIn, findAll, currentRouteName, currentURL } from '@ember/test-helpers';
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
  ASSIGNMENT_LIST_RESPONSE,
  ASSIGNMENT_DATA_RESPONSE,
} from 'vault/tests/helpers/oidc-config';
import { capabilitiesStub, overrideResponse } from 'vault/tests/helpers/stubs';

const searchSelect = create(ss);
const flashMessage = create(fm);

// in congruency with backend verbiage 'applications' are referred to as 'clients'
// throughout the codebase and the term 'applications' only appears in the UI

module('Acceptance | oidc-config clients', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    oidcConfigHandlers(this.server);
    this.store = this.owner.lookup('service:store');
    return authPage.login();
  });

  module('keys', function () {
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
        .dom('[data-test-oidc-key-linked-block="default"] [data-test-item]')
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
      assert.strictEqual(
        currentURL(),
        '/vault/access/oidc/keys/default/clients',
        'navigates to key applications list '
      );

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
      await click('[data-test-oidc-key-linked-block="test-key"] [data-test-oidc-key-menu-link="edit"]');
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
      assert.strictEqual(
        findAll('[data-test-oidc-client-linked-block]').length,
        1,
        'it lists only one client'
      );

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
        assert.strictEqual(
          json.verification_ttl,
          86400,
          'request made with correct args to accurate endpoint'
        );
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
      assert.strictEqual(
        findAll('[data-test-component="info-table-row"]').length,
        9,
        'renders all info rows'
      );

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
      assert.strictEqual(
        findAll('[data-test-component="info-table-row"]').length,
        9,
        'renders all info rows'
      );
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
      assert.strictEqual(
        findAll('[data-test-component="info-table-row"]').length,
        4,
        'renders all info rows'
      );
    });
  });

  module('assignments', function () {
    test('it renders only allow_all when no assignments are configured', async function (assert) {
      assert.expect(3);

      //* clear out test state
      await clearRecord(this.store, 'oidc/assignment', 'test-assignment');

      await visit(OIDC_BASE_URL + '/assignments');
      assert.strictEqual(currentURL(), '/vault/access/oidc/assignments');
      assert.dom('[data-test-tab="assignments"]').hasClass('active', 'assignments tab is active');
      assert
        .dom('[data-test-oidc-assignment-linked-block="allow_all"]')
        .hasClass('is-disabled', 'renders default allow all assignment and is disabled.');
    });

    test('it renders empty state when no clients are configured', async function (assert) {
      assert.expect(5);
      this.server.get('/identity/oidc/client', () => overrideResponse(404));

      await visit(OIDC_BASE_URL);
      assert.strictEqual(currentURL(), '/vault/access/oidc');
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

    test('it creates an assignment inline, creates a client, updates client to limit access, deletes client', async function (assert) {
      assert.expect(21);

      //* clear out test state
      await clearRecord(this.store, 'oidc/client', 'test-app');
      await clearRecord(this.store, 'oidc/client', 'my-webapp'); // created by oidc-provider-test
      await clearRecord(this.store, 'oidc/assignment', 'assignment-inline');

      // create a client with allow all access
      await visit(OIDC_BASE_URL + '/clients/create');
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.clients.create',
        'navigates to create form'
      );
      await fillIn('[data-test-input="name"]', 'test-app');
      await click('[data-test-toggle-group="More options"]');
      // toggle ttls to false, testing it sets correct default duration
      await click('[data-test-input="idTokenTtl"]');
      await click('[data-test-input="accessTokenTtl"]');
      await click(SELECTORS.clientSaveButton);
      assert.strictEqual(
        flashMessage.latestMessage,
        'Successfully created the application test-app.',
        'renders success flash upon client creation'
      );
      assert.strictEqual(
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
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.clients.client.edit',
        'navigates to edit page from details'
      );
      await fillIn('[data-test-input="redirectUris"] [data-test-string-list-input="0"]', 'some-url.com');

      // limit access & create new assignment inline
      await click('[data-test-oidc-radio="limited"]');
      await clickTrigger();
      await fillIn('.ember-power-select-search input', 'assignment-inline');
      await searchSelect.options.objectAt(0).click();
      await click('[data-test-search-select="entities"] .ember-basic-dropdown-trigger');
      await searchSelect.options.objectAt(0).click();
      await click('[data-test-search-select="groups"] .ember-basic-dropdown-trigger');
      await searchSelect.options.objectAt(0).click();
      await click(SELECTORS.assignmentSaveButton);
      assert.strictEqual(
        flashMessage.latestMessage,
        'Successfully created the assignment assignment-inline.',
        'renders success flash upon assignment creating'
      );
      await click(SELECTORS.clientSaveButton);
      assert.strictEqual(
        flashMessage.latestMessage,
        'Successfully updated the application test-app.',
        'renders success flash upon client updating'
      );
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.clients.client.details',
        'navigates back to details on update'
      );
      assert.dom('[data-test-value-div="Redirect URI"]').hasText('some-url.com', 'shows updated attribute');
      assert
        .dom('[data-test-value-div="Assignment"]')
        .hasText('assignment-inline', 'updated to limited assignment');

      // edit back to allow_all
      await click(SELECTORS.clientEditButton);
      assert.dom(SELECTORS.clientSaveButton).hasText('Update', 'form button renders correct text');
      await click('[data-test-oidc-radio="allow-all"]');
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
      assert.strictEqual(
        flashMessage.latestMessage,
        'Application deleted successfully',
        'renders success flash upon deleting client'
      );
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.clients.index',
        'navigates back to list view after delete'
      );
      // delete last client
      await click('[data-test-oidc-client-linked-block]');
      assert.strictEqual(currentRouteName(), 'vault.cluster.access.oidc.clients.client.details');
      await click(SELECTORS.clientDeleteButton);
      await click(SELECTORS.confirmActionButton);

      //TODO this part of the test has a race condition
      //because other tests could have created clients - there is no guarantee that this will be the last
      //client in the list to redirect to the call to action
      //assert.strictEqual(
      //currentRouteName(),
      //'vault.cluster.access.oidc.index',
      //'redirects to call to action if only existing client is deleted'
      //);

      //* clean up test state
      await clearRecord(this.store, 'oidc/assignment', 'assignment-inline');
    });

    test('it creates, updates, and deletes an assignment', async function (assert) {
      assert.expect(14);
      await visit(OIDC_BASE_URL + '/assignments');

      //* ensure clean test state
      await clearRecord(this.store, 'oidc/assignment', 'test-assignment');

      // create a new assignment
      await click(SELECTORS.assignmentCreateButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.assignments.create',
        'navigates to create form'
      );
      assert.dom('[data-test-oidc-assignment-title]').hasText('Create Assignment', 'Form title renders');
      await fillIn('[data-test-input="name"]', 'test-assignment');
      await click('[data-test-component="search-select"]#entities .ember-basic-dropdown-trigger');
      await click('.ember-power-select-option');
      await click(SELECTORS.assignmentSaveButton);
      assert.strictEqual(
        flashMessage.latestMessage,
        'Successfully created the assignment test-assignment.',
        'renders success flash upon creating the assignment'
      );
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.assignments.assignment.details',
        'navigates to the assignments detail view after save'
      );

      // assert default values in assignment details view are correct
      assert.dom('[data-test-value-div="Name"]').hasText('test-assignment');
      assert.dom('[data-test-value-div="Entities"]').hasText('test-entity', 'shows the entity name.');

      // edit assignment
      await click(SELECTORS.assignmentEditButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.assignments.assignment.edit',
        'navigates to the assignment edit page from details'
      );
      assert.dom('[data-test-oidc-assignment-title]').hasText('Edit Assignment', 'Form title renders');
      await click('[data-test-component="search-select"]#groups .ember-basic-dropdown-trigger');
      await click('.ember-power-select-option');
      assert.dom('[data-test-oidc-assignment-save]').hasText('Update');
      await click(SELECTORS.assignmentSaveButton);
      assert.strictEqual(
        flashMessage.latestMessage,
        'Successfully updated the assignment test-assignment.',
        'renders success flash upon updating the assignment'
      );

      assert
        .dom('[data-test-value-div="Entities"]')
        .hasText('test-entity', 'it still shows the entity name.');
      assert.dom('[data-test-value-div="Groups"]').hasText('test-group', 'shows updated group name id.');

      // delete the assignment
      await click(SELECTORS.assignmentDeleteButton);
      await click(SELECTORS.confirmActionButton);
      assert.strictEqual(
        flashMessage.latestMessage,
        'Assignment deleted successfully',
        'renders success flash upon deleting assignment'
      );
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.assignments.index',
        'navigates back to assignment list view after delete'
      );
    });

    test('it navigates to and from an assignment from the list view', async function (assert) {
      assert.expect(6);
      this.server.get('/identity/oidc/assignment', () =>
        overrideResponse(200, { data: ASSIGNMENT_LIST_RESPONSE })
      );
      this.server.get('/identity/oidc/assignment/test-assignment', () =>
        overrideResponse(200, { data: ASSIGNMENT_DATA_RESPONSE })
      );
      await visit(OIDC_BASE_URL + '/assignments');
      assert
        .dom('[data-test-oidc-assignment-linked-block="test-assignment"]')
        .exists('displays linked block for test-assignment');

      await click(SELECTORS.assignmentCreateButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.assignments.create',
        'assignments index toolbar navigates to create form'
      );
      await click(SELECTORS.assignmentCancelButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.assignments.index',
        'create form navigates back to assignment index on cancel'
      );

      await click('[data-test-popup-menu-trigger]');
      await click('[data-test-oidc-assignment-menu-link="edit"]');
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.assignments.assignment.edit',
        'linked block popup menu navigates to edit'
      );
      await click(SELECTORS.assignmentCancelButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.assignments.assignment.details',
        'edit form navigates back to assignment details on cancel'
      );
      // navigate to details from index page
      await visit('/vault/access/oidc/assignments');
      await click('[data-test-popup-menu-trigger]');
      await click('[data-test-oidc-assignment-menu-link="details"]');
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.oidc.assignments.assignment.details',
        'popup menu navigates to assignment details'
      );
    });

    test('it hides assignment delete and edit when no permission', async function (assert) {
      assert.expect(5);
      this.server.get('/identity/oidc/assignment', () =>
        overrideResponse(null, { data: ASSIGNMENT_LIST_RESPONSE })
      );
      this.server.get('/identity/oidc/assignment/test-assignment', () =>
        overrideResponse(null, { data: ASSIGNMENT_DATA_RESPONSE })
      );
      this.server.post('/sys/capabilities-self', () =>
        capabilitiesStub(OIDC_BASE_URL + '/assignment/test-assignment', ['read'])
      );

      await visit(OIDC_BASE_URL + '/assignments');
      await click('[data-test-oidc-assignment-linked-block="test-assignment"]');
      assert
        .dom('[data-test-oidc-assignment-title]')
        .hasText('test-assignment', 'renders assignment name as title');
      assert.dom(SELECTORS.assignmentDetailsTab).hasClass('active', 'details tab is active');
      assert.dom(SELECTORS.assignmentDeleteButton).doesNotExist('delete option is hidden');
      assert.dom(SELECTORS.assignmentEditButton).doesNotExist('edit button is hidden');
      assert.strictEqual(
        findAll('[data-test-component="info-table-row"]').length,
        3,
        'renders all assignment info rows'
      );
    });
  });
});
