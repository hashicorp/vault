/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  auth_method_accessors: () => ['auth_userpass_bb95c2b1'],
  auth_method_types: () => ['userpass'],
  identity_entity_ids: () => ['f831667b-7392-7a1c-c0fc-33d48cb1c57d'],
  identity_group_ids: () => ['34db6b52-591e-bc22-8af0-4add5e167326'],
  mfa_method_ids: null,
  name: () => (Math.random() + 1).toString(36).substring(2),
  namespace_id: 'root',

  afterCreate(record, server) {
    const rawMethodIds =
      record.mfa_method_ids ||
      record.mfaMethodIds ||
      record.attrs?.mfa_method_ids ||
      record.attrs?.mfaMethodIds;

    if (rawMethodIds && rawMethodIds.length) {
      record.mfa_method_ids = rawMethodIds;
    } else {
      // aggregate all existing methods and choose a random one
      const methods = ['Totp', 'Duo', 'Okta', 'Pingid'].reduce((methods, type) => {
        const collection = server.db[`mfa${type}Methods`];
        const records = collection ? collection.where({}) : [];
        if (records.length) {
          methods.push(...records);
        }
        return methods;
      }, []);
      // if no methods were found create one since it is a required for login enforcements
      if (!methods.length) {
        methods.push(
          server.create('mfa-totp-method', {
            id: (Math.random() + 1).toString(36).substring(2),
          })
        );
      }
      const method = methods.length ? methods[Math.floor(Math.random() * methods.length)] : null;
      const methodId = method ? method.id : null;
      if (methodId) {
        record.mfa_method_ids = [methodId];
      }
    }

    const authMethodAccessors = record.auth_method_accessors ||
      record.authMethodAccessors ||
      record.attrs?.auth_method_accessors ||
      record.attrs?.authMethodAccessors || ['auth_userpass_bb95c2b1'];

    const authMethodTypes = record.auth_method_types ||
      record.authMethodTypes ||
      record.attrs?.auth_method_types ||
      record.attrs?.authMethodTypes || ['userpass'];

    const identityEntityIds = record.identity_entity_ids ||
      record.identityEntityIds ||
      record.attrs?.identity_entity_ids ||
      record.attrs?.identityEntityIds || ['f831667b-7392-7a1c-c0fc-33d48cb1c57d'];

    const identityGroupIds = record.identity_group_ids ||
      record.identityGroupIds ||
      record.attrs?.identity_group_ids ||
      record.attrs?.identityGroupIds || ['34db6b52-591e-bc22-8af0-4add5e167326'];

    if (server.db.mfaLoginEnforcements && record.id) {
      server.db.mfaLoginEnforcements.update(record.id, {
        id: record.id,
        name: record.name || record.attrs?.name,
        mfa_method_ids: record.mfa_method_ids,
        auth_method_accessors: authMethodAccessors,
        auth_method_types: authMethodTypes,
        identity_entity_ids: identityEntityIds,
        identity_group_ids: identityGroupIds,
      });
    }
  },
});
