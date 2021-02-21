<template>
  <array-crud
    :rules="rules"
    dispatch="relaydns"
    :items="$store.state.config.RelayDNS"
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