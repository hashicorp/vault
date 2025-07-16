/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { baseResourceFactory } from 'vault/resources/base-factory';
import { service } from '@ember/service';
import { supportedTypes } from 'vault/utils/auth-form-helpers';
import engineDisplayData from 'vault/helpers/engines-display-data';

import type { Mount } from 'vault/mount';
import type VersionService from 'vault/services/version';
import type NamespaceService from 'vault/services/namespace';
import type { PathInfo } from 'vault/utils/openapi-helpers';

export default class AuthMethodResource extends baseResourceFactory<Mount>() {
  @service declare readonly version: VersionService;
  @service declare readonly namespace: NamespaceService;

  id: string;
  declare paths: PathInfo;

  constructor(data: Mount, context: unknown) {
    super(data, context);
    // strip trailing slash from path for id since it is used in routing
    this.id = data.path.replace(/\/$/, '');
  }

  // namespaces introduced types with a `ns_` prefix for built-in engines
  // so we need to strip that to normalize the type
  get methodType() {
    return this.type.replace(/^ns_/, '');
  }

  get icon() {
    // methodType refers to the backend type (e.g., "aws", "azure")
    const engineData = engineDisplayData(this.methodType);
    return engineData?.glyph || 'users';
  }

  get directLoginLink() {
    const ns = this.namespace.path;
    const nsQueryParam = ns ? `namespace=${encodeURIComponent(ns)}&` : '';
    const isSupported = supportedTypes(this.version.isEnterprise).includes(this.methodType);
    return isSupported
      ? `${window.origin}/ui/vault/auth?${nsQueryParam}with=${encodeURIComponent(this.path)}`
      : '';
  }

  // used when the `auth` prefix is important,
  // currently only when setting perf mount filtering
  get apiPath() {
    return `auth/${this.path}`;
  }

  get localDisplay() {
    return this.local ? 'local' : 'replicated';
  }

  get supportsUserLockoutConfig() {
    return ['approle', 'ldap', 'userpass'].includes(this.methodType);
  }
}
