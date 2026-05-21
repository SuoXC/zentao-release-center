import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/projects'
    },
    {
      path: '/projects',
      name: 'Projects',
      component: () => import('../views/Projects.vue')
    },
    {
      path: '/project',
      name: 'ProjectDetail',
      component: () => import('../views/ProjectDetail.vue')
    },
    {
      path: '/release',
      name: 'ReleaseDetail',
      component: () => import('../views/ReleaseDetail.vue')
    }
  ]
})

export default router
