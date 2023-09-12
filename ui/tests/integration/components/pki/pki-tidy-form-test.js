/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/page/pki-tidy-form';

module('Integration | Component | pki tidy form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.14.1+ent';
    this.server.post('/sys/capabilities-self', () => {});
    this.onSave = () => {};
    this.onCancel = () => {};
    this.manualTidy = this.store.createRecord('pki/tidy', { backend: 'pki-manual-tidy' });
    this.store.pushPayload('pki/tidy', {
      modelName: 'pki/tidy',
      id: 'pki-auto-tidy',
    });
    this.autoTidy = this.store.peekRecord('pki/tidy', 'pki-auto-tidy');
  });

  test('it hides or shows fields depending on auto-tidy toggle', async function (assert) {
    assert.expect(37);
    this.version.version = '1.14.1+ent';
    const sectionHeaders = [
      'Universal operations',
      'ACME operations',
      'Issuer operations',
      'Cross-cluster operations',
    ];

    await render(
      hbs`
      <PkiTidyForm
      @tidy={{this.autoTidy}}
      @tidyType="auto"
      @onSave={{this.onSave}}
      @onCancel={{this.onCancel}}
    />
    `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.toggleInput('intervalDuration')).isNotChecked('Automatic tidy is disabled');
    assert.dom(`[data-test-ttl-form-label="Automatic tidy disabled"]`).exists('renders disabled label text');

    this.autoTidy.eachAttribute((attr) => {
      if (attr === 'enabled' || attr === 'intervalDuration') return;
      assert.dom(SELECTORS.inputByAttr(attr)).doesNotExist(`does not render ${attr} when auto tidy disabled`);
    });

    sectionHeaders.forEach((group) => {
      assert.dom(SELECTORS.tidySectionHeader(group)).doesNotExist(`does not render ${group} header`);
    });

    // ENABLE AUTO TIDY
    await click(SELECTORS.toggleInput('intervalDuration'));
    assert.dom(SELECTORS.toggleInput('intervalDuration')).isChecked('Automatic tidy is enabled');
    assert.dom(`[data-test-ttl-form-label="Automatic tidy enabled"]`).exists('renders enabled text');

    this.autoTidy.eachAttribute((attr) => {
      const skipFields = ['enabled', 'tidyAcme', 'intervalDuration'];
      if (skipFields.includes(attr)) return; // combined with duration ttl or asserted elsewhere
      assert.dom(SELECTORS.inputByAttr(attr)).exists(`renders ${attr} when auto tidy enabled`);
    });

    sectionHeaders.forEach((group) => {
      assert.dom(SELECTORS.tidySectionHeader(group)).exists(`renders ${group} header`);
    });
  });

  test('it renders all attribute fields, including enterprise', async function (assert) {
    assert.expect(25);
    this.version.version = '1.14.1+ent';
    this.autoTidy.enabled = true;
    const skipFields = ['enabled', 'tidyAcme', 'intervalDuration']; // combined with duration ttl or asserted separately
    await render(
      hbs`
      <PkiTidyForm
      @tidy={{this.autoTidy}}
      @tidyType="auto"
      @onSave={{this.onSave}}
      @onCancel={{this.onCancel}}
    />
    `,
      { owner: this.engine }
    );

    this.autoTidy.eachAttribute((attr) => {
      if (skipFields.includes(attr)) return;
      assert.dom(SELECTORS.inputByAttr(attr)).exists(`renders ${attr} for auto tidyType`);
    });

    await render(
      hbs`
      <PkiTidyForm
      @tidy={{this.manualTidy}}
      @tidyType="manual"
      @onSave={{this.onSave}}
      @onCancel={{this.onCancel}}
    />
    `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.toggleInput('intervalDuration')).doesNotExist('hides automatic tidy toggle');

    this.manualTidy.eachAttribute((attr) => {
      if (skipFields.includes(attr)) return;
      assert.dom(SELECTORS.inputByAttr(attr)).exists(`renders ${attr} for manual tidyType`);
    });
  });

  test('it hides enterprise fields for OSS', async function (assert) {
    assert.expect(7);
    this.version.version = '1.14.1';
    this.autoTidy.enabled = true;

    const enterpriseFields = [
      'tidyRevocationQueue',
      'tidyCrossClusterRevokedCerts',
      'revocationQueueSafetyBuffer',
    ];

    // tidyType = auto
    await render(
      hbs`
      <PkiTidyForm
      @tidy={{this.autoTidy}}
      @tidyType="auto"
      @onSave={{this.onSave}}
      @onCancel={{this.onCancel}}
    />
    `,
      { owner: this.engine }
    );

    assert
      .dom(SELECTORS.tidySectionHeader('Cross-cluster operations'))
      .doesNotExist(`does not render ent header`);

    enterpriseFields.forEach((entAttr) => {
      assert.dom(SELECTORS.inputByAttr(entAttr)).doesNotExist(`does not render ${entAttr} for auto tidyType`);
    });

    // tidyType = manual
    await render(
      hbs`
      <PkiTidyForm
      @tidy={{this.manualTidy}}
      @tidyType="manual"
      @onSave={{this.onSave}}
      @onCancel={{this.onCancel}}
    />
    `,
      { owner: this.engine }
    );

    enterpriseFields.forEach((entAttr) => {
      assert
        .dom(SELECTORS.inputByAttr(entAttr))
        .doesNotExist(`does not render ${entAttr} for manual tidyType`);
    });
  });

  test('it should change the attributes on the model', async function (assert) {
    assert.expect(12);
    this.server.post('/pki-auto-tidy/config/auto-tidy', (schema, req) => {
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          acme_account_safety_buffer: '60s',
          enabled: true,
          interval_duration: '10s',
          issuer_safety_buffer: '20s',
          pause_duration: '30s',
          revocation_queue_safety_buffer: '40s',
          safety_buffer: '50s',
          tidy_acme: true,
          tidy_cert_store: true,
          tidy_cross_cluster_revoked_certs: true,
          tidy_expired_issuers: true,
          tidy_move_legacy_ca_bundle: true,
          tidy_revocation_queue: true,
          tidy_revoked_cert_issuer_associations: true,
          tidy_revoked_certs: true,
        },
        'response contains updated model values'
      );
    });
    await render(
      hbs`
      <PkiTidyForm
      @tidy={{this.autoTidy}}
      @tidyType="auto"
      @onSave={{this.onSave}}
      @onCancel={{this.onCancel}}
    />
    `,
      { owner: this.engine }
    );

    assert.dom(SELECTORS.toggleInput('intervalDuration')).isNotChecked('Automatic tidy is disabled');
    assert.dom(SELECTORS.toggleLabel('Automatic tidy disabled')).exists('auto tidy has disabled label');
    assert.false(this.autoTidy.enabled, 'enabled is false on model');

    // enable auto-tidy
    await click(SELECTORS.toggleInput('intervalDuration'));
    await fillIn(SELECTORS.intervalDuration, 10);

    assert.dom(SELECTORS.toggleInput('intervalDuration')).isChecked('toggle enabled auto tidy');
    assert.dom(SELECTORS.toggleLabel('Automatic tidy enabled')).exists('auto tidy has enabled label');

    assert.dom(SELECTORS.toggleInput('acmeAccountSafetyBuffer')).isNotChecked('ACME tidy is disabled');
    assert.dom(SELECTORS.toggleLabel('Tidy ACME disabled')).exists('ACME label has correct disabled text');
    assert.false(this.autoTidy.tidyAcme, 'tidyAcme is false on model');

    await click(SELECTORS.toggleInput('acmeAccountSafetyBuffer'));
    await fillIn(SELECTORS.acmeAccountSafetyBuffer, 60);
    assert.true(this.autoTidy.tidyAcme, 'tidyAcme toggles to true');

    const fillInValues = {
      issuerSafetyBuffer: 20,
      pauseDuration: 30,
      revocationQueueSafetyBuffer: 40,
      safetyBuffer: 50,
    };
    this.autoTidy.eachAttribute(async (attr, { type }) => {
      const skipFields = ['enabled', 'tidyAcme', 'intervalDuration', 'acmeAccountSafetyBuffer']; // combined with duration ttl or asserted separately
      if (skipFields.includes(attr)) return;
      if (type === 'boolean') {
        await click(SELECTORS.inputByAttr(attr));
      }
      if (type === 'string') {
        await fillIn(SELECTORS.toggleInput(attr), `${fillInValues[attr]}`);
      }
    });

    assert.dom(SELECTORS.toggleInput('acmeAccountSafetyBuffer')).isChecked('ACME tidy is enabled');
    assert.dom(SELECTORS.toggleLabel('Tidy ACME enabled')).exists('ACME label has correct enabled text');
    await click(SELECTORS.tidySave);
  });

  test('it updates auto-tidy config', async function (assert) {
    assert.expect(4);
    this.server.post('/pki-auto-tidy/config/auto-tidy', (schema, req) => {
      assert.ok(true, 'Request made to update auto-tidy');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          enabled: false,
          tidy_acme: false,
        },
        'response contains auto-tidy params'
      );
    });
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    this.onCancel = () => assert.ok(true, 'onCancel callback fires on save success');

    await render(
      hbs`
      <PkiTidyForm
        @tidy={{this.autoTidy}}
        @tidyType="auto"
        @onSave={{this.onSave}}
        @onCancel={{this.onCancel}}
      />
    `,
      { owner: this.engine }
    );

    await click(SELECTORS.tidySave);
    await click(SELECTORS.tidyCancel);
  });

  test('it saves and performs manual tidy', async function (assert) {
    assert.expect(4);

    this.server.post('/pki-manual-tidy/tidy', (schema, req) => {
      assert.ok(true, 'Request made to perform manual tidy');
      assert.propEqual(
        JSON.parse(req.requestBody),
        { tidy_acme: false },
        'response contains manual tidy params'
      );
    });
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    this.onCancel = () => assert.ok(true, 'onCancel callback fires on save success');

    await render(
      hbs`
      <PkiTidyForm
        @tidy={{this.manualTidy}}
        @tidyType="manual"
        @onSave={{this.onSave}}
        @onCancel={{this.onCancel}}
      />
    `,
      { owner: this.engine }
    );

    await click(SELECTORS.tidySave);
    await click(SELECTORS.tidyCancel);
  });
});
