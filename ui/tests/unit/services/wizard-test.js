/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import localStorage from 'vault/lib/local-storage';

const DISMISSED_WIZARD_KEY = 'dismissed-wizards';

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
});
