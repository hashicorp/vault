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
import { kvDataPath, kvMetadataPath } from 'vault/utils/kv-path';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';

module('Integration | Component | kv-v2 | Page::Secret::Metadata::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.dataId = kvDataPath(this.backend, this.path);
    this.metadataId = kvMetadataPath(this.backend, this.path);

    this.metadataModel = (withCustom = false) => {
      const metadata = withCustom
        ? this.server.create('kv-metadatum', 'withCustomMetadata')
        : this.server.create('kv-metadatum');
      metadata.id = this.metadataId;
      this.store.pushPayload('kv/metadata', {
        modelName: 'kv/metadata',
        ...metadata,
      });
      return this.store.peekRecord('kv/metadata', this.metadataId);
    };

    this.metadata = this.metadataModel();

    // empty secret model always exists for permissions
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.dataId,
      custom_metadata: null,
    });
    this.secret = this.store.peekRecord('kv/data', this.dataId);

    // this is the route model, not an ember data model
    this.model = {
      backend: this.backend,
      path: this.path,
      secret: this.secret,
      metadata: this.metadata,
    };
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path },
    ];
  });

  test('it renders metadata details', async function (assert) {
    assert.expect(8);
    this.metadata = this.metadataModel();
    await render(
      hbs`
       <Page::Secret::Metadata::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.title).includesText(this.model.path, 'renders secret path as page title');
    assert.dom(PAGE.emptyStateTitle).hasText('No custom metadata', 'renders the correct empty state');
    assert.dom(PAGE.metadata.deleteMetadata).exists();
    assert.dom(PAGE.metadata.editBtn).exists();

    // Metadata details
    assert
      .dom(PAGE.infoRowValue('Last updated'))
      .hasTextContaining('Mar', 'Displays updated date with formatting');
    assert.dom(PAGE.infoRowValue('Maximum versions')).hasText('15');
    assert.dom(PAGE.infoRowValue('Check-and-Set required')).hasText('Yes');
    assert
      .dom(PAGE.infoRowValue('Delete version after'))
      .hasText('3 hours 25 minutes 19 seconds', 'correctly shows and formats the timestamp.');
  });

  test('it renders custom metadata from secret model', async function (assert) {
    assert.expect(2);
    this.metadata = this.metadataModel();
    this.secret.customMetadata = { hi: 'there' };
    await render(
      hbs`
       <Page::Secret::Metadata::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.emptyStateTitle).doesNotExist();
    assert.dom(PAGE.infoRowValue('hi')).hasText('there', 'renders custom metadata from secret');
  });

  test('it renders custom metadata from metadata model', async function (assert) {
    assert.expect(4);
    this.metadata = this.metadataModel({ withCustom: true });
    await render(
      hbs`
       <Page::Secret::Metadata::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.emptyStateTitle).doesNotExist();
    // Metadata details
    assert.dom(PAGE.infoRowValue('foo')).hasText('abc');
    assert.dom(PAGE.infoRowValue('bar')).hasText('123');
    assert.dom(PAGE.infoRowValue('baz')).hasText('5c07d823-3810-48f6-a147-4c06b5219e84');
  });
});
