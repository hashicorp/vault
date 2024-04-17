/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PERMISSIONS_BANNER_STATES } from 'vault/services/permissions';

const TEXT = {
  titleReadFail: 'Resultant ACL check failed',
  titleNoAccess: 'You do not have access to this namespace',
  messageReadFail:
    "Links might be shown that you don't have access to. Contact your administrator to update your policy.",
  messageNoAccess:
    'Log into the namespace directly, or contact your administrator if you think you should have access.',
};
module('Integration | Component | resultant-acl-banner', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders correctly by default', async function (assert) {
    await render(hbs`<ResultantAclBanner />`);

    assert.dom('[data-test-resultant-acl-banner] .hds-alert__title').hasText(TEXT.titleReadFail);
    assert.dom('[data-test-resultant-acl-banner] .hds-alert__description').hasText(TEXT.messageReadFail);
    assert.dom('[data-test-resultant-acl-reauthenticate]').doesNotExist('Does not show reauth link');
  });

  test('it renders correctly with set namespace', async function (assert) {
    const nsService = this.owner.lookup('service:namespace');
    nsService.setNamespace('my-ns');
    this.set('failType', undefined);

    await render(hbs`<ResultantAclBanner @isEnterprise={{true}} @failType={{this.failType}} />`);

    assert
      .dom('[data-test-resultant-acl-banner] .hds-alert__title')
      .hasText(TEXT.titleReadFail, 'title correct for default fail type');
    assert
      .dom('[data-test-resultant-acl-banner] .hds-alert__description')
      .hasText(TEXT.messageReadFail, 'message correct for default fail type');
    assert
      .dom('[data-test-resultant-acl-reauthenticate]')
      .hasText('Log into my-ns namespace', 'Shows reauth link with given namespace');

    this.set('failType', PERMISSIONS_BANNER_STATES.noAccess);
    assert
      .dom('[data-test-resultant-acl-banner] .hds-alert__title')
      .hasText(TEXT.titleNoAccess, 'title correct for no access failtype');
    assert
      .dom('[data-test-resultant-acl-banner] .hds-alert__description')
      .hasText(TEXT.messageNoAccess, 'message correct for no access failtype');
  });

  test('it renders correctly with default namespace', async function (assert) {
    await render(hbs`<ResultantAclBanner @isEnterprise={{true}} />`);
    assert
      .dom('[data-test-resultant-acl-reauthenticate]')
      .hasText('Log into root namespace', 'Shows reauth link with default namespace');
  });

  test('it goes away when dismiss button clicked', async function (assert) {
    await render(hbs`<ResultantAclBanner />`);
    assert.dom('[data-test-resultant-acl-banner]').exists('Shows banner initially');
    await click('.hds-dismiss-button');
    assert.dom('[data-test-resultant-acl-banner]').doesNotExist('Hides banner after dismiss');
  });
});
