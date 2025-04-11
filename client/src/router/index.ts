import {createRouter, createWebHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AuthLayout from "@/views/Auth/Layout.vue";
import Login from "@/views/Auth/Login.vue";
import Register from "@/views/Auth/Register.vue";
import GameLayout from "@/views/Game/Layout.vue";
import Index from "@/views/Game/IndexView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/auth',
      name: 'auth',
      component: AuthLayout,
      children: [
        {
          path: 'login',
          name: 'auth.login',
          component: Login,
        },
        {
          path: 'register',
          name: 'auth.register',
          component: Register,
        }
      ]
    },
    {
      path: "/game",
      name: "game",
      component: GameLayout,
      children: [
        {
          path: "",
          name: "game.index",
          component: Index,
        }
      ]
    }
  ],
})

export default router
