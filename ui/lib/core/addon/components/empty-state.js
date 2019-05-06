import Component from '@ember/component';
import layout from '../templates/components/empty-state';

/**
 * @module EmptyState
 * `EmptyState` components are used to render a helpful message and any necessary content when a user
 * encounters a state that would usually be blank.
 *
 * @example
 * ```js
 * <EmptyState @title="You don't have an secrets yet" @message="An explanation of why you don't have any secrets but also you maybe want to create one." />
 * ```
 *
 * @param title=null{String} - A short label for the empty state
 * @param message=null{String} - A description of why a user might be seeing the empty state and possibly instructions for actions they may take.
 *
 */

export default Component.extend({
  layout,
  tagName: '',
  title: null,
  message: null,
});
