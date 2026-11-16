/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import type { PlaywrightWorkerArgs, TestInfo } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';

const MODES = ['dr', 'performance'];

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
