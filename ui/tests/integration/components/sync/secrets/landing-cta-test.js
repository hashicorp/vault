/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const { cta } = PAGE;

module('Integration | Component | sync | Secrets::LandingCta', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
  });

  test('it should render CTA if feature is activated', async function (assert) {
    await render(hbs`<Secrets::LandingCta @isActivated={{true}} /> `, {
      owner: this.engine,
    });

    assert.dom(cta.button).exists('it renders create destination button');
    assert
      .dom(GENERAL.textBody('sync-step-2'))
      .hasText(
        'Step 2: Select secrets from Vault and sync secrets to your destination.',
        'Sync instructions are not limited to a specific secrets engine'
      );
  });

  test('it shows engine-neutral sync instructions before activation', async function (assert) {
    await render(hbs`<Secrets::LandingCta @isActivated={{false}} />`, {
      owner: this.engine,
    });

    assert
      .dom(GENERAL.textBody('sync-step-2'))
      .hasText(
        'Step 2: Select secrets from Vault and sync secrets to your destination.',
        'Sync instructions are not limited to a specific secrets engine'
      );
  });
});
