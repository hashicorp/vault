/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('namespace workflow', async ({ page }) => {
  await test.step('create namespace', async () => {
    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Namespaces' }).click();
    await page.getByRole('link', { name: 'Create namespace' }).click();
    await page.getByRole('textbox', { name: 'Path' }).fill('testNamespace');
    await page.getByRole('button', { name: 'Save' }).click();
  });

  await test.step('should display the new namespace in the namespace picker and switch to it', async () => {
    await page.getByRole('button', { name: 'root' }).click();
    await page.getByRole('option', { name: 'testNamespace' }).click();
  });

  await test.step('should switch to the new namespace and display the correct header', async () => {
    await expect(page.locator('#app-main-content').getByText('testNamespace')).toBeVisible();
  });

  await test.step('delete namespace', async () => {
    await page.getByRole('button', { name: 'testNamespace' }).click();
    await page.getByRole('option', { name: 'root' }).click();
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Namespaces' }).click();
    await page.getByRole('button', { name: 'More options' }).click();
    await page.getByRole('button', { name: 'Delete' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
  });
});

test('namespace wizard workflow', async ({ page }) => {
  await page.goto('dashboard');

  await test.step('Navigate to namespaces wizard', async () => {
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Namespaces' }).click();
    await page.getByRole('button', { name: 'New to Namespaces?' }).click();
    const modal = page.getByRole('dialog', { name: 'Welcome to Namespaces' });
    await expect(modal).toBeVisible();
    await page.getByRole('button', { name: 'Guided start' }).click();
    await expect(page.getByRole('heading', { name: 'Namespaces Guided Start' })).toBeVisible();
  });

  await test.step('Should show step 1 selection options', async () => {
    await page.getByRole('heading', { name: 'What best describes your' }).click();

    await expect(
      page.getByRole('heading', {
        name: 'What best describes your access policy between teams and applications?',
      })
    ).toContainText('What best describes your access policy between teams and applications?');
    await expect(page.getByText('Flexible/shared access: our')).toBeVisible();
    await expect(
      page.getByText('Strict isolation required: our policy mandates hard boundaries (separate')
    ).toBeVisible();
  });

  await test.step('Should show flexible/shared access information in Step 1 if it is selected', async () => {
    await page.getByText('Flexible/shared access:').click();
    await expect(page.getByRole('heading', { name: 'Your recommended setup' })).toContainText(
      'Your recommended setup'
    );
    await page.getByRole('heading', { name: 'Single namespace' }).click();
    await expect(page.getByRole('heading', { name: 'Single namespace' })).toContainText('Single namespace');
    await expect(page.getByText('Your organization should be')).toContainText(
      'Your organization should be comfortable with your current setup of one global namespace. You can always add more namespaces later.'
    );
  });

  await test.step('Should navigate to "Apply changes" step once "next" is clicked for flexible/shared access selection', async () => {
    await page.getByRole('button', { name: 'Next' }).click();

    await expect(page.getByRole('heading', { name: "No action needed, you're all set." })).toBeVisible();
    await expect(
      page.getByRole('heading', { name: 'Next up: build out your access lists and identities' })
    ).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Why use ACL and identities?' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Set up identities' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Learn more about namespaces' })).toBeVisible();
  });

  await page.getByRole('button', { name: 'Back' }).click();

  await test.step('Should navigate to "Map out namespaces" once "next" is clicked for strict isolation selection', async () => {
    await page.getByText('Strict isolation required: our policy mandates hard boundaries (separate').click();
    await page.getByRole('button', { name: 'Next' }).click();
    await expect(page.getByRole('heading', { name: 'Map out your namespaces' })).toContainText(
      'Map out your namespaces'
    );
    await page.getByText('Create the namespaces you').click();
    await expect(page.getByText('Create the namespaces you')).toContainText(
      'Create the namespaces you need using the 3-layer structure, starting with the global level. Refresh the preview to update. These changes will only be applied on the next step, once you select the implementation method.'
    );
    await page.getByRole('textbox', { name: 'Global' }).fill('global');
    await page.getByRole('textbox', { name: 'Org' }).fill('org-1');
    await page.getByRole('textbox', { name: 'Project' }).fill('project-2');
    await page.getByRole('button', { name: 'Add' }).first().click();
    await page.getByRole('button', { name: 'Next' }).click();
  });

  await test.step('Should navigate and display "Apply changes" once "Map out your namespaces" is complete', async () => {
    await expect(
      page.getByText('Terraform automation Recommended Manage configurations by Infrastructure as')
    ).toBeVisible();
    await expect(
      page.getByText('variable "global_child_namespaces" { type = set(string) default = ["org-1"] }')
    ).toContainText(
      'variable "global_child_namespaces" { type = set(string) default = ["org-1"] } variable "global_org-1_child_namespaces" { type = set(string) default = ["project-2"] } resource "vault_namespace" "global" { path = "global" } resource "vault_namespace" "global_children" { for_each = var.global_child_namespaces namespace = vault_namespace.global.path path = each.key } resource "vault_namespace" "global_org-1_children" { for_each = var.global_org-1_child_namespaces namespace = vault_namespace.global_children["org-1"].path_fq path = each.key }'
    );
    await page.getByText('API/CLI Manage namespaces').click();
    await expect(page.getByText('API/CLI Manage namespaces')).toContainText(
      'API/CLI Manage namespaces directly via the Vault CLI or REST API. Best for quick updates, custom scripting, or terminal-based workflows.'
    );
    await expect(page.getByText('curl \\ --header "X-Vault-')).toContainText(
      'curl \\ --header "X-Vault-Token: $VAULT_TOKEN" \\ --request PUT \\ $VAULT_ADDR/v1/sys/namespaces/global curl \\ --header "X-Vault-Token: $VAULT_TOKEN" \\ --header "X-Vault-Namespace: /global" \\ --request PUT \\ $VAULT_ADDR/v1/sys/namespaces/org-1 curl \\ --header "X-Vault-Token: $VAULT_TOKEN" \\ --header "X-Vault-Namespace: /global/org-1" \\ --request PUT \\ $VAULT_ADDR/v1/sys/namespaces/project-2'
    );
    await page.getByRole('radio', { name: 'Vault UI workflow Apply' }).check();
    await expect(page.getByText('Apply changes immediately.')).toContainText(
      'Apply changes immediately. Note: Changes made in the UI will be overwritten by any future updates made via Infrastructure as Code (Terraform).'
    );
    await page.getByRole('tab', { name: 'Map out namespaces (complete)' }).click();
    await expect(page.getByRole('heading', { name: 'Map out your namespaces' })).toContainText(
      'Map out your namespaces'
    );
  });
});
