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
  wait().then(() => {
    assert.equal(component.consoleInputValue, "list this/thing/here", 'populates console input with previous command on up after enter');
  });
  wait().then(() => component.down());
  return wait().then(() => {
    assert.equal(component.consoleInputValue, "", 'populates console input with next command on down');
  });
});

test('it cycles through history with more than one command', function(assert) {

  this.render(hbs`{{console/ui-panel}}`);

  component.consoleInput('list this/thing/here').enter();
  wait().then(() => component.consoleInput('read that/thing/there').enter());
  wait().then(() => component.consoleInput('qwerty').enter());

  wait().then(() => component.up());
  wait().then(() => {
    assert.equal(component.consoleInputValue, "qwerty", 'populates console input with previous command on up after enter');
  });
  wait().then(() => component.up());
  wait().then(() => {
    assert.equal(component.consoleInputValue, "read that/thing/there", 'populates console input with previous command on up');
  });
  wait().then(() => component.up());
  wait().then(() => {
    assert.equal(component.consoleInputValue, "list this/thing/here", 'populates console input with previous command on up');
  });
  wait().then(() => component.up());
  wait().then(() => {
    assert.equal(component.consoleInputValue, "qwerty", 'populates console input with initial command if cycled through all previous commands');
  });
  wait().then(() => component.down());
  return wait().then(() => {
    assert.equal(component.consoleInputValue, "", 'clears console input if down pressed after history is on most recent command');
  });
});