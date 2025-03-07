import {
	ActionTree,
	createStore,
	MutationTree,
	Store as VuexStore,
} from 'vuex';
import { state, IAppState } from './state';

// Import from libary mutators and actions

// atomic mutations
const mutations: MutationTree<IAppState> = {
};

// list of callable actions
// call mutator: await commit('MUTATOR', value);
const actions: ActionTree<IAppState, IAppState> = {
};

const store = createStore<IAppState>({
	state: state,
	mutations: mutations,
	actions: actions,
});
export default store;

// Allow usage of this.$store in components
declare module '@vue/runtime-core' {
    interface ComponentCustomProperties  {
        $store: VuexStore
    }
}
