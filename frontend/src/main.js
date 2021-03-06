import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';

import VueForm from 'vue-form';
import vuetify from '@/plugins/vuetify';
import '@/plugins/vue-session';
import '@/plugins/vue-notification';

Vue.config.productionTip = false;

// Enable vue-form
Vue.use(VueForm);

import firebase from 'firebase/app';
import 'firebase/auth'; // for authentication
import 'firebase/analytics'; // for authentication

// Fetch and load the store settings
store.dispatch('getSettings').then(() => {
  // Initialize Firebase with a "default" Firebase project
  const firebaseConfig = {
    apiKey: store.state.settings.apiKey,
    authDomain: store.state.settings.authDomain,
    projectId: store.state.settings.projectId,
    storageBucket: store.state.settings.storageBucket,
    appId: store.state.settings.appId,
    measurementId: store.state.settings.measurementId
  };

  const defaultProject = firebase.initializeApp(firebaseConfig);
  firebase.analytics();

  firebase.auth().useDeviceLanguage();
  firebase.auth().onAuthStateChanged(user => {
    store.dispatch('fetchUser', user);
    if (user) {
      console.debug('User is signed in');
    } else {
      console.debug('No user is signed in');
    }

    new Vue({
      vuetify,
      router,
      store,
      render: h => h(App)
    }).$mount('#app');
  });
});
