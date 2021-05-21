// package: valvedgrpc
// file: valved.proto

import * as valved_pb from "./valved_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ValvedSvcOpen = {
  readonly methodName: string;
  readonly service: typeof ValvedSvc;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof valved_pb.OpenValveRequest;
  readonly responseType: typeof valved_pb.OpenValveReply;
};

type ValvedSvcClose = {
  readonly methodName: string;
  readonly service: typeof ValvedSvc;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof valved_pb.CloseValveRequest;
  readonly responseType: typeof valved_pb.CloseValveReply;
};

type ValvedSvcStatus = {
  readonly methodName: string;
  readonly service: typeof ValvedSvc;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof valved_pb.StatusValveRequest;
  readonly responseType: typeof valved_pb.StatusValveReply;
};

export class ValvedSvc {
  static readonly serviceName: string;
  static readonly Open: ValvedSvcOpen;
  static readonly Close: ValvedSvcClose;
  static readonly Status: ValvedSvcStatus;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class ValvedSvcClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  open(
    requestMessage: valved_pb.OpenValveRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: valved_pb.OpenValveReply|null) => void
  ): UnaryResponse;
  open(
    requestMessage: valved_pb.OpenValveRequest,
    callback: (error: ServiceError|null, responseMessage: valved_pb.OpenValveReply|null) => void
  ): UnaryResponse;
  close(
    requestMessage: valved_pb.CloseValveRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: valved_pb.CloseValveReply|null) => void
  ): UnaryResponse;
  close(
    requestMessage: valved_pb.CloseValveRequest,
    callback: (error: ServiceError|null, responseMessage: valved_pb.CloseValveReply|null) => void
  ): UnaryResponse;
  status(
    requestMessage: valved_pb.StatusValveRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: valved_pb.StatusValveReply|null) => void
  ): UnaryResponse;
  status(
    requestMessage: valved_pb.StatusValveRequest,
    callback: (error: ServiceError|null, responseMessage: valved_pb.StatusValveReply|null) => void
  ): UnaryResponse;
}

