/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { baseResourceFactory } from 'vault/resources/base-factory';

interface RoleAwsData {
  name?: string;
  credential_type?: string;
  role_arns?: string[];
  policy_arns?: string[];
  policy_document?: string;
  backend?: string;
  canRead?: boolean;
  canEdit?: boolean;
  canDelete?: boolean;
  canGenerate?: boolean;
}

export default class RoleAwsResource extends baseResourceFactory<RoleAwsData>() {
  declare name?: string;
  declare backend?: string;

  get id() {
    return this.name;
  }
}
