/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

export function capabilitiesStub(requestPath, capabilitiesArray) {
  // sample of capabilitiesArray: ['read', 'update']
  return {
    request_id: '40f7e44d-af5c-9b60-bd20-df72eb17e294',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      [requestPath]: capabilitiesArray,
      capabilities: capabilitiesArray,
    },
    wrap_info: null,
    warnings: null,
    auth: null,
  };
}

export const noopStub = (response) => {
  return function () {
    return [response, { 'Content-Type': 'application/json' }, JSON.stringify({})];
  };
};

/**
 * allowAllCapabilitiesStub mocks the response from capabilities-self
 * that allows the user to do any action (root user)
 * Example usage assuming setupMirage(hooks) was called:
 * this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub(['read']));
 */
export function allowAllCapabilitiesStub(capabilitiesList = ['root']) {
  return function (_, { requestBody }) {
    const { paths } = JSON.parse(requestBody);
    const specificCapabilities = paths.reduce((obj, path) => {
      return {
        ...obj,
        [path]: capabilitiesList,
      };
    }, {});
    return {
      ...specificCapabilities,
      capabilities: capabilitiesList,
      request_id: 'mirage-795dc9e1-0321-9ac6-71fc',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: { ...specificCapabilities, capabilities: capabilitiesList },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  };
}

/**
 * returns a response with the given httpStatus and data based on status
 * @param {number} httpStatus 403, 404, 204, or 200 (default)
 * @param {object} payload what to return in the response if status is 200
 * @returns {Response}
 */
export function overrideResponse(httpStatus = 200, payload = {}) {
  if (httpStatus === 403) {
    return new Response(
      403,
      { 'Content-Type': 'application/json' },
      JSON.stringify({ errors: ['permission denied'] })
    );
  }
  if (httpStatus === 404) {
    return new Response(404, { 'Content-Type': 'application/json' });
  }
  if (httpStatus === 204) {
    return new Response(204, { 'Content-Type': 'application/json' });
  }
  return new Response(httpStatus, { 'Content-Type': 'application/json' }, payload);
}
