/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
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
    await render(hbs`<PkiPageHeader @backend={{this.model}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-header-title] span').hasClass('hs-icon', 'Correct icon renders in title');
    assert.dom('[data-test-header-title]').hasText(this.mount, 'Mount path renders in title');
  });

  test('it should render tabs', async function (assert) {
    await render(hbs`<PkiPageHeader @backend={{this.model}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-secret-list-tab="Overview"]').hasText('Overview', 'Overview tab renders');
    assert.dom('[data-test-secret-list-tab="Roles"]').hasText('Roles', 'Roles tab renders');
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
});
