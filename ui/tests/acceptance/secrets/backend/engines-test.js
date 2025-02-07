/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, find, findAll, currentRouteName, visit } from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { UNSUPPORTED_ENGINES, mountableEngines } from 'vault/helpers/mountable-secret-engines';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';

const SELECTORS = {
  backendLink: (path) =>
    path ? `[data-test-secrets-backend-link="${path}"]` : '[data-test-secrets-backend-link]',
};

module('Acceptance | secret-engine list view', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  test('it allows you to disable an engine', async function (assert) {
    // first mount an engine so we can disable it.
    const enginePath = `alicloud-disable-${this.uid}`;
    await runCmd(mountEngineCmd('alicloud', enginePath));
    await visit('/vault/secrets');
    assert.dom(SELECTORS.backendLink(enginePath)).exists();
    const row = SELECTORS.backendLink(enginePath);
    await click(`${row} ${GENERAL.menuTrigger}`);
    await click(`${row} ${GENERAL.confirmTrigger}`);
    await click(GENERAL.confirmButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends page'
    );
    assert.dom(SELECTORS.backendLink(enginePath)).doesNotExist('does not show the disabled engine');
  });

  test('it adds disabled css styling to unsupported secret engines', async function (assert) {
    assert.expect(16);
    const allEnginesArray = mountableEngines();
    for (const engineObject of allEnginesArray) {
      const engine = engineObject.type;
      const enginePath = `${engine}-${this.uid}`;
      await runCmd(mountEngineCmd(engine, enginePath));
      await visit('/vault/cluster/dashboard');
      await visit('/vault/secrets');

      if (UNSUPPORTED_ENGINES.includes(engine)) {
        assert
          .dom(PAGE.backends.link(enginePath))
          .doesNotHaveClass(
            'linked-block',
            `the linked-block class is not added to the unsupported ${engine}, which effectively disables it.`
          );
      } else {
        assert
          .dom(PAGE.backends.link(enginePath))
          .hasClass('linked-block', `linked-block class is added to supported ${engine} engines.`);
      }
      // cleanup
      await runCmd(deleteEngineCmd(enginePath));
    }
  });

  test('it filters by name and engine type', async function (assert) {
    assert.expect(5);
    const enginePath1 = `aws-1-${this.uid}`;
    const enginePath2 = `aws-2-${this.uid}`;

    await await runCmd(mountEngineCmd('aws', enginePath1));
    await await runCmd(mountEngineCmd('aws', enginePath2));
    await visit('/vault/secrets');
    // filter by type
    await clickTrigger('#filter-by-engine-type');
    await click(GENERAL.searchSelect.option());

    const rows = findAll(SELECTORS.backendLink());
    const rowsAws = Array.from(rows).filter((row) => row.innerText.includes('aws'));

    assert.strictEqual(rows.length, rowsAws.length, 'all rows returned are aws');
    // filter by name
    await clickTrigger('#filter-by-engine-name');
    const firstItemToSelect = find(GENERAL.searchSelect.option()).innerText;
    await click(GENERAL.searchSelect.option());
    const singleRow = document.querySelectorAll('[data-test-secrets-backend-link]');
    assert.strictEqual(singleRow.length, 1, 'returns only one row');
    assert.dom(singleRow[0]).includesText(firstItemToSelect, 'shows the filtered by name engine');
    // clear filter by engine name
    await click(`#filter-by-engine-name ${GENERAL.searchSelect.removeSelected}`);
    const rowsAgain = document.querySelectorAll('[data-test-secrets-backend-link]');
    assert.ok(rowsAgain.length > 1, 'filter has been removed');

    // verify overflow style exists on engine name
    assert.dom('[data-test-secret-path]').hasClass('overflow-wrap', 'secret engine name has overflow class ');

    // cleanup
    await runCmd(deleteEngineCmd(enginePath1));
    await runCmd(deleteEngineCmd(enginePath2));
  });
});
