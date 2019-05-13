import config from '../config/environment';

export function initialize(/* application */) {
  // attach mount hooks to the environment config
  // context will be the router DSL
  config.addRootMounts = function() {
    this.mount('replication');
  };
}

export default {
  initialize,
};
