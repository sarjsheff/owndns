<template>
  <crud-view
    :items="$store.state.config.Users || {}"
    :form="form"
    :add="addItem"
    :del="delItem"
    :save="saveItem"
    @reset="reset()"
    @edit="edit($event)"
    title="Users"
    subtitle=""
  >
    <template v-slot:list="{ items, openEditdialog, openDeletedialog }">
      <v-simple-table>
        <template v-slot:default>
          <thead>
            <tr>
              <th class="text-left">Username</th>
              <th class="text-left">isAdmin</th>
              <th class="text-right"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="key in Object.keys(items)" :key="key">
              <td>{{ key }}</td>
              <td>
                <v-icon v-if="items[key].IsAdmin === true">mdi-check</v-icon
                ><v-icon v-else>mdi-cancel</v-icon>
              </td>
              <td class="text-right">
                <v-btn icon>
                  <v-icon color="grey lighten-1" @click="openEditdialog(key)">
                    mdi-pencil
                  </v-icon>
                </v-btn>
                <v-btn icon>
                  <v-icon color="grey lighten-1" @click="openDeletedialog(key)">
                    mdi-delete
                  </v-icon>
                </v-btn>
              </td>
            </tr>
          </tbody>
        </template>
      </v-simple-table>
    </template>
    <template v-slot:form="{ form, editdialog }">
      <v-row>
        <v-col>
          <v-text-field
            v-if="editdialog"
            label="Username"
            :value="form.id"
            disabled
          />
          <v-text-field
            v-else
            label="Username"
            v-model="form.name.Username"
            autocomplete="username"
            :rules="usernameRules"
          />
        </v-col>
      </v-row>
      <v-row v-if="!editdialog">
        <v-col
          ><v-text-field
            label="Password"
            v-model="form.name.Password"
            type="password"
            autocomplete="new-password"
            :rules="passwordRules"
        /></v-col>
      </v-row>
      <v-row v-if="!editdialog">
        <v-col
          ><v-text-field
            label="Repeat password"
            v-model="form.name.RePassword"
            type="password"
            autocomplete="new-password"
            :rules="repasswordRule"
        /></v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-switch
            v-model="form.name.IsAdmin"
            inset
            label="Is admin"
          ></v-switch>
        </v-col>
      </v-row>
    </template>
  </crud-view>
</template>

<script>
import CrudView from "./CrudView.vue";
export default {
  components: { CrudView },
  // props: ["items", "set"],
  data() {
    return {
      form: { id: undefined, name: {} },
      dispatch: "users",
      usernameRules: [
        (value) => {
          return (value || "").trim().length > 3;
        },
      ],
      passwordRules: [
        (value) => {
          return (value || "").trim().length > 5;
        },
      ],
      repasswordRules: [
        (value) => {
          return (
            (value || "").trim().length > 5 && value === this.form.password
          );
        },
      ],
    };
  },
  computed: {
    repasswordRule() {
      return [
        (value) => {
          return (
            (value || "").trim().length > 5 && value === this.form.name.Password
          );
        },
      ];
    },
  },
  methods: {
    reset() {
      this.form.id = undefined;
      this.form.name = {};
    },
    edit(id) {
      this.form.id = id;
      this.form.name = this.$store.state.config.Users[id];
    },
    // save() {
    //   this.$store.dispatch(`${this.dispatch}/save`, this.form);
    // },
    saveItem() {
      return this.$store.dispatch(`${this.dispatch}/save`, this.form);
    },
    // add() {
    //   this.$store.dispatch(`${this.dispatch}/add`, this.form.name);
    // },
    addItem() {
      return this.$store.dispatch(`${this.dispatch}/add`, this.form.name);
    },
    // del(id) {
    //   this.$store.dispatch(`${this.dispatch}/del`, id);
    // },
    delItem(id) {
      return this.$store.dispatch(`${this.dispatch}/del`, id);
    },
  },
};
</script>