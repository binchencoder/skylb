/**
 * @fileoverview The base view class.
 */

goog.provide('skylb.View');

goog.require('goog.Disposable');
goog.require('goog.dom');
goog.require('goog.dom.dataset');


/**
 * Constructs a base view class.
 * @param {string} name The name of the view.
 * @param {string} title The title of the view.
 * @param {RegExp} regex The regex of the view.
 * @param {skylb.XhrMgr} xhrMgr The xhr manager.
 * @param {skylb.ViewMgr} viewMgr The view manager.
 *
 * @constructor
 * @extends {goog.Disposable}
 */
skylb.View = function(name, title, regex, xhrMgr, viewMgr) {
  goog.base(this);

  /**
   * @private {string} The name of the view.
   */
  this.name_ = name;

  /**
   * @private {string} The title of the view.
   */
  this.title_ = title;

  /**
   * @private {RegExp} The regex of the view URL.
   */
  this.regex_ = regex;

  /**
   * @protected {skylb.XhrMgr} The Xhr manager.
   */
  this.xhrMgr = xhrMgr;

  /**
   * @protected {skylb.ViewMgr} The view manager.
   */
  this.viewMgr = viewMgr;

  /**
   * @protected {Object}
   */
  this.data = {};

  /**
   * @protected {Array}
   */
  this.match = [];

  /**
   * @protected {goog.Uri.QueryData}
   */
  this.query = null;
};
goog.inherits(skylb.View, goog.Disposable);


/** @override */
skylb.View.prototype.disposeInternal = function() {
  goog.base(this, 'disposeInternal');
};


/**
 * @return {string} The name of the view.
 */
skylb.View.prototype.getName = function() {
  return this.name_;
};


/**
 * @return {string} The title of the view.
 */
skylb.View.prototype.getTitle = function() {
  return this.title_;
};


/**
 * @return {RegExp} the URL regex.
 */
skylb.View.prototype.getUrlRegex = function() {
  return this.regex_;
};


/**
 * @param {Object} data The data.
 */
skylb.View.prototype.setData = function(data) {
  this.data = data;
};


/**
 * @param {Array} match The regex match.
 */
skylb.View.prototype.setMatch = function(match) {
  this.match = match;
};


/**
 * @param {Array} query The query.
 */
skylb.View.prototype.setQuery = function(query) {
  this.query = query;
};


/**
 * @return {string} the URL.
 */
skylb.View.prototype.getUrl = goog.abstractMethod;


/**
 * Load the view.
 * @param {Function} callback The callback function if it's successfully
 *     loaded.
 */
skylb.View.prototype.load = goog.abstractMethod;


/**
 * Reset the view.
 */
skylb.View.prototype.reset = function() {
  this.data = {};
  this.match = [];
  this.query = null;
};

/**
 * Update the view size.
 */
skylb.View.prototype.updateViewSize = function() {};
