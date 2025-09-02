/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Helper from '@ember/component/helper';

export default Helper.extend({
  router: service(),

  compute([routeName, ...models], { replace = false }) {
    return () => {
      const router = this.router;
      const method = replace ? router.replaceWith : router.transitionTo;
      return method.call(router, routeName, ...models);
    };
  },
});
