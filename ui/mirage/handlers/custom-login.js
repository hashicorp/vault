/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { sanitizePath } from 'core/utils/sanitize-path';

export default function (server) {
  // LIST, READ and DELETE requests for default-auth (login customizations)
  server.get('sys/config/ui/login/default-auth', (schema, req) => {
    // API expects { data: { list: true } } as query params when making LIST requests
    if (req.queryParams.list) {
      const records = schema.db['loginRules'];
      if (records) {
        const keys = records.map(({ name }) => name);
        const key_info = records.reduce((obj, record) => {
          const { name, namespace_path, disable_inheritance } = record;
          obj[name] = { namespace_path, disable_inheritance, name };
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
    const findRule = (ns) => schema.db['loginRules'].findBy({ namespace_path: ns });
    // the namespace header shouldn't have a trailing slash, but sanitize just in case
    //  if no namespace is passed then it's the root namespace (which does not have a trailing slash)
    const nsPath = sanitizePath(nsHeader) ? `${nsHeader}/` : 'root';
    let rule = findRule(nsPath);

    if (!rule && nsHeader?.includes('/')) {
      // for simplicity, tests only support namespaces nested one level, e.g. "test-ns/child"
      const [parent] = nsHeader.split('/');
      const parentRule = findRule(`${parent}/`);
      rule = parentRule?.disable_inheritance ? null : parentRule;
    }

    // Fallback to root namespace settings to simulate inheritance if no rule exists or parent has disabled inheritance
    rule = rule || findRule('root');

    const { default_auth_type, backup_auth_types, disable_inheritance } = rule || {};
    return { data: { default_auth_type, backup_auth_types, disable_inheritance } };
  });
}
