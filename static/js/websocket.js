var socket;
var my_position;
var uname;
var user_id;

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
        if(data.Type != 3){
            console.log(data);
        }
        

        switch (data.Type) {
        case 0: // JOIN
            if (data.User == uname) {
                //You joined the chat room.
                my_position = data.GameUser.Position
                user_id = data.GameUser.UserId
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
            if(data.NowPosition == my_position){
                $("#your_turn").show();
                $(".quantity").show(); 
                $("#add_point").attr("min",data.MaxPoint)
                $("#add_point").attr("max",data.Detail.Point)
                if(data.MaxPoint <= data.Detail.RoundPoint){
                    $("#check").show(); 
                }else{
                    $("#check").hide();
                }
            }else{
                $(".quantity").hide();
            }
            var roundhtml = "<p>"+data.GM.GameStatus+"</p>"
            + "<p>当前轮底池"+ data.AllPointInRound + "</p>"
            + "<p>当前位置"+ data.NowPosition + "</p>"
            + "<p>最小Point" + data.MaxPoint + "</p>"
            + "<p>小盲"+ data.GM.BigBindPosition + "</p>"
            + "<p>第一轮底池："+ data.GM.Pot1st + "</p>"
            + "<p>第二轮底池："+ data.GM.Pot2nd + "</p>"
            + "<p>第三轮底池："+ data.GM.Pot3rd + "</p>"
            + "<p>第四轮底池："+ data.GM.Pot4th + "</p>"
            $("#RoundInfo").html(roundhtml);
            //渲染回合内容
            break;
        case 8://玩家操作
            var ophtml = "<p>"+ data.Name + " " + data.GameMatchLog.Operation + data.GameMatchLog.PointNumber + "</p>";
            $("#UserOp").append(ophtml);
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

    var postOperation = function (type,operation,point){
        let msg = {
            "Type"      :   type,
            "Operation" :   operation,
            "Point"     :   point,
            "UserId"    :   user_id,
            "Position"  :   my_position,
            "Name"      :   uname
        }
        let send_json = JSON.stringify(msg);
        socket.send(send_json);
        $("#your_turn").hide();
    }


    $('#add_point_button').click(function () {
        var add_point = $('#add_point').val();
        postMsg(add_point,'add_point');
    });

    $('#start_game').click(function () {
        postMsg('start','game_op');
        $("#UserOp").append("");
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


    $('#call').click(function () {
        //call
        postOperation('user_op','call',0);
    });

    $('#raise').click(function () {
        //raise
        var add_point = $('#add_point').val();
        console.log(add_point)
        postOperation('user_op','raise',parseInt(add_point));
    });

    
    $('#check').click(function () {
        //check
        postOperation('user_op','check',0);
    });

    
    $('#fold').click(function () {
        //check
        postOperation('user_op','fold',0);
    });
    
    $('#allin').click(function () {
        //check
        postOperation('user_op','allin',0);
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
