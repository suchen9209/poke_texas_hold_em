$(document).ready(function () {
    $('#login_button').click(function () {
        $.ajax({
            url:"/v2/login",
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
});
