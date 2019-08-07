/**
 * @fileoverview The LogsView class.
 */

goog.provide('skylb.LogsView');

goog.require('goog.array');
goog.require('goog.dom');
goog.require('goog.dom.forms');
goog.require('goog.events');
goog.require('goog.net.EventType');
goog.require('goog.net.XhrIo.ResponseType');
goog.require('proto.proto.GetAllServicesRequest');
goog.require('proto.proto.GetAllServicesResponse');
goog.require('proto.proto.GetLogsRequest');
goog.require('proto.proto.GetLogsResponse');
goog.require('skylb.View');
goog.require('skylb.logs_html.LogsTemplate');


/**
 * Constructs a LogsView class.
 *
 * @param {skylb.XhrMgr} xhrMgr The Xhr manager.
 * @param {skylb.ViewMgr} viewMgr The view manager.
 * @constructor
 * @extends {skylb.View}
 */
skylb.LogsView = function(xhrMgr, viewMgr) {
  goog.base(this, 'logs', 'Logs Viewer', /^\/logs$/, xhrMgr, viewMgr);

  /** @private{proto.proto.GetAllServicesResponse} */
  this.services_ = null;
};
goog.inherits(skylb.LogsView, skylb.View);


/** @override */
skylb.LogsView.prototype.getUrl = function() {
  return '/logs';
};


/** @override */
skylb.LogsView.prototype.load = function(callback) {
  /** @type{Element} */
  var container = this.viewMgr.getContainer();
  goog.dom.removeChildren(container);

  var req = new proto.proto.GetAllServicesRequest();
  var bytes = req.serializeBinary();

  var xhr = this.xhrMgr.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);
  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var data = xhr.getResponse();
    this.xhrMgr.returnXhr(xhr);

    this.services_ = proto.proto.GetAllServicesResponse.deserializeBinary(data);
    this.loadLogs_(null);
  }, undefined, this);

  xhr.send('/_/get-all-services', 'POST', bytes);

  if (callback) {
    callback();
  }
};


/**
 * @param {proto.proto.GetLogsRequest|null} req the request.
 * @private
 */
skylb.LogsView.prototype.loadLogs_ = function(req) {
  /** @type{Element} */
  var container = this.viewMgr.getContainer();
  goog.dom.removeChildren(container);

  if (!req) {
    req = new proto.proto.GetLogsRequest();
    req.setServiceId(-1);
  }
  var bytes = req.serializeBinary();

  var xhr = this.xhrMgr.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);
  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var data = xhr.getResponse();
    this.xhrMgr.returnXhr(xhr);

    var resp = proto.proto.GetLogsResponse.deserializeBinary(data);
    var logs = resp.getLogsList();
    var options = {
      year: 'numeric',
      month: 'numeric',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
      timeZoneName: 'short',
      hour12: false
    };
    goog.array.forEach(logs, function(ele, idx) {
      var t = new Intl.DateTimeFormat('zh-Hans-CN-u-ca', options)
          .format(new Date(ele.getOpTime()));
      ele.setOpTime(t);
    });

    var t = new skylb.logs_html.LogsTemplate();
    t.render(container, {
      'operator': resp.getOperator(),
      'service_id': resp.getServiceId(),
      'logs': logs,
      'services': this.services_.getServicesList()
    });

    var btn = goog.dom.getElementByClass('btn-filter', container);
    goog.events.listen(btn, goog.events.EventType.CLICK, function(e) {
      var opEle = goog.dom.getElementByClass('filter-operator', container);
      var op = goog.dom.forms.getValue(opEle);
      var svcEle = goog.dom.getElementByClass('filter-service', container);
      var svc = goog.dom.forms.getValue(svcEle);
      req = new proto.proto.GetLogsRequest();
      req.setOperator(op);
      req.setServiceId(svc);
      this.loadLogs_(req);
    }, undefined, this);
  }, undefined, this);

  xhr.send('/_/get-logs', 'POST', bytes);
};
