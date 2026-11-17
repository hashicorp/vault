/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { highlightText, holdForVideo, test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

import type { Page } from '@playwright/test';

const AUTH = 'ttl-max-core-userpass';
const TOO_LARGE = (label: string) => `${label} is too large. Maximum allowed value is 106,751 days.`;
const LEASE_ID = 'ttl-max-core/renewable/lease1234';
const K8S = 'ttl-max-core-k8s';
const K8S_ROLE = 'ttl-max-role';

// Setup and cleanup live in hooks so the mounts are removed even when a test fails or retries.
test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  if (!(await api.get(`/v1/sys/auth/${AUTH}/tune`)).ok()) {
    await expect(await api.post(`/v1/sys/auth/${AUTH}`, { data: { type: 'userpass' } })).toBeOK();
  }
  // No Kubernetes cluster exists; the config and role only let the credentials form load.
  await ensureMount(api, K8S, 'kubernetes');
  await expect(
    await api.post(`/v1/${K8S}/config`, {
      data: { kubernetes_host: 'https://127.0.0.1:1', disable_local_ca_jwt: true },
    })
  ).toBeOK();
  await expect(
    await api.post(`/v1/${K8S}/roles/${K8S_ROLE}`, {
      data: { allowed_kubernetes_namespaces: '*', service_account_name: 'sa' },
    })
  ).toBeOK();
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/auth/${AUTH}`);
  await api.delete(`/v1/sys/mounts/${K8S}`);
  await api.dispose();
});

// The container of a picker's error message, so a recording highlights only that text.
const errorOf = (page: Page, label: string) => page.getByText(TOO_LARGE(label)).locator('..');

// Records the requests a page sends, so a test can prove nothing was saved.
// `bodies` holds the JSON payload of each recorded request, for checking what was sent.
const recordWrites = (page: Page, pattern: RegExp) => {
  const writes: string[] = [];
  const bodies: Array<Record<string, unknown>> = [];
  page.on('request', (request) => {
    if (request.method() !== 'GET' && pattern.test(request.url())) {
      writes.push(`${request.method()} ${new URL(request.url()).pathname}`);
      bodies.push(request.postDataJSON() ?? {});
    }
  });
  return { writes, bodies };
};

// Pickers are found by their label with the same data-test attributes the QUnit suite uses, because a
// toggled-off picker has no accessible group name.
const turnOn = (page: Page, label: string) => page.locator(`[data-test-toggle-label="${label}"]`).click();

// Clicking first moves the recorded cursor to the input; fill alone types without it.
const fillTtl = async (page: Page, label: string, value: string, unit?: string) => {
  const input = page.locator(`[data-test-ttl-value="${label}"]`);
  await input.click();
  await input.fill(value);
  if (unit) await page.locator(`[data-test-ttl-unit="${label}"] select`).selectOption(unit);
};

// An over-max TTL should be explained where it was typed and must never reach the API.
test('tuning an auth method blocks durations longer than Vault supports', async ({ page }) => {
  const { writes, bodies } = recordWrites(page, /\/tune$/);
  const open = async () => {
    await page.goto(`settings/auth/configure/${AUTH}/options`);
    await expect(page.getByRole('button', { name: 'Update options' })).toBeVisible();
  };

  await test.step('an over-max default lease TTL is flagged', async () => {
    await open();
    await turnOn(page, 'Default Lease TTL');
    await fillTtl(page, 'Default Lease TTL', '9999999999', 'd');
    await page.getByRole('button', { name: 'Update options' }).click();

    await expect(page.getByText(TOO_LARGE('Default Lease TTL'))).toBeVisible();
    await expect(page.getByText('There is an error with this form.')).toBeVisible();
    await highlightText(page, errorOf(page, 'Default Lease TTL'), TOO_LARGE('Default Lease TTL'));
    expect(writes).toEqual([]);
  });

  await test.step('an over-max lockout duration is flagged', async () => {
    await open();
    await turnOn(page, 'Lockout duration');
    await fillTtl(page, 'Lockout duration', '9999999999', 'd');
    await page.getByRole('button', { name: 'Update options' }).click();

    await expect(page.getByText(TOO_LARGE('Lockout duration'))).toBeVisible();
    await highlightText(page, errorOf(page, 'Lockout duration'), TOO_LARGE('Lockout duration'));
    expect(writes).toEqual([]);
  });

  await test.step('an over-max value turned off again does not block saving', async () => {
    await open();
    await turnOn(page, 'Lockout duration');
    await fillTtl(page, 'Lockout duration', '9999999999', 'd');
    await turnOn(page, 'Lockout duration');
    await page.getByRole('button', { name: 'Update options' }).click();

    await expect.poll(() => writes).toEqual([`POST /v1/sys/mounts/auth/${AUTH}/tune`]);
    expect(bodies[0]?.user_lockout_config).toEqual({ lockout_duration: '0' });
    await holdForVideo(page);
  });
});

test('wrapping data blocks a TTL longer than Vault supports', async ({ page }) => {
  const { writes } = recordWrites(page, /wrapping\/wrap$/);
  await page.goto('tools/wrap');

  await turnOn(page, 'Wrap TTL');
  await fillTtl(page, 'Wrap TTL', '9999999999', 'd');
  await page.getByRole('button', { name: 'Wrap data' }).click();

  await expect(page.getByText(TOO_LARGE('Wrap TTL'))).toBeVisible();
  await highlightText(page, errorOf(page, 'Wrap TTL'), TOO_LARGE('Wrap TTL'));
  expect(writes).toEqual([]);

  await fillTtl(page, 'Wrap TTL', '30', 'm');
  await page.getByRole('button', { name: 'Wrap data' }).click();
  await expect.poll(() => writes).toEqual(['POST /v1/sys/wrapping/wrap']);
  await holdForVideo(page);
});

test('renewing a lease blocks an increment longer than Vault supports', async ({ page }) => {
  // The test server has no source of renewable leases (they need `vault server -dev-leased-kv`), so only the
  // lease lookup and renew calls are stubbed; the page and the request it sends are the real ones.
  await page.route('**/v1/sys/leases/lookup', (route) =>
    route.fulfill({
      json: {
        data: {
          id: LEASE_ID,
          renewable: true,
          issue_time: '2026-01-01T00:00:00.000Z',
          expire_time: '2026-01-02T00:00:00.000Z',
          last_renewal: null,
          ttl: 3600,
        },
      },
    })
  );
  await page.route('**/v1/sys/leases/renew', (route) =>
    route.fulfill({ json: { lease_id: LEASE_ID, renewable: true, lease_duration: 3600 } })
  );
  const { writes, bodies } = recordWrites(page, /leases\/renew$/);
  const open = async () => {
    await page.goto(`access/leases/show/${LEASE_ID}`);
    await expect(page.getByRole('button', { name: 'Renew lease' })).toBeVisible();
  };

  await test.step('an over-max increment is flagged and the lease is not renewed', async () => {
    await open();
    await turnOn(page, 'Increment');
    await fillTtl(page, 'Increment', '9999999999', 'd');
    await page.getByRole('button', { name: 'Renew lease' }).click();

    await expect(page.getByText(TOO_LARGE('Increment'))).toBeVisible();
    await highlightText(page, errorOf(page, 'Increment'), TOO_LARGE('Increment'));
    expect(writes).toEqual([]);
  });

  await test.step('an over-max increment turned off again is not sent', async () => {
    await open();
    await turnOn(page, 'Increment');
    await fillTtl(page, 'Increment', '9999999999', 'd');
    await turnOn(page, 'Increment');
    await page.getByRole('button', { name: 'Renew lease' }).click();

    await expect(page.getByText('was successfully renewed.')).toBeVisible();
    expect(writes).toEqual(['POST /v1/sys/leases/renew']);
    expect(bodies[0]?.increment).toBe('');
    await holdForVideo(page);
  });
});

// A TTL that was turned off keeps being sent as it always was; only one Vault would reject is replaced,
// so the role's default is used instead.
test('kubernetes credentials use the role default when an over-limit TTL is turned off', async ({ page }) => {
  const { writes, bodies } = recordWrites(page, new RegExp(`/${K8S}/creds/${K8S_ROLE}$`));
  // The configured cluster does not exist, so Vault would answer with a connection error. Return the
  // credentials a real cluster would, so the recording ends on the page a user would see.
  await page.route(new RegExp(`/${K8S}/creds/${K8S_ROLE}$`), (route) =>
    route.fulfill({
      json: {
        lease_id: `${K8S}/creds/${K8S_ROLE}/lease1234`,
        lease_duration: 3600,
        renewable: false,
        data: {
          service_account_name: 'sa',
          service_account_namespace: 'default',
          service_account_token: 'example-token',
        },
      },
    })
  );
  const open = async () => {
    await page.goto(`secrets-engines/${K8S}/kubernetes/roles/${K8S_ROLE}/credentials`);
    await page.getByLabel('Kubernetes namespace').fill('default');
  };

  await test.step('an over-max TTL is flagged and nothing is requested', async () => {
    await open();
    await turnOn(page, 'Time-to-Live (TTL)');
    await fillTtl(page, 'Time-to-Live (TTL)', '9999999999', 'd');
    await page.getByRole('button', { name: 'Generate credentials' }).click();

    await expect(page.getByText(TOO_LARGE('Time-to-Live (TTL)'))).toBeVisible();
    await highlightText(page, errorOf(page, 'Time-to-Live (TTL)'), TOO_LARGE('Time-to-Live (TTL)'));
    expect(writes).toEqual([]);
  });

  await test.step('turned off again, the role default is used', async () => {
    await turnOn(page, 'Time-to-Live (TTL)');
    await page.getByRole('button', { name: 'Generate credentials' }).click();

    await expect.poll(() => writes).toEqual([`POST /v1/${K8S}/creds/${K8S_ROLE}`]);
    expect(bodies[0]?.ttl).toBe('');
    await expect(page.getByText('Lease ID')).toBeVisible();
    await holdForVideo(page);
  });

  await test.step('a valid TTL that was turned off is still sent as before', async () => {
    await open();
    await turnOn(page, 'Time-to-Live (TTL)');
    await fillTtl(page, 'Time-to-Live (TTL)', '2', 'h');
    await turnOn(page, 'Time-to-Live (TTL)');
    await page.getByRole('button', { name: 'Generate credentials' }).click();

    await expect.poll(() => writes.length).toBe(2);
    expect(bodies[1]?.ttl).toBe('2h');
    await expect(page.getByText('Lease ID')).toBeVisible();
    await holdForVideo(page);
  });
});
