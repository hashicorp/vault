import { moduleForComponent, test } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import wait from 'ember-test-helpers/wait';
import uiPanel from 'vault/tests/pages/components/console/ui-panel';
import hbs from 'htmlbars-inline-precompile';

const component = create(uiPanel);

moduleForComponent('console/ui-panel', 'Integration | Component | console/ui panel', {
  integration: true,

  beforeEach(){
    component.setContext(this);
  },

  afterEach(){
    component.removeContext();
  },
});

test('it renders', function(assert) {

  this.render(hbs`{{console/ui-panel}}`);

  assert.ok(component.hasInput);

});

test('it clears console input on enter', function(assert) {

  this.render(hbs`{{console/ui-panel}}`);

  component.consoleInput('list this/thing/here').enter();

  return wait().then(() => {
    assert.equal(component.consoleInputValue, "", 'empties input field on enter');
  });
});

test('it adds command to history on enter', function(assert) {

  this.render(hbs`{{console/ui-panel}}`);

  component.consoleInput('list this/thing/here').enter();
  wait().then(() => component.up());
  return wait().then(() => {
    assert.equal(component.consoleInputValue, "list this/thing/here", 'populates console input with previous command on up after enter');
  });
});