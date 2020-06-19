import Pretender from 'pretender';

const noop = response => {
  return function() {
    return [response, { 'Content-Type': 'application/json' }, JSON.stringify({})];
  };
};

/*
This helper uses Pretender to mock a server. It intercepts all
requests so that we don't make a real, external call to the Vault
API.
*/
export default function(options = { usePassthrough: false }) {
  return new Pretender(function() {
    let fn = noop();
    if (options.usePassthrough) {
      fn = this.passthrough;
    }
    this.post('/v1/**', fn);
    this.put('/v1/**', fn);
    this.get('/v1/**', fn);
    this.delete('/v1/**', fn || noop(204));
  });
}
