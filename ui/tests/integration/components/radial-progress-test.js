import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import radialProgress from 'vault/tests/pages/components/radial-progress';

const component = create(radialProgress);

module('Integration | Component | radial progress', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    // We have to manually round the circumference, strokeDash, and strokeDashOffset because
    // ie11 truncates decimals differently than other browsers.
    let circumference = ((19 / 2) * Math.PI * 2).toFixed(2);
    await render(hbs`{{radial-progress progressDecimal=0.5}}`);

    assert.equal(component.viewBox, '0 0 20 20');
    assert.equal(component.height, '20');
    assert.equal(component.width, '20');
    assert.equal(component.strokeWidth, '1');
    assert.equal(component.r, 19 / 2);
    assert.equal(component.cx, 10);
    assert.equal(component.cy, 10);
    assert.equal(Number(component.strokeDash).toFixed(2), circumference);
    assert.equal(Number(component.strokeDashOffset).toFixed(3), (circumference * 0.5).toFixed(3));
  });
});
