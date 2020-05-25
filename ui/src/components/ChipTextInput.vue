<template>
  <v-combobox
          :value="value"
          @input="val => $emit('input', val)"
          :items="candidates"
          chips
          clearable
          :label="placeholder"
          multiple
          flat
          solo-inverted
          :prepend-inner-icon="icon"
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
            value: Array,  // for v-model
            icon: String,
        },
        methods: {
            remove(item) {
                this.value.splice(this.value.indexOf(item), 1);
            },
        },
    }
</script>

<style scoped>

</style>