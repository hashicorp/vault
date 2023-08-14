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
import { duration } from 'core/helpers/format-duration';
import { createSecretsEngine, generateBreadcrumbs } from 'vault/tests/helpers/ldap';

const selectors = {
  rotateAction: '[data-test-toolbar-rotate-action] button',
  confirmRotate: '[data-test-confirm-button]',
  configAction: '[data-test-toolbar-config-action]',
  configCta: '[data-test-config-cta]',
  mountConfig: '[data-test-mount-config]',
  pageError: '[data-test-page-error]',
  fieldValue: (label) => `[data-test-value-div="${label}"]`,
};

module('Integration | Component | ldap | Page::Configuration', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');

    this.backend = createSecretsEngine(this.store);
    this.breadcrumbs = generateBreadcrumbs(this.backend.id);

    this.store.pushPayload('ldap/config', {
      modelName: 'ldap/config',
      backend: 'ldap-test',
      ...this.server.create('ldap-config'),
    });
    this.config = this.store.peekRecord('ldap/config', 'ldap-test');

    this.renderComponent = () => {
      return render(
        hbs`<Page::Configuration
          @backendModel={{this.backend}}
          @configModel={{this.config}}
          @configError={{this.error}}
          @breadcrumbs={{this.breadcrumbs}}
        />`,
        {
          owner: this.engine,
        }
      );
    };
  });

  test('it should render tab page header, config cta and mount config', async function (assert) {
    this.config = null;

    await this.renderComponent();

    assert.dom('.title svg').hasClass('flight-icon-folder-users', 'LDAP icon renders in title');
    assert.dom('.title').hasText('ldap-test', 'Mount path renders in title');
    assert
      .dom(selectors.rotateAction)
      .doesNotExist('Rotate root action is hidden when engine is not configured');
    assert.dom(selectors.configAction).hasText('Configure LDAP', 'Toolbar action has correct text');
    assert.dom(selectors.configCta).exists('Config cta renders');
    assert.dom(selectors.mountConfig).exists('Mount config renders');
  });

  test('it should render config fetch error', async function (assert) {
    this.config = null;
    this.error = { httpStatus: 403, message: 'Permission denied' };

    await this.renderComponent();

    assert.dom(selectors.pageError).exists('Config fetch error is rendered');
  });

  test('it should render display fields', async function (assert) {
    await this.renderComponent();

    assert.dom(selectors.fieldValue('Administrator Distinguished Name')).hasText(this.config.binddn);
    assert.dom(selectors.fieldValue('URL')).hasText(this.config.url);
    assert.dom(selectors.fieldValue('Schema')).hasText(this.config.schema);
    assert.dom(selectors.fieldValue('Password Policy')).hasText(this.config.password_policy);
    assert.dom(selectors.fieldValue('Userdn')).hasText(this.config.userdn);
    assert.dom(selectors.fieldValue('Userattr')).hasText(this.config.userattr);
    assert
      .dom(selectors.fieldValue('Connection Timeout'))
      .hasText(duration([this.config.connection_timeout]));
    assert.dom(selectors.fieldValue('Request Timeout')).hasText(duration([this.config.request_timeout]));
    assert.dom(selectors.fieldValue('CA Certificate')).hasText(this.config.certificate);
    assert.dom(selectors.fieldValue('Start TLS')).includesText('No');
    assert.dom(selectors.fieldValue('Insecure TLS')).includesText('No');
    assert.dom(selectors.fieldValue('Client TLS Certificate')).hasText(this.config.client_tls_cert);
    assert.dom(selectors.fieldValue('Client TLS Key')).hasText(this.config.client_tls_key);
  });

  test('it should rotate root password', async function (assert) {
    assert.expect(1);

    this.server.post(`/${this.config.backend}/rotate-root`, () => {
      assert.ok(true, 'Request made to rotate root password');
    });

    await this.renderComponent();
    await click(selectors.rotateAction);
    await click(selectors.confirmRotate);
  });
});
