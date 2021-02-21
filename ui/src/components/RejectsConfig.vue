<template>
  <array-crud
    :rules="rules"
    dispatch="rejects"
    :items="$store.state.config.RejectNames"
    title="List of reject patterns"
    subtitle="If string Reject is found in the DNS request name, an empty response is formed or the redirect address is transmitted."
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