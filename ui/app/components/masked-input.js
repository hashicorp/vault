import Ember from 'ember';
const { computed } = Ember;

export default Ember.Component.extend({
	tagName: "",
	value: null,
	shouldObscure: computed("isMasked", "isFocused", function(){
		if(this.get('isFocused') === true){
			return false;
		}
		return this.get('isMasked');
	}),
	isMasked: true,
	isFocused: false,
	onKeyDown(){},
	onChange(){},
	actions: {
		toggleMask(){
			this.toggleProperty('isMasked');
		},
		setFocus(isFocused){
			this.set('isFocused', isFocused);
		}
	}
});
