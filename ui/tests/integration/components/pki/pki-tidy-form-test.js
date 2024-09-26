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
import { PKI_TIDY_FORM } from 'vault/tests/helpers/pki/pki-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { convertToSeconds } from 'core/utils/duration-utils';

module('Integration | Component | pki tidy form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';
    this.server.post('/sys/capabilities-self', () => {});
    this.onSave = () => {};
    this.onCancel = () => {};
    this.manualTidy = this.store.createRecord('pki/tidy', { backend: 'pki-manual-tidy' });
    this.autoTidyServerDefaults = {
      enabled: false,
      interval_duration: '12h',
      safety_buffer: '3d',
      issuer_safety_buffer: '365d',
      min_startup_backoff_duration: '5m',
      max_startup_backoff_duration: '15m',
    };
    this.store.pushPayload('pki/tidy', {
      modelName: 'pki/tidy',
      id: 'pki-auto-tidy',
      // setting defaults here to simulate how this form works in the app.
      // on init, we retrieve these from the server and pre-populate form (instead of explicitly set on the model)
      ...this.autoTidyServerDefaults,
    });
    this.autoTidy = this.store.peekRecord('pki/tidy', 'pki-auto-tidy');
    this.numTidyAttrs = Object.keys(this.autoTidy.allByKey).length;
  });

  test('it hides or shows fields depending on auto-tidy toggle', async function (assert) {
    const sectionHeaders = [
      'Automatic tidy settings',
      'Universal operations',
      'ACME operations',
      'Issuer operations',
      'Cross-cluster operations',
    ];
    const loopAssertCount = this.numTidyAttrs * 2 - 3; // loop skips 3 params
    const headerAssertCount = sectionHeaders.length * 2;
    assert.expect(loopAssertCount + headerAssertCount + 4);

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
    assert.dom(GENERAL.toggleInput('enabled')).isNotChecked();
    assert
      .dom(GENERAL.ttl.toggle('enabled'))
      .hasText('Automatic tidy disabled Automatic tidy operations will not run.');

    this.autoTidy.eachAttribute((attr) => {
      if (attr === 'enabled') return;
      assert
        .dom(PKI_TIDY_FORM.inputByAttr(attr))
        .doesNotExist(`does not render ${attr} when auto tidy disabled`);
    });

    sectionHeaders.forEach((group) => {
      assert.dom(PKI_TIDY_FORM.tidySectionHeader(group)).doesNotExist(`does not render ${group} header`);
    });

    // ENABLE AUTO TIDY
    await click(GENERAL.toggleInput('enabled'));
    assert.dom(GENERAL.toggleInput('enabled')).isChecked();
    assert.dom(GENERAL.ttl.toggle('enabled')).hasText('Automatic tidy enabled');

    this.autoTidy.eachAttribute((attr) => {
      const skipFields = ['enabled', 'tidyAcme'];
      if (skipFields.includes(attr)) return; // combined with duration ttl or asserted elsewhere
      assert.dom(PKI_TIDY_FORM.inputByAttr(attr)).exists(`renders ${attr} when auto tidy enabled`);
    });

    sectionHeaders.forEach((group) => {
      assert.dom(PKI_TIDY_FORM.tidySectionHeader(group)).exists(`renders ${group} header`);
    });
  });

  test('it renders all attribute fields, including enterprise', async function (assert) {
    assert.expect(35);
    this.autoTidy.enabled = true;
    const skipFields = ['enabled', 'tidyAcme']; // combined with duration ttl or asserted separately
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
      assert.dom(PKI_TIDY_FORM.inputByAttr(attr)).exists(`renders ${attr} for auto tidyType`);
    });

    // MANUAL TIDY
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
    assert.dom(GENERAL.toggleInput('enabled')).doesNotExist('hides automatic tidy toggle');

    this.manualTidy.eachAttribute((attr) => {
      if (skipFields.includes(attr)) return;
      // auto tidy fields we shouldn't see in the manual tidy form
      if (this.manualTidy.autoTidyConfigFields.includes(attr)) {
        assert
          .dom(PKI_TIDY_FORM.inputByAttr(attr))
          .doesNotExist(`${attr} should not appear on manual tidyType`);
      } else {
        assert.dom(PKI_TIDY_FORM.inputByAttr(attr)).exists(`renders ${attr} for manual tidyType`);
      }
    });
  });

  test('it hides enterprise fields for CE', async function (assert) {
    this.version.type = 'community';
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
      .dom(PKI_TIDY_FORM.tidySectionHeader('Cross-cluster operations'))
      .doesNotExist(`does not render ent header`);

    enterpriseFields.forEach((entAttr) => {
      assert
        .dom(PKI_TIDY_FORM.inputByAttr(entAttr))
        .doesNotExist(`does not render ${entAttr} for auto tidyType`);
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
        .dom(PKI_TIDY_FORM.inputByAttr(entAttr))
        .doesNotExist(`does not render ${entAttr} for manual tidyType`);
    });
  });

  test('it should change the attributes on the model', async function (assert) {
    assert.expect(12);
    // ttl picker defaults to seconds, unless unit is set by default value (set in beforeEach hook)
    // on submit, any user inputted values should be converted to seconds for the payload
    const fillInValues = {
      acmeAccountSafetyBuffer: { time: 680, unit: 'h' },
      intervalDuration: { time: 10, unit: 'h' },
      issuerSafetyBuffer: { time: 20, unit: 'd' },
      maxStartupBackoffDuration: { time: 30, unit: 'm' },
      minStartupBackoffDuration: { time: 10, unit: 'm' },
      pauseDuration: { time: 30, unit: 's' },
      revocationQueueSafetyBuffer: { time: 40, unit: 's' },
      safetyBuffer: { time: 50, unit: 'd' },
    };
    const calcValue = (param) => {
      const { time, unit } = fillInValues[param];
      return `${convertToSeconds(time, unit)}s`;
    };
    this.server.post('/pki-auto-tidy/config/auto-tidy', (schema, req) => {
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          acme_account_safety_buffer: '48h',
          enabled: true,
          min_startup_backoff_duration: calcValue('minStartupBackoffDuration'),
          max_startup_backoff_duration: calcValue('maxStartupBackoffDuration'),
          interval_duration: calcValue('intervalDuration'),
          issuer_safety_buffer: calcValue('issuerSafetyBuffer'),
          pause_duration: calcValue('pauseDuration'),
          revocation_queue_safety_buffer: calcValue('revocationQueueSafetyBuffer'),
          safety_buffer: calcValue('safetyBuffer'),
          tidy_acme: true,
          tidy_cert_metadata: true,
          tidy_cert_store: true,
          tidy_cmpv2_nonce_store: true,
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

    assert.dom(GENERAL.toggleInput('enabled')).isNotChecked();
    assert.dom(GENERAL.ttl.toggle('enabled')).hasTextContaining('Automatic tidy disabled');
    assert.false(this.autoTidy.enabled, 'enabled is false on model');

    // enable auto-tidy
    await click(GENERAL.toggleInput('enabled'));
    assert.dom(GENERAL.toggleInput('enabled')).isChecked();
    assert.dom(GENERAL.ttl.toggle('enabled')).hasText('Automatic tidy enabled');

    assert.dom(PKI_TIDY_FORM.toggleInput('acmeAccountSafetyBuffer')).isNotChecked('ACME tidy is disabled');
    assert
      .dom(PKI_TIDY_FORM.toggleLabel('Tidy ACME disabled'))
      .exists('ACME label has correct disabled text');
    assert.false(this.autoTidy.tidyAcme, 'tidyAcme is false on model');

    await click(PKI_TIDY_FORM.toggleInput('acmeAccountSafetyBuffer'));
    await fillIn(PKI_TIDY_FORM.acmeAccountSafetyBuffer, 2); // units are days based on defaultValue
    assert.dom(PKI_TIDY_FORM.toggleInput('acmeAccountSafetyBuffer')).isChecked('ACME tidy is enabled');
    assert.dom(PKI_TIDY_FORM.toggleLabel('Tidy ACME enabled')).exists('ACME label has correct enabled text');
    assert.true(this.autoTidy.tidyAcme, 'tidyAcme toggles to true');

    this.autoTidy.eachAttribute(async (attr, { type }) => {
      const skipFields = ['enabled', 'tidyAcme', 'acmeAccountSafetyBuffer']; // combined with duration ttl or asserted separately
      if (skipFields.includes(attr)) return;
      // all params right now are either a boolean or TTL, this if/else will need to be updated if that changes
      if (type === 'boolean') {
        await click(PKI_TIDY_FORM.inputByAttr(attr));
      } else {
        const { time } = fillInValues[attr];
        await fillIn(PKI_TIDY_FORM.toggleInput(attr), `${time}`);
      }
    });

    await click(PKI_TIDY_FORM.tidySave);
  });

  test('it updates auto-tidy config', async function (assert) {
    assert.expect(4);
    this.server.post('/pki-auto-tidy/config/auto-tidy', (schema, req) => {
      assert.ok(true, 'Request made to update auto-tidy');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          ...this.autoTidyServerDefaults,
          acme_account_safety_buffer: '720h',
          tidy_acme: false,
        },
        'response contains default auto-tidy params'
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

    await click(PKI_TIDY_FORM.tidySave);
    await click(PKI_TIDY_FORM.tidyCancel);
  });

  test('it saves and performs manual tidy', async function (assert) {
    assert.expect(4);

    this.server.post('/pki-manual-tidy/tidy', (schema, req) => {
      assert.ok(true, 'Request made to perform manual tidy');
      assert.propEqual(
        JSON.parse(req.requestBody),
        { acme_account_safety_buffer: '720h', tidy_acme: false },
        'response contains manual tidy params'
      );
      return { id: 'pki-manual-tidy' };
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

    await click(PKI_TIDY_FORM.tidySave);
    await click(PKI_TIDY_FORM.tidyCancel);
  });
});
