/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { highlightText, holdForVideo, test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

const AWS_MOUNT = 'ttl-max-aws';
const AWS_ROLE = 'federated';
const TTL_ERROR = 'TTL is too large. Maximum allowed value is 106,751 days.';

// Setup and cleanup live in hooks so the mount is removed even when the test fails or retries.
test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, AWS_MOUNT, 'aws');
  // A role is enough to reach the credentials form; no AWS account is contacted.
  await expect(
    await api.post(`/v1/${AWS_MOUNT}/roles/${AWS_ROLE}`, {
      data: {
        credential_type: 'federation_token',
        policy_document:
          '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:ListBucket","Resource":"*"}]}',
      },
    })
  ).toBeOK();
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/mounts/${AWS_MOUNT}`);
  await api.dispose();
});

test('the ttl picker flags values longer than Vault supports', async ({ page }) => {
  await page.goto(`secrets-engines/${AWS_MOUNT}/credentials/${AWS_ROLE}`);
  const ttl = page.getByRole('group', { name: 'TTL' });
  const units = ttl.getByLabel('Number of units');

  // Clicking first moves the recorded cursor to the input; fill alone types without it.
  await units.click();
  await units.fill('9999999999');
  await ttl.getByRole('combobox').selectOption('d');
  await expect(page.getByText(TTL_ERROR)).toBeVisible();
  await highlightText(page, page.locator('form'), TTL_ERROR);

  await units.click();
  await units.fill('1');
  await expect(page.getByText(TTL_ERROR)).toBeHidden();
  await holdForVideo(page);
});
