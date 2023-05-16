/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

module.exports = [
  // example redirect:
  // {
  //   source: '/vault/docs/some/path',
  //   destination: '/vault/docs/some/other/path',
  //   permanent: true,
  // },
  {
    source: '/vault/docs/plugins/plugin-portal',
    destination: '/vault/integrations',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/autoauth/:slug',
    destination: '/vault/docs/agentandproxy/autoauth/:slug',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/:slug',
    destination: '/vault/docs/agentandproxy/agent/:slug',
    permanent: true,
  }
]
