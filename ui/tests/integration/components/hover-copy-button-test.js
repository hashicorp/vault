import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import copyButton from 'vault/tests/pages/components/hover-copy-button';
import { triggerSuccess } from '../../helpers/ember-cli-clipboard';
const component = create(copyButton);

module('Integration | Component | hover copy button', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  test('it shows success message in tooltip', async function(assert) {
    this.set('copyValue', 'foo');
    await render(
      hbs`<div class="has-copy-button" tabindex="-1">{{hover-copy-button copyValue=copyValue}}</div>`
    );

    component.focusContainer();
    assert.ok(component.buttonIsVisible);
    component.mouseEnter();
    assert.equal(component.tooltipText, 'Copy', 'shows copy');
    triggerSuccess(this, '[data-test-hover-copy-button]');
    assert.equal(component.tooltipText, 'Copied!', 'shows success message');
  });

  test('it has the correct class when alwaysShow is true', async function(assert) {
    this.set('copyValue', 'foo');
    await render(hbs`{{hover-copy-button alwaysShow=true copyValue=copyValue}}`);
    assert.ok(component.buttonIsVisible);
    assert.ok(component.wrapperClass.includes('hover-copy-button-static'));
  });
});
