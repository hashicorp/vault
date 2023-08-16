/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { createSecretsEngine, generateBreadcrumbs } from 'vault/tests/helpers/ldap';
import sinon from 'sinon';

module('Integration | Component | ldap | Page::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');

    this.backendModel = createSecretsEngine(this.store);
    this.breadcrumbs = generateBreadcrumbs(this.backendModel.id);

    const pushPayload = (type) => {
      this.store.pushPayload(`ldap/${type}`, {
        modelName: `ldap/${type}`,
        backend: 'ldap-test',
        ...this.server.create(`ldap-${type}`),
      });
    };

    ['role', 'library'].forEach((type) => {
      pushPayload(type);
      if (type === 'role') {
        pushPayload(type);
      }
      const key = type === 'role' ? 'roles' : 'libraries';
      this[key] = this.store.peekAll(`ldap/${type}`);
    });

    this.renderComponent = () => {
      return render(
        hbs`<Page::Overview
          @promptConfig={{this.promptConfig}}
          @backendModel={{this.backendModel}}
          @roles={{this.roles}}
          @libraries={{this.libraries}}
          @librariesStatus={{(array)}}
          @breadcrumbs={{this.breadcrumbs}}
        />`,
        {
          owner: this.engine,
        }
      );
    };
  });

  test('it should render tab page header and config cta', async function (assert) {
    this.promptConfig = true;

    await this.renderComponent();

    assert.dom('.title svg').hasClass('flight-icon-folder-users', 'LDAP icon renders in title');
    assert.dom('.title').hasText('ldap-test', 'Mount path renders in title');
    assert.dom('[data-test-toolbar-action="config"]').hasText('Configure LDAP', 'Toolbar action renders');
    assert.dom('[data-test-config-cta]').exists('Config cta renders');
  });

  test('it should render overview cards', async function (assert) {
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await this.renderComponent();

    assert.dom('[data-test-roles-count]').hasText('2', 'Roles card renders with correct count');
    assert.dom('[data-test-libraries-count]').hasText('1', 'Libraries card renders with correct count');
    assert
      .dom('[data-test-overview-card-container="Accounts checked-out"]')
      .exists('Accounts checked-out card renders');

    await click('[data-test-component="search-select"] .ember-power-select-trigger');
    await click('.ember-power-select-option');
    await click('[data-test-generate-credential-button]');

    const didTransition = transitionStub.calledWith(
      'vault.cluster.secrets.backend.ldap.roles.role.credentials',
      this.roles[0].type,
      this.roles[0].name
    );
    assert.true(didTransition, 'Transitions to credentials route when generating credentials');
  });
});
