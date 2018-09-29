import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, blur, find, triggerEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import waitForError from 'vault/tests/helpers/wait-for-error';

module('Integration | Component | wrap ttl', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.lastOnChangeCall = null;
    this.set('onChange', val => {
      this.lastOnChangeCall = val;
    });
  });

  test('it requires `onChange`', async function(assert) {
    let promise = waitForError();
    render(hbs`{{wrap-ttl}}`);
    let err = await promise;
    assert.ok(err.message.includes('`onChange` handler is a required attr in'), 'asserts without onChange');
  });

  test('it renders', async function(assert) {
    await render(hbs`{{wrap-ttl onChange=(action onChange)}}`);
    assert.equal(this.lastOnChangeCall, '30m', 'calls onChange with 30m default on first render');
    assert.equal(find('label[for=wrap-response]').textContent.trim(), 'Wrap response');
  });

  test('it nulls out value when you uncheck wrapResponse', async function(assert) {
    await render(hbs`{{wrap-ttl onChange=(action onChange)}}`);
    await click('#wrap-response');
    await triggerEvent('#wrap-response', 'change');
    assert.equal(this.lastOnChangeCall, null, 'calls onChange with null');
  });

  test('it sends value changes to onChange handler', async function(assert) {
    await render(hbs`{{wrap-ttl onChange=(action onChange)}}`);

    await fillIn('[data-test-wrap-ttl-picker] input', '20');
    assert.equal(this.lastOnChangeCall, '20m', 'calls onChange correctly on time input');

    await fillIn('#unit', 'h');
    await blur('#unit');
    assert.equal(this.lastOnChangeCall, '20h', 'calls onChange correctly on unit change');

    await fillIn('#unit', 'd');
    await blur('#unit');
    assert.equal(this.lastOnChangeCall, '480h', 'converts days to hours correctly');
  });
});
