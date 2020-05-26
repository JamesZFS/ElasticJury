<template>
  <!-- 案件详情页 -->
  <v-container>
    <v-card
            shaped
            :elevation="5"
            class="my-4"
    >
      <v-btn text absolute
             @click="onClickRelating"
             class="ma-4"
             style="right: 0"
             color="primary">
        相关案件
      </v-btn>
      <v-card-title class="justify-center headline">
        案件 {{id}}
      </v-card-title>

      <v-expansion-panels
              v-model="openedPanels"
              hover
              multiple
              flat
              accordion
              class="mb-10"
      >
        <v-expansion-panel>
          <v-expansion-panel-header class="font-weight-bold">
            法官
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-chip v-for="(judge, index) in judges" class="ma-1" :key="index">
              <v-avatar left>
                <v-icon>mdi-account-circle</v-icon>
              </v-avatar>
              {{judge}}
            </v-chip>
          </v-expansion-panel-content>
        </v-expansion-panel>

        <v-expansion-panel>
          <v-expansion-panel-header class="font-weight-bold">
            法条
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-chip label v-for="(law, index) in laws" class="ma-1" :key="index">
              {{law}}
            </v-chip>
          </v-expansion-panel-content>
        </v-expansion-panel>

        <v-expansion-panel>
          <v-expansion-panel-header class="font-weight-bold">
            标签
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-chip label v-for="(tag, index) in tags" class="ma-1" :key="index">
              {{tag}}
            </v-chip>
          </v-expansion-panel-content>
        </v-expansion-panel>

        <v-expansion-panel>
          <v-expansion-panel-header class="font-weight-bold">
            全文内容
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <p>{{detail}}</p>
          </v-expansion-panel-content>
        </v-expansion-panel>

        <v-expansion-panel>
          <v-expansion-panel-header class="font-weight-bold">
            案件结构
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-row>
              <v-col cols="8">
                <!--       search box         -->
                <v-text-field
                        v-model="search"
                        label="搜索..."
                        flat
                        clearable
                        class="mt-n5"
                />
                <v-treeview
                        dense
                        activatable
                        :items="tree"
                        :search="search"
                        :filter="filter"
                        class="scroll-view"
                >
                  <template v-slot:prepend="{item, open}">
                    <v-icon v-if="item.children">
                      {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
                    </v-icon>
                    <v-icon v-else>mdi-circle-medium</v-icon>
                  </template>
                  <template v-slot:label="{item}">
                    <!--    key-value pair      -->
                    <div v-if="item.attributes" style="cursor: pointer" @click="onClickAttribute(item.attributes)">
                      <span class="font-weight-bold">{{item.attributes.key}}</span>
                      <span v-if="item.attributes.value.length > 0"> : {{item.attributes.value.slice(0, 30)}}</span>
                      <span v-if="item.attributes.value.length > 30">...</span>
                    </div>
                  </template>
                </v-treeview>
              </v-col>
              <v-divider vertical/>
              <v-col>
                <h3 class="font-weight-bold mb-2">{{sideView.header}}</h3>
                <p class="scroll-view">{{sideView.content}}</p>
              </v-col>
            </v-row>
          </v-expansion-panel-content>
        </v-expansion-panel>

      </v-expansion-panels>
    </v-card>
  </v-container>
</template>

<script>
    import {getCaseDetail} from "../api";
    import {xmlToTree} from "../utils";

    export default {
        name: "Detail",
        data: () => ({
            openedPanels: [0, 1, 2, 4],
            id: 0,
            judges: [],
            laws: [],
            tags: [],
            detail: '',
            tree: [], // xml tree parsed as js object
            sideView: {
                header: '',
                content: '',
            },
            search: null,
        }),
        async created() {
            // get case id from route path
            this.id = parseInt(this.$route.params.id)
            document.title = `案件 ${this.id}`
            let data = await getCaseDetail(this.id)
            // convert xml tree into js object for display
            data.tree = xmlToTree(data.tree)
            // load case detail & xml tree
            Object.assign(this, data)
        },
        methods: {
            onClickAttribute(attr) {
                this.sideView.header = attr.key
                this.sideView.content = attr.value
            },
            onClickRelating() {
                let routeData = this.$router.resolve(`/?misc=${this.detail.slice(0, 200)}`); // limit 200 chars
                window.open(routeData.href, '_blank');
            },
            filter(item, search) {
                if (item.attributes) {
                    return item.attributes.key.indexOf(search) >= 0 || item.attributes.value.indexOf(search) >= 0
                } else {
                    return false
                }
            },
        }
    }
</script>

<style scoped>
  .scroll-view {
    height: 90vh;
    overflow-y: scroll;
  }
</style>