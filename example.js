/*swipe*/
function Swipe(m,e){var f=function(){};var u=function(C){setTimeout(C||f,0)};var B={addEventListener:!!window.addEventListener,touch:("ontouchstart" in window)||window.DocumentTouch&&document instanceof DocumentTouch,transitions:(function(C){var E=["transitionProperty","WebkitTransition","MozTransition","OTransition","msTransition"];for(var D in E){if(C.style[E[D]]!==undefined){return true}}return false})(document.createElement("swipe"))};if(!m){return}var c=m.children[0];var s,d,r,g;e=e||{};var k=parseInt(e.startSlide,10)||0;var v=e.speed||300;e.continuous=e.continuous!==undefined?e.continuous:true;function n(){s=c.children;g=s.length;if(s.length<2){e.continuous=false}if(B.transitions&&e.continuous&&s.length<3){c.appendChild(s[0].cloneNode(true));c.appendChild(c.children[1].cloneNode(true));s=c.children}d=new Array(s.length);r=m.getBoundingClientRect().width||m.offsetWidth;c.style.width=(s.length*r)+"px";var D=s.length;while(D--){var C=s[D];C.style.width=r+"px";C.setAttribute("data-index",D);if(B.transitions){C.style.left=(D*-r)+"px";q(D,k>D?-r:(k<D?r:0),0)}}if(e.continuous&&B.transitions){q(i(k-1),-r,0);q(i(k+1),r,0)}if(!B.transitions){c.style.left=(k*-r)+"px"}m.style.visibility="visible"}function o(){if(e.continuous){a(k-1)}else{if(k){a(k-1)}}}function p(){if(e.continuous){a(k+1)}else{if(k<s.length-1){a(k+1)}}}function i(C){return(s.length+(C%s.length))%s.length}function a(G,D){if(k==G){return}if(B.transitions){var F=Math.abs(k-G)/(k-G);if(e.continuous){var C=F;F=-d[i(G)]/r;if(F!==C){G=-F*s.length+G}}var E=Math.abs(k-G)-1;while(E--){q(i((G>k?G:k)-E-1),r*F,0)}G=i(G);q(k,r*F,D||v);q(G,0,D||v);if(e.continuous){q(i(G-F),-(r*F),0)}}else{G=i(G);j(k*-r,G*-r,D||v)}k=G;u(e.callback&&e.callback(k,s[k]))}function q(C,E,D){l(C,E,D);d[C]=E}function l(D,G,F){var C=s[D];var E=C&&C.style;if(!E){return}E.webkitTransitionDuration=E.MozTransitionDuration=E.msTransitionDuration=E.OTransitionDuration=E.transitionDuration=F+"ms";E.webkitTransform="translate("+G+"px,0)translateZ(0)";E.msTransform=E.MozTransform=E.OTransform="translateX("+G+"px)"}function j(G,F,C){if(!C){c.style.left=F+"px";return}var E=+new Date;var D=setInterval(function(){var H=+new Date-E;if(H>C){c.style.left=F+"px";if(A){x()}e.transitionEnd&&e.transitionEnd.call(event,k,s[k]);clearInterval(D);return}c.style.left=(((F-G)*(Math.floor((H/C)*100)/100))+G)+"px"},4)}var A=e.auto||0;var w;function x(){w=setTimeout(p,A)}function t(){A=0;clearTimeout(w)}var h={};var y={};var z;var b={handleEvent:function(C){switch(C.type){case"touchstart":this.start(C);break;case"touchmove":this.move(C);break;case"touchend":u(this.end(C));break;case"webkitTransitionEnd":case"msTransitionEnd":case"oTransitionEnd":case"otransitionend":case"transitionend":u(this.transitionEnd(C));break;case"resize":u(n);break}if(e.stopPropagation){C.stopPropagation()}},start:function(C){var D=C.touches[0];h={x:D.pageX,y:D.pageY,time:+new Date};z=undefined;y={};c.addEventListener("touchmove",this,false);c.addEventListener("touchend",this,false)},move:function(C){if(C.touches.length>1||C.scale&&C.scale!==1){return}if(e.disableScroll){C.preventDefault()}var D=C.touches[0];y={x:D.pageX-h.x,y:D.pageY-h.y};if(typeof z=="undefined"){z=!!(z||Math.abs(y.x)<Math.abs(y.y))}if(!z){C.preventDefault();t();if(e.continuous){l(i(k-1),y.x+d[i(k-1)],0);l(k,y.x+d[k],0);l(i(k+1),y.x+d[i(k+1)],0)}else{y.x=y.x/((!k&&y.x>0||k==s.length-1&&y.x<0)?(Math.abs(y.x)/r+1):1);l(k-1,y.x+d[k-1],0);l(k,y.x+d[k],0);l(k+1,y.x+d[k+1],0)}}},end:function(E){var G=+new Date-h.time;var D=Number(G)<250&&Math.abs(y.x)>20||Math.abs(y.x)>r/2;var C=!k&&y.x>0||k==s.length-1&&y.x<0;if(e.continuous){C=false}var F=y.x<0;if(!z){if(D&&!C){if(F){if(e.continuous){q(i(k-1),-r,0);q(i(k+2),r,0)}else{q(k-1,-r,0)}q(k,d[k]-r,v);q(i(k+1),d[i(k+1)]-r,v);k=i(k+1)}else{if(e.continuous){q(i(k+1),r,0);q(i(k-2),-r,0)}else{q(k+1,r,0)}q(k,d[k]+r,v);q(i(k-1),d[i(k-1)]+r,v);k=i(k-1)}e.callback&&e.callback(k,s[k])}else{if(e.continuous){q(i(k-1),-r,v);q(k,0,v);q(i(k+1),r,v)}else{q(k-1,-r,v);q(k,0,v);q(k+1,r,v)}}}c.removeEventListener("touchmove",b,false);c.removeEventListener("touchend",b,false)},transitionEnd:function(C){if(parseInt(C.target.getAttribute("data-index"),10)==k){t();e.transitionEnd&&e.transitionEnd.call(C,k,s[k]);A=e.auto||0;x()}}};n();if(A){x()}if(B.addEventListener){if(B.touch){c.addEventListener("touchstart",b,false)}if(B.transitions){c.addEventListener("webkitTransitionEnd",b,false);c.addEventListener("msTransitionEnd",b,false);c.addEventListener("oTransitionEnd",b,false);c.addEventListener("otransitionend",b,false);c.addEventListener("transitionend",b,false)}window.addEventListener("resize",b,false)}else{window.onresize=function(){n()}}return{setup:function(){n()},slide:function(D,C){t();a(D,C)},prev:function(){t();o()},next:function(){t();p()},stop:function(){t()},getPos:function(){return k},getNumSlides:function(){return g},kill:function(){t();c.style.width="";c.style.left="";var D=s.length;while(D--){var C=s[D];C.style.width="";C.style.left="";if(B.transitions){l(D,0,0)}}if(B.addEventListener){c.removeEventListener("touchstart",b,false);c.removeEventListener("webkitTransitionEnd",b,false);c.removeEventListener("msTransitionEnd",b,false);c.removeEventListener("oTransitionEnd",b,false);c.removeEventListener("otransitionend",b,false);c.removeEventListener("transitionend",b,false);window.removeEventListener("resize",b,false)}else{window.onresize=null}}}}if(window.jQuery||window.Zepto){(function(a){a.fn.Swipe=function(b){return this.each(function(){a(this).data("Swipe",new Swipe(a(this)[0],b))})}})(window.jQuery||window.Zepto)};
/*微信相关*/
$(function () {
//tab
 $(".weui_tab .weui_navbar_item").click(function(){
        $(this).addClass("tab-green").siblings().removeClass('tab-green');
        });
 $(".weui_tab_red .weui_navbar_item").click(function(){
        $(this).addClass("tab-red").siblings().removeClass('tab-red');
        });
 $(".weui_tab_blue .weui_navbar_item").click(function(){
        $(this).addClass("tab-blue").siblings().removeClass('tab-blue');
        }); 

 $('#search_input').focus(function(){//获得焦点
 var $weuiSearchBar = $('#search_bar');
$weuiSearchBar.addClass('weui_search_focusing');
$('#search-fixed').addClass('search-fixed');
 });
 $('#search_input').blur(function(){//失去焦点
  var $weuiSearchBar = $('#search_bar');
                    $weuiSearchBar.removeClass('weui_search_focusing');
                    $('#search-fixed').removeClass('search-fixed');
                    if($(this).val()){
                        $('#search_text').hide();
                    }else{
                        $('#search_text').show();
                    }
 });
  $('#search_input').focus(function(){
 var $searchShow = $("#search_show");
                    if($(this).val()){
                        $searchShow.show();
                    }else{
                        $searchShow.hide();
                    }
 }); 
   $('#search_cancel').tap(function(){
   $("#search_show").hide();
                    $('#search_input').val('');
 }); 
  $('#search_clear').tap(function(){
$("#search_show").hide();
                    $('#search_input').val('');
 });  
});

/***select**/
+function(a){"use strict";a.Template7=a.t7=function(){function a(a){return"[object Array]"===Object.prototype.toString.apply(a)}function c(a){return"function"==typeof a}function e(a){var d,e,f,g,h,i,j,k,b=a.replace(/[{}#}]/g,"").split(" "),c=[];for(e=0;e<b.length;e++)if(g=b[e],0===e)c.push(g);else if(0===g.indexOf('"'))if(2===g.match(/"/g).length)c.push(g);else{for(d=0,f=e+1;f<b.length;f++)if(g+=" "+b[f],b[f].indexOf('"')>=0){d=f,c.push(g);break}d&&(e=d)}else if(g.indexOf("=")>0){if(h=g.split("="),i=h[0],j=h[1],2!==j.match(/"/g).length){for(d=0,f=e+1;f<b.length;f++)if(j+=" "+b[f],b[f].indexOf('"')>=0){d=f;break}d&&(e=d)}k=[i,j.replace(/"/g,"")],c.push(k)}else c.push(g);return c}function f(b){var d,f,h,i,j,k,l,m,n,p,q,r,s,t,u,w,c=[];if(!b)return[];for(h=b.split(/({{[^{^}]*}})/),d=0;d<h.length;d++)if(i=h[d],""!==i)if(i.indexOf("{{")<0)c.push({type:"plain",content:i});else{if(i.indexOf("{/")>=0)continue;if(i.indexOf("{#")<0&&i.indexOf(" ")<0&&i.indexOf("else")<0){c.push({type:"variable",contextName:i.replace(/[{}]/g,"")});continue}for(j=e(i),k=j[0],l=[],m={},f=1;f<j.length;f++)n=j[f],a(n)?m[n[0]]="false"===n[1]?!1:n[1]:l.push(n);if(i.indexOf("{#")>=0){for(p="",q="",r=0,t=!1,u=!1,w=0,f=d+1;f<h.length;f++)if(h[f].indexOf("{{#")>=0&&w++,h[f].indexOf("{{/")>=0&&w--,h[f].indexOf("{{#"+k)>=0)p+=h[f],u&&(q+=h[f]),r++;else if(h[f].indexOf("{{/"+k)>=0){if(!(r>0)){s=f,t=!0;break}r--,p+=h[f],u&&(q+=h[f])}else h[f].indexOf("else")>=0&&0===w?u=!0:(u||(p+=h[f]),u&&(q+=h[f]));t&&(s&&(d=s),c.push({type:"helper",helperName:k,contextName:l,content:p,inverseContent:q,hash:m}))}else i.indexOf(" ")>0&&c.push({type:"helper",helperName:k,contextName:l,hash:m})}return c}var h,g=function(a){function c(a,b){return a.content?h(a.content,b):function(){return""}}function d(a,b){return a.inverseContent?h(a.inverseContent,b):function(){return""}}function e(a,b){var c,d,g,h,i,e=0;for(0===a.indexOf("../")?(e=a.split("../").length-1,g=b.split("_")[1]-e,b="ctx_"+(g>=1?g:1),d=a.split("../")[e].split(".")):0===a.indexOf("@global")?(b="$.Template7.global",d=a.split("@global.")[1].split(".")):0===a.indexOf("@root")?(b="ctx_1",d=a.split("@root.")[1].split(".")):d=a.split("."),c=b,h=0;h<d.length;h++)i=d[h],0===i.indexOf("@")?h>0?c+="[(data && data."+i.replace("@","")+")]":c="(data && data."+a.replace("@","")+")":isFinite(i)?c+="["+i+"]":0===i.indexOf("this")?c=i.replace("this",b):c+="."+i;return c}function g(a,b){var d,c=[];for(d=0;d<a.length;d++)0===a[d].indexOf('"')?c.push(a[d]):c.push(e(a[d],b));return c.join(", ")}function h(a,h){var i,j,k,l,o,p,q;if(h=h||1,a=a||b.template,"string"!=typeof a)throw new Error("Template7: Template must be a string");if(i=f(a),0===i.length)return function(){return""};for(j="ctx_"+h,k="(function ("+j+", data) {\n",1===h&&(k+="function isArray(arr){return Object.prototype.toString.apply(arr) === '[object Array]';}\n",k+="function isFunction(func){return (typeof func === 'function');}\n",k+='function c(val, ctx) {if (typeof val !== "undefined") {if (isFunction(val)) {return val.call(ctx);} else return val;} else return "";}\n'),k+="var r = '';\n",l=0;l<i.length;l++)if(o=i[l],"plain"!==o.type){if("variable"===o.type&&(p=e(o.contextName,j),k+="r += c("+p+", "+j+");"),"helper"===o.type)if(o.helperName in b.helpers)q=g(o.contextName,j),k+="r += ($.Template7.helpers."+o.helperName+").call("+j+", "+(q&&q+", ")+"{hash:"+JSON.stringify(o.hash)+", data: data || {}, fn: "+c(o,h+1)+", inverse: "+d(o,h+1)+", root: ctx_1});";else{if(o.contextName.length>0)throw new Error('Template7: Missing helper: "'+o.helperName+'"');p=e(o.helperName,j),k+="if ("+p+") {",k+="if (isArray("+p+")) {",k+="r += ($.Template7.helpers.each).call("+j+", "+p+", {hash:"+JSON.stringify(o.hash)+", data: data || {}, fn: "+c(o,h+1)+", inverse: "+d(o,h+1)+", root: ctx_1});",k+="}else {",k+="r += ($.Template7.helpers.with).call("+j+", "+p+", {hash:"+JSON.stringify(o.hash)+", data: data || {}, fn: "+c(o,h+1)+", inverse: "+d(o,h+1)+", root: ctx_1});",k+="}}"}}else k+="r +='"+o.content.replace(/\r/g,"\\r").replace(/\n/g,"\\n").replace(/'/g,"\\'")+"';";return k+="\nreturn r;})",eval.call(window,k)}var b=this;b.template=a,b.compile=function(a){return b.compiled||(b.compiled=h(a)),b.compiled}};return g.prototype={options:{},helpers:{"if":function(a,b){return c(a)&&(a=a.call(this)),a?b.fn(this,b.data):b.inverse(this,b.data)},unless:function(a,b){return c(a)&&(a=a.call(this)),a?b.inverse(this,b.data):b.fn(this,b.data)},each:function(b,d){var g,e="",f=0;if(c(b)&&(b=b.call(this)),a(b)){for(d.hash.reverse&&(b=b.reverse()),f=0;f<b.length;f++)e+=d.fn(b[f],{first:0===f,last:f===b.length-1,index:f});d.hash.reverse&&(b=b.reverse())}else for(g in b)f++,e+=d.fn(b[g],{key:g});return f>0?e:d.inverse(this)},"with":function(a,b){return c(a)&&(a=a.call(this)),b.fn(a)},join:function(a,b){return c(a)&&(a=a.call(this)),a.join(b.hash.delimiter||b.hash.delimeter)},js:function(a){var c;return c=a.indexOf("return")>=0?"(function(){"+a+"})":"(function(){return ("+a+")})",eval.call(this,c).call(this)},js_compare:function(a,b){var c,d;return c=a.indexOf("return")>=0?"(function(){"+a+"})":"(function(){return ("+a+")})",d=eval.call(this,c).call(this),d?b.fn(this,b.data):b.inverse(this,b.data)}}},h=function(a,b){var c,d;return 2===arguments.length?(c=new g(a),d=c.compile()(b),c=null,d):new g(a)},h.registerHelper=function(a,b){g.prototype.helpers[a]=b},h.unregisterHelper=function(a){g.prototype.helpers[a]=void 0,delete g.prototype.helpers[a]},h.compile=function(a,b){var c=new g(a,b);return c.compile()},h.options=g.prototype.options,h.helpers=g.prototype.helpers,h}()}($);+function(a){"use strict";a.openPopup=function(b){a.closePopup(),b=a(b),b.addClass("weui-popup-container-visible");var d=b.find(".weui-popup-modal");d.width(),d.addClass("weui-popup-modal-visible")},a.closePopup=function(b,c){a(".weui-popup-modal-visible").removeClass("weui-popup-modal-visible").transitionEnd(function(){a(this).parent().removeClass("weui-popup-container-visible"),c&&a(this).parent().remove()}).trigger("close")},a(document).on("click",".close-popup",function(){a.closePopup()}),a(document).on("click",".open-popup",function(){a(a(this).data("target")).popup()}),a.fn.popup=function(){return this.each(function(){a.openPopup(this)})}}($);+function(a){"use strict";var b,c=function(b,c){var e,d=this;this.config=c,this.$input=a(b),e=a.t7.compile("<div class='weui-picker-modal weui-select-modal'>"+c.toolbarTemplate+(c.multi?c.checkboxTemplate:c.radioTemplate)+"</div>"),this.$input.prop("readOnly",!0),this.$input.click(function(){d.parseInitValue();var b=d.dialog=a.openPicker(e({items:c.items,title:c.title,closeText:c.closeText}));b.on("change",function(){var f=b.find("input:checked"),g=f.map(function(){return a(this).val()}),h=f.map(function(){return a(this).data("title")});d.updateInputValue(g,h),c.autoClose&&!c.multi&&a.closePicker()})}),a(document).on("click",function(){})};c.prototype.updateInputValue=function(a,b){var c,d,e;this.config.multi?(c=a.join(this.config.split),d=b.join(this.config.split)):(c=a[0],d=b[0]),this.$input.val(d).data("values",c),this.$input.attr("value",d).attr("data-values",c),e={values:c,titles:d},this.$input.trigger("change",e),this.config.onChange&&this.config.onChange.call(this,e)},c.prototype.parseInitValue=function(){var c,d,e,a=this.$input.val(),b=this.config.items;if(void 0!==a&&null!=a&&""!==a)for(c=this.config.multi?a.split(this.config.split):[a],d=0;d<b.length;d++)for(b[d].checked=!1,e=0;e<c.length;e++)b[d].title===c[e]&&(b[d].checked=!0)},a.fn.select=function(d){var e=a.extend({},b,d);if(e.items&&e.items.length)return e.items=e.items.map(function(a){return"string"==typeof a?{title:a,value:a}:a}),this.each(function(){var d,b=a(this);return b.data("weui-select")||b.data("weui-select",new c(this,e)),d=b.data("weui-select")})},b=a.fn.select.prototype.defaults={items:[],title:"请选择",multi:!1,closeText:"关闭",autoClose:!0,onChange:void 0,split:",",toolbarTemplate:'<div class="toolbar">      <div class="toolbar-inner">      <a href="javascript:;" class="picker-button close-picker">{{closeText}}</a>      <h1 class="title">{{title}}</h1>      </div>      </div>',radioTemplate:'<div class="weui_cells weui_cells_radio">        {{#items}}        <label class="weui_cell weui_check_label" for="weui-select-id-{{this.title}}">          <div class="weui_cell_bd weui_cell_primary">            <p>{{this.title}}</p>          </div>          <div class="weui_cell_ft">            <input type="radio" class="weui_check" name="weui-select" id="weui-select-id-{{this.title}}" value="{{this.value}}" {{#if this.checked}}checked="checked"{{/if}} data-title="{{this.title}}">            <span class="weui_icon_checked"></span>          </div>        </label>        {{/items}}      </div>',checkboxTemplate:'<div class="weui_cells weui_cells_checkbox">        {{#items}}        <label class="weui_cell weui_check_label" for="weui-select-id-{{this.title}}">          <div class="weui_cell_bd weui_cell_primary">            <p>{{this.title}}</p>          </div>          <div class="weui_cell_ft">            <input type="checkbox" class="weui_check" name="weui-select" id="weui-select-id-{{this.title}}" value="{{this.value}}" {{#if this.checked}}checked="checked"{{/if}} data-title="{{this.title}}" >            <span class="weui_icon_checked"></span>          </div>        </label>        {{/items}}      </div>'}}($);