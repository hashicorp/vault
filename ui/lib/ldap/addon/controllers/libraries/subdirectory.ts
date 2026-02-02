/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { LdapLibrary } from 'vault/secrets/ldap';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface RouteModel {
  secretsEngine: SecretsEngineResource;
  path_to_library: string;
  libraries: Array<LdapLibrary>;
}

export default class LdapLibrariesSubdirectoryController extends Controller {
  @tracked breadcrumbs: Array<Breadcrumb> = [];

  declare model: RouteModel;
}
