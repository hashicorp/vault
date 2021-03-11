import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import sinon from 'sinon';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | box-radio', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('type', 'aws');
    this.set('displayName', 'An Option');
    this.set('mountType', '');
    this.set('disabled', false);
  });

  test('it renders and triggers on radio select', async function(assert) {
    const spy = sinon.spy();
    this.set('onRadioChange', spy);
    await render(hbs`<BoxRadio
      @key={{type}}
      @glyph={{type}}
      @displayName={{displayName}}
      @onRadioChange={{onRadioChange}}
      @disabled={{disabled}}
    />`);

    assert.dom(this.element).hasText('An Option', 'shows the display name of the option');
    assert.dom('.tooltip').doesNotExist('tooltip does not exist when disabled is false');
    await click('[data-test-mount-key="aws"]');
    assert.ok(spy.calledWith('aws'), 'calls the radio change function when option clicked');
  });

  test('it renders and triggers on box radio click', async function(assert) {
    const spy = sinon.spy();
    this.set('onRadioChange', spy);
    await render(hbs`<BoxRadio
      @key={{type}}
      @glyph={{type}}
      @displayName={{displayName}}
      @onRadioChange={{onRadioChange}}
      @disabled={{disabled}}
    />`);

    assert.dom(this.element).hasText('An Option', 'shows the display name of the option');
    assert.dom('.tooltip').doesNotExist('tooltip does not exist when disabled is false');
    await click('[data-test-box-radio-input="aws"]');
    assert.ok(spy.calledWith('aws'), 'calls the radio change function when radio selected');
  });

  test('it renders correctly when disabled', async function(assert) {
    const spy = sinon.spy();
    this.set('onRadioChange', spy);
    await render(hbs`<BoxRadio
      @key={{type}}
      @glyph={{type}}
      @displayName={{displayName}}
      @onRadioChange={{onRadioChange}}
      @disabled=true
    />`);

    assert.dom(this.element).hasText('An Option', 'shows the display name of the option');
    assert.dom('.ember-basic-dropdown-trigger').exists('tooltip exists');
    await click('[data-test-mount-type="aws"]');
    assert.ok(spy.notCalled, 'does not call the radio change function when option is clicked');
  });
});
