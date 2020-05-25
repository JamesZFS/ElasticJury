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
              @focus="focus = true"
              @blur="focus = false"
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
            history: Array,
        },
        data: () => ({
            loading: false,
            turn: 0,
            candidates: [],
            menu: null,
            focus: false,
            lastWord: '',
            chosen: [],
        }),
        methods: {
            async onInput(val) { // maybe auto complete
                this.menu = false
                this.$emit('input', val) // notify parent
                let turn = ++this.turn
                await sleep(500)
                if (turn !== this.turn || !this.focus) return // only refresh when this is the most recent call
                if (!val || val[val.length - 1] === ' ') {  // ignore blank input or when typing ' '
                    this.lastWord = ''
                    this.candidates = this.history // show suggestions
                    this.menu = this.candidates ? true : false
                    return
                }
                let i = val.length - 1
                while (i >= 0 && val[i] !== ' ') --i
                this.lastWord = val.slice(++i)
                // perform searching (let parent component handle)
                this.loading = true
                // noinspection JSCheckFunctionSignatures
                this.candidates = await this.onAssociate(this.lastWord)
                this.loading = false
                this.menu = this.candidates ? true : false
            },
            onClickCandidate(candidate) {
                // update last word part
                let newValue = this.value
                    ? this.value.slice(0, this.value.length - this.lastWord.length) + candidate
                    : candidate
                this.$emit('input', newValue) // notify parent to update
                this.chosen.push(candidate)
            }
        }
    }
</script>

<style scoped>

</style>