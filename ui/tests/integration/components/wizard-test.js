/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  welcome: '[data-test-welcome]',
  guidedSetup: '[data-test-guided-setup]',
};

module('Integration | Component | Wizard', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.steps = [
      { title: 'First step' },
      { title: 'Another stage' },
      { title: 'Almost done' },
      { title: 'Finale' },
    ];
    this.currentStep = 0;
    this.canProceed = true;
    this.welcomeDocLink = 'test';
    this.onDismiss = sinon.spy();
    this.onStepChange = sinon.spy();
  });

  test('it shows welcome content initially, then hides it when entering wizard', async function (assert) {
    await render(hbs`<Wizard
    @title="Example Wizard"
    @currentStep={{this.currentStep}}
    @steps={{this.steps}}
    @welcomeDocLink={{this.welcomeDocLink}}
    @onStepChange={{this.onStepChange}}
    @onDismiss={{this.onDismiss}}
  >
    <:welcome>
      <div>Some welcome content</div>
    </:welcome>
    </Wizard>`);

    // Assert welcome content is rendered and guided setup content is not
    assert.dom(SELECTORS.welcome).exists('Welcome content is rendered initially');
    assert.dom(SELECTORS.welcome).hasTextContaining('Some welcome content');
    assert
      .dom(SELECTORS.guidedSetup)
      .doesNotExist('guidedSetup content is not rendered when welcome is displayed');

    await click(GENERAL.button('Guided setup'));

    // Assert welcome content is no longer rendered and that guided setup content is rendered
    assert.dom(SELECTORS.welcome).doesNotExist('Welcome content is hidden after entering wizard');
    assert.dom(SELECTORS.guidedSetup).exists('guidedSetup content is now rendered');
    assert.dom(SELECTORS.guidedSetup).hasTextContaining('First step');
  });

  test('it shows custom submit block when provided', async function (assert) {
    // Go to final step
    this.currentStep = 3;
    this.onCustomSubmit = sinon.spy();

    await render(hbs`<Wizard
      @title="Example Wizard"
      @currentStep={{this.currentStep}}
      @steps={{this.steps}}
      @onStepChange={{this.onStepChange}}
      @onDismiss={{this.onDismiss}}
    >
      <:submit>
        <Hds::Button @text="Custom Submit" {{on "click" this.onCustomSubmit}} data-test-custom-submit />
      </:submit>
    </Wizard>`);

    assert.dom('[data-test-custom-submit]').exists('Custom submit button is rendered');
    assert.dom(GENERAL.submitButton).doesNotExist('Default submit button is not rendered');
    await click('[data-test-custom-submit]');
    assert.true(this.onCustomSubmit.calledOnce, 'Custom submit handler is called');
  });

  test('it shows default submit button when custom submit block is not provided', async function (assert) {
    // Go to final step
    this.currentStep = 3;

    await render(hbs`<Wizard
      @title="Example Wizard"
      @showWelcome={{this.showWelcome}}
      @currentStep={{this.currentStep}}
      @steps={{this.steps}}
      @onStepChange={{this.onStepChange}}
      @onDismiss={{this.onDismiss}}
    >
      <:quickstart>
        <div>Quickstart content</div>
      </:quickstart>
    </Wizard>`);

    assert
      .dom(GENERAL.submitButton)
      .exists('Default submit button is rendered when no custom submit provided');
  });

  test('it renders next button when not on final step', async function (assert) {
    await render(hbs`<Wizard
      @title="Example Wizard"
      @canProceed={{this.canProceed}}
      @currentStep={{this.currentStep}}
      @steps={{this.steps}}
      @onStepChange={{this.onStepChange}}
      @onDismiss={{this.onDismiss}}
    >
    </Wizard>`);

    assert.dom(GENERAL.button('Next')).exists('Next button is rendered when not on final step');
    assert.dom(GENERAL.submitButton).doesNotExist('Submit button is not rendered when not on final step');
    await click(GENERAL.button('Next'));
    assert.true(this.onStepChange.calledOnce, 'onStepChange is called');
    // Go to final step
    this.set('currentStep', 3);
    assert.dom(GENERAL.button('next')).doesNotExist('Next button is not rendered when on the final step');
    assert.dom(GENERAL.submitButton).exists('Submit button is rendered on final step');
  });

  test('it renders back button when not on first step', async function (assert) {
    await render(hbs`<Wizard
      @title="Example Wizard"
      @showWelcome={{this.showWelcome}}
      @currentStep={{this.currentStep}}
      @steps={{this.steps}}
      @onStepChange={{this.onStepChange}}
      @onDismiss={{this.onDismiss}}
    >
    </Wizard>`);

    assert.dom(GENERAL.backButton).doesNotExist('Back button is not rendered on the first step');
    this.set('currentStep', 2);
    assert.dom(GENERAL.backButton).exists('Back button is shown when not on first step');
    await click(GENERAL.backButton);
    assert.true(this.onStepChange.calledOnce, 'onStepChange is called');
  });

  test('it dismisses wizard when exit button is clicked within guided setup', async function (assert) {
    await render(hbs`<Wizard
      @title="Example Wizard"
      @showWelcome={{this.showWelcome}}
      @currentStep={{this.currentStep}}
      @steps={{this.steps}}
      @onStepChange={{this.onStepChange}}
      @onDismiss={{this.onDismiss}}
    >
    </Wizard>`);

    assert.dom(GENERAL.cancelButton).exists('Exit button is shown within guided setup');
    await click(GENERAL.cancelButton);
    assert.true(this.onDismiss.calledOnce, 'onDismiss is called when exit button is clicked');
  });
});
