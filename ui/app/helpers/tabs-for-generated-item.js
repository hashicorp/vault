import { helper as buildHelper } from '@ember/component/helper';
import { capitalize } from '@ember/string';
import { pluralize } from 'ember-inflector';

export function tabsForGeneratedItem([model, paths]) {
  if (model.paths) {
    paths = model.paths.paths.filter(path => path.navigation);
  }

  console.log(paths);

  let tabs = paths
    .map(path => {
      if (path.itemType === model.paths.itemType && path.operations.includes('list')) {
        return;
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
