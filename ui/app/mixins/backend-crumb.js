import { computed } from '@ember/object';
import Mixin from '@ember/object/mixin';

export default Mixin.create({
  backendCrumb: computed('backend', function() {
    const backend = this.get('backend');

    if (backend === undefined) {
      throw new Error('backend-crumb mixin requires backend to be set');
    }

    return {
      label: backend,
      text: backend,
      path: 'vault.cluster.secrets.backend.list-root',
      model: backend,
    };
  }),
});
