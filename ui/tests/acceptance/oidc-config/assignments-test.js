import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn, findAll, currentRouteName } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import { create } from 'ember-cli-page-object';
import fm from 'vault/tests/pages/components/flash-message';
import {
  OIDC_BASE_URL,
  SELECTORS,
  clearRecord,
  overrideCapabilities,
  overrideMirageResponse,
  ASSIGNMENT_LIST_RESPONSE,
  ASSIGNMENT_DATA_RESPONSE,
} from 'vault/tests/helpers/oidc-config';
const flashMessage = create(fm);
const ASSIGNMENTS_URL = OIDC_BASE_URL.concat('/assignments');

// OIDC_BASE_URL = '/vault/access/oidc'

module('Acceptance | oidc-config/assignments', function (hooks) {
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

  test('it renders only allow_all when no assignments are configured', async function (assert) {
    assert.expect(3);

    //* clear out test state
    await clearRecord(this.store, 'oidc/assignment', 'test-assignment');

    await visit(ASSIGNMENTS_URL);
    assert.equal(currentURL(), '/vault/access/oidc/assignments');
    assert
      .dom('[data-test-oidc-assignment-linked-block]')
      .includesText('allow_all', 'displays default assignment');
    assert.dom('[data-test-oidc-assignment-linked-block]').hasClass('is-disabled', 'Allow all is disabled.');
  });

  test('it creates, updates, and deletes an assignment', async function (assert) {
    assert.expect(12);
    await visit(ASSIGNMENTS_URL);
    // create a new assignment
    await click(SELECTORS.assignmentCreateButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.assignments.create',
      'navigates to create form'
    );
    await fillIn('[data-test-input="name"]', 'test-assignment');
    await click('[data-test-component="search-select"]#entities .ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
    await click(SELECTORS.assignmentSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully created the OIDC assignment test-assignment.',
      'renders success flash upon creating the assignment'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.assignments.assignment.details',
      'navigates to the assignments detail view after save'
    );
    // assert default values in assignment details view are correct
    assert.dom('[data-test-value-div="Name"]').hasText('test-assignment');
    assert.dom('[data-test-value-div="Entities"]').hasText('1234-12345', 'shows the entity id.');

    // edit assignment
    await click(SELECTORS.assignmentEditButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.assignments.assignment.edit',
      'navigates to the assignment edit page from details'
    );
    await click('[data-test-component="search-select"]#groups .ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
    assert.dom('[data-test-oidc-assignment-save]').hasText('Update');
    await click(SELECTORS.assignmentSaveButton);
    assert.equal(
      flashMessage.latestMessage,
      'Successfully updated the OIDC assignment test-assignment.',
      'renders success flash upon updating the assignment'
    );

    assert.dom('[data-test-value-div="Entities"]').hasText('1234-12345', 'it still shows the entity id.');
    assert.dom('[data-test-value-div="Groups"]').hasText('abcdef-123', 'shows updated group name id.');

    // delete the assignment
    await click(SELECTORS.assignmentDeleteButton);
    await click(SELECTORS.confirmDeleteButton);
    assert.equal(
      flashMessage.latestMessage,
      'Assignment deleted successfully',
      'renders success flash upon deleting assignment'
    );
    assert.equal(
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
    await visit(ASSIGNMENTS_URL);
    let linkedBlock = document.querySelectorAll('[data-test-oidc-assignment-linked-block]')[1];
    assert.dom(linkedBlock).hasText('test-assignment', 'displays linked block for assignment');

    await click(SELECTORS.assignmentCreateButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.assignments.create',
      'assignments index toolbar navigates to create form'
    );
    await click(SELECTORS.assignmentCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.assignments.index',
      'create form navigates back to assignment index on cancel'
    );

    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-assignment-menu-link="edit"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.assignments.assignment.edit',
      'linked block popup menu navigates to edit'
    );
    await click(SELECTORS.assignmentCancelButton);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.assignments.assignment.details',
      'edit form navigates back to assignment details on cancel'
    );
    // navigate to details from index page
    await visit('/vault/access/oidc/assignments');
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-oidc-assignment-menu-link="details"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.oidc.assignments.assignment.details',
      'popup menu navigates to details'
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

    await visit(ASSIGNMENTS_URL);
    let linkedBlock = document.querySelectorAll('[data-test-oidc-assignment-linked-block]')[1];
    await click(linkedBlock);
    assert
      .dom('[data-test-oidc-assignment-title]')
      .hasText('test-assignment', 'renders assignment name as title');
    assert.dom(SELECTORS.assignmentDetailsTab).hasClass('active', 'details tab is active');
    assert.dom(SELECTORS.assignmentDeleteButton).doesNotExist('delete option is hidden');
    assert.dom(SELECTORS.assignmentEditButton).doesNotExist('edit button is hidden');
    assert.equal(
      findAll('[data-test-component="info-table-row"]').length,
      3,
      'renders all assignment info rows'
    );
  });
});
