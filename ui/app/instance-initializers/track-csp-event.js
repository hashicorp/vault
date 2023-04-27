/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export function initialize(appInstance) {
  const service = appInstance.lookup('service:csp-event');
  service.attach();
}

export default {
  name: 'track-csp-event',
  initialize,
};
