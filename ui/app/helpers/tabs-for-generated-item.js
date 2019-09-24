import { helper as buildHelper } from '@ember/component/helper';
import { capitalize } from '@ember/string';
import { pluralize } from 'ember-inflector';

export function tabsForGeneratedItem([model, paths]) {
  debugger;
  if (model.paths) {
    paths = model.paths.paths.filter(path => path.navigation);
  }

  console.log(paths);

  let tabs = paths
    .map(path => {
      //we want to show the list of subItems always, but never the list of parentItems
      if (path.itemType === model.paths.itemType && path.operations.includes('list')) {
        return;
      }
      let types = path.itemType.split('~*');
      //we have a parentType (e.g. role~*secret-id)
      //we need to add the parentID (e.g. role~*my-role~*secret-id)
      if (types.length == 2 && model.id) {
        path.itemType = types.join(`~*${model.id}~*`);
      }
      return {
        label: capitalize(pluralize(path.itemName)),
        routeParams: ['vault.cluster.access.method.item.list', path.itemType],
      };
    })
    .compact();

  tabs.unshift({
    label: model.paths ? capitalize(model.paths.itemType) : 'Configuration',
    routeParams: ['vault.cluster.access.method.item.show', model.paths.itemType, model.id],
  });

  return tabs;
}

export default buildHelper(tabsForGeneratedItem);
