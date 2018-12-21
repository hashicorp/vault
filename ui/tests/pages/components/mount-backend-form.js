import { clickable, collection, fillable, text, value } from 'ember-cli-page-object';
import fields from './form-field';
import errorText from './alert-banner';

export default {
  ...fields,
  ...errorText,
  header: text('[data-test-mount-form-header]'),
  submit: clickable('[data-test-mount-submit]'),
  next: clickable('[data-test-mount-next]'),
  back: clickable('[data-test-mount-back]'),
  path: fillable('[data-test-input="path"]'),
  toggleOptions: clickable('[data-test-toggle-group="Method Options"]'),
  pathValue: value('[data-test-input="path"]'),
  types: collection('[data-test-mount-type-radio] input', {
    select: clickable(),
    mountType: value(),
  }),
  type: fillable('[name="mount-type"]'),
  async selectType(type) {
    return this.types.filterBy('mountType', type)[0].select();
  },
  async mount(type, path) {
    await this.selectType(type);
    if (path) {
      return this.next()
        .path(path)
        .submit();
    }
    return this.next().submit();
  },
};
