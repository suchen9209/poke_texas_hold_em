var socket;

$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host + '/ws/join?uname=' + $('#uname').text());
    // Message received on the socket
    socket.onmessage = function (event) {
        var data = JSON.parse(event.data);
        var li = document.createElement('li');

        console.log(data);

        switch (data.Type) {
        case 0: // JOIN
            if (data.User == $('#uname').text()) {
                li.innerText = 'You joined the chat room.';
            } else {
                li.innerText = data.User + ' joined the chat room.';
            }
            break;
        case 1: // LEAVE
            li.innerText = data.User + ' left the chat room.';
            break;
        case 2: // MESSAGE
            var username = document.createElement('strong');
            var content = document.createElement('span');

            username.innerText = data.User;
            content.innerText = data.Content;

            li.appendChild(username);
            li.appendChild(document.createTextNode(': '));
            li.appendChild(content);

            break;
        case 3://发牌
            let pos_str = "#pos" + data.Position;
            $(pos_str + " .user_name").html(data.User);
            let card_html = "<div class='card_info'>"+data.Card.Value+" "+ data.Card.Color+"</div>";
            $(pos_str + " .user_card").append(card_html)
            break;
        }

        $('#chatbox li').first().before(li);
    };

    // Send messages.
    var postConecnt = function () {
        var uname = $('#uname').text();
        var content = $('#sendbox').val();
        socket.send(content);
        $('#sendbox').val('');
    }

    var postMsg = function (message,type){
        let msg = {
            "Message" : message,
            "Type"    : type
        }
        let send_json = JSON.stringify(msg);
        socket.send(send_json);
    }

    $('#sendbtn').click(function () {
        postConecnt();
    });

    $('#start_game').click(function () {
        postMsg('start','game_op');
    });


});
