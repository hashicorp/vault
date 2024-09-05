import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Service from '@ember/service';
import Sinon from 'sinon';
import { resolve } from 'rsvp';

class FakeRouter extends Service {
  transitionTo = Sinon.stub().returns(resolve());
  transitionToExternal = Sinon.stub().returns(resolve());
}
module('Integration | Helper | engine-transition-to', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.owner.unregister('service:app-router');
    this.owner.register('service:app-router', FakeRouter);
    this.router = this.owner.lookup('service:app-router');
  });

  test('it does not call transition on render', async function (assert) {
    await render(
      hbs`<button data-test-btn {{on "click" (engine-transition-to "vault.cluster")}}>Click</button>`
    );
    assert.true(this.router.transitionTo.notCalled, 'transitionTo not called on render');
    assert.true(this.router.transitionToExternal.notCalled, 'transitionToExternal not called on render');
  });

  test('it calls transitionTo correctly', async function (assert) {
    await render(
      hbs`<button data-test-btn {{on "click" (engine-transition-to "vault.cluster" "foobar" "baz")}}>Click</button>`
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
      hbs`<button data-test-btn {{on "click" (engine-transition-to "vault.cluster" "foobar" "baz" external=true)}}>Click</button>`
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
});
