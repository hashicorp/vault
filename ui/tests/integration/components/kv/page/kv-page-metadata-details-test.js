/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvDataPath } from 'vault/utils/kv-path';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { baseSetup, metadataModel } from 'vault/tests/helpers/kv/kv-run-commands';
import { dateFormat } from 'core/helpers/date-format';

module('Integration | Component | kv-v2 | Page::Secret::Metadata::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    baseSetup(this);
    this.dataId = kvDataPath(this.backend, this.path);
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
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path },
    ];
    this.canDeleteMetadata = true;
    this.canReadCustomMetadata = true;
    this.canReadMetadata = true;
    this.canUpdateMetadata = true;

    this.renderComponent = () => {
      return render(
        hbs`
       <Page::Secret::Metadata::Details
        @breadcrumbs={{this.breadcrumbs}}
        @canDeleteMetadata={{this.canDeleteMetadata}}
        @canReadMetadata={{this.canReadMetadata}}
        @canUpdateMetadata={{this.canReadMetadata}}
        @customMetadata={{or this.model.metadata.customMetadata this.model.secret.customMetadata}}
        @metadata={{this.model.metadata}}
        @path={{this.model.path}}
      />
      `,
        { owner: this.engine }
      );
    };
  });

  test('it renders metadata details', async function (assert) {
    assert.expect(8);
    await this.renderComponent();

    assert.dom(PAGE.title).includesText(this.model.path, 'renders secret path as page title');
    assert.dom(PAGE.emptyStateTitle).hasText('No custom metadata', 'renders the correct empty state');
    assert.dom(PAGE.metadata.deleteMetadata).exists();
    assert.dom(PAGE.metadata.editBtn).exists();

    // Metadata details
    const expectedTime = dateFormat([this.metadata.updatedTime, 'MMM d, yyyy hh:mm aa'], {});
    assert
      .dom(PAGE.infoRowValue('Last updated'))
      .hasTextContaining(expectedTime, 'Displays updated date with formatting');
    assert.dom(PAGE.infoRowValue('Maximum versions')).hasText('15');
    assert.dom(PAGE.infoRowValue('Check-and-Set required')).hasText('Yes');
    assert
      .dom(PAGE.infoRowValue('Delete version after'))
      .hasText('3 hours 25 minutes 19 seconds', 'correctly shows and formats the timestamp.');
  });

  test('it renders custom metadata from secret model', async function (assert) {
    assert.expect(2);
    this.secret.customMetadata = { hi: 'there' };
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).doesNotExist();
    assert.dom(PAGE.infoRowValue('hi')).hasText('there', 'renders custom metadata from secret');
  });

  test('it renders custom metadata from metadata model', async function (assert) {
    assert.expect(4);
    this.model.metadata = metadataModel(this, { withCustom: true });
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).doesNotExist();
    // Metadata details
    assert.dom(PAGE.infoRowValue('foo')).hasText('abc');
    assert.dom(PAGE.infoRowValue('bar')).hasText('123');
    assert.dom(PAGE.infoRowValue('baz')).hasText('5c07d823-3810-48f6-a147-4c06b5219e84');
  });

  test('it renders custom metadata from metadata if secret data exists', async function (assert) {
    assert.expect(4);
    this.secret.customMetadata = { hi: 'there' };
    this.model.metadata = metadataModel(this, { withCustom: true });
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).doesNotExist();
    // Metadata details
    assert.dom(PAGE.infoRowValue('foo')).hasText('abc');
    assert.dom(PAGE.infoRowValue('bar')).hasText('123');
    assert.dom(PAGE.infoRowValue('baz')).hasText('5c07d823-3810-48f6-a147-4c06b5219e84');
  });

  test('it hides delete modal when no permissions', async function (assert) {
    this.canDeleteMetadata = false;
    assert.dom(PAGE.metadata.deleteMetadata).doesNotExist();
  });

  test('it hides edit action when no permissions', async function (assert) {
    this.canUpdateMetadata = false;
    assert.dom(PAGE.metadata.editBtn).doesNotExist();
  });
});
