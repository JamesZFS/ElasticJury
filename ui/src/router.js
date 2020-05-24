import Vue from 'vue'
import Router from 'vue-router'
import Home from "./views/Home";
import Detail from "./views/Detail";
import NotFound from "./views/NotFound";

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: '/',
            component: Home
        },
        {
            path: '/detail/:id',
            component: Detail
        },
        {
            path: '*',
            component: NotFound
        }
    ]
})
