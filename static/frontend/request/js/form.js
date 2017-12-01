const awesomeplete = require("./awesomplete-gh-pages/awesomplete");

(function() {
    let input = document.getElementById("songSelect");
    let submitButton = document.getElementById("submit");
    let songSelectAuto = new awesomeplete(input);

    fetchSingersCount();
    setTitles(songSelectAuto);

    songSelectAuto.replace = setChoice;
    songSelectAuto.maxItems = 10;
    songSelectAuto.autoFirst = true;
    submitButton.addEventListener("click", listenForSubmit);
}());

function checkKey(e) {
    if(e.keyCode == 13){
        $("#submit").click();
    }
}

function setTitles(songSelect) {
    $.getJSON("/api/getsongslist", function(data, status, jqXHR) {
        list = [];
        for (let song of data) {
            let label = song.artist + " - " + song.title;
            let value = song.id;
            list.push({
                label: label,
                value: value,
            });
        }
        songSelect.list = list;
    });
}

function setChoice(text) {
    let songInput = document.getElementById("songSelect");
    let songChoice = document.getElementById("songSelection");
    let songChoiceId = document.getElementById("songSelectionId");

    songInput.value = text.label;
    songChoice.value = text.label;
    songChoiceId.value = text.value;
}

function listenForSubmit() {
    let songChoiceId = document.getElementById("songSelectionId");
    let songChoiceBox = document.getElementById("songSelect");
    let songChoice = document.getElementById("songSelection");
    if (songChoiceId.value < 0) {
        swal({
            title: "NotLikeThis",
            text: 'To ensure that we have the files for all of your requests, you must select a song from the list.\n' +
                  'Please use the text box to search available tracks, then click on your desired choice.',
            type: "error",
            confirmButtonText: "Return"
        });
    } else if (songChoice.value != songChoiceBox.value) {
        swal({
            title: "Are you sure you want to do this?",
            text: "To ensure that we have the files for all of your requests, you must select a song from the list.\n" +
                  "You have typed '" + songChoiceBox.value + "' in the song selection box, but '" +
                  songChoice.value + "' will be submitted. Are you sure you want to do this?",
            type: "warning",
            showConfirmButton: true,
            showCancelButton: true,
            closeOnConfirm: false,
            confirmButtonText: "I know what I'm doing!"
        }, function() {
            processSubmit(songChoiceId.value);
        });
    } else if (document.getElementById("singer1").value == "") {
        swal({
            title: "NotLikeThis",
            text: 'You need to request for at least one singer...',
            type: "error",
            confirmButtonText: "Return"
        }, function() {
            processSubmit(songChoiceId.value);
        });
    } else {
        processSubmit(songChoiceId.value);
    }
}

function processSubmit(id) {
    submitToServer();
    swal({
        title: "All Done!",
        text: "Your request has been submitted!\nThis message will disappear in 5 seconds",
        type: "success",
        timer: 5000,
        showConfirmButton: false
    }, function() {
        location.reload();
    });
}

function fetchSingersCount() {
    let dataStore = document.getElementById("noSingers");
    $.get("/api/nosingers", function(data, status, jqXHR) {
        dataStore.value = data;
        insertSingers(data);
    });
}

function insertSingers(count) {
    let singersDiv = document.getElementById("singers");
    for (var i = 0; i < count; i++) {
        let id = "singer" + i.toString();
        let pholder = "Singer " + i.toString();

        let elem = document.createElement("input");
        elem.id = id;
        elem.placeholder = pholder;
        elem.onkeyup = checkKey;
        singersDiv.appendChild(elem);
    }
}

function submitToServer() {
    let noSingers = document.getElementById("noSingers").value;
    let songId = document.getElementById("songSelectionId").value;
    let singers = [];
    for (let i = 0; i < noSingers; i++) {
        let elemId = "singer" + i.toString();
        let singerText = document.getElementById(elemId).value;
        singers.push(singerText);
    }

    let res = {
        songId : songId,
        singers : singers
    };
    console.log(res);
    $.post('/api/addrequest', res, function(resp) {
        return;
    }, 'json');
    console.log(res);
}