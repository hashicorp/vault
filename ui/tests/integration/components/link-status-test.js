/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { statuses } from '../../../mirage/handlers/hcp-link';

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
      <div id="modal-wormhole"></div>
      <LinkStatus @status={{undefined}} />
    `);

    assert.dom('.navbar-status').doesNotExist('Banner is hidden for missing status message');
  });

  test('it does not render banner in oss version', async function (assert) {
    this.owner.lookup('service:version').set('version', '1.13.0');

    await render(hbs`
      <div id="modal-wormhole"></div>
      <LinkStatus @status={{get this.statuses 0}} />
    `);

    assert.dom('.navbar-status').doesNotExist('Banner is hidden in oss');
  });

  test('it renders connected status', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <LinkStatus @status={{get this.statuses 0}} />
    `);

    assert.dom('.navbar-status').hasClass('connected', 'Correct banner class renders for connected state');
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
      <div id="modal-wormhole"></div>
      <LinkStatus @status={{get this.statuses 1}} />
    `);

    assert.dom('.navbar-status').hasClass('warning', 'Correct banner class renders for error state');
    assert
      .dom('[data-test-link-status]')
      .hasText(
        'There was an error connecting to HCP. Click here for more information.',
        'Banner copy renders for error state'
      );

    await click('[data-test-link-status] button');
    assert
      .dom('[data-test-link-status-timestamp]')
      .hasText('2022-09-21T11:25:02.196835-07:00', 'Timestamp renders');
    assert
      .dom('[data-test-link-status-error]')
      .hasText('unable to establish a connection with HCP', 'Error renders');

    // connecting error
    await render(hbs`
      <div id="modal-wormhole"></div>
      <LinkStatus @status={{get this.statuses 3}} />
    `);
    assert
      .dom('[data-test-link-status-error]')
      .hasText('principal does not have the permission to register as a provider', 'Error renders');

    // this shouldn't happen but placeholders should render if disconnected/connecting status is returned without timestamp and/or error
    await render(hbs`
      <div id="modal-wormhole"></div>
      <LinkStatus @status="connecting" />
    `);
    assert.dom('[data-test-link-status-timestamp]').hasText('Not available', 'Timestamp placeholder renders');
    assert.dom('[data-test-link-status-error]').hasText('Not available', 'Error placeholder renders');
  });
});
