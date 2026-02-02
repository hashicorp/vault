/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
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

const SYNC_DESTINATIONS = syncDestinations();
module('Integration | Component | sync | Secrets::Page::Destinations::CreateAndEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);
  setupDataStubs(hooks);

  hooks.beforeEach(function () {
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.apiPath = 'sys/sync/destinations/:type/:name';

    this.generateForm = (isNew = false, type = 'aws-sm') => {
      const { defaultValues } = findDestination(type);
      let data = defaultValues;

      if (!isNew) {
        if (type !== 'aws-sm') {
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
    assert.dom(PAGE.breadcrumbs).hasText('Secrets Sync Select Destination Create Destination');
    await click(PAGE.cancelButton);
    const transition = this.transitionStub.calledWith('vault.cluster.sync.secrets.destinations.create');
    assert.true(transition, 'transitions to vault.cluster.sync.secrets.destinations.create on cancel');
  });

  test('create: it renders headers and fieldGroups subtext', async function (assert) {
    this.generateForm(true);
    assert.expect(4);

    await this.renderComponent();
    assert
      .dom(PAGE.form.fieldGroupHeader('Credentials'))
      .hasText('Credentials', 'renders credentials section on create');
    assert
      .dom(PAGE.form.fieldGroupHeader('Advanced configuration'))
      .hasText('Advanced configuration', 'renders advanced configuration section on create');
    assert
      .dom(PAGE.form.fieldGroupSubtext('Credentials'))
      .hasText('Connection credentials are sensitive information used to authenticate with the destination.');
    assert
      .dom(PAGE.form.fieldGroupSubtext('Advanced configuration'))
      .hasText('Configuration options for the destination.');
  });

  test('edit: it renders breadcrumbs and navigates back to details on cancel', async function (assert) {
    this.generateForm();
    assert.expect(2);

    await this.renderComponent();
    assert.dom(PAGE.breadcrumbs).hasText('Secrets Sync Destinations Destination Edit Destination');

    await click(PAGE.cancelButton);
    const transition = this.transitionStub.calledWith('vault.cluster.sync.secrets.destinations.destination');
    assert.true(transition, 'transitions to vault.cluster.sync.secrets.destinations.destination on cancel');
  });

  test('edit: it renders headers and fieldGroup subtext', async function (assert) {
    this.generateForm();
    assert.expect(4);

    await this.renderComponent();
    assert
      .dom(PAGE.form.fieldGroupHeader('Credentials'))
      .hasText('Credentials', 'renders credentials section on edit');
    assert
      .dom(PAGE.form.fieldGroupHeader('Advanced configuration'))
      .hasText('Advanced configuration', 'renders advanced configuration section on edit');
    assert
      .dom(PAGE.form.fieldGroupSubtext('Credentials'))
      .hasText(
        'Connection credentials are sensitive information and the value cannot be read. Enable the input to update.'
      );
    assert
      .dom(PAGE.form.fieldGroupSubtext('Advanced configuration'))
      .hasText('Configuration options for the destination.');
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
    await click(PAGE.kvObjectEditor.deleteRow());
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
    await click(PAGE.enableField('access_key_id'));
    await click(PAGE.inputByAttr('access_key_id')); // click on input but do not change value
    await click(PAGE.enableField('secret_access_key'));
    await fillIn(PAGE.inputByAttr('secret_access_key'), 'new-secret');
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
      .dom(PAGE.messageError)
      .hasText(
        `Error 1 error occurred: * couldn't create store node in syncer: failed to create store: unable to initialize store of type "azure-kv": failed to parse azure key vault URI: parse "my-unprasableuri": invalid URI for request`
      );
  });

  test('it renders warning validation only when editing vercel-project team_id', async function (assert) {
    const type = 'vercel-project';
    this.generateForm(true, type); // new destination

    assert.expect(2);

    await this.renderComponent();
    await typeIn(PAGE.inputByAttr('team_id'), 'id');
    assert
      .dom(PAGE.validationWarningByAttr('team_id'))
      .doesNotExist('does not render warning validation for new vercel-project destination');

    this.generateForm(false, type); // existing destination
    await this.renderComponent();
    await PAGE.form.fillInByAttr('team_id', '');
    await typeIn(PAGE.inputByAttr('team_id'), 'edit');
    assert
      .dom(PAGE.validationWarningByAttr('team_id'))
      .hasText(
        'Team ID should only be updated if the project was transferred to another account.',
        'it renders validation warning'
      );
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

        for (const field of this.formFields) {
          assert.dom(PAGE.fieldByAttr(field.name)).exists();
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
          await fillIn(PAGE.inputByAttr(field.name), 'blah');

          assert
            .dom(PAGE.inputByAttr(field.name))
            .hasClass('masked-font', `it renders ${field.name} for ${destination} with masked font`);
          assert
            .dom(PAGE.form.enableInput(field.name))
            .doesNotExist(`it does not render enable input for ${field.name}`);
        });
      });

      test('it saves destination and transitions to details', async function (assert) {
        this.generateForm(true, type);
        assert.expect(4);

        const name = 'my-name';
        const path = `sys/sync/destinations/${type}/my-name`;

        this.server.post(path, (schema, req) => {
          const payload = JSON.parse(req.requestBody);

          assert.ok(true, `makes request: POST ${path}`);
          assert.notPropContains(payload, { name, type }, 'name and type do not exist in payload');
          // instead of looping through all attrs, just grab the second one (first is 'name')
          const testAttr = this.formFields[1].name;
          assert.propContains(payload, { [testAttr]: `my-${testAttr}` }, 'payload contains expected attrs');
          return payload;
        });

        await this.renderComponent();

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
        assert.expect(Object.keys(validationAssertions).length);

        await this.renderComponent();
        await click(GENERAL.submitButton);

        // only asserts validations for presence, refactor if validations change
        for (const attr in validationAssertions) {
          const { message } = validationAssertions[attr].find((v) => v.type === 'presence');
          assert.dom(PAGE.validationErrorByAttr(attr)).hasText(message, `renders validation: ${message}`);
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

        for (const field of this.formFields) {
          if (editable.includes(field.name)) {
            if (maskedParams.includes(field.name)) {
              // Enable inputs with sensitive values
              await click(PAGE.form.enableInput(field.name));
            }
            await PAGE.form.fillInByAttr(field.name, `new-${field.name}-value`);
          } else {
            assert.dom(PAGE.inputByAttr(field.name)).isDisabled(`${field.name} is disabled`);
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
