<template>
  <v-container>
    <v-expand-transition>
      <div v-if="displayWelcome">
        <v-img
                src="../assets/logo.svg"
                class="my-3"
                contain
                height="200"
        />
        <v-row justify="center" align="center" style="height: 20vh">
          <h1 class="display-2 font-weight-bold mb-3">
            Welcome to ElasticJury
          </h1>
        </v-row>
      </div>
    </v-expand-transition>

    <div class="mx-auto my-5 search-box">
      <ChipTextInput
              placeholder="全文检索词..."
              icon="mdi-magnify"
              v-model="words.inputs"
              :candidates="words.candidates"
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
              icon="mdi-bookmark-multiple-outline"
              v-model="tags.inputs"
              :candidates="tags.candidates"
      />
    </div>

    <v-row justify="center" class="mb-2">
      <v-radio-group v-model="mode" row class="mt-1 mr-8">
        <v-radio label="And" value="AND"/>
        <v-radio label="Or" value="OR"/>
      </v-radio-group>
      <v-btn @click="onSearch" color="primary" :disabled="!searchAble">Search</v-btn>
      <v-btn @click="onPing" class="ml-10" color="secondary">Ping</v-btn>
    </v-row>

    <v-skeleton-loader v-if="loading" type="table"/>
    <CaseList
            v-else
            :items="result.infos"
            :weights="result.weightsToDisplay"
            @click="onClickCase"
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
      共找到{{resultLength}}条相关结果
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
    import CaseList from "../components/CaseList";
    import {getCaseInfo, ping, searchCaseId} from "../api";

    export default {
        name: 'Home',
        components: {ChipTextInput, CaseList},
        data: () => ({
            displayWelcome: true,
            loading: false,
            curPage: 0,
            casesPerPage: 10,
            notFoundTip: false,
            foundTip: false,
            mode: 'AND',
            words: {
                inputs: [],
                candidates: ['调解', '协议', '当事人'],
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
                weights: [],
                infos: [],
                weightsToDisplay: [],
            }
        }),
        computed: {
            pageCount() {
                return Math.ceil(this.result.ids.length / this.casesPerPage)
            },
            searchAble() {
                return this.words.inputs.length > 0 || this.judges.inputs.length > 0 ||
                    this.laws.inputs.length > 0 || this.tags.inputs.length > 0
            },
            resultLength() {
                return this.result.ids.length
            }
        },
        mounted() {
            // possibly parse url into search param
            let query = this.$route.query
            if (query.hasOwnProperty('word') || query.hasOwnProperty('judge') ||
                query.hasOwnProperty('law') || query.hasOwnProperty('tag')) {
                this.displayWelcome = false;
                // parse param from route and do search
                this.parseParams(query)
                this.doSearch()
            }
        },
        methods: {
            async setPage(page) {
                this.loading = true
                this.curPage = page
                // load results when page changes
                let idsToLoad = this.result.ids.slice((this.curPage - 1) * this.casesPerPage, this.curPage * this.casesPerPage)
                let resp = await getCaseInfo(idsToLoad)
                this.result.weightsToDisplay = this.result.weights.slice((this.curPage - 1) * this.casesPerPage, this.curPage * this.casesPerPage)
                this.result.infos = Object.values(resp)
                this.loading = false
            },
            parseParams(query) {
                this.words.inputs = query.word ? query.word.split(',') : []
                this.judges.inputs = query.judge ? query.judge.split(',') : []
                this.laws.inputs = query.law ? query.law.split(',') : []
                this.tags.inputs = query.tag ? query.tag.split(',') : []
                this.mode = query.mode || 'AND'
            },
            dumpParams() {
                let query = {};
                if (this.words.inputs.length > 0) query.word = this.words.inputs.join(',')
                if (this.judges.inputs.length > 0) query.judge = this.judges.inputs.join(',')
                if (this.laws.inputs.length > 0) query.law = this.laws.inputs.join(',')
                if (this.tags.inputs.length > 0) query.tag = this.tags.inputs.join(',')
                query.mode = this.mode
                return query
            },
            async doSearch() {
                this.loading = true;
                let resp = await searchCaseId(
                    this.words.inputs,
                    this.judges.inputs,
                    this.laws.inputs,
                    this.tags.inputs,
                    this.mode
                )
                if (resp.count === 0) {
                    // no result:
                    this.result.ids = []
                    this.result.infos = []
                    this.notFoundTip = true
                } else {
                    let pairs = Object.entries(resp.result)
                        // eslint-disable-next-line
                        .sort(([_id1, weight1], [_id2, weight2]) => weight2 - weight1) // sort by weight desc
                    this.result.ids = []
                    this.result.weights = []
                    pairs.forEach(([id, weight]) => {
                        this.result.ids.push(parseInt(id))
                        this.result.weights.push(weight)
                    })
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
            onClickCase(index) {
                alert(`Clicked ${index}`);
            },
        }
    }
</script>

<style scoped>
  .search-box {
    max-width: 800px;
  }
</style>
