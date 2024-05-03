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
    this.showSecretsSync = false;

    this.renderComponent = async () => {
      await render(hbs`<Clients::Counts::NavBar @showSecretsSync={{this.showSecretsSync}} />`);
    };
  });

  test('it renders default tabs', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.tab('overview')).hasText('Overview');
    assert.dom(GENERAL.tab('token')).hasText('Entity/Non-entity clients');
    assert.dom(GENERAL.tab('acme')).hasText('ACME clients');
  });

  test('it shows secrets sync tab if showSecretsSync is true', async function (assert) {
    this.showSecretsSync = true;
    await this.renderComponent();

    assert.dom(GENERAL.tab('sync')).exists();
  });

  test('it should not show secrets sync tab if showSecretsSync is false', async function (assert) {
    this.showSecretsSync = false;
    await this.renderComponent();

    assert.dom(GENERAL.tab('sync')).doesNotExist();
  });
});
