import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find, findAll } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';

let handler = (data, e) => {
  if (e && e.preventDefault) e.preventDefault();
  return;
};

module('Integration | Component | switch', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });
    this.set('handler', sinon.spy(handler));

    await render(hbs`<Switch
      @onChange={{handler}}
      @inputId="thing"
    />`);

    assert.equal(findAll('label')[0].textContent.trim(), '');

    // Template block usage:
    await render(hbs`
      <Switch
        @onChange={{handler}}
        @inputId="thing"
      >
        <span id="test-value" class="has-text-grey">template block text</span>
      </Switch>
    `);
    assert.dom('[data-test-switch-label]').exists('switch label exists');
    assert.equal(find('#test-value').textContent.trim(), 'template block text', 'yielded text renders');
  });

  test('it has the correct classes', async function(assert) {
    this.set('handler', sinon.spy(handler));
    await render(hbs`
      <Switch
        @onChange={{handler}}
        @inputId="thing"
        @round={{true}}
      >
      template block text
      </Switch>
    `);
    assert.dom('.switch.is-rounded').exists('switch has is-rounded class');
    // await pauseTest();
    await render(hbs`
      <Switch
        @onChange={{handler}}
        @inputId="thing"
        @size="small"
      >
        template block text
      </Switch>
    `);
    assert.dom('.switch.is-small').exists('switch has is-small class');

    await render(hbs`
      <Switch
        @onChange={{handler}}
        @inputId="thing"
        @size="large"
      >
        template block text
      </Switch>
    `);
    assert.dom('.switch.is-large').exists('switch has is-large class');
  });
});
