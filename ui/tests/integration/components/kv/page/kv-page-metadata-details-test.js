/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { dateFormat } from 'core/helpers/date-format';

module('Integration | Component | kv-v2 | Page::Secret::Metadata::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.metadata = {
      custom_metadata: null,
      max_versions: 15,
      cas_required: true,
      delete_version_after: '4h30m',
      updated_time: '2025-09-16T22:19:59.916935Z',
    };
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.capabilities = {
      canReadMetadata: true,
      canUpdateMetadata: true,
      canDeleteMetadata: true,
      canReadData: true,
    };
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: this.path },
    ];

    this.renderComponent = () =>
      render(
        hbs`
          <Page::Secret::Metadata::Details
            @backend={{this.backend}}
            @breadcrumbs={{this.breadcrumbs}}
            @capabilities={{this.capabilities}}
            @metadata={{this.metadata}}
            @path={{this.path}}
          />
        `,
        { owner: this.engine }
      );
  });

  test('it should render page title and toolbar elements', async function (assert) {
    await this.renderComponent();

    assert.dom(PAGE.title).includesText(this.path, 'renders secret path as page title');
    assert.dom(PAGE.secretTab('Overview')).exists('renders Secrets tab');
    assert.dom(PAGE.secretTab('Secret')).exists('renders Secret tab');
    assert.dom(PAGE.secretTab('Metadata')).exists('renders Metadata tab');
    assert.dom(PAGE.secretTab('Paths')).exists('renders Paths tab');
    assert.dom(PAGE.secretTab('Version History')).exists('renders Version History tab');
  });

  test('it renders metadata details', async function (assert) {
    assert.expect(7);
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).hasText('No custom metadata', 'renders the correct empty state');
    assert.dom(PAGE.metadata.deleteMetadata).exists();
    assert.dom(PAGE.metadata.editBtn).exists();

    // Metadata details
    const expectedTime = dateFormat([this.metadata.updated_time, 'MMM d, yyyy hh:mm aa'], {});
    assert
      .dom(PAGE.infoRowValue('Last updated'))
      .hasTextContaining(expectedTime, 'Displays updated date with formatting');
    assert.dom(PAGE.infoRowValue('Maximum versions')).hasText('15');
    assert.dom(PAGE.infoRowValue('Check-and-Set required')).hasText('Yes');
    assert
      .dom(PAGE.infoRowValue('Delete version after'))
      .hasText('4 hours 30 minutes', 'correctly shows and formats the timestamp.');
  });

  test('it renders empty state if cannot read metadata but can read data', async function (assert) {
    // this.metadata = null;
    this.capabilities.canReadMetadata = false;
    this.capabilities.canReadData = true;
    await this.renderComponent();
    assert
      .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
      .hasText('Request custom metadata?');
  });

  test('it renders custom metadata', async function (assert) {
    assert.expect(3);
    this.metadata.custom_metadata = { foo: 'bar', bar: 'baz' };
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).doesNotExist();
    // Metadata details
    assert.dom(PAGE.infoRowValue('foo')).hasText('bar');
    assert.dom(PAGE.infoRowValue('bar')).hasText('baz');
  });

  test('it hides delete modal when no permissions', async function (assert) {
    this.capabilities.canDeleteMetadata = false;
    await this.renderComponent();
    assert.dom(PAGE.metadata.deleteMetadata).doesNotExist();
  });

  test('it hides edit action when no permissions', async function (assert) {
    this.capabilities.canUpdateMetadata = false;
    await this.renderComponent();
    assert.dom(PAGE.metadata.editBtn).doesNotExist();
  });
});
