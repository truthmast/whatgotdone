<template>
  <div>
    <h1>Recent Entries</h1>
    <p>Check out what other users got done this week:</p>

    <PartialJournal v-bind:key="item.key" v-bind:entry="item" v-for="item in recentEntries"/>
  </div>
</template>

<script>
import PartialJournal from "../components/PartialJournal.vue";

import refreshRecent from "../controllers/Recent.js";

export default {
  name: "Recent",
  components: {
    PartialJournal
  },
  computed: {
    recentEntries() {
      return this.$store.state.recentEntries;
    }
  },
  created() {
    refreshRecent();
  }
};
</script>

<style scoped>
h1 {
  text-align: left;
}

p {
  text-align: left;
}
</style>