import Ember from 'ember';

export default class ControlGroupError extends Ember.Error {
  constructor() {
    super();
    this.message = 'Control Group encountered';
  }
}
