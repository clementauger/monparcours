const Nightmare = require('nightmare');
const moment = require('moment');
const chai = require('chai');
const expect = chai.expect;

const key = process.env.AKEY;

var nm = function(){
  return Nightmare({
    show: !true,
    openDevTools: !true,
    webPreferences: {
        images: false
    },
    typeInterval: 25
  })
}

describe('test monparcours navigation', () => {
  it('should navigate using menu', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
      .goto('http://localhost:8100')
      .wait(".menu dl dt a")
      .wait(".foot dl dt a")
      .wait(".menu dl dt a.hl")
      .click(".foot dl dt:nth-child(1) a")
      .wait(".foot dl dt a.hl")
      .evaluate(() => location.hash)
      .end()
      .then(hash => {
        expect(hash).to.equal('#!/contact')
        done()
      }).catch(console.log)
  })
  it('should navigate using menu', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
      .goto('http://localhost:8100')
      .wait(".menu dl dt a.hl")
      .click(".menu dl dt:nth-child(1) a")
      .wait(".menu dl dt a.hl")
      .evaluate(() => location.hash)
      .end()
      .then(hash => {
        expect(hash).to.equal('#!/accueil')
        done()
      }).catch(console.log)
  })
})


function gotoCreatePage(nightmare){
  return nightmare
    .goto('http://localhost:8100')
    .wait(".menu dl dt a")
    .wait(".foot dl dt a")
    .wait(".menu dl dt a.hl")
    .click(".menu dl dt a[href='#!/creer']")
}
function fillCreateForm(nightmare){
  return nightmare
    .wait(500)
    .type(".page-create input[name='title']", "titre")
    .type(".page-create input[name='organizer']", "organizer")
    .click(".page-create .bt-calendar")
    .click(".pika-lendar .pika-row td:nth-child(2) .pika-day")
    .type(".page-create input[name='protest']", "protest")
    .type(".page-create textarea[name='description']", "description")
    .click(".page-create .map-container .leaflet-tile-container")
    .wait(500)
    .type(".page-create .step:nth-child(1) input[name='place']", "place")
    .select(".page-create .step:nth-child(1) .time-cpnt .hours", "5")
    .select(".page-create .step:nth-child(1) .time-cpnt .minutes", "30")
    .type(".page-create .step:nth-child(1) textarea[name='details']", "details")
}

describe('test monparcours create', () => {

  it('should navigate to create page', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    gotoCreatePage(nightmare)
      .evaluate(() => location.hash)
      .end()
      .then(hash => {
        expect(hash).to.equal('#!/creer')
        done()
      }).catch(console.log)
  })

  it('should display the create page with a disabled button at startup', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
      .exists(".page-create .bt-save:disabled")
      .end()
      .then((r) => {
        expect(r).to.equal(true)
        done()
      }).catch(console.log)
  })

  it('should enable the create page once a location exists', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
      .click(".page-create .map-container .leaflet-tile-container")
      .wait(100)
      .exists(".page-create .bt-save:disabled")
      .end()
      .then((r) => {
        expect(r).to.equal(false)
        done()
      }).catch(console.log)
  })

  it('should return error if the form is incomplete', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
      .wait(".page-create .map-container .leaflet-tile-container")
      .click(".page-create .map-container .leaflet-tile-container")
      .click(".page-create .bt-save")
      .wait(".page-create input.error")
      .exists(".page-create input.error")
      .end()
      .then((r) => {
        expect(r).to.equal(true)
        done()
      }).catch(console.log)
  })

  it('should save the protest', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view")
      .evaluate(() => location.hash)
      .end()
      .then((r) => {
        expect(r).to.have.string('/voir/');
        done()
      }).catch(console.log)
  })

  it('should set the copy content value', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view .copy-cpnt .content")
      .wait(".foot")
      .evaluate(() => {
        return [location.toString(),document.querySelector(".page-view .copy-cpnt .content").innerHTML]
      })
      .end()
      .then((v) => {
        expect(v[0]).to.eq(v[1]);
        done()
      })
  })

  it('should set the title value', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view span.title")
      .evaluate(() => {
        return document.querySelector(".page-view span.title").innerHTML
      })
      .end()
      .then((v) => {
        expect(v).to.eq("titre");
        done()
      })
  })

  it('should set the organizer value', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view span.organizer")
      .evaluate(() => {
        return document.querySelector(".page-view span.organizer").innerHTML
      })
      .end()
      .then((v) => {
        expect(v).to.eq("organizer");
        done()
      })
  })

  it('should set the protest value', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view span.protest")
      .evaluate(() => {
        return document.querySelector(".page-view span.protest").innerHTML
      })
      .end()
      .then((v) => {
        expect(v).to.eq("protest");
        done()
      })
  })

  it('should set the description value', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view span.description")
      .evaluate(() => {
        return document.querySelector(".page-view span.description").innerHTML
      })
      .end()
      .then((v) => {
        expect(v).to.eq("description");
        done()
      })
  })

  it('should set the gather_at value', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view span.gather_at")
      .evaluate(() => {
        return document.querySelector(".page-view span.gather_at").innerHTML
      })
      .end()
      .then((v) => {
        expect(v).to.match(/\d+\s\w+\s\d+/)
        done()
      })
  })

  it('should set the step', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view .steps .step")
      .evaluate(() => {
        return document.querySelectorAll(".page-view .steps .step").length
      })
      .end()
      .then((v) => {
        expect(v).to.eq(1)
        done()
      })
  })

  it('should set the step place', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view .steps .step")
      .evaluate(() => {
        return document.querySelector(".page-view .steps .step .place").innerHTML
      })
      .end()
      .then((v) => {
        expect(v).to.eq("lieu")
        done()
      })
  })

  it('should set the step details', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view .steps .step")
      .evaluate(() => {
        return document.querySelector(".page-view .steps .step .details").innerHTML
      })
      .end()
      .then((v) => {
        expect(v).to.eq("details")
        done()
      })
  })

  it('should set the step time', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view .steps .step")
      .evaluate(() => {
        return document.querySelector(".page-view .steps .step .hours").innerHTML + ":"+
        document.querySelector(".page-view .steps .step .minutes").innerHTML;
      })
      .end()
      .then((v) => {
        expect(v).to.eq("05:30")
        done()
      })
  })

  it('should record private protest', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create #publiccb")
      .wait(500)
      .type(".page-create input[name=password]", "testpwd")
      .click(".page-create .bt-save")
      .wait(".page-view .require-pwd")
      .type(".page-view .require-pwd input", "testpwd")
      .click(".page-view .require-pwd button")
      .wait(".page-view span.protest")
      .evaluate(() => {
        return document.querySelector(".page-view span.protest").innerHTML;
      })
      .end()
      .then((v) => {
        expect(v).to.eq("protest")
        done()
      })
  })
})

describe('test monparcours mines', () => {

  it('should display the correct number of protests', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view")
      .click(".menu dl dt:nth-child(3) a")
      .wait(".page-mines")
      .evaluate(() => document.querySelectorAll(".page-mines table tr").length)
      .end()
      .then((r) => {
        expect(r).to.eq(2);
        done()
      })
  })

  it('should display the correct values', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare = gotoCreatePage(nightmare)
      .wait(".page-create .bt-save")
    fillCreateForm(nightmare)
      .click(".page-create .bt-save")
      .wait(".page-view")
      .click(".menu dl dt:nth-child(3) a")
      .wait(".page-mines")
      .evaluate(() => document.querySelector(".page-mines table tr:nth-child(2) td:nth-child(2)").innerHTML)
      .end()
      .then((r) => {
        expect(r).to.eq("titre");
        done()
      })
  })

})

describe('test monparcours search', () => {

  it('should display search page', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare = gotoCreatePage(nightmare)
      .click(".menu dl dt:nth-child(4) a")
      .wait(".page-search .bt-save")
      .click(".page-search .bt-save")
      .wait(".page-search .results")
      .evaluate(() => document.querySelectorAll(".page-search .protests-cpnt table tr").length)
      .end()
      .then((r) => {
        expect(r).to.be.above(1);
        done()
      })
  })

  it('should display empty search results', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare = gotoCreatePage(nightmare)
      .click(".menu dl dt:nth-child(4) a")
      .wait(".page-search .bt-save")
      .type(".page-search input[name=title]", "pastitre")
      .click(".page-search .bt-save")
      .wait(".page-search .empty")
      .end()
      .then(() => {
        done()
      })
  })

})

describe('test monparcours contact', () => {

  it('should display contact page', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare = gotoCreatePage(nightmare)
      .click(".foot dl dt:nth-child(1) a")
      .wait(".page-contact")
      .click(".page-contact .bt-save")
      .wait(".page-contact input[name=returnaddr].error")
      .end()
      .then(() => {
        done()
      })
  })

  it('should submit contact page', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare = gotoCreatePage(nightmare)
      .click(".foot dl dt:nth-child(1) a")
      .wait(".page-contact")
      .type(".page-contact input[name=returnaddr]", "returnaddr")
      .type(".page-contact input[name=subject]", "subject")
      .type(".page-contact textarea[name=body]", "body")
      .type(".page-contact input[name=captchasolution]", "magic")
      .click(".page-contact .bt-save")
      .wait(".page-contact h2.ok")
      .end()
      .then(() => {
        done()
      })
  })

})

describe('test monparcours admin', () => {

  it('should display admin menu', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare.goto('http://localhost:8100?ia=y')
      .wait(".menu a[href='#!/admin']")
      .click(".menu a[href='#!/admin']")
      .wait(".page-admin .bt-save")
      .end()
      .then(() => {
        done()
      })
  })

  it('should login', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare.goto('http://localhost:8100?ia=y')
      .wait(".menu a[href='#!/admin']")
      .click(".menu a[href='#!/admin']")
      .wait(".page-admin .bt-save")
      .type(".page-admin textarea[name=key]", key)
      .click(".page-admin .bt-save")
      .wait(".page-admin h3.ok")
      .end()
      .then(() => {
        done()
      })
  })

  it('should read contact messages', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare.goto('http://localhost:8100?ia=y')
      .wait(".menu a[href='#!/admin']")
      .click(".menu a[href='#!/admin']")
      .wait(".page-admin .bt-save")
      .type(".page-admin textarea[name=key]", key)
      .click(".page-admin .bt-save")
      .wait(".foot a[href='#!/contacts']")
      .click(".foot a[href='#!/contacts']")
      .wait(".page-contacts .contacts-cpnt table tr")
      .wait(10*5)
      .evaluate(() => document.querySelectorAll(".page-contacts .contacts-cpnt table tr").length)
      .end()
      .then((r) => {
        expect(r).to.be.above(1);
        done()
      })
  })

  it('should delete contact messages', function(done) {
    this.timeout('10s')

    var nightmare = nm()
    nightmare.goto('http://localhost:8100?ia=y')
      .wait(".menu a[href='#!/admin']")
      .click(".menu a[href='#!/admin']")
      .wait(".page-admin .bt-save")
      .type(".page-admin textarea[name=key]", key)
      .click(".page-admin .bt-save")
      .wait(".foot a[href='#!/contacts']")
      .click(".foot a[href='#!/contacts']")
      .wait(".page-contacts .contacts-cpnt table tr")
      .wait(200)
      .click(".page-contacts .contacts-cpnt table tr:nth-child(2)")
      .wait(".page-viewcontact .bt-trash")
      .click(".page-viewcontact .bt-trash")
      .wait(200)
      .evaluate(() => document.querySelectorAll(".page-contacts .contacts-cpnt table tr").length)
      .end()
      .then((r) => {
        expect(r).to.eq(1);
        done()
      })
  })

})
