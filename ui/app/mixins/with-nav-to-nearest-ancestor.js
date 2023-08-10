/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Mixin from '@ember/object/mixin';
import { task } from 'ember-concurrency';
import { ancestorKeysForKey } from 'core/utils/key-utils';

// This mixin is currently used in a controller and a component, but we
// don't see cancellation of the task as the while loop runs in either

// Controller in Ember are singletons so there's no cancellation there
// during the loop. For components, it might be expected that the task would
// be cancelled when we transitioned to a new route and a rerender occured, but this is not
// the case since we are catching the error. Since Ember's route transitions are lazy
// and we're catching any 404s, the loop continues until the transtion succeeds, or exhausts
// the ancestors array and transitions to the root
export default Mixin.create({
  navToNearestAncestor: task(function* (key) {
    const ancestors = ancestorKeysForKey(key);
    let errored = false;
    let nearest = ancestors.pop();
    while (nearest) {
      try {
        const transition = this.transitionToRoute('vault.cluster.secrets.backend.list', nearest);
        transition.data.isDeletion = true;
        yield transition.promise;
      } catch (e) {
        // in the route error event handler, we're only throwing when it's a 404,
        // other errors will be in the route and will not be caught, so the task will complete
        errored = true;
        nearest = ancestors.pop();
      } finally {
        if (!errored) {
          nearest = null;
          // eslint-disable-next-line
          return;
        }
        errored = false;
      }
    }
    yield this.transitionToRoute('vault.cluster.secrets.backend.list-root');
  }),
});
