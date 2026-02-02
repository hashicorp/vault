/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { createSecretsEngine, generateBreadcrumbs } from 'vault/tests/helpers/ldap/ldap-helpers';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { LDAP_SELECTORS } from 'vault/tests/helpers/ldap/ldap-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | ldap | Page::Libraries', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'ldap-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);
    this.secretsEngine = createSecretsEngine();
    this.breadcrumbs = generateBreadcrumbs(this.backend);
    this.libraries = ['foo', 'bar', 'foo/'].map((name) =>
      this.server.create('ldap-library', { name, completeLibraryName: name })
    );
    this.capabilities = this.libraries.reduce((capabilities, { name }) => {
      const path = this.owner
        .lookup('service:capabilities')
        .pathFor('ldapLibrary', { backend: this.backend, name });
      capabilities[path] = { canRead: true, canUpdate: true, canDelete: true };
      return capabilities;
    }, {});
    this.promptConfig = false;

    this.renderComponent = () =>
      render(
        hbs`<Page::Libraries
          @libraries={{this.libraries}}
          @capabilities={{this.capabilities}}
          @promptConfig={{this.promptConfig}}
          @secretsEngine={{this.secretsEngine}}
          @breadcrumbs={{this.breadcrumbs}}
        />`,
        { owner: this.engine }
      );

    setRunOptions({
      rules: {
        list: { enabled: false },
      },
    });
  });

  test('it should render tab page header', async function (assert) {
    this.promptConfig = true;

    await this.renderComponent();

    assert.dom(GENERAL.icon('folder-users')).hasClass('hds-icon-folder-users', 'LDAP icon renders in title');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('ldap-test', 'Mount path renders in title');
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
      .dom(GENERAL.emptyStateTitle)
      .hasText('There are no libraries matching "baz"', 'Filter message renders');

    await fillIn('[data-test-filter-input]', 'foo');
    assert.dom('[data-test-list-item-content]').exists({ count: 2 }, 'List is filtered with correct results');

    await fillIn('[data-test-filter-input]', '');
    assert
      .dom('[data-test-list-item-content]')
      .exists({ count: 3 }, 'All libraries are displayed when filter is cleared');
  });
});
