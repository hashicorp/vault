import { module, skip, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import uiPanel from 'vault/tests/pages/components/console/ui-panel';
import hbs from 'htmlbars-inline-precompile';

const component = create(uiPanel);

module('Integration | Component | console/ui panel', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  test('it renders', async function(assert) {
    await render(hbs`{{console/ui-panel}}`);
    assert.ok(component.hasInput);
  });

  test('it clears console input on enter', async function(assert) {
    await render(hbs`{{console/ui-panel}}`);
    component.consoleInput('list this/thing/here').enter();
    return settled().then(() => {
      assert.equal(component.consoleInputValue, '', 'empties input field on enter');
    });
  });

  test('it clears the log when using clear command', async function(assert) {
    await render(hbs`{{console/ui-panel}}`);
    component.consoleInput('list this/thing/here').enter();
    component.consoleInput('list this/other/thing').enter();
    component.consoleInput('read another/thing').enter();
    settled().then(() => {
      assert.notEqual(component.logOutput, '', 'there is output in the log');
      component.consoleInput('clear').enter();
    });

    settled().then(() => component.up());
    return settled().then(() => {
      assert.equal(component.logOutput, '', 'clears the output log');
      assert.equal(
        component.consoleInputValue,
        'clear',
        'populates console input with previous command on up after enter'
      );
    });
  });

  test('it adds command to history on enter', async function(assert) {
    await render(hbs`{{console/ui-panel}}`);
    component.consoleInput('list this/thing/here').enter();
    settled().then(() => component.up());
    settled().then(() => {
      assert.equal(
        component.consoleInputValue,
        'list this/thing/here',
        'populates console input with previous command on up after enter'
      );
    });
    settled().then(() => component.down());
    return settled().then(() => {
      assert.equal(component.consoleInputValue, '', 'populates console input with next command on down');
    });
  });

  skip('it cycles through history with more than one command', function(assert) {
    this.render(hbs`{{console/ui-panel}}`);
    component.consoleInput('list this/thing/here').enter();
    settled().then(() => component.consoleInput('read that/thing/there').enter());
    settled().then(() => component.consoleInput('qwerty').enter());

    settled().then(() => component.up());
    settled().then(() => {
      assert.equal(
        component.consoleInputValue,
        'qwerty',
        'populates console input with previous command on up after enter'
      );
    });
    settled().then(() => component.up());
    settled().then(() => {
      assert.equal(
        component.consoleInputValue,
        'read that/thing/there',
        'populates console input with previous command on up'
      );
    });
    settled().then(() => component.up());
    settled().then(() => {
      assert.equal(
        component.consoleInputValue,
        'list this/thing/here',
        'populates console input with previous command on up'
      );
    });
    settled().then(() => component.up());
    settled().then(() => {
      assert.equal(
        component.consoleInputValue,
        'qwerty',
        'populates console input with initial command if cycled through all previous commands'
      );
    });
    settled().then(() => component.down());
    return settled().then(() => {
      assert.equal(
        component.consoleInputValue,
        '',
        'clears console input if down pressed after history is on most recent command'
      );
    });
  });
});
