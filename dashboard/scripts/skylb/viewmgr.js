/**
 * @fileoverview The view manager.
 */

goog.provide('skylb.ViewMgr');

goog.require('goog.Disposable');
goog.require('goog.Uri');
goog.require('goog.events');
goog.require('goog.events.EventType');
goog.require('goog.structs.Map');


/**
 * The view manager.
 *
 * @param {goog.dom.ViewportSizeMonitor} vsm The view port manager.
 * @param {Element} container The view container.
 * @constructor
 * @extends {goog.Disposable}
 */
skylb.ViewMgr = function(vsm, container) {
  goog.base(this);

  /**
   * @private {proto.proto.UserInfo}
   */
  this.viewer_ = null;

  /**
   * The view container element.
   * @private {Element}
   */
  this.container_ = container;

  /**
   * @private {goog.structs.Map<string, !skylb.View>}
   */
  this.views_ = new goog.structs.Map();

  /**
   * @private {goog.structs.Map<string, Object>}
   */
  this.data_ = new goog.structs.Map();

  /**
   * The current view.
   * @private {skylb.View}
   */
  this.view_ = null;

  /**
   * The viewport size monitor.
   * @private {goog.dom.ViewportSizeMonitor}
   */
  this.vsm_ = vsm;

  /**
   * The last URL.
   * @private {string}
   */
  this.lastUrl_ = null;

  /**
   * @type {string}
   */
  this.defaultUrl = "";

  goog.events.listen(this.vsm_, goog.events.EventType.RESIZE, function(e) {
    if (this.view_) {
      this.view_.updateViewSize();
    }
  }, undefined, this);

  goog.events.listen(window, goog.events.EventType.POPSTATE, function(event) {
    this.loadViewByUrl(window.location, false);
  }, undefined, this);
};
goog.inherits(skylb.ViewMgr, goog.Disposable);


/**
 * Get the current viewer.
 * @return {proto.proto.UserInfo} The current viewer.
 */
skylb.ViewMgr.prototype.getViewer = function() {
  return this.viewer_;
};


/**
 * Set the current viewer.
 * @param {proto.proto.UserInfo} viewer The current viewer.
 */
skylb.ViewMgr.prototype.setViewer = function(viewer) {
  this.viewer_ = viewer;
};


/**
 * @return {Element} the view container.
 */
skylb.ViewMgr.prototype.getContainer = function() {
  return this.container_;
};


/**
 * @return {goog.math.Size} the viewport size.
 */
skylb.ViewMgr.prototype.getViewportSize = function() {
  return this.vsm_.getSize();
};


/**
 * Register a view.
 *
 * @param {skylb.View} view The view to register.
 */
skylb.ViewMgr.prototype.registerView = function(view) {
  this.views_.set(view.getName(), view);
};


/**
 * Load a view by the given view name.
 *
 * @param {string} name The view name.
 * @param {boolean} pushState Whether to push history state.
 * @param {Object} data The data object.
 */
skylb.ViewMgr.prototype.loadViewByName =
    function(name, pushState, data) {
  var view = this.views_.get(name);
  if (view) {
    this.view_ = view;
    view.reset();
    view.setData(data);

    view.load(goog.bind(function() {
      var url = view.getUrl();
      if (pushState && url != this.lastUrl_) {
        this.lastUrl_ = url;
        this.pushState_(url, view.getTitle());
      }
    }, this));
  } else {
    // TODO(zhwang): show error.
  }
};


/**
 * Load a view by the given URL.
 *
 * @param {string} url The view URL.
 * @param {boolean} pushState Whether to push history state.
 */
skylb.ViewMgr.prototype.loadViewByUrl = function(url, pushState) {
  var uri = goog.Uri.parse(url);
  var path = uri.getPath();
  if (path == '' || path == '/') {
    path = this.defaultUrl;
  }

  /** @type{goog.Uri.QueryData} */
  var query = uri.getQueryData();

  this.views_.forEach(function(view, key) {
    var regex = view.getUrlRegex();
    var m = path.match(regex);
    if (m) {
      this.view_ = view;
      view.reset();
      view.setMatch(m);
      view.setQuery(query);

      view.load(goog.bind(function() {
        var url = view.getUrl();
        if (pushState && url != this.lastUrl_) {
          this.lastUrl_ = url;
          this.pushState_(url, view.getTitle());
        }
      }, this));
      return;
    }
  }, this);
};


/**
 * Refresh the current view.
 */
skylb.ViewMgr.prototype.refreshView = function() {
  this.loadViewByUrl(window.location, /* pushState */ false);
};


/**
 * Push the history state.
 * @param {string} url The url to push.
 * @param {string} title The title of the page.
 * @private
 */
skylb.ViewMgr.prototype.pushState_ = function(url, title) {
  if (window.history && window.history.pushState) {
    window.history.pushState(null, title, url);
  }
};
