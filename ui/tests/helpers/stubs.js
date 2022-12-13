export function capabilitiesStub(requestPath, capabilitiesArray) {
  // sample of capabilitiesArray: ['read', 'update']
  return {
    request_id: '40f7e44d-af5c-9b60-bd20-df72eb17e294',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      capabilities: capabilitiesArray,
      [requestPath]: capabilitiesArray,
    },
    wrap_info: null,
    warnings: null,
    auth: null,
  };
}

/**
 * allowAllCapabilitiesStub mocks the response from capabilities-self
 * that allows the user to do any action (root user)
 * EXAMPLE USAGE:
 * this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub);
 */
export function allowAllCapabilitiesStub() {
  return {
    request_id: '40f7e44d-af5c-9b60-bd20-df72eb17e294',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      capabilities: ['root'],
    },
    wrap_info: null,
    warnings: null,
    auth: null,
  };
}
