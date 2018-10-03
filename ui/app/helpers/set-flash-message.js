import { inject as service } from '@ember/service';
import Helper from '@ember/component/helper';

export default Helper.extend({
  flashMessages: service(),

  compute([message, type]) {
    return () => {
      this.get('flashMessages')[type || 'success'](message);
    };
  },
});
