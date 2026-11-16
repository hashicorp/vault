/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import type { APIRequestContext, Page, PlaywrightWorkerArgs, TestInfo } from '@playwright/test';
import { highlightElement, test } from '../../fixtures/demo';
import { ensureMount, ensureNamespace, rootApi } from '../../fixtures/root-api';
import { BasePage } from '../../pages/base';

const MODES = ['dr', 'performance'];

// These workflows create activation tokens, so do not retain secret-bearing traces.
test.use({ trace: 'off' });

// Replication outlives the test run, so always return the cluster to disabled.
const disableReplication = async (playwright: PlaywrightWorkerArgs['playwright'], testInfo: TestInfo) => {
  const api = await rootApi(playwright, testInfo);
  for (const mode of MODES) {
    await api.post(`/v1/sys/replication/${mode}/primary/disable`);
  }
  await api.dispose();
};

test.beforeAll(async ({ playwright }, testInfo) => disableReplication(playwright, testInfo));
test.afterAll(async ({ playwright }, testInfo) => disableReplication(playwright, testInfo));

// Characterizes the replication routes and flows against a real Enterprise cluster.
// Replication needs transactional storage, so this runs on the raft persona.
// The routes are checked through their URLs and what the pages show.
test('replication routes and workflows', async ({ page }) => {
  test.setTimeout(180_000);

  await test.step('replication starts disabled and offers to enable it', async () => {
    await page.goto('replication');
    await expect(page).toHaveURL(/\/ui\/vault\/replication$/);
    await expect(page.getByRole('heading', { name: 'Enable replication' })).toBeVisible();
    await expect(page.locator('[data-test-replication-type-select="dr"]')).toBeVisible();
    await expect(page.locator('[data-test-replication-type-select="performance"]')).toBeVisible();
  });

  await test.step('an unknown replication mode redirects to the replication index', async () => {
    await page.goto('replication/not-a-mode');
    await expect(page).toHaveURL(/\/ui\/vault\/replication$/);
  });

  await test.step('manage and secondaries redirect to the mode details while disabled', async () => {
    await page.goto('replication/performance/manage');
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance$/);
    await page.goto('replication/performance/secondaries/add');
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance$/);
  });

  await test.step('the sidebar links to the overview and each mode', async () => {
    await page.goto('replication');
    await page.getByRole('link', { name: 'Disaster recovery' }).click();
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/dr$/);
    await page.getByRole('link', { name: 'Performance replication' }).click();
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance$/);
    await page.getByRole('link', { name: 'Overview' }).click();
    await expect(page).toHaveURL(/\/ui\/vault\/replication$/);
  });

  await test.step('enables performance replication as primary', async () => {
    await page.goto('replication/performance');
    await page.locator('[data-test-replication-enable-form]').getByRole('button', { name: 'Enable' }).click();
    await expect(page.getByRole('heading', { name: 'Performance', exact: true })).toBeVisible({
      timeout: 60_000,
    });
    await expect(page.locator('[data-test-badge="primary"]')).toBeVisible();
  });

  await test.step('the mode tabs link to details, manage and secondaries', async () => {
    await page.locator('[data-test-replication-link="manage"]').click();
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance\/manage$/);
    await page.locator('[data-test-replication-link="secondaries"]').click();
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance\/secondaries$/);
    await page.locator('[data-test-replication-link="details"]').click();
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance$/);
  });

  await test.step('adds a secondary and returns to the secondaries list', async () => {
    await page.goto('replication/performance/secondaries');
    await page.locator('[data-test-secondary-add]').click();
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance\/secondaries\/add$/);
    await page.locator('[data-test-input="Secondary ID"]').fill('demo-secondary');
    await page.getByRole('button', { name: 'Generate token' }).click();
    await expect(page.getByText('Copy your token')).toBeVisible();
    await page.getByRole('button', { name: 'Copy token' }).click();
    await page.getByRole('button', { name: 'Close' }).click();
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance\/secondaries$/);
    await expect(page.locator('[data-test-secondary-name="demo-secondary"]')).toBeVisible();
  });

  await test.step('shows the path filter config page for a secondary', async () => {
    await page.goto('replication/performance/secondaries/config/show/demo-secondary');
    await expect(page).toHaveURL(
      /\/ui\/vault\/replication\/performance\/secondaries\/config\/show\/demo-secondary$/
    );
    await expect(page.getByText('No config yet')).toBeVisible();
  });

  await test.step('the revoke page asks for confirmation', async () => {
    await page.goto('replication/performance/secondaries/revoke');
    await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance\/secondaries\/revoke$/);
    await expect(page.getByRole('heading', { name: 'Revoke a secondary token' })).toBeVisible();
  });

  await test.step('disables performance replication', async () => {
    await page.goto('replication/performance/manage');
    await page.getByRole('button', { name: 'Disable Replication' }).click();
    await page.locator('[data-test-confirmation-modal-input="Disable Replication?"]').fill('Performance');
    await page.locator('[data-test-confirm-button="Disable Replication?"]').click();
    await expect(page.getByRole('heading', { name: 'Enable replication' })).toBeVisible({ timeout: 60_000 });
  });
});

test.describe('path filter configuration', () => {
  const mount = 'edm-filter-kv';
  const authMount = 'edm-filter-auth';
  const namespace = 'edm-filter-ns';
  const filterPath = '/v1/sys/replication/performance/primary/paths-filter';
  const secondaries = 'replication/performance/secondaries';
  const crudId = 'edm-filter-crud';
  const deniedCases = [
    { operation: 'update', mode: 'Deny', id: 'edm-filter-denied-update' },
    { operation: 'delete', mode: 'Include everything', id: 'edm-filter-denied-delete' },
  ];
  const tokenId = 'edm-filter-token';
  let api: APIRequestContext;

  test.describe.configure({ timeout: 120_000 });

  const personaHeaders = async (page: Page): Promise<Record<string, string>> => ({
    'X-Vault-Token': await page.evaluate(() => {
      const key = Object.keys(localStorage).find((key) => key.startsWith('vault-token'));
      const value = key && localStorage.getItem(key);
      if (!value) throw new Error('No stored persona session');
      const stored: unknown = JSON.parse(value);
      if (
        !stored ||
        typeof stored !== 'object' ||
        !('token' in stored) ||
        typeof stored.token !== 'string' ||
        !stored.token
      ) {
        throw new Error('Persona session has no token');
      }
      return stored.token;
    }),
  });

  const readFilter = async (page: Page, id: string): Promise<{ mode: string; paths: string[] }> => {
    const response = await page.request.get(`${filterPath}/${id}`, {
      headers: await personaHeaders(page),
    });
    await expect(response).toBeOK();
    return (await response.json()).data;
  };

  test.beforeAll(async ({ playwright }, testInfo) => {
    api = await rootApi(playwright, testInfo);
    await expect(await api.post('/v1/sys/replication/performance/primary/enable', { data: {} })).toBeOK();
    await expect
      .poll(
        async () => {
          const response = await api.get('/v1/sys/replication/status');
          await expect(response).toBeOK();
          return (await response.json()).data.performance.mode;
        },
        { timeout: 60_000 }
      )
      .toBe('primary');

    await ensureMount(api, mount, 'kv');
    await ensureNamespace(api, namespace);
    const mounts = await api.get('/v1/sys/internal/ui/mounts');
    await expect(mounts).toBeOK();
    if (!(await mounts.json()).data.auth[`${authMount}/`]) {
      await expect(await api.post(`/v1/sys/auth/${authMount}`, { data: { type: 'userpass' } })).toBeOK();
    }
    const headers = { 'X-Vault-Namespace': namespace };
    const namespacedMounts = await api.get('/v1/sys/internal/ui/mounts', { headers });
    await expect(namespacedMounts).toBeOK();
    if (!(await namespacedMounts.json()).data.secret['kv/']) {
      await expect(await api.post('/v1/sys/mounts/kv', { headers, data: { type: 'kv' } })).toBeOK();
    }
    for (const id of [crudId, ...deniedCases.map(({ id }) => id)]) {
      await expect(
        await api.post('/v1/sys/replication/performance/primary/secondary-token', { data: { id } })
      ).toBeOK();
      await expect(await api.delete(`${filterPath}/${id}`)).toBeOK();
      if (id !== crudId) {
        await expect(
          await api.post(`${filterPath}/${id}`, { data: { mode: 'allow', paths: [`${mount}/`] } })
        ).toBeOK();
      }
    }
  });

  test.afterAll(async () => {
    if (!api) return;
    try {
      await expect(await api.post('/v1/sys/replication/performance/primary/disable')).toBeOK();
      await expect(await api.delete(`/v1/sys/mounts/${mount}`)).toBeOK();
      await expect(await api.delete(`/v1/sys/auth/${authMount}`)).toBeOK();
      await expect(await api.delete(`/v1/sys/namespaces/${namespace}`)).toBeOK();
    } finally {
      await api.dispose();
    }
  });

  // A filter must retain the chosen mode and paths after saving, revisiting, and removing it.
  test('persists path filter create, edit and delete outcomes', async ({ page }) => {
    const basePage = new BasePage(page);

    await test.step('an absent filter offers creation and saves the selected root mount', async () => {
      await page.goto(`${secondaries}/config/show/${crudId}`);
      await expect(page.getByText('No config yet', { exact: true })).toBeVisible();
      await page.getByRole('link', { name: 'Create config', exact: true }).click();
      await page.getByRole('radio', { name: /^Allow/ }).check();
      await page.getByRole('combobox', { name: 'Paths in this allow' }).click();
      await page.getByRole('option', { name: `${mount}/`, exact: true }).click();
      const saved = page.waitForResponse(
        (response) =>
          response.url().endsWith(`${filterPath}/${crudId}`) &&
          ['PUT', 'POST'].includes(response.request().method())
      );
      await page.getByRole('button', { name: 'Create', exact: true }).click();
      const response = await saved;
      expect(response.ok()).toBe(true);
      expect(response.request().postDataJSON()).toEqual({ mode: 'allow', paths: [`${mount}/`] });
      const persisted = await readFilter(page, crudId);
      expect(persisted).toEqual({ mode: 'allow', paths: [`${mount}/`] });
      await test.info().attach('filter-create-readback', {
        body: JSON.stringify({ method: response.request().method(), status: response.status(), persisted }),
        contentType: 'application/json',
      });
      await expect(page).toHaveURL(new RegExp(`/config/show/${crudId}$`));
      await page.reload();
      await expect(page.locator('[data-test-mount-config-mode]')).toContainText('allow');
      await expect(page.locator('[data-test-mount-config-paths]')).toContainText(`${mount}/`);
      await highlightElement(page, page.locator('[data-test-mount-config-paths]'));
    });

    await test.step('creating an existing filter redirects to edit without losing its values', async () => {
      await page.goto(`${secondaries}/config/create/${crudId}`);
      await expect(page).toHaveURL(new RegExp(`/config/edit/${crudId}$`));
      await expect(page.getByRole('radio', { name: /^Allow/ })).toBeChecked();
      await expect(page.locator('[data-test-selected-option="0"]')).toContainText(`${mount}/`);
    });

    await test.step('an edit persists auth and namespace mount selections in deny mode', async () => {
      await page.goto(`${secondaries}/config/edit/${crudId}`);
      await basePage.dismissFlashMessages();
      await page.getByRole('radio', { name: /^Deny/ }).check();
      await page.getByRole('button', { name: 'delete selection', exact: true }).click();
      await page.getByRole('combobox', { name: 'Paths in this deny' }).click();
      await page.getByRole('option', { name: `auth/${authMount}/`, exact: true }).click();
      await page.getByRole('combobox', { name: 'Paths in this deny' }).click();
      await page.locator('[data-test-component="search-select"] input[type="search"]').fill(namespace);
      await page.getByRole('option', { name: `${namespace}/kv/`, exact: true }).click();
      const paths = [`auth/${authMount}/`, `${namespace}/kv/`];
      const saved = page.waitForResponse(
        (response) =>
          response.url().endsWith(`${filterPath}/${crudId}`) &&
          ['PUT', 'POST'].includes(response.request().method())
      );
      await page.getByRole('button', { name: 'Update', exact: true }).click();
      const response = await saved;
      expect(response.ok()).toBe(true);
      expect(response.request().postDataJSON()).toEqual({ mode: 'deny', paths });
      const persisted = await readFilter(page, crudId);
      expect(persisted.mode).toBe('deny');
      expect([...persisted.paths].sort()).toEqual([...paths].sort());
      await test.info().attach('filter-edit-readback', {
        body: JSON.stringify({ method: response.request().method(), status: response.status(), persisted }),
        contentType: 'application/json',
      });
      await expect(page).toHaveURL(new RegExp(`/config/show/${crudId}$`));
      await page.reload();
      await expect(page.locator('[data-test-mount-config-mode]')).toContainText('deny');
      for (const path of paths) {
        await expect(page.locator('[data-test-mount-config-paths]')).toContainText(path);
      }
      await highlightElement(page, page.locator('[data-test-mount-config-paths]'));
    });

    await test.step('Include everything removes the persisted configuration', async () => {
      await page.goto(`${secondaries}/config/edit/${crudId}`);
      await page.getByRole('radio', { name: /^Include everything/ }).check();
      await expect(page.locator('[data-test-remove-warning]')).toContainText('remove any existing filter');
      const deleted = page.waitForResponse(
        (response) =>
          response.url().endsWith(`${filterPath}/${crudId}`) && response.request().method() === 'DELETE'
      );
      await page.getByRole('button', { name: 'Update', exact: true }).click();
      expect((await deleted).ok()).toBe(true);
      await expect(page).toHaveURL(/\/replication\/performance\/secondaries$/);
      const readback = await page.request.get(`${filterPath}/${crudId}`, {
        headers: await personaHeaders(page),
      });
      expect(readback.status()).toBe(404);
      await test.info().attach('filter-delete-readback', {
        body: JSON.stringify({ status: readback.status() }),
        contentType: 'application/json',
      });
      await page.goto(`${secondaries}/config/show/${crudId}`);
      await page.reload();
      await expect(page.getByText('No config yet', { exact: true })).toBeVisible();
      await highlightElement(page, page.locator('[data-test-empty-state-title]'));
    });
  });

  test.describe('credential and permission outcomes', () => {
    // Secondary-token creation has its own filter write, separate from the configuration routes.
    test('secondary token creation persists its path filter', async ({ page }) => {
      await page.goto(`${secondaries}/add`);
      await page.getByLabel('Secondary ID', { exact: true }).fill(tokenId);
      await page.getByRole('radio', { name: /^Allow/ }).check();
      await page.getByRole('combobox', { name: 'Paths in this allow' }).click();
      await page.getByRole('option', { name: `${mount}/`, exact: true }).click();
      const saved = page.waitForResponse(
        (response) =>
          response.url().endsWith(`${filterPath}/${tokenId}`) &&
          ['PUT', 'POST'].includes(response.request().method())
      );
      await page.getByRole('button', { name: 'Generate token', exact: true }).click();
      const response = await saved;
      expect(response.ok()).toBe(true);
      expect(response.request().postDataJSON()).toEqual({ mode: 'allow', paths: [`${mount}/`] });
      await expect(page.getByRole('heading', { name: 'Copy your token' })).toBeVisible();
      const persisted = await readFilter(page, tokenId);
      expect(persisted).toEqual({ mode: 'allow', paths: [`${mount}/`] });
      await test.info().attach('secondary-token-filter-readback', {
        body: JSON.stringify({ status: response.status(), persisted }),
        contentType: 'application/json',
      });
      await page.goto(`${secondaries}/config/show/${tokenId}`);
      await page.reload();
      await expect(page.locator('[data-test-mount-config-mode]')).toContainText('allow');
      await expect(page.locator('[data-test-mount-config-paths]')).toContainText(`${mount}/`);
    });

    // A token that can open the form must not change a filter without write/delete permission.
    for (const { operation, mode, id } of deniedCases) {
      test(`a denied path filter ${operation} preserves the saved configuration`, async ({ page }) => {
        const policyName = `${id}-readonly`;
        await expect(
          await api.put(`/v1/sys/policies/acl/${policyName}`, {
            data: {
              policy: `
              path "*" { capabilities = ["read", "list"] }
              path "sys/capabilities-self" { capabilities = ["update"] }
              path "sys/replication/performance/primary/secondary-token" { capabilities = ["update", "sudo"] }
              path "sys/replication/performance/primary/paths-filter/*" { capabilities = ["read"] }
            `,
            },
          })
        ).toBeOK();
        const created = await api.post('/v1/auth/token/create', {
          data: { policies: [policyName], ttl: '15m' },
        });
        await expect(created).toBeOK();
        const { auth } = await created.json();
        try {
          await page.goto('dashboard');
          await page.evaluate(
            ({ token, policies }) => {
              const key = Object.keys(localStorage).find((key) => key.startsWith('vault-token'));
              const value = key && localStorage.getItem(key);
              if (!key || !value) throw new Error('No persona session to replace');
              localStorage.setItem(key, JSON.stringify({ ...JSON.parse(value), token, policies }));
            },
            { token: auth.client_token, policies: auth.policies }
          );
          const headers = await personaHeaders(page);
          const lookup = await page.request.get('/v1/auth/token/lookup-self', { headers });
          await expect(lookup).toBeOK();
          const { data } = await lookup.json();
          expect(data.policies).toContain(policyName);
          expect(data.policies).not.toContain('root');

          await test.step(`a denied ${operation} leaves the filter unchanged`, async () => {
            await page.goto(`${secondaries}/config/edit/${id}`);
            await expect(page.getByRole('button', { name: 'Update', exact: true })).toBeVisible();
            await page.getByRole('radio', { name: new RegExp(`^${mode}`) }).check();
            const rejected = page.waitForResponse(
              (response) =>
                response.url().endsWith(`${filterPath}/${id}`) &&
                ['PUT', 'POST', 'DELETE'].includes(response.request().method())
            );
            await page.getByRole('button', { name: 'Update', exact: true }).click();
            const response = await rejected;
            expect(response.status()).toBe(403);
            const persisted = await readFilter(page, id);
            expect(persisted).toEqual({ mode: 'allow', paths: [`${mount}/`] });
            await test.info().attach(`denied-filter-${response.request().method().toLowerCase()}`, {
              body: JSON.stringify({ status: response.status(), persisted }),
              contentType: 'application/json',
            });
            await expect(
              page.getByText(
                new RegExp(
                  `There was an error ${
                    operation === 'delete' ? 'deleting' : 'saving'
                  } the config for ${id}:.*permission denied`,
                  'is'
                )
              )
            ).toBeVisible();
            await page.reload();
            await expect(page.getByRole('radio', { name: /^Allow/ })).toBeChecked();
            await expect(page.locator('[data-test-selected-option="0"]')).toContainText(`${mount}/`);
          });
        } finally {
          await expect(
            await api.post('/v1/auth/token/revoke', { data: { token: auth.client_token } })
          ).toBeOK();
          await expect(await api.delete(`/v1/sys/policies/acl/${policyName}`)).toBeOK();
        }
      });
    }
  });
});

// Adding a secondary needs sudo on secondary-token, while revoking needs only update on revoke-secondary.
// Each token's ACL must decide which secondary actions the pages offer and which routes they allow.
test.describe('secondary capability gating', () => {
  const mount = 'edm-caps-kv';
  const secondaries = 'replication/performance/secondaries';
  const filterPath = '/v1/sys/replication/performance/primary/paths-filter';
  const filteredId = 'edm-caps-filtered';
  const unfilteredId = 'edm-caps-unfiltered';
  const basePolicy = `
    path "*" { capabilities = ["read", "list"] }
    path "sys/capabilities-self" { capabilities = ["update"] }
  `;
  const cases = [
    {
      name: 'a token without secondary permissions cannot add or revoke secondaries',
      policyName: 'edm-caps-none',
      policy: basePolicy,
      canAdd: false,
      canRevoke: false,
    },
    {
      name: 'a token without sudo on secondary tokens can only revoke secondaries',
      policyName: 'edm-caps-no-sudo',
      policy: `${basePolicy}
        path "sys/replication/performance/primary/secondary-token" { capabilities = ["update"] }
        path "sys/replication/performance/primary/revoke-secondary" { capabilities = ["update"] }
      `,
      canAdd: false,
      canRevoke: true,
    },
    {
      name: 'a non-root token with secondary permissions can add and revoke secondaries',
      policyName: 'edm-caps-full',
      policy: `${basePolicy}
        path "sys/replication/performance/primary/secondary-token" { capabilities = ["update", "sudo"] }
        path "sys/replication/performance/primary/revoke-secondary" { capabilities = ["update"] }
      `,
      canAdd: true,
      canRevoke: true,
    },
  ];
  let api: APIRequestContext;

  test.describe.configure({ timeout: 120_000 });

  const readKnownSecondaries = async (): Promise<string[]> => {
    const response = await api.get('/v1/sys/replication/performance/status');
    await expect(response).toBeOK();
    return [...(await response.json()).data.known_secondaries].sort();
  };

  test.beforeAll(async ({ playwright }, testInfo) => {
    api = await rootApi(playwright, testInfo);
    await expect(await api.post('/v1/sys/replication/performance/primary/enable', { data: {} })).toBeOK();
    await expect
      .poll(
        async () => {
          const response = await api.get('/v1/sys/replication/status');
          await expect(response).toBeOK();
          return (await response.json()).data.performance.mode;
        },
        { timeout: 60_000 }
      )
      .toBe('primary');

    await ensureMount(api, mount, 'kv');
    for (const id of [filteredId, unfilteredId]) {
      await expect(
        await api.post('/v1/sys/replication/performance/primary/secondary-token', { data: { id } })
      ).toBeOK();
    }
    await expect(
      await api.post(`${filterPath}/${filteredId}`, { data: { mode: 'allow', paths: [`${mount}/`] } })
    ).toBeOK();
  });

  test.afterAll(async () => {
    if (!api) return;
    try {
      await expect(await api.post('/v1/sys/replication/performance/primary/disable')).toBeOK();
      await expect(await api.delete(`/v1/sys/mounts/${mount}`)).toBeOK();
    } finally {
      await api.dispose();
    }
  });

  for (const { name, policyName, policy, canAdd, canRevoke } of cases) {
    test(name, async ({ page }) => {
      await expect(await api.put(`/v1/sys/policies/acl/${policyName}`, { data: { policy } })).toBeOK();
      const created = await api.post('/v1/auth/token/create', {
        data: { policies: [policyName], ttl: '15m' },
      });
      await expect(created).toBeOK();
      const { auth } = await created.json();
      try {
        await page.goto('dashboard');
        const sessionToken = await page.evaluate(
          ({ token, policies }) => {
            const key = Object.keys(localStorage).find((key) => key.startsWith('vault-token'));
            const value = key && localStorage.getItem(key);
            if (!key || !value) throw new Error('No persona session to replace');
            localStorage.setItem(key, JSON.stringify({ ...JSON.parse(value), token, policies }));
            return JSON.parse(localStorage.getItem(key) ?? '{}').token;
          },
          { token: auth.client_token, policies: auth.policies }
        );
        expect(sessionToken).toBe(auth.client_token);
        const lookup = await page.request.get('/v1/auth/token/lookup-self', {
          headers: { 'X-Vault-Token': auth.client_token },
        });
        await expect(lookup).toBeOK();
        const { data } = await lookup.json();
        expect(data.policies).toContain(policyName);
        expect(data.policies).not.toContain('root');

        await test.step(`the dashboard ${canAdd ? 'offers' : 'hides'} adding a secondary`, async () => {
          await page.goto('replication/performance');
          await expect(page.locator('[data-test-badge="primary"]')).toBeVisible();
          const knownSecondariesCard = page.getByRole('heading', { name: /Known secondaries$/ });
          await expect(knownSecondariesCard).toBeVisible();
          const addLink = page.locator('[data-test-link-to="Add secondary"]');
          if (canAdd) {
            await expect(addLink).toBeVisible();
            await highlightElement(page, addLink);
          } else {
            await expect(addLink).toHaveCount(0);
            await highlightElement(page, knownSecondariesCard);
          }
        });

        await test.step('the secondaries list offers only the permitted actions', async () => {
          await page.goto(secondaries);
          for (const id of [filteredId, unfilteredId]) {
            await expect(page.locator(`[data-test-secondary-name="${id}"]`)).toBeVisible();
          }
          const addLink = page.locator('[data-test-secondary-add]');
          const revokeLink = page.getByRole('link', { name: 'Revoke secondary', exact: true });
          await (canAdd ? expect(addLink).toBeVisible() : expect(addLink).toHaveCount(0));
          await (canRevoke ? expect(revokeLink).toBeVisible() : expect(revokeLink).toHaveCount(0));
          const row = page.locator(`[data-test-secondary-name="${filteredId}"]`);
          await highlightElement(page, row);
          await row.getByRole('button', { name: 'Secondary popup nav menu' }).click();
          await expect(page.getByRole('link', { name: 'Path filter config', exact: true })).toBeVisible();
          const revokeItem = page.getByRole('button', { name: 'Revoke', exact: true });
          await (canRevoke ? expect(revokeItem).toBeVisible() : expect(revokeItem).toHaveCount(0));
        });

        await test.step(`the add page ${canAdd ? 'opens' : 'redirects to the mode details'}`, async () => {
          await page.goto(`${secondaries}/add`);
          if (canAdd) {
            await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance\/secondaries\/add$/);
            await expect(page.getByLabel('Secondary ID', { exact: true })).toBeVisible();
          } else {
            await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance$/);
            await expect(page.locator('[data-test-badge="primary"]')).toBeVisible();
          }
        });

        await test.step(`the revoke page ${canRevoke ? 'opens' : 'redirects to the overview'}`, async () => {
          await page.goto(`${secondaries}/revoke`);
          if (canRevoke) {
            await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance\/secondaries\/revoke$/);
            await expect(page.getByRole('heading', { name: 'Revoke a secondary token' })).toBeVisible();
          } else {
            await expect(page).toHaveURL(/\/ui\/vault\/replication$/);
            await expect(page.getByRole('heading', { name: 'Replication', exact: true })).toBeVisible();
          }
        });

        await test.step(`path filter config ${
          canAdd ? 'opens' : 'redirects to the mode details'
        }`, async () => {
          await page.goto(`${secondaries}/config/create/${unfilteredId}`);
          if (canAdd) {
            await expect(page).toHaveURL(new RegExp(`/config/create/${unfilteredId}$`));
            await expect(page.getByRole('button', { name: 'Create', exact: true })).toBeVisible();
          } else {
            await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance$/);
            await expect(page.locator('[data-test-badge="primary"]')).toBeVisible();
          }
          await page.goto(`${secondaries}/config/edit/${filteredId}`);
          if (canAdd) {
            await expect(page).toHaveURL(new RegExp(`/config/edit/${filteredId}$`));
            await expect(page.getByRole('radio', { name: /^Allow/ })).toBeChecked();
            await expect(page.getByRole('button', { name: 'Update', exact: true })).toBeVisible();
          } else {
            await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance$/);
            await expect(page.locator('[data-test-badge="primary"]')).toBeVisible();
          }
        });

        await test.step('the manage page still offers every primary action', async () => {
          await page.goto('replication/performance/manage');
          await expect(page).toHaveURL(/\/ui\/vault\/replication\/performance\/manage$/);
          for (const action of ['Disable Replication', 'Demote', 'Recover', 'Reindex']) {
            await expect(page.getByRole('button', { name: action, exact: true })).toBeVisible();
          }
          await highlightElement(page, page.locator('[data-test-demote-replication]'));
        });

        await test.step('the known secondaries are unchanged', async () => {
          const knownSecondaries = await readKnownSecondaries();
          expect(knownSecondaries).toEqual([filteredId, unfilteredId].sort());
          await test.info().attach('known-secondaries-readback', {
            body: JSON.stringify({ policy: policyName, knownSecondaries }),
            contentType: 'application/json',
          });
        });
      } finally {
        await expect(
          await api.post('/v1/auth/token/revoke', { data: { token: auth.client_token } })
        ).toBeOK();
        await expect(await api.delete(`/v1/sys/policies/acl/${policyName}`)).toBeOK();
      }
    });
  }
});
