<!DOCTYPE html>
<html>

<head>
    <title>{{ .title }}</title>
    <style type="text/css">
    /*公共样式*/
    body,
    h1,
    h2,
    h3,
    h4,
    p,
    ul,
    ol,
    li,
    form,
    button,
    input,
    textarea,
    th,
    td {
        margin: 0;
        padding: 0
    }

    body,
    button,
    input,
    select,
    textarea {
        font: 12px/1.5 Microsoft YaHei UI Light, tahoma, arial, "\5b8b\4f53";
        *line-height: 1.5;
        -ms-overflow-style: scrollbar
    }

    h1,
    h2,
    h3,
    h4 {
        font-size: 100%
    }

    ul,
    ol {
        list-style: none
    }

    a {
        text-decoration: none
    }

    a:hover {
        text-decoration: underline
    }

    img {
        border: 0
    }

    button,
    input,
    select,
    textarea {
        font-size: 100%
    }

    table {
        border-collapse: collapse;
        border-spacing: 0
    }

    /*rem*/
    html {
        font-size: 62.5%;
    }

    body {
        font: 16px/1.5 "microsoft yahei", 'tahoma';
    }

    body .mobile-page {
        font-size: 1.6rem;
    }

    /*浮动*/
    .fl {
        float: left;
    }

    .fr {
        float: right;
    }

    .clearfix:after {
        content: '';
        display: block;
        height: 0;
        clear: both;
        visibility: hidden;
    }

    body {
        background-color: #F5F5F5;
    }

    .mobile-page {
        max-width: 600px;
    }

    .mobile-page .admin-img,
    .mobile-page .user-img {
        width: 45px;
        height: 45px;
    }

    i.triangle-admin,
    i.triangle-user {
        width: 0;
        height: 0;
        position: absolute;
        top: 10px;
        display: inline-block;
        border-top: 10px solid transparent;
        border-bottom: 10px solid transparent;
    }

    .mobile-page i.triangle-admin {
        left: 4px;
        border-right: 12px solid #fff;
    }

    .mobile-page i.triangle-user {
        right: 4px;
        border-left: 12px solid #9EEA6A;
    }

    .mobile-page .admin-group,
    .mobile-page .user-group {
        padding: 6px;
        display: flex;
        display: -webkit-flex;
    }

    .mobile-page .admin-group {
        justify-content: flex-start;
        -webkit-justify-content: flex-start;
    }

    .mobile-page .user-group {
        justify-content: flex-end;
        -webkit-justify-content: flex-end;
    }

    .mobile-page .admin-reply,
    .mobile-page .user-reply {
        display: inline-block;
        padding: 8px;
        border-radius: 4px;
        background-color: #fff;
        margin: 0 15px 12px;
    }

    .mobile-page .admin-reply {
        box-shadow: 0px 0px 2px #ddd;
    }

    .mobile-page .user-reply {
        text-align: left;
        background-color: #9EEA6A;
        box-shadow: 0px 0px 2px #bbb;
    }

    .mobile-page .user-msg,
    .mobile-page .admin-msg {
        width: 75%;
        position: relative;
    }

    .mobile-page .user-msg {
        text-align: right;
    }

    /*界面*/
    .interface {
        width: 1000px;
        height: 700px;
    }

    .personnel-list {
        float: left;
        width: 200px;
        height: 500px;
        background-color: #bbbbbb;
        border-style: solid;
        border-color: #000000;
        overflow: scroll;
    }

    /*聊天框*/
    .chat-with {
        float: left;
        width: 600px;
        height: 400px;
        background-color: #bbbbbb;
        border-style: solid;
        border-color: #000000;
        overflow: scroll;
    }

    .send-msg {
        float: left;
        width: 600px;
        height: 100px;
        background-color: #bbbbbb;
        border-style: solid;
        border-color: #000000;
        overflow: scroll;
    }
    </style>
</head>

<body>
    <div class="mobile-page">
        <div class="interface">
            <div class="personnel-list">
                <ul class="personnel-list-ul">
                </ul>
                <!-- 在线列表 -->
                <!-- 用户列表
            进入的时候拉取用户列表
            有人加入的时候添加
            有人退出以后删除 -->
            </div>
            <div class="chat-with">
                <div class="admin-group">
                    <div class="admin-img">
                        管理员
                    </div>
                    <div class="admin-msg">
                        <i class="triangle-admin"></i>
                        <span class="admin-reply">欢迎加入聊天~</span>
                    </div>
                </div>
            </div>
            <div class="send-msg">
                <!-- <input type="text" name="msg" placeholder="你想要发送的消息" value="" size="35">
                <button type="submit"> send</button> -->
                <form onsubmit="return doSubmit();">
                    <input type="text" name="msg" placeholder="你想要发送的消息" value="" size="35"/>
                    <input type="button" name="button" value="send" />
                </form>
            </div>
        </div>
        <script src="http://91vh.com/js/jquery-2.1.4.min.js"></script>
        <script type="text/javascript">
        function currentTime() {
            let timeStamp = (new Date()).valueOf();

            return timeStamp
        }

        function randomNumber(minNum, maxNum) {
            switch (arguments.length) {
                case 1:
                    return parseInt(Math.random() * minNum + 1, 10);
                    break;
                case 2:
                    return parseInt(Math.random() * (maxNum - minNum + 1) + minNum, 10);
                    break;
                default:
                    return 0;
                    break;
            }
        }


        function sendId() {

            let timeStamp = currentTime();
            let randId = randomNumber(100000, 999999);
            let id = timeStamp + "-" + randId;

            return id
        }

        function msg(name, msg) {
            let html = '<div class="admin-group">' +
                '<div class="admin-img" >' + name + '</div>' +
                // '<img class="admin-img" src="http://localhost/public/img/aa.jpg" />'+
                '<div class="admin-msg">' +
                '<i class="triangle-admin"></i>' +
                '<span class="admin-reply">' + msg + '</span>' +
                '</div>' +
                '</div>';
            return html
        }

        function myMsg(name, msg) {
            let html = '<div class="user-group">' +
                '<div class="user-msg">' +
                '<span class="user-reply">' + msg + '</span>' +
                '<i class="triangle-user"></i>' +
                '</div>' +
                '<div class="user-img" >' + name + '</div>' +
                // '<img class="user-img" src="http://localhost/public/img/cc.jpg" />'+
                '</div>';
            return html
        }

        function userDiv(name) {

            let html = '<div id="' + name + '">' +
                name +
                '</div>';
            return html

        }

        function addChatWith(msg) {
            $(".chat-with").append(msg);
            // 页面滚动条置底
            $('.chat-with').animate({ scrollTop: document.body.clientHeight + 10000 + 'px' }, 80);

        }

        function addUserList(name) {
            music = "<li id=\"" + name + "\">" + name + "</li>";
            $(".personnel-list-ul").append(music);
        }

        function delUserList(name) {
            $("#" + name).remove();
        }


        // 连接webSocket
        ws = new WebSocket("ws://127.0.0.1:8089/acc");

        ws.onopen = function(evt) {
            console.log("Connection open ...");

            // // 连接以后
            // person = prompt("请输入你的名字", "hello-" + currentTime());
            // if (person != null) {
            //     console.log("用户准备登陆:" + person);
            //     ws.send('{"seq":"' + sendId() + '","cmd":"login","data":{"userId":"' + person + '","appId":101}}');
            // }

            person = randomNumber(10000, 99999)
            console.log("用户准备登陆:" + person);
            ws.send('{"seq":"' + sendId() + '","cmd":"login","data":{"userId":"' + person + '","appId":101}}');

            // 定时心跳
            setInterval(heartbeat, 30 * 1000)
        };

        // 收到消息
        ws.onmessage = function(evt) {
            console.log("Received Message: " + evt.data);
            data_array = JSON.parse(evt.data);
            console.log(data_array);

            if (data_array.cmd === "msg") {
                data = data_array.response.data
                addChatWith(msg(data.from, data.msg))
            } else if (data_array.cmd === "enter") {
                data = data_array.response.data
                addChatWith(msg("管理员", "欢迎" + data.from + "加入~"))
                addUserList(data.from)
            } else if (data_array.cmd === "exit") {
                data = data_array.response.data
                addChatWith(msg("管理员", data.from + "悄悄的离开了~"))
                delUserList(data.from)
            }


        };

        ws.onclose = function(evt) {
            console.log("Connection closed.");
        };

        // 心跳
        function heartbeat() {
            console.log("定时心跳:" + person);
            ws.send('{"seq":"' + sendId() + '","cmd":"heartbeat","data":{}}');

        }

        // 点击按钮事件
        // $("button").click(function() {
        $("input[name='button']").click(function() {
            sendMsg();
        });

// 回车提交
        function doSubmit() {
            sendMsg();
            return false;
        }

        function sendMsg() {
            let msg = $("input[name='msg']").val()
            console.log("button 点击:" + msg);
            if (msg !== "") {

                $.ajax({
                    type: "POST",
                    url: 'http://127.0.0.1:8080/user/sendMessageAll',
                    data: {
                        appId: "101",
                        userId: person,
                        msgId: sendId(),
                        message: msg,
                    },
                    contentType: "application/x-www-form-urlencoded",
                    success: function(data) {
                        console.log(data);
                        addChatWith(myMsg(person, msg))
                        $("input[name='msg']").val("");
                    }
                });
            }
        }

        setTimeout(function() { getUserList(); }, 500); // 1秒后将会调用执行

        function getUserList() {
            $.ajax({
                type: "GET",
                url: "http://127.0.0.1:8080/user/list",
                dataType: "json",
                success: function(data) {
                    console.log("user list:" + data.code + "userList:" + data.data.userList);
                    if (data.code != 200) {
                        return false
                    }
                    var music = "";
                    //i表示在data中的索引位置，n表示包含的信息的对象
                    $.each(data.data.userList, function(i, n) {
                        //获取对象中属性为optionsValue的值
                        let name = n
                        if (n == person) {
                            name = name + "(自己)"
                        }
                        music += "<li id=\"" + n + "\">" + name + "</li>";
                    });
                    $(".personnel-list-ul").append(music);

                    return false
                }
            });
        }
        </script>
    </div>
</body>

</html>