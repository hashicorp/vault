/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { randomUUID } from 'crypto';
import { expect, type APIRequestContext } from '@playwright/test';
import { test as base } from '../../fixtures/demo';
import { rootApi, ensureMount } from '../../fixtures/root-api';

const test = base.extend<{ rootApi: APIRequestContext; backendPath: string }>({
  rootApi: async ({ playwright }, use, testInfo) => {
    const api = await rootApi(playwright, testInfo);
    await use(api);
    await api.dispose();
  },
  backendPath: async ({ rootApi }, use) => {
    const path = `aws-ui-${randomUUID()}`;
    await ensureMount(rootApi, path, 'aws');
    try {
      await use(path);
    } finally {
      await expect(await rootApi.delete(`/v1/sys/mounts/${path}`)).toBeOK();
    }
  },
});

const roleTypes = ['iam_user', 'assumed_role', 'federation_token', 'session_token'] as const;
const policy = {
  Version: '2012-10-17',
  Statement: [{ Effect: 'Allow', Action: 's3:GetObject', Resource: '*' }],
};

for (const credentialType of roleTypes) {
  test(`aws ${credentialType} role create, edit, read-back and details deletion`, async ({
    page,
    rootApi,
    backendPath,
  }) => {
    const name = 'workflow-role';
    const roleUrl = `/v1/${backendPath}/roles/${name}`;
    const arn =
      credentialType === 'iam_user'
        ? 'arn:aws:iam::aws:policy/ReadOnlyAccess'
        : 'arn:aws:iam::123456789012:role/Developer';

    await test.step('create the role through its rendered inputs', async () => {
      await page.goto(`secrets-engines/${backendPath}/create`);
      await expect(page.getByRole('heading', { name: 'Create an AWS Role', exact: true })).toBeVisible();
      await page.getByLabel('Role name', { exact: true }).fill(name);
      await page.getByLabel('Credential type', { exact: true }).selectOption(credentialType);
      if (credentialType === 'iam_user' || credentialType === 'assumed_role') {
        await page.locator('[data-test-string-list-input="0"]').fill(arn);
      }
      if (credentialType !== 'session_token') {
        await page
          .locator('[data-test-component="code-mirror-modifier"] [contenteditable="true"]')
          .fill(JSON.stringify(policy));
      }
      const saved = page.waitForResponse(
        (response) => response.url().endsWith(roleUrl) && response.request().method() === 'POST'
      );
      await page.getByRole('button', { name: 'Create role', exact: true }).click();
      expect((await saved).ok()).toBe(true);
      await expect(page).toHaveURL(new RegExp(`/show/${name}$`));
    });

    await test.step('independently read persisted values and reload details', async () => {
      const response = await rootApi.get(roleUrl);
      await expect(response).toBeOK();
      const { data } = await response.json();
      expect(data.credential_type).toBe(credentialType);
      if (credentialType === 'iam_user') expect(data.policy_arns).toEqual([arn]);
      if (credentialType === 'assumed_role') expect(data.role_arns).toEqual([arn]);
      if (credentialType !== 'session_token') expect(JSON.parse(data.policy_document)).toEqual(policy);
      await page.reload();
      await expect(page.locator('[data-test-row-value="Credential type"]').first()).toHaveText(
        credentialType
      );
      if (credentialType === 'iam_user' || credentialType === 'assumed_role') {
        await expect(
          page
            .locator(`[data-test-row-value="${credentialType === 'iam_user' ? 'Policy ARNs' : 'Role ARNs'}"]`)
            .first()
        ).toContainText(arn);
      }
    });

    await test.step('edit the type and verify persistence after reloading', async () => {
      const updatedType = credentialType === 'session_token' ? 'federation_token' : credentialType;
      const updatedPolicy = {
        ...policy,
        Statement: [{ ...policy.Statement[0], Action: 's3:ListBucket' }],
      };
      const updatedArn = arn.replace(/ReadOnlyAccess|Developer/, 'Updated');
      await page.getByRole('link', { name: 'Edit role', exact: true }).click();
      await page.getByLabel('Credential type', { exact: true }).selectOption(updatedType);
      if (credentialType === 'iam_user' || credentialType === 'assumed_role') {
        await page.locator('[data-test-string-list-input="0"]').fill(updatedArn);
      }
      await page
        .locator('[data-test-component="code-mirror-modifier"] [contenteditable="true"]')
        .fill(JSON.stringify(updatedPolicy));
      const saved = page.waitForResponse(
        (response) => response.url().endsWith(roleUrl) && response.request().method() === 'POST'
      );
      await page.getByRole('button', { name: 'Save', exact: true }).click();
      expect((await saved).ok()).toBe(true);
      await expect(page).toHaveURL(new RegExp(`/show/${name}$`));
      const response = await rootApi.get(roleUrl);
      await expect(response).toBeOK();
      const { data } = await response.json();
      expect(data.credential_type).toBe(updatedType);
      expect(data.policy_arns || []).toEqual(credentialType === 'iam_user' ? [updatedArn] : []);
      expect(data.role_arns || []).toEqual(credentialType === 'assumed_role' ? [updatedArn] : []);
      expect(JSON.parse(data.policy_document)).toEqual(updatedPolicy);
      await page.reload();
      await expect(page.locator('[data-test-row-value="Credential type"]').first()).toHaveText(updatedType);
    });

    await test.step('delete from details and verify the role stays deleted', async () => {
      const deleted = page.waitForResponse(
        (response) => response.url().endsWith(roleUrl) && response.request().method() === 'DELETE'
      );
      await page.getByRole('button', { name: 'Delete role', exact: true }).click();
      await page.getByRole('button', { name: 'Confirm', exact: true }).click();
      expect((await deleted).ok()).toBe(true);
      await expect(page).toHaveURL(new RegExp(`/${backendPath}/list$`));
      expect((await rootApi.get(roleUrl)).status()).toBe(404);
      await page.reload();
      await expect(page.locator(`[data-test-secret-link="${name}"]`)).toHaveCount(0);
      await expect(
        page.getByRole('heading', { name: 'No roles in this backend', exact: true })
      ).toBeVisible();
    });
  });
}

test('aws credential generation renders real backend validation failures and permits retry', async ({
  page,
  rootApi,
  backendPath,
}) => {
  const name = 'assumed-role';
  await expect(
    await rootApi.post(`/v1/${backendPath}/roles/${name}`, {
      data: {
        credential_type: 'assumed_role',
        role_arns: ['arn:aws:iam::123456789012:role/Allowed'],
      },
    })
  ).toBeOK();
  await page.goto(`secrets-engines/${backendPath}/credentials/${name}`);
  await expect(page.getByRole('heading', { name: 'Generate AWS Credentials', exact: true })).toBeVisible();
  await expect(page.getByText('Generating credentials of type: assumed_role', { exact: true })).toBeVisible();
  const roleArn = page.getByLabel('Role ARN', { exact: true });
  await roleArn.fill('arn:aws:iam::123456789012:role/NotAllowed');

  for (let attempt = 0; attempt < 2; attempt++) {
    const generated = page.waitForResponse(
      (response) =>
        response.url().endsWith(`/v1/${backendPath}/creds/${name}`) && response.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Generate', exact: true }).click();
    expect((await generated).status()).toBe(400);
    await expect(page.getByRole('alert')).toContainText('not in allowed role arns');
    await expect(roleArn).toHaveValue('arn:aws:iam::123456789012:role/NotAllowed');
    await expect(page.getByRole('button', { name: 'Generate', exact: true })).toBeEnabled();
  }
});
