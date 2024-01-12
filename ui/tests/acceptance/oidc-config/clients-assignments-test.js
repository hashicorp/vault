/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn, findAll, currentRouteName } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import fm from 'vault/tests/pages/components/flash-message';
import {
  OIDC_BASE_URL, // -> '/vault/access/oidc'
  SELECTORS,
  clearRecord,
  overrideCapabilities,
  overrideMirageResponse,
  ASSIGNMENT_LIST_RESPONSE,
  ASSIGNMENT_DATA_RESPONSE,
} from 'vault/tests/helpers/oidc-config';
const searchSelect = create(ss);
const flashMessage = create(fm);

module('Acceptance | oidc-config clients and assignments', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'oidcConfig';
  });

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    return authPage.login();
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

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
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(404));

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

    assert.dom('[data-test-value-div="Entities"]').hasText('test-entity', 'it still shows the entity name.');
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
      overrideMirageResponse(null, ASSIGNMENT_LIST_RESPONSE)
    );
    this.server.get('/identity/oidc/assignment/test-assignment', () =>
      overrideMirageResponse(null, ASSIGNMENT_DATA_RESPONSE)
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
      overrideMirageResponse(null, ASSIGNMENT_LIST_RESPONSE)
    );
    this.server.get('/identity/oidc/assignment/test-assignment', () =>
      overrideMirageResponse(null, ASSIGNMENT_DATA_RESPONSE)
    );
    this.server.post('/sys/capabilities-self', () =>
      overrideCapabilities(OIDC_BASE_URL + '/assignment/test-assignment', ['read'])
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
