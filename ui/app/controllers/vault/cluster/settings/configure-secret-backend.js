import { isPresent } from '@ember/utils';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

const CONFIG_ATTRS = {
  // ssh
  configured: false,

  // aws root config
  iamEndpoint: null,
  stsEndpoint: null,
  accessKey: null,
  secretKey: null,
  region: '',
};

export default Controller.extend(CONFIG_ATTRS, {
  queryParams: ['tab'],
  tab: '',
  flashMessages: service(),
  loading: false,
  reset() {
    this.get('model').rollbackAttributes();
    this.setProperties(CONFIG_ATTRS);
  },
  actions: {
    saveConfig(options = { delete: false }) {
      const isDelete = options.delete;
      if (this.get('model.type') === 'ssh') {
        this.set('loading', true);
        this.get('model')
          .saveCA({ isDelete })
          .then(() => {
            this.set('loading', false);
            this.send('refreshRoute');
            this.set('configured', !isDelete);
            if (isDelete) {
              this.get('flashMessages').success('SSH Certificate Authority Configuration deleted!');
            } else {
              this.get('flashMessages').success('SSH Certificate Authority Configuration saved!');
            }
          });
      }
    },

    save(method, data) {
      this.set('loading', true);
      const hasData = Object.keys(data).some(key => {
        return isPresent(data[key]);
      });
      if (!hasData) {
        return;
      }
      this.get('model')
        .save({
          adapterOptions: {
            adapterMethod: method,
            data,
          },
        })
        .then(() => {
          this.get('model').send('pushedData');
          this.reset();
          this.get('flashMessages').success('The backend configuration saved successfully!');
        })
        .finally(() => {
          this.set('loading', false);
        });
    },
  },
});
