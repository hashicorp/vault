/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | auth-method/configuration', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.createModel = (path, type) => {
      this.model = this.store.createRecord('auth-method', { path, type });
      this.model.set('config', this.store.createRecord('mount-config'));
    };
    this.renderComponent = async () => await render(hbs`<AuthMethod::Configuration @model={{this.model}} />`);
  });

  test('it renders direct link for supported method', async function (assert) {
    this.createModel('token/', 'token');
    await this.renderComponent();
    assert.dom(GENERAL.infoRowValue('UI login link')).hasText(`${window.origin}/ui/vault/auth?with=token%2F`);
  });

  test('it does not render direct link for unsupported method', async function (assert) {
    this.createModel('my-approle/', 'approle');
    await this.renderComponent();
    assert.dom(GENERAL.infoRowValue('UI login link')).doesNotExist();
  });

  test('it renders direct link if within a namespace', async function (assert) {
    this.owner.lookup('service:namespace').set('path', 'foo/bar');
    this.createModel('token/', 'token');
    await this.renderComponent();
    assert
      .dom(GENERAL.infoRowValue('UI login link'))
      .hasText(`${window.origin}/ui/vault/auth?namespace=foo%2Fbar&with=token%2F`);
  });
});
