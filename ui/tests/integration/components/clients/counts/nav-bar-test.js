/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | clients/counts/nav-bar', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = async () => {
      await render(hbs`<Clients::Counts::NavBar />`);
    };
  });

  test('it renders default tabs', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.tab('overview')).hasText('Overview');
    assert.dom(GENERAL.tab('client list')).hasText('Client list');
  });
});
