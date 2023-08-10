/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export function capabilitiesStub(requestPath, capabilitiesArray) {
  // sample of capabilitiesArray: ['read', 'update']
  return {
    [requestPath]: capabilitiesArray,
    capabilities: capabilitiesArray,
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
