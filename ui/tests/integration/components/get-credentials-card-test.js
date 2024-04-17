/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import Service from '@ember/service';
import { click, render } from '@ember/test-helpers';
import { selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const storeService = Service.extend({
  query(modelType) {
    return new Promise((resolve, reject) => {
      switch (modelType) {
        case 'database/role':
          resolve([{ id: 'my-role', backend: 'database' }]);
          break;
        default:
          reject({ httpStatus: 404, message: 'not found' });
          break;
      }
      reject({ httpStatus: 404, message: 'not found' });
    });
  },
});

module('Integration | Component | get-credentials-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.router.transitionTo = sinon.stub();

    this.owner.unregister('service:store');
    this.owner.register('service:store', storeService);
    this.set('title', 'Get Credentials');
    this.set('searchLabel', 'Role to use');
    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
      },
    });
  });

  hooks.afterEach(function () {
    this.router.transitionTo.reset();
  });

  test('it shows a disabled button when no item is selected', async function (assert) {
    assert.expect(2);
    await render(hbs`<GetCredentialsCard @title={{this.title}} @searchLabel={{this.searchLabel}}/>`);
    assert.dom('[data-test-get-credentials]').isDisabled();
    assert.dom('[data-test-get-credentials]').hasText('Get credentials', 'Button has default text');
  });

  test('it shows button that can be clicked to credentials route when an item is selected', async function (assert) {
    assert.expect(4);
    const models = ['database/role'];
    this.set('models', models);
    await render(
      hbs`<GetCredentialsCard @title={{this.title}} @searchLabel={{this.searchLabel}} @placeholder="Search for a role..." @models={{this.models}} />`
    );
    assert
      .dom('[data-test-component="search-select"]#search-input-role')
      .exists('renders search select component by default');
    assert
      .dom('[data-test-component="search-select"]#search-input-role')
      .hasText('Search for a role...', 'renders placeholder text passed to search select');
    await clickTrigger();
    await selectChoose('', 'my-role');
    assert.dom('[data-test-get-credentials]').isEnabled();
    await click('[data-test-get-credentials]');
    assert.propEqual(
      this.router.transitionTo.lastCall.args,
      ['vault.cluster.secrets.backend.credentials', 'my-role'],
      'transitionTo is called with correct route and role name'
    );
  });
});
