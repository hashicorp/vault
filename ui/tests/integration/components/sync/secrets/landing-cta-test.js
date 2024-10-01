/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import sinon from 'sinon';

const { cta } = PAGE;

module('Integration | Component | sync | Secrets::LandingCta', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
  });

  test('it should render promotional copy if feature is not activated', async function (assert) {
    await render(hbs`<Secrets::LandingCta @isActivated={{false}} /> `, {
      owner: this.engine,
    });

    assert
      .dom(cta.summary)
      .hasText(
        'This feature allows you to sync secrets to platforms and tools across your stack to get secrets when and where you need them. Learn more about Secrets Sync'
      );
    assert.dom(cta.link).hasText('Learn more about Secrets Sync');
    assert.dom(cta.button).doesNotExist('does not render create destination button');
  });

  test('it should render CTA copy and action if feature is activated', async function (assert) {
    await render(hbs`<Secrets::LandingCta @isActivated={{true}} /> `, {
      owner: this.engine,
    });

    assert
      .dom(cta.summary)
      .hasText(
        'This feature allows you to sync secrets to platforms and tools across your stack to get secrets when and where you need them. Learn more about Secrets Sync'
      );
    assert.dom(cta.link).hasText('Learn more about Secrets Sync');
    assert.dom(cta.button).exists('it renders create destination button');
  });
});
