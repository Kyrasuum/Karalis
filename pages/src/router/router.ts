import { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
    {
        path: '/:catchAll(.*)*',
        name: "main",
        component: () => import('/src/pages/main.vue'),
        meta: {
            requiresAuth: false,
            level: 0
        }
    },
];

export default routes;
