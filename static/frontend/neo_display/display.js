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
    this.cur = {};
    this.source = new window.EventSource('/api/queuestream');
    this.active = true;

    let client = this;

    //Listen for SSEs
    this.source.addEventListener('queue',  function(e) {
        client.queue = JSON.parse(e.data);
        client.queuedisplay = makeQueueDivs(client.queue, client.queuedisplay, client.cur);
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
}

function setActive(newState) {
    if (newState === true) {
        $('#overlay').css("visibility", "hidden")
    } else {
        $('#overlay').css("visibility", "visible")
    }
}

function setNowPlaying(nowPlaying) {
    let title = nowPlaying.songtitle;
    let artist = nowPlaying.songartist;
    let singers = nowPlaying.singers.join(", ");

    $('#title').text(title);
    $('#artist').text(artist);
    $('#singers').text(singers);
}


function makeQueueDivs(queue, prevQueueDivs, nowPlaying) {
    let newQueueDisplay = {};
    queue.map((q_entry, index) => {
        let newPos = index * 55;
        itemid = q_entry.reqid;
        if (prevQueueDivs.hasOwnProperty(itemid)) {
            //Thingy is already in queue
            prevQueueDivs[itemid].animate({
                top: newPos
            }, 1000);
            if (nowPlaying.reqid != q_entry.reqid) {
                newQueueDisplay[q_entry.reqid] = prevQueueDivs[itemid]
            }
        } else {
            //Make a new div
            let newdiv = $('<div/>')
                .attr("id", itemid)
                .addClass("queueitem");
            $('#queue').append(newdiv);

            //Add relevant text
            newdiv.append(
                $('<div>')
                    .addClass("itemtitle")
                    .append("h1")
                    .text(q_entry.songtitle)
            );
            newdiv.append(
                $('<div>')
                    .addClass("itemartist")
                    .append("h2")
                    .text(q_entry.songartist)
            );
            newdiv.append(
                $('<div>')
                    .addClass("itemsingers")
                    .append("h3")
                    .text(q_entry.singers)
            );

            //Sort out positioning
            newdiv.css("top", $(window).height());
            newdiv.animate({
                top: newPos
            }, 1000);
            if (nowPlaying.reqid != q_entry.reqid) {
                newQueueDisplay[q_entry.reqid] = newdiv
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