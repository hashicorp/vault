/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const statuses = [
  'connected',
  'disconnected since 2022-09-21T11:25:02.196835-07:00; error: unable to establish a connection with HCP',
  'connecting since 2022-09-21T11:25:02.196835-07:00; error: unable to establish a connection with HCP',
  'connecting since 2022-09-21T11:25:02.196835-07:00; error: principal does not have the permission to register as a provider',
  'connecting since 2022-09-21T11:25:02.196835-07:00; error: could not obtain a token with the supplied credentials',
];
let index = null;

export default function (server) {
  const handleResponse = (req, props) => {
    const xhr = req.passthrough();
    xhr.onreadystatechange = () => {
      if (xhr.readyState === 4 && xhr.status < 300) {
        // XMLHttpRequest response prop only has a getter -- redefine as writable and set value
        Object.defineProperty(xhr, 'response', {
          writable: true,
          value: JSON.stringify({
            ...JSON.parse(xhr.responseText),
            ...props,
          }),
        });
      }
    };
  };

  server.get('sys/seal-status', (schema, req) => {
    // return next status from statuses array
    if (index === null || index === statuses.length - 1) {
      index = 0;
    } else {
      index++;
    }
    return handleResponse(req, { hcp_link_status: statuses[index] });
  });
  // enterprise only feature initially
  server.get('sys/health', (schema, req) => handleResponse(req, { version: '1.12.0-dev1+ent' }));
}
