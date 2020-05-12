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
            :items="result.info"
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
                info: [],
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
        methods: {
            async setPage(page) {
                this.curPage = page;
                // load results when page changes
                let idsToLoad = this.result.ids.slice((this.curPage - 1) * this.casesPerPage, this.curPage * this.casesPerPage);
                let resp = await getCaseInfo(idsToLoad);
                this.result.info = Object.values(resp);
            },
            async onSearch() {
                this.displayWelcome = false;
                this.loading = true;
                let resp = await searchCaseId(
                    this.words.inputs,
                    this.judges.inputs,
                    this.laws.inputs,
                    this.tags.inputs,
                    this.mode
                );
                if (resp.count === 0) {
                    // no result:
                    this.result.ids = []
                    this.result.info = []
                    this.notFoundTip = true
                } else {
                    this.result.ids = Object.entries(resp.result)
                        // .sort(([_id1, val1], [_id2, val2]) => val2 - val1)
                        .map(([id]) => parseInt(id));
                    await this.setPage(1)
                    this.foundTip = true
                }
                this.loading = false;
            },
            async onPing() {
                let resp = await ping();
                alert(resp);
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
