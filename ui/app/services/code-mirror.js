import Service from '@ember/service';

// This service chiefly exists now for testing purposes.
export default class CodeMirror extends Service {
  _instances = Object.create(null);

  instanceFor(id) {
    return this._instances[id];
  }

  registerInstance(id, instance) {
    this._instances[id] = instance;

    return instance;
  }

  unregisterInstance(id) {
    delete this._instances[id];
  }
}
