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
import { duration } from 'core/helpers/format-duration';

module('Integration | Component | ldap | Page::Library::Details::Configuration', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');

    this.store.pushPayload('ldap/library', {
      modelName: 'ldap/library',
      backend: 'ldap-test',
      ...this.server.create('ldap-library', { name: 'test-library' }),
    });
    this.model = this.store.peekRecord('ldap/library', 'test-library');
    this.renderComponent = () => {
      return render(hbs`<Page::Library::Details::Configuration @model={{this.model}} />`, {
        owner: this.engine,
      });
    };
  });

  test('it should render configuration details', async function (assert) {
    await this.renderComponent();

    const fields = [
      { label: 'Library name', key: 'name' },
      { label: 'TTL', key: 'ttl' },
      { label: 'Max TTL', key: 'max_ttl' },
      { label: 'Check-in enforcement', key: 'disable_check_in_enforcement' },
    ];
    fields.forEach((field) => {
      const { label, key } = field;
      const value = label.includes('TTL') ? duration([this.model[key]]) : this.model[key];
      const method = key === 'disable_check_in_enforcement' ? 'includesText' : 'hasText';

      assert.dom(`[data-test-row-label="${label}"]`).hasText(label, `${label} info row label renders`);
      assert.dom(`[data-test-value-div="${label}"]`)[method](value, `${label} info row label renders`);
    });

    assert
      .dom('[data-test-check-in-icon]')
      .hasClass('flight-icon-check-circle', 'Correct icon renders for enabled check in enforcement');
    assert
      .dom('[data-test-check-in-icon]')
      .hasClass('icon-true', 'Correct class renders for enabled check in enforcement');

    this.model.disable_check_in_enforcement = 'Disabled';
    await this.renderComponent();

    assert
      .dom('[data-test-check-in-icon]')
      .hasClass('flight-icon-x-square', 'Correct icon renders for disabled check in enforcement');
    assert
      .dom('[data-test-check-in-icon]')
      .hasClass('icon-false', 'Correct class renders for disabled check in enforcement');
  });
});
