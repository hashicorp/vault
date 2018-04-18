import Adapter from './pki';

export default Adapter.extend({
  url(_, snapshot) {
    const backend = snapshot.attr('backend');
    return `/v1/${backend}/root/sign-intermediate`;
  },
});
