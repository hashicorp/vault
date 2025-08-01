/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Sinon from 'sinon';

module('Integration | Helper | transition-to', function (hooks) {
  setupRenderingTest(hooks);
  // using 'kv' here for testing, but this could be any Ember engine in the app
  // sets this.engine, which we use to set context for the component testing service:app-router
  setupEngine(hooks, 'kv');

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.router.reopen({
      transitionTo: Sinon.stub(),
      transitionToExternal: Sinon.stub(),
    });
  });

  test('it does not call transition on render', async function (assert) {
    await render(hbs`<button data-test-btn {{on "click" (transition-to "vault.cluster")}}>Click</button>`);

    assert.true(this.router.transitionTo.notCalled, 'transitionTo not called on render');
    assert.true(this.router.transitionToExternal.notCalled, 'transitionToExternal not called on render');
  });

  test('it calls transitionTo correctly', async function (assert) {
    await render(
      hbs`<button data-test-btn {{on "click" (transition-to "vault.cluster" "foobar" "baz")}}>Click</button>`
    );
    await click('[data-test-btn]');

    assert.true(this.router.transitionTo.calledOnce, 'transitionTo called once on click');
    assert.deepEqual(
      this.router.transitionTo.args[0],
      ['vault.cluster', 'foobar', 'baz'],
      'transitionTo called with positional params'
    );
    assert.true(this.router.transitionToExternal.notCalled, 'transitionToExternal not called');
  });

  test('it calls transitionToExternal correctly', async function (assert) {
    await render(
      hbs`<button data-test-btn {{on "click" (transition-to "vault.cluster" "foobar" "baz" external=true)}}>Click</button>`
    );
    await click('[data-test-btn]');

    assert.true(this.router.transitionToExternal.calledOnce, 'transitionToExternal called');
    assert.deepEqual(
      this.router.transitionToExternal.args[0],
      ['vault.cluster', 'foobar', 'baz'],
      'transitionToExternal called with positional params'
    );
    assert.true(this.router.transitionTo.notCalled, 'transitionTo not called');
  });

  // This test is confusing (and admittedly not ideal) because stubbing routers gets strange,
  // but if you go into the TransitionTo class and console.log owner.lookup('service:router') in get router()
  // you'll see the getter returns 'service:app-router' (because of the context setup)
  // so although we're asserting this.router, the TransitionTo helper is using "service:app-router" under the hood.
  // This test passing, indirectly means the helper works as expected. Failures might be something like "global failure: TypeError: this.router is undefined"
  test('it uses service:app-router when base router undefined', async function (assert) {
    await render(
      hbs`<button data-test-btn {{on "click" (transition-to "vault.cluster" "foobar" "baz" external=true)}}>Click</button>`,
      { owner: this.engine }
    );
    await click('[data-test-btn]');
    assert.true(this.router.transitionToExternal.calledOnce, 'transitionToExternal called');
  });
});
