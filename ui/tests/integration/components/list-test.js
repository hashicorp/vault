/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, find, findAll, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { v4 as uuidv4 } from 'uuid';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { createSecretsEngine } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | secret-engine/list', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
      },
    }));
    this.version = this.owner.lookup('service:version');
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashMessages.registerTypes(['success', 'danger']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    this.flashDangerSpy = sinon.spy(this.flashMessages, 'danger');
    this.uid = uuidv4();
    // generate a model of cubbyhole, kv, and nomad
    this.secretEngineModels = [
      createSecretsEngine(undefined, 'cubbyhole', 'cubbyhole-test'),
      createSecretsEngine(undefined, 'kv', 'kv-test'),
      createSecretsEngine(undefined, 'aws', 'aws-1'),
      createSecretsEngine(undefined, 'aws', 'aws-2'),
      createSecretsEngine(undefined, 'nomad', 'nomad-test'),
      createSecretsEngine(undefined, 'badType', 'external-test'),
    ];
  });

  test('it allows you to disable an engine', async function (assert) {
    const enginePath = 'kv-test';
    this.server.delete(`sys/mounts/${enginePath}`, () => {
      assert.true(true, 'Request is made to delete engine');
      return overrideResponse(204);
    });
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);

    assert.dom(SES.secretsBackendLink(enginePath)).exists('shows the link for the kvv2 secrets engine');
    const row = SES.secretsBackendLink(enginePath);
    await click(`${row} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('disable-engine'));
    await click(GENERAL.confirmButton);

    assert.true(
      this.flashSuccessSpy.calledWith(`The kv Secrets Engine at ${enginePath}/ has been disabled.`),
      'Flash message shows that engine was disabled.'
    );
  });

  test('hovering over the icon of an external unrecognized engine type sets unrecognized tooltip text', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);

    await selectChoose(GENERAL.searchSelect.trigger('filter-by-engine-name'), 'external-test');
    await triggerEvent('.hds-tooltip-button', 'mouseenter');

    assert
      .dom('.hds-tooltip-container')
      .hasText(
        'This plugin is not supported by the UI. Please use the CLI to manage this engine.',
        'shows tooltip text for unsupported engine'
      );
  });

  test('hovering over the icon of an unsupported engine sets unsupported tooltip text', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);

    await selectChoose(GENERAL.searchSelect.trigger('filter-by-engine-type'), 'nomad');
    await triggerEvent('.hds-tooltip-button', 'mouseenter');

    assert
      .dom('.hds-tooltip-container')
      .hasText(
        'The UI only supports configuration views for these secret engines. The CLI must be used to manage other engine resources.',
        'shows tooltip text for unsupported engine'
      );
  });

  test('hovering over the icon of a supported engine sets engine name as tooltip', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);
    await selectChoose(GENERAL.searchSelect.trigger('filter-by-engine-name'), 'aws-1');

    await triggerEvent('.hds-tooltip-button', 'mouseenter');

    assert.dom('.hds-tooltip-container').hasText('AWS', 'shows tooltip text for supported engine with name');
  });

  test('hovering over the icon of a kv engine shows engine name and version', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}}/>`);

    await selectChoose(GENERAL.searchSelect.trigger('filter-by-engine-name'), `kv-test`);

    await triggerEvent('.hds-tooltip-button', 'mouseenter');
    assert
      .dom('.hds-tooltip-container')
      .hasText('KV version 1', 'shows tooltip text for kv engine with version');
  });

  test('it adds disabled css styling to unsupported secret engines', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);
    assert
      .dom(SES.secretsBackendLink('nomad-test'))
      .doesNotHaveClass(
        'linked-block',
        `the linked-block class is not added to the unsupported nomad engine, which effectively disables it.`
      );

    assert
      .dom(SES.secretsBackendLink('aws-1'))
      .hasClass('linked-block', `linked-block class is added to supported aws engines.`);
  });

  test('it filters by name and engine type', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);
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
    assert.true(rowsAgain.length > 1, 'filter has been removed');
  });

  test('it applies overflow styling', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);
    // not using the secret-engine-selector "secretPath" because I want to return the first node of a querySelectorAll
    const firstSecretEngine = document.querySelectorAll('[data-test-secret-path]')[0];
    assert.dom(firstSecretEngine).hasClass('overflow-wrap', 'secret engine name has overflow class ');
  });
});
