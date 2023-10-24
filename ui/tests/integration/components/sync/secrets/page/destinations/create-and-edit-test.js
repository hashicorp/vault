/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { click, render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { syncDestinations } from 'vault/helpers/sync-destinations';
import { underscore } from '@ember/string';

const SYNC_DESTINATIONS = syncDestinations();
module('Integration | Component | sync | Secrets::Page::Destinations::CreateAndEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.renderFormComponent = () => {
      return render(hbs` <Secrets::Page::Destinations::CreateAndEdit @destination={{this.model}} />`, {
        owner: this.engine,
      });
    };
  });

  test('it navigates back on cancel', async function (assert) {
    assert.expect(1);
    const type = SYNC_DESTINATIONS[0].type;
    this.model = this.store.createRecord(`sync/destinations/${type}`, { type });

    await this.renderFormComponent();

    await click(PAGE.form.cancelButton);
    const transition = this.transitionStub.calledWith('vault.cluster.sync.secrets.destinations.create');
    assert.true(transition, 'transitions to vault.cluster.sync.secrets.destinations.create on cancel');
  });

  // module runs for each destination type
  for (const destination of SYNC_DESTINATIONS) {
    const { displayName, type } = destination;

    module(`destination: ${type}`, function (hooks) {
      hooks.beforeEach(function () {
        this.model = this.store.createRecord(`sync/destinations/${type}`, { type });
      });

      test('it renders destination form', async function (assert) {
        assert.expect(this.model.formFields.length + 1);

        await this.renderFormComponent();

        assert.dom(PAGE.title).hasTextContaining(`Create destination for ${displayName}`);
        for (const attr of this.model.formFields) {
          assert.dom(PAGE.inputByAttr(attr.name)).exists();
        }
      });

      test('it saves destination and transitions to details', async function (assert) {
        assert.expect(4);
        const name = 'my-name';
        this.server.post(`sys/sync/destinations/${type}/${name}`, (schema, req) => {
          const payload = JSON.parse(req.requestBody);
          assert.ok(true, `makes request: POST sys/sync/destinations/${type}/${name}`);
          assert.notPropContains(payload, { name: 'my-name', type }, 'name and type do not exist in payload');

          // instead of looping through all attrs, just grab the second one (first is 'name')
          const testAttr = this.model.formFields[1].name;
          assert.propContains(
            payload,
            { [underscore(testAttr)]: `my-${testAttr}` },
            'payload contains expected attrs'
          );
          return payload;
        });

        await this.renderFormComponent();

        for (const attr of this.model.formFields) {
          await PAGE.form.fillInByAttr(attr.name, `my-${attr.name}`);
        }
        await click(PAGE.form.saveButton);
        const actualArgs = this.transitionStub.lastCall.args;
        const expectedArgs = ['vault.cluster.sync.secrets.destinations.destination.details', type, name];
        assert.propEqual(actualArgs, expectedArgs, 'transitionTo called with expected args');
      });

      test('it validates inputs', async function (assert) {
        const validations = this.model._validations;
        assert.expect(Object.keys(validations).length);

        await this.renderFormComponent();

        await click(PAGE.form.saveButton);

        // only asserts validations for presence, may want to refactor if validations change
        for (const attr in validations) {
          const { message } = validations[attr].find((v) => v.type === 'presence');
          assert.dom(PAGE.validation(attr)).hasText(message, `renders validation: ${message}`);
        }
      });
    });
  }
});
