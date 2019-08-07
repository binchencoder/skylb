/**
 * @fileoverview The main app of skylb.
 */
goog.provide('skylb.main');

goog.require('goog.dom');
goog.require('goog.dom.ViewportSizeMonitor');
goog.require('goog.dom.classlist');
goog.require('goog.events');
goog.require('goog.events.EventType');
goog.require('goog.net.EventType');
goog.require('goog.net.XhrIo');
goog.require('proto.proto.GetCurrentUserRequest');
goog.require('proto.proto.GetCurrentUserResponse');
goog.require('skylb.HomeView');
goog.require('skylb.LogsView');
goog.require('skylb.ServiceView');
goog.require('skylb.UsersManagementView');
goog.require('skylb.ViewMgr');
goog.require('skylb.XhrMgr');


/**
 * The main app.
 * @export
 */
skylb.main = function() {
  /** @type{goog.dom.ViewportSizeMonitor} */
  var vsm = new goog.dom.ViewportSizeMonitor();

  /** @type {skylb.XhrMgr} */
  var xhrMgr = new skylb.XhrMgr();

  /** @type {skylb.ViewMgr} */
  var viewMgr = new skylb.ViewMgr(vsm, goog.dom.getElement('view-container'));

  skylb.registerViews(xhrMgr, viewMgr);
  skylb.installEventHandlers(viewMgr);

  var req = new proto.proto.GetCurrentUserRequest();
  var payload = req.serializeBinary();
  var xhr = xhrMgr.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);
  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var data = xhr.getResponse();
    xhrMgr.returnXhr(xhr);

    var resp = proto.proto.GetCurrentUserResponse.deserializeBinary(data);
    if (resp.getErrorMsg()) {
      alert("Failed to get current user, please reload and retry.");
      return;
    }

    viewMgr.setViewer(resp.getUser());
    viewMgr.loadViewByUrl(window.location, /* pushState */ false);
  });
  xhr.send('/_/get-current-user', 'POST', payload);
};


/**
 * Registers all views.
 *
 * @param {skylb.XhrMgr} xhrMgr The Xhr manager.
 * @param {skylb.ViewMgr} viewMgr The view manager.
 */
skylb.registerViews = function(xhrMgr, viewMgr) {
  viewMgr.registerView(new skylb.HomeView(xhrMgr, viewMgr));
  viewMgr.registerView(new skylb.LogsView(xhrMgr, viewMgr));
  viewMgr.registerView(new skylb.ServiceView(xhrMgr, viewMgr));
  viewMgr.registerView(new skylb.UsersManagementView(xhrMgr, viewMgr));
};


/**
 * Installs event handlers.
 *
 * @param {skylb.ViewMgr} viewMgr The view manager.
 */
skylb.installEventHandlers = function(viewMgr) {
  var topPanel = goog.dom.getElement('top-panel');

  var logo = goog.dom.getElementsByTagNameAndClass('img', 'logo', topPanel);
  if (logo && logo.length > 0) {
    goog.events.listen(logo[0], goog.events.EventType.CLICK, function() {
      viewMgr.loadViewByName('', /* pushState */ true, {});
    });
  }

  var users = goog.dom.getElementsByTagNameAndClass(
      'img', 'skylb-users', topPanel);
  if (users && users.length > 0) {
    goog.events.listen(users[0], goog.events.EventType.CLICK, function() {
      viewMgr.loadViewByName('users', /* pushState */ true, {});
    });
  }

  var logs = goog.dom.getElementsByTagNameAndClass(
      'img', 'skylb-logs', topPanel);
  if (logs && logs.length > 0) {
    goog.events.listen(logs[0], goog.events.EventType.CLICK, function() {
      viewMgr.loadViewByName('logs', /* pushState */ true, {});
    });
  }

  var meDropdown = goog.dom.getElement('me-dropdown');
  goog.events.listen(goog.dom.getElement('loginUser'),
      goog.events.EventType.CLICK, function() {
        goog.dom.classlist.toggle(meDropdown, 'hidden');
      });

  var btnLogout = goog.dom.getElementByClass('btn-logout', meDropdown);
  goog.events.listen(btnLogout, goog.events.EventType.CLICK, function() {
    goog.dom.classlist.add(meDropdown, 'hidden');
    window.location = "/logout";
  });
};
