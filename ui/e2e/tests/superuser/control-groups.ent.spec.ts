/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, test } from 'e2e/fixtures/demo';
import { rootApi } from 'e2e/fixtures/root-api';

import type { APIRequestContext, Browser, TestInfo } from '@playwright/test';

const AUTHORIZER_PASSWORD = 'authorizer-password';

// Opens a page with no stored session so it can sign in as a different identity than the persona.
const newSignedOutPage = async (browser: Browser, testInfo: TestInfo) => {
  // An explicit empty state stops the context from inheriting the persona's session.
  const context = await browser.newContext({
    baseURL: testInfo.project.use.baseURL,
    storageState: { cookies: [], origins: [] },
  });
  return { context, page: await context.newPage() };
};

const write = async (api: APIRequestContext, path: string, data: object, method: 'post' | 'put' = 'post') => {
  const response = await api[method](`/v1/${path}`, { data });
  await expect(response).toBeOK();
  return response;
};

// Configure the control group max TTL through the UI and confirm Vault stored it.
test('control group configuration saves and persists the maximum TTL', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);

  await test.step('open the configure page from the approval workflow page', async () => {
    await page.goto('access/control-groups');
    await expect(page.getByRole('heading', { name: 'Approval workflow' })).toBeVisible();
    await page.getByRole('link', { name: 'Configure' }).click();
    await expect(page).toHaveURL(/access\/control-groups\/configure$/);
  });

  await test.step('save a maximum TTL of 2 hours', async () => {
    await page.getByRole('group', { name: /Maximum TTL/ }).click();
    await page.getByLabel('TTL unit for Maximum TTL').selectOption('h');
    await page.getByRole('textbox', { name: 'Number of units' }).fill('2');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page.getByText('The configuration was saved successfully.')).toBeVisible();
  });

  await test.step('Vault stores the value and the UI shows it after a reload', async () => {
    const config = await api.get('/v1/sys/config/control-group');
    await expect(config).toBeOK();
    expect((await config.json()).data.max_ttl).toBe(7200);

    await page.reload();
    await expect(page.getByLabel('TTL unit for Maximum TTL')).toHaveValue('h');
    await expect(page.getByRole('textbox', { name: 'Number of units' })).toHaveValue('2');
  });

  await test.step('remove the configuration', async () => {
    await expect(await api.delete('/v1/sys/config/control-group')).toBeOK();
  });
});

// A requester hits a control-group-protected secret, an authorizer approves the accessor in the UI,
// and the requester then uses the approved request to read the secret.
test('control group request is authorized by a group member and unlocked by the requester', async ({
  browser,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const id = `${Date.now()}`;
  const kvMount = `cg-kv-${id}`;
  const authorizer = `cg-authorizer-${id}`;
  const managers = `cg-managers-${id}`;
  let accessor = '';

  await test.step('seed a protected secret, an authorizer group member and a requester token', async () => {
    await write(api, `sys/mounts/${kvMount}`, { type: 'kv' });
    await write(api, `${kvMount}/foo`, { bar: 'baz' });

    if (!(await api.get('/v1/sys/auth/userpass')).ok()) {
      await write(api, 'sys/auth/userpass', { type: 'userpass' });
    }
    const userpass = await (await api.get('/v1/sys/auth/userpass')).json();
    await write(api, `auth/userpass/users/${authorizer}`, {
      password: AUTHORIZER_PASSWORD,
      token_policies: 'default',
    });
    const entity = await (await write(api, 'identity/entity', { name: authorizer })).json();
    await write(api, 'identity/entity-alias', {
      name: authorizer,
      canonical_id: entity.data.id,
      mount_accessor: userpass.data.accessor,
    });

    await write(
      api,
      `sys/policies/acl/cg-authorizer-${id}`,
      {
        policy: `
          path "sys/control-group/authorize" { capabilities = ["update"] }
          path "sys/control-group/request" { capabilities = ["update"] }
        `,
      },
      'put'
    );
    await write(api, 'identity/group', {
      name: managers,
      member_entity_ids: [entity.data.id],
      policies: [`cg-authorizer-${id}`],
    });
    await write(
      api,
      `sys/policies/acl/cg-requester-${id}`,
      {
        policy: `
          path "${kvMount}/foo" {
            capabilities = ["create", "read", "update", "delete", "list"]
            control_group = {
              max_ttl = "24h"
              factor "ops_manager" {
                identity {
                  group_names = ["${managers}"]
                  approvals = 1
                }
              }
            }
          }
        `,
      },
      'put'
    );
  });

  const requesterToken = (
    await (await write(api, 'auth/token/create', { policies: [`cg-requester-${id}`], ttl: '1h' })).json()
  ).auth.client_token;

  const requester = await newSignedOutPage(browser, testInfo);
  const approver = await newSignedOutPage(browser, testInfo);

  await test.step('the requester is redirected to the accessor page for the locked secret', async () => {
    const { page: rp } = requester;
    await rp.goto('auth?with=token');
    await rp.getByRole('textbox', { name: 'Token' }).fill(requesterToken);
    await rp.getByRole('button', { name: 'Sign in' }).click();
    await expect(rp).toHaveURL(/dashboard/);
    await rp.goto(`secrets-engines/${kvMount}/show/foo`);

    await expect(rp).toHaveURL(/access\/control-groups\/.+/);
    await expect(
      rp.getByText('Someone is requesting access to a path locked by a Control Group')
    ).toBeVisible();
    await expect(rp.getByText('Awaiting authorization.')).toBeVisible();
    accessor = (await rp
      .getByRole('alertdialog', { name: 'Accessor' })
      .getByRole('button')
      .textContent())!.trim();
    expect(accessor).not.toBe('');
  });

  await test.step('the authorizer looks up the accessor and authorizes the request', async () => {
    const { page: ap } = approver;
    await ap.goto('auth?with=userpass');
    await ap.getByRole('textbox', { name: 'Username' }).fill(authorizer);
    await ap.getByRole('textbox', { name: 'Password' }).fill(AUTHORIZER_PASSWORD);
    await ap.getByRole('button', { name: 'Sign in' }).click();
    await expect(ap).toHaveURL(/dashboard/);

    await ap.goto('access/control-groups');
    await ap.getByRole('textbox', { name: 'Accessor' }).fill(accessor);
    await ap.getByRole('button', { name: 'Lookup' }).click();
    await expect(ap).toHaveURL(new RegExp(`access/control-groups/${accessor}$`));
    await expect(ap.getByText(`Someone is requesting access to ${kvMount}/foo`)).toBeVisible();

    await ap.getByRole('button', { name: 'Authorize' }).click();
    await expect(ap.getByText('Thanks!')).toBeVisible();
    await expect(ap.getByText('Already approved by')).toContainText(authorizer);
  });

  await test.step('Vault reports the request as approved', async () => {
    const lookup = await api.post('/v1/sys/control-group/request', { data: { accessor } });
    await expect(lookup).toBeOK();
    const { data } = await lookup.json();
    expect(data.approved).toBe(true);
    expect(data.authorizations.map((a: { entity_name: string }) => a.entity_name)).toContain(authorizer);
  });

  await test.step('the requester sees the approval and opens the unlocked secret', async () => {
    const { page: rp } = requester;
    await rp.goto(`access/control-groups/${accessor}`);
    await rp.getByRole('button', { name: 'Visit' }).click();
    await expect(rp).toHaveURL(new RegExp(`secrets-engines/${kvMount}/show/foo`));
    await rp.getByRole('button', { name: 'show value' }).click();
    await expect(rp.getByText('baz', { exact: true })).toBeVisible();
  });

  await requester.context.close();
  await approver.context.close();
});
