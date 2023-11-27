/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { statuses } from '../../../mirage/handlers/hcp-link';

const SELECTORS = {
  modalOpen: '[data-test-link-status] button',
  modalClose: '[data-test-icon="x"]',
  bannerConnected: '.hds-alert [data-test-icon="info"]',
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

    assert.dom(SELECTORS.bannerConnected).exists('Success banner renders for connected state');
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
      .hasText('Error connecting to HCP', 'Banner title renders for error state');
    assert
      .dom('[data-test-link-error]')
      .hasText(
        'Since 2022-09-21T11:25:02.196835-07:00, unable to establish a connection with HCP. Check the logs for more information.',
        'Error renders for case 1'
      );

    // unable to establish connection error
    await render(hbs`
      <LinkStatus @status={{get this.statuses 2}} />
    `);
    assert
      .dom('[data-test-link-error]')
      .hasText(
        'Since 2022-09-21T11:25:02.196835-07:00, unable to establish a connection with HCP. Check the logs for more information.',
        'Error renders for case 2'
      );

    // no permissions error
    await render(hbs`
      <LinkStatus @status={{get this.statuses 3}} />
    `);
    assert
      .dom('[data-test-link-error]')
      .hasText(
        'Since 2022-09-21T11:25:02.196835-07:00, principal does not have the permission to register as a provider. Check the logs for more information.',
        'Error renders for case 3'
      );

    // could not obtain token error
    await render(hbs`
      <LinkStatus @status={{get this.statuses 4}} />
    `);
    assert
      .dom('[data-test-link-error]')
      .hasText(
        'Since 2022-09-21T11:25:02.196835-07:00, could not obtain a token with the supplied credentials. Check the logs for more information.',
        'Error renders for case 3'
      );

    await render(hbs`
      <LinkStatus @status="connecting" />
    `);
    assert.dom('[data-test-link-error]').doesNotExist('No errors rendered when link is in connected state');
  });
});
