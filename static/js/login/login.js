$(document).ready(function () {
    $('#login_button').click(function () {
        $.ajax({
            url:"/v2/user/login",
            type:"POST",
            data:{
                uname : $("#login_name").val(),
                password:$("#login_pwd").val()
            },
            success:function (data){
                console.log(data)
            }

        });
    });

    $('#register_button').click(function () {
        $.ajax({
            url:"/v2/user/register",
            type:"POST",
            data:{
                uname : $("#register_name").val(),
                password:$("#register_pwd").val(),
                password2:$("#register_pwd2").val()
            },
            success:function (data){
                console.log(data)
            }

        });
    });

    $("#show_login").on('click',()=>{
        $("#register_panel").hide();
        $("#login_panel").show();
    });
    $("#show_register").on('click',()=>{
        $("#login_panel").hide();
        $("#register_panel").show();
    });
});
