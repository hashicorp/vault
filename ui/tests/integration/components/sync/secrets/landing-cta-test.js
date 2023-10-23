/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';

module('Integration | Component | sync | Secrets::LandingCta', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
  });

  test('it should render promotional copy for community version', async function (assert) {
    this.version.version = '1.16.0';
    await render(
      hbs`
     <Secrets::LandingCta />
    `,
      { owner: this.engine }
    );

    assert
      .dom(PAGE.cta.summary)
      .hasText(
        'This enterprise feature allows you to sync secrets to platforms and tools across your stack to get secrets when and where you need them.'
      );
    assert.dom(PAGE.cta.button).hasText('Learn more about secrets sync');
  });

  test('it should render enterprise copy', async function (assert) {
    this.version.version = '1.16.0+ent';
    await render(
      hbs`
     <Secrets::LandingCta />
    `,
      { owner: this.engine }
    );

    assert
      .dom(PAGE.cta.summary)
      .hasText(
        'Sync secrets to platforms and tools across your stack to get secrets when and where you need them. Secrets sync tutorial'
      );
    assert.dom(PAGE.cta.button).hasText('Create first destination');
  });
});
