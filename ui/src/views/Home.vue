<template>
  <!-- 主页/搜索/结果页 -->
  <v-container>
    <v-expand-transition>
      <div v-if="displayWelcome">
        <v-img
                src="../assets/logo.png"
                class="mt-5"
                contain
                height="250"
        />
        <v-row justify="center" align="center" class="text-center" style="height: 20vh">
          <h1 class="display-2 font-weight-bold">
            ElasticJury
          </h1>
        </v-row>
      </div>
    </v-expand-transition>

    <div class="mx-auto my-5 search-box">
      <AutocompleteInput
              placeholder="综合搜索..."
              icon="mdi-search-web"
              v-model="misc.input"
              :candidates="misc.candidates"
              :on-associate="onAssociateMisc"
      />
      <ChipTextInput
              placeholder="法官名..."
              icon="mdi-account-multiple"
              v-model="judges.inputs"
              :candidates="judges.candidates"
      />
      <ChipTextInput
              placeholder="法条..."
              icon="mdi-book-open-page-variant"
              v-model="laws.inputs"
              :candidates="laws.candidates"
      />
      <ChipTextInput
              placeholder="标签..."
              icon="mdi-tag-multiple-outline"
              v-model="tags.inputs"
              :candidates="tags.candidates"
      />
    </div>

    <v-row justify="center" class="mb-2">
      <v-btn @click="onSearch" color="primary" :disabled="!searchAble">Search</v-btn>
      <v-btn @click="onPing" class="ml-10" color="secondary">Ping</v-btn>
    </v-row>

    <v-skeleton-loader v-if="loading" type="table"/>
    <CaseInfoList
            v-else
            :items="result.infos"
            @click-case="onClickCase"
            @click-judge="judge => pushInput(judges.inputs, judge)"
            @click-law="law => pushInput(laws.inputs, law)"
            @click-tag="tag => pushInput(tags.inputs, tag)"
    />
    <v-pagination
            v-if="resultLength > 0"
            v-model="curPage"
            :total-visible="10"
            :length="pageCount"
            @input="setPage"
            circle
            class="my-5"
    />

    <v-snackbar
            v-model="notFoundTip"
            color="error"
            :timeout="3000"
    >
      没有找到相关结果
      <v-btn
              dark
              text
              @click="notFoundTip = false"
      >
        Close
      </v-btn>
    </v-snackbar>

    <v-snackbar
            v-model="foundTip"
            color="success"
            :timeout="3000"
    >
      共找到 {{resultLength}} 条相关结果
      <v-btn
              dark
              text
              @click="foundTip = false"
      >
        Close
      </v-btn>
    </v-snackbar>

  </v-container>
</template>

<script>
    import ChipTextInput from "../components/ChipTextInput";
    import CaseInfoList from "../components/CaseInfoList";
    import AutocompleteInput from "../components/AutocompleteInput";
    import {getCaseInfo, ping, searchCaseId} from "../api";
    import {sleep} from "../utils";

    export default {
        name: 'Home',
        components: {AutocompleteInput, ChipTextInput, CaseInfoList},
        data: () => ({
            displayWelcome: true,
            loading: false,
            curPage: 0,
            casesPerPage: 10,
            notFoundTip: false,
            foundTip: false,
            misc: {
                input: '',
                candidates: ['调解', '公司', '当事人'],
            },
            judges: {
                inputs: [],
                candidates: ['黄琴英', '高原', '张成镇'],
            },
            laws: {
                inputs: [],
                candidates: ['《中华人民共和国民法通则》', '《中华人民共和国民事诉讼法》', '《中华人民共和国担保法》'],
            },
            tags: {
                inputs: [],
                candidates: ['民事案件', '一审案件'],
            },
            result: {
                ids: [],
                infos: [],
            }
        }),
        computed: {
            pageCount() {
                return Math.ceil(this.result.ids.length / this.casesPerPage)
            },
            searchAble() {
                return this.judges.inputs.length > 0 || this.laws.inputs.length > 0 ||
                    this.tags.inputs.length > 0 || this.misc.input
            },
            resultLength() {
                return this.result.ids.length
            }
        },
        created() {
            // possibly parse url into search param
            let query = this.$route.query
            if (query.hasOwnProperty('misc') || query.hasOwnProperty('judge') ||
                query.hasOwnProperty('law') || query.hasOwnProperty('tag')) {
                this.displayWelcome = false;
                // parse param from route and do search
                this.parseParams(query)
                this.doSearch()
            }
        },
        methods: {
            async setPage(to) {
                this.loading = true
                this.curPage = to
                // load results when page changes
                let idsToLoad = this.result.ids.slice((to - 1) * this.casesPerPage, to * this.casesPerPage)
                this.result.infos = await getCaseInfo(idsToLoad)
                this.loading = false
            },
            parseParams(query) {
                this.misc.input = query.misc || ''
                this.judges.inputs = query.judge ? query.judge.split(',') : []
                this.laws.inputs = query.law ? query.law.split(',') : []
                this.tags.inputs = query.tag ? query.tag.split(',') : []
            },
            dumpParams() {
                let query = {};
                if (this.misc.input) query.misc = this.misc.input.slice(0, 200) // limit 200 chars
                if (this.judges.inputs.length > 0) query.judge = this.judges.inputs.join(',')
                if (this.laws.inputs.length > 0) query.law = this.laws.inputs.join(',')
                if (this.tags.inputs.length > 0) query.tag = this.tags.inputs.join(',')
                return query
            },
            pushInput(inputs, item) {
                if (!inputs.includes(item)) {
                    inputs.push(item)
                }
            },
            async doSearch() {
                this.loading = true;
                let resp = await searchCaseId(
                    this.misc.input,
                    this.judges.inputs,
                    this.laws.inputs,
                    this.tags.inputs,
                )
                if (resp.length === 0) {
                    // no result:
                    this.result.ids = []
                    this.result.infos = []
                    this.notFoundTip = true
                } else {
                    this.result.ids = resp
                    await this.setPage(1)
                    this.foundTip = true
                }
                this.loading = false
            },
            onSearch() {
                this.displayWelcome = false
                this.$router.push({query: this.dumpParams()})
                this.doSearch()
            },
            async onPing() {
                let resp = await ping()
                alert(resp)
            },
            onClickCase(id) {
                let routeData = this.$router.resolve(`detail/${id}`);
                window.open(routeData.href, '_blank');
            },
            async onAssociateMisc(word) {
                console.log('associating:', word)
                await sleep(500)
                return [word, word + '123', word + '456']
            }
        }
    }
</script>

<style scoped>
  .search-box {
    max-width: 800px;
  }
</style>
