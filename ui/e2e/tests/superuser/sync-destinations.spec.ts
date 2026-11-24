/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect, Page } from '@playwright/test';
import { BasePage } from '../../pages/base';
import {
  SYNC_DESTINATION_AWS_WIF_RESPONSE,
  SYNC_DESTINATION_AZURE_WIF_RESPONSE,
  SYNC_DESTINATION_GCP_WIF_RESPONSE,
} from '../../../tests/helpers/sync/mocks';

/**
 * Navigate to the sync destination creation form
 * @param page - The Playwright page object
 * @param type - The destination type ('aws-sm', 'azure-kv', 'gcp-sm')
 */
async function openCreateDestinationForm(page: Page, type: string) {
  await page.goto('dashboard');
  await page.getByRole('link', { name: 'Secrets', exact: true }).click();
  await page.getByRole('link', { name: 'Secrets sync' }).click();

  const enableButton = page.getByRole('button', { name: 'Enable' });
  // waitFor auto-waits for the page to render after navigation; catch means secrets sync is already activated
  const needsActivation = await enableButton
    .waitFor({ state: 'visible', timeout: 100 })
    .then(() => true)
    .catch(() => false);
  if (needsActivation) {
    await enableButton.click();
    await page.getByRole('checkbox', { name: "I've read the above linked" }).check();
    await page.getByRole('button', { name: 'Confirm' }).click();
  }

  await page.getByRole('link', { name: 'Create first destination' }).click();
  await page.locator(`[data-test-select-destination="${type}"]`).click();
}

test('sync destination wif workflow for aws', async ({ page }) => {
  const basePage = new BasePage(page);

  // Set up route mocks before any navigation
  await page.route('**/v1/sys/sync/destinations/aws-sm/test-aws', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(SYNC_DESTINATION_AWS_WIF_RESPONSE),
    });
  });

  await test.step('navigate to open new create destination form for AWS Secrets Manager', async () => {
    await openCreateDestinationForm(page, 'aws-sm');
  });
  await expect(page.getByRole('radio', { name: 'Workload Identity Federation' })).toBeVisible();
  await page.getByRole('radio', { name: 'Workload Identity Federation' }).check();
  await page.getByRole('button', { name: 'Create destination' }).click();

  await expect(page.getByText('Name is required.')).toBeVisible();
  await expect(page.getByText('Role ARN is required.')).toBeVisible();
  await expect(page.getByText('Identity token audience is required.')).toBeVisible();

  await page.getByRole('textbox', { name: 'Name' }).fill('test-aws');
  await page.getByRole('textbox', { name: 'Role ARN' }).fill('arn:aws:iam::111111111111:role/wif_test');
  await page
    .getByRole('textbox', { name: 'identity_token_audience' })
    .fill('vault-test.wif-test.sbx.hashidemos.io/v1/identity/oidc/secrets-sync');

  // Set up response watchers before clicking to avoid missing fast responses
  const postResponse = page.waitForResponse(
    (resp) =>
      resp.url().includes('/v1/sys/sync/destinations/aws-sm/test-aws') && resp.request().method() === 'POST'
  );
  const getResponse = page.waitForResponse(
    (resp) =>
      resp.url().includes('/v1/sys/sync/destinations/aws-sm/test-aws') && resp.request().method() === 'GET'
  );

  await page.getByRole('button', { name: 'Create destination' }).click();
  await Promise.all([postResponse, getResponse]);
  await expect(page.getByText('Connection successful', { exact: true })).toBeVisible();
  await expect(page.getByText('You have successfully created a sync destination')).toBeVisible();
  await basePage.dismissFlashMessages();

  // Verify that the details page shows the correct information from the mocked response
  await expect(page.getByRole('heading', { name: 'test-aws' })).toBeVisible();
  await expect(page.getByText('WIF', { exact: true })).toBeVisible();
});

test('sync destination wif workflow for azure', async ({ page }) => {
  const basePage = new BasePage(page);

  await page.route('**/v1/sys/sync/destinations/azure-kv/test-azure', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(SYNC_DESTINATION_AZURE_WIF_RESPONSE),
    });
  });

  await test.step('navigate to open new create destination form for Azure Key Vault', async () => {
    await openCreateDestinationForm(page, 'azure-kv');
  });
  await expect(page.getByRole('radio', { name: 'Workload Identity Federation' })).toBeVisible();
  await page.getByRole('radio', { name: 'Workload Identity Federation' }).check();
  await page.getByRole('button', { name: 'Create destination' }).click();

  await expect(page.getByText('Name is required.')).toBeVisible();
  await expect(page.getByText('Key Vault URI is required.')).toBeVisible();
  await expect(page.getByText('Tenant ID is required.')).toBeVisible();
  await expect(page.getByText('Client ID is required.')).toBeVisible();
  await expect(page.getByText('Identity token audience is required.')).toBeVisible();

  await page.getByRole('textbox', { name: 'Name' }).fill('test-azure');
  await page.getByRole('textbox', { name: 'Key Vault URI' }).fill('https://test-keyvault.vault.azure.net');
  await page.getByRole('textbox', { name: 'Tenant ID' }).fill('11111111-1111-1111-1111-111111111111');
  await page.getByRole('textbox', { name: 'Client ID' }).fill('test-client-id');
  await page
    .getByRole('textbox', { name: 'identity_token_audience' })
    .fill('https://test-audience.azure.com');

  const postResponse = page.waitForResponse(
    (resp) =>
      resp.url().includes('/v1/sys/sync/destinations/azure-kv/test-azure') &&
      resp.request().method() === 'POST'
  );
  const getResponse = page.waitForResponse(
    (resp) =>
      resp.url().includes('/v1/sys/sync/destinations/azure-kv/test-azure') &&
      resp.request().method() === 'GET'
  );

  await page.getByRole('button', { name: 'Create destination' }).click();
  await Promise.all([postResponse, getResponse]);
  await expect(page.getByText('Connection successful', { exact: true })).toBeVisible();
  await expect(page.getByText('You have successfully created a sync destination')).toBeVisible();
  await basePage.dismissFlashMessages();

  await expect(page.getByRole('heading', { name: 'test-azure' })).toBeVisible();
  await expect(page.getByText('WIF', { exact: true })).toBeVisible();
});

test('sync destination wif workflow for gcp', async ({ page }) => {
  const basePage = new BasePage(page);

  await page.route('**/v1/sys/sync/destinations/gcp-sm/test-gcp', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(SYNC_DESTINATION_GCP_WIF_RESPONSE),
    });
  });

  await test.step('navigate to open new create destination form for Google Secret Manager', async () => {
    await openCreateDestinationForm(page, 'gcp-sm');
  });
  await expect(page.getByRole('radio', { name: 'Workload Identity Federation' })).toBeVisible();
  await page.getByRole('radio', { name: 'Workload Identity Federation' }).check();
  await page.getByRole('button', { name: 'Create destination' }).click();

  await expect(page.getByText('Name is required.')).toBeVisible();
  await expect(page.getByText('Project ID is required.')).toBeVisible();
  await expect(page.getByText('Service account email is required.')).toBeVisible();
  await expect(page.getByText('Identity token audience is required.')).toBeVisible();

  await page.getByRole('textbox', { name: 'Name' }).fill('test-gcp');
  await page.getByRole('textbox', { name: 'Project ID' }).fill('test-gcp-project');
  await page
    .getByRole('textbox', { name: 'Service account email' })
    .fill('test-sa@test-gcp-project.iam.gserviceaccount.com');
  await page.getByRole('textbox', { name: 'identity_token_audience' }).fill('https://test-audience.gcp.com');

  const postResponse = page.waitForResponse(
    (resp) =>
      resp.url().includes('/v1/sys/sync/destinations/gcp-sm/test-gcp') && resp.request().method() === 'POST'
  );
  const getResponse = page.waitForResponse(
    (resp) =>
      resp.url().includes('/v1/sys/sync/destinations/gcp-sm/test-gcp') && resp.request().method() === 'GET'
  );

  await page.getByRole('button', { name: 'Create destination' }).click();
  await Promise.all([postResponse, getResponse]);
  await expect(page.getByText('Connection successful', { exact: true })).toBeVisible();
  await expect(page.getByText('You have successfully created a sync destination')).toBeVisible();
  await basePage.dismissFlashMessages();

  await expect(page.getByRole('heading', { name: 'test-gcp' })).toBeVisible();
  await expect(page.getByText('WIF', { exact: true })).toBeVisible();
});
