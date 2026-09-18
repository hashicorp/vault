/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

// The raft persona boots real Raft storage; this page does not render for inmem storage.
test('raft storage overview renders the single-node cluster', async ({ page }) => {
  await page.goto('dashboard');
  await page.getByRole('link', { name: 'Raft storage' }).click();

  await expect(page.getByRole('heading', { name: 'Raft storage' })).toBeVisible();

  await test.step('renders the current node as leader and voter', async () => {
    const rows = page.locator('[data-raft-row]');
    await expect(rows).toHaveCount(1);
    await expect(rows.getByText(/127\.0\.0\.1:\d+/)).toBeVisible();
    // the single node in a freshly bootstrapped cluster is both the leader and a voter
    await expect(rows.getByText('Yes')).toHaveCount(2);
  });
});
