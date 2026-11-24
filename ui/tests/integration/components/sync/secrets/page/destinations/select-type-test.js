/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { click, render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { destinationTypes, syncDestinations } from 'vault/helpers/sync-destinations';

module('Integration | Component | sync | Secrets::Page::Destinations::SelectType', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');

  test('it transitions to selected type', async function (assert) {
    const types = destinationTypes();
    assert.expect(types.length);
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await render(
      hbs`
     <Secrets::Page::Destinations::SelectType />
    `,
      { owner: this.engine }
    );

    for (const type of types) {
      await click(PAGE.selectType(type));
      const transition = transitionStub.calledWith(
        'vault.cluster.sync.secrets.destinations.create.destination',
        type
      );
      assert.true(transition, `transitionTo called with param: ${type}`);
    }
  });

  module('dark mode icon resolution', function (hooks) {
    hooks.beforeEach(function () {
      this.theme = this.owner.lookup('service:theme');
    });

    hooks.afterEach(function () {
      this.theme.setTheme('light');
    });

    test('icons have -color suffix stripped in dark mode', async function (assert) {
      this.theme.setTheme('dark');
      await render(hbs`<Secrets::Page::Destinations::SelectType />`, { owner: this.engine });

      for (const destination of syncDestinations()) {
        const expectedIcon = destination.icon.replace(/-color$/, '');
        assert
          .dom(`${PAGE.selectType(destination.type)} [data-test-icon="${expectedIcon}"]`)
          .exists(`${destination.type} renders monochrome icon "${expectedIcon}" in dark mode`);
      }
    });

    test('icons keep -color suffix in light mode', async function (assert) {
      this.theme.setTheme('light');
      await render(hbs`<Secrets::Page::Destinations::SelectType />`, { owner: this.engine });

      for (const destination of syncDestinations()) {
        assert
          .dom(`${PAGE.selectType(destination.type)} [data-test-icon="${destination.icon}"]`)
          .exists(`${destination.type} renders colour icon "${destination.icon}" in light mode`);
      }
    });
  });
});
