import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import copyButton from 'vault/tests/pages/components/hover-copy-button';
const component = create(copyButton);

module('Integration | Component | hover copy button', function(hooks) {
  setupRenderingTest(hooks);

  // ember-cli-clipboard helpers don't like the new style
  test('it shows success message in tooltip', async function(assert) {
    this.set('copyValue', 'foo');
    await render(
      hbs`<div class="has-copy-button" tabindex="-1">
      <HoverCopyButton @copyValue={{copyValue}} />
      </div>`
    );
    await settled();
    await component.focusContainer();
    await settled();
    assert.ok(component.buttonIsVisible);
    await component.mouseEnter();
    await settled();
    assert.equal(component.tooltipText, 'Copy', 'shows copy');
  });

  test('it has the correct class when alwaysShow is true', async function(assert) {
    this.set('copyValue', 'foo');
    await render(hbs`{{hover-copy-button alwaysShow=true copyValue=copyValue}}`);
    assert.ok(component.buttonIsVisible);
    assert.ok(component.wrapperClass.includes('hover-copy-button-static'));
  });
});
