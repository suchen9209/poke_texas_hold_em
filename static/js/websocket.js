var socket;
var my_position;
var uname;

function get_card_html(color,value){
    let html = '<div class="card">'+
    '<div class="face front puker-'+color+value+'" title="puker-spade6"></div>'+
    '</div>';
    return html;
    
}

function show_op_button(){
    console.log("show button");
}

$(document).ready(function () {
    // Create a socket
    uname = $('#uname').text()
    socket = new WebSocket('ws://' + window.location.host + '/ws/join?uname=' + uname);
    // Message received on the socket
    socket.onmessage = function (event) {
        var data = JSON.parse(event.data);
        var li = document.createElement('li');
        let pos_str;

        console.log(data);

        switch (data.Type) {
        case 0: // JOIN
            if (data.User == uname) {
                //You joined the chat room.
                my_position = data.GameUser.Position
            }
            // pos_str = "#pos" + data.GameUser.Position;
            // $(pos_str + " .user_name").html(data.User);

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
            pos_str = "#pos" + data.Position;
            // $(pos_str + " .user_name").html(data.User);
            let card_html = get_card_html(data.Card.Color,data.Card.Value);
            $(pos_str + " .user_card").append(card_html)
            break;
        case 4://公共牌
            let public_card_html = get_card_html(data.Card.Color,data.Card.Value);
            $("#public_card_table").append(public_card_html)
            break;
        case 5://清理场面上的牌
            $(".user_card").html("")
            $("#public_card_table").html("")
            break;
        case 6://更新用户信息
            let info = data.Info;
            for (const key in info) {
                if (Object.hasOwnProperty.call(info, key)) {
                    const element = info[key];
                    let pos_str = "#pos" + info[key].Position;
                    $(pos_str + " .user_name").html(info[key].Name);
                    $(pos_str + " .user_point").html(info[key].Point);
                }
            }
            break;
        case 7://回合信息
            console.log(data)
            if(data.NowPosition == my_position){
                alert("Your turn")
            }
            //渲染回合内容
            break;
        }

        $('#chatbox li').first().before(li);
    };

    var postMsg = function (message,type){
        let msg = {
            "Message" : message,
            "Type"    : type
        }
        let send_json = JSON.stringify(msg);
        socket.send(send_json);
    }

    $('#add_point_button').click(function () {
        var add_point = $('#add_point').val();
        postMsg(add_point,'add_point');
    });

    $('#start_game').click(function () {
        postMsg('start','game_op');
    });


    $('#show_card3').click(function () {
        postMsg('show_card3','game_op');
    });


    $('#show_card4').click(function () {
        postMsg('show_card4','game_op');
    });


    $('#show_card5').click(function () {
        postMsg('show_card5','game_op');
    });

    $('.quantity').each(function() {
        var spinner = $(this),
          input = spinner.find('input[type="number"]'),
          btnUp = spinner.find('.quantity-up'),
          btnDown = spinner.find('.quantity-down'),
          min = input.attr('min'),
          max = input.attr('max');
          step = input.attr('step');

          step = parseInt(step)
      
        btnUp.click(function() {
          var oldValue = parseFloat(input.val());
          if (oldValue >= max) {
            var newVal = oldValue;
          } else {
            var newVal = oldValue + step;
          }
          spinner.find("input").val(newVal);
          spinner.find("input").trigger("change");
        });
      
        btnDown.click(function() {
          var oldValue = parseFloat(input.val());
          if (oldValue <= min) {
            var newVal = oldValue;
          } else {
            var newVal = oldValue - step;
          }
          spinner.find("input").val(newVal);
          spinner.find("input").trigger("change");
        });
      
      });


});
