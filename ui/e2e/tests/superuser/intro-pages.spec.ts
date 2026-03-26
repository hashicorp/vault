/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import { DISMISSED_WIZARD_KEY } from '../../../app/utils/constants/wizard';
import { BasePage } from '../../pages/base';

test('intro pages workflow', async ({ page }) => {
  const basePage = new BasePage(page);
  await page.goto('dashboard');

  // remove the dismissed wizards from local storage so that the intro pages will show up on login
  await page.evaluate(
    ({ key }) => {
      localStorage.removeItem(key);
    },
    { key: DISMISSED_WIZARD_KEY }
  );

  await test.step('Verify secrets intro page content and workflow', async () => {
    await page.getByRole('link', { name: 'Secrets', exact: true }).click();
    await expect(page.getByRole('heading', { name: 'Welcome to Secrets engines' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Enable a Secret engine' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Skip' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'View documentation' })).toBeVisible();

    // verify clicking enable routes to enablement page
    await page.getByRole('link', { name: 'Enable a Secret engine' }).click();
    await expect(page.getByText('Enable secrets engine')).toBeVisible();
    await expect(page.getByLabel('KV - enabled engine type')).toBeVisible();

    // navigate back to secrets intro page and click skip and verify the intro page is dismissed
    await page.getByLabel('Secrets Navigation Links').getByRole('link', { name: 'Secrets engines' }).click();
    await page.getByRole('button', { name: 'Skip' }).click();

    // intro page is dismissed and not visible
    await expect(page.getByRole('heading', { name: 'Welcome to Secrets engines' })).not.toBeVisible();

    // user sees secrets engines list and can see the 'new to secrets engines' button
    await expect(page.getByRole('heading', { name: 'Secrets engines' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'cubbyhole/' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'New to Secret engines?' })).toBeVisible();

    // clicking enable from banner routes to enablement page
    await page.getByRole('button', { name: 'New to Secret engines?' }).click();
    await page.getByRole('link', { name: 'Enable a Secret engine' }).click();
    await expect(page.getByText('Enable secrets engine')).toBeVisible();

    // click button and close the banner and assert the banner is closed
    await page.getByLabel('breadcrumbs').getByRole('link', { name: 'Secrets engines' }).click();
    await page.getByRole('button', { name: 'Skip' }).click();
    await page.getByRole('button', { name: 'New to Secret engines?' }).click();
    await page.getByRole('button', { name: 'Close' }).click();
    await expect(page.getByText('Welcome to Secrets engines')).not.toBeVisible();

    // enable an engine and verify the welcome button doesn't render on secrets engines page
    await basePage.enableEngine('transit', 'test-transit');
    await page.getByLabel('Secrets Navigation Links').getByRole('link', { name: 'Secrets engines' }).click();
    await expect(page.getByRole('button', { name: 'New to Secret engines?' })).not.toBeVisible();

    await basePage.disableEngine('test-transit');
    await page.getByRole('link', { name: 'Back to main navigation' }).click();
  });

  await test.step('Verify auth methods intro page content and workflow', async () => {
    await page.getByRole('link', { name: 'Access control', exact: true }).click();
    await page.getByRole('link', { name: 'Authentication methods' }).click();
    await expect(page.getByRole('heading', { name: 'Welcome to Authentication methods' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Enable a new method' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Skip' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'View documentation' })).toBeVisible();

    // verify clicking enable routes to enablement page
    await page.getByRole('link', { name: 'Enable a new method' }).click();
    await expect(page.getByText('Enable an Authentication Method')).toBeVisible();
    await expect(page.getByLabel('Userpass - enabled engine type')).toBeVisible();

    // navigate back to auth methods intro page and click skip and verify the intro page is dismissed
    await page.getByRole('link', { name: 'Authentication methods' }).click();
    await page.getByRole('button', { name: 'Skip' }).click();
    await expect(page.getByRole('heading', { name: 'Welcome to Authentication methods' })).not.toBeVisible();

    // user sees auth methods list and can see the 'new to auth methods' button
    await expect(page.getByRole('heading', { name: 'Authentication methods' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'token/' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'New to Auth methods?' })).toBeVisible();

    // clicking enable from banner routes to enablement page
    await page.getByRole('button', { name: 'New to Auth methods?' }).click();
    await page.getByRole('link', { name: 'Enable a new method' }).click();
    await expect(page.getByText('Enable an Authentication Method')).toBeVisible();

    // click button and close the banner and assert the banner is closed
    await page.getByRole('link', { name: 'Authentication methods' }).click();
    await page.getByRole('button', { name: 'Skip' }).click();
    await page.getByRole('button', { name: 'New to Auth methods?' }).click();
    await page.getByRole('button', { name: 'Close' }).click();
    await expect(page.getByText('Welcome to Authentication methods')).not.toBeVisible();

    // enable a method and verify the welcome button doesn't render on auth methods page
    await page.getByRole('link', { name: 'Enable new method' }).click();
    await page.getByLabel('Userpass - enabled engine type').click();
    await page.getByRole('button', { name: 'Enable method' }).click();
    await page.getByRole('link', { name: 'Authentication methods' }).click();
    await expect(page.getByRole('button', { name: 'New to Auth methods?' })).not.toBeVisible();

    await basePage.dismissFlashMessages();

    await page.getByRole('link', { name: 'Back to main navigation' }).click();
  });

  await test.step('verify namespace intro page content and workflow', async () => {
    await page.getByRole('link', { name: 'Access control', exact: true }).click();
    await page.getByRole('link', { name: 'Namespaces' }).click();

    await expect(page.getByRole('heading', { name: 'Welcome to Namespaces' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Guided start' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Skip' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'View documentation' })).toBeVisible();

    // verify clicking create routes to guided start page
    await page.getByRole('button', { name: 'Guided start' }).click();
    await expect(page.getByText('Namespaces Guided Start')).toBeVisible();
    await expect(page.getByRole('tab', { name: 'Select setup (current)' })).toBeVisible();
    await page.getByRole('button', { name: 'Exit' }).click();

    await expect(page.getByRole('heading', { name: 'Welcome to Namespaces' })).not.toBeVisible();
    await expect(page.locator('h1')).toHaveText('Namespaces');
    await expect(page.getByRole('link', { name: 'Create namespace' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'New to Namespaces?' })).toBeVisible();
    await page.getByRole('link', { name: 'Create namespace' }).click();
    await page.getByRole('textbox', { name: 'Path' }).fill('testNs');
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByRole('button', { name: 'New to Namespaces?' })).not.toBeVisible();
  });

  await test.step('create policy for test namespace for testing ACL intro page', async () => {
    await page.getByRole('link', { name: 'ACL policies' }).click();
    await page.getByRole('link', { name: 'Create ACL policy' }).click();
    await page.getByRole('textbox', { name: 'Policy name' }).fill('testns');
    await page.getByRole('textbox', { name: 'Resource path' }).fill('testNs/*');
    await page.getByRole('checkbox', { name: 'create' }).check();
    await page.getByRole('checkbox', { name: 'read' }).check();
    await page.getByRole('checkbox', { name: 'list' }).check();
    await page.getByRole('button', { name: 'Create policy' }).click();

    await page.getByRole('button', { name: 'root' }).click();
    await page.getByRole('option', { name: 'testNs' }).click();
  });

  await test.step('Verify access control intro page content and workflow', async () => {
    await page.goto('dashboard?namespace=testNs');

    //nav to access control intro page and verify the content
    await page.getByRole('link', { name: 'Access control', exact: true }).click();
    await expect(page.getByRole('heading', { name: 'Welcome to ACL Policies' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Create a policy' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Skip' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'View documentation' })).toBeVisible();

    await page.getByRole('button', { name: 'Skip' }).click();
    await expect(page.getByRole('heading', { name: 'Welcome to ACL Policies' })).not.toBeVisible();
    await expect(page.getByRole('heading', { name: 'ACL policies' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'New to ACL policies?' })).toBeVisible();

    await page.getByRole('button', { name: 'New to ACL policies?' }).click();
    await expect(page.getByText('Welcome to ACL Policies')).toBeVisible();
    await page.getByRole('link', { name: 'Create a policy' }).click();
    await page.getByRole('textbox', { name: 'Policy name' }).fill('testpolicy');
    await page.getByRole('textbox', { name: 'Resource path' }).fill('testpath');
    await page.getByText('create', { exact: true }).click();
    await page.getByRole('button', { name: 'Create policy' }).click();
    await expect(page.getByRole('heading', { name: 'testpolicy' })).toBeVisible();

    await page.getByLabel('breadcrumbs').getByRole('link', { name: 'ACL policies' }).click();
    await expect(page.getByRole('link', { name: 'testpolicy', exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: 'New to ACL policies?' })).not.toBeVisible();
  });

  await test.step('cleanup', async () => {
    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Namespaces' }).click();
    await page.getByRole('button', { name: 'More options' }).click();
    await page.getByRole('button', { name: 'Delete' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
    await page.getByRole('link', { name: 'ACL policies' }).click();
    await page.getByRole('link', { name: 'testns Policy nav menu' }).getByLabel('Policy nav menu').click();
    await page.getByRole('button', { name: 'Delete' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
    await page.getByRole('link', { name: 'Authentication methods' }).click();
    await page
      .getByRole('link', { name: 'Type of auth mount userpass/' })
      .getByLabel('Overflow options')
      .click();
    await page.getByRole('button', { name: 'Disable' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
  });
});
