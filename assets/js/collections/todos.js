/*global Backbone */
var app = app || {};

(function () {
	'use strict';


	// Todo Collection
	// ---------------

	app.Todos = Backbone.Collection.extend({
		// Reference to this collection's model.
		model: app.Todo,

		// Filter down the list of all todo items that are finished.
		completed: function () {
			return this.where({completed: true});
		},

		// Filter down the list to only todo items that are still not finished.
		remaining: function () {
			return this.where({completed: false});
		},

		// Todos are sorted by their original insertion order.
		comparator: 'createdAt',

		url: '/',
	});
})();
