<template>
  <crud-view
    :items="items"
    :form="form"
    @reset="reset()"
    @edit="edit($event)"
    @save="save()"
    @add="add()"
    @del="del($event)"
  >
    <template v-slot:form="{ form }">
      <v-row>
        <v-col
          ><v-text-field label="Name" v-model="form.name" :rules="rules"
        /></v-col>
      </v-row>
    </template>
  </crud-view>
</template>

<script>
import CrudView from "./CrudView.vue";
export default {
  components: { CrudView },
  props: ["items", "set", "rules"],
  data() {
    return {
      form: { id: undefined, name: undefined },
    };
  },
  methods: {
    reset() {
      this.form.id = undefined;
      this.form.name = undefined;
    },
    edit(id) {
      this.form.id = id;
      this.form.name = this.items[id];
      console.log(this.items[id]);
    },
    save() {
      this.$store.dispatch(
        this.set,
        this.items.map((e, idx) => {
          if (idx === this.form.id) {
            return this.form.name;
          } else {
            return e;
          }
        })
      );
    },
    add() {
      this.$store.dispatch(this.set, [...this.items, this.form.name]);
    },
    del(id) {
      this.$store.dispatch(
        this.set,
        this.items.filter((e, idx) => {
          return idx != id;
        })
      );
    },
  },
};
</script>