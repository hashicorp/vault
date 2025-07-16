/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { getRelativePath, sanitizePath } from 'core/utils/sanitize-path';
import { tracked } from '@glimmer/tracking';
import { buildWaiter } from '@ember/test-waiters';

const waiter = buildWaiter('namespaces');
const ROOT_NAMESPACE = '';
export default class NamespaceService extends Service {
  @service auth;
  @service flags;
  @service store;

  // populated by the query param on the cluster route
  @tracked path = '';
  // list of namespaces available to the current user under the
  // current namespace
  @tracked accessibleNamespaces = null;

  get userRootNamespace() {
    return this.auth.authData?.userRootNamespace;
  }

  get inRootNamespace() {
    return this.path === ROOT_NAMESPACE;
  }

  // the top-level namespace is "admin" for HVD managed clusters accessing the UI
  // (similar to "root" for self-managed clusters)
  // this getter checks if the user is specifically at the administrative namespace level
  get inHvdAdminNamespace() {
    return this.flags.isHvdManaged && this.path === 'admin';
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

  setNamespace(path) {
    // If a user explicitly logs in to the 'root' namespace, the path is set to 'root'.
    // The root namespace doesn't have a set path, so when verifying the selected namespace, it returns null.
    // Adding a check here, so if the namespace is 'root', it'll be set to an empty string to match the root namespace.
    if (!path || sanitizePath(path) === 'root') {
      this.path = '';
      return;
    }
    this.path = path;
  }

  @task({ drop: true })
  *findNamespacesForUser() {
    const waiterToken = waiter.beginAsync();
    // uses the adapter and the raw response here since
    // models get wiped when switching namespaces and we
    // want to keep track of these separately
    const store = this.store;
    const adapter = store.adapterFor('namespace');
    const userRoot = this.auth.authData.userRootNamespace;
    try {
      const ns = yield adapter.findAll(store, 'namespace', null, {
        adapterOptions: {
          forUser: true,
          namespace: userRoot,
        },
      });
      const keys = ns.data.keys || [];
      this.accessibleNamespaces = keys.map((n) => {
        let fullNS = n;
        // if the user's root isn't '', then we need to construct
        // the paths so they connect to the user root to the list
        // otherwise using the current ns to grab the correct leaf
        // node in the graph doesn't work
        if (userRoot) {
          fullNS = `${userRoot}/${n}`;
        }
        return fullNS.replace(/\/$/, '');
      });
    } catch (e) {
      //do nothing here
    } finally {
      waiter.endAsync(waiterToken);
    }
  }

  reset() {
    this.accessibleNamespaces = null;
  }
}
