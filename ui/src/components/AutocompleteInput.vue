<template>
  <v-menu
          v-model="menu"
          close-on-click
          offset-y
  >
    <template v-slot:activator="{ on }">
      <v-text-field
              :value="value"
              @input="onInput"
              :label="placeholder"
              :prepend-inner-icon="icon"
              :loading="loading"
              flat
              solo-inverted
              clearable
              hide-details
              class="mb-8"
      />
    </template>
    <v-list dense>
      <v-list-item
              v-for="(candidate, index) in candidates"
              :key="index"
              @click="onClickCandidate(candidate)"
      >
        {{ candidate }}
      </v-list-item>
    </v-list>
  </v-menu>
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
            menu: null,
            lastWord: ''
        }),
        methods: {
            async onInput(val) { // maybe auto complete
                this.menu = false
                let turn = ++this.turn
                await sleep(1000)
                if (turn !== this.turn) return // only refresh when this is the most recent call
                if (!val || val[val.length - 1] === ' ') return // ignore blank input or when typing ' '
                let i = val.length - 1
                while (i >= 0 && val[i] !== ' ') --i
                this.lastWord = val.slice(++i)
                this.$emit('input', val) // notify parent
                // perform searching (let parent component handle)
                this.loading = true
                // noinspection JSCheckFunctionSignatures
                this.candidates = await this.onAssociate(this.lastWord)
                this.loading = false
                this.menu = true
            },
            onClickCandidate(candidate) {
                // update last word part
                let newValue = this.value.slice(0, this.value.length - this.lastWord.length) + candidate
                this.$emit('input', newValue) // notify parent to update
            }
        }
    }
</script>

<style scoped>

</style>