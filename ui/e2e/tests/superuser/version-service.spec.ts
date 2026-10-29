/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, test } from '../../fixtures/demo';

test('version service reflects backend metadata and licensed features', async ({ page }) => {
  const healthResponsePromise = page.waitForResponse(
    (response) => new URL(response.url()).pathname === '/v1/sys/health'
  );
  const sealStatusResponsePromise = page.waitForResponse(
    (response) => new URL(response.url()).pathname === '/v1/sys/seal-status'
  );
  const licenseFeaturesResponsePromise = page.waitForResponse(
    (response) => new URL(response.url()).pathname === '/v1/sys/license/features'
  );

  await page.goto('dashboard');

  const healthResponse = await healthResponsePromise;
  const sealStatusResponse = await sealStatusResponsePromise;
  const licenseFeaturesResponse = await licenseFeaturesResponsePromise;

  expect(healthResponse.status()).toBe(200);
  expect(sealStatusResponse.status()).toBe(200);
  expect(licenseFeaturesResponse.status()).toBe(200);

  const health = await healthResponse.json();
  const sealStatus = await sealStatusResponse.json();
  const licenseFeatures = await licenseFeaturesResponse.json();
  const isEnterprise = health.has_chroot_namespace || health.enterprise;
  const displayedVersion = isEnterprise ? sealStatus.version.split('+')[0] : sealStatus.version;

  expect(isEnterprise, 'the Enterprise backend reports its edition').toBe(true);
  expect(sealStatus.cluster_name, 'seal status returns the cluster name').toBeTruthy();
  await expect(page.getByRole('heading', { name: `Vault v${displayedVersion}` })).toBeVisible();
  expect(licenseFeatures.features).toContain('Namespaces');

  await page.goto('access/namespaces/create');
  await expect(page.getByRole('heading', { name: 'Create a Namespace' })).toBeVisible();
  await expect(page.getByRole('textbox', { name: 'Path' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Save' })).toBeVisible();
});
