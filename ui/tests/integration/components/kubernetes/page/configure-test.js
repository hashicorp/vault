/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, waitUntil, find, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { Response } from 'miragejs';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | kubernetes | Page::Configure', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.newModel = this.store.createRecord('kubernetes/config', { backend: 'kubernetes-new' });
    this.existingConfig = {
      kubernetes_host: 'https://192.168.99.100:8443',
      kubernetes_ca_cert: '-----BEGIN CERTIFICATE-----\n.....\n-----END CERTIFICATE-----',
      service_account_jwt: 'test-jwt',
      disable_local_ca_jwt: true,
    };
    this.store.pushPayload('kubernetes/config', {
      modelName: 'kubernetes/config',
      backend: 'kubernetes-edit',
      ...this.existingConfig,
    });
    this.editModel = this.store.peekRecord('kubernetes/config', 'kubernetes-edit');
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: 'kubernetes', route: 'overview' },
      { label: 'Configure' },
    ];
    this.expectedInferred = {
      disable_local_ca_jwt: false,
      kubernetes_ca_cert: null,
      kubernetes_host: null,
      service_account_jwt: null,
    };
    setRunOptions({
      rules: {
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it should display proper options when toggling radio cards', async function (assert) {
    await render(hbs`<Page::Configure @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
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

    let status = 404;
    this.server.get('/:path/check', () => {
      assert.ok(
        waitUntil(() => find('[data-test-config] button').disabled),
        'Button is disabled while request is in flight'
      );
      return new Response(status, {});
    });

    await render(hbs`<Page::Configure @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    await click('[data-test-config] button');
    assert
      .dom('[data-test-icon="x-square-fill"]')
      .hasClass('has-text-danger', 'Icon is displayed for error state with correct styling');
    const error =
      'Vault could not infer a configuration from your environment variables. Check your configuration file to edit or delete them, or configure manually.';
    assert.dom('[data-test-config] span').hasText(error, 'Error text is displayed');
    assert.dom('[data-test-config-save]').isDisabled('Save button is disabled in error state');

    status = 204;
    await click('[data-test-radio-card="manual"]');
    await click('[data-test-radio-card="local"]');
    await click('[data-test-config] button');
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

    this.server.post('/:path/config', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      assert.deepEqual(json, this.existingConfig, 'Values are passed to create endpoint');
      return new Response(204, {});
    });

    const stub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await render(hbs`<Page::Configure @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    await click('[data-test-radio-card="manual"]');
    await fillIn('[data-test-input="kubernetesHost"]', this.existingConfig.kubernetes_host);
    await fillIn('[data-test-input="serviceAccountJwt"]', this.existingConfig.service_account_jwt);
    await fillIn('[data-test-input="kubernetesCaCert"]', this.existingConfig.kubernetes_ca_cert);
    await click('[data-test-config-save]');
    assert.ok(
      stub.calledWith('vault.cluster.secrets.backend.kubernetes.configuration'),
      'Transitions to configuration route on save success'
    );
  });

  test('it should edit existing manual config', async function (assert) {
    assert.expect(6);

    const stub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await render(hbs`<Page::Configure @model={{this.editModel}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    assert.dom('[data-test-radio-card="manual"] input').isChecked('Manual config radio card is checked');
    assert
      .dom('[data-test-input="kubernetesHost"]')
      .hasValue(this.existingConfig.kubernetes_host, 'Host field is populated');
    assert
      .dom('[data-test-input="serviceAccountJwt"]')
      .hasValue(this.existingConfig.service_account_jwt, 'JWT field is populated');
    assert
      .dom('[data-test-input="kubernetesCaCert"]')
      .hasValue(this.existingConfig.kubernetes_ca_cert, 'Cert field is populated');

    await fillIn('[data-test-input="kubernetesHost"]', 'http://localhost:1212');
    await click('[data-test-config-cancel]');

    assert.ok(
      stub.calledWith('vault.cluster.secrets.backend.kubernetes.configuration'),
      'Transitions to configuration route when cancelling edit'
    );
    assert.strictEqual(
      this.editModel.kubernetesHost,
      this.existingConfig.kubernetes_host,
      'Model values are rolled back on cancel'
    );
  });

  test('it should display inferred success message when editing model using local values', async function (assert) {
    this.store.pushPayload('kubernetes/config', {
      modelName: 'kubernetes/config',
      backend: 'kubernetes-edit-2',
      disable_local_ca_jwt: false,
    });
    this.model = this.store.peekRecord('kubernetes/config', 'kubernetes-edit-2');

    await render(hbs`<Page::Configure @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

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

    this.server.post('/:path/config', () => {
      assert.ok(true, 'Save request made after confirmation');
      return new Response(204, {});
    });

    await render(
      hbs`
            <Page::Configure @model={{this.editModel}} @breadcrumbs={{this.breadcrumbs}} />
    `,
      { owner: this.engine }
    );
    await click('[data-test-config-save]');
    assert
      .dom('[data-test-edit-config-body]')
      .hasText(
        'Making changes to your configuration may affect how Vault will reach the Kubernetes API and authenticate with it. Are you sure?',
        'Confirm modal renders'
      );
    await click('[data-test-config-confirm]');
  });

  test('it should validate form and show errors', async function (assert) {
    await render(hbs`<Page::Configure @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    await click('[data-test-radio-card="manual"]');
    await click('[data-test-config-save]');

    assert
      .dom('[data-test-field-validation="kubernetesHost"] [data-test-inline-error-message]')
      .hasText('Kubernetes host is required', 'Error renders for required field');
    assert.dom('[data-test-alert]').hasText('There is an error with this form.', 'Alert renders');
  });

  test('it should save inferred config', async function (assert) {
    assert.expect(2);

    this.server.get('/:path/check', () => new Response(204, {}));
    this.server.post('/:path/config', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      assert.deepEqual(json, this.expectedInferred, 'Values are passed to create endpoint');
      return new Response(204, {});
    });

    const stub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await render(hbs`<Page::Configure @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    await click('[data-test-config] button');
    await click('[data-test-config-save]');

    assert.ok(
      stub.calledWith('vault.cluster.secrets.backend.kubernetes.configuration'),
      'Transitions to configuration route on save success'
    );
  });

  test('it should unset manual config values when saving local cluster option', async function (assert) {
    assert.expect(1);

    this.server.get('/:path/check', () => new Response(204, {}));
    this.server.post('/:path/config', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      assert.deepEqual(json, this.expectedInferred, 'Manual config values are unset in server payload');
      return new Response(204, {});
    });

    await render(hbs`<Page::Configure @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    await click('[data-test-radio-card="manual"]');
    await fillIn('[data-test-input="kubernetesHost"]', this.existingConfig.kubernetes_host);
    await fillIn('[data-test-input="serviceAccountJwt"]', this.existingConfig.service_account_jwt);
    await fillIn('[data-test-input="kubernetesCaCert"]', this.existingConfig.kubernetes_ca_cert);

    await click('[data-test-radio-card="local"]');
    await click('[data-test-config] button');
    await click('[data-test-config-save]');
  });
});
