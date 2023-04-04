/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

'use strict';
const fs = require('fs');
module.exports = {
  name: require('./package').name,

  isDevelopingAddon() {
    return true;
  },

  postBuild(result) {
    // We gitignore the contents of our output directory http/web_ui
    // but we need to keep the folder structure for the Vault build
    fs.writeFileSync(result.directory + '/.gitkeep', '');
  },
};
