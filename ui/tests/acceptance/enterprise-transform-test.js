import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL, click } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { typeInSearch, selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';

import authPage from 'vault/tests/pages/auth';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import transformationsPage from 'vault/tests/pages/secrets/backend/transform/transformations';
import rolesPage from 'vault/tests/pages/secrets/backend/transform/roles';
import searchSelect from 'vault/tests/pages/components/search-select';

const searchSelectComponent = create(searchSelect);

const mount = async () => {
  let path = `transform-${Date.now()}`;
  await mountSecrets.enable('transform', path);
  return path;
};

const newTransformation = async (backend, name, submit = false) => {
  const transformationName = name || 'foo';
  await transformationsPage.visitCreate({ backend });
  await transformationsPage.name(transformationName);
  await clickTrigger('#template');
  await selectChoose('#template', '.ember-power-select-option', 0);
  // Don't automatically choose role because we might be testing that
  if (submit) {
    await transformationsPage.submit();
  }
  return transformationName;
};

const newRole = async (backend, name) => {
  const roleName = name || 'bar';
  await rolesPage.visitCreate({ backend });
  await rolesPage.name(roleName);
  await clickTrigger('#transformations');
  await selectChoose('#transformations', '.ember-power-select-option', 0);
  await rolesPage.submit();
  return roleName;
};

module('Acceptance | Enterprise | Transform secrets', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it enables Transform secrets engine and shows tabs', async function(assert) {
    let backend = `transform-${Date.now()}`;
    await mountSecrets.enable('transform', backend);
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/list`,
      'mounts and redirects to the transformations list page'
    );
    assert.ok(transformationsPage.isEmpty, 'renders empty state');
    assert
      .dom('.is-active[data-test-tab="Transformations"]')
      .exists('Has Transformations tab which is active');
    assert.dom('[data-test-tab="Roles"]').exists('Has Roles tab');
    assert.dom('[data-test-tab="Templates"]').exists('Has Templates tab');
    assert.dom('[data-test-tab="Alphabets"]').exists('Has Alphabets tab');
  });

  test('it can create a transformation and add itself to the role attached', async function(assert) {
    let backend = await mount();
    const transformationName = 'foo';
    const roleName = 'foo-role';
    await transformationsPage.createLink({ backend });
    assert.equal(currentURL(), `/vault/secrets/${backend}/create`, 'redirects to create transformation page');
    await transformationsPage.name(transformationName);

    assert.dom('[data-test-input="type"').hasValue('fpe', 'Has type FPE by default');
    assert.dom('[data-test-input="tweak_source"]').exists('Shows tweak source when FPE');
    await transformationsPage.type('masking');
    assert
      .dom('[data-test-input="masking_character"]')
      .exists('Shows masking character input when changed to masking type');
    assert.dom('[data-test-input="tweak_source"]').doesNotExist('Does not show tweak source when masking');
    await clickTrigger('#template');
    assert.equal(searchSelectComponent.options.length, 2, 'list shows two builtin options by default');
    await selectChoose('#template', '.ember-power-select-option', 0);
    await clickTrigger('#allowed_roles');
    await typeInSearch(roleName);
    await selectChoose('#allowed_roles', '.ember-power-select-option', 0);
    await transformationsPage.submit();
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/show/${transformationName}`,
      'redirects to show transformation page after submit'
    );
    await click(`[data-test-secret-breadcrumb="${backend}"]`);
    assert.equal(currentURL(), `/vault/secrets/${backend}/list`, 'Links back to list view from breadcrumb');
  });

  test('it can create a role and add itself to the transformation attached', async function(assert) {
    const roleName = 'my-role';
    let backend = await mount();
    // create transformation without role
    await newTransformation(backend, 'a-transformation', true);
    await click(`[data-test-secret-breadcrumb="${backend}"]`);
    assert.equal(currentURL(), `/vault/secrets/${backend}/list`, 'Links back to list view from breadcrumb');
    await click('[data-test-tab="Roles"]');
    assert.equal(currentURL(), `/vault/secrets/${backend}/list?tab=role`, 'links to role list page');
    // create role with transformation attached
    await rolesPage.createLink();
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/create?itemType=role`,
      'redirects to create role page'
    );
    await rolesPage.name(roleName);
    await clickTrigger('#transformations');
    assert.equal(searchSelectComponent.options.length, 1, 'lists the transformation');
    await selectChoose('#transformations', '.ember-power-select-option', 0);
    await rolesPage.submit();
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/show/role/${roleName}`,
      'redirects to show role page after submit'
    );
    await click(`[data-test-secret-breadcrumb="${backend}"]`);
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/list?tab=role`,
      'Links back to role list view from breadcrumb'
    );
  });

  test('it adds a role to a transformation when added to a role', async function(assert) {
    const roleName = 'role-test';
    let backend = await mount();
    let transformation = await newTransformation(backend, 'b-transformation', true);
    await newRole(backend, roleName);
    await transformationsPage.visitShow({ backend, id: transformation });
    assert.dom('[data-test-row-value="Allowed roles"]').hasText(roleName);
  });

  test('it shows a message if an update fails after save', async function(assert) {
    const roleName = 'role-remove';
    let backend = await mount();
    // Create transformation
    let transformation = await newTransformation(backend, 'c-transformation', true);
    // create role
    await newRole(backend, roleName);
    await transformationsPage.visitShow({ backend, id: transformation });
    assert.dom('[data-test-row-value="Allowed roles"]').hasText(roleName);
    // Edit transformation
    await click('[data-test-edit-link]');
    assert.dom('.modal.is-active').exists('Confirmation modal appears');
    await rolesPage.modalConfirm();
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/edit/${transformation}`,
      'Correctly links to edit page for secret'
    );
    // remove role
    await click('#allowed_roles [data-test-selected-list-button="delete"]');
    await transformationsPage.save();
    assert.dom('.flash-message.is-info').exists('Shows info message since role could not be updated');
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/show/${transformation}`,
      'Correctly links to show page for secret'
    );
    assert
      .dom('[data-test-row-value="Allowed roles"]')
      .doesNotExist('Allowed roles are no longer on the transformation');
  });
});
