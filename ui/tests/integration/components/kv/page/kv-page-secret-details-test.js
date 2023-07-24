/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { click, find, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvDataPath } from 'vault/utils/kv-path';
import { SELECTORS, parseJsonEditor } from 'vault/tests/helpers/kv/kv-general-selectors';

module('Integration | Component | kv | Page::Secret::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.version = 2;
    this.id = kvDataPath(this.backend, this.path);
    this.secretData = { foo: 'bar' };
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      secret_data: this.secretData,
      created_time: '2023-07-20T02:12:17.379762Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: this.version,
    });
    this.secret = this.store.peekRecord('kv/data', this.id);

    // this is the route model, not an ember data model
    this.model = {
      backend: this.backend,
      path: this.path,
      secret: this.secret,
    };
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'secrets' },
      { label: this.model.path },
    ];
  });

  test('it renders secret details and toggles json view', async function (assert) {
    assert.expect(6);
    await render(
      hbs`
       <Page::Secret::Details
        @secretPath={{this.model.path}}
        @secret={{this.model.secret}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(SELECTORS.pageTitle).includesText(this.model.path, 'renders secret path as page title');
    assert.dom(SELECTORS.infoRowValue('foo')).exists('renders row for secret data');
    assert.dom(SELECTORS.infoRowValue('foo')).hasText('***********');
    await click(SELECTORS.toggleMasked);
    assert.dom(SELECTORS.infoRowValue('foo')).hasText('bar', 'renders secret value');
    await click(SELECTORS.toggleJson);
    assert.propEqual(parseJsonEditor(find), this.secretData, 'json editor renders secret data');
    assert.dom(SELECTORS.tooltipTrigger).includesText(this.version, 'renders version');
  });

  test('it renders deleted empty state', async function (assert) {
    assert.expect(2);
    this.secret.deletionTime = '2023-07-23T02:12:17.379762Z';
    await render(
      hbs`
       <Page::Secret::Details
        @secretPath={{this.model.path}}
        @secret={{this.model.secret}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.emptyStateTitle).hasText('Version 2 of this secret has been deleted');
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        'This version has been deleted but can be undeleted. View other versions of this secret by clicking the Version History tab above.'
      );
  });

  test('it renders destroyed empty state', async function (assert) {
    assert.expect(2);
    this.secret.destroyed = true;
    await render(
      hbs`
       <Page::Secret::Details
        @secretPath={{this.model.path}}
        @secret={{this.model.secret}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.emptyStateTitle).hasText('Version 2 of this secret has been permanently destroyed');
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        'A version that has been permanently deleted cannot be restored. You can view other versions of this secret in the Version History tab above.'
      );
  });
});
