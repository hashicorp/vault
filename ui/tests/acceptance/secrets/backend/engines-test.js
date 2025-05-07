/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, find, findAll, currentRouteName, visit, currentURL } from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { UNSUPPORTED_ENGINES, mountableEngines } from 'vault/helpers/mountable-secret-engines';

module('Acceptance | secret-engine list view', function (hooks) {
  setupApplicationTest(hooks);

  const createSecret = async (path, key, value, enginePath) => {
    await click(SES.createSecretLink);
    await fillIn('[data-test-secret-path]', path);

    await fillIn('[data-test-secret-key]', key);
    await fillIn(GENERAL.inputByAttr(key), value);
    await click('[data-test-secret-save]');
    await click(SES.crumb(enginePath));
  };

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  test('it allows you to disable an engine', async function (assert) {
    // first mount an engine so we can disable it.
    const enginePath = `alicloud-disable-${this.uid}`;
    await runCmd(mountEngineCmd('alicloud', enginePath));
    await visit('/vault/secrets');
    // to reduce flakiness, searching by engine name first in case there are pagination issues
    await selectChoose(GENERAL.searchSelect.trigger('filter-by-engine-name'), enginePath);
    assert.dom(SES.secretsBackendLink(enginePath)).exists('the alicloud engine is mounted');

    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('disable-engine'));
    await click(GENERAL.confirmButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends list page'
    );
    assert.dom(SES.secretsBackendLink(enginePath)).doesNotExist('does not show the disabled engine');
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
          .dom(SES.secretsBackendLink(enginePath))
          .doesNotHaveClass(
            'linked-block',
            `the linked-block class is not added to the unsupported ${engine}, which effectively disables it.`
          );
      } else {
        assert
          .dom(SES.secretsBackendLink(enginePath))
          .hasClass('linked-block', `linked-block class is added to supported ${engine} engines.`);
      }
      // cleanup
      await runCmd(deleteEngineCmd(enginePath));
    }
  });

  test('it filters by name and engine type', async function (assert) {
    const enginePath1 = `aws-1-${this.uid}`;
    const enginePath2 = `aws-2-${this.uid}`;

    await await runCmd(mountEngineCmd('aws', enginePath1));
    await await runCmd(mountEngineCmd('aws', enginePath2));
    await visit('/vault/secrets');
    // filter by type
    await clickTrigger('#filter-by-engine-type');
    await click(GENERAL.searchSelect.option());

    const rows = findAll(SES.secretsBackendLink());
    const rowsAws = Array.from(rows).filter((row) => row.innerText.includes('aws'));

    assert.strictEqual(rows.length, rowsAws.length, 'all rows returned are aws');
    // filter by name
    await clickTrigger('#filter-by-engine-name');
    const firstItemToSelect = find(GENERAL.searchSelect.option()).innerText;
    await click(GENERAL.searchSelect.option());
    const singleRow = document.querySelectorAll(SES.secretsBackendLink());
    assert.strictEqual(singleRow.length, 1, 'returns only one row');
    assert.dom(singleRow[0]).includesText(firstItemToSelect, 'shows the filtered by name engine');
    // clear filter by engine name
    await click(`#filter-by-engine-name ${GENERAL.searchSelect.removeSelected}`);
    const rowsAgain = document.querySelectorAll(SES.secretsBackendLink());
    assert.ok(rowsAgain.length > 1, 'filter has been removed');

    // cleanup
    await runCmd(deleteEngineCmd(enginePath1));
    await runCmd(deleteEngineCmd(enginePath2));
  });

  test('it applies overflow styling', async function (assert) {
    await visit('/vault/secrets');
    // not using the secret-engine-selector "secretPath" because I want to return the first node of a querySelectorAll
    const firstSecretEngine = document.querySelectorAll('[data-test-secret-path]')[0];
    assert.dom(firstSecretEngine).hasClass('overflow-wrap', 'secret engine name has overflow class ');
  });

  test('it allows navigation to a non-nested secret with pagination', async function (assert) {
    assert.expect(2);

    const enginePath1 = `kv-v1-${this.uid}`;
    const secretPath = 'secret-9';
    await runCmd(mountEngineCmd('kv', enginePath1));

    // check kv1
    await visit('/vault/secrets');
    await click(SES.secretsBackendLink(enginePath1));
    for (let i = 0; i <= 15; i++) {
      await createSecret(`secret-${i}`, 'foo', 'bar', enginePath1);
    }

    // navigate and check that details view is shown from non-nested secrets
    await click(GENERAL.pagination.next);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath1}/list?page=2`,
      'After clicking next page in navigates to the second page.'
    );
    await click(SES.secretLink(secretPath));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath1}/show/${secretPath}`,
      'After clicking a non-nested secret, it navigates to the details view.'
    );

    // cleanup
    await runCmd(deleteEngineCmd(enginePath1));
  });

  test('it allows navigation to a nested secret with pagination', async function (assert) {
    assert.expect(2);

    const enginePath1 = `kv-v1-${this.uid}`;
    const parentPath = 'nested';

    await runCmd(mountEngineCmd('kv', enginePath1));

    // check kv1
    await visit('/vault/secrets');
    await click(SES.secretsBackendLink(enginePath1));
    for (let i = 0; i <= 15; i++) {
      await createSecret(`${parentPath}/secret-${i}`, 'foo', 'bar', enginePath1);
    }

    // navigate and check that the children list view is shown from nested secrets
    await click(SES.secretLink(`${parentPath}/`));

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath1}/list/${parentPath}/`,
      'After clicking a nested secret it navigates to the children list view.'
    );

    await click(GENERAL.pagination.next);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath1}/list/${parentPath}/?page=2`,
      'After clicking next page it navigates to the second page.'
    );

    // cleanup
    await runCmd(deleteEngineCmd(enginePath1));
  });
});
