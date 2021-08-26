import { helper as buildHelper } from '@ember/component/helper';

export function splitObject(originalObject, array) {
  let object1 = {};
  let object2 = {};
  // convert object to key's array
  let keys = Object.keys(originalObject);
  // iterate over keys to see if they match values in the array
  keys.forEach(key => {
    array.forEach(item => {
      if (key === item) {
        object1[key] = originalObject[key];
      }
    });
    if (!array.includes(key)) {
      object2[key] = originalObject[key];
    }
  });
  return [object1, object2];
}

export default buildHelper(splitObject);
