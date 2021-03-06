<template>
  <div>
    <h1>{{ username }}</h1>

    <h2>About</h2>

    <vue-markdown
      :linkify="false"
      :html="false"
      :anchorAttributes="{rel: 'ugc' }"
      :source="aboutMarkdown"
      class="user-bio"
    ></vue-markdown>

    <p
      v-if="profileLoaded && !aboutMarkdown"
      class="no-bio-message"
    >This user has not yet created a public bio.</p>

    <template v-if="twitterHandle || emailAddress">
      <h2>Contact</h2>

      <ul>
        <li v-if="emailAddress">
          <a :href="'mailto:' + emailAddress" class="email-address">{{ emailAddress }}</a> (Email)
        </li>
        <li v-if="twitterHandle">
          <a
            :href="'https://twitter.com/' + twitterHandle"
            class="twitter-handle"
          >@{{ twitterHandle }}</a> (Twitter)
        </li>
      </ul>
    </template>

    <b-button v-if="canEdit" to="/profile/edit" variant="primary" class="edit-btn float-right">Edit</b-button>

    <h2>Recent entries</h2>

    <PartialJournal v-bind:key="item.key" v-bind:entry="item" v-for="item in recentEntries" />

    <p
      v-if="entriesLoaded && recentEntries.length === 0"
      class="no-entries-message"
    >This user has not submitted any recent updates.</p>
  </div>
</template>

<script>
import Vue from "vue";
import VueMarkdown from "vue-markdown";
import PartialJournal from "../components/PartialJournal.vue";

Vue.use(VueMarkdown);

export default {
  name: "UserProfile",
  components: {
    VueMarkdown,
    PartialJournal
  },
  data() {
    return {
      aboutMarkdown: "",
      twitterHandle: null,
      emailAddress: null,
      recentEntries: [],
      profileLoaded: false,
      entriesLoaded: false
    };
  },
  computed: {
    username: function() {
      return this.$route.params.username;
    },
    loggedInUsername: function() {
      return this.$store.state.username;
    },
    canEdit: function() {
      return this.loggedInUsername && this.loggedInUsername === this.username;
    }
  },
  methods: {
    clear: function() {
      this.aboutMarkdown = "";
      this.twitterHandle = null;
      this.emailAddress = null;
      this.recentEntries = [];
      this.profileLoaded = false;
      this.entriesLoaded = false;
    },
    loadProfile: function() {
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/user/${this.username}`;
      this.$http
        .get(url)
        .then(result => {
          this.aboutMarkdown = result.data.aboutMarkdown;
          this.twitterHandle = result.data.twitterHandle;
          this.emailAddress = result.data.emailAddress;
          this.profileLoaded = true;
        })
        .catch(error => {
          if (error.response && error.response.status == 404) {
            // A 404 for a user profile is equivalent to retrieving an empty profile.
            this.profileLoaded = true;
          }
        });
    },
    loadRecentEntries: function() {
      this.recentEntries = [];
      const url = `${process.env.VUE_APP_BACKEND_URL}/api/entries/${this.username}`;
      this.$http.get(url).then(result => {
        for (const entry of result.data) {
          this.recentEntries.push({
            key: `${this.username}/${entry.date}`,
            author: this.username,
            date: new Date(entry.date),
            lastModified: new Date(entry.lastModified),
            markdown: entry.markdown
          });
        }
        // Sort newest to oldest.
        this.recentEntries.sort((a, b) => b.date - a.date);
        this.entriesLoaded = true;
      });
    }
  },
  created() {
    this.loadProfile();
    this.loadRecentEntries();
  },
  watch: {
    $route(to, from) {
      if (to.params.username != from.params.username) {
        this.clear();
        this.loadProfile();
        this.loadRecentEntries();
      }
    }
  }
};
</script>

<style scoped>
* {
  text-align: left;
}

h1 {
  margin-bottom: 50px;
}

h2 {
  clear: both;
  margin-top: 40px;
  margin-bottom: 30px;
}

.no-bio-message {
  font-style: italic;
}

.no-entries-message {
  font-style: italic;
}
</style>