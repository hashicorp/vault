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
    source: '/vault/docs/agent/autoauth',
    destination: '/vault/docs/agentandproxy/autoauth',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/autoauth/:slug',
    destination: '/vault/docs/agentandproxy/autoauth/:slug',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/template',
    destination: '/vault/docs/agentandproxy/agent/template',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/winsvc',
    destination: '/vault/docs/agentandproxy/agent/winsvc',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/versions',
    destination: '/vault/docs/agentandproxy/agent/versions',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/apiproxy',
    destination: '/vault/docs/agentandproxy/agent/apiproxy',
    permanent: true,
  },
  {
    source: '/vault/docs/agent',
    destination: '/vault/docs/agentandproxy/agent',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/caching',
    destination: '/vault/docs/agentandproxy/agent/caching',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/caching/:slug',
    destination: '/vault/docs/agentandproxy/agent/caching/:slug',
    permanent: true,
  }
]
