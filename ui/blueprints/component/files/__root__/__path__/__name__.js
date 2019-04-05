/**
 * @module <%= classifiedModuleName %>
 * <%= classifiedModuleName %> components are used to...
 * 
 * @example
 * ```js
 * <<%= classifiedModuleName %> @param1={param1} @param2={param2} />
 * ```
 *
 * @param param1 {String} - param1 is...
 * @param [param2=value] {String} - param2 is... //brackets mean it is optional and = sets the default value
 */
import Component from '@ember/component';
<%= importTemplate %>
export default Component.extend({<%= contents %>
});
