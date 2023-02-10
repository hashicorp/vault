import { module, test } from 'qunit';
import Sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, typeIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import waitForError from 'vault/tests/helpers/wait-for-error';

module('Integration | Component | wrap ttl', function (hooks) {
  setupRenderingTest(hooks);

  test('it requires `onChange`', async function (assert) {
    const promise = waitForError();
    render(hbs`<WrapTtl />`);
    const err = await promise;
    assert.ok(err.message.includes('`onChange` handler is a required attr in'), 'asserts without onChange');
  });

  test('it renders', async function (assert) {
    const changeSpy = Sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`<WrapTtl @onChange={{this.onChange}} />`);
    assert.ok(changeSpy.calledWithExactly('30m'), 'calls onChange with 30m default on render');
    assert.dom('[data-test-ttl-form-label]').hasText('Wrap response');
  });

  test('it nulls out value when you uncheck wrapResponse', async function (assert) {
    const changeSpy = Sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`<WrapTtl @onChange={{this.onChange}} />`);
    await click('[data-test-ttl-form-label]');
    assert.ok(changeSpy.calledWithExactly(null), 'calls onChange with null');
  });

  test('it sends value changes to onChange handler', async function (assert) {
    const changeSpy = Sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`<WrapTtl @onChange={{this.onChange}} />`);
    // for testing purposes we need to input unit first because it keeps seconds value
    await typeIn('[data-test-select="ttl-unit"]', 'h');
    assert.ok(changeSpy.calledWithExactly('30h'), 'calls onChange correctly on time input');
    await typeIn('[data-test-ttl-value]', '20');
    assert.ok(changeSpy.calledWithExactly('20h'), 'calls onChange correctly on unit change');
  });
});
