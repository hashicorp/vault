import Pretender from 'pretender';
import { noopStub } from './stubs';

/**
 * DEPRECATED prefer to use `setupMirage` along with stubs in vault/tests/helpers/stubs
 */
export default function (options = { usePassthrough: false }) {
  return new Pretender(function () {
    let fn = noopStub();
    if (options.usePassthrough) {
      fn = this.passthrough;
    }
    this.post('/v1/**', fn);
    this.put('/v1/**', fn);
    this.get('/v1/**', fn);
    this.delete('/v1/**', fn || noopStub(204));
  });
}
