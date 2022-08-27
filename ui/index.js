var app = angular.module('XRAY', [], function($interpolateProvider) {

});

app.filter('toid', function() {
    return function(domain) {
       return domain.replace( /[^a-z0-9_]/g, '_' )
    }
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
    scope.firstTimeUpdate = false;

    scope.applyFilters = function(data) {
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

            data.targets = filtered;
        }
    };

    scope.update = function() {
        $.get('/targets', function(data) {
            if( data.stats.Progress < 100.0 || scope.firstTimeUpdate == false ) {
                var start = new Date(data.stats.Start),
                    stop = new Date(data.stats.Stop),
                    dur = new Date(null);

                dur.setSeconds( (stop-start) / 1000 );
                scope.duration = dur.toISOString().substr(11, 8);
            }
            
            scope.ntargets = Object.keys(scope.targets).length;

            scope.applyFilters(data);

            scope.targets = data.targets;
            scope.domain = data.domain;
            scope.stats = data.stats;
            
            document.title = "XRAY ( " + scope.domain + " | " + scope.stats.Progress.toFixed(2) + "% )";

            scope.$apply();
            scope.firstTimeUpdate = true;

            $('.htoggle').each(function() {
                $(this).click(function(e){
                    $( $(this).attr('href') ).toggle();
                    return false;
                });
            });
        });
    }

    setInterval( scope.update, 500 );
}]);
