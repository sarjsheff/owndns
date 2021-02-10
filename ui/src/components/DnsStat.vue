<template>
  <v-card>
    <v-card-title>Statistics</v-card-title>
    <v-container>
      <v-expansion-panels>
        <v-expansion-panel
          v-for="ip in Object.keys($store.state.stat)"
          :key="ip"
        >
          <v-expansion-panel-header>{{ ip }}</v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-simple-table fixed-header height="300px">
              <template v-slot:default>
                <thead>
                  <tr>
                    <th class="text-left">Name</th>
                    <th class="text-left">Query count</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="dname in Object.keys($store.state.stat[ip]).sort(
                      (a, b) => {
                        return (
                          $store.state.stat[ip][b] - $store.state.stat[ip][a]
                        );
                      }
                    )"
                    :key="dname"
                  >
                    <td>{{ dname }}</td>
                    <td>{{ $store.state.stat[ip][dname] }}</td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-expansion-panel-content>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-container>
  </v-card>
</template>

<script>
export default {
  mounted() {
    this.$store.dispatch("getstat");
  },
};
</script>