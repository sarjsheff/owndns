import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    version: undefined,
    stat: {},
    lastlog: [],
    config: {}
  },
  mutations: {
    setconfig(state, config) {
      state.config = config;
    },
    setlastlog(state, lastlog) {
      state.lastlog = lastlog;
    },
    setversion(state, version) {
      state.version = version;
    },
    setstat(state, stat) {
      state.stat = stat;
    },
  },
  actions: {
    getconfig(context) {
      fetch("/config")
        .then((data) => data.json())
        .then((json) => {
          context.commit("setconfig", json);
        })
        .catch((e) => {
          console.log(e);
          //context.commit("setlastlog", []);
        });
    },
    saveconfig(context, config) {
      fetch(`/config/save`, {
        method: "POST",
        cache: "no-cache",
        headers: {
          "Content-Type": "application/json",
          // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        referrerPolicy: "no-referrer",
        body: JSON.stringify(config),
      })
        .then((data) => data.json())
        .then(() => {
          context.dispatch("getconfig");
        });

    },
    setrelaydns(context, names) {
      context.dispatch("saveconfig", { ...context.state.config, RelayDns: names });
    },
    setrejectnames(context, names) {
      context.dispatch("saveconfig", { ...context.state.config, RejectNames: names });
    },
    getlastlog(context) {
      fetch("/log.json")
        .then((data) => data.json())
        .then((json) => {
          context.commit("setlastlog", json);
        })
        .catch((e) => {
          console.log(e);
          context.commit("setlastlog", []);
        });
    },
    getversion(context) {
      fetch("/version.json")
        .then((data) => data.json())
        .then((json) => {
          console.log(json);
          context.commit("setversion", json ? json.version || "" : "");
        })
        .catch((e) => {
          console.log(e);
          context.commit("setversion", "?");
        });
    },
    getstat(context) {
      fetch("/stat.json")
        .then((data) => data.json())
        .then((json) => {
          console.log(json);
          context.commit("setstat", json || {});
        })
        .catch((e) => {
          console.log(e);
          context.commit("setstat", {});
        });
    },
  },
  modules: {
    // rejects: arraycrud("rejects"),
    // relaydns: arraycrud("relaydns"),
  },
});
