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
        if (e == null || e == undefined || e.data == "<nil>") {
            console.log("Nothing left in queue");
            admin.nowPlaying = {};
            $('#nowPlaying').text("Nothing left in queue");
            return;
        }
        try {
            admin.nowPlaying = JSON.parse(e.data);
            admin.setPlaying();
        } catch(error) {
            console.log(e);
            console.log(error);
        }
    });
    this.source.addEventListener('active', function(e) {
        admin.active = JSON.parse(e.data).active;
    });


    //Functions for updating display elements
    this.setPlaying = function() {
        $('#nowPlaying').text(admin.nowPlaying.title + " - " + admin.nowPlaying.artist);
    };

    //Setup the singer removal modal
    let remove_modal = document.getElementById('removeModal');
    let remove_btn = document.getElementById("removeActivate");
    let remove_span = document.getElementById("close_remove");

    remove_btn.onclick = function() {remove_modal.style.display = "block";};
    remove_span.onclick = function() {remove_modal.style.display = "none";};
    window.onclick = function(event) {
        if (event.target == remove_modal) {
            remove_modal.style.display = "none";
        }
    };

    //Setup the reset modal
    let reset_modal = document.getElementById('resetModal');
    let reset_btn = document.getElementById("resetActivate");
    let reset_span = document.getElementById("close_reset");

    reset_btn.onclick = function() {reset_modal.style.display = "block";};
    reset_span.onclick = function() {reset_modal.style.display = "none";};
    window.onclick = function(event) {
        if (event.target == reset_modal) {
            reset_modal.style.display = "none";
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
    $('#sendDelete').click(function() {jQuery.ajax({
          type: "POST",
          url: '/admin/remove_singer',
          data: JSON.stringify({"singer": $('#removeName').val()}),
          contentType: 'application/json'
      });
      remove_modal.style.display = "none";
    });
    $('#sendReset').click(function() {
        jQuery.post('/admin/reset_queue');
        reset_modal.style.display = "none";
    });
    $("input[name='singersInput'").change(() => {
        let value = $("input[name='singersInput'").val();
        jQuery.post('/admin/singers/'+value);
    });
        
    $(document).keypress(function(e) {
        //Also advance on space bar
        if (e.which == 32) {
          let remove_modal = document.getElementById('removeModal');
          if (remove_modal.style.display == "none") {
            jQuery.post('/admin/advance');
          }
        }
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
