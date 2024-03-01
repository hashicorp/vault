/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import hbs from 'htmlbars-inline-precompile';
import { click, render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import sinon from 'sinon';

const { cta } = PAGE;

module('Integration | Component | sync | Secrets::LandingCta', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.featureEnabled = true;

    this.renderComponent = () =>
      render(
        hbs`
          <Secrets::LandingCta @featureEnabled={{this.featureEnabled}}/>
        `,
        { owner: this.engine }
      );
  });

  test('it should render promotional copy for community version', async function (assert) {
    await this.renderComponent();

    assert
      .dom(cta.summary)
      .hasText(
        'This enterprise feature allows you to sync secrets to platforms and tools across your stack to get secrets when and where you need them. Learn more about secrets sync'
      );
    assert.dom(cta.link).hasText('Learn more about secrets sync');
  });

  test('it should render enterprise copy and action', async function (assert) {
    this.version.type = 'enterprise';

    await this.renderComponent();

    assert
      .dom(cta.summary)
      .hasText(
        'Sync secrets to platforms and tools across your stack to get secrets when and where you need them. Secrets sync tutorial'
      );
    assert.dom(cta.link).hasText('Secrets sync tutorial');

    await click(cta.button);

    assert.propEqual(
      this.transitionStub.lastCall.args,
      ['vault.cluster.sync.secrets.destinations.create'],
      'Transitions to destinations create route'
    );
  });

  test('it should render enable feature modal', async function (assert) {
    assert.expect(3);

    this.version.type = 'enterprise';
    this.featureEnabled = false;

    this.server.patch('/sys/sync/config', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.false(payload.disabled, 'Correct value sent to enable feature in config request');
    });

    await this.renderComponent();

    await click(cta.button);
    assert.dom(cta.featureConfirm).isDisabled('Confirm button is disabled when checkbox is unchecked');
    await click(cta.featureCheck);
    await click(cta.featureConfirm);
    assert.propEqual(
      this.transitionStub.lastCall.args,
      ['vault.cluster.sync.secrets.destinations.create'],
      'Transitions to destinations create route after enable feature success'
    );
  });
});
