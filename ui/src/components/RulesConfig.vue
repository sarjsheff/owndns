<template>
  <object-crud set="setrules"></object-crud>
</template>

<script>
import ObjectCrud from "./ObjectCrud.vue";
export default {
  components: { ObjectCrud },
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