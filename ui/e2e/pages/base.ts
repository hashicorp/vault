/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Page } from '@playwright/test';

export class BasePage {
  constructor(protected page: Page) {}
  // remove all flash messages by clicking the dismiss button until there are no more
  // this is useful if many are rendered over top of a button preventing click in a test
  async dismissFlashMessages() {
    const locator = this.page.getByRole('button', { name: 'Dismiss' });
    // use a while loop because clicking one might cause the next one to shift or re-render.
    while ((await locator.count()) > 0) {
      await locator.first().click();
    }
  }

  async goToSecrets() {
    await this.page.getByRole('link', { name: 'Secrets', exact: true }).click();
    const skipButton = this.page.getByRole('button', { name: 'Skip' });
    if (await skipButton.isVisible()) {
      await skipButton.click();
    }
  }

  /**
   * Enable a secrets engine with a dynamic path
   * @param engineType - The type of engine to enable (e.g., 'KV', 'Transit', 'PKI Certificates')
   * @param path - The mount path for the engine
   * @param options - Optional configuration for the engine with variable key-value pairs
   */
  async enableEngine(
    engineType: string,
    path: string,
    options?: {
      defaultLeaseTtl?: { unit: number; option: string };
      maxLeaseTtl?: { unit: number; option: string };
      external?: boolean;
      pluginVersion?: string;
      skipEnable?: boolean;
    }
  ) {
    await this.page.goto('dashboard');
    await this.goToSecrets();

    // Click "Enable new engine"
    await this.page.getByRole('link', { name: 'Enable new engine' }).click();
    await this.page.getByRole('heading', { name: engineType }).click();

    if (options?.external) {
      // Prerequisite: mock plugin catalog endpoint in the test so the External plugin option is available.
      await this.page.locator('label:nth-child(2) > .hds-form-radio-card__control-wrapper').click();
      await this.page.getByLabel('Plugin version Required').selectOption(options.pluginVersion);
    }

    if (options?.defaultLeaseTtl) {
      await this.page.locator('label').filter({ hasText: 'Default Lease TTL Vault will' }).click();
      await this.page
        .getByLabel('TTL unit for Default Lease TTL')
        .selectOption(options.defaultLeaseTtl.option as string);
      await this.page
        .getByRole('group', { name: 'Default Lease TTL Lease will' })
        .getByLabel('Number of units')
        .fill(options.defaultLeaseTtl.unit.toString());
    }

    if (options?.maxLeaseTtl) {
      await this.page
        .getByLabel('TTL unit for Max Lease TTL')
        .selectOption(options.maxLeaseTtl.option as string);
      await this.page
        .getByRole('group', { name: 'Max Lease TTL Lease will' })
        .getByLabel('Number of units')
        .fill(options.maxLeaseTtl.unit.toString());
    }

    // Fill in the path
    await this.page.getByRole('textbox', { name: 'Path' }).fill(path);

    // Enable the engine
    if (!options?.skipEnable) {
      await this.page.getByRole('button', { name: 'Enable engine' }).click();
    }
  }
}
