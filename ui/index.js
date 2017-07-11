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
            if( scope.stats.Progress < 100.0 || scope.stats.Progress == 0.0 ) {
                var start = new Date(data.stats.Start),
                    stop = new Date(data.stats.Stop),
                    dur = new Date(null);

                dur.setSeconds( (stop-start) / 1000 );
                scope.duration = dur.toISOString().substr(11, 8);;
            }

            if( $('#show_empty').is(':checked') == false ) {
                var filtered = {};
                for( var ip in data.targets ) {
                    var t = data.targets[ip];
                    if( t.Info != null && t.Info.ports.length > 0 ) {
                        filtered[ip] = t;
                    }
                }

                data.targets = filtered;
            }

            var search = $('#search').val();
            if( search != "" ) {
                search = search.toLowerCase();

                var filtered = {};
                for( var ip in data.targets ) {
                    var t = data.targets[ip];
                    var txt = JSON.stringify(t).toLowerCase();
                    if( txt.search(search) >= 0 ) {
                        filtered[ip] = t;
                    }
                }

                console.log(filtered);

                data.targets = filtered;
            }

            scope.targets = data.targets;
            scope.ntargets = Object.keys(scope.targets).length;
            scope.domain = data.domain;
            scope.stats = data.stats;
            
            document.title = "XRAY ( " + scope.domain + " | " + scope.stats.Progress.toFixed(2) + "% )";

            if( $('#show_charts').is(':checked') == true ) {
                $('#charts').show();

                var countries_chart_data = {
                    datasets:[{
                        label: 'Hosts/Countries',
                        data: []
                    }],     

                    labels: []
                };
                var countries_chart_opts = {
                };

                var ports_chart_data = {
                    datasets:[{
                        label: 'Hosts/Port',
                        data: []
                    }],     

                    labels: []
                };
                var ports_chart_opts = {
                };
                
                for( var ip in data.targets ) {
                    var t = data.targets[ip];

                    if( t.Info != null && t.Info.country_code != "" ) {
                        var c = t.Info.country_code;
                        var idx = countries_chart_data.labels.indexOf(c);
                        
                        if( idx == -1 ) {
                            countries_chart_data.labels.push(c);
                            countries_chart_data.datasets[0].data.push(1)
                        } else {
                            countries_chart_data.datasets[0].data[idx] += 1
                        }
                    }

                    if( t.Info != null && t.Info.ports != [] ) {
                        for( var i in t.Info.ports ) {
                            var port = t.Info.ports[i];
                            var sport = ""+port;
                            var idx = ports_chart_data.labels.indexOf(sport);
                            
                            if( idx == -1 ) {
                                ports_chart_data.labels.push(sport);
                                ports_chart_data.datasets[0].data.push(1)
                            } else {
                                ports_chart_data.datasets[0].data[idx] += 1
                            }
                        }
                    }
                }
                
                var countries_chart = new Chart( 'countries_chart',{
                    type: 'horizontalBar',
                    data: countries_chart_data,
                    options: countries_chart_opts
                });

                var ports_chart = new Chart( 'ports_chart',{
                    type: 'bar',
                    data: ports_chart_data,
                    options: ports_chart_opts
                });
            } else {
                $('#charts').hide();
            }

            scope.$apply();
        });
    }

    setInterval( scope.update, 500 );
}]);
