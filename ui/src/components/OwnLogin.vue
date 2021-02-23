<template>
  <v-app>
    <v-overlay v-if="startup">
      <v-progress-circular indeterminate size="64"></v-progress-circular>
    </v-overlay>
    <v-container fill-height v-else>
      <v-row align="center" justify="center">
        <v-col xs="12" sm="5" md="5" lg="4" xl="4">
          <v-card>
            <v-overlay :value="loading">
              <v-progress-circular
                indeterminate
                size="64"
              ></v-progress-circular>
            </v-overlay>
            <v-card-title>OwnDNS</v-card-title>
            <v-card-text>
              <v-row>
                <v-col
                  ><v-text-field
                    label="Имя пользователя"
                    v-model="username"
                    prepend-icon="mdi-face"
                    placeholder=" "
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col
                  ><v-text-field
                    placeholder=" "
                    label="Пароль"
                    v-model="password"
                    prepend-icon="mdi-lock"
                    type="password"
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-btn depressed block color="primary" @click="login()"
                    >Войти</v-btn
                  >
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-app>
</template>
<script>
import sjcl from "sjcl";

async function digestMessage(message) {
  if (crypto.subtle) {
    const msgUint8 = new TextEncoder().encode(message);
    const hashBuffer = await crypto.subtle.digest("SHA-256", msgUint8);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray
      .map((b) => b.toString(16).padStart(2, "0"))
      .join("");
    return hashHex;
  } else {
    return sjcl.codec.hex.fromBits(sjcl.hash.sha256.hash(message));
  }
}

export default {
  data() {
    return {
      username: undefined,
      password: undefined,
      loading: false,
      startup: true,
    };
  },
  mounted() {
    fetch("/user.json")
      .then((data) => data.json())
      .then((json) => {
        if (json.ok) {
          this.$store.dispatch("setuser", json.user);
        }
        this.startup = false;
      })
      .catch((err) => {
        console.log(err);
        this.startup = false;
      });
  },
  methods: {
    login() {
      this.loading = true;
      this.hash.then((p) => {
        fetch("/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            // 'Content-Type': 'application/x-www-form-urlencoded',
          },
          body: JSON.stringify({
            username: this.username, //document.getElementById("username").value,
            password: p,
          }),
        })
          .then((data) => data.json())
          .then((json) => {
            if (json.ok) {
              this.$store.dispatch("setuser", json.user);
              this.loading = false;
            } else {
              this.loading = false;
            }
          })
          .catch((err) => {
            console.log(err);
            this.loading = false;
          });
      });
    },
  },
  computed: {
    hash() {
      return digestMessage(this.password);
    },
  },
};
</script>

<style scoped>
.loading {
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0px;
  top: 0px;
  background-color: rgba(10, 10, 10, 0.8);

  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.loader {
  border: 16px solid #f3f3f3;
  /* Light grey */
  border-top: 16px solid #3498db;
  /* Blue */
  border-radius: 50%;
  width: 80px;
  height: 80px;
  animation: spin 2s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}
</style>
