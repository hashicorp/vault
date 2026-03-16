/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import { BasePage } from '../../pages/base';

test('kubernetes secrets workflow', async ({ page }) => {
  const basePage = new BasePage(page);

  await test.step('enable kubernetes secrets engine', async () => {
    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Secrets', exact: true }).click();
    await page.getByRole('link', { name: 'Enable new engine' }).click();
    await page.getByLabel('Kubernetes - enabled engine').click();
    await page.getByRole('button', { name: 'Enable engine' }).click();
    await expect(page.locator('section')).toContainText('Kubernetes not configured');
    await basePage.dismissFlashMessages();
  });

  await test.step('configure kubernetes secrets engine', async () => {
    await page.getByRole('link', { name: 'Configure Kubernetes' }).click();
    await page.getByRole('button', { name: 'Get config values' }).click();
    await expect(page.locator('section')).toContainText(
      'Vault could not infer a configuration from your environment variables. Check your configuration file to edit or delete them, or configure manually.'
    );
    await page.getByRole('radio', { name: 'Manual configuration Generate' }).check();
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page.locator('#error-kubernetes_host')).toContainText('Kubernetes host is required');
    await page.getByRole('textbox', { name: 'Kubernetes host' }).fill('https://192.168.99.100:8443');
    await page.getByRole('textbox', { name: 'Service account JWT' }).fill('test-jwt');
    await page.getByRole('textbox', { name: 'Kubernetes CA Certificate' }).fill('-----CERTIFICATE-----');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(
      page.locator('.info-table-row').filter({
        hasText: 'Kubernetes host https://192.168.99.100:8443',
      })
    ).toBeVisible();
    await expect(
      page.locator('.info-table-row').filter({
        hasText: 'Certificate PEM Format -----CERTIFICATE-----',
      })
    ).toBeVisible();
    await page.getByRole('link', { name: 'Exit configuration' }).click();
    await expect(page.locator('section')).toContainText(
      'Roles Create Role The number of Vault roles being used to generate Kubernetes credentials. None'
    );
    await expect(page.locator('section')).toContainText(
      'Generate credentials Quickly generate credentials by typing the role name. Type to find a role... Generate'
    );
    await basePage.dismissFlashMessages();
  });

  await test.step('edit kubernetes configuration', async () => {
    await page.getByRole('button', { name: 'Manage' }).click();
    await page.getByRole('link', { name: 'Configure' }).click();
    await page.getByRole('link', { name: 'Edit configuration' }).click();
    await expect(page.getByRole('textbox', { name: 'Service account JWT' })).toBeEmpty();
    await page.getByRole('textbox', { name: 'Kubernetes host' }).fill('https://127.0.0.1:8443');
    await page.getByRole('textbox', { name: 'Kubernetes CA Certificate' }).fill('-----NEW CERT-----');
    await page.getByRole('button', { name: 'Save' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
    await expect(
      page.locator('.info-table-row').filter({
        hasText: 'Kubernetes host https://127.0.0.1:8443',
      })
    ).toBeVisible();
    await expect(
      page.locator('.info-table-row').filter({
        hasText: 'Certificate PEM Format -----NEW CERT-----',
      })
    ).toBeVisible();
    await page.getByRole('link', { name: 'Exit configuration' }).click();
    await basePage.dismissFlashMessages();
  });

  await test.step('create kubernetes role', async () => {
    await page.getByRole('link', { name: 'Roles' }).click();
    await expect(page.locator('section')).toContainText(
      'No roles yet When created, roles will be listed here. Create a role to start generating service account tokens.'
    );
    await page.getByRole('link', { name: 'Create role' }).click();
    await expect(page.locator('section')).toContainText(
      'Choose an option above To configure a Vault role, choose what should be generated in Kubernetes by Vault.'
    );
    await page.getByRole('radio', { name: 'Generate token only using' }).check();
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page.locator('#error-name')).toContainText('Name is required');
    await page.getByRole('textbox', { name: 'Role name' }).fill('test-role');
    await page.getByRole('textbox', { name: 'Service account name' }).fill('foo');
    await page.getByRole('textbox', { name: 'Allowed Kubernetes namespaces' }).fill('*');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(
      page.locator('.info-table-row').filter({
        hasText: 'Role name test-role',
      })
    ).toBeVisible();

    await page.getByRole('link', { name: 'Roles' }).click();
    await expect(page.getByRole('link', { name: 'test-role', exact: true })).toBeVisible();
    await basePage.dismissFlashMessages();
  });

  await test.step('edit kubernetes role', async () => {
    await page.getByRole('link', { name: 'test-role', exact: true }).click();
    await expect(
      page.locator('.info-table-row').filter({
        hasText: 'Service account name foo',
      })
    ).toBeVisible();
    await page.getByRole('button', { name: 'Manage' }).click();
    await page.getByRole('link', { name: 'Edit role' }).click();
    await page.getByRole('radio', { name: 'Generate token, service' }).check();
    await page.getByRole('textbox', { name: 'Kubernetes role name' }).click();
    await page.getByRole('textbox', { name: 'Kubernetes role name' }).fill('admin');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(
      page.locator('.info-table-row').filter({
        hasText: 'Kubernetes role name admin',
      })
    ).toBeVisible();

    await page.getByRole('button', { name: 'Manage' }).click();
    await page.getByRole('link', { name: 'Edit role' }).click();
    await page.getByRole('radio', { name: 'Generate entire Kubernetes' }).check();
    await page.getByLabel('Role rules template').selectOption('5');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page.locator('section')).toContainText(
      'Generated role rules Role rules rules: - apiGroups: [""] resources: ["secrets", "services"] verbs: ["get", "watch", "list", "create", "delete", "deletecollection", "patch", "update"]'
    );
    await basePage.dismissFlashMessages();
  });

  await test.step('generate credentials from kubernetes role', async () => {
    // mock since we aren't connected to a kubernetes cluster
    await page.route('**/v1/kubernetes/creds/test-role', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          lease_duration: 3600,
          lease_id: 'kubernetes/creds/test-role/aWczfcfJ7NKUdiirJrPXIs38',
          data: {
            service_account_name: 'default',
            service_account_namespace: 'default',
            service_account_token: 'eyJhbGciOiJSUzI1NiIsImtpZCI6Imlr',
          },
        }),
      });
    });
    await page.getByRole('button', { name: 'Manage' }).click();
    await page.getByRole('link', { name: 'Generate credentials' }).click();
    await page.getByRole('textbox', { name: 'Kubernetes namespace' }).fill('user');
    await page.getByRole('button', { name: 'Generate credentials' }).click();
    await expect(page.getByLabel('Warning')).toContainText(
      "Warning You won't be able to access these credentials later, so please copy them now."
    );
    await page.getByRole('button', { name: 'Done' }).click();
    await page.getByRole('link', { name: 'Roles' }).click();
    await basePage.dismissFlashMessages();
  });

  await test.step('filter kubernetes roles', async () => {
    await page.getByRole('link', { name: 'Create role' }).click();
    await page.getByRole('radio', { name: 'Generate token only using' }).check();
    await page.getByRole('textbox', { name: 'Role name' }).fill('foo');
    await page.getByRole('textbox', { name: 'Allowed Kubernetes namespaces' }).fill('*');
    await page.getByRole('textbox', { name: 'Service account name' }).fill('bar');
    await page.getByRole('button', { name: 'Save' }).click();
    await page.getByRole('link', { name: 'Roles' }).click();
    await page.getByRole('textbox', { name: 'Search by path' }).fill('test');
    await page.getByRole('button', { name: 'Search' }).click();
    await expect(page.getByRole('link', { name: 'test-role', exact: true })).toBeVisible();
    await expect(page.getByRole('link', { name: 'foo', exact: true })).not.toBeVisible();
    await page.getByRole('textbox', { name: 'Search by path' }).fill('foo');
    await page.getByRole('button', { name: 'Search' }).click();
    await expect(page.getByRole('link', { name: 'foo', exact: true })).toBeVisible();
    await expect(page.getByRole('link', { name: 'test-role', exact: true })).not.toBeVisible();
  });

  await test.step('delete kubernetes role', async () => {
    await page.getByRole('button', { name: 'More options' }).click();
    await page.getByRole('button', { name: 'Delete' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
    await expect(page.locator('section')).toContainText('There are no roles matching "foo"');
    await page.getByRole('textbox', { name: 'Search by path' }).click();
    await page.getByRole('textbox', { name: 'Search by path' }).fill('');
    await page.getByRole('button', { name: 'Search' }).click();
    await expect(page.getByRole('link', { name: 'test-role', exact: true })).toBeVisible();
    await basePage.dismissFlashMessages();
  });

  await test.step('kubernetes overview', async () => {
    await page.getByRole('link', { name: 'Overview' }).click();
    await expect(page.locator('section')).toContainText(
      'Roles View Roles The number of Vault roles being used to generate Kubernetes credentials. 1'
    );
    await page.getByRole('link', { name: 'View Roles' }).click();
    await expect(page.getByRole('link', { name: 'test-role', exact: true })).toBeVisible();
    await page.getByRole('link', { name: 'Overview' }).click();
    await page.getByText('Type to find a role...').click();
    await page.getByRole('option', { name: 'test-role' }).click();
    await page.getByRole('button', { name: 'Generate' }).click();
    await expect(page.getByRole('paragraph')).toContainText(
      'This will generate credentials using the role test-role.'
    );
  });
});
