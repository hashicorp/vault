import Model from '@ember-data/model';

export default class GeneratedItemModel extends Model {
  allFields = [];

  get fieldGroups() {
    const groups = {
      default: [],
    };
    const fieldGroups = [];
    this.constructor.eachAttribute((name, attr) => {
      // if the attr comes in with a fieldGroup from OpenAPI,
      if (attr.options.fieldGroup) {
        if (groups[attr.options.fieldGroup]) {
          groups[attr.options.fieldGroup].push(attr);
        } else {
          groups[attr.options.fieldGroup] = [attr];
        }
      } else {
        // otherwise just add that attr to the default group
        groups.default.push(attr);
      }
    });
    for (const group in groups) {
      fieldGroups.push({ [group]: groups[group] });
    }
    return fieldGroups;
  }
}
