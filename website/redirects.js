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
    destination: '/vault/docs/agent-and-proxy/autoauth',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/autoauth/:slug',
    destination: '/vault/docs/agent-and-proxy/autoauth/:slug',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/template',
    destination: '/vault/docs/agent-and-proxy/agent/template',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/winsvc',
    destination: '/vault/docs/agent-and-proxy/agent/winsvc',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/versions',
    destination: '/vault/docs/agent-and-proxy/agent/versions',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/apiproxy',
    destination: '/vault/docs/agent-and-proxy/agent/apiproxy',
    permanent: true,
  },
  {
    source: '/vault/docs/agent',
    destination: '/vault/docs/agent-and-proxy/agent',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/caching',
    destination: '/vault/docs/agent-and-proxy/agent/caching',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/caching/:slug',
    destination: '/vault/docs/agent-and-proxy/agent/caching/:slug',
    permanent: true,
  },
  {
    source: '/vault/docs/:version(v1\.[4-9]\.x)/agent-and-proxy/:slug',
    destination: '/vault/docs/:version/agent/:slug',
    permanent: true,
  }
]