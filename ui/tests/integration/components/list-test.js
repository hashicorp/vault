/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, findAll, triggerEvent, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { v4 as uuidv4 } from 'uuid';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
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
      createSecretsEngine(undefined, 'aws', 'aws-1', 'v1.0.0'),
      createSecretsEngine(undefined, 'aws', 'aws-2', 'v2.0.0'),
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
    await fillIn(GENERAL.inputSearch('secret-engine-path'), 'external-test');

    const engineTooltip = document.querySelector(GENERAL.tooltip('Backend type'));
    await triggerEvent(engineTooltip, 'mouseenter');

    assert
      .dom(engineTooltip.nextSibling)
      .hasText(
        `This engine's type is not recognized by the UI. Please use the CLI to manage this engine.`,
        'shows tooltip text for unknown engine'
      );
  });

  test('hovering over the icon of an unsupported engine sets unsupported tooltip text', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);
    await fillIn(GENERAL.inputSearch('secret-engine-path'), 'nomad');

    const engineTooltip = document.querySelector(GENERAL.tooltip('Backend type'));
    await triggerEvent(engineTooltip, 'mouseenter');

    assert
      .dom(engineTooltip.nextSibling)
      .hasText(
        'The UI only supports configuration views for these secret engines. The CLI must be used to manage other engine resources.',
        'shows tooltip text for unsupported engine'
      );
  });

  test('hovering over the icon of a supported engine sets engine name as tooltip', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);
    await fillIn(GENERAL.inputSearch('secret-engine-path'), 'aws-1');

    const engineTooltip = document.querySelector(GENERAL.tooltip('Backend type'));
    await triggerEvent(engineTooltip, 'mouseenter');

    assert.dom(engineTooltip.nextSibling).hasText('AWS', 'shows tooltip text for supported engine with name');
  });

  test('hovering over the icon of a kv engine shows engine name and version', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}}/>`);
    await fillIn(GENERAL.inputSearch('secret-engine-path'), `kv-test`);

    const engineTooltip = document.querySelector(GENERAL.tooltip('Backend type'));
    await triggerEvent(engineTooltip, 'mouseenter');
    assert
      .dom(engineTooltip.nextSibling)
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

  test('it filters by engine path and engine type', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);
    // filter by type
    await click(GENERAL.toggleInput('filter-by-engine-type'));
    await click(GENERAL.checkboxByAttr('aws'));

    const rows = findAll(SES.secretsBackendLink());
    const rowsAws = Array.from(rows).filter((row) => row.innerText.includes('aws'));
    assert.strictEqual(rows.length, rowsAws.length, 'all rows returned are aws');

    // clear filter by type
    await click(GENERAL.button('Clear all'));
    assert.true(document.querySelectorAll(SES.secretsBackendLink()).length > 1, 'filter has been removed');

    // filter by path
    await fillIn(GENERAL.inputSearch('secret-engine-path'), 'kv');
    const singleRow = document.querySelectorAll(SES.secretsBackendLink());
    assert.dom(singleRow[0]).includesText('kv', 'shows the filtered by path engine');

    // clear filter by engine path
    await fillIn(GENERAL.inputSearch('secret-engine-path'), '');
    const rowsAgain = document.querySelectorAll(SES.secretsBackendLink());
    assert.true(rowsAgain.length > 1, 'search filter text has been removed');
  });

  test('it filters by engine version', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);
    // select engine type
    await click(GENERAL.toggleInput('filter-by-engine-type'));
    await click(GENERAL.checkboxByAttr('aws'));

    // filter by version
    await click(GENERAL.toggleInput('filter-by-engine-version'));
    await click(GENERAL.checkboxByAttr('v2.0.0'));
    const singleRow = document.querySelectorAll(SES.secretsBackendLink());
    assert.dom(singleRow[0]).includesText('aws-2', 'shows the single engine filtered by version');
  });

  test('it applies overflow styling', async function (assert) {
    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngineModels}} />`);
    // not using the secret-engine-selector "secretPath" because I want to return the first node of a querySelectorAll
    const firstSecretEngine = document.querySelectorAll('[data-test-secret-path]')[0];
    assert.dom(firstSecretEngine).hasClass('overflow-wrap', 'secret engine name has overflow class ');
  });
});
