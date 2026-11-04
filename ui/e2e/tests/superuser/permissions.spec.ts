/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, type Page } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureNamespace, rootApi } from '../../fixtures/root-api';

const NAMESPACE = 'permissions-ns';
const RESULTANT_ACL_URL = '/v1/sys/internal/ui/resultant-acl';

const waitForResultantAcl = (page: Page) =>
  page.waitForResponse(
    (response) => response.url().includes(RESULTANT_ACL_URL) && response.request().method() === 'GET'
  );

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureNamespace(api, NAMESPACE);
  await api.dispose();
});

// The sidebar and namespace banner are driven by the resultant-acl payload, so a token whose
// policy grants everything must see every access-control link and no banner, even in a child
// namespace, where the ACL is still read from the token's own root namespace.
test('superuser resultant ACL grants full access-control navigation in a child namespace', async ({
  page,
}) => {
  const aclResponse = waitForResultantAcl(page);
  await page.goto(`dashboard?namespace=${NAMESPACE}`);

  await test.step('resultant ACL is read from the user root namespace', async () => {
    const response = await aclResponse;
    expect(response.status()).toBe(200);
    const headers = await response.request().allHeaders();
    expect(headers['x-vault-namespace'] ?? '').toBe('');
    const { data } = await response.json();
    expect(data.root).toBe(false);
    expect(Object.keys(data.glob_paths)).toContain('');
  });

  await test.step('no namespace access banner is shown', async () => {
    await expect(page.getByText('You do not have access to this namespace')).toHaveCount(0);
    await expect(page.getByText('Resultant ACL check failed')).toHaveCount(0);
  });

  await test.step('every access-control link is shown', async () => {
    await page.getByRole('link', { name: 'Access control' }).click();
    for (const name of [
      'ACL policies',
      'Leases',
      'Authentication methods',
      'Multi-factor authentication',
      'OIDC provider',
      'Namespaces',
      'Groups',
      'Entities',
    ]) {
      await expect(page.getByRole('link', { name, exact: true })).toBeVisible();
    }
  });
});
