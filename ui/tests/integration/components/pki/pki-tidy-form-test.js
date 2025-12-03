/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { PKI_TIDY_FORM } from 'vault/tests/helpers/pki/pki-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { convertToSeconds } from 'core/utils/duration-utils';
import PkiTidyForm from 'vault/forms/secrets/pki/tidy';
import TidyGroupsHelper from 'pki/helpers/tidy-groups';
import tidyFieldLabel from 'pki/helpers/tidy-field-label';
import sinon from 'sinon';

module('Integration | Component | pki tidy form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.backend = 'pki-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';

    this.tidyGroupsHelper = new TidyGroupsHelper(this.owner);
    const groupsReducer = (fields, group) => [...fields, ...Object.values(group)[0]];
    this.allFields = this.tidyGroupsHelper.compute([]).reduce(groupsReducer, []);
    this.autoFields = this.tidyGroupsHelper.compute(['auto']).reduce(groupsReducer, []);

    this.onSave = sinon.spy();
    this.onCancel = sinon.spy();

    const { secrets } = this.owner.lookup('service:api');
    this.tidyStub = sinon.stub(secrets, 'pkiTidy').resolves();
    this.autoTidyStub = sinon.stub(secrets, 'pkiConfigureAutoTidy').resolves();
    // setting defaults here to simulate how this form works in the app.
    // on init, we retrieve these from the server and pre-populate form
    this.autoTidyServerDefaults = {
      enabled: false,
      acme_account_safety_buffer: '48h',
      interval_duration: '12h',
      safety_buffer: '3d',
      issuer_safety_buffer: '365d',
      min_startup_backoff_duration: '5m',
      max_startup_backoff_duration: '15m',
      revocation_queue_safety_buffer: '10s',
      pause_duration: '0',
    };

    this.renderComponent = (type = 'auto') => {
      this.tidyType = type;
      const formData = type === 'auto' ? this.autoTidyServerDefaults : {};
      const schemaKey = type === 'auto' ? 'PkiConfigureAutoTidyRequest' : 'PkiTidyRequest';
      this.form = new PkiTidyForm(schemaKey, formData, { isNew: type === 'manual' });
      return render(
        hbs`<PkiTidyForm @form={{this.form}} @tidyType={{this.tidyType}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
        { owner: this.engine }
      );
    };
  });

  test('it hides or shows fields depending on auto-tidy toggle', async function (assert) {
    const sectionHeaders = [
      'Automatic tidy settings',
      'Universal operations',
      'ACME operations',
      'Issuer operations',
      'Cross-cluster operations',
    ];
    const headerAssertCount = sectionHeaders.length * 2;

    await this.renderComponent();
    // the form isn't created until we render so expect assertions need to be after render
    const loopAssertCount = this.allFields.length * 2 - 3; // loop skips 3 params
    assert.expect(loopAssertCount + headerAssertCount + 4);

    assert.dom(GENERAL.toggleInput('enabled')).isNotChecked();
    assert
      .dom(GENERAL.ttl.toggle('enabled'))
      .hasText('Automatic tidy disabled Automatic tidy operations will not run.');

    this.allFields.forEach((field) => {
      if (field !== 'enabled') {
        assert
          .dom(PKI_TIDY_FORM.inputByAttr(field))
          .doesNotExist(`does not render ${field} when auto tidy disabled`);
      }
    });

    sectionHeaders.forEach((group) => {
      assert.dom(PKI_TIDY_FORM.tidySectionHeader(group)).doesNotExist(`does not render ${group} header`);
    });

    // ENABLE AUTO TIDY
    await click(GENERAL.toggleInput('enabled'));
    assert.dom(GENERAL.toggleInput('enabled')).isChecked();
    assert.dom(GENERAL.ttl.toggle('enabled')).hasText('Automatic tidy enabled');

    this.allFields.forEach((field) => {
      const skipFields = ['enabled', 'tidy_acme'];
      // combined with duration ttl or asserted elsewhere
      if (!skipFields.includes(field)) {
        assert.dom(PKI_TIDY_FORM.inputByAttr(field)).exists(`renders ${field} when auto tidy enabled`);
      }
    });

    sectionHeaders.forEach((group) => {
      assert.dom(PKI_TIDY_FORM.tidySectionHeader(group)).exists(`renders ${group} header`);
    });
  });

  test('it renders all attribute fields, including enterprise', async function (assert) {
    assert.expect(35);

    await this.renderComponent();
    await click(GENERAL.toggleInput('enabled'));

    const skipFields = ['enabled', 'tidy_acme']; // combined with duration ttl or asserted separately
    this.allFields.forEach((field) => {
      if (!skipFields.includes(field)) {
        assert.dom(PKI_TIDY_FORM.inputByAttr(field)).exists(`renders ${field} for auto tidyType`);
      }
    });

    // MANUAL TIDY
    await this.renderComponent('manual');

    assert.dom(GENERAL.toggleInput('enabled')).doesNotExist('hides automatic tidy toggle');

    this.allFields.forEach((field) => {
      if (!skipFields.includes(field)) {
        // auto tidy fields we shouldn't see in the manual tidy form
        if (this.autoFields.includes(field)) {
          assert
            .dom(PKI_TIDY_FORM.inputByAttr(field))
            .doesNotExist(`${field} should not appear on manual tidyType`);
        } else {
          assert.dom(PKI_TIDY_FORM.inputByAttr(field)).exists(`renders ${field} for manual tidyType`);
        }
      }
    });
  });

  test('it hides enterprise fields for CE', async function (assert) {
    this.version.type = 'community';

    const enterpriseFields = [
      'tidy_revocation_queue',
      'tidy_cross_cluster_revoked_certs',
      'tidy_cmpv2_nonce_store',
      'revocation_queue_safety_buffer',
    ];

    // tidyType = auto
    await this.renderComponent();
    await click(GENERAL.toggleInput('enabled'));

    assert
      .dom(PKI_TIDY_FORM.tidySectionHeader('Cross-cluster operations'))
      .doesNotExist(`does not render ent header`);

    enterpriseFields.forEach((entAttr) => {
      assert
        .dom(PKI_TIDY_FORM.inputByAttr(entAttr))
        .doesNotExist(`does not render ${entAttr} for auto tidyType`);
    });

    // tidyType = manual
    await this.renderComponent('manual');

    enterpriseFields.forEach((entAttr) => {
      assert
        .dom(PKI_TIDY_FORM.inputByAttr(entAttr))
        .doesNotExist(`does not render ${entAttr} for manual tidyType`);
    });
  });

  test('it should update form values', async function (assert) {
    assert.expect(11);
    // ttl picker defaults to seconds, unless unit is set by default value (set in beforeEach hook)
    // on submit, any user inputted values should be converted to seconds for the payload
    const fillInValues = {
      acme_account_safety_buffer: { time: 680, unit: 'h' },
      interval_duration: { time: 10, unit: 'h' },
      issuer_safety_buffer: { time: 20, unit: 'd' },
      max_startup_backoff_duration: { time: 30, unit: 'm' },
      min_startup_backoff_duration: { time: 10, unit: 'm' },
      pause_duration: { time: 30, unit: 's' },
      revocation_queue_safety_buffer: { time: 40, unit: 's' },
      safety_buffer: { time: 50, unit: 'd' },
    };
    const calcValue = (param) => {
      const { time, unit } = fillInValues[param];
      return `${convertToSeconds(time, unit)}s`;
    };

    await this.renderComponent();

    assert.dom(GENERAL.toggleInput('enabled')).isNotChecked();
    assert.dom(GENERAL.ttl.toggle('enabled')).hasTextContaining('Automatic tidy disabled');
    assert.false(this.form.data.enabled, 'enabled is false on form');

    // enable auto-tidy
    await click(GENERAL.toggleInput('enabled'));
    assert.dom(GENERAL.toggleInput('enabled')).isChecked();
    assert.dom(GENERAL.ttl.toggle('enabled')).hasText('Automatic tidy enabled');

    assert.dom(PKI_TIDY_FORM.toggleInput('acme_account_safety_buffer')).isNotChecked('ACME tidy is disabled');
    assert
      .dom(PKI_TIDY_FORM.toggleLabel('Tidy ACME disabled'))
      .exists('ACME label has correct disabled text');

    await click(PKI_TIDY_FORM.toggleInput('acme_account_safety_buffer'));
    await fillIn(PKI_TIDY_FORM.acmeAccountSafetyBuffer, 2); // units are days based on defaultValue
    assert.dom(PKI_TIDY_FORM.toggleInput('acme_account_safety_buffer')).isChecked('ACME tidy is enabled');
    assert.dom(PKI_TIDY_FORM.toggleLabel('Tidy ACME enabled')).exists('ACME label has correct enabled text');
    assert.true(this.form.data.tidy_acme, 'tidy_acme toggles to true');

    for (const field of this.allFields) {
      const skipFields = ['enabled', 'tidy_acme', 'acme_account_safety_buffer']; // combined with duration ttl or asserted separately
      if (!skipFields.includes(field)) {
        // all params right now are either a boolean or TTL, this if/else will need to be updated if that changes
        if (Object.keys(fillInValues).includes(field)) {
          const { time } = fillInValues[field];
          await fillIn(GENERAL.ttl.input(tidyFieldLabel(field)), `${time}`);
        } else {
          await click(PKI_TIDY_FORM.inputByAttr(field));
        }
      }
    }

    await click(PKI_TIDY_FORM.tidySave);
    const payload = {
      acme_account_safety_buffer: '48h',
      enabled: true,
      min_startup_backoff_duration: calcValue('min_startup_backoff_duration'),
      max_startup_backoff_duration: calcValue('max_startup_backoff_duration'),
      interval_duration: calcValue('interval_duration'),
      issuer_safety_buffer: calcValue('issuer_safety_buffer'),
      pause_duration: calcValue('pause_duration'),
      revocation_queue_safety_buffer: calcValue('revocation_queue_safety_buffer'),
      safety_buffer: calcValue('safety_buffer'),
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
    };
    assert.true(this.autoTidyStub.calledWith(this.backend, payload), 'API called with correct payload');
  });

  test('it updates auto-tidy config', async function (assert) {
    assert.expect(3);

    await this.renderComponent();

    await click(PKI_TIDY_FORM.tidySave);
    assert.true(
      this.autoTidyStub.calledWith(this.backend, this.autoTidyServerDefaults),
      'API called with correct payload'
    );
    assert.true(this.onSave.called, 'onSave called on save success');

    await click(PKI_TIDY_FORM.tidyCancel);
    assert.true(this.onCancel.called, 'onCancel called on click');
  });

  test('it saves and performs manual tidy', async function (assert) {
    assert.expect(3);

    await this.renderComponent('manual');

    await click(PKI_TIDY_FORM.tidySave);
    assert.true(this.tidyStub.calledWith(this.backend), 'API called with correct payload');
    assert.true(this.onSave.called, 'onSave called on save success');

    await click(PKI_TIDY_FORM.tidyCancel);
    assert.true(this.onCancel.called, 'onCancel called on click');
  });
});
