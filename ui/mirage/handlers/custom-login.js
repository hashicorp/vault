/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  // LIST, READ and DELETE requests for default-auth (login customizations)
  server.get('sys/config/ui/login/default-auth', (schema, req) => {
    // API expects { data: { list: true } } as query params when making LIST requests
    if (req.queryParams.list) {
      const records = schema.db['loginRules'];
      if (records) {
        const keys = records.map(({ name }) => name);
        const key_info = records.reduce((obj, record) => {
          const { name, namespace, disable_inheritance } = record;
          // TBD, but likely only limited information will be returned about the record from the LIST request
          obj[name] = { namespace, disable_inheritance };
          return obj;
        }, {});
        return {
          data: { keys, key_info },
        };
      }
      return new Response(404, {}, { errors: [] });
    }
  });

  server.get('sys/config/ui/login/default-auth/:name', (schema, req) => {
    // req.params come in as: { name: "Login rule name" }
    const record = schema.db['loginRules'].findBy(req.params);
    if (record) {
      delete record.id; // "name" is the id
      return { data: record };
    }
    return new Response(404, {}, { errors: [] });
  });

  server.delete('sys/config/ui/login/default-auth/:name', (schema, req) => {
    const record = schema.db['loginRules'].findBy(req.params);
    if (record) {
      schema.db['loginRules'].remove(record);
      return new Response(204); // No content
    }
    return new Response(404, {}, { errors: [] });
  });

  // UNAUTHENTICATED READ ONLY for login form display logic
  server.get('sys/internal/ui/default-auth-methods', (schema, req) => {
    const nsHeader = req.requestHeaders['X-Vault-Namespace'];
    // if no namespace is passed, assume root
    const findRule = (ns = '') => schema.db['loginRules'].findBy({ namespace_path: ns });

    let rule = findRule(nsHeader || '');

    if (!rule && nsHeader?.includes('/')) {
      // for simplicity, tests only nest namespaces one level, e.g. "test-ns/child"
      const [parent] = nsHeader.split('/');
      const parentRule = findRule(parent);
      rule = parentRule?.disable_inheritance ? null : parentRule;
    }

    // Fallback to root namespace settings to simulate inheritance if no rule exists or parent has disabled inheritance
    rule = rule || findRule();

    const { default_auth_type, backup_auth_types, disable_inheritance } = rule || {};
    return { data: { default_auth_type, backup_auth_types, disable_inheritance } };
  });
}
