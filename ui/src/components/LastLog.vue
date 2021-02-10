<template>
  <v-card>
    <v-card-title>Last requests</v-card-title>
    <v-simple-table dense>
      <template v-slot:default>
        <thead>
          <tr>
            <th class="text-left">Date</th>
            <th class="text-left">Client IP</th>
            <th class="text-left">Name</th>
            <th class="text-left">Answer</th>
            <th class="text-left"></th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(item, i) in $store.state.lastlog"
            :key="`${i}#${item.Date}`"
          >
            <td>{{ item.Date }}</td>
            <td>{{ item.Ip }}</td>
            <td>{{ item.Q.Name }}</td>
            <td>
              <v-chip
                class="ma-2"
                label
                small
                v-for="(d, i) in item.Answer"
                :key="i"
              >
                {{ answer(d) }}
              </v-chip>
            </td>
            <td><v-icon small v-if="item.Rejects">mdi-cancel</v-icon></td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
  </v-card>
</template>

<script>
export default {
  mounted() {
    this.$store.dispatch("getlastlog");
  },
  methods: {
    answer(d) {
      if (d.A != undefined) {
        return `${d.Hdr.Name} A ${d.A}`;
      } else if (d.Target != undefined) {
        return `${d.Hdr.Name} Target ${d.Target}`;
      } else {
        return d;
      }
    },
  },
};
</script>