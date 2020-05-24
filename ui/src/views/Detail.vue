<template>
  <!-- 案件详情页 -->
  <v-container>
    <v-app-bar
            app
            color="primary"
            dark
    >
      <span class="headline">案件 {{id}}</span>
      <v-spacer></v-spacer>
      <div class="d-flex align-center">
        <v-img
                alt="Vuetify Logo"
                class="shrink mr-2"
                contain
                src="https://cdn.vuetifyjs.com/images/logos/vuetify-logo-dark.png"
                transition="scale-transition"
                width="40"
        />
        <span class="headline">ElasticJury</span>
      </div>
    </v-app-bar>

    <v-card
            shaped
            class="my-4"
    >
      <v-card-title class="justify-center">
        案件 {{id}}
      </v-card-title>

      <v-divider/>

      <v-card-subtitle>
        <span class="font-weight-bold">法官：</span>
        <v-chip v-for="judge in judges" class="ma-1">
          <v-avatar left>
            <v-icon>mdi-account-circle</v-icon>
          </v-avatar>
          {{judge}}
        </v-chip>
      </v-card-subtitle>

      <v-divider/>

      <v-card-subtitle>
        <span class="font-weight-bold">法条：</span>
        <v-chip label v-for="law in laws" class="ma-1">
          {{law}}
        </v-chip>
      </v-card-subtitle>

      <v-divider/>

      <v-card-subtitle>
        <span class="font-weight-bold">标签：</span>
        <v-chip label v-for="tag in tags" class="ma-1">
          {{tag}}
        </v-chip>
      </v-card-subtitle>

      <v-divider/>

      <v-treeview
              open-on-click
              hoverable
              :items="tree"
              class="my-3"
      >
        <template v-slot:prepend="{item, open}">
          <v-icon v-if="item.children">
            {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
          </v-icon>
          <v-icon v-else>mdi-circle-medium</v-icon>
        </template>
        <template v-slot:label="{item}">
          <!--    key-value pair      -->
          <div v-if="item.attributes">
            <span class="font-weight-bold">{{item.attributes.key}}</span>
            <span v-if="item.attributes.value.trim().length > 0"> : {{item.attributes.value}}</span>
          </div>
        </template>
      </v-treeview>

    </v-card>
  </v-container>
</template>

<script>
    import {getCaseDetail} from "../api";
    import convert from "xml-js";

    export default {
        name: "Detail",
        data: () => ({
            id: 0,
            judges: [],
            laws: [],
            tags: [],
            detail: '',
            tree: [], // xml tree parsed as js object
        }),
        async created() {
            // get case id from route path
            this.id = parseInt(this.$route.params.id)
            // load case detail & xml tree
            Object.assign(this, await getCaseDetail(this.id))
            // convert xml tree into js object for display
            this.tree = convert.xml2js(this.tree, {
                compact: false,
                ignoreComment: true,
                elementsKey: 'children',
            }).children[0].children
        }
    }
</script>

<style scoped>

</style>