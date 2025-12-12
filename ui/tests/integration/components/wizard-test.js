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
  welcome: '[data-test-welcome-content]',
  quickstart: '[data-test-quickstart-content]',
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
    this.step = 0;
    this.showWelcome = false;
    this.onDismiss = sinon.spy();
    this.onStepChange = sinon.spy();
  });

  test('it shows welcome content initially, then hides it when entering wizard', async function (assert) {
    this.set('showWelcome', true);
    await render(hbs`<Wizard
    @title="Example Wizard"
    @showWelcome={{this.showWelcome}}
    @currentStep={{this.step}}
    @steps={{this.steps}}
    @onStepChange={{this.onStepChange}}
    @onDismiss={{this.onDismiss}}
  >
    <:welcome>
      <div>Some welcome content</div>
      {{!-- TODO: This will change once the welcome page structure is defined in a follow up PR --}}
      <Hds::Button @text="Dismiss" {{on "click" this.onDismiss}} />
      <Hds::Button @text="Enter wizard" {{on "click" (fn (mut this.showWelcome) false)}} data-test-enter-wizard-button />
    </:welcome>
    <:quickstart>
      <div data-test-quickstart-content>Quickstart content</div>
    </:quickstart>
    </Wizard>`);

    // Assert welcome content is rendered and quickstart content is not
    assert.dom(SELECTORS.welcome).exists('Welcome content is rendered initially');
    assert.dom(SELECTORS.welcome).hasTextContaining('Some welcome content');
    assert
      .dom(SELECTORS.quickstart)
      .doesNotExist('Quickstart content is not rendered when welcome is displayed');

    await click('[data-test-enter-wizard-button]');

    // Assert welcome content is no longer rendered and that quickstart content is rendered
    assert.dom(SELECTORS.welcome).doesNotExist('Welcome content is hidden after entering wizard');
    assert.dom(SELECTORS.quickstart).exists('Quickstart content is now rendered');
    assert.dom(SELECTORS.quickstart).hasTextContaining('Quickstart content');
  });

  test('it shows custom submit block when provided', async function (assert) {
    // Go to final step
    this.set('step', 3);
    this.onCustomSubmit = sinon.spy();

    await render(hbs`<Wizard
      @title="Example Wizard"
      @showWelcome={{this.showWelcome}}
      @currentStep={{this.step}}
      @steps={{this.steps}}
      @onStepChange={{this.onStepChange}}
      @onDismiss={{this.onDismiss}}
    >
      <:quickstart>
        <div data-test-quickstart-content>Quickstart content</div>
      </:quickstart>
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
    this.set('step', 3);

    await render(hbs`<Wizard
      @title="Example Wizard"
      @showWelcome={{this.showWelcome}}
      @currentStep={{this.step}}
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
      @showWelcome={{this.showWelcome}}
      @currentStep={{this.step}}
      @steps={{this.steps}}
      @onStepChange={{this.onStepChange}}
      @onDismiss={{this.onDismiss}}
    >
      <:quickstart>
        <div>Quickstart content</div>
      </:quickstart>
    </Wizard>`);

    assert.dom(GENERAL.nextButton).exists('Next button is rendered when not on final step');
    assert.dom(GENERAL.submitButton).doesNotExist('Submit button is not rendered when not on final step');
    await click(GENERAL.nextButton);
    assert.true(this.onStepChange.calledOnce, 'onStepChange is called');
    // Go to final step
    this.set('step', 3);
    assert.dom(GENERAL.nextButton).doesNotExist('Next button is not rendered when on the final step');
    assert.dom(GENERAL.submitButton).exists('Submit button is rendered on final step');
  });

  test('it renders back button when not on first step', async function (assert) {
    await render(hbs`<Wizard
      @title="Example Wizard"
      @showWelcome={{this.showWelcome}}
      @currentStep={{this.step}}
      @steps={{this.steps}}
      @onStepChange={{this.onStepChange}}
      @onDismiss={{this.onDismiss}}
    >
      <:quickstart>
        <div>Quickstart content</div>
      </:quickstart>
    </Wizard>`);

    assert.dom(GENERAL.backButton).doesNotExist('Back button is not rendered on the first step');
    this.set('step', 2);
    assert.dom(GENERAL.backButton).exists('Back button is shown when not on first step');
    await click(GENERAL.backButton);
    assert.true(this.onStepChange.calledOnce, 'onStepChange is called');
  });

  test('it dismisses wizard when exit button is clicked within quickstart', async function (assert) {
    await render(hbs`<Wizard
      @title="Example Wizard"
      @showWelcome={{this.showWelcome}}
      @currentStep={{this.step}}
      @steps={{this.steps}}
      @onStepChange={{this.onStepChange}}
      @onDismiss={{this.onDismiss}}
    >
      <:quickstart>
        <div>Quickstart content</div>
      </:quickstart>
    </Wizard>`);

    assert.dom(GENERAL.cancelButton).exists('Exit button is shown within quickstart');
    await click(GENERAL.cancelButton);
    assert.true(this.onDismiss.calledOnce, 'onDismiss is called when exit button is clicked');
  });
});
