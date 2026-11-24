/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';

const ACL_POLICY = 'path "secret/*" { capabilities = ["read"] }';
const SENTINEL_POLICY = Buffer.from('main = rule { true }').toString('base64');

// The persona can read policies named readable-* and only list the others (see readonly.hcl).
const POLICY_TYPES = [
  { type: 'acl', extension: 'hcl', data: { policy: ACL_POLICY } },
  {
    type: 'egp',
    extension: 'sentinel',
    data: { policy: SENTINEL_POLICY, enforcement_level: 'soft-mandatory', paths: ['*'] },
  },
  {
    type: 'rgp',
    extension: 'sentinel',
    data: { policy: SENTINEL_POLICY, enforcement_level: 'soft-mandatory' },
  },
];

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  for (const { type, data } of POLICY_TYPES) {
    for (const name of [`readable-${type}`, `hidden-${type}`]) {
      await expect(await api.post(`/v1/sys/policies/${type}/${name}`, { data })).toBeOK();
    }
  }
  await api.dispose();
});

// Downloading reads the policy body, so it must not be offered for policies the token cannot read.
for (const { type, extension } of POLICY_TYPES) {
  test(`read-only user can only download ${type.toUpperCase()} policies it can read`, async ({ page }) => {
    await page.goto(`policies/${type}`);
    const row = (name: string) => page.getByRole('row').filter({ hasText: name });
    await expect(row(`hidden-${type}`)).toBeVisible();

    await test.step('an unreadable policy has no row menu', async () => {
      await expect(row(`hidden-${type}`).getByRole('button', { name: 'Policy options' })).toHaveCount(0);
    });

    await test.step('a readable policy downloads', async () => {
      await row(`readable-${type}`).getByRole('button', { name: 'Policy options' }).click();
      const download = page.waitForEvent('download');
      await page.getByRole('button', { name: 'Download policy' }).click();
      expect((await download).suggestedFilename()).toBe(`readable-${type}.${extension}`);
    });
  });
}
