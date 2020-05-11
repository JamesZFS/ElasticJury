<template>
  <!--            @input="$emit('input', $event.target.value)"-->
  <v-combobox
          :value="value"
          @input="onInput"
          :items="candidates"
          chips
          clearable
          :label="placeholder"
          multiple
          flat
          solo-inverted
          prepend-inner-icon="mdi-magnify"
  >
    <template v-slot:selection="{ attrs, item, select, selected }">
      <v-chip
              v-bind="attrs"
              :input-value="selected"
              close
              @click="select"
              @click:close="remove(item)"
      >
        <strong>{{ item }}</strong>
      </v-chip>
    </template>
  </v-combobox>
</template>

<script>
    export default {
        props: {
            placeholder: String,
            candidates: Array,
            value: Array  // for v-model
        },
        methods: {
            remove(item) {
                this.value.splice(this.value.indexOf(item), 1);
                this.value = [...this.value]
            },
            onInput(value) {
                this.$emit('input', value)
            }
        },
    }
</script>

<style scoped>

</style>