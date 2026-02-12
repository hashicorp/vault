/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, waitFor } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  intro: '[data-test-intro]',
  guidedStart: '[data-test-guided-start]',
};

module('Integration | Component | Wizard', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.canProceed = true;
    this.currentStep = 0;
    this.isModal = false;
    this.steps = [
      { title: 'First step' },
      { title: 'Another stage' },
      { title: 'Almost done' },
      { title: 'Finale' },
    ];
    this.title = 'Example Wizard';
    this.wizardState = undefined;
    this.updateWizardState = sinon.spy();
    this.onDismiss = sinon.spy();
    this.onStepChange = sinon.spy();
    this.wizardId = 'test-wizard';
    this.wizardService = this.owner.lookup('service:wizard');
    this.wizardService.setIntroVisible(this.wizardId, false);

    this.renderComponent = () => {
      return render(hbs`
        <Wizard
          @wizardId={{this.wizardId}}
          @canProceed={{this.canProceed}}
          @currentStep={{this.currentStep}}
          @isModal={{@this.isIntroModal}}
          @steps={{this.steps}}
          @title="Example Wizard"
          @wizardState={{this.wizardState}}
          @updateWizardState={{this.updateWizardState}}
          @onStepChange={{this.onStepChange}}
          @onDismiss={{this.onDismiss}}
      >
        <:intro>
          <div>Some intro content</div>
        </:intro>
        <:introActions>
          <div> Some actions </div>
        </:introActions>
        <:submit>
          <Hds::Button @text="Custom Submit" data-test-custom-submit />
        </:submit>
        <:exit>
          <Hds::Button @text="Custom Exit" data-test-custom-exit />
        </:exit>
      </Wizard>`);
    };
  });

  test('it shows intro content initially, then hides it when entering wizard', async function (assert) {
    this.wizardService.setIntroVisible(this.wizardId, true);
    await this.renderComponent();

    // Assert intro content is rendered and guided start content is not
    assert.dom(SELECTORS.intro).exists('intro content is rendered initially');
    assert.dom(SELECTORS.intro).hasTextContaining('Some intro content');
    assert
      .dom(SELECTORS.guidedStart)
      .doesNotExist('guidedStart content is not rendered when intro is displayed');

    // Use wizard service to hide the intro
    this.wizardService.setIntroVisible(this.wizardId, false);
    await waitFor(SELECTORS.guidedStart);

    // Assert intro content is no longer rendered and that guided start content is rendered
    assert.dom(SELECTORS.intro).doesNotExist('intro content is hidden after entering wizard');
    assert.dom(SELECTORS.guidedStart).exists('guidedStart content is now rendered');
    assert.dom(SELECTORS.guidedStart).hasTextContaining('First step');
  });

  test('it shows custom submit block when provided', async function (assert) {
    // Start wizard and go to final step
    this.currentStep = 3;

    await this.renderComponent();

    assert.dom('[data-test-custom-submit]').exists('Custom submit button is rendered');
    assert.dom(GENERAL.submitButton).doesNotExist('Default submit button is not rendered');
  });

  test('it shows default submit button when custom submit block is not provided', async function (assert) {
    // Start wizard and go to final step
    this.currentStep = 3;

    await render(hbs`
      <Wizard
        @wizardId={{this.wizardId}}
        @canProceed={{this.canProceed}}
        @currentStep={{this.currentStep}}
        @isModal={{@this.isIntroModal}}
        @steps={{this.steps}}
        @title="Example Wizard"
        @wizardState={{this.wizardState}}
        @updateWizardState={{this.updateWizardState}}
        @onStepChange={{this.onStepChange}}
        @onDismiss={{this.onDismiss}}
      >
        <:intro>
          <div>Some intro content</div>
        </:intro>
        <:introActions>
          <div> Some actions </div>
        </:introActions>
      </Wizard>`);

    assert
      .dom(GENERAL.submitButton)
      .exists('Default submit button is rendered when no custom submit provided');
  });

  test('it shows custom exit block when provided', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-custom-exit]').exists('Custom exit button is rendered');
    assert.dom(GENERAL.cancelButton).doesNotExist('Default exit button is not rendered');
  });

  test('it shows default exit button when custom exit block is not provided', async function (assert) {
    await render(hbs`
      <Wizard
        @wizardId={{this.wizardId}}
        @canProceed={{this.canProceed}}
        @currentStep={{this.currentStep}}
        @isModal={{@this.isIntroModal}}
        @showIntro={{this.showIntro}}
        @steps={{this.steps}}
        @title="Example Wizard"
        @wizardState={{this.wizardState}}
        @updateWizardState={{this.updateWizardState}}
        @onStepChange={{this.onStepChange}}
        @onDismiss={{this.onDismiss}}
      >
        <:intro>
          <div>Some intro content</div>
        </:intro>
        <:introActions>
          <div> Some actions </div>
        </:introActions>
      </Wizard>`);

    assert.dom(GENERAL.cancelButton).exists('Default exit button is rendered when no custom exit provided');
    await click(GENERAL.cancelButton);
    assert.true(this.onDismiss.calledOnce, 'onDismiss is called when exit button is clicked');
  });

  test('it renders next button when not on final step', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.button('Next')).exists('Next button is rendered when not on final step');
    assert
      .dom('[data-test-custom-submit]')
      .doesNotExist('Custom submit button is not rendered when not on final step');
    await click(GENERAL.button('Next'));
    assert.true(this.onStepChange.calledOnce, 'onStepChange is called');
    // Go to final step
    this.set('currentStep', 3);
    assert.dom(GENERAL.button('Next')).doesNotExist('Next button is not rendered when on the final step');
    assert.dom('[data-test-custom-submit]').exists('Custom submit button is rendered on final step');
  });

  test('it renders back button when not on first step', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.backButton).doesNotExist('Back button is not rendered on the first step');
    this.set('currentStep', 2);
    assert.dom(GENERAL.backButton).exists('Back button is shown when not on first step');
    await click(GENERAL.backButton);
    assert.true(this.onStepChange.calledOnce, 'onStepChange is called');
  });
});
