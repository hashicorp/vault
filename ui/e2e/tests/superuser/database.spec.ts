/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { randomUUID } from 'crypto';
import { expect, type APIRequestContext, type Page } from '@playwright/test';
import { test as base } from '../../fixtures/demo';
import { rootApi, ensureMount } from '../../fixtures/root-api';

const test = base.extend<{ rootApi: APIRequestContext; backendPath: string }>({
  rootApi: async ({ playwright }, use, testInfo) => {
    const api = await rootApi(playwright, testInfo);
    await use(api);
    await api.dispose();
  },
  backendPath: async ({ rootApi }, use) => {
    const path = `database-ui-${randomUUID()}`;
    await ensureMount(rootApi, path, 'database');
    try {
      await use(path);
    } finally {
      await expect(await rootApi.delete(`/v1/sys/mounts/${path}`)).toBeOK();
    }
  },
});

const PLUGIN = 'mongodb-database-plugin';
const CONNECTION_URL = 'mongodb://127.0.0.1:4321/horses';

const isConfigRequest = (
  url: string,
  backendPath: string,
  name: string,
  method: string,
  requestMethod: string
) => url.endsWith(`/v1/${backendPath}/config/${name}`) && requestMethod === method;

// A real database is not reachable from the test environment, so connections are always created
// with verify_connection unchecked. Vault then stores the config without dialing the database.
const createConnection = async (page: Page, backendPath: string, name: string) => {
  await page.goto(`secrets-engines/${backendPath}/create`);
  await page.locator('[data-test-input="plugin_name"]').selectOption(PLUGIN);
  await page.locator('[data-test-input="name"]').fill(name);
  await page.locator('[data-test-input="connection_url"]').fill(CONNECTION_URL);
  await page.locator('[data-test-input="username"]').fill('user');
  await page.locator('[data-test-input="password"]').fill('so-secure');
  await page.locator('[data-test-input="verify_connection"]').click();
  const saved = page.waitForResponse((r) =>
    isConfigRequest(r.url(), backendPath, name, 'POST', r.request().method())
  );
  await page.getByRole('button', { name: 'Create database', exact: true }).click();
  const response = await saved;
  expect(response.ok()).toBe(true);
  return response;
};

test('database connection create, show, edit, list and delete', async ({ page, rootApi, backendPath }) => {
  const name = 'horses-db';
  const connectionUrl = `/v1/${backendPath}/config/${name}`;

  await test.step('shows the empty state before any connection exists', async () => {
    await page.goto(`secrets-engines/${backendPath}/list`);
    await expect(
      page.getByRole('heading', { name: 'No connections in this backend', exact: true })
    ).toBeVisible();
  });

  await test.step('creates the connection and skips root credential rotation', async () => {
    await page.goto(`secrets-engines/${backendPath}/create`);
    await expect(page.getByRole('heading', { name: 'No plugin selected' }).first()).toBeVisible();
    await page.locator('[data-test-input="plugin_name"]').selectOption(PLUGIN);
    await page.locator('[data-test-input="name"]').fill(name);
    await page.locator('[data-test-input="connection_url"]').fill(CONNECTION_URL);
    await page.locator('[data-test-input="username"]').fill('user');
    await page.locator('[data-test-input="password"]').fill('so-secure');
    await expect(page.locator('[data-test-input="verify_connection"]')).toBeChecked();
    await page.locator('[data-test-input="verify_connection"]').click();
    const saved = page.waitForResponse((r) =>
      isConfigRequest(r.url(), backendPath, name, 'POST', r.request().method())
    );
    await page.getByRole('button', { name: 'Create database', exact: true }).click();
    expect((await saved).ok()).toBe(true);
    await expect(page.getByRole('dialog')).toContainText('Rotate your root credentials?');
    await page.getByRole('button', { name: 'Enable without rotating', exact: true }).click();
    await expect(page).toHaveURL(new RegExp(`/show/${name}$`));
  });

  await test.step('independently reads the persisted config and reloads the details page', async () => {
    const response = await rootApi.get(connectionUrl);
    await expect(response).toBeOK();
    const { data } = await response.json();
    expect(data.plugin_name).toBe(PLUGIN);
    expect(data.connection_details.connection_url).toBe(CONNECTION_URL);
    expect(data.connection_details.username).toBe('user');
    expect(data.connection_details.password).toBeUndefined();
    await page.reload();
    await expect(page.locator('span[data-test-row-value="Connection URL"]')).toHaveText(CONNECTION_URL);
    await expect(page.locator('[data-test-row-value="Password"]')).toHaveCount(0);
  });

  await test.step('lists the connection', async () => {
    await page.goto(`secrets-engines/${backendPath}/list`);
    await expect(page.getByRole('link', { name: new RegExp(`^${name} `) })).toBeVisible();
  });

  await test.step('edits without sending the unchanged password', async () => {
    await page.goto(`secrets-engines/${backendPath}/show/${name}`);
    await page.getByRole('link', { name: 'Edit configuration', exact: true }).click();
    await expect(page).toHaveURL(new RegExp(`/edit/${name}$`));
    await expect(page.locator('[data-test-input="name"]')).toHaveAttribute('readonly', '');
    await expect(page.locator('[data-test-input="plugin_name"]')).toHaveAttribute('readonly', '');
    await expect(page.locator('[data-test-input="password"]')).toHaveCount(0);
    await expect(page.getByText('Update password', { exact: true })).toBeVisible();
    const updatedUrl = 'mongodb://127.0.0.1:4321/zebras';
    await page.locator('[data-test-input="connection_url"]').fill(updatedUrl);
    const saved = page.waitForResponse((r) =>
      isConfigRequest(r.url(), backendPath, name, 'POST', r.request().method())
    );
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    const response = await saved;
    expect(response.ok()).toBe(true);
    const sent = response.request().postDataJSON();
    expect(sent.connection_url).toBe(updatedUrl);
    expect(sent.password).toBeUndefined();
    await expect(page).toHaveURL(new RegExp(`/show/${name}$`));
    const { data } = await (await rootApi.get(connectionUrl)).json();
    expect(data.connection_details.connection_url).toBe(updatedUrl);
    await page.reload();
    await expect(page.locator('span[data-test-row-value="Connection URL"]')).toHaveText(updatedUrl);
  });

  await test.step('requires typing the connection name to delete and stays deleted', async () => {
    await page.goto(`secrets-engines/${backendPath}/show/${name}`);
    await page.getByRole('button', { name: 'Delete connection', exact: true }).click();
    await expect(page.getByRole('dialog')).toContainText('Delete connection?');
    await page.getByRole('dialog').getByRole('textbox').fill(name);
    const deleted = page.waitForResponse((r) =>
      isConfigRequest(r.url(), backendPath, name, 'DELETE', r.request().method())
    );
    await page.getByRole('button', { name: 'Confirm', exact: true }).click();
    expect((await deleted).ok()).toBe(true);
    await expect(page).toHaveURL(new RegExp(`/${backendPath}/list$`));
    expect((await rootApi.get(connectionUrl)).status()).toBe(404);
    await page.reload();
    await expect(page.getByRole('link', { name: new RegExp(`^${name} `) })).toHaveCount(0);
    await expect(
      page.getByRole('heading', { name: 'No connections in this backend', exact: true })
    ).toBeVisible();
  });
});

test('database connection edit sends a new password only when the toggle is enabled', async ({
  page,
  backendPath,
}) => {
  const name = 'password-db';
  await createConnection(page, backendPath, name);
  await page.getByRole('button', { name: 'Enable without rotating', exact: true }).click();
  await expect(page).toHaveURL(new RegExp(`/show/${name}$`));

  await page.getByRole('link', { name: 'Edit configuration', exact: true }).click();
  await page.getByText('Update password', { exact: true }).click();
  await page.getByLabel('New password', { exact: true }).fill('new-secret');
  const saved = page.waitForResponse((r) =>
    isConfigRequest(r.url(), backendPath, name, 'POST', r.request().method())
  );
  await page.getByRole('button', { name: 'Save', exact: true }).click();
  const response = await saved;
  expect(response.ok()).toBe(true);
  expect(response.request().postDataJSON().password).toBe('new-secret');
  await expect(page).toHaveURL(new RegExp(`/show/${name}$`));
});

test('database connection create renders the real backend error and permits retry', async ({
  page,
  rootApi,
  backendPath,
}) => {
  // With verify_connection left checked, Vault dials the unreachable database and rejects the save.
  const name = 'unreachable-db';
  await page.goto(`secrets-engines/${backendPath}/create`);
  await page.locator('[data-test-input="plugin_name"]').selectOption(PLUGIN);
  await page.locator('[data-test-input="name"]').fill(name);
  await page.locator('[data-test-input="connection_url"]').fill('mongodb://127.0.0.1:1/none');
  await page.locator('[data-test-input="username"]').fill('user');
  await page.locator('[data-test-input="password"]').fill('so-secure');

  for (let attempt = 0; attempt < 2; attempt++) {
    const saved = page.waitForResponse((r) =>
      isConfigRequest(r.url(), backendPath, name, 'POST', r.request().method())
    );
    await page.getByRole('button', { name: 'Create database', exact: true }).click();
    expect((await saved).ok()).toBe(false);
    await expect(
      page.getByRole('alert').filter({ hasText: 'error creating database object' }).first()
    ).toBeVisible();
    await expect(page.locator('[data-test-input="name"]')).toHaveValue(name);
    await expect(page.getByRole('button', { name: 'Create database', exact: true })).toBeEnabled();
  }
  expect((await rootApi.get(`/v1/${backendPath}/config/${name}`)).status()).toBe(404);
});

test('database connection reset and rotate root call the backend and surface its response', async ({
  page,
  backendPath,
}) => {
  // Vault dials the unreachable database before failing, which exceeds the default timeout.
  test.setTimeout(180_000);
  // Reset and rotate-root need a reachable database for a success response. The UI contract under
  // test is that each action hits the right endpoint and renders the backend's result.
  const name = 'actions-db';
  await createConnection(page, backendPath, name);
  await page.getByRole('button', { name: 'Enable without rotating', exact: true }).click();
  await expect(page).toHaveURL(new RegExp(`/show/${name}$`));

  await test.step('reset connection', async () => {
    const reset = page.waitForResponse(
      (r) => r.url().endsWith(`/v1/${backendPath}/reset/${name}`) && r.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Reset connection', exact: true }).click();
    await page.getByRole('button', { name: 'Confirm', exact: true }).click();
    const response = await reset;
    if (response.ok()) {
      await expect(page.getByText('Successfully reset connection')).toBeVisible();
    } else {
      await expect(page.getByRole('alert').filter({ hasText: /\S/ }).first()).toBeVisible();
    }
  });

  await test.step('rotate root credentials', async () => {
    const rotate = page.waitForResponse(
      (r) => r.url().endsWith(`/v1/${backendPath}/rotate-root/${name}`) && r.request().method() === 'POST',
      { timeout: 60_000 }
    );
    await page.getByRole('button', { name: 'Rotate root credentials', exact: true }).click();
    await page.getByRole('button', { name: 'Confirm', exact: true }).click();
    const response = await rotate;
    if (response.ok()) {
      await expect(page.getByText('Successfully rotated credentials')).toBeVisible();
    } else {
      await expect(page.getByRole('alert').filter({ hasText: /\S/ }).first()).toBeVisible();
    }
  });
});

test('database connection create can attempt root credential rotation from the modal', async ({
  page,
  backendPath,
}) => {
  test.setTimeout(120_000);
  const name = 'rotate-db';
  await createConnection(page, backendPath, name);
  const rotate = page.waitForResponse(
    (r) => r.url().endsWith(`/v1/${backendPath}/rotate-root/${name}`) && r.request().method() === 'POST',
    { timeout: 60_000 }
  );
  await page.getByRole('button', { name: 'Rotate and enable', exact: true }).click();
  const response = await rotate;
  if (response.ok()) {
    await expect(
      page.getByText(`Successfully rotated root credentials for connection "${name}"`)
    ).toBeVisible();
  } else {
    await expect(page.getByText(/Error rotating root credentials/)).toBeVisible();
  }
  // Either way the user lands on the connection details page.
  await expect(page).toHaveURL(new RegExp(`/show/${name}$`));
});

const ROLE_CONNECTION = 'role-conn';

const isRoleRequest = (
  url: string,
  backendPath: string,
  kind: string,
  name: string,
  requestMethod: string,
  method: string
) => url.endsWith(`/v1/${backendPath}/${kind}/${name}`) && requestMethod === method;

// Roles are saved against a connection, so one is seeded through the API. Vault stores the config
// without dialing the database because verify_connection is false.
const seedConnection = async (rootApi: APIRequestContext, backendPath: string) => {
  await expect(
    await rootApi.post(`/v1/${backendPath}/config/${ROLE_CONNECTION}`, {
      data: {
        plugin_name: PLUGIN,
        connection_url: CONNECTION_URL,
        username: 'user',
        password: 'so-secure',
        verify_connection: false,
      },
    })
  ).toBeOK();
};

const allowedRoles = async (rootApi: APIRequestContext, backendPath: string) => {
  const { data } = await (await rootApi.get(`/v1/${backendPath}/config/${ROLE_CONNECTION}`)).json();
  return data.allowed_roles as string[];
};

const selectConnection = async (page: Page) => {
  await page.locator('[data-test-component="search-select"]#database .ember-basic-dropdown-trigger').click();
  await page.locator('.ember-power-select-option', { hasText: ROLE_CONNECTION }).click();
};

// The code editor loads lazily once it scrolls into view, so scroll before looking for its input.
const statementEditor = async (page: Page) => {
  await page.locator('[data-test-component="code-mirror-modifier"]').first().scrollIntoViewIfNeeded();
  return page.locator('[data-test-component="code-mirror-modifier"] [contenteditable="true"]').first();
};

test('database dynamic role create, show, edit, list and delete', async ({ page, rootApi, backendPath }) => {
  test.setTimeout(120_000);
  await seedConnection(rootApi, backendPath);
  const name = 'dynamic-role';
  const roleUrl = `/v1/${backendPath}/roles/${name}`;
  const statement = '{"db":"admin","roles":[{"role":"readWrite"}]}';

  await test.step('shows the empty state before any role exists', async () => {
    await page.goto(`secrets-engines/${backendPath}/list?tab=role`);
    await expect(page.getByRole('heading', { name: /No roles in this backend/ })).toBeVisible();
  });

  await test.step('creates the role and adds it to the connection allowed roles', async () => {
    await page.goto(`secrets-engines/${backendPath}/create?itemType=role`);
    await page.locator('[data-test-input="name"]').fill(name);
    await selectConnection(page);
    await page.locator('[data-test-input="type"]').selectOption('dynamic');
    await (await statementEditor(page)).fill(statement);
    const saved = page.waitForResponse((r) =>
      isRoleRequest(r.url(), backendPath, 'roles', name, r.request().method(), 'POST')
    );
    await page.locator('form').getByRole('button', { name: 'Create role', exact: true }).click();
    expect((await saved).ok()).toBe(true);
    await expect(page).toHaveURL(new RegExp(`/show/role/${name}$`));
  });

  await test.step('independently reads the persisted role and reloads the details page', async () => {
    const response = await rootApi.get(roleUrl);
    await expect(response).toBeOK();
    const { data } = await response.json();
    expect(data.db_name).toBe(ROLE_CONNECTION);
    expect(data.creation_statements).toEqual([statement]);
    expect(await allowedRoles(rootApi, backendPath)).toContain(name);
    await page.reload();
    await expect(page.locator('span[data-test-row-value="Type of role"]')).toHaveText('dynamic');
    await expect(page.locator('span[data-test-row-value="Role name"]')).toHaveText(name);
    await expect(page.getByRole('button', { name: /Generate credentials/ })).toBeVisible();
  });

  await test.step('lists the role', async () => {
    await page.goto(`secrets-engines/${backendPath}/list?tab=role`);
    await expect(page.locator('.list-item-row').filter({ hasText: name })).toBeVisible();
  });

  await test.step('edits the creation statement and persists it', async () => {
    await page.goto(`secrets-engines/${backendPath}/show/role/${name}`);
    await page.getByRole('link', { name: 'Edit role', exact: true }).click();
    await expect(page).toHaveURL(new RegExp(`/edit/role/${name}`));
    const updated = '{"db":"admin","roles":[{"role":"read"}]}';
    await (await statementEditor(page)).fill(updated);
    const saved = page.waitForResponse((r) =>
      isRoleRequest(r.url(), backendPath, 'roles', name, r.request().method(), 'POST')
    );
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    expect((await saved).ok()).toBe(true);
    await expect(page).toHaveURL(new RegExp(`/show/role/${name}$`));
    const { data } = await (await rootApi.get(roleUrl)).json();
    expect(data.creation_statements).toEqual([updated]);
    expect(data.db_name).toBe(ROLE_CONNECTION);
  });

  await test.step('deletes the role, removes it from allowed roles and stays deleted', async () => {
    await page.goto(`secrets-engines/${backendPath}/show/role/${name}`);
    await page.getByRole('button', { name: 'Delete role', exact: true }).click();
    const deleted = page.waitForResponse((r) =>
      isRoleRequest(r.url(), backendPath, 'roles', name, r.request().method(), 'DELETE')
    );
    await page.getByRole('button', { name: 'Confirm', exact: true }).click();
    expect((await deleted).ok()).toBe(true);
    await expect(page).toHaveURL(new RegExp(`/${backendPath}/list\\?tab=role`));
    expect((await rootApi.get(roleUrl)).status()).toBe(404);
    expect(await allowedRoles(rootApi, backendPath)).not.toContain(name);
    await page.reload();
    await expect(page.locator('.list-item-row').filter({ hasText: name })).toHaveCount(0);
  });
});

test('database static role create, edit and rotate never expose the password', async ({
  page,
  rootApi,
  backendPath,
}) => {
  test.setTimeout(120_000);
  await seedConnection(rootApi, backendPath);
  const name = 'static-role';
  const roleUrl = `/v1/${backendPath}/static-roles/${name}`;

  await test.step('creates a static role without rotating immediately', async () => {
    await page.goto(`secrets-engines/${backendPath}/create?itemType=role`);
    await page.locator('[data-test-input="name"]').fill(name);
    await selectConnection(page);
    await page.locator('[data-test-input="type"]').selectOption('static');
    await page.locator('[data-test-input="username"]').fill('static_user');
    await page.locator('[data-test-toggle-label="toggle-skip_import_rotation"]').click();
    const saved = page.waitForResponse((r) =>
      isRoleRequest(r.url(), backendPath, 'static-roles', name, r.request().method(), 'POST')
    );
    await page.locator('form').getByRole('button', { name: 'Create role', exact: true }).click();
    const response = await saved;
    expect(response.ok()).toBe(true);
    expect(response.request().postDataJSON().username).toBe('static_user');
    expect(response.request().postDataJSON().skip_import_rotation).toBe(true);
    await expect(page).toHaveURL(new RegExp(`/show/role/${name}$`));
  });

  await test.step('independently reads the persisted role and reloads the details page', async () => {
    const response = await rootApi.get(roleUrl);
    await expect(response).toBeOK();
    const { data } = await response.json();
    expect(data.db_name).toBe(ROLE_CONNECTION);
    expect(data.username).toBe('static_user');
    expect(data.skip_import_rotation).toBe(true);
    expect(data.password).toBeUndefined();
    expect(await allowedRoles(rootApi, backendPath)).toContain(name);
    await page.reload();
    await expect(page.locator('span[data-test-row-value="Type of role"]')).toHaveText('static');
    await expect(page.locator('span[data-test-row-value="Username"]')).toHaveText('static_user');
    await expect(page.locator('span[data-test-row-value="Password"]')).toHaveCount(0);
    await expect(page.getByRole('button', { name: 'Rotate credentials', exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: /Get credentials/ })).toBeVisible();
  });

  await test.step('lists the role as a static role', async () => {
    await page.goto(`secrets-engines/${backendPath}/list?tab=role`);
    await expect(page.locator('.list-item-row').filter({ hasText: name })).toBeVisible();
    await expect(page.locator('.list-item-row').filter({ hasText: name })).toContainText('static');
  });

  await test.step('edits the rotation period and sends the username but never the password', async () => {
    await page.goto(`secrets-engines/${backendPath}/show/role/${name}`);
    await page.getByRole('link', { name: 'Edit role', exact: true }).click();
    await expect(page).toHaveURL(new RegExp(`/edit/role/${name}`));
    await page.locator('[data-test-ttl-value="Rotation period"]').fill('2');
    const saved = page.waitForResponse((r) =>
      isRoleRequest(r.url(), backendPath, 'static-roles', name, r.request().method(), 'POST')
    );
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    const response = await saved;
    expect(response.ok()).toBe(true);
    const sent = response.request().postDataJSON();
    expect(sent.username).toBe('static_user');
    expect(sent.rotation_period).toBe('172800s');
    expect(sent.password).toBeUndefined();
    await expect(page).toHaveURL(new RegExp(`/show/role/${name}$`));
    await expect(page.locator('span[data-test-row-value="Password"]')).toHaveCount(0);
    const { data } = await (await rootApi.get(roleUrl)).json();
    expect(data.rotation_period).toBe(172800);
    expect(data.username).toBe('static_user');
    expect(data.password).toBeUndefined();
  });

  await test.step('edits the password and never displays it on the details page', async () => {
    const password = 'Distinct-Pw-VAULT-35458';
    await page.goto(`secrets-engines/${backendPath}/edit/role/${name}`);
    await page.getByRole('button', { name: 'Enable input', exact: true }).click();
    await page.locator('[data-test-input="password"]').fill(password);
    const saved = page.waitForResponse((r) =>
      isRoleRequest(r.url(), backendPath, 'static-roles', name, r.request().method(), 'POST')
    );
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    const response = await saved;
    expect(response.ok()).toBe(true);
    expect(response.request().postDataJSON().password).toBe(password);
    await expect(page).toHaveURL(new RegExp(`/show/role/${name}$`));
    await expect(page.locator('span[data-test-row-value="Username"]')).toHaveText('static_user');
    await expect(page.getByText(password)).toHaveCount(0);
    await expect(page.locator('span[data-test-row-value="Password"]')).toHaveCount(0);
    await page.reload();
    await expect(page.locator('span[data-test-row-value="Username"]')).toHaveText('static_user');
    await expect(page.getByText(password)).toHaveCount(0);
    await expect(page.locator('span[data-test-row-value="Password"]')).toHaveCount(0);
  });

  await test.step('rotate credentials calls the rotate-role endpoint and surfaces its response', async () => {
    // Rotating needs a reachable database for success, so the contract is the endpoint and the rendered result.
    const rotate = page.waitForResponse(
      (r) => r.url().endsWith(`/v1/${backendPath}/rotate-role/${name}`) && r.request().method() === 'POST',
      { timeout: 90_000 }
    );
    await page.getByRole('button', { name: 'Rotate credentials', exact: true }).click();
    const response = await rotate;
    if (response.ok()) {
      await expect(page.getByText(`Credentials for ${name} role were rotated`)).toBeVisible();
    } else {
      await expect(page.getByRole('alert').filter({ hasText: /\S/ }).first()).toBeVisible();
    }
  });

  await test.step('rotate credentials from the list menu calls the same endpoint', async () => {
    await page.goto(`secrets-engines/${backendPath}/list?tab=role`);
    await page
      .locator('.list-item-row')
      .filter({ hasText: name })
      .getByRole('button', { name: `Manage database ${name}` })
      .click();
    const rotate = page.waitForResponse(
      (r) => r.url().endsWith(`/v1/${backendPath}/rotate-role/${name}`) && r.request().method() === 'POST',
      { timeout: 90_000 }
    );
    await page.getByRole('button', { name: 'Rotate credentials', exact: true }).click();
    await rotate;
  });

  await test.step('deletes the static role', async () => {
    await page.goto(`secrets-engines/${backendPath}/show/role/${name}`);
    await page.getByRole('button', { name: 'Delete role', exact: true }).click();
    const deleted = page.waitForResponse((r) =>
      isRoleRequest(r.url(), backendPath, 'static-roles', name, r.request().method(), 'DELETE')
    );
    await page.getByRole('button', { name: 'Confirm', exact: true }).click();
    expect((await deleted).ok()).toBe(true);
    expect((await rootApi.get(roleUrl)).status()).toBe(404);
    expect(await allowedRoles(rootApi, backendPath)).not.toContain(name);
  });
});

test('database role create validates required fields without sending a request', async ({
  page,
  rootApi,
  backendPath,
}) => {
  await seedConnection(rootApi, backendPath);
  const writes: string[] = [];
  page.on('request', (request) => {
    if (request.method() === 'POST' && request.url().includes(`/v1/${backendPath}/`)) {
      writes.push(request.url());
    }
  });
  await page.goto(`secrets-engines/${backendPath}/create?itemType=role`);
  await page.locator('form').getByRole('button', { name: 'Create role', exact: true }).click();
  await expect(page.getByText('Role name is required.')).toBeVisible();
  await expect(page.getByText('There was an error submitting this form.')).toBeVisible();
  expect(writes).toEqual([]);
});
