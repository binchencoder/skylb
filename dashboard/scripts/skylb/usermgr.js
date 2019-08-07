/**
 * @fileoverview The UsersManagementView class.
 */

goog.provide('skylb.UsersManagementView');

goog.require('goog.array');
goog.require('goog.dom');
goog.require('goog.dom.classlist');
goog.require('goog.dom.forms');
goog.require('goog.events');
goog.require('goog.net.EventType');
goog.require('goog.net.XhrIo.ResponseType');
goog.require('goog.ui.Dialog');
goog.require('proto.proto.GetUsersRequest');
goog.require('proto.proto.GetUsersResponse');
goog.require('proto.proto.UpsertUserRequest');
goog.require('proto.proto.UpsertUserResponse');
goog.require('proto.proto.UserInfo');
goog.require('skylb.View');
goog.require('skylb.usermgr_html.AddOrEditUserTemplate');
goog.require('skylb.usermgr_html.ManageUsersTemplate');


/**
 * Constructs a UsersManagementView class.
 *
 * @param {skylb.XhrMgr} xhrMgr The Xhr manager.
 * @param {skylb.ViewMgr} viewMgr The view manager.
 * @constructor
 * @extends {skylb.View}
 */
skylb.UsersManagementView = function(xhrMgr, viewMgr) {
  goog.base(this, 'users', 'Users Management', /^\/users$/, xhrMgr, viewMgr);
};
goog.inherits(skylb.UsersManagementView, skylb.View);


/** @override */
skylb.UsersManagementView.prototype.getUrl = function() {
  return '/users';
};


/** @override */
skylb.UsersManagementView.prototype.load = function(callback) {
  /** @type{Element} */
  var container = this.viewMgr.getContainer();
  goog.dom.removeChildren(container);

  var req = new proto.proto.GetUsersRequest();
  var bytes = req.serializeBinary();

  var xhr = this.xhrMgr.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);
  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var data = xhr.getResponse();
    this.xhrMgr.returnXhr(xhr);

    var resp = proto.proto.GetUsersResponse.deserializeBinary(data);
    var users = resp.getUsersList();

    var t = new skylb.usermgr_html.ManageUsersTemplate();
    t.render(container, {'users': users});

    var divs = goog.dom.getElementsByClass('user-card');
    goog.array.forEach(divs, function(div, idx) {
      if (goog.dom.classlist.contains(div, 'create-user')) {
        goog.events.listen(div, goog.events.EventType.CLICK, function(e) {
          this.createEditUserDialog_(new proto.proto.UserInfo(), true);
        }, undefined, this);
      } else {
        var edit = goog.dom.getElementByClass('edit', div);
        goog.events.listen(edit, goog.events.EventType.CLICK, function(e) {
          this.createEditUserDialog_(users[idx], false);
        }, undefined, this);
      }
    }, this);
  }, undefined, this);

  xhr.send('/_/get-users', 'POST', bytes);

  if (callback) {
    callback();
  }
};


/**
 * Starts the create/edit user dialog.
 * @param {proto.proto.UserInfo} userInfo The user info.
 * @param {boolean} isNew True to create a new user.
 * @private
 */
skylb.UsersManagementView.prototype.createEditUserDialog_ =
    function(userInfo, isNew) {
  var userDialog = new goog.ui.Dialog("edit-user-dialog");
  userDialog.setDisposeOnHide(true);
  userDialog.setModal(true);
  userDialog.setHasTitleCloseButton(true);
  userDialog.setButtonSet(null);
  var title = userInfo.getLoginName() ? 'Edit User' : "Create User";
  userDialog.setTitle(title);

  var container = userDialog.getContentElement();

  var t = new skylb.usermgr_html.AddOrEditUserTemplate();
  t.render(container, {'user': userInfo});

  var saveBtn = goog.dom.getElementByClass('save-button', container);
  goog.events.listen(saveBtn, goog.events.EventType.CLICK, function(e) {
    var errEle = goog.dom.getElementByClass('user-edit-errors', container);
    goog.dom.classlist.add(errEle, 'invisible');
    this.createEditUser_(userDialog, isNew);
  }, undefined, this);

  userDialog.setVisible(true);
};

/**
 * Create or edit user.
 * @param {goog.ui.Dialog} dialog The dialog.
 * @param {boolean} isNew True to create a new user.
 * @private
 */
skylb.UsersManagementView.prototype.createEditUser_ =
    function(dialog, isNew) {
  var info = new proto.proto.UserInfo();
  var container = dialog.getContentElement();

  var loginname = goog.dom.getElementByClass('username-input', container);
  info.setLoginName(goog.dom.forms.getValue(loginname));

  var version = goog.dom.getElementByClass('version-input', container);
  info.setVersion(goog.dom.forms.getValue(version));

  var disabled = goog.dom.getElementByClass('is-disabled', container);
  info.setDisabled(disabled.checked);

  var isAdmin = goog.dom.getElementByClass('is-admin-input', container);
  info.setIsAdmin(isAdmin.checked);
  var isDev = goog.dom.getElementByClass('is-dev-input', container);
  info.setIsDev(isDev.checked);
  var isOps = goog.dom.getElementByClass('is-ops-input', container);
  info.setIsOps(isOps.checked);

  var req = new proto.proto.UpsertUserRequest();
  req.setUser(info);
  req.setIsNew(isNew);

  var payload = req.serializeBinary();
  var xhr = this.xhrMgr.getXhrIo();
  xhr.setResponseType(goog.net.XhrIo.ResponseType.ARRAY_BUFFER);
  goog.events.listen(xhr, goog.net.EventType.SUCCESS, function(e) {
    var data = xhr.getResponse();
    this.xhrMgr.returnXhr(xhr);

    var resp = proto.proto.UpsertUserResponse.deserializeBinary(data);
    if (resp.getErrorMsg()) {
      this.showErrorMsg(container, resp.getErrorMsg());
    } else {
      dialog.setVisible(false);
      this.load(undefined);
    }
  }, undefined, this);

  xhr.send('/_/upsert-user', 'POST', payload);
};

/**
 * Display the given error message.
 * @param {Element} container The container element.
 * @param {string} msg The message.
 */
skylb.UsersManagementView.prototype.showErrorMsg =
    function(container, msg) {
  var errEle = goog.dom.getElementByClass('user-edit-errors', container);
  goog.dom.setTextContent(errEle, msg);
  goog.dom.classlist.remove(errEle, 'invisible');
};
