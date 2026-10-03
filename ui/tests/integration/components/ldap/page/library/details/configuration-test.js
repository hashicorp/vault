/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { duration } from 'core/helpers/format-duration';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupDetailsTest } from 'vault/tests/helpers/ldap/ldap-helpers';

module('Integration | Component | ldap | Page::Library::Details::Configuration', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);
  setupDetailsTest(hooks);

  hooks.beforeEach(function () {
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
      const value = this.model.library[key];
      let formattedValue = value;
      if (key === 'disable_check_in_enforcement') {
        formattedValue = value ? 'Disabled' : 'Enabled';
      } else if (['max_ttl', 'ttl'].includes(key)) {
        formattedValue = duration([value]);
      }
      const method = key === 'disable_check_in_enforcement' ? 'includesText' : 'hasText';

      assert.dom(GENERAL.infoRowLabel(label)).hasText(label, `${label} info row label renders`);
      assert.dom(GENERAL.infoRowValue(label))[method](formattedValue, `${label} info row value renders`);
    });

    assert
      .dom('[data-test-check-in-icon]')
      .hasClass('hds-icon-check-circle', 'Correct icon renders for enabled check in enforcement');
    assert
      .dom('[data-test-check-in-icon]')
      .hasClass('icon-true', 'Correct class renders for enabled check in enforcement');

    this.model.library.disable_check_in_enforcement = 'Disabled';
    await this.renderComponent();

    assert
      .dom('[data-test-check-in-icon]')
      .hasClass('hds-icon-x-square', 'Correct icon renders for disabled check in enforcement');
    assert
      .dom('[data-test-check-in-icon]')
      .hasClass('icon-false', 'Correct class renders for disabled check in enforcement');
  });
});
