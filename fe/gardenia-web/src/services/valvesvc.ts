import {
  CloseValveReply,
  CloseValveRequest,
  OpenValveReply,
  OpenValveRequest,
  StatusValveReply,
  StatusValveRequest,
} from "@/grpc/valved_pb";
import { ServiceError, ValvedSvcClient } from "@/grpc/valved_pb_service";

export interface OperationResult {
  Message?: string;
  Error?: string;
}

export class ValvedSvc {
  private svcClient: ValvedSvcClient;

  constructor(address = "http://localhost:14000") {
    this.svcClient = new ValvedSvcClient(address);
  }

  Open(callback: (response: OperationResult | undefined) => void) {
    const request = new OpenValveRequest();
    this.svcClient.open(
      request,
      (err: ServiceError | null, r: OpenValveReply | null) => {
        console.log(r?.getMessage());
        const response = { Message: r?.getMessage(), Error: err?.message };
        callback(response);
      }
    );
  }

  Close(callback: (response: OperationResult | undefined) => void) {
    const request = new CloseValveRequest();
    this.svcClient.close(
      request,
      (err: ServiceError | null, r: CloseValveReply | null) => {
        console.log(r?.getMessage());
        const response = { Message: r?.getMessage(), Error: err?.message };
        callback(response);
      }
    );
  }

  Status(callback: (status: boolean, error: string | undefined) => void) {
    const request = new StatusValveRequest();
    this.svcClient.status(
      request,
      (err: ServiceError | null, r: StatusValveReply | null) => {
        console.log(r?.getStatus());
        const st = r?.getStatus() == 0;
        callback(st, err?.message);
      }
    );
  }
}
