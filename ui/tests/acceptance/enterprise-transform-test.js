/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL, click, settled, currentRouteName } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { typeInSearch, selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';

import authPage from 'vault/tests/pages/auth';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import transformationsPage from 'vault/tests/pages/secrets/backend/transform/transformations';
import rolesPage from 'vault/tests/pages/secrets/backend/transform/roles';
import templatesPage from 'vault/tests/pages/secrets/backend/transform/templates';
import alphabetsPage from 'vault/tests/pages/secrets/backend/transform/alphabets';
import searchSelect from 'vault/tests/pages/components/search-select';
import { runCmd } from '../helpers/commands';
import { allEngines } from 'vault/helpers/mountable-secret-engines';

const searchSelectComponent = create(searchSelect);

const mount = async () => {
  const path = `transform-${Date.now()}`;
  await mountSecrets.enable('transform', path);
  await settled();
  return path;
};

const newTransformation = async (backend, name, submit = false) => {
  const transformationName = name || 'foo';
  await transformationsPage.visitCreate({ backend });
  await settled();
  await transformationsPage.name(transformationName);
  await settled();
  await clickTrigger('#template');
  await selectChoose('#template', '.ember-power-select-option', 0);
  await settled();
  // Don't automatically choose role because we might be testing that
  if (submit) {
    await transformationsPage.submit();
    await settled();
  }
  return transformationName;
};

const newRole = async (backend, name) => {
  const roleName = name || 'bar';
  await rolesPage.visitCreate({ backend });
  await settled();
  await rolesPage.name(roleName);
  await settled();
  await clickTrigger('#transformations');
  await settled();
  await selectChoose('#transformations', '.ember-power-select-option', 0);
  await settled();
  await rolesPage.submit();
  await settled();
  return roleName;
};

module('Acceptance | Enterprise | Transform secrets', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it transitions to list route after mount success', async function (assert) {
    assert.expect(1);
    const engine = allEngines().find((e) => e.type === 'transform');

    // delete any previous mount with same name
    await runCmd([`delete sys/mounts/${engine.type}`]);
    await mountSecrets.visit();
    await mountSecrets.selectType(engine.type);
    await mountSecrets.next().path(engine.type);
    await mountSecrets.submit();

    assert.strictEqual(
      currentRouteName(),
      `vault.cluster.secrets.backend.list-root`,
      `${engine.type} navigates to list view`
    );

    // cleanup
    await runCmd([`delete sys/mounts/${engine.type}`]);
  });

  test('it enables Transform secrets engine and shows tabs', async function (assert) {
    const backend = `transform-${Date.now()}`;
    await mountSecrets.enable('transform', backend);
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/list`,
      'mounts and redirects to the transformations list page'
    );
    assert.ok(transformationsPage.isEmpty, 'renders empty state');
    assert
      .dom('.active[data-test-secret-list-tab="Transformations"]')
      .exists('Has Transformations tab which is active');
    assert.dom('[data-test-secret-list-tab="Roles"]').exists('Has Roles tab');
    assert.dom('[data-test-secret-list-tab="Templates"]').exists('Has Templates tab');
    assert.dom('[data-test-secret-list-tab="Alphabets"]').exists('Has Alphabets tab');
  });

  test('it can create a transformation and add itself to the role attached', async function (assert) {
    const backend = await mount();
    const transformationName = 'foo';
    const roleName = 'foo-role';
    await settled();
    await transformationsPage.createLink({ backend });
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/create`,
      'redirects to create transformation page'
    );
    await transformationsPage.name(transformationName);
    await settled();
    assert.dom('[data-test-input="type"]').hasValue('fpe', 'Has type FPE by default');
    assert.dom('[data-test-input="tweak_source"]').exists('Shows tweak source when FPE');
    await transformationsPage.type('masking');
    await settled();
    assert
      .dom('[data-test-input="masking_character"]')
      .exists('Shows masking character input when changed to masking type');
    assert.dom('[data-test-input="tweak_source"]').doesNotExist('Does not show tweak source when masking');
    await clickTrigger('#template');
    await settled();
    assert.strictEqual(searchSelectComponent.options.length, 2, 'list shows two builtin options by default');
    await selectChoose('#template', '.ember-power-select-option', 0);
    await settled();

    await clickTrigger('#allowed_roles');
    await settled();
    await typeInSearch(roleName);
    await settled();
    await selectChoose('#allowed_roles', '.ember-power-select-option', 0);
    await settled();
    await transformationsPage.submit();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/show/${transformationName}`,
      'redirects to show transformation page after submit'
    );
    await click(`[data-test-secret-breadcrumb="${backend}"]`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/list`,
      'Links back to list view from breadcrumb'
    );
  });

  test('it can create a role and add itself to the transformation attached', async function (assert) {
    const roleName = 'my-role';
    const backend = await mount();
    // create transformation without role
    await newTransformation(backend, 'a-transformation', true);
    await click(`[data-test-secret-breadcrumb="${backend}"]`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/list`,
      'Links back to list view from breadcrumb'
    );
    await click('[data-test-secret-list-tab="Roles"]');
    assert.strictEqual(currentURL(), `/vault/secrets/${backend}/list?tab=role`, 'links to role list page');
    // create role with transformation attached
    await rolesPage.createLink();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/create?itemType=role`,
      'redirects to create role page'
    );
    await rolesPage.name(roleName);
    await clickTrigger('#transformations');
    assert.strictEqual(searchSelectComponent.options.length, 1, 'lists the transformation');
    await selectChoose('#transformations', '.ember-power-select-option', 0);
    await rolesPage.submit();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/show/role/${roleName}`,
      'redirects to show role page after submit'
    );
    await click(`[data-test-secret-breadcrumb="${backend}"]`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/list?tab=role`,
      'Links back to role list view from breadcrumb'
    );
  });

  test('it adds a role to a transformation when added to a role', async function (assert) {
    const roleName = 'role-test';
    const backend = await mount();
    const transformation = await newTransformation(backend, 'b-transformation', true);
    await newRole(backend, roleName);
    await transformationsPage.visitShow({ backend, id: transformation });
    await settled();
    assert.dom('[data-test-row-value="Allowed roles"]').hasText(roleName);
  });

  test('it shows a message if an update fails after save', async function (assert) {
    const roleName = 'role-remove';
    const backend = await mount();
    // Create transformation
    const transformation = await newTransformation(backend, 'c-transformation', true);
    // create role
    await newRole(backend, roleName);
    await settled();
    await transformationsPage.visitShow({ backend, id: transformation });
    assert.dom('[data-test-row-value="Allowed roles"]').hasText(roleName);
    // Edit transformation
    await click('[data-test-edit-link]');
    assert.dom('.modal.is-active').exists('Confirmation modal appears');
    await rolesPage.modalConfirm();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/edit/${transformation}`,
      'Correctly links to edit page for secret'
    );
    // remove role
    await settled();
    await click('#allowed_roles [data-test-selected-list-button="delete"]');

    await transformationsPage.save();
    await settled();
    assert.dom('.flash-message.is-info').exists('Shows info message since role could not be updated');
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/show/${transformation}`,
      'Correctly links to show page for secret'
    );
    assert
      .dom('[data-test-row-value="Allowed roles"]')
      .doesNotExist('Allowed roles are no longer on the transformation');
  });

  test('it allows creation and edit of a template', async function (assert) {
    const templateName = 'my-template';
    const backend = await mount();
    await click('[data-test-secret-list-tab="Templates"]');

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/list?tab=template`,
      'links to template list page'
    );
    await settled();
    await templatesPage.createLink();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/create?itemType=template`,
      'redirects to create template page'
    );
    await templatesPage.name(templateName);
    await templatesPage.pattern(`(\\d{4})`);
    await clickTrigger('#alphabet');
    await settled();
    assert.ok(searchSelectComponent.options.length > 0, 'lists built-in alphabets');
    await selectChoose('#alphabet', '.ember-power-select-option', 0);
    assert.dom('#alphabet .ember-power-select-trigger').doesNotExist('Alphabet input no longer searchable');
    await templatesPage.submit();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/show/template/${templateName}`,
      'redirects to show template page after submit'
    );
    await templatesPage.editLink();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/edit/template/${templateName}`,
      'Links to template edit page'
    );
    await settled();
    assert.dom('[data-test-input="name"]').hasAttribute('readonly');
  });

  test('it allows creation and edit of an alphabet', async function (assert) {
    const alphabetName = 'vowels-only';
    const backend = await mount();
    await click('[data-test-secret-list-tab="Alphabets"]');

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/list?tab=alphabet`,
      'links to alphabet list page'
    );
    await alphabetsPage.createLink();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/create?itemType=alphabet`,
      'redirects to create alphabet page'
    );
    await alphabetsPage.name(alphabetName);
    await alphabetsPage.alphabet('aeiou');
    await alphabetsPage.submit();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/show/alphabet/${alphabetName}`,
      'redirects to show alphabet page after submit'
    );
    assert.dom('[data-test-row-value="Name"]').hasText(alphabetName);
    assert.dom('[data-test-row-value="Alphabet"]').hasText('aeiou');
    await alphabetsPage.editLink();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/edit/alphabet/${alphabetName}`,
      'Links to alphabet edit page'
    );
    assert.dom('[data-test-input="name"]').hasAttribute('readonly');
  });
});
