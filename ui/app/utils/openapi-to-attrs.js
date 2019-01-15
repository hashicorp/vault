import DS from 'ember-data';
const { attr } = DS;

export const expandOpenApiProps = function(props) {
  let attrs = {};
  // expand all attributes
  for (let prop in props) {
    let details = props[prop];
    let editType = details.type;
    if (details.format === 'seconds') {
      editType = 'ttl';
    } else if (details.items) {
      editType = details.items.type + details.type.capitalize();
    }
    attrs[prop.camelize()] = attr({
      editType: editType || details.type,
    });
  }
  return attrs;
};
