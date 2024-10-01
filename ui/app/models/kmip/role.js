/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import apiPath from 'vault/utils/api-path';
import lazyCapabilities from 'vault/macros/lazy-capabilities';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';
import {
  operationFields,
  operationFieldsWithoutSpecial,
  tlsFields,
} from 'vault/utils/model-helpers/kmip-role-fields';
import { removeManyFromArray } from 'vault/helpers/remove-from-array';

@withExpandedAttributes()
export default class KmipRoleModel extends Model {
  @attr({ readOnly: true }) backend;
  @attr({ readOnly: true }) scope;

  get editableFields() {
    return Object.keys(this.allByKey).filter((k) => !['backend', 'scope', 'role'].includes(k));
  }

  get fieldGroups() {
    const tls = tlsFields();
    const groups = [{ TLS: tls }];
    // op fields are shown in OperationFieldDisplay
    const opFields = operationFields(this.editableFields);
    // not op fields, tls fields, or role/backend/scope
    const defaultFields = this.editableFields.filter((f) => ![...opFields, ...tls].includes(f));
    if (defaultFields.length) {
      groups.unshift({ default: defaultFields });
    }
    return this._expandGroups(groups);
  }

  get operationFormFields() {
    const objects = [
      'operationCreate',
      'operationActivate',
      'operationGet',
      'operationLocate',
      'operationRekey',
      'operationRevoke',
      'operationDestroy',
    ];

    const attributes = ['operationAddAttribute', 'operationGetAttributes'];
    const server = ['operationDiscoverVersions'];
    const others = removeManyFromArray(operationFieldsWithoutSpecial(this.editableFields), [
      ...objects,
      ...attributes,
      ...server,
    ]);
    const groups = [
      { 'Managed Cryptographic Objects': objects },
      { 'Object Attributes': attributes },
      { Server: server },
    ];
    if (others.length) {
      groups.push({
        Other: others,
      });
    }
    return this._expandGroups(groups);
  }

  @lazyCapabilities(apiPath`${'backend'}/scope/${'scope'}/role/${'id'}`, 'backend', 'scope', 'id') updatePath;
}
