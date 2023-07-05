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
import { createSecretsEngine, generateBreadcrumbs } from 'vault/tests/helpers/ldap';

const selectors = {
  configAction: '[data-test-toolbar-action="config"]',
  configCta: '[data-test-config-cta]',
};

module('Integration | Component | ldap | Page::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');

    this.backend = createSecretsEngine(this.store);
    this.breadcrumbs = generateBreadcrumbs(this.backend.id);

    this.renderComponent = () => {
      return render(
        hbs`<Page::Overview
          @promptConfig={{this.promptConfig}}
          @backendModel={{this.backend}}
          @roles={{this.roles}}
          @libraries={{this.libraries}}
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
    assert.dom(selectors.configAction).hasText('Configure LDAP', 'Correct toolbar action renders');
    assert.dom(selectors.configCta).exists('Config cta renders');
  });

  // TODO:JLR add test to check that card components render once created and that Create role action renders in toolbar
});
