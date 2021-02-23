import Vue from "vue";
import Vuex from "vuex";
import owndns from "../plugins/owndns"

Vue.use(owndns);
Vue.use(Vuex);

function authCheck(context) {
  return (data) => {
    //console.log("FETCH STATUS", data.status);
    if (data.status == 403) {
      context.commit("setuser", undefined);
      return undefined;
    } else {
      return data.json();
    }
  }
}


const arraycrud = (name) => {
  return {
    namespaced: true,
    actions: {
      save(context, item) {
        return new Promise((res, rej) => {
          fetch(`/${name}/save`, {
            method: "POST",
            cache: "no-cache",
            headers: {
              "Content-Type": "application/json",
              // 'Content-Type': 'application/x-www-form-urlencoded',
            },
            referrerPolicy: "no-referrer",
            body: JSON.stringify(item),
          })
            .then(authCheck(context))
            .then((json) => {
              if (json.ok) {
                context.dispatch("getconfig", undefined, { root: true });
                res();
              } else {
                rej(json);
              }
            }).catch(e => { rej(e) });
        });
      },
      add(context, item) {
        return new Promise((res, rej) => {
          fetch(`/${name}/add`, {
            method: "POST",
            cache: "no-cache",
            headers: {
              "Content-Type": "application/json",
              // 'Content-Type': 'application/x-www-form-urlencoded',
            },
            referrerPolicy: "no-referrer",
            body: JSON.stringify({ name: item }),
          })
            .then(authCheck(context))
            .then((json) => {
              if (json.ok) {
                context.dispatch("getconfig", undefined, { root: true });
                res();
              } else {
                rej(json);
              }
            }).catch(e => { rej(e) });
        });
      },
      del(context, id) {
        return new Promise((res, rej) => {
          fetch(`/${name}/del`, {
            method: "POST",
            cache: "no-cache",
            headers: {
              "Content-Type": "application/json",
              // 'Content-Type': 'application/x-www-form-urlencoded',
            },
            referrerPolicy: "no-referrer",
            body: JSON.stringify({ id: id }),
          })
            .then(authCheck(context))
            .then((json) => {
              if (json.ok) {
                context.dispatch("getconfig", undefined, { root: true });
                res();
              } else {
                rej(json);
              }
            }).catch(e => { rej(e) });
        });
      }


    }
  }
}

export default new Vuex.Store({
  state: {
    version: undefined,
    stat: {},
    lastlog: [],
    config: {},
    ws: undefined,
    user: undefined
  },
  mutations: {
    setuser(state, user) {
      state.user = user;
    },
    setws(state, ws) {
      state.ws = ws;
    },
    resetws(state) {
      state.ws = undefined;
    },
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
    logout(context) {
      fetch("/logout").then(() => context.commit("setuser", undefined)).catch(() => context.commit("setuser", undefined))
    },
    setuser(context, user) {
      console.log("Must REFRESH");
      context.dispatch("getversion");
      context.dispatch("getconfig");
      context.commit("setuser", user);
    },
    wsconnect(context) {
      var ws = new WebSocket("ws://localhost:8081/ws");
      ws.onopen = function () {
        // console.log("OPEN", evt);
        context.commit("setws", ws);
        ws.send({ Op: "Version" });
      };
      ws.onclose = function (evt) {
        console.log("CLOSE", evt);
        context.commit("resetws");
        setTimeout(() => context.dispatch("wsconnect"), 1000);
      };
      ws.onmessage = function (evt) {
        console.log("RESPONSE: ", evt);
      };
      ws.onerror = function (evt) {
        console.log("ERROR: ", evt);
      };
    },
    getconfig(context) {
      fetch("/config", { cache: "no-cache" })
        .then(authCheck(context))
        .then((json) => {
          context.commit("setconfig", json);
        })
        .catch((e) => {
          console.log(e);
          //context.commit("setlastlog", []);
        });
    },
    // saveconfig(context, config) {
    //   fetch(`/config/save`, {
    //     method: "POST",
    //     cache: "no-cache",
    //     headers: {
    //       "Content-Type": "application/json",
    //       // 'Content-Type': 'application/x-www-form-urlencoded',
    //     },
    //     referrerPolicy: "no-referrer",
    //     body: JSON.stringify(config),
    //   })
    //     .then((data) => data.json())
    //     .then(() => {
    //       context.dispatch("getconfig");
    //     });

    // },
    // setrelaydns(context, names) {
    //   context.dispatch("saveconfig", { ...context.state.config, RelayDns: names });
    // },
    // setrejectnames(context, names) {
    //   context.dispatch("saveconfig", { ...context.state.config, RejectNames: names });
    // },
    getlastlog(context) {
      fetch("/log.json")
        .then(authCheck(context))
        .then((json) => {
          context.commit("setlastlog", json);
        })
        .catch((e) => {
          console.log(e);
          context.commit("setlastlog", []);
        });
    },
    getversion(context) {
      fetch(`/version.json`, { cache: "no-cache" })
        .then(authCheck(context))
        .then((json) => {
          context.commit("setversion", json ? json.version || "" : "");
        })
        .catch((e) => {
          console.log("versionerror", e);
          context.commit("setversion", "?");
        });
    },
    getstat(context) {
      fetch("/stat.json")
        .then(authCheck(context))
        .then((json) => {
          console.log(json);
          context.commit("setstat", json || {});
        })
        .catch((e) => {
          console.log(e);
          context.commit("setstat", {});
        });
    },
    setvalue(context, data) {
      fetch(`/setvalue`, {
        method: "POST",
        cache: "no-cache",
        headers: {
          "Content-Type": "application/json",
          // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        referrerPolicy: "no-referrer",
        body: JSON.stringify(data),
      })
        .then(authCheck(context))
        .then(() => {
          context.dispatch("getconfig");
        });
    }
  },
  modules: {
    rejects: arraycrud("rejects"),
    relaydns: arraycrud("relaydns"),
    rules: arraycrud("rules"),
    users: arraycrud("users"),
  },
});
