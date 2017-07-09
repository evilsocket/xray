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

app.controller('XRayController', ['$scope', function (scope) {
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

    scope.update = function() {
        $.get('/targets', function(data) {
            var start = new Date(data.stats.Start),
                stop = new Date(data.stats.Stop),
                dur = new Date(null);

            dur.setSeconds( (stop-start) / 1000 );

            scope.duration = dur.toISOString().substr(11, 8);;
            scope.ntargets = Object.keys(scope.targets).length;

            scope.domain = data.domain;
            scope.stats = data.stats;
            scope.targets = data.targets;
            document.title = "XRAY ( " + scope.domain + " | " + scope.stats.Progress + "% )";
            scope.$apply();
        });
    }

    setInterval( scope.update, 500 );
}]);
