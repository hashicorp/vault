/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { statuses } from '../../../mirage/handlers/hcp-link';

const SELECTORS = {
  modalOpen: '[data-test-link-status] button',
  modalClose: '[data-test-icon="x"]',
  bannerSuccess: '.hds-alert [data-test-icon="check-circle"]',
  bannerWarning: '.hds-alert [data-test-icon="alert-triangle"]',
  banner: '[data-test-link-status]',
};

module('Integration | Component | link-status', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  // this can be removed once feature is released for OSS
  hooks.beforeEach(function () {
    this.owner.lookup('service:version').set('version', '1.13.0+ent');
    this.statuses = statuses;
  });

  test('it does not render banner when status is not present', async function (assert) {
    await render(hbs`
      <LinkStatus @status={{undefined}} />
    `);

    assert.dom(SELECTORS.banner).doesNotExist('Banner is hidden for missing status message');
  });

  test('it does not render banner in oss version', async function (assert) {
    this.owner.lookup('service:version').set('version', '1.13.0');

    await render(hbs`
      <LinkStatus @status={{get this.statuses 0}} />
    `);

    assert.dom(SELECTORS.banner).doesNotExist('Banner is hidden in oss');
  });

  test('it renders connected status', async function (assert) {
    await render(hbs`
      <LinkStatus @status={{get this.statuses 0}} />
    `);

    assert.dom(SELECTORS.bannerSuccess).exists('Success banner renders for connected state');
    assert
      .dom('[data-test-link-status]')
      .hasText('This self-managed Vault is linked to HCP.', 'Banner copy renders for connected state');
    assert
      .dom('[data-test-link-status] a')
      .hasAttribute('href', 'https://portal.cloud.hashicorp.com/sign-in', 'HCP sign in link renders');
  });

  test('it should render error states', async function (assert) {
    // disconnected error
    await render(hbs`
      <LinkStatus @status={{get this.statuses 1}} />
    `);

    assert.dom(SELECTORS.bannerWarning).exists('Warning banner renders for error state');
    assert
      .dom('[data-test-link-status]')
      .hasText(
        'There was an error connecting to HCP. Click here for more information.',
        'Banner copy renders for error state'
      );

    await click(SELECTORS.modalOpen);
    assert
      .dom('[data-test-link-status-timestamp]')
      .hasText('2022-09-21T11:25:02.196835-07:00', 'Timestamp renders');
    assert
      .dom('[data-test-link-status-error]')
      .hasText('unable to establish a connection with HCP', 'Error renders');

    await click(SELECTORS.modalClose);
    // connecting error
    await render(hbs`
      <LinkStatus @status={{get this.statuses 3}} />
    `);
    await click(SELECTORS.modalOpen);
    assert
      .dom('[data-test-link-status-error]')
      .hasText('principal does not have the permission to register as a provider', 'Error renders');
    await click(SELECTORS.modalClose);

    // this shouldn't happen but placeholders should render if disconnected/connecting status is returned without timestamp and/or error
    await render(hbs`
      <LinkStatus @status="connecting" />
    `);
    await click(SELECTORS.modalOpen);

    assert.dom('[data-test-link-status-timestamp]').hasText('Not available', 'Timestamp placeholder renders');
    assert.dom('[data-test-link-status-error]').hasText('Not available', 'Error placeholder renders');
  });
});
