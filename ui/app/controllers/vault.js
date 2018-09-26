import Controller from '@ember/controller';

export default Controller.extend({
  queryParams: [
    {
      wrappedToken: 'wrapped_token',
    },
  ],
  wrappedToken: '',
});
