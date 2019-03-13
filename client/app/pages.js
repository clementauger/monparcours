
var ftomate = document.getElementsByName("ftomate")[0].value;

function errorReader(prefix, model) {
  return function () {
    var args = Array.prototype.slice.call(arguments);
    var path = args.map(function(x){
      if (x === parseInt(x, 10)){
        return "["+x+"]";
      }
      return "."+x;
    }).join("");
    path = ""+prefix+""+path;
    return model[path] || "";
  }
}

function protestModel( ){
  this.id = null;
  this.title = "";
  this.protest = "";
  this.organizer = "";
  this.password = "";
  this.gather_at = null;
  this.public = true;
  this.steps = [];

  this.updateMarkers = function(map, dragend, ro) {
    var steps = this.steps;
    steps.map(function(v, i){
      if(!v.marker || !v.marker.addTo ){
        v.marker = L.marker([v.lat, v.lng], {bubblingMouseEvents:false, draggable:!ro});
        v.marker.on('dragend', dragend(v));
      }
      if (v.highlight) {
        v.marker.setIcon(yellowIcon)
      } else {
        v.marker.setIcon(blueIcon)
      }
      v.marker.addTo(map);
      v.dist = 0;
      if (i>0) {
        var from = steps[i-1].marker.getLatLng();
        var to = v.marker.getLatLng();
        v.dist = (from.distanceTo(to)/1000)
      }
    })
  };

  this.highlightByIcon = function(domIcon) {
    this.steps.map(function(v){
      if(v.marker && v.marker._icon==domIcon){
        v.highlight = true;
        v.marker.setIcon( yellowIcon );
      }
    })
  };
  this.isHighlightedByIcon = function(domIcon) {
    var i = false;
    this.steps.map(function(v){
      if(v.marker && v.marker._icon==domIcon){
        i = v.highlight;
      }
    })
    return i
  };
  this.hightlightAtIndex = function(i) {
    if(this.steps[i]) {
      this.steps[i].title = "hled";
      this.steps[i].highlight = true;
      this.steps[i].marker.setIcon(yellowIcon);
    }
  };
  this.hasHighLights = function() {
    return this.steps.filter(function(v){
      return v.highlight;
    }).map(function(v){
      return v.highlight;
    }).length>0;
  };
  this.unHighLightAll = function() {
    this.steps.map(function(v){
      v.marker.setIcon(blueIcon)
      v.highlight = false;
    })
  };
  this.addStep = function(lat, lng) {
    var step = {
      lat: lat,
      lng: lng,
      title: "",
      gather_at: new Date(),
      details: "",
      highlight:false,
      dist:0,
      marker: null,
    };
    this.steps.push(step);
    return step;
  };
  this.rmStepAtIndex = function(i) {
    if (this.steps[i]) {
      this.steps[i].marker.remove();
      this.steps.splice(i,1);
    }
  };
  this.wholeDistance = function(i) {
    var d = 0.0;
    this.steps.map(function(v){
      d += v.dist;
    })
    return d;
  };

  this.reset = function(){
    this.title = "";
    this.protest = "";
    this.description = "";
    this.organizer = "";
    this.gather_at = null;
    this.public = true;
    this.password = "";
    this.steps = [];
    this.steps.map(function(v){
      this.steps.push({
          lat: v.lat,
          lng: v.lng,
          place: v.place,
          gather_at: v.gather_at,
          details: v.details,
      })
    })
  }

  this.setData = function(data){
    this.id = data.id || 0;
    this.title = data.title || "";
    this.protest = data.protest || "";
    this.description = data.description || "";
    this.organizer = data.organizer || "";
    this.gather_at = data.gather_at || null;
    this.public = !!data.public;
    this.password = data.password || "";
    this.steps = [];
    var that = this;
    (data.steps || []).map(function(v){
      that.steps.push({
          lat: v.lat,
          lng: v.lng,
          place: v.place,
          gather_at: v.gather_at,
          details: v.details,
          dist: 0.0,
      })
    })
  }

}

var myModel = {
  RW: new protestModel(),
  RO: new protestModel(),
  selectedMessage: {},
}

var myAPI = {
  createProtest: function(data){
    return new Promise(function(resolve, reject) {
      m.request({
        method: "POST",
        url: "/protests/create",
        data: data,
        withCredentials: true,
        headers: { "X-tomate": ftomate }
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
  getProtestByID: function(id){
    return new Promise(function(resolve, reject) {
      m.request({
        method: "GET",
        url: "/protests/"+id,
        headers: { "X-tomate": ftomate },
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
  getProtestByIDAndPwd: function(id, password){
    return new Promise(function(resolve, reject) {
      m.request({
        method: "POST",
        url: "/protests/"+id,
        data: {pwd: password},
        headers: { "X-tomate": ftomate },
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
  getProtestsByAuthor: function(id){
    return new Promise(function(resolve, reject) {
      m.request({
        method: "GET",
        url: "/protests/by_author/"+id,
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
  searchProtests: function(data){
    return new Promise(function(resolve, reject) {
      m.request({
        method: "POST",
        url: "/protests/search",
        headers: { "X-tomate": ftomate },
        data: data
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
  login: function(model){
    return new Promise(function(resolve, reject) {
      m.request({
        method: "POST",
        url: "/admin/login",
        headers: { "X-tomate": ftomate },
        data: model,
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
  captchaNew: function(){
    return new Promise(function(resolve, reject) {
      m.request({
        url: "/captcha/new",
        headers: { "X-tomate": ftomate },
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
  createContact: function(data){
    return new Promise(function(resolve, reject) {
      m.request({
        method: "POST",
        url: "/contacts/create",
        headers: { "X-tomate": ftomate },
        data: data,
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
  getContacts: function(){
    return new Promise(function(resolve, reject) {
      m.request({
        method: "GET",
        url: "/contacts/list",
        headers: { "X-tomate": ftomate },
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
  delContact: function(id){
    return new Promise(function(resolve, reject) {
      m.request({
        method: "POST",
        url: "/contacts/delete/"+id,
        headers: { "X-tomate": ftomate }
      })
      .then(resolve)
      .catch(function(e){
        console.log(e);
        reject(e);
      })
    })
  },
}

function Head(menu){
  var open = false;
  var onbodyClick = function(){ open = false; m.redraw(); }
  return {
    view: function() {
      var path =  (location.hash || "#!"+Object.keys(menu)[0]).substr(2);
      var c = "menu";
      if (open) { c+= " open" }
      return [
        m("h1", {class: "title"}, "Mon parcours"),
        m("div", {class: c},
          m("div", {class: "bg"}, ""),
          m("i", {class: "icon icon-list-nested", onclick:function(e){ open=!open;m.redraw(); return false;}}, ""),
          m("dl", {},
            Object.keys(menu).map(function(key) {
              var c = key==path?"hl":"";
              return m("dt", m("a", {
                href: "#!"+key,class:c,
                onclick: function(e) {
                  open = false;
                }
              }, menu[key]))
            })
          )
        )
      ]
    },
    oncreate: function(){
      document.getElementById("body").addEventListener("click", onbodyClick, false)
    },
    onremove: function(){
      document.getElementById("body").removeEventListener("click", onbodyClick)
    }
  }
}

function Foot(menu){
  return {
    view: function() {
      var path =  (location.hash || "#!"+Object.keys(menu)[0]).substr(2);
      var c = "foot";
      return m("div", {class: c},
        m("dl", {},
          Object.keys(menu).map(function(key) {
            var c = key==path?"hl":"";
            return m("dt", m("a", {
              href: "#!"+key,class:c,
              onclick: function(e) {
                open = false;
              }
            }, menu[key]))
          })
        )
      )
    }
  }
}

var WelcomePage = {
    view: function() {
      var content  = document.getElementById("text-accueil").innerHTML;
      return m("div", {class:"column"}, m.trust(content))
    }
}

var CreatePage = function(){
  var model = myModel.RW;
  var errModel = {};
  var mymap = {};

  var geocMarker = L.marker([0, 0], {bubblingMouseEvents:false, draggable:false});

  return {
    view: function(opts) {

      var map = m(MapComponent, { model: model, readonly:false});
      var protest = m(ProtestComponent, {model: model, errModel:errModel, readonly:false});
      var geocoder = m(Geocoder, {
        id:"geocoder",
        onselect: function(v){
          if(v&&v.name){
            var coords = [v.properties.lat, v.properties.lon];
            geocMarker.setLatLng(coords);
            geocMarker.addTo(map.state.map);
            map.state.map.setView(coords, 13);
          } else{
            geocMarker.remove();
          }
          m.redraw();
        },
        confirm: "créer une étape",
        onconfirm: function(v){
          geocMarker.remove();

          var lat = parseFloat(v.properties.lat)
          var lon = parseFloat(v.properties.lon)
          var step = model.addStep(lat, lon);
          step.marker = L.marker([lat, lon], {bubblingMouseEvents:false, draggable:true});
          step.place = v.name;
          step.gather_at = new Date();

          m.redraw();
        }
      });
      var steps = m(StepsComponent, {model: model, errModel:errModel, readonly:false});

      var save = function(e){
        var btn = e.target;
        btn.disabled=true;
        myModel.RO.setData(model)
        myModel.RO.author_id = getCookie("rnd");
        e.preventDefault();
        errModel = {};
        myAPI.createProtest(myModel.RO)
          .then(function(data) {
            btn.disabled=false;
            errModel = {};
            model.reset()
            myModel.RO.reset()
            myModel.RO.setData(data)
            self.location = "#!/voir/" + data.id;
          })
          .catch(function(err) {
            btn.disabled=false;
            errModel = err.response || {};
            m.redraw();
          })
        return false;
      }

      return m("div", {class:"column page-create"}, [
        protest,
        geocoder,
        m("div", {style:{position:"relative"}}, [
          map, steps
        ]),
        m("button", {
          class:"button bt-save float-right",
          disabled: model.steps.length<1,
          onclick:save,
        }, "sauvegarder")
      ])
    },
  }
}

var ViewPage = function(){
  var model = myModel.RO;
  var pwdModel = {password:""}
  return {
    oninit:function(opts){
      var attrs = opts.attrs || {};
      myAPI.getProtestByID(attrs.id)
       .then(function(data) {
         model.steps.map(function(v){v.marker.remove();})
         model.setData(data);
         if(model.steps.length) model.startPt = [model.steps[0].lat,model.steps[0].lng];
         m.redraw();
       })
    },
    view: function(opts) {
      var attrs = opts.attrs || {};

      function checkPwd(e){
        myAPI.getProtestByIDAndPwd(attrs.id, pwdModel.password)
         .then(function(data) {
           model.setData(data);
           m.redraw()
         })
      }

      if (model.id==-1) {
        return m("div", {class:"column page-view"},
          m("div", {class:"require-pwd"},
            m("label", {}, ""),
            txtF(pwdModel, "password", "Mot de passe", null, {}),
            m("button",{ onclick:checkPwd }, "soumettre")
          )
        )
      }
      return m("div", {class:"column page-view"}, [
        m(CopyComponent, {value: location.toString()}),
        m(ProtestComponent, {model: model, readonly:true}),
        m("div", {style:{position:"relative"}}, [
          m(MapComponent, { model: model, readonly:true}),
          m(StepsComponent, {model: model, readonly:true})
        ]),
      ]);
    },
  }
}

var NoticePage = function(){
  var content  = document.getElementById("text-notice");
  return {
    view: function() {
      return m("div", {class:"column page-notice"}, m.trust(content.innerHTML))
    },
    oncreate: function() {
      var ifr = document.createElement("iframe");
      ifr.setAttribute("width", "560")
      ifr.setAttribute("height", "315")
      ifr.setAttribute("frameborder", "0")
      ifr.setAttribute("src", "https://www.youtube-nocookie.com/embed/brsJwdH9dso?start=314")
      ifr.setAttribute("allow", "accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture")
      ifr.setAttribute("allowfullscreen", "")

      document.querySelector(".page-notice #yt").appendChild(ifr);

      // .innerHTML = ''+
      // '<iframe width="560" height="315" frameborder="0" ' +
      // 'src="https://www.youtube-nocookie.com/embed/brsJwdH9dso?start=314" ' +
      // 'allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" ' +
      // 'allowfullscreen '
      // '></iframe>';
    },
  }
}

var MinesPage = function(){
  var model = [];
  var onclick = function(v){
    window.location = "#!/voir/" + v.id;
  };
  return {
    oninit:function(vnode){
    myAPI.getProtestsByAuthor(getCookie("rnd"))
     .then(function(data) {
       model.splice(0)
       data.map(function(v){
         model.push(v)
       })
       m.redraw()
     })
    },
    view: function() {
      var emptyRes =  m("div", {class:"empty"},
        m("h3", {}, "il n'y a pas de résultats")
      );
      var protests = m(ProtestsListComponent, {protests: model, onclick:onclick});
      return m("div", {class:"column page-mines"},
        model.length<1 ? null : protests,
        model.length>0 ? null :emptyRes
      )
    }
  }
}

var SearchPage = function(){
  var model = {
    title:"",
    protest:"",
    organizer:"",
    date_start:moment().subtract(7, 'days').toDate(),
    date_end:moment().add(14, 'days').toDate(),
    lat:0,
    log:0
  }
  var onclick = function(v){
    window.open("#!/voir/" + v.id);
  }
  var results = [];
  var hasSearched = false;
  return {
    view: function() {

      var geocoder = m(Geocoder, {
        id:"geocoder",
        onselect: function(v){
          model.lat = parseFloat(v.properties.lat);
          model.lng = parseFloat(v.properties.lon);
          m.redraw();
        }
      });
      var search = function(){
        hasSearched = true;
        var data = {};
        if (model.title) { data.title = model.title; }
        if (model.protest) { data.protest = model.protest; }
        if (model.organizer) { data.organizer = model.organizer; }
        if (model.date_start) { data.date_start = model.date_start; }
        if (model.date_end) { data.date_end = model.date_end; }
        if (model.lat) { data.lat = model.lat; }
        if (model.lng) { data.lng = model.lng; }
        myAPI.searchProtests(data)
         .then(function(data) {
           results.splice(0)
           data.map(function(v){
             results.push(v);
           })
           m.redraw()
         })
      }

      var emptyRes = m("div", {class:"empty"},
        m("h3", {}, "il n'y a pas de résultats")
      );

      var protests = m("div", {class:"results"},
        m(ProtestsListComponent, {protests: results, onclick:onclick})
      );

      return m("div", {class:"column page-search"},
        m("div", {}, "saisissez des critères pour effectuer une recherche"),
        txtF(model, "title", "Titre", null, {}),
        txtF(model, "organizer", "Organisateur", null, {}),
        txtF(model, "protest", "Mouvement", null, {}),
        dateF(model, "date_start", "Debut de période", null, {isMobile: isMobile()}),
        dateF(model, "date_end", "Fin de période", null, {isMobile: isMobile()}),
        geocoder,
        m("button", {class:"bt-save", onclick:search}, "rechercher"),
        hasSearched && results.length<1 ? emptyRes : null,
        results.length<1 ? null : protests
      )
    },
  }
}

var ContactPage = function(){
  var model = {
    returnaddr:"",
    subject:"",
    body:"",
    captchaid:"",
    captchasolution:"",
  }
  var errModel = {}
  var confirmed = false;
  var updateCaptcha = function(){
    return myAPI.captchaNew()
     .then(function(data) {
       model.captchaid=data.id;
       model.captchasolution="";
       m.redraw()
     })
  }
  return {
    oninit: updateCaptcha,
    view: function() {
      if (confirmed) {
        return m("div", {class:"column page-contact"},
          m("h2", {class:"ok"}, "Votre demande de contact est enregistrée, merci!")
        )
      }
      var contactError = errorReader("contactmessageinput.contactmessage", errModel)
      var captchaError = errorReader("contactmessageinput.captchainput", errModel)
      var send = function(e){
        var btn = e.target;
        btn.disabled=true;
        errModel = {};
        myAPI.createContact(model)
         .then(function(data) {
           errModel = {};
           confirmed = true;
           btn.disabled=false;
           updateCaptcha()
           m.redraw()
         })
         .catch(function(e){
           btn.disabled=false;
           errModel = e.response;
           updateCaptcha()
           m.redraw()
         })
      }
      return m("div", {class:"column page-contact"},
        m("div", {}, "saisissez votre demande de contact"),
        txtF(model, "returnaddr", "Adresse de retour", contactError, {}),
        txtF(model, "subject", "Sujet", contactError, {}),
        txtA(model, "body", "Message", contactError, {}),

        m("div", {}, [
          m("div", {}, "saisissez le captcha"),
          m("input", {type:"hidden", value: model.captchaid}, ""),
          m("img", {src: "/captcha/"+model.captchaid+".png"}, ""),
          txtF(model, "captchasolution", "Solution au captcha", captchaError, {}),
        ]),

        m("button", {class:"bt-save", onclick:send,}, "envoyer")
      )
    },
  }
}

var AdminPage = function(){
  var model = {
    key:"",
  }
  var errModel = {}
  var confirmed = false;
  return {
    view: function() {
      var contactError = errorReader("login", errModel);
      if (confirmed) {
        return m("div", {class:"column page-admin"}, [
          m("h3", {class:"ok"}, "Vous êtes connecté.")
        ])
      }
      var login = function(e){
        var btn = e.target;
        btn.disabled=true;
        errModel = {};
        confirmed = false;
        myAPI.login(model)
         .then(function(data) {
           confirmed = true;
           errModel = {};
           btn.disabled=false;
           m.redraw()
         })
         .catch(function(e){
           btn.disabled=false;
           if (e.response) errModel = e.response;
           m.redraw()
         });
      }
      return m("div", {class:"column page-admin"},
        m("div", {}, "saisissez votre clef"),
        txtA(model, "key", "Clef", contactError, {}),
        m("button", {class :"bt-save", onclick:login}, "connexion")
      )
    },
  }
}

var ContactsPage = function(){
  var messages = [];
  var onclick = function(v){
    myModel.selectedMessage = v;
    window.location = "#!/contacts/voir/" + v.id;
  };
  return {
    oninit:function(){
      myAPI.getContacts()
       .then(function(data) {
         messages.splice(0)
         data.map(function(v){
           messages.push(v)
         })
         m.redraw()
       })
    },
    view: function() {
      return m("div", {class:"column page-contacts"},
        m(ContactsListComponent, {messages: messages, onclick:onclick})
      )
    }
  }
}

var ContactViewPage = function(){
  var onreturn = function(v){
    window.location = "#!/contacts/";
  };
  var ontrash = function(v){
    myAPI.delContact(v.id)
     .then(function(data) {
       myModel.selectedMessage = {};
       onreturn()
     })
  };
  return {
    view: function() {
      return m("div", {class:"column page-viewcontact"},
        m(ContactViewComponent, {message: myModel.selectedMessage, onreturn:onreturn, ontrash:ontrash})
      )
    }
  }
}
