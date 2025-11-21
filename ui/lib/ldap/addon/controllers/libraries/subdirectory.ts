/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import type { Breadcrumb } from 'vault/vault/app-types';
import type LdapLibraryModel from 'vault/models/ldap/library';
import type SecretEngineModel from 'vault/models/secret-engine';

interface RouteModel {
  backendModel: SecretEngineModel;
  path_to_library: string;
  libraries: Array<LdapLibraryModel>;
}

export default class LdapLibrariesSubdirectoryController extends Controller {
  @tracked breadcrumbs: Array<Breadcrumb> = [];

  declare model: RouteModel;
}
