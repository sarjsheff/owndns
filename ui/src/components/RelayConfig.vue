<template>
  <array-crud
    set="setrelaydns"
    :rules="rules"
    :items="$store.state.config.RelayDns"
    title="List of relay DNS"
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
          return (
            (value || "")
              .trim()
              .match(/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}$/) != null
          );
        },
      ],
    };
  },
};
</script>