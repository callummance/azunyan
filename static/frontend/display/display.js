let queueSongHeight = 60;

//Speed at which songs move around
let maxVelocity = 4;
let deltaV = 0.2;

//Z-Orders
let vidZ = 10;

//Song display options
let songBgColour = 0x222222;
let songTextColour = 0xFFFFFF;
let songArtistColour = 0xAAAAAA;
let singerTextColour = 0xFFFFFF;

let nowPlayingBGColour = 0xAAAAAA;
let nowPlayingTitleColour = 0x000000;
let nowPlayingArtistColour = 0x444444;
let nowPlayingSingerColour = 0x000000;



(function() {
    DisplayClient();
}());


function DisplayClient() {
    this.queue = [];
    this.queuedisplay = {};
    this.partialqueuedisplay = {};
    this.cur = {};
    this.source = new window.EventSource('/api/queuestream');
    this.active = true;
    this.queueHasScrolled = false;
    this.noSingers = 0;

    var client = this;

    //Get the number of singers
    $.get({url: "/api/nosingers", async: false}).done(function(data) {
      client.noSingers = data;
    });
    //Listen for SSEs
    this.source.addEventListener('queue',  function(e) {
        client.queue = JSON.parse(e.data);
        client.queuedisplay = makeQueueDivs(client.queue, client.queuedisplay, client.cur, $('#queue'));
        client.partialqueuedisplay = makePartialQueueDivs(client.queue, client.partialqueuedisplay, client.cur, $('#waiting'), client.noSingers)
    });
    this.source.addEventListener('cur', function(e) {
        client.cur = JSON.parse(e.data);
        setNowPlaying(client.cur)
    });
    this.source.addEventListener('active', function(e) {
        client.active = JSON.parse(e.data);
        setActive(client.active)
    });
    this.source.addEventListener('message', function(e) {
        client.message.text = JSON.parse(e.data).message;
    });


    $("#queue").scroll(function() {
      client.queueHasScrolled = true;
    })
    setTimeout(function() {
      handleScrollQueue($("#queue"), client);
    }, 5000);
}

function handleScrollQueue(targetDiv, client) {
  let curScrollPosition = targetDiv.scrollTop();
  let contentHeight = targetDiv[0].scrollHeight;
  let divHeight = targetDiv.height();
  let viewportBottom = curScrollPosition + divHeight;

  if (client.queueHasScrolled) {
    client.queueHasScrolled = false;
    setTimeout(function() {
      handleScrollQueue($("#queue"), client);
    }, 2000);
  } else if (viewportBottom >= contentHeight) {
    //We are already at the bottom, scroll back up
    targetDiv.animate({
      scrollTop: 0
    });
    setTimeout(function() {
      handleScrollQueue($("#queue"), client);
    }, 10000);
  } else {
    //Scroll down to the next set of items
    targetDiv.animate({
      scrollTop: viewportBottom
    });
    setTimeout(function() {
      handleScrollQueue($("#queue"), client);
    }, 1000);
  }

}

function setActive(newState) {
    if (newState === true) {
        $('#overlay').css("visibility", "hidden")
    } else {
        $('#overlay').css("visibility", "visible")
    }
}

function setNowPlaying(nowPlaying) {
    let title = nowPlaying.title;
    let artist = nowPlaying.artist;
    let singers = nowPlaying.singers.join(", ");
    let sid = nowPlaying.sid;
    let coverImg = "/i/cover/" + sid

    $('#title').text(title);
    $('#artist').text(artist);
    $('#singers').text(singers);
    $('#cover').attr("src", coverImg);
}


function makeQueueDivs(queue, prevQueueDivs, nowPlaying, targetDiv) {
    let newQueueDisplay = {};
    queue.complete.map((q_entry, index) => {
        let newPos = index * 55;
        itemid = q_entry.ids.join(".");
        if (prevQueueDivs.hasOwnProperty(itemid)) {
            //Thingy is already in queue
            prevQueueDivs[itemid].animate({
                top: newPos
            }, 1000);
            if (nowPlaying.ids != q_entry.ids) {
                newQueueDisplay[itemid] = prevQueueDivs[itemid]
            }
        } else {
            //Make a new div
            let newdiv = $('<div/>')
                .attr("id", itemid)
                .addClass("queueitem");
            targetDiv.append(newdiv);

            //Add relevant text
            newdiv.append(
                $('<div>')
                    .addClass("itemtitle")
                    .append("h1")
                    .text(q_entry.title)
            );
            newdiv.append(
                $('<div>')
                    .addClass("itemartist")
                    .append("h2")
                    .text(q_entry.artist)
            );
            newdiv.append(
                $('<div>')
                    .addClass("itemsingers")
                    .append("h3")
                    .text(q_entry.singers.join(", "))
            );

            //Sort out positioning
            newdiv.css("top", $(window).height());
            newdiv.animate({
                top: newPos
            }, 1000);
            if (nowPlaying.ids != q_entry.ids) {
                newQueueDisplay[itemid] = newdiv
            }
        }
    });
    Object.keys(prevQueueDivs).map(divid => {
        if (!(divid in newQueueDisplay)) {
           $('#' + divid).fadeOut(() => {
                $('#' + divid).remove();
            });
        }
    });
    return newQueueDisplay;
}

function makePartialQueueDivs(queue, prevQueueDivs, nowPlaying, targetDiv, noSingers) {
    let newQueueDisplay = {};
    if (queue.partial != null) {
      queue.partial.map((q_entry, index) => {
          itemid = q_entry.ids.join(".");
          if (!prevQueueDivs.hasOwnProperty(itemid)) {
              //Make a new div
              let newdiv = $('<div/>')
                  .attr("id", itemid)
                  .addClass("partialItem");
              targetDiv.append(newdiv);

              let detailsDiv = $('<div/>')
                  .addClass("partialSongDetails");
              newdiv.append(detailsDiv);

              //Add relevant text
              detailsDiv.append(
                  $('<div>')
                      .addClass("partialitemtitle")
                      .append("h1")
                      .text(q_entry.title)
              );
              detailsDiv.append(
                  $('<div>')
                      .addClass("partialitemartist")
                      .append("h2")
                      .text(q_entry.artist)
              );
              newdiv.append(
                  $('<div>')
                      .addClass("partialitemsingerscount")
                      .append("h3")
              );
              if (nowPlaying.ids != q_entry.ids) {
                  newQueueDisplay[itemid] = newdiv
              }
          } else {
            targetDiv.append(prevQueueDivs[itemid])
            if (nowPlaying.ids != q_entry.ids) {
                newQueueDisplay[q_entry.ids] = prevQueueDivs[itemid]
            }
          }

        //Either way, update the singer count
        newQueueDisplay[itemid].find(".partialitemsingerscount")
            .text(q_entry.singers.length + "/" + noSingers)
      });
    }
    Object.keys(prevQueueDivs).map(divid => {
        if (!(divid in newQueueDisplay)) {
           $('#' + divid).fadeOut(() => {
                $('#' + divid).remove();
            });
        }
    });
    return newQueueDisplay;
}


//UTILITY FUNCTIONS
//------------------------------------------------------------------
function mergeQueues(arr1, arr2) {
    let res = [];
    while (arr1.length > 0 && arr2.length > 0) {
        if (arr1[0].id == arr2[0].id) {
            res.push(arr1.pop());
            arr2.pop();
        } else if (arr1[0].priority > arr2[0].priority) {
            res.push(arr1.pop());
        } else {
            res.push(arr2.pop());
        }
    }
    return res.concat(arr1).concat(arr2);
}

function calcOpacity(count, index, tick) {
    let val = count * ((Math.sin(tick / count + (2 * index * Math.PI) / count) + 1) / 2 - (count-1) / count) + 0.1;
    return count == 1 ? 1 : (val > 0 ? val : 0);
}
