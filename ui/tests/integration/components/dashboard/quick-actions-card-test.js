/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, findAll, click, typeIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { selectChoose } from 'ember-power-select/test-support';
import sinon from 'sinon';
import { DASHBOARD } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';

module('Integration | Component | dashboard/quick-actions-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.api = this.owner.lookup('service:api').secrets;
    const router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(router, 'transitionTo');

    this.secretsEngines = [
      new SecretsEngineResource({
        accessor: 'kubernetes_f3400dee',
        path: 'kubernetes-test/',
        type: 'kubernetes',
      }),
      new SecretsEngineResource({ accessor: 'database_f3400dee', path: 'database-test/', type: 'database' }),
      new SecretsEngineResource({ accessor: 'pki_i1234dd', path: 'apki-test/', type: 'pki' }),
      new SecretsEngineResource({ accessor: 'secrets_j2350ii', path: 'secrets-test/', type: 'kv' }),
      new SecretsEngineResource({ accessor: 'nomad_123hh', path: 'nomad/', type: 'nomad' }),
      new SecretsEngineResource({ accessor: 'pki_f3400dee', path: 'pki-0-test/', type: 'pki' }),
      new SecretsEngineResource({
        accessor: 'pki_i1234dd',
        path: 'pki-1-test/',
        description: 'pki-1-path-description',
        type: 'pki',
      }),
      new SecretsEngineResource({
        accessor: 'secrets_j2350ii',
        path: 'kv-v2-test/',
        options: { version: 2 },
        type: 'kv',
      }),
      new SecretsEngineResource({
        accessor: 'secrets_j2350ii',
        path: 'kv-v1-test/',
        options: { version: 1 },
        type: 'kv',
      }),
    ];

    this.renderComponent = () =>
      render(hbs`<Dashboard::QuickActionsCard @secretsEngines={{this.secretsEngines}} />`);
  });

  hooks.afterEach(function () {
    this.transitionStub.restore();
  });

  test('it does not include kvv1 mounts', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.superSelect('secrets-engines'));

    findAll('.ember-power-select-option').forEach((o) => {
      assert.dom(o).doesNotHaveTextContaining('kv-v1-test');
    });
  });

  test('it should show quick action empty state if no engine is selected', async function (assert) {
    await this.renderComponent();
    assert.dom('.title').hasText('Quick actions');
    assert.dom(GENERAL.superSelect('secrets-engines')).exists({ count: 1 });
    assert.dom(DASHBOARD.emptyState('no-mount-selected')).exists({ count: 1 });
  });

  test('it selects a pki role and issues a leaf certificate', async function (assert) {
    const backend = 'pki-0-test';
    this.apiStub = sinon.stub(this.api, 'pkiListRoles').resolves({ keys: ['some-role'] });

    await this.renderComponent();
    await selectChoose(GENERAL.superSelect('secrets-engines'), backend);
    await selectChoose(GENERAL.superSelect('actions'), 'Issue certificate');

    assert.true(this.apiStub.calledWith(backend, 'true'), 'Request made to fetch options');
    assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
    assert.dom(DASHBOARD.subtitle('param')).hasText('Role to use');

    await selectChoose(GENERAL.superSelect('params'), 'some-role');
    assert.dom(DASHBOARD.actionButton('Issue leaf certificate')).exists({ count: 1 });

    await click(DASHBOARD.actionButton('Issue leaf certificate'));
    const [route, backendParam, roleParam] = this.transitionStub.lastCall.args;
    assert.strictEqual(
      route,
      'vault.cluster.secrets.backend.pki.roles.role.generate',
      'transition is called with expected route'
    );
    assert.strictEqual(backendParam, backend, 'transition has expected backend param');
    assert.strictEqual(roleParam, 'some-role', 'transition has expected role param');
  });

  test('it views a pki certificate', async function (assert) {
    const backend = 'pki-0-test';
    this.apiStub = sinon.stub(this.api, 'pkiListCerts').resolves({ keys: ['51:1c:39:42:ba'] });

    await this.renderComponent();
    await selectChoose(GENERAL.superSelect('secrets-engines'), backend);
    await selectChoose(GENERAL.superSelect('actions'), 'View certificate');

    assert.true(this.apiStub.calledWith(backend, 'true'), 'Request made to fetch options');
    assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
    assert.dom(DASHBOARD.subtitle('param')).hasText('Certificate serial number');
    assert.dom(DASHBOARD.actionButton('View certificate')).exists({ count: 1 });

    await selectChoose(GENERAL.superSelect('params'), '.ember-power-select-option', 0);
    await click(DASHBOARD.actionButton('View certificate'));

    const [route, backendParam, certParam] = this.transitionStub.lastCall.args;
    assert.strictEqual(
      route,
      'vault.cluster.secrets.backend.pki.certificates.certificate.details',
      'transition is called with expected route'
    );
    assert.strictEqual(backendParam, backend, 'transition has expected backend param');
    assert.strictEqual(certParam, '51:1c:39:42:ba', 'transition has expected cert param');
  });

  test('it views a pki issuer', async function (assert) {
    const backend = 'pki-0-test';
    this.apiStub = sinon
      .stub(this.api, 'pkiListIssuers')
      .resolves({ key_info: { '101709a1': { issuer_name: 'test' } }, keys: ['101709a1'] });

    await this.renderComponent();
    await selectChoose(GENERAL.superSelect('secrets-engines'), backend);
    await selectChoose(GENERAL.superSelect('actions'), 'View issuer');

    assert.true(this.apiStub.calledWith(backend, 'true'), 'Request made to fetch options');
    assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
    assert.dom(DASHBOARD.subtitle('param')).hasText('Issuer');
    assert.dom(DASHBOARD.actionButton('View issuer')).exists({ count: 1 });

    await selectChoose(GENERAL.superSelect('params'), '.ember-power-select-option', 0);
    await click(DASHBOARD.actionButton('View issuer'));

    const [route, backendParam, issuerParam] = this.transitionStub.lastCall.args;
    assert.strictEqual(
      route,
      'vault.cluster.secrets.backend.pki.issuers.issuer.details',
      'transition is called with expected route'
    );
    assert.strictEqual(backendParam, backend, 'transition has expected backend param');
    assert.strictEqual(issuerParam, 'test', 'transition has expected issuer param');
  });

  test('it selects a role and generates credentials for a database', async function (assert) {
    const backend = 'database-test';
    this.staticStub = sinon.stub(this.api, 'databaseListStaticRoles').resolves({ keys: ['static-role'] });
    this.dynamicStub = sinon.stub(this.api, 'databaseListRoles').resolves({ keys: ['dynamic-role'] });

    await this.renderComponent();

    await selectChoose(GENERAL.superSelect('secrets-engines'), backend);
    await selectChoose(GENERAL.superSelect('actions'), 'Generate credentials for database');

    assert.true(this.staticStub.calledWith(backend, 'true'), 'Request made to fetch static roles');
    assert.true(this.dynamicStub.calledWith(backend, 'true'), 'Request made to fetch dynamic roles');
    assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
    assert.dom(DASHBOARD.subtitle('param')).hasText('Role to use');
    assert.dom(DASHBOARD.actionButton('Generate credentials')).exists({ count: 1 });

    await click(GENERAL.superSelect('params'));
    assert.dom(GENERAL.searchSelect.option(0)).hasText('static-role', 'Static roles render in dropdown');
    assert.dom(GENERAL.searchSelect.option(1)).hasText('dynamic-role', 'Dynamic roles render in dropdown');

    await selectChoose(GENERAL.superSelect('params'), '.ember-power-select-option', 1);
    await click(DASHBOARD.actionButton('Generate credentials'));

    const [route, backendParam, issuerParam] = this.transitionStub.lastCall.args;
    assert.strictEqual(
      route,
      'vault.cluster.secrets.backend.credentials',
      'transition is called with expected route'
    );
    assert.strictEqual(backendParam, backend, 'transition has expected backend param');
    assert.strictEqual(issuerParam, 'dynamic-role', 'transition has expected role param');
  });

  test('it should show correct actions for kv', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.superSelect('secrets-engines'));
    assert.strictEqual(
      findAll('li.ember-power-select-option').length,
      5,
      'renders only kv v2, pki and db engines'
    );
    await selectChoose(GENERAL.superSelect('secrets-engines'), 'kv-v2-test');
    assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
    await selectChoose(GENERAL.superSelect('actions'), 'Find KV secrets');
    assert.dom(DASHBOARD.kvSearchSelect).exists('Shows option to search fo KVv2 secret');
    assert.dom(DASHBOARD.actionButton('Read secrets')).exists({ count: 1 });
  });

  test('it should render InputSearch when no items are returned for selected action', async function (assert) {
    const backend = 'pki-0-test';
    this.apiStub = sinon.stub(this.api, 'pkiListIssuers').rejects();

    await this.renderComponent();
    await selectChoose(GENERAL.superSelect('secrets-engines'), backend);
    await selectChoose(GENERAL.superSelect('actions'), 'Issue certificate');

    assert.dom(DASHBOARD.paramInputLabel).hasText('Role to use', 'Label renders in param input');
    assert
      .dom(DASHBOARD.paramInput)
      .hasAttribute('placeholder', 'Enter role name', 'Placeholder renders in param input');

    await typeIn(DASHBOARD.paramInput, 'my-role');
    await click(DASHBOARD.actionButton('Issue leaf certificate'));
    assert.true(
      this.transitionStub.calledWith(
        'vault.cluster.secrets.backend.pki.roles.role.generate',
        backend,
        'my-role'
      ),
      'Transitions with correct params'
    );
  });
});
