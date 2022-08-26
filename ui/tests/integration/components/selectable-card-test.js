import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const TOTAL = 15;
const CARD_TITLE = 'Connections';

module('Integration | Component selectable-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('total', TOTAL);
    this.set('cardTitle', CARD_TITLE);
  });

  test('it shows the card total', async function (assert) {
    await render(hbs`<SelectableCard @total={{total}} @cardTitle={{cardTitle}}/>`);
    let titleNumber = this.element.querySelector('.title-number').innerText;

    assert.equal(titleNumber, 15);
  });

  test('it returns card title, ', async function (assert) {
    await render(hbs`<SelectableCard @total={{1}} @cardTitle={{cardTitle}}/>`);
    let titleText = this.element.querySelector('.title').innerText;
    assert.equal(titleText, 'Connections');
  });
});
