/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import { METRICS_DATA_RESPONSE } from '../../../tests/helpers/billing/stubs';
import { formatInTimeZone } from 'date-fns-tz';

test('billing metrics dashboard workflow', async ({ page }) => {
  await test.step('navigate to billing metrics page and mock data', async () => {
    await page.route('**/sys/license/features', async (route) =>
      route.fulfill({ json: { features: ['Consumption Billing'] } })
    );
    await page.route('**/sys/billing/overview', async (route) =>
      route.fulfill({ json: METRICS_DATA_RESPONSE })
    );

    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Billing metrics' }).click();
    await page.waitForResponse('**/sys/billing/overview');
  });

  await test.step('display billing metrics summary panel', async () => {
    await expect(page.getByRole('button', { name: 'From start of January' })).toContainText(
      'From start of January 2026'
    );
    await expect(page.getByRole('heading', { name: 'Billing metrics' })).toBeVisible();
    await expect(page.getByText('Data reflects usage across')).toContainText(
      'Data reflects usage across this Vault cluster. Billing metrics are used in license utilization.'
    );
    await expect(page.getByRole('heading', { name: 'Summary' })).toBeVisible();

    await expect(page.locator('section')).toContainText(
      'Summary Secrets 10 PKI units 100.1234 Data protection calls 420 Managed keys 430 KMIP Enabled Plugins 100'
    );
  });

  await test.step('display billing metrics details by metric panel', async () => {
    await expect(page.getByRole('heading', { name: 'Details by metric' })).toBeVisible();

    await expect(page.locator('section')).toContainText(
      'Secrets Highest number of static secrets, static roles, and dynamic roles managed on the cluster during the month. Secrets replicated to this cluster are not counted. Total 210 KV Secrets 10 Dynamic roles 130 Static roles 70'
    );
    await expect(page.locator('section')).toContainText(
      'Credential units Certificates, tokens, and other credentials issued during the month, adjusted by their duration. Total 200.3702 PKI units 100.1234 SSH OTP units 50.1234 SSH certificate units 50.1234'
    );
    await expect(page.locator('section')).toContainText(
      'Data protection calls Total number of data elements processed. Total 420 Transform 220 Transit 200'
    );
    await expect(page.locator('section')).toContainText(
      'Managed keys Highest number of cryptographic keys managed on the cluster during the month. Keys replicated to this cluster are not counted. Total 430 TOTP 220 KMSE 210'
    );
  });

  await test.step('change the billing period date', async () => {
    await page.getByRole('button', { name: 'From start of January' }).click();
    await page.getByRole('option', { name: 'From start of Dec' }).click();
    await expect(page.locator('section')).toContainText(
      'Summary Secrets 2 PKI units 100.1234 Data protection calls 220 Managed keys 220 KMIP Not enabled Plugins 100'
    );
    await expect(page.locator('section')).toContainText(
      'Secrets Highest number of static secrets, static roles, and dynamic roles managed on the cluster during the month. Secrets replicated to this cluster are not counted. Total 192 KV Secrets 2 Dynamic roles 125 Static roles 65'
    );
    await expect(page.getByRole('button', { name: 'From start of December' })).toContainText(
      'From start of December 2025'
    );
  });
});

test('billing metrics dashboard api returns two months of data', async ({ page }) => {
  await page.route('**/sys/license/features', async (route) =>
    route.fulfill({ json: { features: ['Consumption Billing'] } })
  );

  await page.goto('dashboard');

  // Set up response listener before clicking the link
  const responsePromise = page.waitForResponse('**/sys/billing/overview**');
  await page.getByRole('link', { name: 'Billing metrics' }).click();

  // Wait for the API response and get the actual months returned
  const response = await responsePromise;
  const data = await response.json();
  const months = data.data?.months || [];

  // Verify we have at least 2 months of data
  expect(months.length).toEqual(2);

  // Verify both months appear in the dropdown options
  const currentMonth = new Date(months[0].month);
  const previousMonth = new Date(months[1].month);

  // Click the date range dropdown to verify both months are available
  await page.getByRole('button', { name: formatInTimeZone(currentMonth, 'UTC', 'MMMM') }).click();
  await expect(page.getByRole('option', { name: formatInTimeZone(currentMonth, 'UTC', 'MMMM') })).toHaveText(
    formatInTimeZone(currentMonth, 'UTC', 'MMMM yyyy')
  );
  await expect(page.getByRole('option', { name: formatInTimeZone(previousMonth, 'UTC', 'MMMM') })).toHaveText(
    formatInTimeZone(previousMonth, 'UTC', 'MMMM yyyy')
  );
});
