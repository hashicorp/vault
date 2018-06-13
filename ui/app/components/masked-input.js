import Ember from 'ember';
const { computed } = Ember;
import autosize from 'autosize';


export default Ember.Component.extend({
	value: null,
	didInsertElement(){
		this._super(...arguments);
		autosize(this.element.querySelector('textarea'));
	},
	didRender(){
		this._super(...arguments);
		autosize.update(this.element.querySelector('textarea'));
	},
	willDestroyElement(){
		this._super(...arguments);
		autosize.destroy(this.element.querySelector('textarea'));
	},
	shouldObscure: computed("isMasked", "isFocused", function(){
		if(this.get('isFocused') === true){
			return false;
		}
		return this.get('isMasked');
	}),
	displayValue: computed("shouldObscure", function(){
		if(this.get("shouldObscure")){
			return "■ ■ ■ ■ ■ ■";
		}
		else{
			return this.get('value');
		}
	}),
	isMasked: true,
	isFocused: false,
	onKeyDown(){},
	onChange(){},
	updateValue(e){
		this.set('value', e.target.value);
		this.get('onChange')();
	},
	actions: {
		toggleMask(){
			this.toggleProperty('isMasked');
		},
		setFocus(isFocused){
			this.set('isFocused', isFocused);
		}
	}
});
