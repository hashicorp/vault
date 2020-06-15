import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const TITLE = 'My Table';
const HEADER = 'Cool Header';
const ITEMS = ['https://127.0.0.1:8201', 'hello', 3];

module('Integration | Enterprise | Component | InfoTable', function(hooks) {
  // uses setupTests() under the hood and gives us Ember's dependency injection system
  // allowing us to look up anything in the application, such as this.owner
  // and grants us this.element so we can make assertions
  setupRenderingTest(hooks);

  // setupRenderingTest also enables us to get and set items in our test context
  // defining the consts above does not grant our templates access to them until we
  // manually set these values below
  hooks.beforeEach(function() {
    this.set('title', TITLE);
    this.set('header', HEADER);
    this.set('items', ITEMS);
  });

  test('it renders', async function(assert) {
    await render(hbs`<InfoTable
        @title={{title}}
        @header={{header}}
        @items={{items}}
      />`);

    // you can use ember-test-selectors to find an element in the dom
    assert.dom('[data-test-info-table]').exists();

    // this.element returns the resulting element from render(hbs) above
    // so you can also search within its markup to make any assertions from QUnit DOM
    assert.dom(this.element).includesText(HEADER, `shows the table header`);

    const rows = this.element.querySelectorAll('.info-table-row');
    assert.equal(rows.length, ITEMS.length, 'renders an InfoTableRow for each item');

    rows.forEach((row, i) => {
      assert.equal(row.innerText, ITEMS[i], 'handles strings and numbers as row values');
    });
  });
});
