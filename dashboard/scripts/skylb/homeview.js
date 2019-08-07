/**
 * @fileoverview The HomeView class.
 */

goog.provide('skylb.HomeView');

goog.require('goog.array');
goog.require('goog.dom');
goog.require('goog.dom.dataset');
goog.require('goog.events');
goog.require('goog.net.EventType');
goog.require('goog.net.XhrIo');
goog.require('proto.proto.GetAllServicesRequest');
goog.require('proto.proto.GetAllServicesResponse');
goog.require('skylb.View');
goog.require('skylb.home_html.HomeTemplate');


/**
 * Constructs a HomeView class.
 *
 * @param {skylb.XhrMgr} xhrMgr The Xhr manager.
 * @param {skylb.ViewMgr} viewMgr The view manager.
 * @constructor
 * @extends {skylb.View}
 */
skylb.HomeView = function(xhrMgr, viewMgr) {
  goog.base(this, '', 'Home', /^$/, xhrMgr, viewMgr);
};
goog.inherits(skylb.HomeView, skylb.View);


/** @override */
skylb.HomeView.prototype.getUrl = function() {
  return '/';
};


/** @override */
skylb.HomeView.prototype.load = function(callback) {
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

    var resp = proto.proto.GetAllServicesResponse.deserializeBinary(data);
    var services = resp.getServicesList();

    var t = new skylb.home_html.HomeTemplate();
    t.render(container, {'services': services});

    var trs = goog.dom.getElementsByClass('service-row', container);
    goog.array.forEach(trs, function(tr, idx) {
      goog.events.listen(tr, goog.events.EventType.CLICK, function(e) {
        var serviceId = goog.dom.dataset.get(tr, 'id');
        this.viewMgr.loadViewByName('service', /* pushState*/ true,
            { "service_id": serviceId });
      }, undefined, this);
    }, this);
  }, undefined, this);

  xhr.send('/_/get-all-services', 'POST', bytes);

  if (callback) {
    callback();
  }
};
