import Vue from 'vue';
import VueRouter from 'vue-router';

// IDE
import VueCodemirror from 'vue-codemirror';
import 'codemirror/lib/codemirror.css';

Vue.use(VueCodemirror);

Vue.use(VueRouter);

// Components
import Home from './components/home/Home';
import Practice from './components/practice/Practice';
import Nomic from './components/nomic/Nomic';

// Internal
import NotFound from './components/system/NotFound';

// Routes
const router = new VueRouter({
  mode: 'history',
  // base: __dirname,
  routes: [
    { name: 'home', path: '/', component: Home },
    { name: 'practice', path: '/practice', component: Practice },
    { name: 'nomic', path: '/nomic', component: Nomic },
    { path: '/404', component: NotFound },
    { path: '*', redirect: '/404' }
  ]
});

new Vue({
  router,
  template:`<router-view></router-view>`,
  mounted: async function () {
    console.log('App mounted');
  }
}).$mount('#app');