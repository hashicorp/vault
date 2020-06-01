import Component from '@ember/component';
import layout from '../templates/components/replication-header';

/**
 * @module ReplicationHeader
 * The `ReplicationHeader` component is a header component used on the Replication Dashboards.
 * It is a contextual component of the Replication Page component.
 *
 * @example
 * ```js
 * <ReplicationHeader
    @data={{model}}
    @title="Secondary"
    @secondaryID="meep_123"
    @isSummaryDashboard=false
    />
 * ```
 * @param {Object} model=null - An Ember data object pulled form the Ember cluster model.
 * @param {String} title=null - The title of the header.
 * @param {String} [secondaryID=null] - The secondaryID pulled off of the model object. 
 * @param {Boolean} isSummaryDashboard=false - if the Dashboard is for when you have both a primary performance and dr cluster, then this is true.
 */

export default Component.extend({
  layout,
  classNames: ['replication-header'],
  isSecondary: null,
  secondaryId: null,
  isSummaryDashboard: false,
  'data-test-replication-header': true,
});
