export default function sortObjects(array, key) {
  if (Array.isArray(array) && array?.every((e) => e[key] && typeof e[key] === 'string')) {
    return array.sort((a, b) => {
      // ignore upper vs lowercase
      const valueA = a[key].toUpperCase();
      const valueB = b[key].toUpperCase();
      if (valueA < valueB) return -1;
      if (valueA > valueB) return 1;
      return 0;
    });
  }
  // if not sortable, return original array
  return array;
}
