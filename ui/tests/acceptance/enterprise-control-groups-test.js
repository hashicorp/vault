/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { settled, currentURL, currentRouteName, visit, waitUntil } from '@ember/test-helpers';
import { module, test, skip } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';

import { storageKey } from 'vault/services/control-group';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import authForm from 'vault/tests/pages/components/auth-form';
import controlGroup from 'vault/tests/pages/components/control-group';
import controlGroupSuccess from 'vault/tests/pages/components/control-group-success';
import { writeSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import authPage from 'vault/tests/pages/auth';

const consoleComponent = create(consoleClass);
const authFormComponent = create(authForm);
const controlGroupComponent = create(controlGroup);
const controlGroupSuccessComponent = create(controlGroupSuccess);

module('Acceptance | Enterprise | control groups', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  const POLICY = `
    path "kv/foo" {
      capabilities = ["create", "read", "update", "delete", "list"]
      control_group = {
        max_ttl = "24h"
        factor "ops_manager" {
            identity {
                group_names = ["managers"]
                approvals = 1
            }
         }
      }
    }

    path "kv-v2-mount/data/foo" {
      capabilities = ["create", "read", "update", "list"]
      control_group = {
        max_ttl = "24h"
        factor "ops_manager" {
            identity {
                group_names = ["managers"]
                approvals = 1
            }
         }
      }
    }

    path "kv-v2-mount/*" {
      capabilities = ["list"]
    }
  `;

  const AUTHORIZER_POLICY = `
    path "sys/control-group/authorize" {
      capabilities = ["update"]
    }

    path "sys/control-group/request" {
      capabilities = ["update"]
    }
  `;

  const ADMIN_USER = 'authorizer';
  const ADMIN_PASSWORD = 'test';
  const setupControlGroup = async (context) => {
    await visit('/vault/secrets');
    await consoleComponent.toggle();
    await settled();
    await consoleComponent.runCommands([
      //enable kv-v1 mount and write a secret
      'write sys/mounts/kv type=kv',
      'write kv/foo bar=baz',

      //enable userpass, create user and associated entity
      'write sys/auth/userpass type=userpass',
      `write auth/userpass/users/${ADMIN_USER} password=${ADMIN_PASSWORD} policies=default`,
      `write identity/entity name=${ADMIN_USER} policies=test`,
      // write policies for control group + authorization
      `write sys/policies/acl/kv-control-group policy=${btoa(POLICY)}`,
      `write sys/policies/acl/authorizer policy=${btoa(AUTHORIZER_POLICY)}`,
      // read out mount to get the accessor
      'read -field=accessor sys/internal/ui/mounts/auth/userpass',
    ]);
    await settled();
    const userpassAccessor = consoleComponent.lastTextOutput;

    await consoleComponent.runCommands([
      // lookup entity id for our authorizer
      `write -field=id identity/lookup/entity name=${ADMIN_USER}`,
    ]);
    await settled();
    const authorizerEntityId = consoleComponent.lastTextOutput;
    await consoleComponent.runCommands([
      // create alias for authorizor and add them to the managers group
      `write identity/alias mount_accessor=${userpassAccessor} entity_id=${authorizerEntityId} name=${ADMIN_USER}`,
      `write identity/group name=managers member_entity_ids=${authorizerEntityId} policies=authorizer`,
      // create a token to request access to kv/foo
      'write -field=client_token auth/token/create policies=kv-control-group',
    ]);
    await settled();
    context.userToken = consoleComponent.lastLogOutput;

    await authPage.login(context.userToken);
    await settled();
    return this;
  };

  test('for v2 secrets it redirects you if you try to navigate to a Control Group restricted path', async function (assert) {
    await consoleComponent.runCommands([
      'write sys/mounts/kv-v2-mount type=kv-v2',
      'delete kv-v2-mount/metadata/foo',
    ]);
    await writeSecret('kv-v2-mount', 'foo', 'bar', 'baz');
    await settled();
    await setupControlGroup(this);
    await settled();
    await visit('/vault/secrets/kv-v2-mount/show/foo');

    assert.ok(
      await waitUntil(() => currentRouteName() === 'vault.cluster.access.control-group-accessor'),
      'redirects to access control group route'
    );
  });

  const workflow = async (assert, context, shouldStoreToken) => {
    const url = '/vault/secrets/kv/show/foo';
    await setupControlGroup(context);
    await settled();
    // as the requestor, go to the URL that's blocked by the control group
    // and store the values
    await visit(url);

    const accessor = controlGroupComponent.accessor;
    const controlGroupToken = controlGroupComponent.token;
    await authPage.logout();
    await settled();
    // log in as the admin, navigate to the accessor page,
    // and authorize the control group request
    await visit('/vault/auth?with=userpass');

    await authFormComponent.username(ADMIN_USER);
    await settled();
    await authFormComponent.password(ADMIN_PASSWORD);
    await settled();
    await authFormComponent.login();
    await settled();
    await visit(`/vault/access/control-groups/${accessor}`);

    // putting here to help with flaky test
    assert.dom('[data-test-authorize-button]').exists();
    await controlGroupComponent.authorize();
    await settled();
    assert.strictEqual(controlGroupComponent.bannerPrefix, 'Thanks!', 'text display changes');
    await settled();
    await authPage.logout();
    await settled();
    await authPage.login(context.userToken);
    await settled();
    if (shouldStoreToken) {
      localStorage.setItem(
        storageKey(accessor, 'kv/foo'),
        JSON.stringify({
          accessor,
          token: controlGroupToken,
          creation_path: 'kv/foo',
          uiParams: {
            url,
          },
        })
      );
      await visit(`/vault/access/control-groups/${accessor}`);

      assert.ok(controlGroupSuccessComponent.showsNavigateMessage, 'shows user the navigate message');
      await controlGroupSuccessComponent.navigate();
      await settled();
      assert.strictEqual(currentURL(), url, 'successfully loads the target url');
    } else {
      await visit(`/vault/access/control-groups/${accessor}`);

      await controlGroupSuccessComponent.token(controlGroupToken);
      await settled();
      await controlGroupSuccessComponent.unwrap();
      await settled();
      assert.ok(controlGroupSuccessComponent.showsJsonViewer, 'shows the json viewer');
    }
  };

  skip('it allows the full flow to work without a saved token', async function (assert) {
    await workflow(assert, this);
    await settled();
  });

  skip('it allows the full flow to work with a saved token', async function (assert) {
    await workflow(assert, this, true);
    await settled();
  });

  test('it displays the warning in the console when making a request to a Control Group path', async function (assert) {
    await setupControlGroup(this);
    await settled();
    await consoleComponent.toggle();
    await settled();
    await consoleComponent.runCommands('read kv/foo');
    await settled();
    const output = consoleComponent.lastLogOutput;
    assert.ok(output.includes('A Control Group was encountered at kv/foo'));
    assert.ok(output.includes('The Control Group Token is'));
    assert.ok(output.includes('The Accessor is'));
    assert.ok(output.includes('Visit /ui/vault/access/control-groups/'));
  });
});
