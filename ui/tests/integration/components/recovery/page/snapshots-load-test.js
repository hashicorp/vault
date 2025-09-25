/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, triggerEvent, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import recoveryHandler from 'vault/mirage/handlers/recovery';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | recovery/snapshots-load', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    recoveryHandler(this.server);

    this.model = {
      configs: [],
      configError: undefined,
    };
    this.breadcrumbs = [
      { label: 'Secrets Recovery', route: 'vault.cluster.recovery.snapshots' },
      { label: 'Upload', route: 'vault.cluster.recovery.snapshots.load' },
    ];

    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';

    this.router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(this.router, 'transitionTo');
  });

  test('it should validate form fields', async function (assert) {
    await render(
      hbs`<Recovery::Page::Snapshots::Load @breadcrumbs={{this.breadcrumbs}} @model={{this.model}} />`
    );

    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('config'))
      .hasText('Please select a config', 'Config error renders.');

    assert.dom(GENERAL.validationErrorByAttr('url')).hasText('Please enter a url', 'Url error renders');

    await click(GENERAL.inputByAttr('manual'));
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('file'))
      .hasText('Please upload a snapshot file', 'File error renders.');
  });

  test('it loads a manual snapshot successfully', async function (assert) {
    assert.expect(3);
    this.server.post('sys/storage/raft/snapshot-load', (schema, req) => {
      assert.true(true, 'request is made to expected endpoint');
      const decoder = new TextDecoder();
      const bodyContent = decoder.decode(req.requestBody);
      assert.strictEqual(bodyContent, 'some content for a file', 'payload contains file data');
      return { data: {} };
    });
    const file = new Blob([['some content for a file']], { type: 'text/plain' });
    file.name = 'snapshot.snap';

    await render(
      hbs`<Recovery::Page::Snapshots::Load @breadcrumbs={{this.breadcrumbs}} @model={{this.model}} />`
    );

    await click(GENERAL.inputByAttr('manual'));
    await triggerEvent('[data-test-file-input]', 'change', { files: [file] });
    await click(GENERAL.submitButton);

    await waitUntil(() => this.transitionStub.called);

    assert.true(
      this.transitionStub.calledWith('vault.cluster.recovery.snapshots'),
      'Route transitions correctly on submit success'
    );
  });

  test('it loads an automated snapshot successfully', async function (assert) {
    assert.expect(3);
    this.server.post('sys/storage/raft/snapshot-auto/snapshot-load/:config_name', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.strictEqual(payload.url, 'test-snapshot-url', 'payload contains url');
      assert.strictEqual(req.params.config_name, 'test-config', 'request url is correct');
      return { data: {} };
    });

    this.model.configs = ['test-config'];

    await render(
      hbs`<Recovery::Page::Snapshots::Load @breadcrumbs={{this.breadcrumbs}} @model={{this.model}} />`
    );

    await click(GENERAL.selectByAttr('config'));
    await click('[data-option-index]');
    await fillIn(GENERAL.inputByAttr('url'), 'test-snapshot-url');
    await click(GENERAL.submitButton);

    await waitUntil(() => this.transitionStub.called);

    assert.true(
      this.transitionStub.calledWith('vault.cluster.recovery.snapshots'),
      'Route transitions correctly on submit success'
    );
  });
});
