import { createApp } from 'vue';
import App from 'src/app.vue';
import router from 'src/router';
import { Quasar } from 'quasar';
import qusarConfig from './quasar.config';
import 'src/styles/index.scss';
import store from 'src/store';

// Import icon libraries
import '@quasar/extras/material-icons/material-icons.css'
import '@quasar/extras/fontawesome-v6/fontawesome-v6.css'
// Import Quasar css
import 'quasar/src/css/index.sass'

const app = createApp(App);
app.use(store);
app.use(router);
app.use(Quasar, qusarConfig);
app.mount('#app');
