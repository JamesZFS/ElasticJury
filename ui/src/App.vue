<template>
  <v-app>
    <v-content>
      <router-view/>
    </v-content>
    <v-footer padless color="grey lighten-2" class="mt-10">
      <v-col class="text-center" cols="12">
        <div>
          {{ new Date().getFullYear() }} — <a href="/"><strong>ElasticJury</strong></a>
        </div>
        <div>
          Powered by <a href="mailto:zhengfs16@mails.tsinghua.edu.cn">Fengshi Zheng</a> and <a href="mailto:zhaocg17@mails.tsinghua.edu.cn">Chenggang Zhao</a>, THUCST
        </div>
        <div>
          Icons made by <a href="https://www.flaticon.com/authors/freepik" title="Freepik">Freepik</a> from
          <a href="https://www.flaticon.com/" title="Flaticon"> www.flaticon.com</a>
        </div>
      </v-col>
    </v-footer>
  </v-app>
</template>

<script>
    export default {
        name: 'App',

        data: () => ({
            //
        }),
    };
</script>
