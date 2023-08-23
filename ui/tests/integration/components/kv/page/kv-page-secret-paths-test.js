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
      secret_data: this.secretData,
      created_time: '2023-07-20T02:12:17.379762Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: 2,
    });

    this.secret = this.store.peekRecord('kv/data', this.dataId);

    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: this.path },
    ];
  });

  test('it renders copyable paths', async function (assert) {
    assert.expect(6);

    const expected = {
      api: `/v1/${this.backend}/data/${this.path}`,
      cli: `-mount=${this.backend} "${this.path}"`,
      apiMeta: `/v1/${this.backend}/metadata/${this.path}`,
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
    assert.dom(PAGE.infoRowValue('API path')).hasText(expected.api);
    assert.dom(PAGE.infoRowValue('CLI path')).hasText(expected.cli);
    assert.dom(PAGE.infoRowValue('API path for metadata')).hasText(expected.apiMeta);
    assert.dom(`${PAGE.infoRowValue('API path')} button`).hasAttribute('data-clipboard-text', expected.api);
    assert.dom(`${PAGE.infoRowValue('CLI path')} button`).hasAttribute('data-clipboard-text', expected.cli);
    assert
      .dom(`${PAGE.infoRowValue('API path for metadata')} button`)
      .hasAttribute('data-clipboard-text', expected.apiMeta);
  });

  test('it renders copyable encoded paths', async function (assert) {
    assert.expect(6);
    this.path = 'my spacey secret';
    this.secret.path = this.path;
    const expected = {
      api: `/v1/${this.backend}/data/${encodeURIComponent(this.path)}`,
      cli: `-mount=${this.backend} "${this.path}"`,
      apiMeta: `/v1/${this.backend}/metadata/${encodeURIComponent(this.path)}`,
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

    assert.dom(PAGE.infoRowValue('API path')).hasText(expected.api);
    assert.dom(PAGE.infoRowValue('CLI path')).hasText(expected.cli);
    assert.dom(PAGE.infoRowValue('API path for metadata')).hasText(expected.apiMeta);
    assert.dom(`${PAGE.infoRowValue('API path')} button`).hasAttribute('data-clipboard-text', expected.api);
    assert.dom(`${PAGE.infoRowValue('CLI path')} button`).hasAttribute('data-clipboard-text', expected.cli);
    assert
      .dom(`${PAGE.infoRowValue('API path for metadata')} button`)
      .hasAttribute('data-clipboard-text', expected.apiMeta);
  });

  test('it renders copyable commands', async function (assert) {
    assert.expect(6);

    await render(
      hbs`
<Page::Secret::Paths
  @path={{this.model.path}}
  @secret={{this.model.secret}}
  @breadcrumbs={{this.breadcrumbs}}
/>
      `,
      { owner: this.engine }
    );

    assert.dom('[data-test-code-snippet]').hasText('');
    assert.dom('[data-test-code-snippet] button').hasAttribute('data-clipboard-text', '');
  });

  test('it renders copyable encoded commands', async function (assert) {
    assert.expect(6);

    await render(
      hbs`
<Page::Secret::Paths
  @path={{this.model.path}}
  @secret={{this.model.secret}}
  @breadcrumbs={{this.breadcrumbs}}
/>
      `,
      { owner: this.engine }
    );
    assert.dom('[data-test-code-snippet]').hasText('');
    assert.dom('[data-test-code-snippet] button').hasAttribute('data-clipboard-text', '');
  });
});
