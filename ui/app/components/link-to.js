import LinkComponent from '@ember/routing/link-component';

LinkComponent.reopen({
  activeClass: 'is-active',
});

export default LinkComponent;
