<template>
  <div id="formdiv">
    <div class="m-auto mb-2 mt-4" w="full" text="left" style="max-width: 500px">
      <code>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
        <span
          class="text-gray-400"
        >{{ timestamp }} UTC</span>
        <br />&nbsp;&nbsp;&nbsp;&nbsp;env service dd [hh:mm:ss]
      </code>
    </div>

    <input
      class="col-span-3 row-span-1"
      id="inputmain"
      v-model="name"
      placeholder="stag back t 14"
      type="text"
      autocomplete="false"
      p="x-4 y-2"
      w="full"
      text="left"
      border="~ rounded gray-700"
      outline="none active:none"
      @keydown.enter="go"
      style="max-width: 500px"
    />

    <div w="full" class="m-auto mb-2 mt-5" style="max-width: 500px">
      <p class="text-true-gray-300 mx-4" text="left">field=regex or field!=regex</p>
    </div>
    <input
      class="col-span-3 row-span-1"
      id="inputregexq"
      v-model="regexq"
      placeholder="record.message=GET .* 200"
      type="text"
      autocomplete="false"
      p="x-4 y-2"
      w="full"
      text="left"
      border="~ rounded gray-700"
      outline="none active:none"
      @keydown.enter="go"
      style="max-width: 500px"
    />

    <div w="full" class="m-auto mb-2 mt-5" style="max-width: 500px">
      <p class="text-true-gray-300 mx-4" text="left">Show only these fields:</p>
    </div>
    <input
      class="col-span-3 row-span-1"
      id="inputfields"
      v-model="fields"
      placeholder="fluentd_time record.function"
      type="text"
      autocomplete="false"
      p="x-4 y-2"
      w="full"
      text="left"
      border="~ rounded gray-700"
      outline="none active:none"
      @keydown.enter="go"
      style="max-width: 500px"
    />

    <div
      class="gap-5 grid grid-rows-1 grid-cols-3 m-auto"
      w="full"
      style="max-width: 500px; max-height: 480px"
    >
      <p class="col-span-1 row-span-1"></p>

      <div class="col-span-1 row-span-1">
        <button
          w="full"
          style="max-width: 310px"
          class="mt-5 p-2 text-mm btn truncate"
          :disabled="!name"
          @click="go"
        >History</button>
      </div>

      <div class="col-span-1 row-span-1">
        <button
          w="full"
          style="max-width: 310px"
          class="mt-5 p-2 text-mm btn truncate"
          :disabled="!name"
          @click="go"
        >{{ mainbtntext }}</button>
      </div>

      <div
        style="max-width: 400px"
        class="mt-4 font-bold col-span-3 row-span-1 my-4"
        :class="respFeedbackColor"
      >{{ respFeedback }}</div>
    </div>

    <div class="my-4" style="margin-left: -40px;">
      <label class="contcb" for="scales">
        &nbsp;Truncate long lines
        <input
          type="checkbox"
          id="scales"
          name="scales"
          @click="truncLogFunc()"
        />
        <span class="checkmark"></span>
      </label>
    </div>
  </div>

  <div class="fixed bottom-10 p-6"></div>
  <!-- <hr w="auto" class="my-6 mx-auto" style="max-width: 500px" /> -->
  <div :class="showhelp" id="helpdiv" class="mt-14 pr-8 pb-10">
    <h3 class="text-lg mb-8">Examples</h3>
    <ul text="left" class="list-none list-inside text-sm">
      <li class="pb-4">
        <pre>example.com backend 23 16:45:23</pre>
      </li>
      <li class="pb-4">
        <pre>example.com backend 23 16:45:23</pre>
      </li>
      <li class="pb-4">
        <pre>exampl back 23 16:45:23</pre>
      </li>
      <li class="pb-4">
        <pre>exampl back t 16:45:</pre>
      </li>
      <li class="pb-4">
        <pre>exampl back t 16:</pre>
      </li>
      <li class="pb-4">
        <pre>exampl back t 10m</pre>
      </li>
    </ul>
  </div>
  <div :class="showlogs" id="logsdiv" class="mt-16">
    <vue-json-pretty
      class="leading-none break-all"
      :class="truncatelog"
      style="max-width: 1500px; font-size: 12px !important;"
      :path="'res'"
      :data="jsonData"
      :deep="2"
      :showDoubleQuotes="false"
      :showLine="true"
    ></vue-json-pretty>
  </div>
</template>

<style>
@media (min-width: 800px) {
  #helpdiv {
    position: absolute;
    left: 48vw;
    max-width: 60%;
  }

  #logsdiv {
    position: absolute;
    right: 0px;
    max-width: 62%;
    min-width: 62%;
  }

  #formdiv {
    top: 80px;
    position: fixed;
    max-width: 42vw;
    min-width: 35vw;
  }
}

.vjs-tree__node.is-highlight,
.vjs-tree__node:hover {
  background-color: #262626 !important;
}

#tokeninput::placeholder {
  font-size: 0.7em;
}

/* The container */
.contcb {
  margin-top: 10px;
  position: relative;
  padding-left: 35px;
  margin-bottom: 12px;
  cursor: pointer;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

/* Hide the browser's default checkbox */
.contcb input {
  position: absolute;
  opacity: 0;
  cursor: pointer;
  height: 0;
  width: 0;
}

/* Create a custom checkbox */
.checkmark {
  position: absolute;
  top: 0px;
  left: 0;
  height: 25px;
  width: 25px;
  background-color: #eee;
}

/* On mouse-over, add a grey background color */
.contcb:hover input ~ .checkmark {
  background-color: #ccc;
}

/* When the checkbox is checked, add a blue background */
.contcb input:checked ~ .checkmark {
  @apply bg-green-700;
}

/* Create the checkmark/indicator (hidden when not checked) */
.checkmark:after {
  content: "";
  position: absolute;
  display: none;
}

/* Show the checkmark when checked */
.contcb input:checked ~ .checkmark:after {
  display: block;
}

/* Style the checkmark/indicator */
.contcb .checkmark:after {
  left: 10px;
  top: 7px;
  width: 5px;
  height: 10px;
  border: solid white;
  border-width: 0 3px 3px 0;
  -webkit-transform: rotate(45deg);
  -ms-transform: rotate(45deg);
  transform: rotate(45deg);
}
</style>

<script lang="js">
import VueJsonPretty from "vue-json-pretty";
import "vue-json-pretty/lib/styles.css";
import { useUserStore } from '~/stores/user'

const user = useUserStore()

export default {
  components: {
    VueJsonPretty,
  },
  data() {
    return {
      jsonData: { msg: "Logs will be here" },
      name: ref(""),
      fields: ref(""),
      regexq: ref(""),
      viteHostname: import.meta.env.VITE_HOSTNAME,
      timestamp: "",
      mainbtntext: "Search",
      showlogs: "invisible",
      showhelp: "block",
      truncatelog: "",
      respFeedback: "",
      respFeedbackColor: "text-red-600",
      pollIntervalId: 1,
      blockform: false,
      polltry: 0
    };
  },
  created() {
    this.getNow();
    setInterval(this.getNow, 5000);
  },
  methods: {
    pollingDone() {
      this.blockform = false
      clearInterval(this.pollIntervalId);
      this.mainbtntext = "Search";
    },
    async poll() {
      this.respFeedbackColor = "text-light-500"
      this.showlogs = "invisible"
      fetch(
        import.meta.env.VITE_PROTOCOL +
        "://" +
        import.meta.env.VITE_HOSTNAME +
        "/api/porter/query/poll?qid=" +
        localStorage.getItem('qid') +
        "&token=" +
        user.token
      )
        .then((response) => {
          if (response.status == 401) {
            this.respFeedback = "Not authenticated. Visit from tg chat or sign in."
            this.respFeedbackColor = "text-red-500"
            this.pollingDone();
            return {};
          } else if (response.status == 403) {
            this.respFeedback = "Invalid token. Visit this link from tg chat."
            this.respFeedbackColor = "text-red-500"
            this.pollingDone();
            return {};
          } else if (response.status == 404) {
            this.respFeedback = "Query request id not found."
            this.respFeedbackColor = "text-red-500"
            this.pollingDone();
            return {};
          } else if (response.status == 500) {
            this.respFeedback = "Internal server error."
            this.respFeedbackColor = "text-red-500"
            this.pollingDone();
            return {};
          } else if (response.status != 200) {
            this.respFeedback = "Failed to fetch data from server. Please try again later."
            this.respFeedbackColor = "text-red-500"
            this.pollingDone();
            return {};
          }
          return response.json()
        })
        .then((data) => {
          if ('msg' in data) {
            this.respFeedback = data.msg
            if (data.status == 0) {
              if (this.polltry == 0) {
                setTimeout(() => this.poll, 100)
                this.polltry++
              } else if (this.polltry == 1) {
                setTimeout(() => this.poll, 200)
                this.polltry++
              }
              this.respFeedbackColor = "text-light-500"
            } else if (data.status == 1) {
              this.pollingDone();
              this.showhelp = "hidden"
              this.respFeedbackColor = "text-light-500"
              this.jsonData = data.result
              this.showlogs = "visible"
            } else if (data.status == 2) {
              this.pollingDone();
              this.respFeedbackColor = "text-red-500"
              this.showlogs = "invisible"
              this.showhelp = "block"
            };
          }
        })
    },
    async go() {
      if (this.blockform == true) {
        return {};
      }
      this.polltry = 0
      console.log("new query request")
      this.blockform = true
      this.showlogs = "invisible"
      var urlToFetch = import.meta.env.VITE_PROTOCOL +
        "://" +
        import.meta.env.VITE_HOSTNAME +
        "/api/porter/query/range?query=" +
        this.name +
        "&token=" +
        user.token
      if (this.fields != "") {
        urlToFetch = urlToFetch + "&fields=" + this.fields
      }
      if (this.regexq != "") {
        urlToFetch = urlToFetch + "&regex=" + this.regexq
      }
      fetch(
        urlToFetch
      )
        .then((response) => {
          if (response.status == 401) {
            this.respFeedback = "Not authenticated. Visit from tg chat or sign in."
            this.respFeedbackColor = "text-red-500"
            this.blockform = false
            return {};
          } else if (response.status == 403) {
            this.respFeedback = "Invalid token. Visit this link from tg chat."
            this.respFeedbackColor = "text-red-500"
            this.blockform = false
            return {};
          } else if (response.status == 400) {
            this.respFeedback = "Bad query."
            this.respFeedbackColor = "text-red-500"
            this.blockform = false
            return {}
          } else if (response.status == 500) {
            this.respFeedback = "Internal server error."
            this.respFeedbackColor = "text-red-500"
            this.blockform = false
            return {}
          } else if (response.status != 200) {
            this.respFeedback = "Failed to fetch data from server. Please try again later."
            this.respFeedbackColor = "text-red-500"
            this.blockform = false
            return {};
          }
          this.respFeedback = "Request accepted."
          this.respFeedbackColor = "text-light-500"
          setTimeout(() => this.poll, 80)
          var pollIntervalId = setInterval(this.poll, 600);
          this.pollIntervalId = pollIntervalId;
          return response.json()
        })
        .then((data) => {
          localStorage.setItem("qid", data.qid);
        })
    },
    getNow: function () {
      const today = new Date();
      const time = today.getUTCHours() + ":" + String(today.getUTCMinutes()).padStart(2, '0');
      this.timestamp = time;
    },
    truncLogFunc() {
      if (this.truncatelog == "truncate") {
        this.truncatelog = ""
      } else {
        this.truncatelog = "truncate"
      }
    }
  }
};
</script>
