
var getNavigatorLanguage = function() {
  if (navigator.languages && navigator.languages.length) {
    return navigator.languages[0];
  } else {
    return navigator.userLanguage || navigator.language || navigator.browserLanguage || 'en';
  }
}
moment.locale( getNavigatorLanguage() );

if (getCookie("rnd")===null){
  setCookie("rnd", uuid(), 30)
} else {
  setCookie("rnd", getCookie("rnd"), 30)
}


var q = location.search ? parseQuery(location.search) : parseQuery(location.hash.substring(location.hash.indexOf("?")));
var isadmin = q.ia || getCookie("ia");
if(isadmin) {
  setCookie("ia", "true", 30)
}


var menuHead = {
  "/accueil": "Accueil",
  "/creer":"Créer un parcours",
  "/mes-parcours":"Mes parcours",
  "/rechercher":"Rechercher",
  "/avant-de-partir":"Avant de partir..."
};
if(isadmin) {
  menuHead["/admin"] = "Admin";
}
var menuFoot = {
  // "/a-propos":"A propos",
  "/contact":"Contact"
};
if(isadmin) {
  menuFoot["/contacts"] = "Contacts";
}

var pages = {
  "/accueil": WelcomePage,
  "/creer": CreatePage,
  "/mes-parcours": MinesPage,
  "/rechercher": SearchPage,
  "/voir/:id": ViewPage,
  "/avant-de-partir": NoticePage,
  "/contact": ContactPage
}
if(isadmin) {
  pages["/admin"] = AdminPage;
  pages["/contacts"] = ContactsPage;
  pages["/contacts/voir/:id"] = ContactViewPage;
}

m.mount(document.getElementById("head"), Head(menuHead))
m.mount(document.getElementById("foot"), Foot(menuFoot))
m.route(document.getElementById("body"), "/accueil", pages)
