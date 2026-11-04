/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, type Page } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureNamespace, rootApi } from '../../fixtures/root-api';

const NAMESPACE = 'readonly-permissions-ns';
const RESULTANT_ACL_URL = '/v1/sys/internal/ui/resultant-acl';
const NO_ACCESS_TITLE = 'You do not have access to this namespace';

const waitForResultantAcl = (page: Page) =>
  page.waitForResponse(
    (response) => response.url().includes(RESULTANT_ACL_URL) && response.request().method() === 'GET'
  );

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureNamespace(api, NAMESPACE);
  await api.dispose();
});

// The read-only policy grants identity and namespace reads but nothing under sys/auth, so the
// resultant ACL must show only the links it covers.
test('read-only resultant ACL limits access-control navigation', async ({ page }) => {
  const aclResponse = waitForResultantAcl(page);
  await page.goto('dashboard');

  await test.step('resultant ACL is loaded for the read-only token', async () => {
    const response = await aclResponse;
    expect(response.status()).toBe(200);
    const { data } = await response.json();
    expect(data.root).toBe(false);
    expect(Object.keys(data.glob_paths)).toContain('identity/');
    expect(Object.keys(data.exact_paths)).not.toContain('sys/auth');
  });

  await test.step('no namespace access banner in the root namespace', async () => {
    await expect(page.getByText(NO_ACCESS_TITLE)).toHaveCount(0);
  });

  await test.step('only links covered by the policy are shown', async () => {
    await page.getByRole('link', { name: 'Access control' }).click();
    for (const name of ['ACL policies', 'Namespaces', 'Groups', 'Entities']) {
      await expect(page.getByRole('link', { name, exact: true })).toBeVisible();
    }
    await expect(page.getByRole('link', { name: 'Authentication methods', exact: true })).toHaveCount(0);
  });
});

// The read-only policy lives in the root namespace and grants nothing inside a child namespace,
// so the ACL, still read from the token's root namespace, must produce the no-access banner.
test('read-only token sees the no-access banner in a namespace its policy does not cover', async ({
  page,
}) => {
  const aclResponse = waitForResultantAcl(page);
  await page.goto(`dashboard?namespace=${NAMESPACE}`);

  await test.step('resultant ACL is read from the user root namespace', async () => {
    const response = await aclResponse;
    expect(response.status()).toBe(200);
    const headers = await response.request().allHeaders();
    expect(headers['x-vault-namespace'] ?? '').toBe('');
  });

  await test.step('no-access banner is shown', async () => {
    await expect(page.getByText(NO_ACCESS_TITLE)).toBeVisible();
  });
});
