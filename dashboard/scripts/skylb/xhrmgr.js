/**
 * @fileoverview The Xhr manager.
 */
goog.provide('skylb.XhrMgr');

goog.require('goog.Disposable');
goog.require('goog.dom');
goog.require('goog.dom.classlist');
goog.require('goog.dom.dataset');
goog.require('goog.dom.forms');
goog.require('goog.events');
goog.require('goog.events.EventType');
goog.require('goog.events.KeyCodes');
goog.require('goog.events.KeyHandler');
goog.require('goog.net.EventType');
goog.require('goog.net.XhrIo');
goog.require('goog.net.XhrIoPool');
goog.require('goog.structs.Map');
goog.require('goog.style');
goog.require('goog.ui.Dialog');
goog.require('proto.proto.LoginRequest');
goog.require('proto.proto.LoginResponse');
goog.require('skylb.loginbox_html.LoginBoxTemplate');


/**
 * A manager class for Xhr.
 *
 * @constructor
 * @extends {goog.Disposable}
 */
skylb.XhrMgr = function() {
  goog.base(this);

  var headers = new goog.structs.Map();
  headers.set(goog.net.XhrIo.CONTENT_TYPE_HEADER, 'application/x-protobuf');
  headers.set('X-Requested-With', 'XMLHttpRequest');

  /**
   * @private {goog.net.XhrIoPool} The XhrIo pool.
   */
  this.pool_ = new goog.net.XhrIoPool(headers, 1, 3);
};
goog.inherits(skylb.XhrMgr, goog.Disposable);


/**
 * The timeout interval for AJAX calls.
 * @const {number}
 */
skylb.TIMEOUT_INTERVAL = 30000;


/** @override */
skylb.XhrMgr.prototype.disposeInternal = function() {
  goog.base(this, 'disposeInternal');
  this.pool_.dispose();
};


/**
 * @return {goog.net.XhrIo} an XhrIo instance from the pool.
 */
skylb.XhrMgr.prototype.getXhrIo = function() {
  var /** @type{goog.net.XhrIo} */ xhr = this.pool_.getObject();
  xhr.setTimeoutInterval(skylb.TIMEOUT_INTERVAL);

  var outerContainer = goog.dom.getElement('outer-container');

  goog.events.listen(xhr, goog.net.EventType.READY, function(e) {
    goog.style.setStyle(outerContainer, 'cursor', 'progress');
  }, undefined, this);

  goog.events.listen(xhr, goog.net.EventType.COMPLETE, function(e) {
    goog.style.setStyle(outerContainer, 'cursor', 'auto');
  }, true, this);

  goog.events.listen(xhr, goog.net.EventType.TIMEOUT, function(e) {
    // TODO(zhwang): show timeout.
    console.log("xhr timed out.");
    this.pool_.releaseObject(xhr);
  }, undefined, this);

  goog.events.listen(xhr, goog.net.EventType.ERROR, function(e) {
    if (xhr.getStatus() == 401) {
      // Session Expired.
      this.pool_.releaseObject(xhr);

      var loginDialog = new goog.ui.Dialog("login-dialog");
      loginDialog.setDisposeOnHide(true);
      loginDialog.setModal(true);
      loginDialog.setButtonSet(null);

      var container = loginDialog.getContentElement();
      loginDialog.setTitle('Login');

      var body = goog.dom.getDocument().body;
      var viewer = goog.dom.dataset.get(body, 'viewer');

      var t = new skylb.loginbox_html.LoginBoxTemplate();
      t.render(container, {'loginname': viewer});

      var loginBtn = goog.dom.getElementByClass('login-button', container);
      goog.events.listen(loginBtn, goog.events.EventType.CLICK, function(e) {
        var errEle = goog.dom.getElementByClass('login-errors', container);
        goog.dom.classlist.add(errEle, 'invisible');
        this.login_(loginDialog);
      }, undefined, this);

      // Listen the key press event.
      var kh = new goog.events.KeyHandler(loginDialog);
      goog.events.listen(kh, goog.events.KeyHandler.EventType.KEY, function(e){
        if (e.keyCode == goog.events.KeyCodes.ENTER) {
          this.login_(loginDialog);
        }
      }, undefined, this);

      loginDialog.setVisible(true);
      return;
    }

    // TODO(zhwang): show error.
    console.log("xhr failed.");
    this.pool_.releaseObject(xhr);
  }, undefined, this);

  return xhr;
};


/**
 * Return the XhrIo instance back to the pool.
 * @param {goog.net.XhrIo} xhr The xhrIo object.
 */
skylb.XhrMgr.prototype.returnXhr = function(xhr) {
  this.pool_.releaseObject(xhr);
};


/**
 * Login the current user.
 * @param {goog.ui.Dialog} dialog The dialog.
 * @private
 */
skylb.XhrMgr.prototype.login_ = function(dialog) {
  var xhr = this.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);

  var container = dialog.getContentElement();

  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var xhr = e.target;
    var data = xhr.getResponse();
    this.returnXhr(xhr);

    var resp = proto.proto.LoginResponse.deserializeBinary(data);
    if (resp.getErrorMsg()) {
      var errEle = goog.dom.getElementByClass('login-errors', container);
      goog.dom.setTextContent(errEle, resp.getErrorMsg());
      goog.dom.classlist.remove(errEle, 'invisible');
    } else {
      dialog.setVisible(false);
    }
  }, undefined, this);

  var loginname = goog.dom.getElementByClass('loginname-input', container);
  var password = goog.dom.getElementByClass('password-input', container);

  var req = new proto.proto.LoginRequest();
  req.setLoginName(goog.dom.forms.getValue(loginname));
  req.setPassword(goog.dom.forms.getValue(password));
  var payload = req.serializeBinary();

  xhr.send('/_/login', 'POST', payload);
};
