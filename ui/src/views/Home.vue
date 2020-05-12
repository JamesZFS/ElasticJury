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

    <div class="mx-auto my-2 search-box">
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
      <v-btn @click="onSearch" class="mr-10">Search</v-btn>
      <v-btn @click="onInfo" class="mr-10">Info</v-btn>
      <v-btn @click="onPing">Ping</v-btn>
    </v-row>

    <v-fade-transition>
      <v-skeleton-loader v-if="loading" type="table"/>
      <CaseList
              v-else
              :items="result.info"
              @click="onClickCase"
      />
    </v-fade-transition>

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
        methods: {
            async onSearch() {
                this.displayWelcome = false;
                this.loading = true;
                console.log(this.words.inputs);
                let resp = await searchCaseId(
                    this.words.inputs,
                    this.judges.inputs,
                    this.laws.inputs,
                    this.tags.inputs,
                );
                // alert(JSON.stringify(resp));
                this.result.ids = Object.entries(resp.result)
                    // .sort(([_id1, val1], [_id2, val2]) => val2 - val1)
                    .map(([id, _val]) => parseInt(id));
                await this.onInfo();
                this.loading = false;
            },
            async onInfo() {
                console.log(this.result.ids);
                let resp = await getCaseInfo(this.result.ids);
                this.result.info = Object.values(resp);
                console.log(resp);
            },
            async onPing() {
                let resp = await ping();
                alert(resp);
            },
            onClickCase(index) {
                alert(`Clicked ${index}`);
            }
        }
    }
</script>

<style scoped>
  .search-box {
    max-width: 800px;
  }
</style>
