/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | page/header', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.breadcrumbs = [
      { label: 'Home', route: 'home' },
      { label: 'Details', route: 'home.details' },
      { label: 'Edit item' },
    ];

    await render(
      hbs`<Page::Header @title="Test title" @subtitle="Test subtitle" @description="Test description" @icon="key" @breadcrumbs={{this.breadcrumbs}}>
        <:breadcrumbs>
          <Page::Breadcrumbs @breadcrumbs={{this.breadcrumbs}} />
        </:breadcrumbs>
        <:badges>
          <Hds::Badge @text="Default badge" data-test-badge>Info Badge</Hds::Badge>
        </:badges>
        <:actions>
          <Hds::Button @variant="primary" @text="Manage" data-test-button />
        </:actions>
      </Page::Header>`
    );
    assert.dom(GENERAL.breadcrumbs).exists('renders passed in breadcrumbs');
    assert.dom(GENERAL.badge()).exists('renders passed in badges');
    assert.dom(GENERAL.button()).exists('renders passed in button');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Test title', 'renders passed in title');
    assert.dom(GENERAL.hdsPageHeaderDescription).hasText('Test description', 'renders passed in description');
    assert.dom(GENERAL.hdsPageHeaderSubtitle).hasText('Test subtitle', 'renders passed in subtitle');
  });
});
