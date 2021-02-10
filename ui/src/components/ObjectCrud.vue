<template>
  <crud-view
    :items="$store.state.config.Rules"
    :form="form"
    @reset="reset()"
    @edit="edit($event)"
    @save="save()"
    @add="add()"
    @del="del($event)"
    title="List of rules"
    subtitle="If string Reject is found in the DNS request name, an empty response is formed or the redirect address is transmitted. If line Accept is found in the name of the request, processing continues in spite of the existing Rejects. The * character handles all possible situations. Accept = * resets all Reject, Reject = * does not process all incoming requests except Accept."
  >
    <template v-slot:list="{ items, openEditdialog, openDeletedialog }">
      <v-simple-table>
        <template v-slot:default>
          <thead>
            <tr>
              <th class="text-left">IP</th>
              <th class="text-left">Reject</th>
              <th class="text-left">Accept</th>
              <th class="text-right"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, id) in items" :key="id">
              <td>{{ item.Ip }}</td>
              <td>{{ item.RejectName }}</td>
              <td>{{ item.AcceptOnly }}</td>
              <td class="text-right">
                <v-btn icon>
                  <v-icon color="grey lighten-1" @click="openEditdialog(id)">
                    mdi-pencil
                  </v-icon>
                </v-btn>
                <v-btn icon>
                  <v-icon color="grey lighten-1" @click="openDeletedialog(id)">
                    mdi-delete
                  </v-icon>
                </v-btn>
              </td>
            </tr>
          </tbody>
        </template>
      </v-simple-table>
    </template>
    <template v-slot:form="{ form }">
      <v-row>
        <v-col>
          <v-text-field label="Ip" v-model="form.name.Ip" :rules="[checkip]" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field label="Reject" v-model="form.name.RejectName" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field label="Accept" v-model="form.name.AcceptOnly" />
        </v-col>
      </v-row>
    </template>
  </crud-view>
</template>

<script>
import CrudView from "./CrudView.vue";
export default {
  components: { CrudView },
  props: ["items", "set"],
  data() {
    return {
      form: { id: undefined, name: {} },
    };
  },
  methods: {
    checkip: (value) => {
      return (
        (value || "").trim().match(/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/) !=
        null
      );
    },

    reset() {
      this.form.id = undefined;
      this.form.name = {};
    },
    edit(id) {
      this.form.id = id;
      this.form.name = this.$store.state.config.Rules[id];
    },
    save() {
      this.$store.dispatch("saveconfig", {
        ...this.$store.state.config,
        Rules: this.$store.state.config.Rules.map((e, idx) => {
          if (idx === this.form.id) {
            return this.form.name;
          } else {
            return e;
          }
        }),
      });
    },
    add() {
      this.$store.dispatch("saveconfig", {
        ...this.$store.state.config,
        Rules: [...(this.$store.state.config.Rules || []), this.form.name],
      });
    },
    del(id) {
      this.$store.dispatch("saveconfig", {
        ...this.$store.state.config,
        Rules: this.$store.state.config.Rules.filter((e, idx) => {
          return idx != id;
        }),
      });
    },
  },
};
</script>