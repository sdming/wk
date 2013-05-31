function hello(s) {
	console.log(s);
}
var debugData;

function loadApiTest(options, modal) {
    var defaultOptions = {
        method: "GET",
        url: "",
        accepts: "",
        data:"",
        contentType:"",
        dataType:"html" //xml, json, script, or html
    };
    options = $.extend({}, defaultOptions, options);

    var headers = {};
    if (options.accepts){
        headers.Accept = options.accepts;
    }

    if (!modal) {
        modal = "#apiModal";
    }
    var m = $(modal);
    m.find(".api-status").text("loading");
    m.find(".api-data").text("");
    m.find(".api-link").text(options.url);
    m.modal({});

    $.ajax({
        url: options.url,
        headers: headers,
        cache: false,
        data: options.data,
        type:options.method,
        contentType :options.contentType,
        dataType: options.dataType
    }).success(function (data, status, jqXHR) {
        debugData = data;
        m.find(".api-status").text(status);
        m.find(".api-data").text(data);
    }).complete(function (jqXHR, status) {
        //console.log("complete");
    }).error(function (jqXHR, status, error) {
        m.find(".api-status").text(status);
        m.find(".api-data").text(error);
    });
}

$(document).ready(function () {
    $('a.api-link').click(function () {
        var method = "GET";
        if ($(this).attr("data-method")) {
            method = $(this).attr("data-method");
        }

        var options = {
            url:$(this).attr("href"),
            method:method
        };
        if ($(this).attr("data-accepts")) {   
            options.accepts = $(this).attr("data-accepts");
        }
        if ($(this).attr("data-content")) {
            options.contentType = $(this).attr("data-content");
        }
        if ($(this).attr("data-data")) {
            options.data = $(this).attr("data-data");
        }
        loadApiTest(options);
        return false;
    });
}); 


function setScrollTopNav() {
     $(window).on('scroll', function () {
         var scrolltop = $(this).scrollTop();

         if (scrolltop >= 215) {
             $('#dynamicTopNavBar').css({ "position": "fixed", "top": "0","z-index":"999"});
         }

         else if (scrolltop <= 210) {
             $('#dynamicTopNavBar').css({ "position": "static", "top":"","z-index":"" });
         }
     });
}
