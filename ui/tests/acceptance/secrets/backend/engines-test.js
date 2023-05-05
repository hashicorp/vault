/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentRouteName, settled } from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import backendsPage from 'vault/tests/pages/secrets/backends';
import authPage from 'vault/tests/pages/auth';
import ss from 'vault/tests/pages/components/search-select';

const searchSelect = create(ss);

const disableEngine = async (enginePath) => {
  await backendsPage.visit();
  await settled();
  const row = backendsPage.rows.filterBy('path', `${enginePath}/`)[0];
  await row.menu();
  await settled();
  await backendsPage.disableButton();
  await settled();
  await backendsPage.confirmDisable();
  await settled();
};

module('Acceptance | secret-engine list view', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('it allows you to disable an engine', async function (assert) {
    // first mount an engine so we can disable it.
    const enginePath = `alicloud-disable-${this.uid}`;
    await mountSecrets.enable('alicloud', enginePath);
    await settled();
    assert.ok(backendsPage.rows.filterBy('path', `${enginePath}/`)[0], 'shows the mounted engine');
    await disableEngine(enginePath);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends page'
    );
    assert.strictEqual(
      backendsPage.rows.filterBy('path', `${enginePath}/`).length,
      0,
      'does not show the disabled engine'
    );
  });

  test('it adds disabled css styling to unsupported secret engines', async function (assert) {
    assert.expect(3);
    // first mount engine that is not supported
    const enginePath = `nomad-${this.uid}`;

    await mountSecrets.enable('nomad', enginePath);
    await settled();
    await backendsPage.visit();
    await settled();

    const rows = document.querySelectorAll('[data-test-auth-backend-link]');
    rows.forEach((node) => {
      if (node.innerText.includes('nomad')) {
        assert
          .dom(node)
          .doesNotHaveClass(
            'linked-block',
            `the linked-block class is not added to unsupported engines, which effectively disables it.`
          );
      } else {
        assert.dom(node).hasClass('linked-block', `linked-block class is added to supported engines.`);
      }
    });
    await disableEngine(enginePath);
  });

  test('it filters by name and engine type', async function (assert) {
    assert.expect(6);
    const enginePath1 = `aws-1-${this.uid}`;
    const enginePath2 = `aws-2-${this.uid}`;

    await mountSecrets.enable('aws', enginePath1);
    await mountSecrets.enable('aws', enginePath2);
    await backendsPage.visit();
    await settled();
    await clickTrigger('#filter-by-engine-type');
    await searchSelect.options.objectAt(0).click();

    assert
      .dom('[data-test-auth-backend-link]')
      .exists({ count: 2 }, 'Filter by aws engine type and 2 are returned');

    const rows = document.querySelectorAll('[data-test-auth-backend-link]');
    rows.forEach((block) => {
      assert.ok(block.innerText.includes('aws'), 'all the rows are engine type aws.');
    });

    await clickTrigger('#filter-by-engine-name');
    await searchSelect.options.objectAt(1).click();
    const row = document.querySelectorAll('[data-test-auth-backend-link]');

    assert.dom(row[0]).includesText('aws-2', 'shows the filtered by name engine');
    // clear filter by engine name
    await searchSelect.deleteButtons.objectAt(1).click();

    assert.dom('[data-test-auth-backend-link]').exists({ count: 2 }, 'Back to filtering by aws only');
    // clear filter by engine type
    await searchSelect.deleteButtons.objectAt(0).click();
    assert
      .dom('[data-test-auth-backend-link]')
      .exists({ count: 4 }, 'No filters shows two aws and cubbyhole and secrets.');

    await disableEngine(enginePath1);
    await disableEngine(enginePath2);
  });
});
