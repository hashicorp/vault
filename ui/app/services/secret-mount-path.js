import Service from '@ember/service';

// this service tracks the path of the currently viewed secret mount
// so that we can access that inside of engines where parent route params
// are not accessible
export default Service.extend({
  currentPath: null,
  update(path) {
    this.set('currentPath', path);
  },
  get() {
    return this.currentPath;
  },
});
