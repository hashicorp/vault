/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, type Page } from '@playwright/test';
import { ALL_ENGINES } from '../../lib/core/addon/utils/all-engines-metadata';

export const findEngineDisplayName = (engineType: string) => {
  const engine = ALL_ENGINES.find((e) => e.type === engineType);
  return engine ? engine.displayName : engineType;
};

const configurableEngines = ALL_ENGINES.filter((engine) => engine.isConfigurable).map(
  (engine) => engine.type
);

export class ConfigurationSettingsPage {
  constructor(protected page: Page) {}

  async navigateToConfiguration(path: string) {
    await this.page.getByRole('button', { name: 'Manage', exact: true }).click();
    await this.page.getByRole('link', { name: 'Configure' }).click();
    await expect(this.page.getByRole('heading', { name: `${path} configuration` })).toContainText(
      `${path} configuration`
    );
  }

  async assertPluginSettingsTabActive(engineType: string) {
    const engineDisplayName = findEngineDisplayName(engineType);
    await expect(this.page.getByRole('link', { name: 'General settings' })).not.toHaveClass(/active/);
    await expect(this.page.getByRole('link', { name: `${engineDisplayName} settings` })).toHaveClass(
      /active/
    );
    await expect(this.page.getByRole('link', { name: `${engineDisplayName} settings` })).toContainText(
      `${engineDisplayName} settings`
    );
  }

  async navigateToGeneralSettings(engineType: string) {
    const engineDisplayName = findEngineDisplayName(engineType);
    await this.page.getByRole('link', { name: 'General settings' }).click();
    await expect(this.page.getByRole('link', { name: 'General settings' })).toHaveClass(/active/);

    if (configurableEngines.includes(engineDisplayName)) {
      await expect(
        this.page.getByRole('link', {
          name: `${engineDisplayName} settings`,
        })
      ).not.toHaveClass(/active/);
    }

    await expect(this.page.getByRole('form', { name: 'general settings form' })).toBeVisible();
  }

  async editAndVerifyGeneralSettings(path: string, engineType: string, isExternalPlugin = false) {
    await expect(this.page.getByText(`Engine type ${engineType}`)).toBeVisible();
    await expect(this.page.getByText('Running version')).toContainText('Running version');
    await expect(this.page.getByRole('group', { name: 'Path' })).toBeVisible();
    await expect(this.page.getByRole('button', { name: `copy ${path}/` })).toContainText(`${path}/`);
    await expect(this.page.getByRole('group', { name: 'Accessor' })).toBeVisible();
    await expect(
      this.page.getByText('Default time-to-live (TTL) How long secrets in this engine stay valid. seconds')
    ).toBeVisible();

    if (!isExternalPlugin) {
      await this.page.getByRole('textbox', { name: 'Description' }).fill('some description');
      await this.page.getByRole('textbox', { name: 'Default time-to-live (TTL)time' }).fill('2');
      await this.page.getByRole('button', { name: 'Save changes' }).click();
      await expect(this.page.getByRole('alert', { name: 'Configuration saved' })).toBeVisible();
      await expect(this.page.getByRole('textbox', { name: 'Description' })).toHaveValue('some description');
      await expect(this.page.getByRole('textbox', { name: 'Default time-to-live (TTL)time' })).toHaveValue(
        '2'
      );
    }
  }

  async verifyUnsavedChangesModalOnNavigateAway(path: string) {
    // make a change to trigger the unsaved changes modal
    await this.page.getByRole('textbox', { name: 'Description' }).fill('unsaved changes test');

    // try to navigate away
    this.page.getByRole('link', { name: 'Back to main navigation' }).click();

    // verify the unsaved changes modal appears
    const modal = this.page.getByRole('dialog', { name: 'Unsaved changes' });
    await expect(modal).toBeVisible();
    await expect(
      this.page.getByText("You've made changes to the following: Description Would you like to apply them?")
    ).toBeVisible();

    // save unsaved changes
    await modal.getByRole('button', { name: 'Save changes' }).click();

    // verify still on the configuration page and description was updated
    await this.page.getByRole('link', { name: 'Secrets', exact: true }).click();
    await this.page.getByRole('link', { name: path }).click();
    await this.page.getByRole('button', { name: 'Manage' }).click();
    await this.page.getByRole('link', { name: 'Configure' }).click();
    await this.page.getByRole('link', { name: 'General settings' }).click();
    await expect(this.page.getByRole('textbox', { name: 'Description' })).toHaveValue('unsaved changes test');

    // make another change to trigger the unsaved changes modal
    await this.page.getByRole('textbox', { name: 'Description' }).fill('unsaved changes test 2');

    // try to navigate away again
    this.page.getByRole('link', { name: 'Back to main navigation' }).click();

    // dismiss the modal
    await modal.getByRole('button', { name: 'Discard changes' }).click();
    await expect(this.page.getByRole('link', { name: 'Dashboard' })).toBeVisible();
  }
}
