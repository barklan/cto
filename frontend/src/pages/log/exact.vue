<template>
    <div class="break-all">
        <vue-json-pretty :path="'res'" :data="jsonData"></vue-json-pretty>
    </div>
</template>

<script lang="js">
import VueJsonPretty from "vue-json-pretty";
import "vue-json-pretty/lib/styles.css";

export default {
    components: {
        VueJsonPretty,
    },
    data() {
        return {
            jsonData: { msg: "Throttling request... " },
        };
    },
    beforeMount() {
        fetch(
            import.meta.env.VITE_PROTOCOL +
            "://" +
            import.meta.env.VITE_HOSTNAME +
            "/api/porter/query/exact?key=" +
            this.$route.query.key
        )
            .then((response) => {
                if (response.status == 404) {
                    return { msg: "No logs found for this query." };
                } else if (response.status != 200) {
                    return { msg: "Failed to fetch data from server. Please try again later." };
                }
                return response.json()
            })
            .then((data) => {
                this.jsonData = data;
            });
    },
};
</script>

<style>
.vjs-tree__node.is-highlight,
.vjs-tree__node:hover {
  background-color: #262626 !important;
}
</style>
