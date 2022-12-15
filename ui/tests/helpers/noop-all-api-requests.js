import Pretender from 'pretender';

const noop = (response) => {
  return function () {
    return [response, { 'Content-Type': 'application/json' }, JSON.stringify({})];
  };
};

function allowCapabilities() {
  return {
    request_id: '40f7e44d-af5c-9b60-bd20-df72eb17e294',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      // capabilities: ['root'],
      capabilities: ['create', 'read', 'update', 'delete', 'list', 'sudo'],
    },
    wrap_info: null,
    warnings: null,
    auth: null,
  };
}

export default function (options = { usePassthrough: false }) {
  return new Pretender(function () {
    let fn = noop();
    if (options.usePassthrough) {
      fn = this.passthrough;
    }
    this.post('/v1/capabilities-self', options.usePassthrough ? fn : allowCapabilities);
    this.post('/v1/**', fn);
    this.put('/v1/**', fn);
    this.get('/v1/**', fn);
    this.delete('/v1/**', fn || noop(204));
  });
}
