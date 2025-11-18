/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { duration } from 'core/helpers/format-duration';
import { createSecretsEngine, generateBreadcrumbs } from 'vault/tests/helpers/ldap/ldap-helpers';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import engineDisplayData from 'vault/helpers/engines-display-data';

const selectors = {
  rotateAction: '[data-test-toolbar-rotate-action]',
  configAction: '[data-test-toolbar-config-action]',
  configCta: '[data-test-config-cta]',
  mountConfig: '[data-test-mount-config]',
  pageError: '[data-test-page-error]',
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

    this.model = {
      backendModel: this.backend,
      promptConfig: true,
      configModel: this.config,
      configError: null,
      engineDisplayData: engineDisplayData(this.backend.type),
    };

    this.renderComponent = () => {
      return render(
        hbs`<Page::Configuration
          @model={{this.model}}
          @breadcrumbs={{this.breadcrumbs}}
        />`,
        {
          owner: this.engine,
        }
      );
    };
    setRunOptions({
      rules: {
        // TODO: fix ConfirmAction rendered in toolbar not a list item
        list: { enabled: false },
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  test('it should render tab page header, config cta and mount config', async function (assert) {
    this.model.configModel = null;

    await this.renderComponent();

    assert.dom(GENERAL.icon('folder-users')).hasClass('hds-icon-folder-users', 'LDAP icon renders in title');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('ldap-test configuration', 'Mount path renders in title');
    assert
      .dom(selectors.rotateAction)
      .doesNotExist('Rotate root action is hidden when engine is not configured');
    assert.dom(selectors.configAction).hasText('Configure LDAP', 'Toolbar action has correct text');
  });

  test('it should render config fetch error', async function (assert) {
    this.model.configModel = null;
    this.model.configError = { httpStatus: 403, message: 'Permission denied' };

    await this.renderComponent();

    assert.dom(GENERAL.pageError.error).exists('Config fetch error is rendered');
  });

  test('it should render display fields', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.infoRowValue('Administrator Distinguished Name')).hasText(this.config.binddn);
    assert.dom(GENERAL.infoRowValue('URL')).hasText(this.config.url);
    assert.dom(GENERAL.infoRowValue('Schema')).hasText(this.config.schema);
    assert.dom(GENERAL.infoRowValue('Password Policy')).hasText(this.config.password_policy);
    assert.dom(GENERAL.infoRowValue('Userdn')).hasText(this.config.userdn);
    assert.dom(GENERAL.infoRowValue('Userattr')).hasText(this.config.userattr);
    assert
      .dom(GENERAL.infoRowValue('Connection Timeout'))
      .hasText(duration([this.config.connection_timeout]));
    assert.dom(GENERAL.infoRowValue('Request Timeout')).hasText(duration([this.config.request_timeout]));
    assert.dom(GENERAL.infoRowValue('CA Certificate')).hasText(this.config.certificate);
    assert.dom(GENERAL.infoRowValue('Start TLS')).includesText('No');
    assert.dom(GENERAL.infoRowValue('Insecure TLS')).includesText('No');
    assert.dom(GENERAL.infoRowValue('Client TLS Certificate')).hasText(this.config.client_tls_cert);
    assert.dom(GENERAL.infoRowValue('Client TLS Key')).hasText(this.config.client_tls_key);
  });

  test('it should rotate root password', async function (assert) {
    assert.expect(1);

    this.server.post(`/${this.config.backend}/rotate-root`, () => {
      assert.ok(true, 'Request made to rotate root password');
    });

    await this.renderComponent();
    await click(selectors.rotateAction);
    await click(GENERAL.confirmButton);
  });
});
