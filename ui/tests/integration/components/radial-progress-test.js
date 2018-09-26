import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import radialProgress from 'vault/tests/pages/components/radial-progress';

const component = create(radialProgress);

module('Integration | Component | radial progress', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  test('it renders', async function(assert) {
    let circumference = (19 / 2) * Math.PI * 2;
    await render(hbs`{{radial-progress progressDecimal=0.5}}`);

    assert.equal(component.viewBox, '0 0 20 20');
    assert.equal(component.height, '20');
    assert.equal(component.width, '20');
    assert.equal(component.strokeWidth, '1');
    assert.equal(component.r, 19 / 2);
    assert.equal(component.cx, 10);
    assert.equal(component.cy, 10);
    assert.equal(component.strokeDash, circumference);
    assert.equal(component.strokeDashOffset, circumference * 0.5);
  });
});
