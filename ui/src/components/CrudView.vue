<template>
  <v-card>
    <v-dialog v-if="deleteitem !== undefined" v-model="deletedialog">
      <v-card>
        <v-card-title class="headline">Confirm delete.</v-card-title>
        <v-card-text>
          Delete
          {{
            typeof items[deleteitem] === "string"
              ? items[deleteitem]
              : deleteitem + 1
          }}?
        </v-card-text>
        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" text @click="closeDeletedialog()"> No </v-btn>
          <v-btn color="primary" text @click="deleteItem()"> Yes </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="iseditdialog">
      <v-card>
        <v-card-title class="headline grey lighten-2">{{
          editdialog ? "Edit" : "Add"
        }}</v-card-title>
        <v-container>
          <v-alert v-if="adderror" dense type="error">{{
            this.adderror
          }}</v-alert>
          <v-form v-model="valid" ref="formref">
            <slot
              name="form"
              v-bind:form="form"
              v-bind:editdialog="editdialog"
            ></slot>
          </v-form>
        </v-container>
        <v-card-actions>
          <v-spacer />
          <v-btn
            v-if="adddialog"
            color="primary"
            :disabled="!valid"
            text
            @click="addItem()"
            >Add</v-btn
          >
          <v-btn
            v-if="editdialog"
            color="primary"
            :disabled="!valid"
            text
            @click="saveItem()"
            >Save</v-btn
          >
          <v-btn color="primary" text @click="closeAdddialog()">Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-toolbar flat>
      <v-toolbar-title>{{ title || "List" }}</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn icon @click="openAdddialog()">
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-toolbar>
    <!-- <v-card-title>Rejects list</v-card-title>
    <v-card-subtitle>Name must contain</v-card-subtitle> -->
    <v-divider></v-divider>
    <v-card-text v-if="subtitle">{{ subtitle }}</v-card-text>
    <slot
      name="list"
      v-bind:items="items"
      v-bind:openEditdialog="openEditdialog"
      v-bind:openDeletedialog="openDeletedialog"
    >
      <v-list>
        <v-list-item v-for="(name, i) in items" :key="i">
          <v-list-item-content>{{ name }}</v-list-item-content>
          <v-list-item-action>
            <v-btn icon>
              <v-icon color="grey lighten-1" @click="openEditdialog(i)"
                >mdi-pencil</v-icon
              >
            </v-btn>
          </v-list-item-action>
          <v-list-item-action>
            <v-btn icon>
              <v-icon color="grey lighten-1" @click="openDeletedialog(i)"
                >mdi-delete</v-icon
              >
            </v-btn>
          </v-list-item-action>
        </v-list-item>
      </v-list>
    </slot>
  </v-card>
</template>

<script>
export default {
  props: ["name", "items", "form", "title", "subtitle", "add", "del", "save"],
  data() {
    return {
      deleteitem: undefined,
      deletedialog: false,
      editdialog: false,
      adddialog: false,
      adderror: undefined,
      delerror: undefined,
      saveerror: undefined,
      valid: true,
      // form: {
      //   id: undefined,
      //   name: undefined,
      // },
    };
  },
  // mounted() {
  //this.$store.dispatch(`${this.name}/get`);
  // },

  computed: {
    iseditdialog() {
      return this.adddialog || this.editdialog;
    },
  },
  methods: {
    openEditdialog(id) {
      this.$emit("edit", id);
      this.editdialog = true;
      this.$nextTick(() => {
        this.$refs.formref.validate();
      });
      // this.form.id = id;
      // this.form.name = this.items[id];
      // console.log(this.items[id]);
    },
    openAdddialog() {
      this.$emit("reset");
      this.adddialog = true;
      this.$nextTick(() => {
        this.$refs.formref.validate();
      });
    },
    closeAdddialog() {
      this.editdialog = false;
      this.adddialog = false;
      this.$emit("reset");
    },
    saveItem() {
      if (this.save) {
        this.save()
          .then(() => {
            this.$emit("save");
            this.closeAdddialog();
          })
          .catch((e) => {
            if (e.error) {
              this.saveerror = e.error;
            }
            this.valid = false;
          });
      } else {
        this.$emit("save");
        this.closeAdddialog();
      }
    },
    addItem() {
      if (this.add) {
        this.add()
          .then(() => {
            this.$emit("add");
            this.closeAdddialog();
          })
          .catch((e) => {
            if (e.error) {
              this.adderror = e.error;
            }
            this.valid = false;
          });
      } else {
        this.$emit("add");
        this.closeAdddialog();
      }
    },
    deleteItem() {
      if (this.del) {
        this.del(this.deleteitem)
          .then(() => {
            this.$emit("del", this.deleteitem);
            this.closeDeletedialog();
          })
          .catch((e) => {
            if (e.error) {
              this.delerror = e.error;
            }
          });
      } else {
        this.$emit("del", this.deleteitem);
        this.closeDeletedialog();
      }
    },

    openDeletedialog(i) {
      this.deleteitem = i;
      this.deletedialog = true;
    },
    closeDeletedialog() {
      this.deletedialog = false;
      this.deleteitem = undefined;
    },
  },
};
</script>