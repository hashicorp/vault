/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { createHmac } from 'crypto';
import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureNamespace, rootApi } from '../../fixtures/root-api';
import { DISMISSED_WIZARD_KEY, WIZARD_ID_MAP } from '../../../app/utils/constants/wizard';

import type { APIRequestContext, Page } from '@playwright/test';

const NS = 'auth-login-ns';
const NS_MOUNT = 'auth-login-ns-up';
const ROOT_MOUNT = 'auth-login-up';
const MFA_MOUNT = 'auth-login-mfa';
const MFA_ENFORCEMENT = 'auth-login-mfa-enforcement';
const MFA_ENTITY = 'auth-login-mfa-entity';
const DEFAULT_AUTH_RULE = 'auth-login-default';
const PASSWORD = 'auth-login-password';

// These flows start signed out, so no root token is stored. Wizards are dismissed as in init.setup.ts.
test.use({ storageState: { cookies: [], origins: [] } });

const BASE32 = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';

// Generates the current 6 digit, 30 second SHA1 TOTP code (RFC 6238) for a base32 secret.
const totpCode = (secret: string) => {
  let bits = '';
  for (const char of secret.replace(/=+$/, '').toUpperCase()) {
    bits += BASE32.indexOf(char).toString(2).padStart(5, '0');
  }
  const key = Buffer.from((bits.match(/.{8}/g) ?? []).map((byte) => parseInt(byte, 2)));
  const counter = Buffer.alloc(8);
  counter.writeBigUInt64BE(BigInt(Math.floor(Date.now() / 1000 / 30)));
  const hmac = createHmac('sha1', key).update(counter).digest();
  const offset = hmac[hmac.length - 1] & 0xf;
  const value = hmac.readUInt32BE(offset) & 0x7fffffff;
  return String(value % 1_000_000).padStart(6, '0');
};

const inNamespace = { 'X-Vault-Namespace': NS };

// Removes everything this spec creates. Failures are ignored because most resources may not exist yet.
const cleanup = async (api: APIRequestContext) => {
  await api.delete(`/v1/sys/config/ui/login/default-auth/${DEFAULT_AUTH_RULE}`);
  await api.delete(`/v1/identity/mfa/login-enforcement/${MFA_ENFORCEMENT}`);
  const methods = await api.get('/v1/identity/mfa/method/totp', { params: { list: 'true' } });
  if (methods.ok()) {
    const { data } = await methods.json();
    for (const id of data.keys as string[]) {
      const method = await api.get(`/v1/identity/mfa/method/totp/${id}`);
      if ((await method.json()).data?.issuer === MFA_ENFORCEMENT) {
        await api.delete(`/v1/identity/mfa/method/totp/${id}`);
      }
    }
  }
  await api.delete(`/v1/identity/entity/name/${MFA_ENTITY}`);
  await api.delete(`/v1/sys/auth/${ROOT_MOUNT}`);
  await api.delete(`/v1/sys/auth/${MFA_MOUNT}`);
  await api.delete(`/v1/sys/auth/${NS_MOUNT}`, { headers: inNamespace });
  // Namespace deletion is asynchronous and fails while the namespace is not empty, so retry it until the namespace
  // is gone. Otherwise a repeated run could reuse a namespace that is still being deleted.
  await expect
    .poll(
      async () => {
        const res = await api.get(`/v1/sys/namespaces/${NS}`);
        if (res.status() !== 404) await api.delete(`/v1/sys/namespaces/${NS}`);
        return res.status();
      },
      { timeout: 30_000 }
    )
    .toBe(404);
};

const enableUserpass = async (api: APIRequestContext, mount: string, user: string, headers = {}) => {
  await expect(await api.post(`/v1/sys/auth/${mount}`, { data: { type: 'userpass' }, headers })).toBeOK();
  await expect(
    await api.post(`/v1/auth/${mount}/users/${user}`, { data: { password: PASSWORD }, headers })
  ).toBeOK();
};

let totpSecret = '';
let totpMethodId = '';
let mfaEntityId = '';

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await cleanup(api);
  await ensureNamespace(api, NS);
  await enableUserpass(api, NS_MOUNT, 'nsuser', inNamespace);
  await enableUserpass(api, ROOT_MOUNT, 'alice');
  await enableUserpass(api, MFA_MOUNT, 'mfauser');

  // The enforcement applies to every user of the mfa mount, which has a single user with a TOTP secret.
  const auth = await (await api.get('/v1/sys/auth')).json();
  const accessor = auth.data[`${MFA_MOUNT}/`].accessor;
  const method = await api.post('/v1/identity/mfa/method/totp', {
    data: { issuer: MFA_ENFORCEMENT, period: 30, key_size: 20, algorithm: 'SHA1', digits: 6 },
  });
  await expect(method).toBeOK();
  totpMethodId = (await method.json()).data.method_id;
  const entity = await api.post('/v1/identity/entity', { data: { name: MFA_ENTITY } });
  await expect(entity).toBeOK();
  mfaEntityId = (await entity.json()).data.id;
  await expect(
    await api.post('/v1/identity/entity-alias', {
      data: { name: 'mfauser', canonical_id: mfaEntityId, mount_accessor: accessor },
    })
  ).toBeOK();
  const generated = await api.post('/v1/identity/mfa/method/totp/admin-generate', {
    data: { method_id: totpMethodId, entity_id: mfaEntityId },
  });
  await expect(generated).toBeOK();
  totpSecret = new URL((await generated.json()).data.url).searchParams.get('secret') as string;
  await expect(
    await api.post(`/v1/identity/mfa/login-enforcement/${MFA_ENFORCEMENT}`, {
      data: { mfa_method_ids: [totpMethodId], auth_method_accessors: [accessor] },
    })
  ).toBeOK();
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await cleanup(api);
  await api.dispose();
});

test.beforeEach(async ({ page }) => {
  await page.addInitScript(({ key, ids }) => localStorage.setItem(key, JSON.stringify(ids)), {
    key: DISMISSED_WIZARD_KEY,
    ids: Object.values(WIZARD_ID_MAP),
  });
});

// Reads the persisted token data for the signed in auth method, which is what the UI uses after a reload.
const storedTokenData = (page: Page, authMethodType: string) =>
  page.evaluate((type) => {
    const key = Object.keys(localStorage).find((k) => k.startsWith(`vault-${type}☃`));
    return key ? JSON.parse(localStorage.getItem(key) as string) : null;
  }, authMethodType);

const signInWithUserpass = async (page: Page, username: string, mount: string) => {
  await page.getByLabel('Method').selectOption('userpass');
  await page.getByRole('textbox', { name: 'Username' }).fill(username);
  await page.getByRole('textbox', { name: 'Password' }).fill(PASSWORD);
  await page.getByRole('button', { name: 'Advanced settings' }).click();
  await page.getByRole('textbox', { name: 'Mount path' }).fill(mount);
};

// Login inside a namespace must look up the new token so the UI knows which namespace the token belongs to.
test('userpass login in a namespace looks up the token and persists its namespace', async ({ page }) => {
  await test.step('open the login page for the namespace', async () => {
    await page.goto(`auth?namespace=${NS}`);
    await expect(page.getByRole('heading', { name: 'Sign in to Vault' })).toBeVisible();
  });

  const lookup = page.waitForRequest((req) => req.url().includes('/v1/auth/token/lookup-self'));

  await test.step('sign in with a userpass user', async () => {
    await signInWithUserpass(page, 'nsuser', NS_MOUNT);
    await page.getByRole('button', { name: 'Sign in' }).click();
    await expect(page).toHaveURL(/\/vault\/dashboard/);
  });

  await test.step('lookup-self was requested with the new token in the namespace', async () => {
    const request = await lookup;
    const stored = await storedTokenData(page, 'userpass');
    expect(request.method()).toBe('GET');
    expect(await request.headerValue('x-vault-namespace')).toBe(NS);
    expect(await request.headerValue('x-vault-token')).toBe(stored.token);
    expect((await request.response())?.status()).toBe(200);
  });

  await test.step('token data returned by lookup-self survives a reload', async () => {
    const stored = await storedTokenData(page, 'userpass');
    expect(stored.displayName).toBe('nsuser');
    expect(stored.namespacePath).toBe(`${NS}/`);
    await page.reload();
    await expect(page).toHaveURL(/\/vault\/dashboard/);
    await page.getByRole('button', { name: 'User menu' }).click();
    await expect(page.getByText('nsuser')).toBeVisible();
  });
});

// Token login reads the token itself, so it must keep working for a token created inside a namespace.
test('token login with a namespace signs in and persists the token namespace', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const created = await api.post('/v1/auth/token/create', {
    data: { display_name: 'auth-login-token', policies: ['default'] },
    headers: inNamespace,
  });
  await expect(created).toBeOK();
  const { client_token } = (await created.json()).auth;
  await api.dispose();

  await test.step('sign in with the token and namespace', async () => {
    await page.goto(`auth?namespace=${NS}`);
    await page.getByLabel('Method').selectOption('token');
    await page.getByRole('textbox', { name: 'Token' }).fill(client_token);
    await page.getByRole('button', { name: 'Sign in' }).click();
    await expect(page).toHaveURL(/\/vault\/dashboard/);
  });

  await test.step('stored token data identifies the namespace and survives a reload', async () => {
    const stored = await storedTokenData(page, 'token');
    expect(stored.displayName).toBe('token-auth-login-token');
    expect(stored.namespacePath).toBe(`${NS}/`);
    await page.reload();
    await expect(page).toHaveURL(/\/vault\/dashboard/);
  });
});

// MFA validation returns the auth data that completes the login, and its error messages must reach the user.
test('mfa login validates the passcode, signs in as the entity and rejects a reused passcode', async ({
  page,
  playwright,
}, testInfo) => {
  let code = '';

  await test.step('sign in and reach the mfa prompt', async () => {
    await page.goto('auth');
    await signInWithUserpass(page, 'mfauser', MFA_MOUNT);
    await page.getByRole('button', { name: 'Sign in' }).click();
    await expect(page.getByRole('heading', { name: 'Multi-factor authentication' })).toBeVisible();
  });

  await test.step('verify the passcode and land on the dashboard', async () => {
    code = totpCode(totpSecret);
    const validate = page.waitForRequest((req) => req.url().includes('/v1/sys/mfa/validate'));
    await page.getByText('Enter your authentication code to log in. TOTP passcode').click();
    await page.getByRole('textbox', { name: 'TOTP passcode' }).fill(code);
    await page.getByRole('button', { name: 'Verify' }).click();
    const request = await validate;
    expect(request.method()).toBe('POST');
    expect(request.postDataJSON()).toEqual({
      mfa_request_id: expect.any(String),
      mfa_payload: { [totpMethodId]: [code] },
    });
    expect((await request.response())?.status()).toBe(200);
    await expect(page).toHaveURL(/\/vault\/dashboard/);
  });

  await test.step('the stored token belongs to the mfa entity', async () => {
    const stored = await storedTokenData(page, 'userpass');
    expect(stored.displayName).toBe('mfauser');
    const api = await rootApi(playwright, testInfo);
    const self = await api.get('/v1/auth/token/lookup-self', { headers: { 'X-Vault-Token': stored.token } });
    await expect(self).toBeOK();
    expect((await self.json()).data.entity_id).toBe(mfaEntityId);
    await api.dispose();
  });

  await test.step('signing in again with the same passcode is rejected and shows the validation error', async () => {
    await page.getByRole('button', { name: 'User menu' }).click();
    await page.getByRole('link', { name: 'Log out' }).click();
    await signInWithUserpass(page, 'mfauser', MFA_MOUNT);
    await page.getByRole('button', { name: 'Sign in' }).click();
    await expect(page.getByRole('heading', { name: 'Multi-factor authentication' })).toBeVisible();
    await page.getByText('Enter your authentication code to log in. TOTP passcode').click();
    await page.getByRole('textbox', { name: 'TOTP passcode' }).fill(code);
    const reused = page.waitForResponse((res) => res.url().includes('/v1/sys/mfa/validate'));
    await page.getByRole('button', { name: 'Verify' }).click();
    expect((await reused).status()).toBe(403);
    await expect(page.getByText('The passcode failed to validate')).toBeVisible();
    expect(await storedTokenData(page, 'userpass')).toBeNull();
    await expect(page.getByRole('heading', { name: 'Multi-factor authentication' })).toBeVisible();
  });
});

// Enterprise admins can preselect the login method; the page must read that setting without a token.
test('default and backup auth methods control which login methods render first', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const configure = async (data: Record<string, unknown>) =>
    expect(
      await api.post(`/v1/sys/config/ui/login/default-auth/${DEFAULT_AUTH_RULE}`, {
        data: { namespace_path: '', ...data },
      })
    ).toBeOK();

  try {
    await test.step('default type renders first and backups sit behind "Sign in with other methods"', async () => {
      await configure({ default_auth_type: 'userpass', backup_auth_types: ['token', 'ldap'] });
      const settings = page.waitForResponse((res) =>
        res.url().includes('/v1/sys/internal/ui/default-auth-methods')
      );
      await page.goto('auth');
      const response = await settings;
      expect(response.status()).toBe(200);
      expect((await response.json()).data).toMatchObject({
        default_auth_type: 'userpass',
        backup_auth_types: ['token', 'ldap'],
      });
      await expect(page.getByRole('tab')).toHaveCount(1);
      await expect(page.getByRole('tab', { name: 'Userpass' })).toHaveAttribute('aria-selected', 'true');
      await page.getByRole('button', { name: 'Sign in with other methods' }).click();
      await expect(page.getByRole('tab')).toHaveCount(2);
      await expect(page.getByRole('tab', { name: 'Token' })).toHaveAttribute('aria-selected', 'true');
    });

    await test.step('without a default, the backup types render first', async () => {
      await api.delete(`/v1/sys/config/ui/login/default-auth/${DEFAULT_AUTH_RULE}`);
      await configure({ default_auth_type: '', backup_auth_types: ['userpass', 'token'] });
      await page.goto('auth');
      await expect(page.getByRole('tab')).toHaveCount(2);
      await expect(page.getByRole('tab', { name: 'Userpass' })).toHaveAttribute('aria-selected', 'true');
      await expect(page.getByRole('button', { name: 'Sign in with other methods' })).toHaveCount(0);
    });
  } finally {
    await api.delete(`/v1/sys/config/ui/login/default-auth/${DEFAULT_AUTH_RULE}`);
    await api.dispose();
  }
});
