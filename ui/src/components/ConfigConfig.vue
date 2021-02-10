<template>
  <v-card>
    <v-card-title>Config</v-card-title>
    <v-container>
      <v-form v-model="valid">
        <v-row>
          <v-col
            ><v-text-field
              label="Redirect ip"
              v-model="redirectip"
              :rules="[rules.ip]"
          /></v-col>
        </v-row>
      </v-form>
    </v-container>
    <v-card-actions>
      <v-spacer />
      <v-btn color="primary" text @click="undo()">Undo</v-btn>
      <v-btn color="primary" text :disabled="!valid" @click="save()"
        >Save</v-btn
      >
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  data() {
    return {
      redirectip: this.$store.state.config.RedirectIp,
      valid: true,
      rules: {
        ip: (value) => {
          return (
            value.trim().length == 0 ||
            value.trim().match(/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/) != null
          );
        },
      },
    };
  },
  methods: {
    save() {
      this.$store.dispatch("saveconfig", {
        ...this.$store.state.config,
        RedirectIp: this.redirectip,
      });
    },
    undo() {
      this.redirectip = this.$store.state.config.RedirectIp;
    },
  },
};
</script>