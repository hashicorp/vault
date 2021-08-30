import { helper as buildHelper } from '@ember/component/helper';

export function splitObject(originalObject, array) {
  let object1 = {};
  let object2 = {};
  // convert object to key's array
  let keys = Object.keys(originalObject);
  keys.forEach(key => {
    if (array.includes(key)) {
      object1[key] = originalObject[key];
    } else {
      object2[key] = originalObject[key];
    }
  });
  return [object1, object2];
}

export default buildHelper(splitObject);
