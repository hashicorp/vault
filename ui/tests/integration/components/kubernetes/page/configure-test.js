/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import KubernetesConfigForm from 'vault/forms/secrets/kubernetes/config';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | kubernetes | Page::Configure', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kubernetes-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.createForm = new KubernetesConfigForm({ disable_local_ca_jwt: false }, { isNew: true });
    this.existingConfig = {
      kubernetes_host: 'https://192.168.99.100:8443',
      kubernetes_ca_cert: '-----BEGIN CERTIFICATE-----\n.....\n-----END CERTIFICATE-----',
      service_account_jwt: 'test-jwt',
      disable_local_ca_jwt: true,
    };
    this.editForm = new KubernetesConfigForm(this.existingConfig);
    this.form = this.createForm;
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: 'kubernetes', route: 'overview' },
      { label: 'Configure' },
    ];
    this.expectedInferred = {
      disable_local_ca_jwt: false,
    };

    const { secrets } = this.owner.lookup('service:api');
    this.checkStub = sinon.stub(secrets, 'kubernetesCheckConfiguration').resolves();
    this.configStub = sinon.stub(secrets, 'kubernetesConfigure').resolves();

    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    setRunOptions({
      rules: {
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });
    this.renderComponent = () =>
      render(hbs`<Page::Configure @form={{this.form}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
  });

  test('it should display proper options when toggling radio cards', async function (assert) {
    await this.renderComponent();

    assert
      .dom('[data-test-radio-card="local"] input')
      .isChecked('Local cluster radio card is checked by default');
    assert
      .dom('[data-test-config] p')
      .hasText(
        'Configuration values can be inferred from the pod and your local environment variables.',
        'Inferred text is displayed'
      );
    assert.dom('[data-test-config] button').hasText('Get config values', 'Get config button renders');
    assert
      .dom('[data-test-config-save]')
      .isDisabled('Save button is disabled when config values have not been inferred');
    assert.dom('[data-test-config-cancel]').hasText('Back', 'Back button renders');

    await click('[data-test-radio-card="manual"]');
    assert.dom('[data-test-field]').exists({ count: 3 }, 'Form fields render');
    assert.dom('[data-test-config-save]').isNotDisabled('Save button is enabled');
    assert.dom('[data-test-config-cancel]').hasText('Back', 'Back button renders');
  });

  test('it should check for inferred config variables', async function (assert) {
    assert.expect(8);

    this.checkStub.rejects(getErrorResponse());
    await this.renderComponent();

    await click('[data-test-config] button');
    assert.true(this.checkStub.calledWith(this.backend), 'Check config request is made');
    assert
      .dom('[data-test-icon="x-square-fill"]')
      .hasClass('has-text-danger', 'Icon is displayed for error state with correct styling');
    const error =
      'Vault could not infer a configuration from your environment variables. Check your configuration file to edit or delete them, or configure manually.';
    assert.dom('[data-test-config] span').hasText(error, 'Error text is displayed');
    assert.dom('[data-test-config-save]').isDisabled('Save button is disabled in error state');

    this.checkStub.resolves();
    await click('[data-test-radio-card="manual"]');
    await click('[data-test-radio-card="local"]');
    await click('[data-test-config] button');
    assert.true(this.checkStub.calledWith(this.backend), 'Check config request is made');
    assert
      .dom('[data-test-icon="check-circle-fill"]')
      .hasClass('has-text-success', 'Icon is displayed for success state with correct styling');
    assert
      .dom('[data-test-config] span')
      .hasText('Configuration values were inferred successfully.', 'Success text is displayed');
    assert.dom('[data-test-config-save]').isNotDisabled('Save button is enabled in success state');
  });

  test('it should create new manual config', async function (assert) {
    assert.expect(2);

    await this.renderComponent();

    await click('[data-test-radio-card="manual"]');
    await fillIn('[data-test-input="kubernetes_host"]', this.existingConfig.kubernetes_host);
    await fillIn('[data-test-input="service_account_jwt"]', this.existingConfig.service_account_jwt);
    await fillIn('[data-test-input="kubernetes_ca_cert"]', this.existingConfig.kubernetes_ca_cert);
    await click('[data-test-config-save]');
    assert.true(
      this.configStub.calledWith(this.backend, this.existingConfig),
      'Create config request is made'
    );
    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.kubernetes.configuration'),
      'Transitions to configuration route on save success'
    );
  });

  test('it should render existing manual config data in form', async function (assert) {
    assert.expect(5);

    this.form = this.editForm;
    await this.renderComponent();

    assert.dom('[data-test-radio-card="manual"] input').isChecked('Manual config radio card is checked');
    assert
      .dom('[data-test-input="kubernetes_host"]')
      .hasValue(this.existingConfig.kubernetes_host, 'Host field is populated');
    assert
      .dom('[data-test-input="service_account_jwt"]')
      .hasValue(this.existingConfig.service_account_jwt, 'JWT field is populated');
    assert
      .dom('[data-test-input="kubernetes_ca_cert"]')
      .hasValue(this.existingConfig.kubernetes_ca_cert, 'Cert field is populated');

    await click('[data-test-config-cancel]');

    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.kubernetes.configuration'),
      'Transitions to configuration route when cancelling edit'
    );
  });

  test('it should display inferred success message when editing model using local values', async function (assert) {
    this.form = this.editForm;
    this.form.data.disable_local_ca_jwt = false;

    await this.renderComponent();

    assert.dom('[data-test-radio-card="local"] input').isChecked('Local cluster radio card is checked');
    assert
      .dom('[data-test-icon="check-circle-fill"]')
      .hasClass('has-text-success', 'Icon is displayed for success state with correct styling');
    assert
      .dom('[data-test-config] span')
      .hasText('Configuration values were inferred successfully.', 'Success text is displayed');
  });

  test('it should show confirmation modal when saving edits', async function (assert) {
    assert.expect(2);

    this.form = this.editForm;
    await this.renderComponent();

    await click('[data-test-config-save]');
    assert
      .dom('[data-test-edit-config-body]')
      .hasText(
        'Making changes to your configuration may affect how Vault will reach the Kubernetes API and authenticate with it. Are you sure?',
        'Confirm modal renders'
      );
    await click('[data-test-config-confirm]');
    assert.true(
      this.configStub.calledWith(this.backend, this.existingConfig),
      'Config is saved after confirming'
    );
  });

  test('it should validate form and show errors', async function (assert) {
    await this.renderComponent();

    await click('[data-test-radio-card="manual"]');
    await click('[data-test-config-save]');

    assert
      .dom(GENERAL.validationErrorByAttr('kubernetes_host'))
      .hasText('Kubernetes host is required', 'Error renders for required field');
    assert.dom('[data-test-alert]').hasText('There is an error with this form.', 'Alert renders');
  });

  test('it should save inferred config', async function (assert) {
    assert.expect(2);

    await this.renderComponent();

    await click('[data-test-config] button');
    await click('[data-test-config-save]');

    assert.true(
      this.configStub.calledWith(this.backend, this.expectedInferred),
      'Request made to save inferred config values'
    );
    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.kubernetes.configuration'),
      'Transitions to configuration route on save success'
    );
  });

  test('it should unset manual config values when saving local cluster option', async function (assert) {
    assert.expect(1);

    await this.renderComponent();

    await click('[data-test-radio-card="manual"]');
    await fillIn('[data-test-input="kubernetes_host"]', this.existingConfig.kubernetes_host);
    await fillIn('[data-test-input="service_account_jwt"]', this.existingConfig.service_account_jwt);
    await fillIn('[data-test-input="kubernetes_ca_cert"]', this.existingConfig.kubernetes_ca_cert);

    await click('[data-test-radio-card="local"]');
    await click('[data-test-config] button');
    await click('[data-test-config-save]');

    assert.true(
      this.configStub.calledWith(this.backend, this.expectedInferred),
      'Manual config values are unset in server payload'
    );
  });
});
