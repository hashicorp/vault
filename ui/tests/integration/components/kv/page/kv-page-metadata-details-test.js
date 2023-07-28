/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvMetadataPath } from 'vault/utils/kv-path';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { SELECTORS } from 'vault/tests/helpers/kv/kv-general-selectors';
import { PAGE } from 'vault/tests/helpers/kv/kv-page-selectors';

module('Integration | Component | kv | Page::Secret::MetadataDetails', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it renders metadata details component and shows empty state when custom_metadata is empty', async function (assert) {
    assert.expect(2);
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    const data = this.server.create('kv-metadatum');
    data.id = kvMetadataPath('kv-engine', 'my-secret');
    this.store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...data,
    });
    this.model = this.store.peekRecord('kv/metadata', data.id);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path, route: 'secret', model: this.model.path },
      { label: 'metadata' },
    ];
    await render(
      hbs`<Page::Secret::MetadataDetails @metadata={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(SELECTORS.emptyStateTitle).hasText('No custom metadata', 'renders the correct empty state');
    assert
      .dom(SELECTORS.infoRowValue('Delete version after'))
      .hasText('3 hours 25 minutes 19 seconds', 'correctly shows and formats the timestamp.');
  });

  test('it renders custom metadata when it exists and user has permissions', async function (assert) {
    assert.expect(3);
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    const data = this.server.create('kv-metadatum', 'withCustomMetadata');
    data.id = kvMetadataPath('kv-engine', 'my-secret');
    this.store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...data,
    });
    this.model = this.store.peekRecord('kv/metadata', data.id);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path, route: 'secret', model: this.model.path },
      { label: 'metadata' },
    ];

    await render(
      hbs`<Page::Secret::MetadataDetails @metadata={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
      {
        owner: this.engine,
      }
    );
    for (const key in this.model.customMetadata) {
      const value = this.model.customMetadata[key];
      assert.dom(SELECTORS.infoRowValue(key)).hasText(value);
    }
  });

  test('it renders correct empty state messages with no READ metadata permissions and no secret.customMetadata is returned', async function (assert) {
    assert.expect(3);
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub('list', 'update'));
    // would not return custom_metadata if they did not have permissions
    const data = this.server.create('kv-metadatum');
    data.id = kvMetadataPath('kv-engine', 'my-secret');
    this.store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...data,
    });
    this.model = this.store.peekRecord('kv/metadata', data.id);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path, route: 'secret', model: this.model.path },
      { label: 'metadata' },
    ];
    await render(
      hbs`<Page::Secret::MetadataDetails @metadata={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
      {
        owner: this.engine,
      }
    );

    const [noCustomMetadata, noMetadata] = findAll(SELECTORS.emptyStateTitle);
    assert
      .dom(noCustomMetadata)
      .exists(
        'You do not have access to read custom metadata',
        'renders the empty state about custom_metadata'
      );
    assert
      .dom(noMetadata)
      .exists(
        'You do not have access to secret metadata',
        'renders the empty state about no secret metadata'
      );
    assert.dom(PAGE.details.editMetadataBtn).doesNotExist('does not render edit metadata button.');
  });
});
