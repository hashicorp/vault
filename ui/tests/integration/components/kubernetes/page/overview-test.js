/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import hbs from 'htmlbars-inline-precompile';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { KUBERNETES_OVERVIEW } from 'vault/tests/helpers/kubernetes/kubernetes-selectors';

module('Integration | Component | kubernetes | Page::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kubernetes_f3400dee',
        path: 'kubernetes-test/',
        type: 'kubernetes',
      },
    });
    this.store.pushPayload('kubernetes/role', {
      modelName: 'kubernetes/role',
      backend: 'kubernetes-test',
      ...this.server.create('kubernetes-role'),
    });
    this.store.pushPayload('kubernetes/role', {
      modelName: 'kubernetes/role',
      backend: 'kubernetes-test',
      ...this.server.create('kubernetes-role'),
    });
    this.backend = this.store.peekRecord('secret-engine', 'kubernetes-test');
    this.roles = this.store.peekAll('kubernetes/role');
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend.id },
    ];
    this.promptConfig = false;
    this.renderComponent = () => {
      return render(
        hbs`<Page::Overview @promptConfig={{this.promptConfig}} @backend={{this.backend}} @roles={{this.roles}} @breadcrumbs={{this.breadcrumbs}} />`,
        { owner: this.engine }
      );
    };
    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
      },
    });
  });

  test('it should display role card', async function (assert) {
    await this.renderComponent();
    assert.dom(KUBERNETES_OVERVIEW.rolesCardTitle).hasText('Roles');
    assert
      .dom(KUBERNETES_OVERVIEW.rolesCardSubTitle)
      .hasText('The number of Vault roles being used to generate Kubernetes credentials.');
    assert.dom(KUBERNETES_OVERVIEW.rolesCardLink).hasText('View Roles');

    this.roles = [];

    await this.renderComponent();
    assert.dom(KUBERNETES_OVERVIEW.rolesCardLink).hasText('Create Role');
  });

  test('it should display correct number of roles in role card', async function (assert) {
    await this.renderComponent();
    assert.dom(KUBERNETES_OVERVIEW.rolesCardNumRoles).hasText('2');

    this.roles = [];

    await this.renderComponent();

    assert.dom(KUBERNETES_OVERVIEW.rolesCardNumRoles).hasText('None');
  });

  test('it should display generate credentials card', async function (assert) {
    await this.renderComponent();

    assert.dom(KUBERNETES_OVERVIEW.generateCredentialsCardTitle).hasText('Generate credentials');
    assert
      .dom(KUBERNETES_OVERVIEW.generateCredentialsCardSubTitle)
      .hasText('Quickly generate credentials by typing the role name.');
  });

  test('it should show options for SearchSelect', async function (assert) {
    await this.renderComponent();
    await clickTrigger();
    assert.strictEqual(this.element.querySelectorAll('.ember-power-select-option').length, 2);
    await typeInSearch('role-0');
    assert.strictEqual(this.element.querySelectorAll('.ember-power-select-option').length, 1);
    assert.dom(KUBERNETES_OVERVIEW.generateCredentialsCardButton).isDisabled();
    await selectChoose('', '.ember-power-select-option', 2);
    assert.dom(KUBERNETES_OVERVIEW.generateCredentialsCardButton).isNotDisabled();
  });

  test('it should show ConfigCta when no config is set up', async function (assert) {
    this.promptConfig = true;

    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('Kubernetes not configured');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Get started by establishing the URL of the Kubernetes API to connect to, along with some additional options.'
      );
    assert.dom(GENERAL.emptyStateActions).hasText('Configure Kubernetes');
  });
});
