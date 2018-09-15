import EmberError from '@ember/error';

export default class ControlGroupError extends EmberError {
  constructor(wrapInfo) {
    let { accessor, creation_path, creation_time, token, ttl } = wrapInfo;
    super();
    this.message = 'Control Group encountered';

    // add items from the wrapInfo object to the error
    this.token = token;
    this.accessor = accessor;
    this.creation_path = creation_path;
    this.creation_time = creation_time;
    this.ttl = ttl;
  }
}
