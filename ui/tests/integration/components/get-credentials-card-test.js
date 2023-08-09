/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import Service from '@ember/service';
import { click, find, render, typeIn } from '@ember/test-helpers';
import { selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

const TITLE = 'Get Credentials';
const SEARCH_LABEL = 'Role to use';

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
    this.set('title', TITLE);
    this.set('searchLabel', SEARCH_LABEL);
  });

  hooks.afterEach(function () {
    this.router.transitionTo.reset();
  });

  test('it shows a disabled button when no item is selected', async function (assert) {
    await render(hbs`<GetCredentialsCard @title={{this.title}} @searchLabel={{this.searchLabel}}/>`);
    assert.dom('[data-test-get-credentials]').isDisabled();
  });

  test('it shows button that can be clicked to credentials route when an item is selected', async function (assert) {
    const models = ['database/role'];
    this.set('models', models);
    await render(
      hbs`<GetCredentialsCard @title={{this.title}} @searchLabel={{this.searchLabel}} @placeholder="Search for a role..." @models={{this.models}} @type="role"/>`
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

  test('it renders input search field when renderInputSearch=true and shows placeholder text', async function (assert) {
    await render(
      hbs`<GetCredentialsCard @title={{this.title}} @renderInputSearch={{true}} @placeholder="secret/" @backend="kv" @type="secret"/>`
    );
    assert
      .dom('[data-test-component="search-select"]')
      .doesNotExist('does not render search select component');
    assert.strictEqual(
      find('[data-test-search-roles] input').placeholder,
      'secret/',
      'renders placeholder text passed to search input'
    );
    await typeIn('[data-test-search-roles] input', 'test');
    assert.dom('[data-test-get-credentials]').isEnabled('submit button enables after typing input text');
    await click('[data-test-get-credentials]');
    assert.propEqual(
      this.router.transitionTo.lastCall.args,
      ['vault.cluster.secrets.backend.show', 'test'],
      'transitionTo is called with correct route and secret name'
    );
  });
});
