/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | resultant-acl-banner', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders correctly by default', async function (assert) {
    await render(hbs`<ResultantAclBanner />`);

    assert.dom('[data-test-resultant-acl-banner] .hds-alert__title').hasText('Resultant ACL check failed');
    assert
      .dom('[data-test-resultant-acl-banner] .hds-alert__description')
      .hasText(
        "Links might be shown that you don't have access to. Contact your administrator to update your policy."
      );
    assert.dom('[data-test-resultant-acl-reauthenticate]').doesNotExist('Does not show reauth link');
  });

  test('it renders correctly with set namespace', async function (assert) {
    const nsService = this.owner.lookup('service:namespace');
    nsService.setNamespace('my-ns');

    await render(hbs`<ResultantAclBanner @isEnterprise={{true}} />`);

    assert.dom('[data-test-resultant-acl-banner] .hds-alert__title').hasText('Resultant ACL check failed');
    assert
      .dom('[data-test-resultant-acl-banner] .hds-alert__description')
      .hasText('You do not have access to resources in this namespace.');
    assert
      .dom('[data-test-resultant-acl-reauthenticate]')
      .hasText('Log into my-ns namespace', 'Shows reauth link with given namespace');
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
