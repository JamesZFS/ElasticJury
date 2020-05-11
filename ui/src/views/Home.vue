<template>
  <v-container>
    <v-row justify="center" class="text-center">
      <v-img
              src="../assets/logo.svg"
              class="my-3"
              contain
              height="200"
      />
    </v-row>
    <v-row justify="center" align="center" style="height: 20vh">
      <h1 class="display-2 font-weight-bold mb-3">
        Welcome to ElasticJury
      </h1>
    </v-row>
    <v-row justify="center" align="center" class="mx-auto mt-2 search-box">
      <ChipTextInput
              placeholder="请输入关键词..."
              v-model="search.inputs"
              :candidates="search.candidates"
      />
    </v-row>

    <v-row justify="center">
      <v-btn @click="onSearch" class="mr-10">Search</v-btn>
      <v-btn @click="onInfo" class="mr-10">Info</v-btn>
      <v-btn @click="onPing">Ping</v-btn>
    </v-row>

  </v-container>
</template>

<script>
    import ChipTextInput from "../components/ChipTextInput";
    import {getCaseInfo, ping, searchCaseId} from "../api";

    export default {
        name: 'Home',
        components: {ChipTextInput},
        data: () => ({
            search: {
                inputs: [],
                candidates: ['调解', '协议', '当事人'],
            },
            result: {
                ids: [],
                info: null,
            }
        }),
        methods: {
            async onSearch() {
                console.log(this.search.inputs);
                let resp = await searchCaseId(this.search.inputs);
                alert(JSON.stringify(resp));
                this.result.ids = Object.entries(resp.result)
                    // .sort(([_id1, val1], [_id2, val2]) => val2 - val1)
                    .map(([id, _val]) => parseInt(id));
            },
            async onInfo() {
                console.log(this.result.ids);
                let resp = await getCaseInfo(this.result.ids);
                this.result.info = resp;
                console.log(resp);
            },
            async onPing() {
                let resp = await ping();
                alert(resp);
            }
        }
    }
</script>

<style scoped>
  /*h1, p {*/
  /*  text-align: center;*/
  /*}*/
  .search-box {
    height: 10vh;
    max-width: 800px;
  }
</style>
