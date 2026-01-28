/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { getRelativePath, sanitizePath } from 'core/utils/sanitize-path';
import { tracked } from '@glimmer/tracking';
import { buildWaiter } from '@ember/test-waiters';

import type AuthService from 'vault/services/auth';
import type ApiService from 'vault/services/api';
import type FlagsService from 'vault/services/flags';

const waiter = buildWaiter('namespaces');

export const ROOT_NAMESPACE = '';
export const ADMINISTRATIVE_NAMESPACE = 'admin';

export interface NamespaceOption {
  path: string;
  label: string;
}

export default class NamespaceService extends Service {
  @service declare readonly api: ApiService;
  @service declare readonly auth: AuthService;
  @service declare readonly flags: FlagsService;

  // populated by the query param on the cluster route
  @tracked path = '';
  // list of namespaces available to the current user under the
  // current namespace
  @tracked accessibleNamespaces: string[] | null = null;

  get userRootNamespace() {
    // If there is no authData then fallback to relevant root depending on the cluster type
    const fallback = this.flags.isHvdManaged ? ADMINISTRATIVE_NAMESPACE : ROOT_NAMESPACE;
    return this.auth?.authData?.userRootNamespace ?? fallback;
  }

  get inRootNamespace() {
    return this.path === ROOT_NAMESPACE;
  }

  // the top-level namespace is "admin" for HVD managed clusters accessing the UI
  // (similar to "root" for self-managed clusters)
  // this getter checks if the user is specifically at the administrative namespace level
  get inHvdAdminNamespace() {
    return this.flags.isHvdManaged && this.path === ADMINISTRATIVE_NAMESPACE;
  }

  get currentNamespace() {
    if (this.inRootNamespace) return 'root';

    const parts = this.path?.split('/');
    return parts[parts.length - 1];
  }

  get relativeNamespace() {
    // relative namespace is the current namespace minus the user's root.
    // so if we're in app/staging/group1 but the user's root is app, the
    // relative namespace is staging/group
    return getRelativePath(this.path, this.userRootNamespace);
  }

  setNamespace(path: string | undefined) {
    // If a user explicitly logs in to the 'root' namespace, the path is set to 'root'.
    // The root namespace doesn't have a set path, so when verifying the selected namespace, it returns null.
    // Adding a check here, so if the namespace is 'root', it'll be set to an empty string to match the root namespace.
    if (!path || sanitizePath(path) === 'root') {
      this.path = '';
      return;
    }
    this.path = path;
  }

  findNamespacesForUser = task({ drop: true }, async () => {
    const waiterToken = waiter.beginAsync();
    const headers = this.api.buildHeaders({ namespace: this.userRootNamespace });
    try {
      const { keys = [] } = await this.api.sys.internalUiListNamespaces(headers);

      this.accessibleNamespaces = keys.map((n) => {
        let fullNS = n;
        // if the user's root isn't '', then we need to construct
        // the paths so they connect to the user root to the list
        // otherwise using the current ns to grab the correct leaf
        // node in the graph doesn't work
        if (this.userRootNamespace) {
          fullNS = `${this.userRootNamespace}/${n}`;
        }
        return fullNS.replace(/\/$/, '');
      });
    } catch (e) {
      //do nothing here
    } finally {
      waiter.endAsync(waiterToken);
    }
  });

  reset() {
    this.accessibleNamespaces = null;
  }

  getOptions(): NamespaceOption[] {
    /* Each namespace option has 2 properties: { path and label }
     *   - path: full namespace path (used to navigate to the namespace)
     *   - label: text displayed inside the namespace picker dropdown (if root, then path is "", else label = path)
     *
     *  Example:
     *   | path           | label          |
     *   | ----           | -----          |
     *   | ''             | 'root'         |
     *   | 'parent'       | 'parent'       |
     *   | 'parent/child' | 'parent/child' |
     */
    const options = (this.accessibleNamespaces || []).map((ns: string) => ({ path: ns, label: ns }));

    // Add the user's root namespace because `sys/internal/ui/namespaces` does not include it.
    // this.userRootNamespace is guaranteed to be defined due to the fallback in the getter.
    if (!options?.find((o) => o.path === this.userRootNamespace)) {
      // the 'root' namespace is technically an empty string so we manually add the 'root' label.
      const label = this.userRootNamespace === ROOT_NAMESPACE ? 'root' : this.userRootNamespace;
      options.unshift({ path: this.userRootNamespace, label });
    }
    return options;
  }
}
