jQuery(document).ready(function($){
  var searchBox = $("#songSelect");
  var resultsBox = $("#resultsbox");
  var searchResults = {};
  var latestResultSet = 0;
  var latestRecieved = 0;

  var searchTimeoutReady = true;

  var lazyInstance = $("#resultsBox .albumimage").Lazy({
    effect: "fadeIn",
    visibleOnly: true,
    chainable: false,
    autoDestroy: false,
    threshold: 0,
    appendScroll: $("#resultsbox"),
    beforeLoad: function(elem, resp) {
      console.log("Now loading element");
      console.log(elem);
    },
  });

  // Set an onclick handler on the results box
  $("#resultsbox").on("click","article", function(event) {
    selectedSong(event.currentTarget.id);
  });

  if (!String.prototype.trim) {
    String.prototype.trim = function () {
      return this.replace(/^[\s\uFEFF\xA0]+|[\s\uFEFF\xA0]+$/g, '');
    };
  }

  (function() {
    searchBox.on("keydown", function(e) {
      if (e.which == 13) {
        requestSearch(searchBox.val())
      }
    })
  }());

  function requestSearch(searchStr) {
    if (!searchTimeoutReady) {
      return;
    } else {
      searchTimeoutReady = false;
      $.getJSON("/api/searchsongs", {
        q: searchStr
      }, function(searchResults) {
        //Update results box
        displayResults(searchResults, latestResultSet++);
      });
      //Set timer until next allowed call in 0.5s
      setTimeout(function() {
        searchTimeoutReady = true
      }, 500);
    }
  }

  function displayResults(results, setNo) {
    newResults = {};
    if (setNo < latestRecieved) { return; }
    results.forEach(function(newSong) {
      if (newSong.ID in searchResults) {
        //We already have that song
        newResults[newSong.ID] = searchResults[newSong.ID];
        // searchResults[newSong.ID].detach();
      } else {
        //Create a new display object for the song

        let songCard = $(document.createElement("article"));
        songCard.addClass("resultcard");
        createResultCard(newSong, songCard);
        newResults[newSong.ID] = songCard;
      }
    });
    resultsBox.empty()
    results.forEach(function(result) {
      resultsBox.append(newResults[result.ID]);
    });
    searchResults = newResults;
    latestRecieved = setNo;
    lazyInstance.update();
  }

  function setLazy(img) {
    lazyInstance.addItems(img);
  }

  function createResultCard(result, card) {
    let imgurl = "/i/cover/" + result.ID;
    let imgobj = $(document.createElement("img"));
    imgobj.addClass("albumimage");
    imgobj.attr("data-src", imgurl);
    let countryobj = $(document.createElement("div"));
    countryobj.addClass("cardlang");
    countryobj.addClass("flag-icon-background");
    countryobj.addClass("flag-icon-" + getFlagName(result.language));
    let infoobj = $(document.createElement("div"));
    infoobj.addClass("cardinfo");
    let titleobj = $(document.createElement("div"));
    titleobj.text(result.title);
    titleobj.addClass("cardtitle");
    let artistobj = $(document.createElement("div"));
    artistobj.text(result.artist);
    artistobj.addClass("cardartist");

    infoobj.append(titleobj);
    infoobj.append(artistobj);

    card.append(imgobj);
    card.append(countryobj);
    card.append(infoobj);
    card.attr("id", result.ID);

    setLazy(imgobj);
  }

  function selectedSong(id){
    let card = searchResults[id].clone(false, false);
    card.addClass("pulse");
    card = $("<div>").append(card);
    let content = $("<div>");
    let textobj = $("<div>");
    textobj.text("Who wants to sing this song?");
    textobj.addClass("modaltextdiv");
    let namebox = $("<input>");
    namebox.attr("placeholder", "Enter your name");
    namebox.addClass("modalnamebox");
    namebox.attr("id", "nameinput");
    content.append(textobj);
    content.append(namebox);
    let modal = new tingle.modal({
      footer: true,
      stickyFooter: false,
      closeMethods: ["overlay", "button", "escape"],
      closeLabel: "Close",
      cssClass: ["modal"]
    })

    modal.setContent(card.append(content).html());
    modal.addFooterBtn("Let's Go!", "tingle-btn tingle-btn--primary tingle-btn--pull-right", function() {
      name = $("#nameinput").val();
      if (name.trim() == "") {
        alert("Please enter a name.\n\nThis version of the karaoke queue uses names to match singers up, as well as to prevent the queue being clogged with requests from people that have left.\nThank you for your understanding.");
      } else {
        modal.close()
        submitSelection(id, name, function() {
        });
      }
      console.log(name + " has requested " + id);
    })
    modal.addFooterBtn("Go back", "tingle-btn tingle-btn--danger tingle-btn--pull-right", function() {
      modal.close();
    });
    modal.open();
  }

  function submitSelection(sid, name, success) {
    let requestObj = {
      songid: sid,
      singer: name
    };
    $.ajax({
      url: "/api/addrequest",
      type: "POST",
      data: JSON.stringify(requestObj),
      contentType: "application/json; charset=utf-8",
    })
    .done(function(data, textStatus) {
      alert("Request Submitted!")
      success();
    })
    .fail(function(data, textStatus, errorThrown) {
      alert("Something went wrong: " + errorThrown)
    })
  }

  function getFlagName(lang) {
    var res = "un"
    switch (lang) {
      case "English":
        res = "gb";
        break;
      case "Japanese":
        res = "jp";
        break;
      case "German":
        res = "de";
        break;
      case "Español":
        res = "es";
        break;
      case "Polish":
        res = "pl";
        break;
      case "Korean":
        res = "kr";
        break;
      case "Swedish":
        res = "se";
        break;
      case "Chinese":
        res = "cn";
        break;
      case "Finnish":
        res = "fi";
        break;
    }
    return res;
  }

  $.fn.textfill = function(maxFontSize) {
    maxFontSize = parseInt(maxFontSize, 10);
    return this.each(function(){
      var ourText = $("span", this),
        parent = ourText.parent(),
        maxHeight = parent.height(),
        maxWidth = parent.width(),
        fontSize = parseInt(ourText.css("fontSize"), 10),
        multiplier = maxWidth/ourText.width(),
        newSize = (fontSize*(multiplier-0.1));
      ourText.css(
        "fontSize",
        (maxFontSize > 0 && newSize > maxFontSize) ?
          maxFontSize :
          newSize
        );
    });
  };

});
