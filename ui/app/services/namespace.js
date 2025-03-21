/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { getRelativePath } from 'core/utils/sanitize-path';
import { tracked } from '@glimmer/tracking';
import { buildWaiter } from '@ember/test-waiters';

const waiter = buildWaiter('namespaces');
const ROOT_NAMESPACE = '';
export default class NamespaceService extends Service {
  @service store;
  @service auth;

  //populated by the query param on the cluster route
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
    if (!path) {
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
