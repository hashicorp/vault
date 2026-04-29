/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import { BasePage } from '../../pages/base';
import { generateActivityResponse } from '../../../mirage/handlers/clients';
import { ACTIVITY_EXPORT_STUB } from '../../../tests/helpers/clients/client-count-helpers';

import type { ActivityData } from '../../../types/vault/client-counts/activity-api';

test('client counts workflow', async ({ page }) => {
  const basePage = new BasePage(page);
  const activityResponse = generateActivityResponse();
  const activityData: ActivityData = activityResponse.data;

  await test.step('navigate to client counts page and mock data', async () => {
    await page.route('**/sys/internal/counters/activity', async (route) =>
      route.fulfill({ json: activityResponse })
    );
    await page.route('**/sys/internal/counters/activity/export**', async (route) =>
      route.fulfill({ body: ACTIVITY_EXPORT_STUB })
    );

    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Client count' }).click();
    await page.waitForResponse('**/sys/internal/counters/activity');
    await page.waitForResponse('**/sys/internal/counters/activity/export**');
    await basePage.dismissFlashMessages();
  });

  await test.step('client usage header', async () => {
    await expect(page.getByText('For billing period: July 2023 - January 2024')).toBeVisible();
    await expect(page.getByText('Dashboard last updated:')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Refresh page' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Export activity data' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'July 2023', exact: true })).toBeVisible();
  });

  await test.step('client usage overview', async () => {
    const { entity_clients, non_entity_clients, acme_clients } = activityData.total;
    const totalClients = (entity_clients + non_entity_clients + acme_clients).toLocaleString();
    await expect(page.getByRole('img', { name: `Total of ${totalClients} Total` })).toBeVisible();
    await expect(page.getByRole('switch', { name: 'Split by client type' })).toBeVisible();
    await expect(page.locator('.lineal-chart')).toBeVisible();
    await expect(page.locator('section')).toContainText(
      `${entity_clients.toLocaleString()} Entity clients ${non_entity_clients.toLocaleString()} Non-entity clients ${acme_clients.toLocaleString()} ACME clients`
    );
  });

  await test.step('client attribution filters', async () => {
    await page.getByRole('button', { name: 'Namespace', exact: true }).click();
    await page.getByRole('option', { name: 'ns1/' }).click();
    await expect(page.getByRole('button', { name: 'Dismiss ns1/' })).toBeVisible();
    await page.getByRole('button', { name: 'Mount path', exact: true }).click();
    await page.getByRole('option', { name: 'auth/token/0/' }).click();
    await expect(page.getByRole('button', { name: 'Dismiss auth/token/0/' })).toBeVisible();
    await page.getByRole('button', { name: 'Mount type', exact: true }).click();
    await page.getByRole('option', { name: 'token' }).click();
    await expect(page.getByRole('button', { name: 'Dismiss token' })).toBeVisible();
    await page.getByRole('button', { name: 'Month' }).click();
    await page.getByRole('option', { name: 'January' }).click();
    await expect(page.getByRole('button', { name: 'Dismiss January' })).toBeVisible();
    await page.getByRole('button', { name: 'Clear filters' }).click();
    await expect(page.locator('section')).toContainText('Filters applied: None');
  });

  await test.step('client attribution table', async () => {
    const ns1 = activityData.by_namespace.find((ns) => ns.namespace_path === 'ns1/');
    const mount = ns1?.mounts[0];
    const mountType =
      mount?.mount_type === 'deleted mount' ? 'Deleted' : mount?.mount_type?.replace(/\/$/, '');
    const count = `${mount?.counts.clients}`;

    await page.getByRole('button', { name: 'Namespace', exact: true }).click();
    await page.getByRole('option', { name: 'ns1/' }).click();
    await expect(page.getByRole('gridcell', { name: 'ns1/' }).first()).toBeVisible();
    await expect(page.getByRole('gridcell', { name: mount?.mount_path }).first()).toBeVisible();
    await expect(page.getByRole('gridcell', { name: mountType }).first()).toBeVisible();
    await expect(page.getByRole('gridcell', { name: count }).first()).toBeVisible();
  });

  await test.step('client list', async () => {
    await page.getByRole('link', { name: 'Client list' }).click();
    await page.getByRole('button', { name: 'Clear filters' }).click();
    await expect(page.getByText('Namespace Mount path Mount type Month Filters applied: None')).toBeVisible();
    await expect(page.getByRole('tab', { name: 'Entity 8' })).toBeVisible();
    await expect(page.getByRole('tab', { name: 'Non-entity 18' })).toBeVisible();
    await expect(page.getByRole('tab', { name: 'ACME 16' })).toBeVisible();
    await expect(page.getByRole('tab', { name: 'Secret sync 8' })).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'copy 5692c6ef-c871-128e-fb06-' })).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'entity' }).first()).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'ns1/' }).first()).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'vK5Bt' })).toBeVisible();
    await expect(page.getByRole('gridcell', { name: '2023-09-15T23:48:09Z' })).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'auth/userpass/0/' })).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'userpass' }).nth(1)).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'auth_userpass_f47ad0b4' }).first()).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'entity_b3e2a7ff' }).first()).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'bob' }).first()).toBeVisible();
    await expect(page.getByRole('gridcell', { name: 'false' }).first()).toBeVisible();
    await expect(page.getByRole('gridcell', { name: '[ "7537e6b7-3b06-65c2-1fb2-' }).first()).toBeVisible();
  });

  await test.step('counts configuration', async () => {
    await page.getByRole('link', { name: 'Configuration' }).click();
    await expect(page.getByText('On', { exact: true })).toBeVisible();
    await expect(page.getByText('48')).toBeVisible();
    await page.getByRole('link', { name: 'Edit configuration' }).click();
    await page.getByRole('textbox', { name: 'Retention period' }).fill('72');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page.getByText('Retention period must be less than or equal to 60.')).toBeVisible();
    await page.getByRole('textbox', { name: 'Retention period' }).fill('60');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page.getByText('60')).toBeVisible();
  });
});
