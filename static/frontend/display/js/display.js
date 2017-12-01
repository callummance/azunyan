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
    initPixi();
}());

function initPixi() {
    let windowHeight = $(window).innerHeight();
    let windowWidth = $(window).innerWidth();

    //Create renderer and add to html document
    let renderer = PIXI.autoDetectRenderer(
        windowWidth, windowHeight,
        {antialias: true, transparent: true, resolution: 1}
    );
    renderer.view.style.position = "absolute";
    renderer.view.style.display = "block";
    renderer.autoResize = true;
    renderer.resize(window.innerWidth, window.innerHeight);
    document.body.appendChild(renderer.view);

    //Create a stage then get the renderer to render it
    let stage = new PIXI.Container();
    renderer.render(stage);

    //Create client
    let client = new DisplayClient(stage);


    //Start the main loop
    renderLoop(stage, renderer, client);
}

function renderLoop(stage, renderer, client) {
    //Get ready for the next frame
    requestAnimationFrame(function() {renderLoop(stage, renderer, client);});

    //Carry out logic
    client.renderFrame(stage);

    //Render the stage
    renderer.render(stage);
}

function DisplayClient(stage) {
    this.stage = stage;
    this.queue = [];
    this.cur = {};
    this.source = new window.EventSource('/api/queuestream');
    this.active = true;
    this.playingIntro = false;
    this.displayQueue = new QueueDisplay(this);
    this.nowPlayingDisplay = new CurSong();
    this.message = new PIXI.Text("", {
        fontFamily: "sans-serif",
        fontSize: 25,
        align: "center",
        wordWrap: true
    });

    let client = this;

    //Load Resources
    PIXI.loader
        .add('spritesheet', '/images/dance_animation/dance.json')
        .add('displayLogo', '/images/disp_logo.png')
        .add('bubble', '/images/bubble.png')
        .load(function(loader, resources) {
            //Create and use animated sprite
            let danceTex = [];
            for (let i = 0; i < 6; i ++) {
                danceTex.push(PIXI.Texture.fromFrame('dance_00' + (i) + ".png"));
            }
            let mikuSprite = new PIXI.extras.AnimatedSprite(danceTex);
            client.stage.addChild(mikuSprite);
            mikuSprite.animationSpeed = 0.3;
            mikuSprite.anchor.set(0,1);
            mikuSprite.y = window.innerHeight;
            mikuSprite.play();

            //Add Logo
            let dispLogo = new PIXI.Sprite(resources.displayLogo.texture);
            dispLogo.anchor.set(0.5, 0);
            dispLogo.height = Math.floor((window.innerWidth/5)/dispLogo.width * dispLogo.height);
            dispLogo.width = Math.floor(window.innerWidth / 5);
            dispLogo.x = Math.floor(window.innerWidth/6);
            client.stage.addChild(dispLogo);

            //Add Speech Bubble and text
            let bubble = new PIXI.Sprite(resources.bubble.texture);
            bubble.anchor.set(0,1);
            bubble.height = mikuSprite.height;
            bubble.width = Math.floor(window.innerWidth/3) - (mikuSprite.width - 30);
            bubble.x = mikuSprite.width - 60;
            bubble.y = window.innerHeight - 6;
            client.stage.addChild(bubble);

            client.message.x = bubble.x + 80;
            client.message.y = bubble.getBounds().top + 20;
            client.message.style.wordWrapWidth = Math.floor(window.innerWidth/3) - (mikuSprite.width - 30) - 100;
            client.stage.addChild(client.message);

            //Add Now Playing display
            let nowPlayingDisplay = new PIXI.Container();
            nowPlayingDisplay.width = Math.floor(window.innerWidth/3) - 20;
            nowPlayingDisplay.height = window.innerHeight - (bubble.height + dispLogo.height + 40);
            nowPlayingDisplay.x = 5;
            nowPlayingDisplay.y = dispLogo.getBounds().bottom + 20;
            client.nowPlayingDisplay.setup(nowPlayingDisplay);
            client.stage.addChild(nowPlayingDisplay);

            //Make sure the intro video is still on top
            if (typeof client.introIndex != "undefined") {
                client.introIndex.bringToFront();
            }
        });

    //Listen for SSEs
    this.source.addEventListener('queue',  function(e) {
        client.queue = JSON.parse(e.data);
        client.displayQueue.updateQueue(client.queue);
    });
    this.source.addEventListener('current', function(e) {
        client.cur = JSON.parse(e.data);
        client.displayQueue.updateCur(client.cur);
        client.nowPlayingDisplay.updateSong(client.cur);
        client.nowPlayingDisplay.active = true;
    });
    this.source.addEventListener('active', function(e) {client.active = JSON.parse(e.data).active});
    this.source.addEventListener('message', function(e) {client.message.text = JSON.parse(e.data).message});



    this.renderFrame = function() {
        if ((!!client.queue.length || !!client.cur) && client.active) {
            //Clean up after intro screen
            if (client.playingIntro) {
                stage.removeChild(client.introIndex);
                client.playingIntro = false;
            }

            client.displayQueue.render();
            client.nowPlayingDisplay.show()

        } else {
            if (!this.playingIntro) {
                client.playingIntro = true;
                client.introIndex = displayWelcome(stage);
            }
        }
    };
}



function QueueDisplay(client) {
    this.queue = [];
    this.cur = null;

    this.containerWidth = Math.floor(2 * window.innerWidth/3);
    this.containerHeight = window.innerHeight - 4;
    this.containerX = Math.floor(window.innerWidth/3);
    this.containerY = 2;

    let thisQueue = this;

    let queueContainer = new PIXI.Container();
    queueContainer.width = thisQueue.containerWidth;
    queueContainer.height = thisQueue.containerHeight;
    queueContainer.position.x = thisQueue.containerX;
    queueContainer.position.y = thisQueue.containerY;
    client.stage.addChild(queueContainer);

    this.render = function() {
        for (let i = 0; i < thisQueue.queue.length; i++) {
            if (thisQueue.queue[i].update()) {
                thisQueue.queue.splice(i, 1);
                this.render();
                break;
            }
        }
    };

    this.updateQueue = function(newQueue) {
        let res = [];
        for (let qi = 0; qi < newQueue.length; qi++) {
            let matchingQueueElements = thisQueue.queue.filter(function(elem) {
                return elem.id == newQueue[qi].queuePos;
            });
            if (matchingQueueElements.length == 0) {
                res.push(new QueuedSong(newQueue[qi], qi, queueContainer));
            } else {
                matchingQueueElements[0].serverUpdate(qi);
                res.push(matchingQueueElements[0]);
            }
        }
        thisQueue.queue = mergeQueues(res, thisQueue.queue);
    };

    this.updateCur = function(newCur) {
        let res = [];
        let matchingQueueElements = thisQueue.queue.filter(function(elem) {
            return elem.id == newCur.queuePos;
        });
        if (matchingQueueElements.length == 0) {
            let newQueueElem = new QueuedSong(newCur, -1, queueContainer);
            newQueueElem.isTrash = true;
            res.push(newQueueElem);
        } else {
            matchingQueueElements[0].serverUpdate(-1);
            matchingQueueElements[0].isTrash = true;
            res.push(matchingQueueElements[0]);
        }
        thisQueue.queue = res.concat(thisQueue.queue);
    }
}

function CurSong() {
    let thisCur = this;
    this.active = false;

   //Request metadata
    this.updateSong = function(songObj) {
        thisCur.title = songObj.title;
        thisCur.artist = songObj.artist;
        thisCur.singers = songObj.singers.filter(function (singer) {
            return singer != "";
        });
    };

    this.show = function() {
        if (typeof thisCur.nowPlayingTitle != "undefined" &&
            typeof thisCur.nowPlayingArtist != "undefined" &&
            typeof thisCur.nowPlayingSingers != "undefined") {

            thisCur.parentContainer.alpha = this.active ? 1 : 0;
            thisCur.nowPlayingTitle.text = thisCur.title;
            thisCur.nowPlayingArtist.text = "by " + thisCur.artist;
            if (typeof thisCur.singers == "undefined" || thisCur.singers.length == 0) {
                thisCur.nowPlayingSingers.text = "DrakeGuildy and the Elemarians"
            } else if (thisCur.singers.length == 1) {
                thisCur.nowPlayingSingers.text = thisCur.singers[0];
            } else {
                let singerNames = "";
                for (let i = 0; i < thisCur.singers.length - 1; i ++) {
                    singerNames += thisCur.singers[i] + ", ";
                }
                singerNames = singerNames.slice(0, -2);
                singerNames += " & " + thisCur.singers[thisCur.singers.length - 1];
                thisCur.nowPlayingSingers.text = singerNames;
            }
            let offset = thisCur.parentContainer.getBounds().top;
            thisCur.nowPlayingTitle.y = thisCur.nowPlayingText.getBounds().bottom + 60 - offset;
            thisCur.nowPlayingArtist.y = thisCur.nowPlayingTitle.getBounds().bottom - offset;
            thisCur.nowPlayingBy.y = thisCur.nowPlayingArtist.getBounds().bottom + 50 - offset;
            thisCur.nowPlayingSingers.y = thisCur.nowPlayingBy.getBounds().bottom + 50 - offset;
        }
    };

    this.setup = function(parentContainer) {
        thisCur.parentContainer = parentContainer;
        thisCur.parentContainer.alpha = this.active ? 1 : 0;

        thisCur.textWidth = parentContainer._width - 10;
        thisCur.textLocX = Math.floor(parentContainer._width / 2);

        thisCur.bg = new PIXI.Graphics();
        thisCur.bg.beginFill(nowPlayingBGColour);
        thisCur.bg.drawRect(0, 0, parentContainer._width, parentContainer._height);
        thisCur.bg.endFill();

        thisCur.nowPlayingText = new PIXI.Text("Now Playing", {
            fontFamily: "sansserif",
            fontSize: 20,
            fill: nowPlayingTitleColour
        });
        thisCur.nowPlayingText.anchor.set(0.5, 0);
        thisCur.nowPlayingText.x = thisCur.textLocX;
        thisCur.nowPlayingText.y = 3;

        thisCur.nowPlayingTitle = new PIXI.Text("", {
            fontFamily: "sansserif",
            fontSize: 50,
            fill: nowPlayingTitleColour,
            align: "center",
            wordWrap: true,
            wordWrapWidth: thisCur.textWidth
        });
        thisCur.nowPlayingTitle.anchor.set(0.5, 0);
        thisCur.nowPlayingTitle.x = thisCur.textLocX;
        thisCur.nowPlayingTitle.y = thisCur.nowPlayingText.getBounds().bottom;

        thisCur.nowPlayingArtist = new PIXI.Text("", {
            fontFamily: "sansserif",
            fontSize: 20,
            fill: nowPlayingArtistColour,
            align: "center",
            wordWrap: true,
            wordWrapWidth: thisCur.textWidth
        });
        thisCur.nowPlayingArtist.anchor.set(0.5, 0);
        thisCur.nowPlayingArtist.x = thisCur.textLocX;
        thisCur.nowPlayingArtist.y = thisCur.nowPlayingTitle.getBounds().bottom;

        thisCur.nowPlayingBy = new PIXI.Text("sung by", {
            fontFamily: "sansserif",
            fontSize: 20,
            fill: nowPlayingArtistColour,
            align: "center",
            wordWrap: true,
            wordWrapWidth: thisCur.textWidth
        });
        thisCur.nowPlayingBy.anchor.set(0.5, 0);
        thisCur.nowPlayingBy.x = thisCur.textLocX;
        thisCur.nowPlayingBy.y = thisCur.nowPlayingArtist.getBounds().bottom + 50;

        thisCur.nowPlayingSingers = new PIXI.Text("sung by", {
            fontFamily: "sansserif",
            fontSize: 50,
            fill: nowPlayingSingerColour,
            align: "center",
            wordWrap: true,
            wordWrapWidth: thisCur.textWidth
        });
        thisCur.nowPlayingSingers.anchor.set(0.5, 0);
        thisCur.nowPlayingSingers.x = thisCur.textLocX;
        thisCur.nowPlayingSingers.y = thisCur.nowPlayingBy.getBounds().bottom + 50;


        parentContainer.addChild(thisCur.bg);
        parentContainer.addChild(thisCur.nowPlayingText);
        parentContainer.addChild(thisCur.nowPlayingTitle);
        parentContainer.addChild(thisCur.nowPlayingArtist);
        parentContainer.addChild(thisCur.nowPlayingBy);
        parentContainer.addChild(thisCur.nowPlayingSingers);
    };
}

function QueuedSong(songObj, position, parentContainer) {
    let thisSong = this;
    //Element metadata
    this.id = songObj.queuePos;
    this.title = songObj.title;
    this.artist = songObj.artist;
    this.singers = songObj.singers.filter(function(singer) {
        return singer != "";
    });
    this.priority = songObj.priority;
    this.parentContainer = parentContainer;

    this.isTrash = false;

    //Attributes to do with location and motion
    this.elemWidth = parentContainer._width - 2;
    this.elemHeight = queueSongHeight;
    this.targetYLoc = 1 + position * (queueSongHeight + 3);
    this.xLoc = 0;
    this.yLoc = window.innerHeight + 1 + this.targetYLoc;
    this.velocity = 0;
    this.tick = 0;

    //Display Elements
    this.songContainer = new PIXI.Container();

    //Updates position
    this.updatePos = function() {
        if (thisSong.yLoc == thisSong.targetYLoc) {
            thisSong.velocity = 0;
        } else if (thisSong.yLoc < thisSong.targetYLoc) {
            thisSong.velocity += (thisSong.velocity < maxVelocity) ? deltaV : 0;
            thisSong.yLoc += (thisSong.yLoc + thisSong.velocity > thisSong.targetYLoc) ? thisSong.targetYLoc - thisSong.yLoc : thisSong.velocity;
            thisSong.velocity = (thisSong.yLoc + thisSong.velocity > thisSong.targetYLoc) ? 0 : thisSong.velocity;
            return false;
        } else {
            thisSong.velocity += (thisSong.velocity < maxVelocity) ? 0 - deltaV : 0;
            thisSong.yLoc += (thisSong.yLoc + thisSong.velocity < thisSong.targetYLoc) ? thisSong.targetYLoc - thisSong.yLoc : thisSong.velocity;
            thisSong.velocity = (thisSong.yLoc + thisSong.velocity < thisSong.targetYLoc) ? 0 : thisSong.velocity;
            return false;
        }
    };

    this.update = function() {
        this.tick += 0.05;
        let res = thisSong.updatePos();
        thisSong.songContainer.position.x = thisSong.xLoc;
        thisSong.songContainer.position.y = thisSong.yLoc;
        for (let i = 0; i < thisSong.singers.length; i++) {
            thisSong.singers[i].alpha = calcOpacity(thisSong.singers.length, i, thisSong.tick);
        }
        return res;
    };

    this.serverUpdate = function(queuePos) {
        this.targetYLoc = 1 + queuePos * (queueSongHeight + 3);
    };

    this.setup = function() {
        thisSong.songContainer.width = thisSong.elemWidth;
        thisSong.songContainer.height = thisSong.elemHeight;

        thisSong.rect = new PIXI.Graphics();
        thisSong.rect.beginFill(songBgColour);
        thisSong.rect.drawRect(0, 0, thisSong.elemWidth, thisSong.elemHeight);
        thisSong.rect.endFill();

        thisSong.title = new PIXI.Text(thisSong.title,
            {
                fontFamily: 'Sans-Serif',
                fontSize: 30,
                fill: songTextColour
            });
        thisSong.title.x = 2;
        thisSong.title.y = 2;

        thisSong.artist = new PIXI.Text(' - ' + thisSong.artist,
            {
                fontFamily: 'Sans-Serif',
                fontSize: 15,
                fill: songArtistColour
            });
        thisSong.artist.x = 10;
        thisSong.artist.y = thisSong.title.getBounds().bottom + 1;


        thisSong.singers = thisSong.singers.map(function(singer) {
            let newSinger = new PIXI.Text(singer,
                {
                    fontFamily: 'Sans-Serif',
                    fontSize: 25,
                    fill: singerTextColour
                });
            newSinger.anchor.set(0, 0.5);
            newSinger.x = Math.floor(2 * thisSong.songContainer._width/3);
            newSinger.y = Math.floor(thisSong.songContainer._height/2);
            return newSinger;
        });

        thisSong.songContainer.addChild(thisSong.rect);
        thisSong.songContainer.addChild(thisSong.title);
        thisSong.songContainer.addChild(thisSong.artist);
        thisSong.songContainer.addChild(...thisSong.singers);

        thisSong.parentContainer.addChild(thisSong.songContainer);

        thisSong.update();
    };

    this.setup()
}

//UTILITY FUNCTIONS
//------------------------------------------------------------------
function displayWelcome(stage) {
    let welcomeScreen = new PIXI.Container();

    let vid = document.createElement('video');
    vid.preload = 'auto';
    vid.loop = true;
    vid.src = '/display/welcomebg.webm';
    let bg = PIXI.Texture.fromVideo(vid);
    let bgSprite = new PIXI.Sprite(bg);

    bgSprite.anchor.x = 0.5;
    bgSprite.anchor.y = 0.5;
    bgSprite.position.x = window.innerWidth/2;
    bgSprite.position.y = window.innerHeight/2;
    bgSprite.width = window.innerWidth;
    bgSprite.height = window.innerHeight;
    bgSprite.zOrder = vidZ;

    let img = PIXI.Texture.fromImage('/images/welcome.png');
    let imgSprite = new PIXI.Sprite(img);

    imgSprite.anchor.x = 0.5;
    imgSprite.anchor.y = 0.5;
    imgSprite.position.x = window.innerWidth/2;
    imgSprite.position.y = window.innerHeight/2;

    let fade = new PIXI.Graphics();
    fade.beginFill(0xFFFFFF, 0.5);
    fade.drawRect(0, 0, window.innerWidth, innerHeight);
    fade.endFill();


    welcomeScreen.addChild(bgSprite);
    welcomeScreen.addChild(fade);
    welcomeScreen.addChild(imgSprite);

    stage.addChild(welcomeScreen);
    welcomeScreen.bringToFront = function() {
        if (welcomeScreen.parent) {
            let parent = welcomeScreen.parent;
            parent.removeChild(welcomeScreen);
            parent.addChild(welcomeScreen);
        }
    };
    return welcomeScreen;
}

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