<template>
  <v-list two-line>
    <template v-for="(item, index) in items">
      <div :key="index">
        <v-btn
                text
                x-large
                class="font-weight-bold"
                @click="$emit('click-case', item.id)"
        >
          案件 {{item.id}}
        </v-btn>

        <v-divider></v-divider>

        <v-list-item>
          <v-list-item-content style="font-size: small">
            <div>

              <div class="my-2">
                <span class="font-weight-bold">法官：</span>
                <v-chip small v-for="(judge, index) in item.judges.slice(0, 10)" class="ma-1" :key="index"
                        @click="$emit('click-judge', judge)">
                  <v-avatar left>
                    <v-icon>mdi-account-circle</v-icon>
                  </v-avatar>
                  {{judge}}
                </v-chip>
                <span v-if="item.judges.length > 10">...</span>
              </div>

              <div class="my-2">
                <span class="font-weight-bold">法条：</span>
                <v-chip label small v-for="(law, index) in item.laws.slice(0, 3)" class="ma-1" :key="index"
                        @click="$emit('click-law', law)">
                  {{law}}
                </v-chip>
                <span v-if="item.laws.length > 3">...</span>
              </div>

              <div class="my-2">
                <span class="font-weight-bold">标签：</span>
                <v-chip label small v-for="(tag, index) in item.tags.slice(0, 8)" class="ma-1" :key="index"
                        @click="$emit('click-tag', tag)">
                  {{tag}}
                </v-chip>
                <span v-if="item.tags.length > 8">...</span>
              </div>

              <div class="my-2">
                <span class="font-weight-bold">细节：</span>
                <span class="grey--text text--darken-2" style="line-height: normal">
                  {{item.detail.slice(0, 300)}}
                </span>
                <span v-if="item.detail.length > 300">...</span>
              </div>

            </div>
          </v-list-item-content>

        </v-list-item>
      </div>
    </template>
  </v-list>
</template>

<script>
    export default {
        props: {
            items: Array,
        },
        data: () => ({}),
        methods: {}
    }
</script>

<style scoped>
</style>