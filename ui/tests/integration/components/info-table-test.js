import { module, test, todo } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

// We'll need to set these variables in our test context as well
const TITLE = 'My Table';
const HEADER = 'Cool Header';
const ITEMS = ['https://127.0.0.1:8201', 'hello', 3];

module('Integration | Enterprise | Component | InfoTable | ember learn', function(hooks) {
  /*
  EXAMPLE 1

  setupRenderingTest gives us access to Ember's dependency injection
  system, allowing us to look up anything in the application, such
  as this.owner. It also grants us this.element so we can make
  assertions based on the markup that specific element (as opposed
  to the whole DOM).
  */
  setupRenderingTest(hooks);

  /*
  setupRenderingTest also enables us to get and set items in our
  test context. Defining the consts above does not grant our
  templates access to them until we manually set these values inside the hooks below.
  */
  hooks.beforeEach(function() {
    this.set('title', TITLE);
    this.set('header', HEADER);
    this.set('items', ITEMS);
  });

  todo('it renders', async function(assert) {
    await render(hbs`<InfoTable
        @title={{title}}
        @header={{header}}
        @items={{items}}
      />`);

    /*
    The next few assertions utilize ember-test-selectors and javascript's querySelectorAll
    */
    assert.dom('[data-test-info-table]').exists();

    /*
    this.element returns the resulting element from render(hbs). This allows you to search within its markup to make any
    assertions from QUnit DOM.
    */
    assert.dom(this.element).includesText(HEADER, `shows the table header`);

    /*
    Notice you can use either a class, id, or data-test-* attribute
    here. In this scenario the class name works well because there'
    nothing special about the info-table-row itself, but it does
    mean if our class names ever changed this test would fail.
    */
    const rows = this.element.querySelectorAll('.info-table-row');
    assert.equal(rows.length, ITEMS.length, 'renders an InfoTableRow for each item');

    rows.forEach((row, i) => {
      assert.equal(row.innerText, ITEMS[i], 'handles strings and numbers as row values');
    });
  });
});
