/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'vault/tests/helpers';

module('Integration | Component | form/v2/wizard', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onCancel = sinon.spy();
    this.onSuccess = sinon.spy();

    this.wizardConfig = {
      title: 'Test Wizard',
      description: 'A test wizard',
      steps: [
        {
          name: 'step1',
          title: 'Step 1',
          description: 'First step',
          formConfig: {
            name: 'step1-form',
            path: '/v1/test/step1',
            payload: {
              name: '',
            },
            submit: sinon.stub().resolves({ id: 'step1-result' }),
            sections: [
              {
                name: 'basic',
                fields: [
                  {
                    name: 'name',
                    label: 'Name',
                    type: 'TextInput',
                    validations: [{ type: 'required' }],
                  },
                ],
              },
            ],
          },
        },
        {
          name: 'step2',
          title: 'Step 2',
          description: 'Second step',
          formConfig: {
            name: 'step2-form',
            path: '/v1/test/step2',
            payload: {
              description: '',
            },
            submit: sinon.stub().resolves({ id: 'step2-result' }),
            sections: [
              {
                name: 'details',
                fields: [
                  {
                    name: 'description',
                    label: 'Description',
                    type: 'TextArea',
                  },
                ],
              },
            ],
          },
        },
      ],
    };
  });

  test('it renders the first step content', async function (assert) {
    await render(hbs`
      <Form::V2::Wizard
        @config={{this.wizardConfig}}
        @onCancel={{this.onCancel}}
        @onSuccess={{this.onSuccess}}
      />
    `);

    assert.dom('label').includesText('Name', 'renders first step field');
    assert.dom(this.element).includesText('Step 1', 'renders first step title');
  });
});
