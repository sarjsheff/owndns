<template>
  <array-crud
    :rules="rules"
    set="setrejectnames"
    :items="$store.state.config.RejectNames"
  ></array-crud>
</template>

<script>
import ArrayCrud from "./ArrayCrud.vue";
export default {
  components: { ArrayCrud },
  data() {
    return {
      rules: [
        (value) => {
          return (value || "").trim().length > 2;
        },
      ],
    };
  },
};
</script>