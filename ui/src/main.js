import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import owndns from './plugins/owndns'
Vue.use(owndns);
import store from './store'

Vue.config.productionTip = false

new Vue({
  owndns,
  vuetify,
  store,
  render: h => h(App)
}).$mount('#app')
