/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { render, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
/* eslint-disable no-useless-escape */

module('Integration | Component | kv-v2 | KvPathsCard', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.isCondensed = false;
    this.clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();

    this.assertClipboard = async (assert, expected, index) => {
      const buttons = findAll(GENERAL.copyButton);
      await click(buttons[index]);
      assert.true(this.clipboardSpy.called, 'clipboard.writeText was called');

      const call = this.clipboardSpy.getCall(index);
      assert.strictEqual(
        call.args[0],
        expected,
        `Expected clipboard text to match: ${expected} (at index ${index})`
      );
    };

    this.renderComponent = async () => {
      return render(
        hbs`<KvPathsCard @backend={{this.backend}} @path={{this.path}} @isCondensed={{this.isCondensed}} />`,
        {
          owner: this.engine,
        }
      );
    };
  });

  hooks.afterEach(function () {
    this.clipboardSpy.restore(); // Restore original clipboard
  });

  test('it renders condensed version', async function (assert) {
    this.isCondensed = true;

    await this.renderComponent();

    assert.dom('[data-test-component="info-table-row"] .helper-text').doesNotExist('subtext does not render');

    assert.dom('[data-test-component="info-table-row"] div').hasClass('is-one-quarter');
    assert.dom(PAGE.paths.codeSnippet('cli')).doesNotExist();
    assert.dom(PAGE.paths.codeSnippet('api')).doesNotExist();

    const paths = [
      { label: 'API path', expected: `/v1/${this.backend}/data/${this.path}` },
      { label: 'Configuration path', expected: `${this.backend}/data/${this.path}` },
      { label: 'CLI path', expected: `-mount="${this.backend}" "${this.path}"` },
    ];
    for (const [index, path] of paths.entries()) {
      assert.dom(PAGE.infoRowValue(path.label)).hasText(path.expected);
      await this.assertClipboard(assert, path.expected, index);
    }
  });

  test('it renders uncondensed version', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-component="info-table-row"] .helper-text').exists('subtext renders');
    assert.dom('[data-test-component="info-table-row"] div').hasClass('is-one-third');
    assert.dom(PAGE.infoRowValue('API path for metadata')).exists();
    assert.dom(PAGE.paths.codeSnippet('cli')).exists();
    assert.dom(PAGE.paths.codeSnippet('api')).exists();
  });

  test('it renders copyable paths', async function (assert) {
    const paths = [
      { label: 'API path', expected: `/v1/${this.backend}/data/${this.path}` },
      { label: 'Configuration path', expected: `${this.backend}/data/${this.path}` },
      { label: 'CLI path', expected: `-mount="${this.backend}" "${this.path}"` },
      { label: 'API path for metadata', expected: `/v1/${this.backend}/metadata/${this.path}` },
      { label: 'Configuration metadata path', expected: `${this.backend}/metadata/${this.path}` },
    ];

    await this.renderComponent();

    for (const [index, path] of paths.entries()) {
      assert.dom(PAGE.infoRowValue(path.label)).hasText(path.expected);
      await this.assertClipboard(assert, path.expected, index);
    }
  });

  test('it renders copyable encoded mount and secret paths', async function (assert) {
    this.path = `my spacey!"secret`;
    this.backend = `my fancy!"backend`;
    const backend = encodeURIComponent(this.backend);
    const path = encodeURIComponent(this.path);
    const paths = [
      {
        label: 'API path',
        expected: `/v1/${backend}/data/${path}`,
      },
      {
        label: 'Configuration path',
        expected: `${backend}/data/${path}`,
      },
      { label: 'CLI path', expected: `-mount="${this.backend}" "${this.path}"` },
      {
        label: 'API path for metadata',
        expected: `/v1/${backend}/metadata/${path}`,
      },
    ];

    await this.renderComponent();

    for (const [index, path] of paths.entries()) {
      assert.dom(PAGE.infoRowValue(path.label)).hasText(path.expected);
      await this.assertClipboard(assert, path.expected, index);
    }
  });

  test('it renders copyable commands', async function (assert) {
    const url = `https://127.0.0.1:8200/v1/${this.backend}/data/${this.path}`;
    const expected = {
      cli: `vault kv get -mount="${this.backend}" "${this.path}"`,
      api: `curl \\ --header \"X-Vault-Token: ...\" \\ --request GET \\ ${url}`,
    };
    await this.renderComponent();

    assert.dom(PAGE.paths.codeSnippet('cli')).hasText(expected.cli);
    assert.dom(PAGE.paths.codeSnippet('api')).hasText(expected.api);
  });

  test('it renders copyable encoded mount and path commands', async function (assert) {
    this.path = `my spacey!"secret`;
    this.backend = `my fancy!"backend`;

    const backend = encodeURIComponent(this.backend);
    const path = encodeURIComponent(this.path);
    const url = `https://127.0.0.1:8200/v1/${backend}/data/${path}`;

    const expected = {
      cli: `vault kv get -mount="${this.backend}" "${this.path}"`,
      api: `curl \\ --header \"X-Vault-Token: ...\" \\ --request GET \\ ${url}`,
    };
    await this.renderComponent();

    assert.dom(PAGE.paths.codeSnippet('cli')).hasText(expected.cli);
    assert.dom(PAGE.paths.codeSnippet('api')).hasText(expected.api);
  });
});
