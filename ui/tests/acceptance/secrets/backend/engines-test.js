/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, settled } from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import backendsPage from 'vault/tests/pages/secrets/backends';
import authPage from 'vault/tests/pages/auth';
import ss from 'vault/tests/pages/components/search-select';

const consoleComponent = create(consoleClass);
const searchSelect = create(ss);

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

    await backendsPage.visit();
    await settled();
    const row = backendsPage.rows.filterBy('path', `${enginePath}/`)[0];
    await row.menu();
    await settled();
    await backendsPage.disableButton();
    await settled();
    await backendsPage.confirmDisable();
    await settled();
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
    assert.expect(2);
    // first mount engine that is not supported
    const enginePath = `nomad-${this.uid}`;

    await mountSecrets.enable('nomad', enginePath);
    await settled();
    await backendsPage.visit();
    await settled();

    const rows = document.querySelectorAll('[data-test-secrets-backend-link]');
    const rowUnsupported = Array.from(rows).filter((row) => row.innerText.includes('nomad'));
    const rowSupported = Array.from(rows).filter((row) => row.innerText.includes('cubbyhole'));
    assert
      .dom(rowUnsupported[0])
      .doesNotHaveClass(
        'linked-block',
        `the linked-block class is not added to unsupported engines, which effectively disables it.`
      );
    assert.dom(rowSupported[0]).hasClass('linked-block', `linked-block class is added to supported engines.`);

    // cleanup
    await consoleComponent.runCommands([`delete sys/mounts/${enginePath}`]);
  });

  test('it filters by name and engine type', async function (assert) {
    assert.expect(4);
    const enginePath1 = `aws-1-${this.uid}`;
    const enginePath2 = `aws-2-${this.uid}`;

    await mountSecrets.enable('aws', enginePath1);
    await mountSecrets.enable('aws', enginePath2);
    await backendsPage.visit();
    await settled();
    // filter by type
    await clickTrigger('#filter-by-engine-type');
    await searchSelect.options.objectAt(0).click();

    const rows = document.querySelectorAll('[data-test-secrets-backend-link]');
    const rowsAws = Array.from(rows).filter((row) => row.innerText.includes('aws'));

    assert.strictEqual(rows.length, rowsAws.length, 'all rows returned are aws');
    // filter by name
    await clickTrigger('#filter-by-engine-name');
    const firstItemToSelect = searchSelect.options.objectAt(0).text;
    await searchSelect.options.objectAt(0).click();
    const singleRow = document.querySelectorAll('[data-test-secrets-backend-link]');
    assert.strictEqual(singleRow.length, 1, 'returns only one row');
    assert.dom(singleRow[0]).includesText(firstItemToSelect, 'shows the filtered by name engine');
    // clear filter by engine name
    await searchSelect.deleteButtons.objectAt(1).click();
    const rowsAgain = document.querySelectorAll('[data-test-secrets-backend-link]');
    assert.ok(rowsAgain.length > 1, 'filter has been removed');

    // cleanup
    await consoleComponent.runCommands([`delete sys/mounts/${enginePath1}`]);
    await consoleComponent.runCommands([`delete sys/mounts/${enginePath2}`]);
  });
});
