/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { createSecretsEngine, generateBreadcrumbs } from 'vault/tests/helpers/ldap/ldap-helpers';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { LDAP_SELECTORS } from 'vault/tests/helpers/ldap/ldap-selectors';

module('Integration | Component | ldap | Page::Libraries', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    this.store = this.owner.lookup('service:store');
    this.backend = createSecretsEngine(this.store);
    this.breadcrumbs = generateBreadcrumbs(this.backend.id);

    for (const name of ['foo', 'bar', 'foo/']) {
      this.store.pushPayload('ldap/library', {
        modelName: 'ldap/library',
        backend: 'ldap-test',
        ...this.server.create('ldap-library', { name }),
      });
    }
    this.libraries = this.store.peekAll('ldap/library');
    this.promptConfig = false;

    this.renderComponent = () => {
      return render(
        hbs`<Page::Libraries
          @promptConfig={{this.promptConfig}}
          @backendModel={{this.backend}}
          @libraries={{this.libraries}}
          @breadcrumbs={{this.breadcrumbs}}
        />`,
        { owner: this.engine }
      );
    };
    setRunOptions({
      rules: {
        list: { enabled: false },
      },
    });
  });

  test('it should render tab page header and config cta', async function (assert) {
    this.promptConfig = true;

    await this.renderComponent();

    assert.dom('.title svg').hasClass('hds-icon-folder-users', 'LDAP icon renders in title');
    assert.dom('.title').hasText('ldap-test', 'Mount path renders in title');
    assert
      .dom('[data-test-toolbar-action="config"]')
      .hasText('Configure LDAP', 'Correct toolbar action renders');
    assert.dom('[data-test-config-cta]').exists('Config cta renders');
  });

  test('it should render create libraries cta', async function (assert) {
    this.libraries = null;

    await this.renderComponent();

    assert
      .dom('[data-test-toolbar-action="library"]')
      .hasText('Create library', 'Toolbar action has correct text');
    assert
      .dom('[data-test-toolbar-action="library"] svg')
      .hasClass('hds-icon-plus', 'Toolbar action has correct icon');
    assert
      .dom('[data-test-filter-input]')
      .doesNotExist('Libraries filter input is hidden when libraries have not been created');
    assert.dom('[data-test-empty-state-title]').hasText('No libraries created yet', 'Title renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'Use libraries to manage a set of highly privileged accounts that can be shared among a team.',
        'Message renders'
      );
  });

  test('it should render libraries list', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-list-item-content] svg').hasClass('hds-icon-folder', 'List item icon renders');
    assert.dom('[data-test-library="foo"]').hasText('foo', 'List item name renders');

    await click(LDAP_SELECTORS.libraryMenu('foo'));
    assert.dom('[data-test-subdirectory]').doesNotExist();
    assert.dom('[data-test-edit]').hasText('Edit', 'Edit link renders in menu');
    assert.dom('[data-test-details]').hasText('Details', 'Details link renders in menu');
    assert.dom('[data-test-delete]').hasText('Delete', 'Details link renders in menu');

    await click(LDAP_SELECTORS.libraryMenu('foo/'));
    assert.dom('[data-test-subdirectory]').hasText('Content', 'Content link renders in menu');
    assert.dom('[data-test-edit]').doesNotExist();
    assert.dom('[data-test-details]').doesNotExist();
    assert.dom('[data-test-delete]').doesNotExist();
  });

  test('it should filter libraries', async function (assert) {
    await this.renderComponent();

    await fillIn('[data-test-filter-input]', 'baz');
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('There are no libraries matching "baz"', 'Filter message renders');

    await fillIn('[data-test-filter-input]', 'foo');
    assert.dom('[data-test-list-item-content]').exists({ count: 2 }, 'List is filtered with correct results');

    await fillIn('[data-test-filter-input]', '');
    assert
      .dom('[data-test-list-item-content]')
      .exists({ count: 3 }, 'All libraries are displayed when filter is cleared');
  });
});
