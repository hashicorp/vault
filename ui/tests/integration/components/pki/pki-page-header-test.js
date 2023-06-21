/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | pki page header test', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'pki_f3400dee',
        path: 'pki-test/',
        type: 'pki',
      },
    });
    this.model = this.store.peekRecord('secret-engine', 'pki-test');
    this.mount = this.model.path.slice(0, -1);
  });

  test('it should render title', async function (assert) {
    await render(hbs`<PkiPageHeader @model={{this.model}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-header-title] svg').hasClass('flight-icon-pki', 'Correct icon renders in title');
    assert.dom('[data-test-header-title]').hasText(this.mount, 'Mount path renders in title');
  });

  test('it should render tabs', async function (assert) {
    await render(hbs`<PkiPageHeader @model={{this.model}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-secret-list-tab="overview"]').hasText('Overview', 'Overview tab renders');
    assert.dom('[data-test-secret-list-tab="roles"]').hasText('Roles', 'Roles tab renders');
    assert.dom('[data-test-secret-list-tab="Issuers"]').hasText('Issuers', 'Issuers tab renders');
    assert.dom('[data-test-secret-list-tab="Keys"]').hasText('Keys', 'Keys tab renders');
    assert
      .dom('[data-test-secret-list-tab="Certificates"]')
      .hasText('Certificates', 'Certificates tab renders');
    assert.dom('[data-test-secret-list-tab="Tidy"]').hasText('Tidy', 'Tidy tab renders');
    assert
      .dom('[data-test-secret-list-tab="Configuration"]')
      .hasText('Configuration', 'Configuration tab renders');
  });

  test('it should render filter for roles', async function (assert) {
    await render(
      hbs`<PkiPageHeader @model={{this.model}} @filterRoles={{true}} @rolesFilterValue="test" />`,
      { owner: this.engine }
    );
    assert.dom('[data-test-nav-input] input').hasValue('test', 'Filter renders with provided value');
  });

  test('it should yield block for toolbar actions', async function (assert) {
    await render(
      hbs`
      <PkiPageHeader @model={{this.model}}>
        <span data-test-yield>It yields!</span>
      </PkiPageHeader>
    `,
      { owner: this.engine }
    );

    assert
      .dom('.toolbar-actions [data-test-yield]')
      .hasText('It yields!', 'Block is yielded for toolbar actions');
  });
});
