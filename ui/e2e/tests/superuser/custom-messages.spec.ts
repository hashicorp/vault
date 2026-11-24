/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import { BasePage } from '../../pages/base';
import fs from 'fs';
import path from 'path';

const keysPath = path.resolve(__dirname, '../../tmp/superuser-keys.json');

enum MessageLocation {
  LoginPage = 'On the login page',
  PostLogin = 'After the user logs in',
}

enum MessageType {
  Alert = 'Alert message',
  Modal = 'Modal',
}

async function navigateToCustomMessages(page: BasePage['page']) {
  await page.goto('dashboard');
  await page.getByRole('link', { name: 'Operational tools', exact: true }).click();
  await page.getByRole('link', { name: 'Custom messages' }).click();
}

// A helper function to create messages which will default to creating
// active messages after user logins in, if start_time is not provided.
async function createMessage(
  page: BasePage['page'],
  title: string,
  message: string,
  location: MessageLocation = MessageLocation.PostLogin,
  type: MessageType = MessageType.Alert,
  startTime = '2020-01-01T00:00',
  endTime?: string,
  linkText?: string,
  linkUrl?: string
) {
  await page.getByRole('button', { name: 'Create message' }).click();
  await page.getByRole('textbox', { name: 'Title', exact: true }).fill(title);
  await page.getByRole('textbox', { name: 'Message', exact: true }).fill(message);
  await page.getByRole('radio', { name: location }).check();
  await page.getByRole('radio', { name: type }).check();
  await page.getByRole('textbox', { name: 'Message starts' }).fill(startTime);
  if (endTime) {
    await page.getByRole('radio', { name: 'Specific date' }).click();
    await page.locator('[data-test-input="end_time"]').fill(endTime);
  }
  if (linkText) {
    await page.getByPlaceholder('Display text (e.g. Learn more)').fill(linkText);
  }
  if (linkUrl) {
    await page.getByPlaceholder('Link URL (e.g. https://www.hashicorp.com/)').fill(linkUrl);
  }
  await page.getByRole('button', { name: 'Create message' }).click();
}

async function deleteMessageFromList(page: BasePage['page'], title: string) {
  await page.getByRole('link', { name: title }).getByLabel('Message popup menu').click();
  await page.getByRole('button', { name: 'Delete' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
}

test.describe('Custom messages workflows', () => {
  test.beforeEach(async ({ page }) => {
    await navigateToCustomMessages(page);
  });

  test('custom messages page is reachable and initially shows empty state', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'Custom messages' })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'No messages yet' })).toBeVisible();

    await test.step('default tab selected is After user logs in', async () => {
      await expect(page.getByRole('link', { name: 'After user logs in' })).toHaveClass(/active/);
    });

    await test.step('filtering messages is disabled on init', async () => {
      await expect(page.getByRole('button', { name: 'Apply filters' })).toBeDisabled();
    });
  });

  test('create and view a message on the login page', async ({ page }) => {
    await createMessage(
      page,
      'Login Page banner',
      'This is a login page banner message.',
      MessageLocation.LoginPage,
      MessageType.Alert
    );

    // Logout and route to login page to verify the login page banner appears
    await page.getByRole('button', { name: 'User menu' }).click();
    await page.getByRole('listitem').filter({ hasText: 'Log out' }).click();
    await expect(page.getByLabel('Login Page Banner')).toBeVisible();

    const { root_token } = JSON.parse(fs.readFileSync(keysPath, 'utf-8'));
    await page.getByRole('textbox', { name: 'Token' }).fill(root_token);
    await page.getByRole('button', { name: 'Sign in' }).click();

    // Navigate back to custom messages list to delete the message
    await expect(page.getByRole('button', { name: 'root' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Dashboard' })).toBeVisible();
    await page.getByRole('link', { name: 'Operational tools', exact: true }).click();
    await page.getByRole('link', { name: 'Custom messages' }).click();
    await expect(page.getByRole('heading', { name: 'Custom messages' })).toBeVisible();

    // Cleanup
    await page.getByRole('link', { name: 'On login page' }).click();
    await deleteMessageFromList(page, 'Login Page banner');
  });

  test('create a modal message on the login page with preview', async ({ page }) => {
    await test.step('can toggle to login page tab', async () => {
      const loginPageTab = page.getByRole('link', { name: 'On login page' });
      await loginPageTab.click();
      await expect(loginPageTab).toHaveClass(/active/);
    });

    await page.getByRole('button', { name: 'Create message' }).click();
    await page.getByRole('radio', { name: 'Modal' }).click();
    await page.getByRole('textbox', { name: 'Title' }).fill('E2E Login Modal');
    await page
      .getByRole('textbox', { name: 'Message', exact: true })
      .fill('Important notice for all users on the login page.');

    await page.getByRole('textbox', { name: 'Message starts' }).fill('2020-01-01T00:00');

    await test.step('can see a preview of a modal', async () => {
      await page.getByRole('button', { name: 'Preview' }).click();
      await expect(page.getByRole('dialog').getByText('E2E Login Modal')).toBeVisible();
      await expect(
        page.getByRole('dialog').getByText('Important notice for all users on the login page.')
      ).toBeVisible();
      await page.getByRole('button', { name: 'Confirm' }).click();
    });

    await page.getByRole('button', { name: 'Create message' }).click();

    await page.getByRole('link', { name: 'Custom messages' }).first().click();
    await page.getByRole('link', { name: 'On login page' }).click();

    // Cleanup
    await deleteMessageFromList(page, 'E2E Login Modal');
  });

  test('create and view multiple messages', async ({ page }) => {
    // Create a an active banner message to be shown after user logs in
    await createMessage(
      page,
      'E2E Post-login banner',
      'This is a post-login banner message.',
      MessageLocation.PostLogin,
      MessageType.Alert
    );
    await navigateToCustomMessages(page);

    // Create an active modal message to be shown after user logs in
    await createMessage(
      page,
      'E2E Post-login Modal',
      'This is a post-login modal message.',
      MessageLocation.PostLogin,
      MessageType.Modal
    );

    // Verify the banner messages and modal appear in the correct locations
    await expect(page.getByLabel('E2E Post-login Banner')).toBeVisible();
    await expect(page.getByLabel('E2E Post-login Modal')).toBeVisible();

    await test.step('banners can only be dismissed after modal(s) are dismissed', async () => {
      await expect(
        page
          .getByLabel('E2E Post-login Banner')
          .getByRole('button', { name: 'Dismiss' })
          .click({ timeout: 1000 })
      ).rejects.toThrow();

      // Close the modal and confirm the banner can be dismissed
      await page.getByRole('button', { name: 'Confirm' }).click();
      await expect(page.getByLabel('E2E Post-login Modal')).not.toBeVisible();
      await page.getByLabel('E2E Post-login banner').getByRole('button', { name: 'Dismiss' }).click();
      await expect(page.getByLabel('E2E Post-login Banner')).not.toBeVisible();
    });

    // Cleanup
    await page.getByRole('button', { name: 'Delete message' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
    await navigateToCustomMessages(page);
    await deleteMessageFromList(page, 'Post-login banner');
  });

  test('message details page shows correct field values', async ({ page }) => {
    await createMessage(
      page,
      'E2E Details Check',
      'Details page verification message.',
      MessageLocation.PostLogin,
      MessageType.Alert,
      '2020-01-01T00:00',
      '',
      'Learn more',
      'https://example.com'
    );

    await expect(page.getByRole('heading', { name: 'E2E Details Check' })).toBeVisible();
    await expect(page.locator('span[data-test-row-value="Message"]')).toHaveText(
      'Details page verification message.'
    );
    await expect(
      page.getByLabel('E2E Details Check').getByRole('link', { name: 'Learn more' })
    ).toBeVisible();
    await expect(
      page.getByLabel('E2E Details Check').getByRole('link', { name: 'Learn more' })
    ).toHaveAttribute('href', 'https://example.com');

    await navigateToCustomMessages(page);
    await deleteMessageFromList(page, 'E2E Details Check');
  });

  test('edit and delete messages', async ({ page }) => {
    const basePage = new BasePage(page);

    await createMessage(
      page,
      'E2E Test Banner',
      'This is a test banner message.',
      MessageLocation.PostLogin,
      MessageType.Alert
    );

    await expect(page.getByRole('heading', { name: 'E2E Test Banner' })).toBeVisible();
    await expect(page.locator('span[data-test-row-value="Message"]')).toHaveText(
      'This is a test banner message.'
    );

    await test.step('can edit and view message updates on details page', async () => {
      await page.getByRole('link', { name: 'Edit message' }).click();
      await expect(page.getByRole('heading', { name: 'Edit message' })).toBeVisible();
      await page.getByRole('textbox', { name: 'Title' }).fill('E2E Test Banner Updated');
      await basePage.dismissFlashMessages();
      await page.getByRole('button', { name: 'Save' }).click();

      await expect(page.getByRole('heading', { name: 'E2E Test Banner Updated' })).toBeVisible();
      await expect(page.getByLabel('E2E Test Banner Updated')).toBeVisible();
    });

    await test.step('can delete message from details page and navigate user back to list view', async () => {
      await page.getByRole('button', { name: 'Delete message' }).click();
      await page.getByRole('button', { name: 'Confirm' }).click();

      await expect(page.getByRole('link', { name: 'E2E Test Banner Updated' })).not.toBeVisible();
      await expect(page.getByLabel('E2E Test Banner Updated')).not.toBeVisible();
    });
  });

  test('list and filter messages', async ({ page }) => {
    // Active message with no expiration date
    await createMessage(
      page,
      'E2E Active Banner',
      'This message is currently active.',
      MessageLocation.PostLogin,
      MessageType.Alert,
      '2020-01-01T00:00'
    );
    await navigateToCustomMessages(page);

    // Active until message with a future expiration date
    await createMessage(
      page,
      'E2E Active until future date Banner',
      'This message will remain active for a while.',
      MessageLocation.PostLogin,
      MessageType.Alert,
      '2020-01-01T00:00',
      '2030-01-01T00:00'
    );
    await navigateToCustomMessages(page);

    // Scheduled message
    await createMessage(
      page,
      'E2E Scheduled Banner',
      'This message is scheduled to be active in the future.',
      MessageLocation.PostLogin,
      MessageType.Alert,
      '2030-01-01T00:00'
    );
    await navigateToCustomMessages(page);

    // Expired message
    await createMessage(
      page,
      'E2E Expired Banner',
      'This message has expired.',
      MessageLocation.PostLogin,
      MessageType.Alert,
      '2020-01-01T00:00',
      '2020-12-31T23:59'
    );
    await navigateToCustomMessages(page);

    await expect(page.locator('[data-test-list-item]')).toHaveCount(4);

    await test.step('displays badge statuses for messages', async () => {
      await expect(page.getByRole('link', { name: 'E2E Active Banner' })).toContainText('Active');
      await expect(page.getByRole('link', { name: 'E2E Active until future date Banner' })).toContainText(
        'Active until'
      );
      await expect(page.getByRole('link', { name: 'E2E Scheduled Banner' })).toContainText('Scheduled');
      await expect(page.getByRole('link', { name: 'E2E Expired Banner' })).toContainText('Inactive');
    });

    await test.step('can filter messages by status', async () => {
      await page.getByLabel('Filter by message status').selectOption('active');
      await page.getByRole('button', { name: 'Apply filters' }).click();

      await expect(page.locator('[data-test-list-item]')).toHaveCount(2);

      await page.getByLabel('Filter by message status').selectOption('inactive');
      await page.getByRole('button', { name: 'Apply filters' }).click();

      await expect(page.locator('[data-test-list-item]')).toHaveCount(2);
      await navigateToCustomMessages(page);
    });

    await test.step('can filter messages by type', async () => {
      await page.getByLabel('Filter by type').selectOption('modal');
      await page.getByRole('button', { name: 'Apply filters' }).click();

      await expect(page.locator('[data-test-list-item]')).toHaveCount(0);

      await page.getByLabel('Filter by type').selectOption('banner');
      await page.getByRole('button', { name: 'Apply filters' }).click();

      await expect(page.locator('[data-test-list-item]')).toHaveCount(4);
      await navigateToCustomMessages(page);
    });

    await test.step('can filter messages by text', async () => {
      await page.getByRole('searchbox', { name: 'Search by message title' }).fill('Scheduled');
      await page.getByRole('button', { name: 'Apply filters' }).click();
      await expect(page.locator('[data-test-list-item]')).toHaveCount(1);
      await page.getByRole('button', { name: 'Clear filters' }).click();
    });

    // Cleanup
    await deleteMessageFromList(page, 'E2E Expired Banner');
    await deleteMessageFromList(page, 'E2E Scheduled Banner');
    await deleteMessageFromList(page, 'E2E Active Banner');
    await deleteMessageFromList(page, 'E2E Active until future date Banner');
  });
});
