/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import localStorage from 'vault/lib/local-storage';
import { DISMISSED_WIZARD_KEY } from 'vault/utils/constants/wizard';

const STEPS = [
  { title: 'Step 1', component: 'wizard/step-1' },
  { title: 'Step 2', component: 'wizard/step-2' },
  { title: 'Step 3', component: 'wizard/step-3' },
];

module('Unit | Service | wizard', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.service = this.owner.lookup('service:wizard');

    // Stub localStorage methods
    this.getItemStub = sinon.stub(localStorage, 'getItem');
    this.setItemStub = sinon.stub(localStorage, 'setItem');
    this.removeItemStub = sinon.stub(localStorage, 'removeItem');
  });

  hooks.afterEach(function () {
    this.getItemStub.restore();
    this.setItemStub.restore();
    this.removeItemStub.restore();
  });

  module('#loadDismissedWizards', function () {
    test('loads dismissed wizards from localStorage', function (assert) {
      this.getItemStub.withArgs(DISMISSED_WIZARD_KEY).returns(['wizard1', 'wizard2']);
      const service = this.owner.lookup('service:wizard');

      assert.deepEqual(
        service.dismissedWizards,
        ['wizard1', 'wizard2'],
        'loads dismissed wizards from localStorage'
      );
      assert.true(this.getItemStub.calledWith(DISMISSED_WIZARD_KEY), 'calls localStorage.getItem');
    });

    test('returns empty array when localStorage has no dismissed wizards', function (assert) {
      this.getItemStub.withArgs(DISMISSED_WIZARD_KEY).returns(null);
      const service = this.owner.lookup('service:wizard');

      assert.deepEqual(service.dismissedWizards, [], 'returns empty array when no wizards dismissed');
    });
  });

  module('#isDismissed', function () {
    test('returns true when wizard is in dismissed list', function (assert) {
      this.service.dismissedWizards = ['onboarding', 'tutorial'];

      assert.true(this.service.isDismissed('onboarding'), 'returns true for dismissed wizard');
      assert.true(this.service.isDismissed('tutorial'), 'returns true for another dismissed wizard');
    });

    test('returns false when wizard is not in dismissed list', function (assert) {
      this.service.dismissedWizards = ['onboarding'];

      assert.false(this.service.isDismissed('tutorial'), 'returns false for non-dismissed wizard');
    });
  });

  module('#dismiss', function () {
    test('adds wizard to dismissed list', function (assert) {
      this.service.dismissedWizards = [];

      this.service.dismiss('onboarding');

      assert.deepEqual(this.service.dismissedWizards, ['onboarding'], 'adds wizard to dismissed list');
      assert.true(this.setItemStub.calledWith(DISMISSED_WIZARD_KEY, ['onboarding']), 'saves to localStorage');
    });

    test('handles multiple dismissals', function (assert) {
      this.service.dismissedWizards = [];

      this.service.dismiss('wizard1');
      this.service.dismiss('wizard2');
      this.service.dismiss('wizard3');

      assert.deepEqual(
        this.service.dismissedWizards,
        ['wizard1', 'wizard2', 'wizard3'],
        'handles multiple dismissals'
      );
      assert.strictEqual(this.setItemStub.callCount, 3, 'calls localStorage.setItem three times');
    });

    test('does not add duplicate wizards', function (assert) {
      this.service.dismissedWizards = ['onboarding'];

      this.service.dismiss('onboarding');

      assert.deepEqual(this.service.dismissedWizards, ['onboarding'], 'does not add duplicate wizard');
      assert.false(this.setItemStub.called, 'does not call localStorage.setItem for duplicate');
    });

    test('preserves existing dismissed wizards', function (assert) {
      this.service.dismissedWizards = ['wizard1', 'wizard2'];

      this.service.dismiss('wizard3');

      assert.deepEqual(
        this.service.dismissedWizards,
        ['wizard1', 'wizard2', 'wizard3'],
        'preserves existing wizards and adds new one'
      );
    });
  });

  module('#reset', function () {
    test('removes specific wizard from dismissed list', function (assert) {
      this.service.dismissedWizards = ['wizard1', 'wizard2', 'wizard3'];

      this.service.reset('wizard2');

      assert.deepEqual(this.service.dismissedWizards, ['wizard1', 'wizard3'], 'removes specified wizard');
      assert.true(
        this.setItemStub.calledWith(DISMISSED_WIZARD_KEY, ['wizard1', 'wizard3']),
        'updates localStorage with remaining wizards'
      );
    });

    test('handles multiple resets', function (assert) {
      this.service.dismissedWizards = ['wizard1', 'wizard2', 'wizard3'];

      this.service.reset('wizard1');
      this.service.reset('wizard3');

      assert.deepEqual(this.service.dismissedWizards, ['wizard2'], 'removes multiple wizards');
      assert.strictEqual(this.setItemStub.callCount, 2, 'calls localStorage.setItem twice');
    });

    test('does nothing when wizard is not in list', function (assert) {
      this.service.dismissedWizards = ['wizard1'];

      this.service.reset('wizard2');

      assert.deepEqual(this.service.dismissedWizards, ['wizard1'], 'list unchanged when wizard not found');
      assert.true(
        this.setItemStub.calledWith(DISMISSED_WIZARD_KEY, ['wizard1']),
        'still updates localStorage'
      );
    });

    test('handles resetting from empty list', function (assert) {
      this.service.dismissedWizards = [];

      this.service.reset('wizard1');

      assert.deepEqual(this.service.dismissedWizards, [], 'list remains empty');
      assert.true(
        this.setItemStub.calledWith(DISMISSED_WIZARD_KEY, []),
        'updates localStorage with empty array'
      );
    });
  });

  module('#resetAll', function () {
    test('clears all dismissed wizards', function (assert) {
      this.service.dismissedWizards = ['wizard1', 'wizard2', 'wizard3'];

      this.service.resetAll();

      assert.deepEqual(this.service.dismissedWizards, [], 'clears all dismissed wizards');
      assert.true(this.removeItemStub.calledWith(DISMISSED_WIZARD_KEY), 'removes key from localStorage');
    });

    test('handles resetAll when list is already empty', function (assert) {
      this.service.dismissedWizards = [];

      this.service.resetAll();

      assert.deepEqual(this.service.dismissedWizards, [], 'list remains empty');
      assert.true(this.removeItemStub.calledWith(DISMISSED_WIZARD_KEY), 'removes key from localStorage');
    });
  });

  module('#isIntroVisible', function () {
    test('returns true by default when wizard not dismissed and intro visibility not set', function (assert) {
      this.service.dismissedWizards = [];

      assert.true(
        this.service.isIntroVisible('onboarding'),
        'returns true for non-dismissed wizard with unset visibility'
      );
    });

    test('returns false by default when wizard is dismissed', function (assert) {
      this.service.dismissedWizards = ['onboarding'];

      assert.false(this.service.isIntroVisible('onboarding'), 'returns false for dismissed wizard');
    });

    test('returns true when intro visibility is set to true', function (assert) {
      this.service.introVisibleState = { onboarding: true };
      assert.true(this.service.isIntroVisible('onboarding'), 'returns true when set');
    });

    test('returns false when intro visibility is set to false', function (assert) {
      this.service.introVisibleState = { onboarding: false };

      assert.false(this.service.isIntroVisible('onboarding'), 'returns false when set to false');
    });
  });

  module('#setIntroVisible', function () {
    test('sets intro visibility to false', function (assert) {
      this.service.setIntroVisible('onboarding', false);

      assert.false(this.service.introVisibleState.onboarding, 'sets intro visibility to false in state');
      assert.false(this.service.isIntroVisible('onboarding'), 'isIntroVisible returns false');
    });

    test('sets intro visibility to true', function (assert) {
      this.service.setIntroVisible('onboarding', true);

      assert.true(this.service.introVisibleState.onboarding, 'sets intro visibility to true in state');
      assert.true(this.service.isIntroVisible('onboarding'), 'isIntroVisible returns true');
    });

    test('handles multiple wizards independently', function (assert) {
      this.service.setIntroVisible('wizard1', false);
      this.service.setIntroVisible('wizard2', true);
      this.service.setIntroVisible('wizard3', false);

      assert.false(this.service.isIntroVisible('wizard1'), 'wizard1 is not visible');
      assert.true(this.service.isIntroVisible('wizard2'), 'wizard2 is visible');
      assert.false(this.service.isIntroVisible('wizard3'), 'wizard3 is not visible');
    });
  });

  /* Step state management */
  module('#getState / #updateState', function () {
    test('getState returns empty object for uninitialized wizard', function (assert) {
      assert.deepEqual(this.service.getState('namespace'), {}, 'returns empty object when not set');
    });

    test('updateState sets a single key without mutating other keys', function (assert) {
      this.service.updateState('namespace', 'choice', 'strict');

      assert.strictEqual(this.service.getState('namespace').choice, 'strict', 'updates the specified key');
    });

    test('successive updates accumulate correctly', function (assert) {
      this.service.updateState('namespace', 'choice', 'strict');
      this.service.updateState('namespace', 'data', ['ns1', 'ns2']);

      const state = this.service.getState('namespace');
      assert.strictEqual(state.choice, 'strict', 'first update persists');
      assert.deepEqual(state.data, ['ns1', 'ns2'], 'second update persists');
    });

    test('updating one wizard does not affect another', function (assert) {
      this.service.updateState('namespace', 'choice', 'strict');

      assert.strictEqual(this.service.getState('namespace').choice, 'strict', 'namespace updated');
      assert.strictEqual(this.service.getState('acl-policy').choice, undefined, 'acl-policy unchanged');
    });
  });

  module('#clearWizardState', function () {
    test('resets state to empty object and step to 0', function (assert) {
      this.service.updateState('namespace', 'choice', 'strict');
      this.service.setCurrentStep('namespace', 2);

      this.service.clearWizardState('namespace');

      assert.deepEqual(this.service.getState('namespace'), {}, 'state is empty after clear');
      assert.strictEqual(this.service.getCurrentStep('namespace'), 0, 'step resets to 0');
    });

    test('preserves step configuration after clear', function (assert) {
      this.service.setSteps('namespace', STEPS);

      this.service.clearWizardState('namespace');

      assert.deepEqual(this.service.getSteps('namespace'), STEPS, 'step config is preserved');
    });
  });

  module('#getCurrentStep / #setCurrentStep', function () {
    test('returns 0 for uninitialized wizard', function (assert) {
      assert.strictEqual(this.service.getCurrentStep('namespace'), 0, 'defaults to 0');
    });

    test('setCurrentStep updates the active step', function (assert) {
      this.service.setCurrentStep('namespace', 2);

      assert.strictEqual(this.service.getCurrentStep('namespace'), 2, 'step is updated');
    });

    test('step changes for one wizard do not affect another', function (assert) {
      this.service.setCurrentStep('namespace', 1);

      assert.strictEqual(this.service.getCurrentStep('namespace'), 1, 'namespace at step 1');
      assert.strictEqual(this.service.getCurrentStep('acl-policy'), 0, 'acl-policy still at step 0');
    });
  });

  module('#getSteps / #setSteps', function () {
    test('getSteps returns empty array for uninitialized wizard', function (assert) {
      assert.deepEqual(this.service.getSteps('namespace'), [], 'returns empty array');
    });

    test('setSteps replaces the step configuration', function (assert) {
      this.service.setSteps('namespace', STEPS);

      assert.deepEqual(this.service.getSteps('namespace'), STEPS, 'step config is stored');
    });

    test('setSteps for one wizard does not affect another', function (assert) {
      this.service.setSteps('namespace', STEPS);

      assert.strictEqual(this.service.getSteps('namespace').length, 3, 'namespace has 3 steps');
      assert.deepEqual(this.service.getSteps('acl-policy'), [], 'acl-policy still has no steps');
    });
  });

  module('#isFinalStep', function () {
    test('returns false for uninitialized wizard (no steps)', function (assert) {
      assert.false(this.service.isFinalStep('namespace'), 'returns false when no steps registered');
    });

    test('returns false when not on the last step', function (assert) {
      this.service.setSteps('namespace', STEPS);

      assert.false(this.service.isFinalStep('namespace'), 'returns false on step 0 of 3');

      this.service.setCurrentStep('namespace', 1);
      assert.false(this.service.isFinalStep('namespace'), 'returns false on step 1 of 3');
    });

    test('returns true when on the last step', function (assert) {
      this.service.setSteps('namespace', STEPS);
      this.service.setCurrentStep('namespace', 2);

      assert.true(this.service.isFinalStep('namespace'), 'returns true on final step');
    });

    test('reflects the new final step after setSteps reduces the step count', function (assert) {
      this.service.setSteps('namespace', STEPS);
      this.service.setCurrentStep('namespace', 1);

      // Reduce to 2 steps — step 1 is now the final step
      this.service.setSteps('namespace', [STEPS[0], STEPS[2]]);

      assert.true(
        this.service.isFinalStep('namespace'),
        'correctly identifies new final step after setSteps'
      );
    });
  });
});
