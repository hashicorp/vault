/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module(
  'Integration | Component | pki | external-pki | ExternalPki::Page::RolesRoleDetails',
  function (hooks) {
    setupRenderingTest(hooks);
    setupEngine(hooks, 'pki');

    hooks.beforeEach(function () {
      this.model = {
        engine: { id: 'pki-external-ca' },
        role: {},
      };
      this.renderComponent = () =>
        render(hbs`<ExternalPki::Page::RolesRoleDetails @model={{this.model}} />`, { owner: this.engine });
    });

    test('it renders role details', async function (assert) {
      this.model.role = {
        name: 'test-role',
        acme_account_name: 'production-account',
        dns_provider_name: 'aws-route53-prod',
        allowed_domains: ['example.com', '*.example.com'],
        allow_subdomains: true,
      };
      await this.renderComponent();
      assert.dom(GENERAL.infoRowValue('ACME account name')).hasText('production-account');
      assert.dom(GENERAL.infoRowValue('DNS provider name')).hasText('aws-route53-prod');
      assert.dom(GENERAL.infoRowValue('Allowed domains')).hasText('example.com,*.example.com');
      assert.dom(GENERAL.infoRowValue('Allow subdomains')).hasText('Yes');
    });
  }
);
