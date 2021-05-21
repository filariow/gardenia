// package: valvedgrpc
// file: valved.proto

var valved_pb = require("./valved_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ValvedSvc = (function () {
  function ValvedSvc() {}
  ValvedSvc.serviceName = "valvedgrpc.ValvedSvc";
  return ValvedSvc;
}());

ValvedSvc.Open = {
  methodName: "Open",
  service: ValvedSvc,
  requestStream: false,
  responseStream: false,
  requestType: valved_pb.OpenValveRequest,
  responseType: valved_pb.OpenValveReply
};

ValvedSvc.Close = {
  methodName: "Close",
  service: ValvedSvc,
  requestStream: false,
  responseStream: false,
  requestType: valved_pb.CloseValveRequest,
  responseType: valved_pb.CloseValveReply
};

ValvedSvc.Status = {
  methodName: "Status",
  service: ValvedSvc,
  requestStream: false,
  responseStream: false,
  requestType: valved_pb.StatusValveRequest,
  responseType: valved_pb.StatusValveReply
};

exports.ValvedSvc = ValvedSvc;

function ValvedSvcClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ValvedSvcClient.prototype.open = function open(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ValvedSvc.Open, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ValvedSvcClient.prototype.close = function close(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ValvedSvc.Close, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ValvedSvcClient.prototype.status = function status(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ValvedSvc.Status, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.ValvedSvcClient = ValvedSvcClient;

