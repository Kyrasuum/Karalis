import { createWebHistory, createRouter } from 'vue-router';
import routes from './router';
import store from 'src/store';
import { until } from 'src/types/net';

const router = createRouter({
	scrollBehavior: () => ({ left: 0, top: 0 }),
	routes,
	history: createWebHistory(import.meta.env.VUE_ROUTER_BASE),
});

//route guard
router.beforeEach(async (to, from) => {
    if (to.matched.some(record => record.meta.requiresAuth)) {
        await until((_: any) => store.state.user != null && store.state.user.internal_id != null);

        if (!store.state.user || !store.state.user.id || !store.state.user.role || store.state.user.id == '') {
            return '';
        }

        if (to.matched.some(record => store.state.user.role.length != 1 || record.meta.level > store.state.user.role[0].level)) {
            return '';
        }
    }
})

export default router;
