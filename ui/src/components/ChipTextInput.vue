<template>
  <v-combobox
          :value="value"
          @input="onInput"
          :search-input.sync="search"
          :items="candidates"
          chips
          clearable
          :label="placeholder"
          multiple
          flat
          deletable-chips
          solo-inverted
          hide-selected
          :prepend-inner-icon="icon"
          :loading="loading"
  >
  </v-combobox>
</template>

<script>
    import {sleep} from "../utils";

    export default {
        props: {
            placeholder: String,
            history: Array,
            value: Array,  // for v-model
            icon: String,
            onAssociate: Function,
        },
        data: () => ({
            loading: false,
            turn: 0,
            candidates: [],
            search: '',
        }),
        methods: {
            onInput(val) {
                this.candidates = []
                this.search = ''
                this.$emit('input', val)
            },
        },
        watch: {
            async search(chip) { // maybe auto complete
                this.candidates = []
                let turn = ++this.turn
                await sleep(1000)
                if (turn !== this.turn) return // only refresh when this is the most recent call
                if (!chip) {
                    this.candidates = this.history
                    return
                }
                // perform searching (let parent component handle)
                this.loading = true
                // noinspection JSCheckFunctionSignatures
                this.candidates = await this.onAssociate(chip)
                this.loading = false
            },
        }
    }
</script>

<style scoped>

</style>