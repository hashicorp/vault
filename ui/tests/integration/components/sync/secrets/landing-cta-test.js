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
    this.version = this.owner.lookup('service:version');
    this.persona = this.owner.lookup('service:persona');
    this.version.type = 'enterprise';
    this.version.features = [];
    this.version.activatedFeatures = [];

    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
  });

  test('it should render promotional copy for community', async function (assert) {
    this.version.type = 'community';
    this.secretsSyncPersona = this.persona.secretsSyncPersona;
    await render(hbs`<Secrets::LandingCta @secretsSyncPersona={{this.secretsSyncPersona}} /> `, {
      owner: this.engine,
    });

    assert
      .dom(cta.summary)
      .hasText(
        'This premium enterprise feature allows you to sync secrets to platforms and tools across your stack to get secrets when and where you need them. Learn more about Secrets Sync'
      );
    assert.dom(cta.link).hasText('Learn more about Secrets Sync');
    assert.dom(cta.button).doesNotExist('does not render create destination button');
  });

  test('it should render promotional copy for enterprise version without feature', async function (assert) {
    this.secretsSyncPersona = this.persona.secretsSyncPersona;

    await render(hbs`<Secrets::LandingCta @secretsSyncPersona={{this.secretsSyncPersona}} /> `, {
      owner: this.engine,
    });

    assert
      .dom(cta.summary)
      .hasText(
        'This premium enterprise feature allows you to sync secrets to platforms and tools across your stack to get secrets when and where you need them. Learn more about Secrets Sync'
      );
    assert.dom(cta.link).hasText('Learn more about Secrets Sync');
    assert.dom(cta.button).doesNotExist('does not render create destination button');
  });

  test('it should render CTA copy but not action when feature exists on enterprise license and is not activated', async function (assert) {
    this.version.features = ['Secrets Sync'];
    this.secretsSyncPersona = this.persona.secretsSyncPersona;
    await render(hbs`<Secrets::LandingCta @secretsSyncPersona={{this.secretsSyncPersona}} /> `, {
      owner: this.engine,
    });
    assert
      .dom(cta.summary)
      .hasText(
        'Sync secrets to platforms and tools across your stack to get secrets when and where you need them. Secrets Sync tutorial'
      );
    assert.dom(cta.link).hasText('Secrets Sync tutorial');
    assert.dom(cta.button).doesNotExist('does not render create destination button');
  });

  test('it should render CTA copy and action when feature exists on enterprise license and is activated', async function (assert) {
    this.version.features = ['Secrets Sync'];
    this.version.activatedFeatures = ['secrets-sync'];
    this.secretsSyncPersona = this.persona.secretsSyncPersona;
    await render(hbs`<Secrets::LandingCta @secretsSyncPersona={{this.secretsSyncPersona}} /> `, {
      owner: this.engine,
    });

    assert
      .dom(cta.summary)
      .hasText(
        'Sync secrets to platforms and tools across your stack to get secrets when and where you need them. Secrets Sync tutorial'
      );
    assert.dom(cta.link).hasText('Secrets Sync tutorial');
    assert.dom(cta.button).exists('it renders create destination button');
  });
});
