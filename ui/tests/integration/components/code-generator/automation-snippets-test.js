/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | code-generator/automation-snippets', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.cliArgs = undefined;
    this.tfvpArgs = undefined;
    this.apiArgs = undefined;
    this.customTabs = undefined;

    this.renderComponent = () => {
      return render(hbs`
        <CodeGenerator::AutomationSnippets
          @cliArgs={{this.cliArgs}}
          @tfvpArgs={{this.tfvpArgs}}
          @apiArgs={{this.apiArgs}}
          @customTabs={{this.customTabs}}
        />`);
    };
  });

  test('it does not render when args are undefined', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.fieldByAttr('terraform')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('cli')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('api')).doesNotExist();
  });

  test('it renders snippets when args are provided', async function (assert) {
    this.tfvpArgs = { resource: 'vault_mount', resourceArgs: { path: '"my-mount"', type: '"kv-v2"' } };
    this.cliArgs = {
      command: 'kv delete ',
      content: '-mount=secret creds',
    };
    this.apiArgs = {
      url: 'sys/mounts/{path}',
      payload: { path: 'my-mount', type: 'kv-v2' },
    };
    await this.renderComponent();
    const expectedTfvp = `resource "vault_mount" "<local identifier>" {
 path = "my-mount"
 type = "kv-v2"
}`;
    const expectedCli = 'vault kv delete -mount=secret creds';
    const expectedApi = `curl \\
  --header "X-Vault-Token: $VAULT_TOKEN" \\
  --request POST \\
  --data '{"path":"my-mount","type":"kv-v2"}' \\
  $VAULT_ADDR/v1/sys/mounts/my-mount
`;
    assert.dom(GENERAL.hdsTab('terraform')).exists();
    assert.dom(GENERAL.hdsTab('cli')).exists();
    assert.dom(GENERAL.hdsTab('api')).exists();
    assert.dom(GENERAL.fieldByAttr('terraform')).hasText(expectedTfvp);
    assert.dom(GENERAL.fieldByAttr('cli')).hasText(expectedCli);
    assert.dom(GENERAL.fieldByAttr('api')).hasText(expectedApi, 'it renders API snippet');
  });

  test('it selects tabs', async function (assert) {
    this.tfvpArgs = { resource: 'vault_mount', resourceArgs: { path: '"my-mount"', type: '"kv-v2"' } };
    this.cliArgs = {
      command: 'kv delete ',
      content: '-mount=secret creds',
    };
    this.apiArgs = {
      url: 'sys/mounts/{path}',
      payload: { path: 'my-mount', type: 'kv-v2' },
    };
    await this.renderComponent();
    assert.dom(GENERAL.hdsTabPanel('terraform')).doesNotHaveAttribute('hidden');
    assert.dom(GENERAL.hdsTabPanel('cli')).hasAttribute('hidden');
    assert.dom(GENERAL.hdsTabPanel('api')).hasAttribute('hidden');

    await click(GENERAL.hdsTab('cli'));
    assert.dom(GENERAL.hdsTabPanel('cli')).doesNotHaveAttribute('hidden');
    assert.dom(GENERAL.hdsTabPanel('terraform')).hasAttribute('hidden');
    assert.dom(GENERAL.hdsTabPanel('api')).hasAttribute('hidden');

    await click(GENERAL.hdsTab('api'));
    assert.dom(GENERAL.hdsTabPanel('api')).doesNotHaveAttribute('hidden');
    assert.dom(GENERAL.hdsTabPanel('terraform')).hasAttribute('hidden');
    assert.dom(GENERAL.hdsTabPanel('cli')).hasAttribute('hidden');
  });

  test('it does not render tabs when only one snippet arg exists', async function (assert) {
    this.apiArgs = {
      url: 'sys/mounts/{path}',
      payload: { path: 'my-mount', type: 'kv-v2' },
    };
    await this.renderComponent();
    assert.dom(GENERAL.hdsTab()).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('api')).exists();
  });

  test('it conditionally renders tabs based on provided args', async function (assert) {
    // Only terraform
    this.tfvpArgs = { resource: 'vault_mount', resourceArgs: { path: '"my-mount"' } };
    await this.renderComponent();
    assert.dom(GENERAL.hdsTab()).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('terraform')).exists();
    assert.dom(GENERAL.fieldByAttr('cli')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('api')).doesNotExist();

    // Add CLI (use set() to trigger DOM update)
    this.set('cliArgs', { command: 'secrets enable', content: '-path=my-mount kv-v2' });
    assert.dom(GENERAL.hdsTab()).exists({ count: 2 }, 'terraform and cli tabs rendered');
    assert.dom(GENERAL.hdsTab('terraform')).exists();
    assert.dom(GENERAL.hdsTab('cli')).exists();
    assert.dom(GENERAL.hdsTab('api')).doesNotExist();

    // Add API (use set() to trigger DOM update)
    this.set('apiArgs', { url: 'sys/mounts/{path}', payload: { path: 'my-mount' } });
    assert.dom(GENERAL.hdsTab()).exists({ count: 3 }, 'all three tabs rendered');
    assert.dom(GENERAL.hdsTab('terraform')).exists();
    assert.dom(GENERAL.hdsTab('cli')).exists();
    assert.dom(GENERAL.hdsTab('api')).exists();
  });

  test('it includes namespace in API snippet for non-root namespaces', async function (assert) {
    this.apiArgs = {
      url: 'sys/mounts/{path}',
      payload: { path: 'my-mount', type: 'kv-v2' },
    };
    const namespace = this.owner.lookup('service:namespace');
    namespace.path = 'admin';
    await this.renderComponent();
    const expectedSnippet = `curl \\
  --header "X-Vault-Token: $VAULT_TOKEN" \\
  --header "X-Vault-Namespace: admin" \\
  --request POST \\
  --data '{"path":"my-mount","type":"kv-v2"}' \\
  $VAULT_ADDR/v1/sys/mounts/my-mount
`;
    assert.dom(GENERAL.fieldByAttr('api')).hasText(expectedSnippet, 'it renders API snippet with namespace');
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
