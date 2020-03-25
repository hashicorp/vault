import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find, findAll } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';

let handler = (data, e) => {
  if (e && e.preventDefault) e.preventDefault();
  return;
};

module('Integration | Component | toggle', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    this.set('handler', sinon.spy(handler));

    await render(hbs`<Toggle
      @onChange={{handler}}
      @name="thing"
    />`);

    assert.equal(findAll('label')[0].textContent.trim(), '');

    await render(hbs`
      <Toggle
        @onChange={{handler}}
        @name="thing"
      >
        <span id="test-value" class="has-text-grey">template block text</span>
      </Toggle>
    `);
    assert.dom('[data-test-toggle-label="thing"]').exists('toggle label exists');
    assert.equal(find('#test-value').textContent.trim(), 'template block text', 'yielded text renders');
  });

  test('it has the correct classes', async function(assert) {
    this.set('handler', sinon.spy(handler));
    await render(hbs`
      <Toggle
        @onChange={{handler}}
        @name="thing"
        @size="small"
      >
        template block text
      </Toggle>
    `);
    assert.dom('.toggle.is-small').exists('toggle has is-small class');
  });

  test('it sets the id of the input correctly', async function(assert) {
    this.set('handler', sinon.spy(handler));
    await render(hbs`
    <Toggle
      @onChange={{handler}}
      @name="my toggle"
    >
      Label
    </Toggle>
    `);
    assert.dom('#toggle-mytoggle').exists('input has correct ID');
    assert.dom('label').hasAttribute('for', 'toggle-mytoggle', 'label has correct for attribute');
  });
});
