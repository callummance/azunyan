(function() {
    let panel = new AdminPanel();
})();

function AdminPanel() {
    //Set sensible defaults
    this.queue = [];
    this.nowPlaying = {};
    this.active = true;

    //Save reference to this
    let admin = this;

    //Load and subscribe to event source
    this.source = new window.EventSource('/api/queuestream');
    this.source.addEventListener('queue', function(e) {
        admin.queue = JSON.parse(e.data);
    });
    this.source.addEventListener('cur', function(e) {
        admin.nowPlaying = JSON.parse(e.data);
        admin.setPlaying();
    });
    this.source.addEventListener('active', function(e) {
        admin.active = JSON.parse(e.data).active;
    });


    //Functions for updating display elements
    this.setPlaying = function() {
        $('#nowPlaying').text(admin.nowPlaying.songtitle + " - " + admin.nowPlaying.songartist);
    };

    //Setup the modal
    let modal = document.getElementById('updateModal');
    let btn = document.getElementById("updateActivate");
    let span = document.getElementsByClassName("close")[0];

    btn.onclick = function() {modal.style.display = "block";};
    span.onclick = function() {modal.style.display = "none";};
    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
    };

    //Setup the modal
    let qmodal = document.getElementById('queueModal');
    let qbtn = document.getElementById("queueActivate");
    let qspan = document.getElementsByClassName("close")[0];

    qbtn.onclick = function() {
        qmodal.style.display = "block";
        loadQueue("queueTable");
    };
    qspan.onclick = function() {qmodal.style.display = "none";};
    window.onclick = function(event) {
        if (event.target == qmodal) {
            qmodal.style.display = "none";
        }
    };

    //Listen for button presses
    $('#advance').click(function() {jQuery.post('/admin/advance');});
    $('#activate').click(function() {jQuery.ajax({
        type: "POST",
        url: '/admin/active',
        data: JSON.stringify({"active": true}),
        contentType: 'application/json'
    });});
    $('#deactivate').click(function() {jQuery.ajax({
        type: "POST",
        url: '/admin/active',
        data: JSON.stringify({"active": false}),
        contentType: 'application/json'
    });});
    $('#activater').click(function() {jQuery.ajax({
        type: "POST",
        url: '/admin/req_active',
        data: JSON.stringify({"active": true}),
        contentType: 'application/json'
    });});
    $('#deactivater').click(function() {jQuery.ajax({
        type: "POST",
        url: '/admin/req_active',
        data: JSON.stringify({"active": false}),
        contentType: 'application/json'
    });});
    $('#sendUpdate').click(function() {jQuery.ajax({
        type: "POST",
        url: '/admin/merge_songs',
        data: $('#jsonField').val(),
        contentType: 'application/json'
    });
        modal.style.display = "none";});
    $(document).keypress(function(e) {
        //Also advance on space bar
        if (e.which == 32) {jQuery.post('/admin/advance');}
    });
}

function loadQueue(id) {
    $.getJSON('/admin/get_queue', function(data) {
        let rows = [];
        $.each(data, function(i, obj) {
            let titleDiv = document.createElement('div');
            titleDiv.className = "cell";
            titleDiv.append(obj.title);

            let singersDiv = document.createElement('div');
            singersDiv.className = "cell";
            singersDiv.append(obj.singers.join(', '));

            let prioDiv = document.createElement('div');
            prioDiv.className = "cell";
            prioDiv.append(obj.priority);

            let modDiv = document.createElement('div');
            let modBox = document.createElement('input');
            modBox.type = "number";
            modBox.value = obj.mod;
            modBox.onblur = function() {
                let id = obj.queuePos;
                let newVal = modBox.value;
                jQuery.ajax({
                    type: "POST",
                    url: '/admin/update_basemod',
                    data: JSON.stringify({
                        id: id,
                        newMod: newVal
                    }),
                    contentType: 'application/json'
                });
            };
            modDiv.className = "cell";
            modDiv.appendChild(modBox);

            let delDiv = document.createElement('div');
            let delBut = document.createElement('button');
            delBut.innerText = "Remove";
            delBut.onclick = function() {
                let id = obj.queuePos;
                jQuery.ajax({
                    type: "DELETE",
                    url: '/admin/remove_queue',
                    data: JSON.stringify({id: id}),
                    contentType: 'application/json'
                });
            };
            delDiv.appendChild(delBut);

            let rowDiv = document.createElement('div');
            rowDiv.className = "row";
            rowDiv.appendChild(titleDiv);
            rowDiv.appendChild(singersDiv);
            rowDiv.appendChild(prioDiv);
            rowDiv.appendChild(modDiv);
            rowDiv.appendChild(delDiv);


            rows.push(rowDiv);
            console.log(data);
        });

        let target = document.getElementById(id);
        while (target.firstChild) {
            target.removeChild(target.firstChild);
        }
        target.innerHTML = `
            <div class="row header">
                <div class="cell">
                    Song Title
                </div>
                <div class="cell">
                    Singers
                </div>
                <div class="cell">
                    Prio
                </div>
                <div class="cell">
                    Mod
                </div>
                <div class="cell">
                    Del
                </div>
            </div>`;
        $.each(rows, function(i, obj) {
            target.appendChild(obj);
        });
    });
}