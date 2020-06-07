import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

// Components
import Home from './components/home/Home';

// Routes
const router = new VueRouter({
  mode: 'history',
  // base: __dirname,
  routes: [
    { name: 'home', path: '/', component: Home }
  ]
});

new Vue({
  router,
  template:`<router-view></router-view>`,
  mounted: async function () {
    console.log('App mounted');
  }
}).$mount('#app');