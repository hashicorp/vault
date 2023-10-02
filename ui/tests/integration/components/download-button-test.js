/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render, resetOnerror, setupOnerror } from '@ember/test-helpers';
import { isPresent } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

const SELECTORS = {
  button: '[data-test-download-button]',
  icon: '[data-test-icon="download"]',
};
module('Integration | Component | download button', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const downloadService = this.owner.lookup('service:download');
    this.downloadSpy = sinon.stub(downloadService, 'miscExtension');

    this.data = 'my data to download';
    this.filename = 'my special file';
    this.extension = 'csv';
  });

  test('it renders', async function (assert) {
    await render(hbs`
     <DownloadButton>
      <Icon @name="download" />
        Download
     </DownloadButton>
   `);
    assert.dom(SELECTORS.button).hasClass('button');
    assert.ok(isPresent(SELECTORS.icon), 'renders yielded icon');
    assert.dom(SELECTORS.button).hasTextContaining('Download', 'renders yielded text');
  });

  test('it downloads with defaults when only passed @data arg', async function (assert) {
    assert.expect(3);

    await render(hbs`
      <DownloadButton
        @data={{this.data}}
      >
        Download
      </DownloadButton>
    `);
    await click(SELECTORS.button);
    const [filename, content, extension] = this.downloadSpy.getCall(0).args;
    assert.ok(filename.includes('Z'), 'filename defaults to ISO string');
    assert.strictEqual(content, this.data, 'called with correct data');
    assert.strictEqual(extension, 'txt', 'called with default extension');
  });

  test('it calls download service with passed in args', async function (assert) {
    assert.expect(3);

    await render(hbs`
      <DownloadButton
        @data={{this.data}}
        @filename={{this.filename}}
        @mime={{this.mime}}
        @extension={{this.extension}}
      >
        Download
      </DownloadButton>
    `);

    await click(SELECTORS.button);
    const [filename, content, extension] = this.downloadSpy.getCall(0).args;
    assert.ok(filename.includes(`${this.filename}-`), 'filename added to ISO string');
    assert.strictEqual(content, this.data, 'called with correct data');
    assert.strictEqual(extension, this.extension, 'called with passed in extension');
  });

  test('it sets download content with arg passed to fetchData', async function (assert) {
    assert.expect(3);
    this.fetchData = () => 'this is fetched data from a parent function';
    await render(hbs`
      <DownloadButton @fetchData={{this.fetchData}} >
        Download
      </DownloadButton>
    `);

    await click(SELECTORS.button);
    const [filename, content, extension] = this.downloadSpy.getCall(0).args;
    assert.ok(filename.includes('Z'), 'filename defaults to ISO string');
    assert.strictEqual(content, this.fetchData(), 'called with fetched data');
    assert.strictEqual(extension, 'txt', 'called with default extension');
  });

  test('it throws error when both data and fetchData are passed as args', async function (assert) {
    assert.expect(1);
    setupOnerror((error) => {
      assert.strictEqual(
        error.message,
        'Assertion Failed: Only pass either @data or @fetchData, passing both means @data will be overwritten by the return value of @fetchData',
        'throws error with incorrect args'
      );
    });
    this.fetchData = () => 'this is fetched data from a parent function';
    await render(hbs`
        <DownloadButton @data={{this.data}} @fetchData={{this.fetchData}} />
      `);
    resetOnerror();
  });
});
