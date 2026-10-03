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
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import sinon from 'sinon';

module('Integration | Component | ldap | Page::Configuration', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.secretsEngine = createSecretsEngine();
    this.owner.lookup('service:secret-mount-path').update(this.secretsEngine.path);
    this.breadcrumbs = generateBreadcrumbs(this.secretsEngine.id);
    this.config = this.server.create('ldap-config');
    this.model = {
      secretsEngine: this.secretsEngine,
      promptConfig: true,
      config: this.config,
      configError: null,
      engineDisplayData: engineDisplayData(this.secretsEngine.type),
    };

    this.renderComponent = () =>
      render(hbs`<Page::Configuration @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });

    setRunOptions({
      rules: {
        // TODO: fix ConfirmAction rendered in toolbar not a list item
        list: { enabled: false },
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  test('it should render tab page header', async function (assert) {
    this.model.config = null;

    await this.renderComponent();

    assert.dom(GENERAL.icon('folder-users')).hasClass('hds-icon-folder-users', 'LDAP icon renders in title');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('ldap-test configuration', 'Mount path renders in title');
    assert
      .dom(GENERAL.confirmTrigger)
      .doesNotExist('Rotate root action is hidden when engine is not configured');
    assert.dom(SES.configure).doesNotExist('"Edit configuration" is hidden when not configured');
  });

  test('it should render config fetch error', async function (assert) {
    this.model.config = null;
    this.model.configError = { status: 403, message: 'Permission denied' };

    await this.renderComponent();

    assert.dom(GENERAL.pageError.error).exists('Config fetch error is rendered');
  });

  test('it should render display fields', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.infoRowValue('Administrator distinguished name')).hasText(this.config.binddn);
    assert.dom(GENERAL.infoRowValue('URL')).hasText(this.config.url);
    assert.dom(GENERAL.infoRowValue('Schema')).hasText(this.config.schema);
    assert.dom(GENERAL.infoRowValue('Password policy')).hasText(this.config.password_policy);
    assert.dom(GENERAL.infoRowValue('Userdn')).hasText(this.config.userdn);
    assert.dom(GENERAL.infoRowValue('Userattr')).hasText(this.config.userattr);
    assert
      .dom(GENERAL.infoRowValue('Connection timeout'))
      .hasText(duration([this.config.connection_timeout]));
    assert.dom(GENERAL.infoRowValue('Request timeout')).hasText(duration([this.config.request_timeout]));
    assert.dom(GENERAL.infoRowValue('CA certificate')).hasText(this.config.certificate);
    assert.dom(GENERAL.infoRowValue('Start TLS')).includesText('No');
    assert.dom(GENERAL.infoRowValue('Insecure TLS')).includesText('No');
    assert.dom(GENERAL.infoRowValue('Client TLS certificate')).hasText(this.config.client_tls_cert);
    assert.dom(GENERAL.infoRowValue('Client TLS key')).hasText(this.config.client_tls_key);
  });

  test('it should rotate root password', async function (assert) {
    assert.expect(1);

    const rotateStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'ldapRotateRootCredentials')
      .resolves();

    await this.renderComponent();
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.true(rotateStub.calledWith(this.secretsEngine.path), 'rotate root called with correct mount path');
  });
});
