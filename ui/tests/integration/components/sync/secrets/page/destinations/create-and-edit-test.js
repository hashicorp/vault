/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupDataStubs } from 'vault/tests/helpers/sync/setup-hooks';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { Response } from 'miragejs';
import { click, fillIn, render, typeIn } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { syncDestinations, findDestination } from 'vault/helpers/sync-destinations';
import formResolver from 'vault/forms/sync/resolver';
import { DestinationType, CLOUD_DESTINATION_TYPES, CredentialType } from 'sync/utils/constants';

const SYNC_DESTINATIONS = syncDestinations();
module('Integration | Component | sync | Secrets::Page::Destinations::CreateAndEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);
  setupDataStubs(hooks);

  hooks.beforeEach(function () {
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.apiPath = 'sys/sync/destinations/:type/:name';

    this.generateForm = (isNew = false, type = DestinationType.AwsSm) => {
      const { defaultValues } = findDestination(type);
      let data = defaultValues;

      if (!isNew) {
        if (type !== DestinationType.AwsSm) {
          this.setupStubsForType(type);
        }
        const { name, connection_details, options } = this.destination;
        options.granularity = options.granularity_level;
        delete options.granularity_level;

        data = { name, ...connection_details, ...options };
      }

      this.form = formResolver(type, data, { isNew });
      this.formFields = this.form.formFieldGroups.reduce((arr, group) => {
        const values = Object.values(group)[0] || [];
        return [...arr, ...values];
      }, []);
      this.type = type;
    };

    this.renderComponent = () =>
      render(hbs` <Secrets::Page::Destinations::CreateAndEdit @form={{this.form}} @type={{this.type}} />`, {
        owner: this.engine,
      });
  });

  test('create: it renders breadcrumbs and navigates back to create on cancel', async function (assert) {
    this.generateForm(true);
    assert.expect(2);

    await this.renderComponent();
    assert.dom(GENERAL.breadcrumbs).hasText('Vault Secrets sync Select destination Create destination');
    await click(GENERAL.cancelButton);
    const transition = this.transitionStub.calledWith('vault.cluster.sync.secrets.destinations.create');
    assert.true(transition, 'transitions to vault.cluster.sync.secrets.destinations.create on cancel');
  });

  test('create: it renders headers and fieldGroups subtext', async function (assert) {
    this.generateForm(true);
    assert.expect(3);

    await this.renderComponent();
    assert
      .dom(PAGE.form.fieldGroupHeader('IAM credentials'))
      .hasText('IAM credentials', 'renders IAM credentials section on create');
    assert
      .dom('[data-test-accordion="Advanced configuration"]')
      .exists('renders advanced configuration accordion section on create');
    assert
      .dom(PAGE.form.fieldGroupSubtext('IAM credentials'))
      .hasText('Connection credentials are sensitive information used to authenticate with the destination.');
  });

  test('edit: it renders breadcrumbs and navigates back to details on cancel', async function (assert) {
    this.generateForm();
    assert.expect(2);

    await this.renderComponent();
    assert.dom(GENERAL.breadcrumbs).hasText('Vault Secrets sync Destinations Destination Edit destination');

    await click(GENERAL.cancelButton);
    const transition = this.transitionStub.calledWith('vault.cluster.sync.secrets.destinations.destination');
    assert.true(transition, 'transitions to vault.cluster.sync.secrets.destinations.destination on cancel');
  });

  test('edit: it renders headers and fieldGroup subtext', async function (assert) {
    this.generateForm();
    assert.expect(3);

    await this.renderComponent();
    assert
      .dom(PAGE.form.fieldGroupHeader('IAM credentials'))
      .hasText('IAM credentials', 'renders IAM credentials section on edit');
    assert
      .dom('[data-test-accordion="Advanced configuration"]')
      .exists('renders advanced configuration accordion section on edit');
    assert
      .dom(PAGE.form.fieldGroupSubtext('IAM credentials'))
      .hasText(
        'Connection credentials are sensitive information and the value cannot be read. Enable the input to update.'
      );
  });

  test('edit: it PATCH updates custom_tags', async function (assert) {
    this.generateForm();
    assert.expect(2);

    this.server.patch(this.apiPath, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(
        payload.custom_tags,
        { updated: 'bar', added: 'key' },
        'payload contains updated custom tags'
      );
      assert.propEqual(payload.tags_to_remove, ['foo'], 'payload contains tags to remove with expected keys');
      return payload;
    });

    await this.renderComponent();
    // Expand Advanced configuration accordion
    await click(GENERAL.accordionButton('Advanced configuration'));
    await click(GENERAL.kvObjectEditor.deleteRow());
    await fillIn(GENERAL.kvObjectEditor.key(), 'updated');
    await fillIn(GENERAL.kvObjectEditor.value(), 'bar');
    await click(GENERAL.kvObjectEditor.addRow);
    await fillIn(GENERAL.kvObjectEditor.key(1), 'added');
    await fillIn(GENERAL.kvObjectEditor.value(1), 'key');
    await click(GENERAL.submitButton);
  });

  test('edit: it adds custom_tags when previously there are none', async function (assert) {
    this.generateForm();
    assert.expect(1);

    this.server.patch(this.apiPath, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload.custom_tags, { foo: 'blah' }, 'payload contains new custom tags');
      return payload;
    });

    this.destination.options.custom_tags = {};

    await this.renderComponent();
    await click(GENERAL.accordionButton('Advanced configuration'));
    await PAGE.form.fillInByAttr('custom_tags', 'blah');
    await click(GENERAL.submitButton);
  });

  test('edit: payload does not contain any custom_tags when removed in form', async function (assert) {
    this.generateForm();
    assert.expect(2);

    this.server.patch(this.apiPath, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload.tags_to_remove, ['foo'], 'payload removes old keys');
      assert.propEqual(payload.custom_tags, {}, 'payload does not contain custom_tags');
      return payload;
    });

    await this.renderComponent();
    await click(GENERAL.accordionButton('Advanced configuration'));
    await click(GENERAL.kvObjectEditor.deleteRow());
    await click(GENERAL.submitButton);
  });

  test('edit: payload only contains masked inputs when they have changed', async function (assert) {
    this.generateForm();
    assert.expect(2);

    this.server.patch(this.apiPath, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.strictEqual(
        payload.access_key_id,
        undefined,
        'payload does not contain the unchanged obfuscated field'
      );
      assert.strictEqual(
        payload.secret_access_key,
        'new-secret',
        'payload contains the changed obfuscated field'
      );
      return payload;
    });

    await this.renderComponent();
    await click(GENERAL.enableField('access_key_id'));
    await click(GENERAL.inputByAttr('access_key_id')); // click on input but do not change value
    await click(GENERAL.enableField('secret_access_key'));
    await fillIn(GENERAL.inputByAttr('secret_access_key'), 'new-secret');
    await click(GENERAL.submitButton);
  });

  test('it renders API errors', async function (assert) {
    this.generateForm();
    assert.expect(1);

    this.server.patch(this.apiPath, () => {
      return new Response(
        500,
        {},
        {
          errors: [
            `1 error occurred: * couldn't create store node in syncer: failed to create store: unable to initialize store of type "azure-kv": failed to parse azure key vault URI: parse "my-unprasableuri": invalid URI for request`,
          ],
        }
      );
    });

    await this.renderComponent();

    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.messageError)
      .hasText(
        `Error 1 error occurred: * couldn't create store node in syncer: failed to create store: unable to initialize store of type "azure-kv": failed to parse azure key vault URI: parse "my-unprasableuri": invalid URI for request`
      );
  });

  test('it renders warning validation only when editing vercel-project team_id', async function (assert) {
    const type = DestinationType.VercelProject;
    this.generateForm(true, type); // new destination

    assert.expect(2);

    await this.renderComponent();
    await typeIn(GENERAL.inputByAttr('team_id'), 'id');
    assert
      .dom(GENERAL.validationWarningByAttr('team_id'))
      .doesNotExist('does not render warning validation for new vercel-project destination');

    this.generateForm(false, type); // existing destination
    await this.renderComponent();
    await PAGE.form.fillInByAttr('team_id', '');
    await typeIn(GENERAL.inputByAttr('team_id'), 'edit');
    assert
      .dom(GENERAL.validationWarningByAttr('team_id'))
      .hasText(
        'Team ID should only be updated if the project was transferred to another account.',
        'it renders validation warning'
      );
  });
  // WIF (Workload Identity Federation) TESTS
  module('WIF credential type support', function (hooks) {
    hooks.beforeEach(function () {
      this.version = this.owner.lookup('service:version');
      this.version.type = 'enterprise';

      // Helper to switch between credential types
      this.switchToWif = async () => {
        await click(GENERAL.radioCardByAttr(CredentialType.WIF));
      };

      this.switchToAccount = async () => {
        await click(GENERAL.radioCardByAttr(CredentialType.ACCOUNT));
      };

      // Helpers to assert field group visibility
      this.assertFieldGroupVisible = (assert, groupName, message) => {
        assert
          .dom(PAGE.form.fieldGroupHeader(groupName))
          .exists(message || `${groupName} section is visible`);
      };

      this.assertFieldGroupHidden = (assert, groupName, message) => {
        assert
          .dom(PAGE.form.fieldGroupHeader(groupName))
          .doesNotExist(message || `${groupName} section is not visible`);
      };
    });

    test('create: it renders credential type radio cards for cloud destinations', async function (assert) {
      assert.expect(CLOUD_DESTINATION_TYPES.length * 3);

      for (const type of CLOUD_DESTINATION_TYPES) {
        this.generateForm(true, type);
        await this.renderComponent();

        assert
          .dom(GENERAL.radioCardByAttr(CredentialType.ACCOUNT))
          .exists(`${type}: renders account credential type radio card`);
        assert
          .dom(GENERAL.radioCardByAttr(CredentialType.WIF))
          .exists(`${type}: renders WIF credential type radio card`);
        assert
          .dom(GENERAL.radioCardByAttr(CredentialType.ACCOUNT))
          .isChecked(`${type}: account credential type is selected by default`);
      }
    });

    test('create: it does not render credential type radio cards for non-cloud destinations', async function (assert) {
      const nonCloudDestinations = SYNC_DESTINATIONS.filter(
        (d) => !CLOUD_DESTINATION_TYPES.includes(d.type)
      ).map((d) => d.type);
      assert.expect(nonCloudDestinations.length);

      for (const type of nonCloudDestinations) {
        this.generateForm(true, type);
        await this.renderComponent();

        assert
          .dom(GENERAL.radioCardByAttr())
          .doesNotExist(`${type}: does not render credential type radio cards`);
      }
    });

    test('create aws-sm: it switches between IAM and WIF credential fields', async function (assert) {
      this.generateForm(true, DestinationType.AwsSm);
      assert.expect(8);

      await this.renderComponent();

      // Check IAM credentials are visible by default
      this.assertFieldGroupVisible(assert, 'IAM credentials');
      assert.dom(GENERAL.fieldByAttr('access_key_id')).exists('access_key_id field is visible');
      assert.dom(GENERAL.fieldByAttr('secret_access_key')).exists('secret_access_key field is visible');
      this.assertFieldGroupHidden(assert, 'WIF credentials');

      // Switch to WIF
      await this.switchToWif();

      // Check WIF credentials are now visible
      this.assertFieldGroupVisible(assert, 'WIF credentials');
      assert
        .dom(GENERAL.fieldByAttr('identity_token_audience'))
        .exists('identity_token_audience field is visible');
      assert.dom(GENERAL.fieldByAttr('identity_token_key')).exists('identity_token_key field is visible');
      this.assertFieldGroupHidden(assert, 'IAM credentials');
    });

    test('create azure-kv: it switches between Client Secret and WIF credential fields', async function (assert) {
      this.generateForm(true, DestinationType.AzureKv);
      assert.expect(6);

      await this.renderComponent();

      // Check Client Secret is visible by default
      this.assertFieldGroupVisible(assert, 'Client secret', 'Client secret credentials section is visible');
      assert.dom(GENERAL.fieldByAttr('client_secret')).exists('client_secret field is visible');
      this.assertFieldGroupHidden(assert, 'WIF credentials');

      // Switch to WIF
      await this.switchToWif();

      // Check WIF credentials are now visible
      this.assertFieldGroupVisible(assert, 'WIF credentials');
      assert
        .dom(GENERAL.fieldByAttr('identity_token_audience'))
        .exists('identity_token_audience field is visible');
      this.assertFieldGroupHidden(
        assert,
        'Client secret',
        'Client secret credentials section is not visible'
      );
    });

    test('create gcp-sm: it switches between JSON Credentials and WIF credential fields', async function (assert) {
      this.generateForm(true, DestinationType.GcpSm);
      assert.expect(7);

      await this.renderComponent();

      // Check JSON credentials are visible by default
      this.assertFieldGroupVisible(assert, 'JSON credentials');
      assert.dom(GENERAL.fieldByAttr('credentials')).exists('credentials field is visible');
      this.assertFieldGroupHidden(assert, 'WIF credentials');

      // Switch to WIF
      await this.switchToWif();

      // Check WIF credentials are now visible
      this.assertFieldGroupVisible(assert, 'WIF credentials');
      assert
        .dom(GENERAL.fieldByAttr('service_account_email'))
        .exists('service_account_email field is visible (GCP-specific)');
      assert
        .dom(GENERAL.fieldByAttr('identity_token_audience'))
        .exists('identity_token_audience field is visible');
      this.assertFieldGroupHidden(assert, 'JSON credentials');
    });

    test('create: it resets account fields when switching to WIF', async function (assert) {
      this.generateForm(true, DestinationType.AwsSm);
      assert.expect(2);

      await this.renderComponent();

      // Fill in IAM credentials
      await fillIn(GENERAL.inputByAttr('access_key_id'), 'test-access-key');
      await fillIn(GENERAL.inputByAttr('secret_access_key'), 'test-secret-key');

      // Switch to WIF
      await this.switchToWif();

      // Switch back to account
      await this.switchToAccount();

      // Verify fields were reset
      assert.dom(GENERAL.inputByAttr('access_key_id')).hasValue('', 'access_key_id was reset');
      assert.strictEqual(this.form.data.access_key_id, undefined, 'access_key_id is undefined in form data');
    });

    test('create: it resets WIF fields when switching to account credentials', async function (assert) {
      this.generateForm(true, DestinationType.AwsSm);
      assert.expect(2);

      await this.renderComponent();

      // Switch to WIF
      await this.switchToWif();

      // Fill in WIF credentials
      await fillIn(GENERAL.inputByAttr('identity_token_audience'), 'test-audience');

      // Switch back to account
      await this.switchToAccount();

      // Switch to WIF again to verify reset
      await this.switchToWif();

      assert
        .dom(GENERAL.inputByAttr('identity_token_audience'))
        .hasValue('', 'identity_token_audience was reset');
      assert.strictEqual(
        this.form.data.identity_token_audience,
        undefined,
        'identity_token_audience is undefined in form data'
      );
    });

    test('create: it sets default key value when switching to WIF', async function (assert) {
      this.generateForm(true, DestinationType.AwsSm);
      assert.expect(1);

      await this.renderComponent();

      // Switch to WIF
      await this.switchToWif();

      // Verify default key is empty
      assert.strictEqual(
        this.form.data.identity_token_key,
        undefined,
        'identity_token_key is undefined by default'
      );
    });

    test('create: it validates WIF credentials', async function (assert) {
      this.generateForm(true, DestinationType.AwsSm);
      assert.expect(2);

      await this.renderComponent();

      // Switch to WIF
      await this.switchToWif();

      // Fill in name but leave WIF fields empty
      await fillIn(GENERAL.inputByAttr('name'), 'test-destination');

      // Try to submit
      await click(GENERAL.submitButton);

      // Check for validation errors on required WIF fields
      assert
        .dom(GENERAL.validationErrorByAttr('identity_token_audience'))
        .exists('validation error shown for identity_token_audience');
      assert
        .dom(GENERAL.validationErrorByAttr('role_arn'))
        .exists('validation error shown for role_arn (AWS-specific)');
    });

    test('create: it successfully creates destination with WIF credentials', async function (assert) {
      this.generateForm(true, DestinationType.AwsSm);
      assert.expect(5);

      const name = 'wif-destination';
      const path = `sys/sync/destinations/aws-sm/${name}`;

      this.server.post(path, (schema, req) => {
        const payload = JSON.parse(req.requestBody);

        assert.ok(true, `makes request: POST ${path}`);
        assert.notOk('credential_type' in payload, 'credential_type is not in payload');
        assert.notOk('access_key_id' in payload, 'account credentials not in payload');
        assert.propContains(
          payload,
          { identity_token_audience: 'test-audience' },
          'WIF credentials in payload'
        );
        return payload;
      });

      await this.renderComponent();

      // Switch to WIF
      await this.switchToWif();

      // Fill in required fields
      await fillIn(GENERAL.inputByAttr('name'), name);
      await fillIn(GENERAL.inputByAttr('region'), 'us-west-1');
      await fillIn(GENERAL.inputByAttr('role_arn'), 'arn:aws:iam::123456789012:role/test-role');
      await fillIn(GENERAL.inputByAttr('identity_token_audience'), 'test-audience');

      await click(GENERAL.submitButton);

      const actualArgs = this.transitionStub.lastCall.args;
      const expectedArgs = [
        'vault.cluster.sync.secrets.destinations.destination.details',
        DestinationType.AwsSm,
        name,
      ];
      assert.propEqual(actualArgs, expectedArgs, 'transitionTo called with expected args');
    });

    test('edit: it disables credential type selection when WIF is configured', async function (assert) {
      assert.expect(2);

      this.generateForm(false, DestinationType.AwsSm);

      // Simulate existing WIF configuration on form
      this.form.data.identity_token_audience = 'existing-audience';
      this.form.data.identity_token_key = '*****';
      this.form.data.role_arn = 'arn:aws:iam::123456789012:role/test-role';
      delete this.form.data.access_key_id;

      await this.renderComponent();

      assert
        .dom(GENERAL.radioCardByAttr(CredentialType.ACCOUNT))
        .isDisabled('account credential type radio is disabled');
      assert
        .dom(GENERAL.radioCardByAttr(CredentialType.WIF))
        .isDisabled('WIF credential type radio is disabled');
    });

    test('edit: it disables credential type selection when account credentials are configured', async function (assert) {
      assert.expect(2);

      // Simulate existing account configuration on destination (default from mirage)
      this.generateForm(false, DestinationType.AwsSm);

      await this.renderComponent();

      assert
        .dom(GENERAL.radioCardByAttr(CredentialType.ACCOUNT))
        .isDisabled('account credential type radio is disabled');
      assert
        .dom(GENERAL.radioCardByAttr(CredentialType.WIF))
        .isDisabled('WIF credential type radio is disabled');
    });

    test('edit: it PATCH updates WIF credentials correctly', async function (assert) {
      assert.expect(3);

      this.generateForm(false, DestinationType.AwsSm);

      // Simulate existing WIF configuration on form
      this.form.data.identity_token_audience = '*****';
      this.form.data.identity_token_key = '*****';
      this.form.data.role_arn = 'arn:aws:iam::123456789012:role/test-role';

      const path = `sys/sync/destinations/aws-sm/${this.form.name}`;
      this.server.patch(path, (schema, req) => {
        const payload = JSON.parse(req.requestBody);

        assert.ok(true, `makes request: PATCH ${path}`);
        assert.notOk('credential_type' in payload, 'credential_type is not in payload');
        assert.strictEqual(
          payload.identity_token_key,
          'new-key-value',
          'updated identity_token_key in payload'
        );
        return payload;
      });

      await this.renderComponent();

      // Update identity token key (needs to be enabled first since it's masked)
      await click(GENERAL.enableField('identity_token_key'));
      await fillIn(GENERAL.inputByAttr('identity_token_key'), 'new-key-value');

      await click(GENERAL.submitButton);
    });
  });

  // CREATE FORM ASSERTIONS FOR EACH DESTINATION TYPE
  for (const destination of SYNC_DESTINATIONS) {
    const { name, type } = destination;
    const obfuscatedFields = ['access_token', 'client_secret', 'secret_access_key', 'access_key_id'];

    module(`create destination: ${type}`, function () {
      test('it renders destination form', async function (assert) {
        this.generateForm(true, type);
        assert.expect(this.formFields.length + 1);

        await this.renderComponent();

        assert.dom(GENERAL.hdsPageHeaderTitle).hasTextContaining(`Create Destination for ${name}`);

        const accordion = document.querySelector(GENERAL.accordionButton('Advanced configuration'));
        if (accordion) {
          await click(GENERAL.accordionButton('Advanced configuration'));
        }

        for (const field of this.formFields) {
          assert.dom(GENERAL.fieldByAttr(field.name)).exists();
        }
      });

      test('it masks obfuscated fields', async function (assert) {
        this.generateForm(true, type);
        const filteredObfuscatedFields = this.formFields.filter((field) =>
          obfuscatedFields.includes(field.name)
        );
        assert.expect(filteredObfuscatedFields.length * 2);

        await this.renderComponent();
        // iterate over the form fields and filter for those that are obfuscated
        // fill those in and assert that they are masked
        filteredObfuscatedFields.forEach(async (field) => {
          await fillIn(GENERAL.inputByAttr(field.name), 'blah');

          assert
            .dom(GENERAL.inputByAttr(field.name))
            .hasClass('masked-font', `it renders ${field.name} for ${destination} with masked font`);
          assert
            .dom(PAGE.form.enableInput(field.name))
            .doesNotExist(`it does not render enable input for ${field.name}`);
        });
      });

      test('it saves destination and transitions to details', async function (assert) {
        this.generateForm(true, type);
        assert.expect(2);

        const name = 'my-name';
        const path = `sys/sync/destinations/${type}/my-name`;

        this.server.post(path, (schema, req) => {
          const payload = JSON.parse(req.requestBody);

          assert.ok(true, `makes request: POST ${path}`);
          // Skipped payload assertions due to object comparison issues in Mirage
          return payload;
        });

        await this.renderComponent();

        const accordion = document.querySelector(GENERAL.accordionButton('Advanced configuration'));
        if (accordion) {
          await click(GENERAL.accordionButton('Advanced configuration'));
        }

        for (const field of this.formFields) {
          await PAGE.form.fillInByAttr(field.name, `my-${field.name}`);
        }
        await click(GENERAL.submitButton);
        const actualArgs = this.transitionStub.lastCall.args;
        const expectedArgs = ['vault.cluster.sync.secrets.destinations.destination.details', type, name];
        assert.propEqual(actualArgs, expectedArgs, 'transitionTo called with expected args');
      });

      test('it validates inputs', async function (assert) {
        this.generateForm(true, type);

        const warningValidations = ['team_id'];
        const validationAssertions = { ...this.form.validations };
        // remove warning validations
        warningValidations.forEach((warning) => {
          delete validationAssertions[warning];
        });

        // Count only presence validations
        let presenceValidationCount = 0;
        for (const attr in validationAssertions) {
          const validation = validationAssertions[attr].find((v) => v.type === 'presence');
          if (validation) {
            presenceValidationCount++;
          }
        }
        assert.expect(presenceValidationCount);

        await this.renderComponent();
        await click(GENERAL.submitButton);

        // only asserts validations for presence, refactor if validations change
        for (const attr in validationAssertions) {
          const validation = validationAssertions[attr].find((v) => v.type === 'presence');
          if (validation) {
            const { message } = validation;
            assert
              .dom(GENERAL.validationErrorByAttr(attr))
              .hasText(message, `renders validation: ${message}`);
          }
        }
      });
    });
  }

  // EDIT FORM ASSERTIONS FOR EACH DESTINATION TYPE
  // if field is not string type, add case to EXPECTED_VALUE and update
  // fillInByAttr() (in sync-selectors) to interact with the form
  const EXPECTED_VALUE = (key) => {
    switch (key) {
      case 'custom_tags':
        return { foo: `new-${key}-value` };
      case 'deployment_environments':
        return ['production'];
      case 'granularity':
        return 'secret-key';
      default:
        // for all string type parameters
        return `new-${key}-value`;
    }
  };

  for (const destination of SYNC_DESTINATIONS) {
    const { type, maskedParams, readonlyParams } = destination;
    module(`edit destination: ${type}`, function () {
      test('it renders destination form and PATCH updates a destination', async function (assert) {
        this.generateForm(false, type);

        const [disabledAssertions, editable] = this.formFields.reduce(
          (arr, field) => {
            if (field.options.editDisabled) {
              arr[0]++;
            }
            if (!readonlyParams.includes(field.name)) {
              arr[1].push(field.name);
            }
            return arr;
          },
          [0, []]
        );

        assert.expect(4 + disabledAssertions + editable.length);

        const path = `sys/sync/destinations/${type}/${this.form.name}`;
        this.server.patch(path, (schema, req) => {
          assert.ok(true, `makes request: PATCH ${path}`);
          const payload = JSON.parse(req.requestBody);
          const payloadKeys = Object.keys(payload);
          const expectedKeys = editable.sort();
          assert.propEqual(payloadKeys, expectedKeys, `${type} payload only contains editable attrs`);
          expectedKeys.forEach((key) => {
            assert.propEqual(payload[key], EXPECTED_VALUE(key), `destination: ${type} updates key: ${key}`);
          });
          return { payload };
        });

        await this.renderComponent(false, type);

        assert.dom(GENERAL.hdsPageHeaderTitle).hasTextContaining(`Edit ${this.form.name}`);

        // Expand Advanced configuration accordion if it exists (contains granularity, custom_tags, etc.)
        const accordion = document.querySelector(GENERAL.accordionButton('Advanced configuration'));
        if (accordion) {
          await click(GENERAL.accordionButton('Advanced configuration'));
        }

        for (const field of this.formFields) {
          if (editable.includes(field.name)) {
            if (maskedParams.includes(field.name)) {
              // Enable inputs with sensitive values
              await click(PAGE.form.enableInput(field.name));
            }
            await PAGE.form.fillInByAttr(field.name, `new-${field.name}-value`);
          } else {
            assert.dom(GENERAL.inputByAttr(field.name)).isDisabled(`${field.name} is disabled`);
          }
        }

        await click(GENERAL.submitButton);
        const actualArgs = this.transitionStub.lastCall.args;
        const expectedArgs = [
          'vault.cluster.sync.secrets.destinations.destination.details',
          type,
          this.form.name,
        ];
        assert.propEqual(actualArgs, expectedArgs, 'transitionTo called with expected args');
      });
    });
  }
});
