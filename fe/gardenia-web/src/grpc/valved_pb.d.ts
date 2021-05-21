// package: valvedgrpc
// file: valved.proto

import * as jspb from "google-protobuf";

export class OpenValveRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpenValveRequest.AsObject;
  static toObject(includeInstance: boolean, msg: OpenValveRequest): OpenValveRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OpenValveRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpenValveRequest;
  static deserializeBinaryFromReader(message: OpenValveRequest, reader: jspb.BinaryReader): OpenValveRequest;
}

export namespace OpenValveRequest {
  export type AsObject = {
  }
}

export class OpenValveReply extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpenValveReply.AsObject;
  static toObject(includeInstance: boolean, msg: OpenValveReply): OpenValveReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OpenValveReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpenValveReply;
  static deserializeBinaryFromReader(message: OpenValveReply, reader: jspb.BinaryReader): OpenValveReply;
}

export namespace OpenValveReply {
  export type AsObject = {
    message: string,
  }
}

export class CloseValveRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseValveRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CloseValveRequest): CloseValveRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CloseValveRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseValveRequest;
  static deserializeBinaryFromReader(message: CloseValveRequest, reader: jspb.BinaryReader): CloseValveRequest;
}

export namespace CloseValveRequest {
  export type AsObject = {
  }
}

export class CloseValveReply extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseValveReply.AsObject;
  static toObject(includeInstance: boolean, msg: CloseValveReply): CloseValveReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CloseValveReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseValveReply;
  static deserializeBinaryFromReader(message: CloseValveReply, reader: jspb.BinaryReader): CloseValveReply;
}

export namespace CloseValveReply {
  export type AsObject = {
    message: string,
  }
}

export class StatusValveRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusValveRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StatusValveRequest): StatusValveRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StatusValveRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusValveRequest;
  static deserializeBinaryFromReader(message: StatusValveRequest, reader: jspb.BinaryReader): StatusValveRequest;
}

export namespace StatusValveRequest {
  export type AsObject = {
  }
}

export class StatusValveReply extends jspb.Message {
  getStatus(): ValveStatusMap[keyof ValveStatusMap];
  setStatus(value: ValveStatusMap[keyof ValveStatusMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusValveReply.AsObject;
  static toObject(includeInstance: boolean, msg: StatusValveReply): StatusValveReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StatusValveReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusValveReply;
  static deserializeBinaryFromReader(message: StatusValveReply, reader: jspb.BinaryReader): StatusValveReply;
}

export namespace StatusValveReply {
  export type AsObject = {
    status: ValveStatusMap[keyof ValveStatusMap],
  }
}

export interface ValveStatusMap {
  OPEN: 0;
  CLOSE: 1;
  INVALID: 2;
}

export const ValveStatus: ValveStatusMap;

