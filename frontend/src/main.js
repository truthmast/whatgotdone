import Vue from 'vue'
import axios from 'axios'
import VueAnalytics from 'vue-analytics'
import VueAxios from 'vue-axios'
import VueMoment from 'vue-moment'
import BootstrapVue from 'bootstrap-vue'
import App from './App.vue'
import store from './store.js'
import router from './router.js'

import 'bootswatch/dist/superhero/bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.config.productionTip = false
Vue.use(VueAxios, axios)
Vue.use(BootstrapVue)

Vue.use(VueMoment)
Vue.prototype.moment = VueMoment

if (process.env.VUE_APP_GOOGLE_ANALYTICS_ID.length > 1) {
  Vue.use(VueAnalytics, {
    id: process.env.VUE_APP_GOOGLE_ANALYTICS_ID,
    router
  });
}

// This callback runs before every route change, including on page load.
router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title(to);
  } else {
    document.title = "What Got Done";
  }
  next();
});

new Vue({
  store,
  router,
  render: h => h(App),
}).$mount('#app')