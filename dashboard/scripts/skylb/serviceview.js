/**
 * @fileoverview The ServiceView class.
 */

goog.provide('skylb.ServiceView');

goog.require('goog.array');
goog.require('goog.dom');
goog.require('goog.dom.classlist');
goog.require('goog.dom.dataset');
goog.require('goog.events');
goog.require('goog.events.EventType');
goog.require('goog.net.EventType');
goog.require('goog.net.XhrIo');
goog.require('goog.object');
goog.require('goog.style');
goog.require('goog.ui.Dialog');
goog.require('goog.ui.LabelInput');
goog.require('skylb.View');
goog.require('skylb.service_html.AddInstanceTemplate');
goog.require('skylb.service_html.InstancesViewerTemplate');
goog.require('skylb.service_html.ServiceTemplate');
goog.require('skylb.service_html.ToggleLameduckTemplate');


/**
 * Constructs a ServiceView class.
 *
 * @param {skylb.XhrMgr} xhrMgr The Xhr manager.
 * @param {skylb.ViewMgr} viewMgr The view manager.
 * @constructor
 * @extends {skylb.View}
 */
skylb.ServiceView = function(xhrMgr, viewMgr) {
  goog.base(this, 'service', 'Service', /^\/service\/(\d+)$/, xhrMgr, viewMgr);

  /** @private{number} */
  this.serviceId_ = null;

  /** @private{Element} */
  this.graphViewer_ = null;

  /** @private{Element} */
  this.svcViewer_ = null;
};
goog.inherits(skylb.ServiceView, skylb.View);

/** @override */
skylb.ServiceView.prototype.getUrl = function() {
  return '/service/'+this.getServiceId_();
};

/**
 * @return {number} the service ID.
 * @private
 */
skylb.ServiceView.prototype.getServiceId_ = function() {
  var serviceId = null;
  if (this.match && this.match.length >= 2) {
    serviceId = this.match[1];
  } else if (this.data) {
    serviceId = this.data['service_id'];
  }
  return serviceId;
};


/** @override */
skylb.ServiceView.prototype.load = function(callback) {
  var serviceId = this.getServiceId_();

  /** @type{Element} */
  var container = this.viewMgr.getContainer();
  goog.dom.removeChildren(container);

  var xhr = this.xhrMgr.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);
  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var data = xhr.getResponse();
    this.xhrMgr.returnXhr(xhr);

    var resp = proto.proto.GetServiceByIdResponse.deserializeBinary(data);
    var t = new skylb.service_html.ServiceTemplate();
    var service = resp.getService();
    t.render(container, {'service': service});

    this.graphViewer_ = goog.dom.getElementByClass('graph-viewer', container);
    this.svcViewer_ = goog.dom.getElementByClass('service-viewer', container);
    this.updateViewSize();

    var nodes = [{
      'id': service.getId(),
      'group': 'target-service',
      'label': service.getName()
    }];
    var edges = [];
    var grpcServices = {};
    goog.array.forEach(service.getIncomingsList(), function(cli) {
      if (!grpcServices[cli.getId()]) {
        grpcServices[cli.getId()] = true;
        nodes.push({
          'id': cli.getId(), 'group':'incomings', 'label': cli.getName()
        });
      }
      edges.push({'from': service.getId(), 'to': cli.getId(), 'arrows':'from'});
    });
    goog.array.forEach(service.getOutgoingsList(), function(out) {
      if (!grpcServices[out.getId()]) {
        grpcServices[out.getId()] = true;
        nodes.push({
          'id': out.getId(), 'group':'outgoings', 'label': out.getName()
        });
      }
      edges.push({'from': service.getId(), 'to': out.getId(), 'arrows':'to'});
    });

    var networkData = {
      'nodes': new vis.DataSet(nodes),
      'edges': new vis.DataSet(edges)
    };
    var options = {};
    var network = new vis.Network(this.graphViewer_, networkData, options);
    network.setOptions({
      'nodes': {
        'shape': 'box',
        'physics': false,
        'font': {
          'face': 'monospace',
          'size': 18
        }
      },
      'groups': {
        'useDefaultGroups': false,
        'target-service': {
          'color':{
            'background':'#D2E5FF',
            'border': '#2B7CE9',
            'highlight': {
              'border': '#2B7CE9',
              'background': '#D2E5FF'
            }
          }
        },
        'incomings': {
          'color':{
            'background':'#FFFF00',
            'border': '#FFA804',
            'highlight': {
              'border': '#FFA500',
              'background': '#FFFFA3'
            }
          }
        },
        'outgoings': {
          'color':{
            'background':'#FB7E81',
            'border': '#FB494D',
            'highlight': {
              'border': '#FA0A10',
              'background': '#FFAFB1'
            }
          }
        }
      },
    });
    network.on('doubleClick', goog.bind(function(params) {
      if (params['nodes'] && params['nodes'][0] &&
          params['nodes'][0] != service.getId()) {
        this.viewMgr.loadViewByName('service', /* pushState*/ true,
            { "service_id": params['nodes'][0] });
      }
    }, this));

    this.installInstViewerListeners_(this.svcViewer_, service.getId());
  }, undefined, this);

  var req = new proto.proto.GetServiceByIdRequest();
  req.setId(serviceId);
  var bytes = req.serializeBinary();
  xhr.send('/_/get-service-by-id', 'POST', bytes);

  if (callback) {
    callback();
  }
};


/** @override */
skylb.ServiceView.prototype.updateViewSize = function() {
  if (this.graphViewer_) {
    var size = this.viewMgr.getViewportSize();
    goog.style.setHeight(this.graphViewer_, size.height * 0.9 - 50);
    goog.style.setWidth(this.graphViewer_, size.width - 400);
  }
};


/**
 * @param {number} serviceId the service ID.
 * @private
 */
skylb.ServiceView.prototype.launchAddInstDialog_ = function(serviceId) {
  var dialog = new goog.ui.Dialog("add-inst-dialog");
  dialog.setDisposeOnHide(true);
  dialog.setModal(true);
  dialog.setHasTitleCloseButton(true);
  dialog.setButtonSet(null);
  dialog.setTitle('Add a new lameduck service instance');

  var container = dialog.getContentElement();

  var t = new skylb.service_html.AddInstanceTemplate();
  t.render(container, {});

  var addrInput = new goog.ui.LabelInput(
      'Service instance address as <host>:<port>');
  addrInput.render(goog.dom.getElementByClass('address-input', container));

  var saveBtn = goog.dom.getElementByClass('save-button', container);
  goog.events.listen(saveBtn, goog.events.EventType.CLICK, function(e) {
    var errEle = goog.dom.getElementByClass('errors', container);
    this.addInstance_(dialog, addrInput, errEle, serviceId);
  }, undefined, this);

  dialog.setVisible(true);
};


/**
 * @param {goog.ui.Dialog} dialog the dialog.
 * @param {goog.ui.LabelInput} input the address input.
 * @param {Element} errEle the error msg element.
 * @param {number} serviceId the service ID.
 * @private
 */
skylb.ServiceView.prototype.addInstance_ =
    function(dialog, input, errEle, serviceId) {
  var addr = input.getValue();
  if (addr.length < 4) {
    goog.dom.setTextContent(errEle, 'Wrong format, should be <host>:<port>');
    goog.dom.classlist.remove(errEle, 'invisible');
    return;
  }

  var xhr = this.xhrMgr.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);
  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var data = xhr.getResponse();
    this.xhrMgr.returnXhr(xhr);

    var resp = proto.proto.AddInstanceResponse.deserializeBinary(data);
    var msg = resp.getErrorMsg();
    if (msg) {
      goog.dom.setTextContent(errEle, msg);
      goog.dom.classlist.remove(errEle, 'invisible');
      return;
    }
    dialog.setVisible(false);
    this.refreshInstances_(serviceId);
  }, undefined, this);

  goog.dom.classlist.add(errEle, 'invisible');
  var req = new proto.proto.AddInstanceRequest();
  req.setId(serviceId);
  req.setAddress(addr);
  var bytes = req.serializeBinary();
  xhr.send('/_/add-instance', 'POST', bytes);
};


/**
 * @param {number} serviceId the service ID.
 * @private
 */
skylb.ServiceView.prototype.refreshInstances_ = function(serviceId) {
  goog.dom.removeChildren(this.svcViewer_);

  var xhr = this.xhrMgr.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);
  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var data = xhr.getResponse();
    this.xhrMgr.returnXhr(xhr);

    var resp = proto.proto.GetServiceByIdResponse.deserializeBinary(data);
    var t = new skylb.service_html.InstancesViewerTemplate();
    var service = resp.getService();
    t.render(this.svcViewer_, {'service': service});

    this.installInstViewerListeners_(this.svcViewer_, service.getId());
  }, undefined, this);

  var req = new proto.proto.GetServiceByIdRequest();
  req.setId(serviceId);
  var bytes = req.serializeBinary();
  xhr.send('/_/get-service-by-id', 'POST', bytes);
};

/**
 * @param {Element} container the service instance viewer element.
 * @param {number} serviceId the service ID.
 * @private
 */
skylb.ServiceView.prototype.installInstViewerListeners_ =
    function(container, serviceId) {
  var btnAddInst = goog.dom.getElementByClass('btn-add-instance', container);
  goog.events.listen(
      btnAddInst,
      goog.events.EventType.CLICK, function(e){
        this.launchAddInstDialog_(serviceId);
      }, undefined, this);

  var btnRefresh = goog.dom.getElementByClass('btn-refresh', container);
  goog.events.listen(
      btnRefresh,
      goog.events.EventType.CLICK, function(e){
        this.refreshInstances_(serviceId);
      }, undefined, this);

  var images = goog.dom.getElementsByTagNameAndClass('img', 'icon', container);
  goog.array.forEach(images, function(img) {
    goog.events.listen(img, goog.events.EventType.CLICK, function(e){
      var addr = goog.dom.dataset.get(img, 'serviceAddress');
      var isLameduck = goog.dom.dataset.get(img, 'serviceIslameduck');
      this.launchToggleLameduckDialog_(serviceId, addr, isLameduck);
    }, undefined, this);
  }, this);
};

/**
 * @param {number} serviceId the service ID.
 * @param {string} addr the instance address.
 * @param {boolean} isLameduck whether it's in lameduck mode at present.
 * @private
 */
skylb.ServiceView.prototype.launchToggleLameduckDialog_ =
    function(serviceId, addr, isLameduck) {
  var dialog = new goog.ui.Dialog("toggle-lameduck-dialog");
  dialog.setDisposeOnHide(true);
  dialog.setModal(true);
  dialog.setHasTitleCloseButton(true);
  dialog.setButtonSet(null);
  dialog.setTitle('Toggle lameduck service instance');

  var container = dialog.getContentElement();

  console.log(serviceId, addr, isLameduck);
  var t = new skylb.service_html.ToggleLameduckTemplate();
  t.render(container, {
    'service_id': serviceId,
    'service_address': addr,
    'is_lameduck': isLameduck,
  });

  var saveBtn = goog.dom.getElementByClass('save-button', container);
  goog.events.listen(saveBtn, goog.events.EventType.CLICK, function(e) {
    var errEle = goog.dom.getElementByClass('errors', container);
    var checkbox = goog.dom.getElementByClass('lameduck-input', container);
    this.toggleLameduckMode_(
        dialog, addr, serviceId, errEle, checkbox.checked);
  }, undefined, this);

  dialog.setVisible(true);
};


/**
 * @param {goog.ui.Dialog} dialog the dialog.
 * @param {string} addr the instance address.
 * @param {number} serviceId the service ID.
 * @param {Element} errEle the error msg element.
 * @param {boolean} isLameduck whether the instance is in duck mode.
 * @private
 */
skylb.ServiceView.prototype.toggleLameduckMode_ =
    function(dialog, addr, serviceId, errEle, isLameduck) {
  console.log("here", this);
  var xhr = this.xhrMgr.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);
  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var data = xhr.getResponse();
    this.xhrMgr.returnXhr(xhr);

    var resp = proto.proto.ToggleLameduckResponse.deserializeBinary(data);
    var msg = resp.getErrorMsg();
    if (msg) {
      goog.dom.setTextContent(errEle, msg);
      goog.dom.classlist.remove(errEle, 'invisible');
      return;
    }
    dialog.setVisible(false);
    this.refreshInstances_(serviceId);
  }, undefined, this);

  goog.dom.classlist.add(errEle, 'invisible');
  var req = new proto.proto.ToggleLameduckRequest();
  req.setId(serviceId);
  req.setAddress(addr);
  req.setLameduck(isLameduck);
  var bytes = req.serializeBinary();
  xhr.send('/_/toggle-lameduck', 'POST', bytes);
};
