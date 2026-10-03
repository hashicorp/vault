/**
 * Copyright IBM Corp. 2016, 2025
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
    await render(hbs`<div class="panel-open"><Console::UiPanel /></div>`);
    assert.ok(component.hasInput);
    assert
      .dom(this.element)
      .hasText(
        "The Vault Web REPL provides an easy way to execute common Vault CLI commands, such as write, read, delete, and list. It does not include KV version 2 write or put commands. For guidance, type `help`. For more detailed documentation, see the HashiCorp Developer site. Examples: → Write secrets to kv v1: write <mount>/my-secret foo=bar → List kv v1 secret keys: list <mount>/ → Read a kv v1 secret: read <mount>/my-secret → Mount a kv v2 secret engine: write sys/mounts/<mount> type=kv options=version=2 → Read a kv v2 secret: kv-get <mount>/secret-path → Read a kv v2 secret's metadata: kv-get <mount>/secret-path -metadata"
      );
  });

  test('it clears console input on enter', async function (assert) {
    await render(hbs`<Console::UiPanel />`);
    await component.runCommands('list this/thing/here', false);
    await settled();
    assert.strictEqual(component.consoleInputValue, '', 'empties input field on enter');
  });

  test('it clears the log when using clear command', async function (assert) {
    await render(hbs`<Console::UiPanel />`);
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
    await render(hbs`<Console::UiPanel />`);

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
    await render(hbs`<Console::UiPanel />`);
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
