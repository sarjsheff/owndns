<template>
  <v-app>
    <v-navigation-drawer permanent app>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="title">OwnDns</v-list-item-title>
          <v-list-item-subtitle>{{
            $store.state.version
          }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>

      <v-divider></v-divider>

      <v-list dense nav>
        <v-list-item
          v-for="item in items"
          :key="item.title"
          link
          :input-value="selected === item.title"
          @click="select(item)"
        >
          <v-list-item-icon>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app color="primary" dark>
      <v-app-bar-nav-icon><v-icon>mdi-dns</v-icon></v-app-bar-nav-icon>
      <v-toolbar-title>{{ selected }}</v-toolbar-title>
      <v-spacer></v-spacer>
    </v-app-bar>

    <v-main>
      <v-container>
        <v-row v-if="selected == 'Statistics'">
          <v-col>
            <dns-stat />
          </v-col>
        </v-row>
        <v-row v-if="selected == 'Logs'">
          <v-col>
            <last-log />
          </v-col>
        </v-row>
        <v-row v-if="selected == 'Rules'">
          <v-col>
            <rules-config />
          </v-col>
        </v-row>
        <v-row v-if="selected == 'Rejects'">
          <v-col>
            <rejects-config />
          </v-col>
        </v-row>
        <v-row v-if="selected == 'Relay DNS'">
          <v-col>
            <relay-config />
          </v-col>
        </v-row>
        <v-row v-if="selected == 'Config'">
          <v-col><config-config /></v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import ConfigConfig from "./components/ConfigConfig.vue";
import DnsStat from "./components/DnsStat";
import LastLog from "./components/LastLog.vue";
import RejectsConfig from "./components/RejectsConfig.vue";
import RelayConfig from "./components/RelayConfig.vue";
import RulesConfig from "./components/RulesConfig.vue";

export default {
  name: "App",

  components: {
    DnsStat,
    LastLog,
    RejectsConfig,
    RelayConfig,
    ConfigConfig,
    RulesConfig,
  },

  data: () => ({
    selected: "Statistics",
    items: [
      { title: "Statistics", icon: "mdi-dns" },
      { title: "Logs", icon: "mdi-dns" },
      { title: "Rules", icon: "mdi-dns" },
      { title: "Rejects", icon: "mdi-dns" },
      { title: "Relay DNS", icon: "mdi-dns" },
      { title: "Config", icon: "mdi-dns" },
    ],
  }),
  mounted() {
    this.$store.dispatch("getversion");
    this.$store.dispatch("getconfig");
  },
  methods: {
    select(item) {
      this.selected = item.title;
    },
  },
};
</script>
