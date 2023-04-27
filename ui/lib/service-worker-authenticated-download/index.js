/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

'use strict';

module.exports = {
  name: require('./package').name,

  isDevelopingAddon() {
    return true;
  },

  serverMiddleware({ app }) {
    app.use((req, res, next) => {
      res.setHeader('Service-Worker-Allowed', '/');
      next();
    });
  },
};
