/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvDataPath } from 'vault/utils/kv-path';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
/* eslint-disable no-useless-escape */

module('Integration | Component | kv-v2 | Page::Secret::Paths', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.dataId = kvDataPath(this.backend, this.path);
    this.secretData = { foo: 'bar' };
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.dataId,
      path: this.path,
      backend: this.backend,
      version: 2,
    });

    this.secret = this.store.peekRecord('kv/data', this.dataId);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: this.path },
    ];

    this.assertClipboard = (assert, element, expected) => {
      assert.dom(element).hasAttribute('data-clipboard-text', expected);
    };
  });

  test('it renders copyable paths', async function (assert) {
    assert.expect(6);

    const paths = [
      { label: 'API path', expected: `/v1/${this.backend}/data/${this.path}` },
      { label: 'CLI path', expected: `-mount="${this.backend}" "${this.path}"` },
      { label: 'API path for metadata', expected: `/v1/${this.backend}/metadata/${this.path}` },
    ];

    await render(
      hbs`
<Page::Secret::Paths
  @path={{this.path}}
  @secret={{this.secret}}
  @breadcrumbs={{this.breadcrumbs}}
/>
      `,
      { owner: this.engine }
    );

    for (const path of paths) {
      assert.dom(PAGE.infoRowValue(path.label)).hasText(path.expected);
      this.assertClipboard(assert, PAGE.paths.copyButton(path.label), path.expected);
    }
  });

  test('it renders copyable encoded mount and secret paths', async function (assert) {
    assert.expect(6);
    this.path = `my spacey!"secret`;
    this.backend = `my fancy!"backend`;
    this.secret.backend = this.backend;
    this.secret.path = this.path;

    const paths = [
      {
        label: 'API path',
        expected: `/v1/${encodeURIComponent(this.backend)}/data/${encodeURIComponent(this.path)}`,
      },
      { label: 'CLI path', expected: `-mount="${this.backend}" "${this.path}"` },
      {
        label: 'API path for metadata',
        expected: `/v1/${encodeURIComponent(this.backend)}/metadata/${encodeURIComponent(this.path)}`,
      },
    ];

    await render(
      hbs`
<Page::Secret::Paths
  @path={{this.path}}
  @secret={{this.secret}}
  @breadcrumbs={{this.breadcrumbs}}
/>
      `,
      { owner: this.engine }
    );

    for (const path of paths) {
      assert.dom(PAGE.infoRowValue(path.label)).hasText(path.expected);
      this.assertClipboard(assert, PAGE.paths.copyButton(path.label), path.expected);
    }
  });

  test('it renders copyable commands', async function (assert) {
    assert.expect(4);
    const url = `http://127.0.0.1:8200/v1/${this.backend}/data/${this.path}?version=2`;
    const expected = {
      cli: `vault kv get -mount="${this.backend}" "${this.path}"`,
      apiDisplay: `curl \\ --header \"X-Vault-Token: ...\" \\ --request GET \\ ${url}`,
      apiCopy: `curl  --header \"X-Vault-Token: ...\"  --request GET \ ${url}`,
    };
    await render(
      hbs`
<Page::Secret::Paths
  @path={{this.path}}
  @secret={{this.secret}}
  @breadcrumbs={{this.breadcrumbs}}
/>
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.paths.codeSnippet('cli')).hasText(expected.cli);
    assert.dom(PAGE.paths.snippetCopy('cli')).hasAttribute('data-clipboard-text', expected.cli);
    assert.dom(PAGE.paths.codeSnippet('api')).hasText(expected.apiDisplay);
    assert.dom(PAGE.paths.snippetCopy('api')).hasAttribute('data-clipboard-text', expected.apiCopy);
  });

  test('it renders copyable encoded commands', async function (assert) {
    assert.expect(4);
    this.path = `my spacey!"secret`;
    this.backend = `my fancy!"backend`;
    this.secret.backend = this.backend;
    this.secret.path = this.path;

    const backend = encodeURIComponent(this.backend);
    const path = encodeURIComponent(this.path);
    const url = `http://127.0.0.1:8200/v1/${backend}/data/${path}?version=2`;

    const expected = {
      cli: `vault kv get -mount="${this.backend}" "${this.path}"`,
      apiDisplay: `curl \\ --header \"X-Vault-Token: ...\" \\ --request GET \\ ${url}`,
      apiCopy: `curl  --header \"X-Vault-Token: ...\"  --request GET \ ${url}`,
    };
    await render(
      hbs`
<Page::Secret::Paths
  @path={{this.path}}
  @secret={{this.secret}}
  @breadcrumbs={{this.breadcrumbs}}
/>
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.paths.codeSnippet('cli')).hasText(expected.cli);
    assert.dom(PAGE.paths.snippetCopy('cli')).hasAttribute('data-clipboard-text', expected.cli);
    assert.dom(PAGE.paths.codeSnippet('api')).hasText(expected.apiDisplay);
    assert.dom(PAGE.paths.snippetCopy('api')).hasAttribute('data-clipboard-text', expected.apiCopy);
  });
});
