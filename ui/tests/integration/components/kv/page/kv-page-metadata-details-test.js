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
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { baseSetup, metadataModel } from 'vault/tests/helpers/kv/kv-run-commands';
import { dateFormat } from 'core/helpers/date-format';

module('Integration | Component | kv-v2 | Page::Secret::Metadata::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    // this.metadata is setup by baseSetup
    baseSetup(this);

    // this is the route model, not an ember data model
    this.model = {
      backend: this.backend,
      path: this.path,
      secret: this.secret,
      metadata: this.metadata,
      canDeleteMetadata: true,
      canReadData: true,
      canReadCustomMetadata: true,
      canReadMetadata: true,
      canUpdateMetadata: true,
    };
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path },
    ];

    this.renderComponent = () => {
      return render(
        hbs`
       <Page::Secret::Metadata::Details
        @backend={{this.model.backend}}
        @breadcrumbs={{this.breadcrumbs}}
        @canDeleteMetadata={{this.model.canDeleteMetadata}}
        @canReadData={{this.model.canReadData}}
        @canReadMetadata={{this.model.canReadMetadata}}
        @canUpdateMetadata={{this.model.canUpdateMetadata}}
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

  test('it renders empty state if cannot read metadata but can read data', async function (assert) {
    this.model.metadata = null;
    await this.renderComponent();
    assert
      .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
      .hasText('Request custom metadata?');
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

  test('it hides delete modal when no permissions', async function (assert) {
    this.model.canDeleteMetadata = false;
    assert.dom(PAGE.metadata.deleteMetadata).doesNotExist();
  });

  test('it hides edit action when no permissions', async function (assert) {
    this.model.canUpdateMetadata = false;
    assert.dom(PAGE.metadata.editBtn).doesNotExist();
  });
});
