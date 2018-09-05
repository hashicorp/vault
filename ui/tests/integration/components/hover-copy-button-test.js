import { moduleForComponent, test } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import copyButton from 'vault/tests/pages/components/hover-copy-button';
import { triggerSuccess } from '../../helpers/ember-cli-clipboard';
const component = create(copyButton);

moduleForComponent('hover-copy-button', 'Integration | Component | hover copy button', {
  integration: true,

  beforeEach() {
    component.setContext(this);
  },

  afterEach() {
    component.removeContext();
  },
});

test('it shows success message in tooltip', function(assert) {
  this.set('copyValue', 'foo');
  this.render(
    hbs`<div class="has-copy-button" tabindex="-1">{{hover-copy-button copyValue=copyValue}}</div>`
  );

  component.focusContainer();
  assert.ok(component.buttonIsVisible);
  component.mouseEnter();
  assert.equal(component.tooltipText, 'Copy', 'shows copy');
  triggerSuccess(this, '[data-test-hover-copy-button]');
  assert.equal(component.tooltipText, 'Copied!', 'shows success message');
});

test('it has the correct class when alwaysShow is true', function(assert) {
  this.set('copyValue', 'foo');
  this.render(hbs`{{hover-copy-button alwaysShow=true copyValue=copyValue}}`);
  assert.ok(component.buttonIsVisible);
  assert.ok(component.wrapperClass.includes('hover-copy-button-static'));
});
