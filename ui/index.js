/*
 * Copyleft 2017, Simone Margaritelli <evilsocket at protonmail dot com>
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   * Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *   * Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *   * Neither the name of ARM Inject nor the names of its contributors may be used
 *     to endorse or promote products derived from this software without
 *     specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */
var app = angular.module('XRAY', [], function($interpolateProvider) {

});

// Convert object into array of objects
app.filter('toArray', function () {
  return function (obj, addKey) {
    if (!angular.isObject(obj)) return obj;
    if ( addKey === false ) {
      return Object.keys(obj).map(function(key) {
        return obj[key];
      });
    } else {
      return Object.keys(obj).map(function (key) {
        var value = obj[key];
        return angular.isObject(value) ?
          Object.defineProperty(value, '$key', { enumerable: false, value: key}) :
          { $key: key, $value: value };
      });
    }
  };
});

// Main controller
app.controller('XRayController', ['$scope', '$http', function (scope, http) {

  // Init variables
  scope.domain = "";
  scope.stats = {
      Start: "",
      Stop: "",
      Total: 0,
      Inputs: 0,
      Eps: 0.0,
      Execs: 0,
      Results: 0,
      Progress: 0.0,
  };
  scope.targets = { };
  scope.ntargets = 0;
  scope.duration = 0;
  scope.firstTimeUpdate = false;
  scope.showEmpty = null;
  scope.searchText = null;
  scope.reverse = true;

  // Sort by property name
  scope.sortBy = function(property) {
    scope.reverse = (scope.order === property) ? !scope.reverse : false;
    scope.order = property;
  };

  // Reset filter and sort variables to original value
  scope.reset = function() {
    scope.reverse = true;
    scope.order = null;
    scope.searchText = null;
  }

  // Data refresh function
  scope.update = function() {

      http.get('/targets').then(function(response) {

        var data = response.data;

        if( data.stats.Progress < 100.0 || scope.firstTimeUpdate == false ) {
            var start = new Date(data.stats.Start),
                stop = new Date(data.stats.Stop),
                dur = new Date(null);

            dur.setSeconds( (stop-start) / 1000 );
            scope.duration = dur.toISOString().substr(11, 8);
        }

        // Update ui variables
        scope.ntargets = Object.keys(scope.targets).length;
        scope.targets = data.targets;
        scope.domain = data.domain;
        scope.stats = data.stats;

        // Update page title
        document.title = "XRAY ( " + scope.domain + " | " + scope.stats.Progress.toFixed(2) + "% )";

        scope.firstTimeUpdate = true;

      });

  }

  // Start refresh interval
  setInterval( scope.update, 500 );

}]);
