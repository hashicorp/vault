/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { settled, waitUntil, click } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { setupApplicationTest } from 'ember-qunit';
import enginesPage from 'vault/tests/pages/secrets/backends';
import authPage from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { v4 as uuidv4 } from 'uuid';

const consoleComponent = create(consoleClass);

module('Acceptance | console', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test("refresh reloads the current route's data", async function (assert) {
    assert.expect(6);
    await enginesPage.visit();
    await settled();
    await consoleComponent.toggle();
    await settled();
    const ids = [uuidv4(), uuidv4(), uuidv4()];
    for (const id of ids) {
      const inputString = `write sys/mounts/console-route-${id} type=kv`;
      await consoleComponent.runCommands(inputString);
      await settled();
    }
    await consoleComponent.runCommands('refresh');
    await settled();
    for (const id of ids) {
      assert.ok(
        enginesPage.rows.findOneBy('path', `console-route-${id}/`),
        'new engine is shown on the page'
      );
    }
    // Clean up
    for (const id of ids) {
      const inputString = `delete sys/mounts/console-route-${id}`;
      await consoleComponent.runCommands(inputString);
      await settled();
    }
    await consoleComponent.runCommands('refresh');
    await settled();
    for (const id of ids) {
      assert.throws(() => {
        enginesPage.rows.findOneBy('path', `console-route-${id}/`);
      }, 'engine was removed');
    }
  });

  test('fullscreen command expands the cli panel', async function (assert) {
    await consoleComponent.toggle();
    await settled();
    await consoleComponent.runCommands('fullscreen');
    await settled();
    const consoleEle = document.querySelector('[data-test-component="console/ui-panel"]');
    // wait for the CSS transition to finish
    await waitUntil(() => consoleEle.offsetHeight === window.innerHeight);
    assert.strictEqual(
      consoleEle.offsetHeight,
      window.innerHeight,
      'fullscreen is the same height as the window'
    );
  });

  test('array output is correctly formatted', async function (assert) {
    await consoleComponent.toggle();
    await settled();
    await consoleComponent.runCommands('read -field=policies /auth/token/lookup-self');
    await settled();
    const consoleOut = document.querySelector('.console-ui-output>pre');
    // wait for the CSS transition to finish
    await waitUntil(() => consoleOut.innerText);
    assert.notOk(consoleOut.innerText.includes('function(){'));
    assert.strictEqual(consoleOut.innerText, '["root"]');
  });

  test('number output is correctly formatted', async function (assert) {
    await consoleComponent.toggle();
    await settled();
    await consoleComponent.runCommands('read -field=creation_time /auth/token/lookup-self');
    await settled();
    const consoleOut = document.querySelector('.console-ui-output>pre');
    // wait for the CSS transition to finish
    await waitUntil(() => consoleOut.innerText);
    assert.strictEqual(consoleOut.innerText.match(/^\d+$/).length, 1);
  });

  test('boolean output is correctly formatted', async function (assert) {
    await consoleComponent.toggle();
    await settled();
    await consoleComponent.runCommands('read -field=orphan /auth/token/lookup-self');
    await settled();
    const consoleOut = document.querySelector('.console-ui-output>pre');
    // have to wrap in a later so that we can wait for the CSS transition to finish
    await waitUntil(() => consoleOut.innerText);
    assert.strictEqual(consoleOut.innerText.match(/^(true|false)$/g).length, 1);
  });

  test('it should open and close console panel', async function (assert) {
    await click('[data-test-console-toggle]');
    assert.dom('[data-test-console-panel]').hasClass('panel-open', 'Sidebar button opens console panel');
    await click('[data-test-console-toggle]');
    assert
      .dom('[data-test-console-panel]')
      .doesNotHaveClass('panel-open', 'Sidebar button closes console panel');
    await click('[data-test-console-toggle]');
    await click('[data-test-console-panel-close]');
    assert
      .dom('[data-test-console-panel]')
      .doesNotHaveClass('panel-open', 'Console panel close button closes console panel');
  });
});
