/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { highlightText, holdForVideo, test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';

const AUTH = 'lockout-config-userpass';

// Setup and cleanup live in hooks so the mount is removed even when a test fails or retries.
test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/auth/${AUTH}`);
  await expect(await api.post(`/v1/sys/auth/${AUTH}`, { data: { type: 'userpass' } })).toBeOK();
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/auth/${AUTH}`);
  await api.dispose();
});

const label = (name: string) => ({
  toggle: `[data-test-toggle-label="${name}"]`,
  value: `[data-test-ttl-value="${name}"]`,
  unit: `[data-test-ttl-unit="${name}"] select`,
});

// Clicking first moves the recorded cursor to the input; fill alone types without it.
const enterDuration = async (
  page: import('@playwright/test').Page,
  name: string,
  value: string,
  unit: string
) => {
  const field = label(name);
  await page.locator(field.toggle).click();
  await page.locator(field.value).click();
  await page.locator(field.value).fill(value);
  await page.locator(field.unit).selectOption(unit);
};

// The counter reset field used to send a key Vault does not read, so it was never stored.
test('lockout settings are stored by Vault and shown when the page is opened again', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const options = `settings/auth/configure/${AUTH}/options`;
  const tune = async () => (await (await api.get(`/v1/sys/auth/${AUTH}/tune`)).json()).data;

  // Disposed in a finally block so the request context is released even when a step fails.
  try {
    await test.step('save a threshold, a lockout duration and a counter reset', async () => {
      await page.goto(options);
      await page.locator('[data-test-input="user_lockout_config.lockout_threshold"]').click();
      await page.locator('[data-test-input="user_lockout_config.lockout_threshold"]').fill('7');
      await enterDuration(page, 'Lockout duration', '10', 'm');
      await enterDuration(page, 'Lockout counter reset', '30', 'm');
      await page.getByRole('button', { name: 'Update options' }).click();
      await expect(page.getByText('The configuration was saved successfully.')).toBeVisible();
    });

    await test.step('Vault stores all three values', async () => {
      expect(await tune()).toMatchObject({
        user_lockout_threshold: 7,
        user_lockout_duration: 600,
        user_lockout_counter_reset_duration: 1800,
      });
    });

    await test.step('opening the page again shows the saved values', async () => {
      await page.goto(options);
      await expect(page.locator(label('Lockout duration').value)).toHaveValue('10');
      await expect(page.locator(label('Lockout counter reset').value)).toHaveValue('30');
      await expect(page.locator(label('Lockout counter reset').unit)).toHaveValue('m');
      await expect(page.locator('[data-test-input="user_lockout_config.lockout_threshold"]')).toHaveValue(
        '7'
      );
      const counterReset = page.getByText('Lockout counter reset', { exact: true });
      await highlightText(page, counterReset.locator('..'), 'Lockout counter reset');
      await holdForVideo(page);
    });

    await test.step('saving again without changes keeps the values', async () => {
      await page.getByRole('button', { name: 'Update options' }).click();
      await expect(page.getByText('The configuration was saved successfully.')).toBeVisible();
      expect(await tune()).toMatchObject({
        user_lockout_threshold: 7,
        user_lockout_duration: 600,
        user_lockout_counter_reset_duration: 1800,
      });
    });
  } finally {
    await api.dispose();
  }
});
