/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, typeIn } from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | get-credentials-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.router.transitionTo = sinon.stub();
    this.dbListStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'databaseListRoles')
      .resolves({ keys: ['my-role'] });
    this.staticListStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'databaseListStaticRoles')
      .resolves({ keys: [] });

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
    await render(
      hbs`<GetCredentialsCard @title={{this.title}} @backend="database"  @searchLabel={{this.searchLabel}} @placeholder="Search for a role..." @models={{this.models}} />`
    );
    assert.dom(GENERAL.button('Get credentials')).isDisabled();
    assert.dom(GENERAL.button('Get credentials')).hasText('Get credentials', 'Button has default text');
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
    assert.dom(GENERAL.button('Get credentials')).isEnabled();
    await click(GENERAL.button('Get credentials'));
    assert.propEqual(
      this.router.transitionTo.lastCall.args,
      ['vault.cluster.secrets.backend.credentials', 'my-role'],
      'transitionTo is called with correct route and role name'
    );
  });

  test('the Overview dropdown displays static and dynamic roles together', async function (assert) {
    this.staticListStub.resolves({ keys: ['10-static', '8-static'] });
    this.dbListStub.resolves({ keys: ['10', '8'] });

    await render(hbs`
      <GetCredentialsCard
        @title={{this.title}}
        @searchLabel={{this.searchLabel}}
        @backend="database"
        @models={{array "database/role"}}
      />
    `);
    await clickTrigger();

    assert.dom(GENERAL.searchSelect.options).exists({ count: 4 }, 'All four roles are offered');
    ['10-static', '8-static', '10', '8'].forEach((role, index) => {
      assert.dom(GENERAL.searchSelect.option(index)).hasText(role, `${role} is available`);
    });
    assert.true(
      this.staticListStub.calledOnceWithExactly('database', true),
      'Static roles are listed through the API client'
    );
    assert.true(
      this.dbListStub.calledOnceWithExactly('database', true),
      'Dynamic roles are listed through the API client'
    );
  });

  for (const roleType of ['static', 'dynamic']) {
    test(`an Overview ${roleType} role can be selected to retrieve credentials`, async function (assert) {
      this.staticListStub.resolves({ keys: ['10-static'] });
      this.dbListStub.resolves({ keys: ['10'] });
      const role = roleType === 'static' ? '10-static' : '10';

      await render(hbs`
        <GetCredentialsCard
          @title={{this.title}}
          @searchLabel={{this.searchLabel}}
          @backend="database"
          @models={{array "database/role"}}
        />
      `);
      await click(GENERAL.searchSelect.trigger('search-input-role'));
      await click(GENERAL.searchSelect.option(roleType === 'static' ? 0 : 1));

      assert.dom(GENERAL.searchSelect.selectedOption()).hasText(role, 'The chosen role is displayed');
      assert.dom(GENERAL.button('Get credentials')).isEnabled('Credentials can be retrieved');
      await click(GENERAL.button('Get credentials'));

      assert.propEqual(
        this.router.transitionTo.lastCall.args,
        ['vault.cluster.secrets.backend.credentials', role],
        'Navigation passes the role name without requiring a role type'
      );
    });

    for (const status of [403, 404]) {
      test(`Overview displays ${roleType} roles when the other role list returns ${status}`, async function (assert) {
        const successful = roleType === 'static' ? this.staticListStub : this.dbListStub;
        const unsuccessful = roleType === 'static' ? this.dbListStub : this.staticListStub;
        const role = `${roleType}-role`;
        successful.resolves({ keys: [role] });
        unsuccessful.rejects(getErrorResponse({ errors: ['permission denied'] }, status));

        await render(hbs`
          <GetCredentialsCard
            @title={{this.title}}
            @searchLabel={{this.searchLabel}}
            @backend="database"
            @models={{array "database/role"}}
          />
        `);
        await click(GENERAL.searchSelect.trigger('search-input-role'));
        await click(GENERAL.searchSelect.option());

        assert
          .dom(GENERAL.searchSelect.selectedOption())
          .hasText(role, 'The accessible role remains selectable');
        assert.dom(GENERAL.button('Get credentials')).isEnabled('A failed list does not block retrieval');
      });
    }
  }

  test('Overview allows entering a role name when both role lists are forbidden', async function (assert) {
    this.staticListStub.rejects(getErrorResponse({ errors: ['permission denied'] }, 403));
    this.dbListStub.rejects(getErrorResponse({ errors: ['permission denied'] }, 403));

    await render(hbs`
      <GetCredentialsCard
        @title={{this.title}}
        @searchLabel={{this.searchLabel}}
        @backend="database"
        @models={{array "database/role"}}
      />
    `);
    await typeIn(GENERAL.inputSearch('search-input-role'), '10-static');

    assert.dom(GENERAL.button('Get credentials')).isEnabled('Manual entry does not require list access');
    await click(GENERAL.button('Get credentials'));
    assert.propEqual(
      this.router.transitionTo.lastCall.args,
      ['vault.cluster.secrets.backend.credentials', '10-static'],
      'The manually entered role can be used to retrieve credentials'
    );
  });

  test('Overview allows entering a role name when both role lists are empty', async function (assert) {
    this.staticListStub.rejects(getErrorResponse({}, 404));
    this.dbListStub.rejects(getErrorResponse({}, 404));

    await render(hbs`
      <GetCredentialsCard
        @title={{this.title}}
        @searchLabel={{this.searchLabel}}
        @backend="database"
        @models={{array "database/role"}}
      />
    `);
    await click(GENERAL.searchSelect.trigger('search-input-role'));
    await fillIn(GENERAL.searchSelect.searchInput, '10-static');
    await click(GENERAL.searchSelect.option());

    assert.dom(GENERAL.searchSelect.selectedOption()).hasText('10-static', 'The entered role is selected');
    assert.dom(GENERAL.button('Get credentials')).isEnabled('An empty list does not block manual entry');
    await click(GENERAL.button('Get credentials'));
    assert.propEqual(
      this.router.transitionTo.lastCall.args,
      ['vault.cluster.secrets.backend.credentials', '10-static'],
      'The manually entered role can be used without a successful list response'
    );
  });
});
