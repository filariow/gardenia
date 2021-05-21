<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png" />
    <ValveCtl />
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import ValveCtl from "@/components/ValveCtl.vue"; // @ is an alias to /src

@Options({
  components: {
    ValveCtl,
  },
})
export default class Home extends Vue {}
</script>
