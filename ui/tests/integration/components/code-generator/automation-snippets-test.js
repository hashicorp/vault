/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | code-generator/automation-snippets', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.cliArgs = undefined;
    this.tfvpArgs = undefined;
    this.customTabs = undefined;

    this.renderComponent = () => {
      return render(hbs`
        <CodeGenerator::AutomationSnippets
          @cliArgs={{this.cliArgs}}
          @tfvpArgs={{this.tfvpArgs}}
          @customTabs={{this.customTabs}}
        />`);
    };
  });

  test('it renders default tabs and snippets', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.fieldByAttr('terraform')).hasClass('language-hcl');
    assert.dom(GENERAL.fieldByAttr('cli')).hasClass('language-shell');
    assert.dom(GENERAL.hdsTab()).exists({ count: 2 });
    assert.dom(GENERAL.hdsTabPanel()).exists({ count: 2 });
    assert
      .dom(GENERAL.hdsTab('terraform'))
      .exists()
      .hasAttribute('aria-selected', 'true')
      .hasText('Terraform Vault Provider');
    assert.dom(GENERAL.hdsTab('cli')).exists().hasAttribute('aria-selected', 'false').hasText('CLI');
    assert.dom(GENERAL.hdsTabPanel('terraform')).doesNotHaveAttribute('hidden');
    assert.dom(GENERAL.hdsTabPanel('cli')).hasAttribute('hidden');
    const expectedTfvp = `resource "<resource name>" "<local identifier>" {

}`;
    assert.dom(GENERAL.fieldByAttr('terraform')).hasText(expectedTfvp, 'it renders empty terraform snippet');
    // Click CLI tab
    const expectedCli = 'vault <command> [args]';
    await click(GENERAL.hdsTab('cli'));
    assert.dom(GENERAL.hdsTab('cli')).hasAttribute('aria-selected', 'true');
    assert.dom(GENERAL.hdsTabPanel('cli')).doesNotHaveAttribute('hidden', 'clicking "CLI" shows cli snippet');
    assert.dom(GENERAL.fieldByAttr('cli')).hasText(expectedCli, 'it renders empty cli snippet');
    assert.dom(GENERAL.hdsTab('terraform')).hasAttribute('aria-selected', 'false');
    assert.dom(GENERAL.hdsTabPanel('terraform')).hasAttribute('hidden');
  });

  test('it renders snippet overrides', async function (assert) {
    this.tfvpArgs = { resource: 'vault_mount', resourceArgs: { path: '"my-mount"', type: '"kv-v2"' } };
    this.cliArgs = {
      command: 'kv delete ',
      content: '-mount=secret creds',
    };
    await this.renderComponent();
    const expectedTfvp = `resource "vault_mount" "<local identifier>" {
 path = "my-mount" 
 type = "kv-v2" 
}`;
    const expectedCli = 'vault kv delete -mount=secret creds';
    await click(GENERAL.hdsTab('cli'));
    assert.dom(GENERAL.fieldByAttr('terraform')).hasText(expectedTfvp);
    assert.dom(GENERAL.fieldByAttr('cli')).hasText(expectedCli);
  });

  test('it includes namespace in snippet for non-root namespaces', async function (assert) {
    this.tfvpArgs = { resource: 'vault_mount', resourceArgs: { path: '"my-mount"', type: '"kv-v2"' } };
    const namespace = this.owner.lookup('service:namespace');
    namespace.path = 'admin';
    await this.renderComponent();
    const expectedSnippet = `resource "vault_mount" "<local identifier>" {
 namespace = "admin"
 path = "my-mount" 
 type = "kv-v2" 
}`;
    assert
      .dom(GENERAL.fieldByAttr('terraform'))
      .hasText(expectedSnippet, 'it renders snippet with namespace');
  });

  test('it renders custom tabs', async function (assert) {
    this.customTabs = [
      {
        key: 'banana',
        label: '🍌 Banana Config',
        snippet: 'banana:\n  ripeness: perfect\n  color: yellow',
        language: 'yaml',
      },
      {
        key: 'pizza',
        label: '🍕 Pizza CLI',
        snippet: 'order pizza --toppings=pepperoni --size=large',
        language: 'shell',
      },
      {
        key: 'magic8ball',
        label: '🎱 Magic 8-Ball Response',
        snippet: `{
  "question": "Should I eat spaghetti for dinner?",
  "answer": "It is certain",
  "confidence": "100%",
  "ask_again_later": true
}`,
        language: 'json',
      },
    ];
    await this.renderComponent();
    assert.dom(GENERAL.hdsTab()).exists({ count: this.customTabs.length });
    assert.dom(GENERAL.hdsTabPanel()).exists({ count: this.customTabs.length });
    this.customTabs.forEach(async ({ key, label, snippet, language }) => {
      assert.dom(GENERAL.hdsTab(key)).hasText(label);
      assert.dom(GENERAL.hdsTabPanel(key)).hasText(snippet);
      assert.dom(GENERAL.fieldByAttr(key)).hasText(snippet, `it renders snippet for tab: ${key}`);
      assert.dom(GENERAL.fieldByAttr(key)).hasClass(`language-${language}`);
    });
  });
});
