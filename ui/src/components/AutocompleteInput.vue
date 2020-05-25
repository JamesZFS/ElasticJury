<template>
  <v-autocomplete
          :value="value"
          @input="val => $emit('input', val)"
          :label="placeholder"
          :prepend-inner-icon="icon"
          :items="candidates"
          @update:search-input="onSync"
          :loading="loading"
          flat
          solo-inverted
          clearable
          hide-no-data
          hide-selected
  >
  </v-autocomplete>
</template>

<script>
    import {sleep} from "../utils";

    export default {
        props: {
            placeholder: String,
            value: String,  // for v-model
            icon: String,
            onAssociate: Function, // arg: current input, expected return: candidate array (async func)
        },
        data: () => ({
            loading: false,
            turn: 0,
            candidates: [],
        }),
        methods: {
            async onSync(word) { // maybe auto complete
                if (!word || word === this.value || word.trim().length === 0) return // ignore blank input
                let turn = ++this.turn
                await sleep(1000)
                if (turn !== this.turn) return // only refresh when this is the most recent call
                // perform searching (let parent component handle)
                this.loading = true
                // noinspection JSCheckFunctionSignatures
                this.candidates = await this.onAssociate(word)
                this.loading = false
            }
        }
    }
</script>

<style scoped>

</style>