<template>
  <div class="valvectl">
    <h1>Valve Controller</h1>
    <div>
      <button @click="switchValve()">{{ btnMsg }}</button>
      <br />
      <br />
      <label>{{ msg }}</label>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { OperationResult, ValvedSvc } from "@/services/valvesvc";

@Options({
  props: {},
})
export default class ValveCtl extends Vue {
  private svcClient = new ValvedSvc();
  private msg = "";
  private status = false;

  created(): void {
    const callback = (status: boolean, error: string | undefined) => {
      if (error) {
        this.msg = error;
        return;
      }

      this.status = status;
    };
    this.svcClient.Status(callback);
  }

  get btnMsg(): string {
    return this.status ? "Switch Off" : "Switch On";
  }

  private switchValve() {
    const callback = (r: OperationResult | undefined) => {
      if (r?.Message) {
        this.status = !this.status;
      }

      this.msg = r?.Error ?? r?.Message ?? "Error: No response from server";
    };

    if (this.status) return this.svcClient.Close(callback);
    return this.svcClient.Open(callback);
  }
}
</script>
