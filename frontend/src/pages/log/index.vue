<template>
  <div id="formdiv">
    <div class="gap-4 grid grid-rows-1 grid-cols-3 m-auto pb-4 h-14" style="max-width: 350px; margin-top: -10px">
      <!-- TODO menu -->
      <button disabled class="btn col-span-1 row-span-1  rounded">Status</button>
      <button disabled class="btn col-span-1 row-span-1  rounded">History</button>
      <button disabled class="btn col-span-1 row-span-1 rounded">Help</button>
    </div>
    <div class="m-auto mb-2 mt-5" w="full" text="left" style="max-width: 500px">
      <code>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
        <span
          class="text-gray-400"
        >{{ timestamp }} UTC</span>
        <br />
        &nbsp;&nbsp;env service dd [hh:mm:ss]
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
      bg="dark-700"
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
      bg="dark-700"
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
      bg="dark-700"
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
      <p class="col-span-2 row-span-1"></p>
      <input
        class="col-span-2 row-span-1"
        id="tokeninput"
        v-model="powertoken"
        placeholder="internal options"
        type="text"
        autocomplete="false"
        p="x-4 y-2"
        w="full"
        text="left"
        bg="dark-700"
        border="~ rounded gray-700"
        outline="none active:none"
        @keydown.enter="go"
        style="max-width: 500px"
      />

      <div class="col-span-1 row-span-1">
        <button
          w="full"
          style="max-width: 310px"
          class="p-2 text-mm btn truncate"
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
  <div :class="showhelp" id="helpdiv" class="mt-6 pr-8 pb-10">
    <h3 class="text-lg mb-4">Main input</h3>
    <ul text="left" class="list-square list-inside text-sm">
      <li class="pb-2">
        <code class="font-extrabold">env</code> - hostname string;
        can be a substring that uniquely identifies hostname;
      </li>
      <li class="pb-2">
        <code class="font-extrabold">dd hh:mm:ss</code>
        <span>
          - only
          <code>dd</code> part is mandatory (can be "t" for today); if some part after that
          is not specified, then range is taken;
        </span>
      </li>
      <!-- <li class="pb-2"> -->
        <!-- <code class="font-extrabold">##m</code> -->
        <!-- <span>- alternative to the above, minutes since now;</span> -->
      <!-- </li> -->
      <li class="pb-2">
        <code class="font-extrabold">service</code> - service name;
        can be a substring that uniquely identifies service;
      </li>
      <!-- <li class="pb-2">
        <code class="font-extrabold">flag</code> - optional, any flag
        (it cannot be user-defined as of now and only single flag on event is possible);
        if not specified records with any flags are included; pre-defined flags are
        <code>err</code> and
        <code>none</code>.
      </li>-->
    </ul>
    <p
      text="left"
      class="text-sm"
    >Any valid query using only main input has time complexity of O(n), where n is input events' rate.</p>
    <h3 class="text-lg my-4">Filter by one field with regex</h3>
    <p text="left" class="text-sm">
      All result set can be filtered by one field's value using a regular expression.
      This rule is applied to the set obtained with main rule.
      Any
      <a
        href="https://github.com/google/re2/wiki/Syntax"
        class="underline"
      >RE2</a> syntax is accepted.
      Time complexity is O(n), where n is either
      the size of the base set or the length of the field's value in question (the latter is guaranteed by
      go's regexp package -
      <a
        href="https://swtch.com/~rsc/regexp/regexp1.html"
        class="underline"
      >more on that here</a>).
    </p>
    <h3 class="text-lg my-4">Include only selected fields</h3>
    <p text="left" class="text-sm">
      Only selected fields (separated by space) are present in result.
      They can be nested to arbitrary depth. Applied to the set
      obtained with above rules. It's complexity is O(n) (should be O(1), but it does not
      remember the schema of any source and unmarshalls keys dynamically for now).
    </p>
  </div>
  <div :class="showlogs" id="logsdiv">
    <vue-json-pretty
      class="leading-none break-all"
      :class="truncatelog"
      style="max-width: 1500px; font-size: 12px !important;"
      :path="'res'"
      :data="jsonData"
      :deep="depth"
    ></vue-json-pretty>
  </div>
</template>

<style>
@media (min-width: 1250px) {
  #helpdiv {
    position: absolute;
    right: 0px;
    max-width: 60%;
  }

  #logsdiv {
    position: absolute;
    right: 0px;
    max-width: 62%;
    min-width: 62%;
  }

  #formdiv {
    position: fixed;
    max-width: 500px;
    min-width: 33%;
  }
}

@media (min-width: 1000px) and (max-width: 1250px) {
  #helpdiv {
    position: absolute;
    right: 0px;
    max-width: 62%;
  }
  #logsdiv {
    position: absolute;
    right: 0px;
    max-width: 62%;
    min-width: 60%;
  }

  #formdiv {
    position: fixed;
    max-width: 33%;
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
      powertoken: ref(""),
      viteHostname: import.meta.env.VITE_HOSTNAME,
      depth: 1,
      timestamp: "",
      mainbtntext: "Search",
      showlogs: "invisible",
      showhelp: "block",
      truncatelog: "",
      respFeedback: "",
      respFeedbackColor: "text-red-600",
      pollIntervalId: 1,
      blockform: false
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
      this.respFeedback = "Polling..."
      this.respFeedbackColor = "text-light-500"
      this.showlogs = "invisible"
      fetch(
        import.meta.env.VITE_PROTOCOL +
        "://" +
        import.meta.env.VITE_HOSTNAME +
        "/api/log/poll?qid=" +
        localStorage.getItem('qid') +
        "&token=" +
        this.$route.query.token
      )
        .then((response) => {
          if (response.status == 401) {
            this.respFeedback = "No token provided. Visit this link from tg chat."
            this.respFeedbackColor = "text-red-500"
            this.pollingDone();
            return {};
          } else if (response.status == 403) {
            this.respFeedback = "Invalid token. Visit this link from tg chat."
            this.respFeedbackColor = "text-red-500"
            this.pollingDone();
            return {};
          } else if (response.status == 400) {
            this.respFeedback = "Bad job id. Please resend your query."
            this.respFeedbackColor = "text-red-500"
            this.pollingDone();
            return {};
          } else if (response.status == 200) {
            this.respFeedback = "Processing..."
            this.respFeedbackColor = "text-light-500"
            return {}
          } else if (response.status == 202) {
            this.respFeedback = "Queued..."
            this.respFeedbackColor = "text-light-500"
            return {}
          } else if (response.status == 404) {
            this.respFeedback = "No logs found for this query."
            this.respFeedbackColor = "text-yellow-500"
            this.pollingDone();
            return {};
          } else if (response.status != 201) {
            this.respFeedback = "Failed to fetch data from server. Please try again later."
            this.respFeedbackColor = "text-red-500"
            this.pollingDone();
            return {};
          }
          this.pollingDone();
          this.respFeedback = "Rendering..."
          this.respFeedbackColor = "text-light-500"
          this.showhelp = "hidden"
          return response.json()
        })
        .then((data) => {
          if (length in data) {
            if (data.length == 100) {
              this.showlogs = "visible"
              this.respFeedback = "Too many matching events. Only 100 most recent are shown."
              this.respFeedbackColor = "text-yellow-500"
            } else {
              this.showlogs = "visible"
              this.respFeedback = data.length + " matching events found."
              this.respFeedbackColor = "text-green-600"
            }
          } else if (this.respFeedback == "Rendering...") {
            this.showlogs = "visible"
            this.respFeedback == ""
          }
          this.jsonData = data;
        })
    },
    BreakSignal() { },
    async go() {
      if (this.blockform == true) {
        return {};
      }
      console.log("new query request")
      this.blockform = true
      this.mainbtntext = "Requesting..."
      this.showlogs = "invisible"
      var urlToFetch = import.meta.env.VITE_PROTOCOL +
        "://" +
        import.meta.env.VITE_HOSTNAME +
        "/api/log/range?query=" +
        this.name +
        "&token=" +
        this.$route.query.token
      if (this.fields != "") {
        urlToFetch = urlToFetch + "&fields=" + this.fields
      }
      if (this.regexq != "") {
        urlToFetch = urlToFetch + "&regex=" + this.regexq
      }
      if (this.powertoken != "") {
        urlToFetch = urlToFetch + "&powertoken=" + this.powertoken
      }
      fetch(
        urlToFetch
      )
        .then((response) => {
          if (response.status == 401) {
            this.respFeedback = "No token provided. Visit this link from tg chat."
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
          } else if (response.status != 202) {
            this.respFeedback = "Failed to fetch data from server. Please try again later."
            this.respFeedbackColor = "text-red-500"
            this.blockform = false
            return {};
          }
          this.respFeedback = "Query job queued."
          this.respFeedbackColor = "text-light-500"
          setTimeout(() => this.poll, 100)
          var pollIntervalId = setInterval(this.poll, 500);
          this.pollIntervalId = pollIntervalId;
          return response.json()
        })
        .then((data) => {
          localStorage.setItem("qid", data.qid);
        })
        .then(() => {
          this.mainbtntext = "Please Wait";
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
