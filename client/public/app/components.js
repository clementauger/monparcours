
function dateComponent(opts) {
  var inp=m("input", {class:"bt-calendar"});

  var picker;

  //todo: add mobile support.

  return {
    view: function(opts){
      var attrs = opts.attrs || {};
      var readOnly = attrs.readonly || false;
      var isMobile = attrs.isMobile || false;

      if(!attrs.value && attrs.bind){
        if (attrs.bind.p in attrs.bind.t) {
          attrs.value = attrs.bind.t[attrs.bind.p]
        }
      }
      if (readOnly) {
        var t = moment(attrs.value).format("LL")
        return m("div.datetime-cpnt.readonly", {},
          m("i.icon-calendar", { }, ""),
          m("label", {class: attrs.bind?attrs.bind.p:""}, attrs.placeholder),
          m("span", {class: attrs.bind?attrs.bind.p:""}, t)
        )
      }

      if(isMobile){
        return m("div.datetime-cpnt.mobile", {class: ""},
          m("i.icon-calendar", {}, ""),
          m("input", {
            type:"datetime-local",
            value: !attrs.value ? attrs.value : moment(attrs.value).format("YYYY-MM-DDTHH:mm:ss"),
            onchange:function(){
              if(this.value){
                var v = moment(this.value).format();
                attrs.bind && (attrs.bind.t[attrs.bind.p] = v);
              }
            }
          },"")
        );
      }

      return m("div.datetime-cpnt", {class: ""},
        m("i.icon-calendar", { }, ""),
        m("input", {type:"hidden", value: moment(attrs.value).format()}, ""),
        m("span.mui-date", { }, [ inp, ])
      );
    },
    oncreate: function(vnode) {
      if(isMobile()){return}

      var attrs = vnode.attrs || {};
      var readOnly = attrs.readonly || false;
      if (readOnly) { return }

      var currentLocaleData = moment.localeData();
      picker = new Pikaday({
        defaultDate : new Date(),
        field       : inp.dom,
        position    : "bottom right",
        onSelect    : function(e){
          attrs.value = e;
          if(attrs.bind){
            attrs.bind.t[attrs.bind.p] = e;
          }
          attrs.onchange && attrs.onchange(e);
          m.redraw()
        },
        i18n: {
          months : currentLocaleData.months(),
          weekdays : currentLocaleData.weekdays(),
          weekdaysShort : currentLocaleData.weekdaysShort()
        },
        format: currentLocaleData.longDateFormat("LL"),
      });
      picker.setDate(attrs.value || new Date())
    },
    onremove: function(opts) {
      if(isMobile()){return}
      var attrs = opts.attrs || {};
      var readOnly = attrs.readonly || false;
      if (readOnly) { return }
      picker.destroy()
    },
  }
}

function timeComponent() {

  return {
    view: function(opts){
      var attrs = opts.attrs || {};
      var readOnly = attrs.readonly || false;
      var value = attrs.value || attrs.bind.t[attrs.bind.p];
      if(!value) {
        value = new Date();
        value.setMinutes(00);
      }else if (!value.getHours) {
        value = moment(value).toDate();
      }
      hours = value.getHours();
      minutes = value.getMinutes();

      var out = new Date();
      out.setHours(hours);
      out.setMinutes(minutes);

      if(hours<10){hours="0"+hours;}
      if(minutes<10){minutes="0"+minutes;}

      if (readOnly) {
        return m("div.time-cpnt.readonly", {},
          m("i.icon-clock", { }, ""),
          m("label", {}, attrs.placeholder),
          m("span", {class:"hours"}, hours),
          m("span", {}, ":"),
          m("span", {class:"minutes"}, minutes)
        )
      }

      if(isMobile()){
        return m("div", {class:"time-cpnt"}, [
          m("i.icon-clock", { }, ""),
          m("input", {
            type:"time",
            value:hours+":"+minutes,
            onchange:function(){
              if(this.value){
                var x = this.value.split(":")
                hours = x[0];
                minutes = x[1];
                out.setHours(hours);
                out.setMinutes(minutes);
                attrs.bind && (attrs.bind.t[attrs.bind.p] = out);
              }
            }
          },"")
        ])
      }

      var hopts = Array(24).fill().map( function(g, i) { return m("option",{value:i},i)});
      var mopts = Array(60).fill().filter(function(g,i){return i%10===0;}).map( function(g, i) { return m("option",{value:i*10},i*10)})

      return m("div", {class:"time-cpnt"}, [
        m("i.icon-clock", { }, ""),
        m("select", {class:"hours", selectedIndex:hours, onchange: function(e){
          hours = e.target.value;
          out.setHours(hours);
          attrs.bind && (attrs.bind.t[attrs.bind.p] = out);
          attrs.onchange && attrs.onchange(out);
        }}, hopts ),
        m("span",{},":"),
        m("select", {class:"minutes", selectedIndex:minutes/10, onchange: function(e){
          minutes = e.target.value;
          out.setMinutes(minutes);
          attrs.bind && (attrs.bind.t[attrs.bind.p] = out);
          attrs.onchange && attrs.onchange(out);
        }}, mopts),
      ]);
    }
  }
}

var Text = function(){
  return {
    view: function(opts) {
      var attrs = opts.attrs || {};
      var readOnly = attrs.readonly || false;
      var value = attrs.value;
      if(!value && attrs.bind){
        value = attrs.bind.t[attrs.bind.p];
      }

      if (readOnly) {
        return m("div.txt-cpnt.readonly", {}, [
          m("label", {class: attrs.bind?attrs.bind.p:""}, attrs.placeholder),
          m("span", {class: attrs.bind?attrs.bind.p:""}, value),
        ])
      }

      var tagAttrs = {
        value: value,
        onkeyup: function(e){
          if(attrs.bind){
            attrs.bind.t[attrs.bind.p] = e.target.value;
          }
          attrs.value = e.target.value;
          attrs.onkeyup && attrs.onkeyup(e);
        },
        placeholder: attrs.placeholder,
        readonly: attrs.readonly,
      }
      if(attrs.bind){
        tagAttrs.name = attrs.bind.p;
      }
      if (attrs.tagn!="textarea") {
        tagAttrs.type = attrs.type || "text";
      }
      if(attrs.error) {
        tagAttrs.class = "error";
      }
      if (attrs.required) {
        tagAttrs.required = true;
      }
      return m(attrs.tagn || "input", tagAttrs)
    },
  }
}

var Checkbox = function(){
  return {
    view: function(opts) {
      var attrs = opts.attrs || {};

      var divAttrs = {
        class: "cb-cpnt",
      }
      var labelAttrs = {
        for: attrs.id,
      }
      var cbAttrs = {
        id: attrs.id,
        value: attrs.value,
        readonly: attrs.readonly,
        type: attrs.type || "checkbox",
        checked:false,
      }
      if(attrs.error) {
        divAttrs.class += " error";
      }
      if (attrs.required) {
        cbAttrs.required = true;
      }
      if (attrs.checked) {
        cbAttrs.checked = !!attrs.checked;
      }
      if(attrs.bind){
        cbAttrs.checked = !!attrs.bind.t[attrs.bind.p];
      }
      if (attrs.id) {
        cbAttrs.id = attrs.id;
        labelAttrs.for = attrs.id;
      }
      cbAttrs.onclick = function(e) {
        // attrs.checked = !attrs.checked;
        if(attrs.bind){
          attrs.bind.t[attrs.bind.p] = !cbAttrs.checked;
        }
        attrs.onclick && attrs.onclick(e, attrs.checked);
      }
      cbAttrs.onchange = function(e) {
        if(attrs.bind){
          attrs.bind.t[attrs.bind.p] = !cbAttrs.checked;
        }
        attrs.onchange && attrs.onchange(e, attrs.checked);
      }
      return m("div",divAttrs,[
        !attrs.labelFirst?null:m("label",labelAttrs,attrs.placeholder),
        m("input",cbAttrs,""),
        attrs.labelFirst?null:m("label",labelAttrs,attrs.placeholder),
      ])
    },
  }
}

function txtF(model, prop, text, errP, opts){
  var attrs = {};
  opts = opts || {};
  opts.placeholder = text;
  if (errP) {opts.error = errP(prop);}
  opts.bind= {t: model, p: prop};
  return m(Text, opts)
}
function txtA(model, prop, text, errP, opts){
  var attrs = {};
  opts = opts || {};
  opts.placeholder = text;
  if (errP) {opts.error = errP(prop);}
  opts.bind= {t: model, p: prop};
  opts.tagn="textarea"
  return m(Text, opts)
}
function dateF(model, prop, text, errP, opts){
  var attrs = {};
  opts = opts || {};
  opts.placeholder = text;
  if (errP) {opts.error = errP(prop);}
  opts.bind= {t: model, p: prop};
  return m(dateComponent, opts)
}
function cbF(model, prop, text, errP, opts){
  var attrs = {};
  opts = opts || {};
  opts.placeholder = text;
  if (errP) {opts.error = errP(prop);}
  opts.bind= {t: model, p: prop};
  return m(Checkbox, opts)
}

var Geocoder = function(opts){
  var attrs = opts.attrs || {};
  var id = attrs.id||"";
  var geoc = new L.Control.Geocoder.Nominatim({
    serviceUrl: "/geocode/"
  });
  var disabled=false;
  var tout;
  var value;
  var defaultText = "saisissez une recherche...";
  var emptyRes = "pas de resultats";
  var loadingRes = "...en cours de chargement...";
  var viewRes = "voir les résultats";
  var results = [ {name: emptyRes} ];
  var selectedIndex = 0;
  var selectedValue = null;
  function fetchResults(value){
    if(disabled) {
      return; // already processing.
    }
    disabled=true;
    selectedIndex = 0;
    m.redraw();
    geoc.geocode(value, function(res){
      results = results.slice(0,1);
      res.map(function(v){results.push(v);});
      results[0].name = results.length>1?viewRes:emptyRes;
      disabled=false;
      m.redraw();
    }, null)
  }
  return {
    view: function(opts) {
      var attrs = opts.attrs || {};

      var divAttrs = {
        class: "geocoder-cpnt",
        id:id
      }
      var inputAttrs = {
        placeholder: defaultText,
        class:"column column-60",
        style:{display:"inline-block","max-width":"65%"},
        type:"text",
        value:value,
        onkeyup:function(){
          clearTimeout(tout);
          value = this.value;
          results[0].name = emptyRes;
          if(value.length){
            results[0].name = loadingRes;
            if(value.length>=4) {
              tout = setTimeout(function(){
                fetchResults(value)
              }, 1050)
            }
          } else {
            results = results.splice(0,1);
          }
        }
      }
      if (disabled) {
        inputAttrs.disabled="disabled";
      }
      var btSearchAttrs = {
        class:"column column-20 float-right",
        style:{display:"inline-block","vertical-align":"bottom","padding": "0 1rem","max-width":"30%"},
        onclick:function(){
          if(value && value.length>=4) {
            fetchResults(value)
          }
        }
      }
      var selectAttrs = {
        disabled: results.length<2,
        selectedIndex: selectedIndex,
        onchange: function(e, i){
          selectedIndex = this.selectedIndex;
          selectedValue = this.selectedIndex==0?null:results[this.selectedIndex];
          attrs.onselect && attrs.onselect(selectedValue, selectedIndex);
        }
      }
      var btConfirmAttrs = {
        disabled: selectedIndex<1,
        style:{"display":"block","width":"100%"},
        onclick:function(){
          attrs.onconfirm && attrs.onconfirm(selectedValue)
          selectedValue = null;
          selectedIndex = 0;
          value = "";
          results = results.slice(0,1)
          results[0].name = "pas de resultats"
          m.redraw()
        }
      }
      var confirmText = attrs.confirm;
      return m("div",divAttrs,[
        m("input",inputAttrs,""),
        m("button",btSearchAttrs, m("i",{class:"demo-icon icon-search"},"")),
        m("select",selectAttrs,
          results.map(function(v, i){
            return m('option',{value:""}, v.name)
          })
        ),
        !confirmText ? null : m("button", btConfirmAttrs, confirmText),
        m("span",{},"license: Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright"),
      ])
    },
  }
}

var CopyComponent = function(){
  var copied = false;
  return {
    view: function(opts) {
      var attrs = opts.attrs || {};
      var value = attrs.value;

      var inputAttrs = {
        type: "text",
        // value: attrs.value,
        readonly: true,
      }
      var btnAttrs = {
        class:"icon-clipboard"
      }
      var c = "copy-cpnt";
      if (copied) {
        c += " copied";
      }
      return m("div", {class:c, onclick: function(){
        clipboardCopy(attrs.value);
        copied = true;
        m.redraw();
        setTimeout(function(){
          copied = false;
          m.redraw();
        },1000*5)
      }}, [
        m("input", inputAttrs),
        m("div", {class:"content"}, attrs.value),
        m("div", {class:"extra"}, [
          m("span", {},copied?"texte copié!":"cliquez pour copier"),
          m("i",{class:"icon-clipboard"})
        ]),
      ])
    },
  }
}

var ProtestsListComponent = function(){
  return {
    view: function(vnode) {
      var protests = vnode.attrs.protests || [];
      var onclick = vnode.attrs.onclick || null;
      return m("div", {class:"protests-cpnt"},
        m("table", {},
          m("tr", {},
            m("th", {}, "Public"),
            m("th", {}, "Titre"),
            m("th", {}, "Date"),
            m("th", {}, "Lieu")
        ),
          protests.map(function(v){
            var steps = v.steps || [];
            return m("tr", {
                onclick:function(e){
                  onclick && onclick(v);
                }
              },
              m("td", {}, v.public?"oui":"non"),
              m("td", {}, v.title),
              m("td", {}, moment(v.gather_at).format("LL")),
              m("td", {}, steps.length ? steps[0].place : "-")
            )
          })
        )
      )
    },
  }
}

var ContactsListComponent = function(){
  return {
    view: function(vnode) {
      var messages = vnode.attrs.messages || [];
      var onclick = vnode.attrs.onclick || null;
      return m("div", {class:"contacts-cpnt"},
        m("table", {},
          m("tr", {},
            m("th", {}, "Date"),
            m("th", {}, "Sujet"),
            m("th", {}, "Contact"),
            m("th", {}, "Message")
          ),
          messages.map(function(v){
            return m("tr", {
                onclick:function(e){
                  onclick && onclick(v);
                }
              },
              m("td", {}, moment(v.created_at).format("LL") ),
              m("td", {}, v.subject),
              m("td", {}, v.returnaddr),
              m("td", {}, v.body.substring(20)+"...")
            )
          })
        )
      )
    },
  }
}

var ContactViewComponent = function(){
  return {
    view: function(vnode) {
      var attrs = vnode.attrs || {};
      var message = attrs.message || [];
      var onreturn = attrs.onreturn || null;
      var ontrash = attrs.ontrash || null;
      return m("div", {class:"contact-cpnt"},
        m("div", {},
          m("span", {}, "Date"),
          moment(message.created_at).format("LL")
        ),
        m("div", {},
          m("span", {}, "Contact"),
          message.returnaddr
        ),
        m("div", {},
          m("span", {}, "Sujet"),
          message.subject
        ),
        m("div", {},
          m("span", {}, "Message"),
          m("br"),
          message.body
        ),
        m("br"),
        m("br"),
        m("div", {},
          m("button", {onclick: function(e){onreturn && onreturn(message, e);}, class:"bt-back"}, "retour"),
          m("button", {onclick: function(e){ontrash && ontrash(message,e);}, class:"bt-trash"}, m("i", {class:"icon-trash-empty"}, ""))
        )
      )
    },
  }
}

var ProtestComponent = function(){

  function doPrint(){
    document.body.classList.add("print")
    window.print();
    document.body.classList.remove("print")
    return false;
  }

  return {
    view: function(opts) {
      var attrs = opts.attrs || {};
      var readOnly = attrs.readonly || false;
      var model = attrs.model || {};
      var errModel = attrs.errModel || {};
      var protestError = errorReader("protest", errModel)

      if (!('public' in model)) {
        model.public = true;
      }

      var cbText = "Ce parcours est public:";

      return m("div", {class:"protest-cpnt"},
        !readOnly ? null : m("i", {class:"print-it icon-print", onclick:doPrint}, ""),
        txtF(model, "title", "Titre", protestError, {readonly:readOnly}),
        txtF(model, "organizer", "Organisateur", protestError, {readonly:readOnly}),
        dateF(model, "gather_at", "Date", protestError, {readonly:readOnly, isMobile: isMobile()}),
        txtF(model, "protest", "Mouvement", protestError, {readonly:readOnly}),
        m("br"),
        txtA(model, "description", "Description", protestError, {readonly:readOnly, style:{width:"450px", resize:"none"}}),
        m("br"),
        m("div", {class:"publicpath"},
          !readOnly ? null : m("b", {}, cbText ),
          !readOnly ? null : m("span", {}, model.public ? "oui": "non"),
          readOnly ? null :cbF(model, "public", cbText, protestError, {id:"publiccb", readonly:readOnly, labelFirst:true})
        ),
        readOnly ? null : m("div", {style:{display:model.public?"none":"block"}} ,
          txtF(model, "password", "Mot de passe pour accéder à ce parcours", protestError, {}),
        ),
        m("br")
      )
    },
  }
}

var MapComponent = function(opts){

  var mymap = null;
  // var provider = L.tileLayer.provider('OpenStreetMap');
  var provider = L.tileLayer('https://{s}.tile.openstreetmap.fr/osmfr/{z}/{x}/{y}.png', {
  	maxZoom: 20,
  	attribution: '&copy; Openstreetmap France | &copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  });
  var c = new L.Control.Coordinates();
  var startPt = [48.8630,  2.3456];

  var line = null;

  var dragend = function(step){
    return function(event){
      var position = event.target.getLatLng();
      step.lat = position.lat;
      step.lng = position.lng;
      m.redraw();
    }
  }

  var drawPath = function(model){
     var coords = [];
     model.steps.map(function(v){
       coords.push([v.lat, v.lng]);
       v.marker.addTo(mymap);
       v.marker.off('dragend');
       v.marker.on('dragend', dragend(v));
       // if(v.highlight){
       //   mymap.setView([v.lat, v.lng], 13);
       // }
     });
     if (coords.length>1){
       if (line === null) {
         line = L.polyline(coords);
         line.addTo(mymap);
       } else {
         line.setLatLngs(coords)
       }
     } else if (line !== null) {
       line.removeFrom(mymap)
       line = null;
     }
  }


  return {
    map: null,
    view: function(vnode) {
      var attrs = vnode.attrs || {};
      var model = attrs.model || {};
      var readOnly = attrs.readonly || false;
      if (mymap){
        model.updateMarkers(mymap, dragend, readOnly);
        drawPath(model);
        if(model.startPt) {
          mymap.setView(model.startPt, 13);
          delete model.startPt;
        }
      }
      return m("div", {class:"map-container"}, "")
    },
    oncreate: function(vnode) {
      var attrs = vnode.attrs || {};
      var model = attrs.model || {};
      var readOnly = attrs.readonly || false;
      mymap = L.map(vnode.dom).setView(startPt, 13);
      provider.addTo(mymap);
      c.addTo(mymap);
      this.map = mymap;

      mymap.on('click', function(e) {
        e.originalEvent.preventDefault();

        var t = e.originalEvent.target;
        if (t.classList.contains("leaflet-marker-icon")) {
          if(model.isHighlightedByIcon(t)) {
            model.unHighLightAll();
          } else {
            model.unHighLightAll();
            model.highlightByIcon(t);
          }
          m.redraw();
          return false;
        }

        if (model.hasHighLights()) {
          model.unHighLightAll()
          m.redraw();
          return
        }

        if(readOnly) return;

        model.addStep(e.latlng.lat, e.latlng.lng);
        model.updateMarkers(mymap, dragend, readOnly);
        drawPath(model);

        c.setCoordinates(e);
        m.redraw();
      });

      model.updateMarkers(mymap, dragend, readOnly);
      drawPath(model);
    },
    onremove: function(vnode) {
      var attrs = vnode.attrs || {};
      var model = attrs.model || {};
      model.steps.map(function(v){v.marker.remove();})
      mymap.remove();
    },
  }
}

var StepsComponent = function(){
  return {
    view: function(opts) {
      var attrs = opts.attrs || {};
      var readOnly = attrs.readonly || false;
      var model = attrs.model || {};
      var errModel = attrs.errModel || {};
      var stepsError = errorReader("protest.steps", errModel)
      var hasHL = model.hasHighLights();
      var stepsClass = "steps"
      if (hasHL) {
        stepsClass+=" hl"
      }
      return m("div", {class:"leftpanel"},
        m("div", {class:"bkgw"}, ""),
        m("div", {class:"leftpanel-container"},[
          m("h3", {}, "vos étapes", m("span", {}, "("+(model.wholeDistance().toFixed(2))+" km)")),
          model.steps.length<2 ? null : m(Checkbox, {
            placeholder:"afficher toutes les étapes",
            checked:hasHL?"":"checked",
            id:"show-steps",
            onclick:function(e, checked){
              checked && model.unHighLightAll();
            }
          }),
          model.steps.length>0 ? null : m("div",{
            onclick:scrollTo(document.getElementById("geocoder")),
          }, m.trust("cliquez sur la <b>carte</b> ou utilisez le <b>geocoder</b> pour ajouter de nouveaux points de rendez vous.")),
          m("div", {class:stepsClass}, model.steps.map(function(v, i) {
            var c = "step ";
            if(hasHL){
              if (!v.highlight) {
                c += " hide";
              } else {
                c += " hl";
              }
            }
            return m("div", {class: c}, [
              m("h5", {onclick: function(e) {
                if(!readOnly){
                  model.unHighLightAll()
                  model.hightlightAtIndex(i);
                  var step = model.steps[i];
                  model.startPt = [step.lat, step.lng];
                  m.redraw();
                }
              }},"étape n°"+(i+1)),
              m(timeComponent, {readonly:readOnly, placeholder:"Heure",  bind:{t:v, p:"gather_at"}, error:stepsError(i, "gather_at") }),
              m(Text, {readonly:readOnly, placeholder:"lieu", bind:{t:v, p:"place"}, error:stepsError(i, "place") }),
              m(Text, {readonly:readOnly, tagn:"textarea", placeholder:"details", bind:{t:v, p:"details"}, error:stepsError(i, "details") }),
              i<1?null:m("div", {class:"dist"}, [
                m("span",{}, (v.dist.toFixed(2))+" km depuis l'étape précédente"),
              ]),
              m("div", {class:"latlng"}, [
                m("span",{},"lat:"+(v.lat.toFixed(4))),
                m("span",{},"lng:"+(v.lng.toFixed(4))),
              ]),
              readOnly ? null : m("div", {style: {"text-align":"right"}}, m("button", { onclick:function(e){model.rmStepAtIndex(i);}, class: "rm-step", }, m("i", {class:"icon-trash-empty"}, "")) ),
            ])
          })),
        ])
      )
    },
  }
}
