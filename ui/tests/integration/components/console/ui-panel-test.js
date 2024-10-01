/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import uiPanel from 'vault/tests/pages/components/console/ui-panel';
import hbs from 'htmlbars-inline-precompile';

const component = create(uiPanel);

module('Integration | Component | console/ui panel', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`{{console/ui-panel}}`);

    assert.ok(component.hasInput);
  });

  test('it clears console input on enter', async function (assert) {
    await render(hbs`{{console/ui-panel}}`);
    await component.runCommands('list this/thing/here', false);
    await settled();
    assert.strictEqual(component.consoleInputValue, '', 'empties input field on enter');
  });

  test('it clears the log when using clear command', async function (assert) {
    await render(hbs`{{console/ui-panel}}`);
    await component.runCommands(
      ['list this/thing/here', 'list this/other/thing', 'read another/thing'],
      false
    );
    await settled();
    assert.notEqual(component.logOutput, '', 'there is output in the log');

    await component.runCommands('clear', false);
    await settled();
    await component.up();
    await settled();
    assert.strictEqual(component.logOutput, '', 'clears the output log');
    assert.strictEqual(
      component.consoleInputValue,
      'clear',
      'populates console input with previous command on up after enter'
    );
  });

  test('it adds command to history on enter', async function (assert) {
    await render(hbs`{{console/ui-panel}}`);

    await component.runCommands('list this/thing/here', false);
    await settled();
    await component.up();
    await settled();
    assert.strictEqual(
      component.consoleInputValue,
      'list this/thing/here',
      'populates console input with previous command on up after enter'
    );
    await component.down();
    await settled();
    assert.strictEqual(component.consoleInputValue, '', 'populates console input with next command on down');
  });

  test('it cycles through history with more than one command', async function (assert) {
    await render(hbs`{{console/ui-panel}}`);
    await component.runCommands(['list this/thing/here', 'read that/thing/there', 'qwerty'], false);
    await settled();
    await component.up();
    await settled();
    assert.strictEqual(
      component.consoleInputValue,
      'qwerty',
      'populates console input with previous command on up after enter'
    );
    await component.up();
    await settled();
    assert.strictEqual(
      component.consoleInputValue,
      'read that/thing/there',
      'populates console input with previous command on up'
    );
    await component.up();
    await settled();
    assert.strictEqual(
      component.consoleInputValue,
      'list this/thing/here',
      'populates console input with previous command on up'
    );
    await component.up();
    await settled();
    assert.strictEqual(
      component.consoleInputValue,
      'qwerty',
      'populates console input with initial command if cycled through all previous commands'
    );
    await component.down();
    await settled();
    assert.strictEqual(
      component.consoleInputValue,
      '',
      'clears console input if down pressed after history is on most recent command'
    );
  });
});
