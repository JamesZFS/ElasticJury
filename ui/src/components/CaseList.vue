<template>
  <v-list three-line>
    <template v-for="(item, index) in items">
      <v-btn
              text
              class="font-weight-bold"
              @click="$emit('click', index)"
      >
        {{'案例' + item.id}}
      </v-btn>

      <v-divider></v-divider>

      <v-list-item>
        <v-list-item-avatar>
          <v-img src="../assets/logo.svg" alt="logo"></v-img>
        </v-list-item-avatar>

        <v-list-item-content style="font-size: small">
          <div>
            <div class="my-2">
              <span class="font-weight-bold">法官：</span>
              {{item.judges}}
            </div>
            <div class="my-2">
              <span class="font-weight-bold">法条：</span>
              {{item.laws}}
            </div>
            <div class="my-2">
              <span class="font-weight-bold">标签：</span>
              {{item.tags}}
            </div>
            <div class="my-2">
              <span class="font-weight-bold">关键词：</span>
              {{item.keywords}}
            </div>
            <div class="my-2">
              <span class="font-weight-bold">细节：</span>
              <v-list-item-subtitle>{{item.detail}}</v-list-item-subtitle>
            </div>
          </div>
        </v-list-item-content>

      </v-list-item>
    </template>
  </v-list>
</template>

<script>
    export default {
        name: "CaseList",
        props: {items: Array},
        data: () => ({}),
        methods: {}
    }
</script>

<style scoped>
</style>