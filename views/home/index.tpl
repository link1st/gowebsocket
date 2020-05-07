<!DOCTYPE html>
<html>

<head>
    <title>{{ .title }}--房间Id({{ .appId }})</title>
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
        height: 600px;
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

    /*room list*/
    .room-list{
        width: 600px;
        height: 200px;
        margin-left: 20px;
    }
    .room-list a{
        color: #428bca;
        text-decoration: none;
        padding-right: 20px;
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
        <div class="room-list">
            <div>
                <b>房间列表:</b><br>
               <a href="/home/index?appId=101">聊天室-Id:101</a>
               <a href="/home/index?appId=102">聊天室-Id:102</a>
               <a href="/home/index?appId=103">聊天室-Id:103</a>
               <a href="/home/index?appId=104">聊天室-Id:104</a>
           </div>
        </div>

        <script src="http://91vh.com/js/jquery-2.1.4.min.js"></script>
        <script type="text/javascript">
        appId = {{ .appId }};

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
        ws = new WebSocket("ws://{{ .webSocketUrl }}/acc");

        ws.onopen = function(evt) {
            console.log("Connection open ...");

            // // 连接以后
            // person = prompt("请输入你的名字", "hello-" + currentTime());
            // if (person != null) {
            //     console.log("用户准备登陆:" + person);
            //     ws.send('{"seq":"' + sendId() + '","cmd":"login","data":{"userId":"' + person + '","appId":'+ appId +'}}');
            // }

           person =  getName();
           // person = randomNumber(10000, 99999)
            console.log("用户准备登陆:" + person);
            ws.send('{"seq":"' + sendId() + '","cmd":"login","data":{"userId":"' + person + '","appId":'+ appId +'}}');

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
                addChatWith(msg("管理员", "欢迎 " + data.from + " 加入~"))
                addUserList(data.from)
            } else if (data_array.cmd === "exit") {
                data = data_array.response.data
                addChatWith(msg("管理员", data.from + " 悄悄的离开了~"))
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
                    url: 'http://{{ .httpUrl }}/user/sendMessageAll',
                    data: {
                        appId: appId,
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
                url: "http://{{ .httpUrl }}/user/list?appId=" + appId,
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

        function getName(){
            let names = ["司马相如","扬雄","班固","张衡","李白","杜甫","白居易","元稹","苏轼","辛弃疾","柳永","李清照","关汉卿","马致远","白朴","郑光祖","罗贯中","施耐庵","吴承恩","曹雪芹","杨朔","魏巍","秦牧","刘伯羽","王勃","杨炯","卢照邻","骆宾王","欧阳洵","褚遂良","虞世南","薛稷","杜审言","崔融","李峤","苏味道","唐代颜真卿","柳公权","欧阳洵","元之赵孟頫","蔡襄","黄庭坚","米芾","苏东坡","王安石","欧阳修","苏东坡","黄庭坚","黄庭坚","张耒","晃无咎","秦观","谢良佐","游酢","杨时","吕大临","李唐","刘松年","马远","夏圭？","杨万里","陆游","范成大","尤袤","关汉卿","马致远","郑光祖","白朴","黄晋","虞集","柳贯","揭俊斯","黄公望","吴镇","倪瓒","王蒙","虞集","杨载","范椁","揭俟斯","高启","张羽","徐贲","杨基","唐伯虎","祝枝山","文征明","周文宾","祝枝山","唐伯虎","文征明","徐祯卿","王时敏","王","王鉴","王原祁","方以智","陈贞慧","冒襄","侯方域","齐国孟尝君","赵国平原君","楚国春申君","魏国信陵君","陆游","杨万里","范成大","尤袤","伏羲","神农","黄帝","黄帝","少昊","颛顼","喾","尧","东伯侯姜桓楚","南伯侯鄂崇禹","西伯侯姬昌","北伯侯崇侯虎","颜回","闵损","冉耕","冉雍","冉求","仲由","宰予","端沐赐","言偃","卜商","颛孙师","曾参","澹台灭明","宓不齐","原宪","公冶长","南宫括","公皙哀","曾蒧","颜无繇","商瞿","高柴","漆雕开","公伯缭","司马耕","樊须","公西赤","巫马施","梁鳣","颜幸","冉孺","曹恤","伯虔","公孙龙","冉季","公祖句兹","秦祖","漆雕哆","颜高","漆雕徒父","壤驷赤","商泽","石作蜀","任不齐","公良孺","后处","秦冉","公夏首","奚容箴","公肩定","颜祖","鄡单","罕父黑","秦商","申党","颜之仆","荣旗","县成","左人郢","燕伋","郑国","秦非","施之常","颜哙","步叔乘","原亢籍","乐欬","廉絜","叔仲会","颜何","狄黑","邦巽","孔忠","公西舆如","公西箴","齐桓公","宋襄公","晋文公","秦穆公","楚庄王","孔子","老子","墨子","张远山公孙接","田开疆","貂蝉","西施","王昭君","杨贵妃","勾践","范蠡","文种","史圣左丘明","商圣范蠡","武圣孙膑","孙膑","庞涓","平原君赵胜","孟尝君田文","信陵君魏无忌","春申君黄歇","白起","王翦","廉颇","李牧","荆轲","专诸","聂政","要离","苏秦","张仪","廉颇","蔺相如","王翦","蒙恬","韩信","张良","萧何","九江王英布","韩王韩信","大梁王彭越","李广","李敢","李陵","贾谊","晁错","司马相如","司马迁","司马相如","杨雄","班固","张衡","吴王刘濞","楚王刘戊","赵王刘遂","胶西王刘印","济南王刘辟光","菑川王刘贤","胶东王刘雄渠","卫清","霍去病","许虔","许劭","邓禹","吴汉","贾复","耿弇","寇恂","岑彭","冯异","朱祜","祭遵","景丹","盖延","铫期","耿纯","臧宫","马武","刘隆为孔融","陈琳","王粲","徐干","阮禹","应瑒","刘桢","颜良","文丑","张颌","高览","淳于琼","华歆","邴原","管宁","张昭","张纮","孙乾","简庸","糜竺","曹豹","诸葛亮","诸葛瑾","诸葛诞","曹操","曹丕","曹植","刘备","关羽","张飞","诸葛亮","关羽","张飞","诸葛亮","蒋琬","董允","费袆","关羽","张飞","赵云","马超","黄忠","阮籍","嵇康","山涛","刘伶","阮咸","向秀","王戎","陈翔","范滂","孔昱","范康","檀敷","张俭","刘表","岑咥","司马朗","司马懿","司马孚","司马旭","司马恂","司马进","司马通","司马敏","蹇硕","曹操","袁绍","鲍鸿","赵融","冯芳","夏牟","淳于琼","张让","赵忠","夏恽","郭胜","孙璋","毕岚","段摇","高望","张恭","韩悝","宋典","粟嵩","冯翎","山子道","王九真","郭凯","王恺","石崇","汝南王亮","楚王玮","赵王伦","齐王冏","河间王颙","成都王颖","长沙王乂","东海王越","王导","谢玄","陆机","陆云","谢灵运","谢惠连","谢眺","斛律光","兰陵王","顾恺之","陆探微","张僧繇","史万岁","韩擒虎","贺若弼","杨素","房玄龄","杜如悔","杜如晦","房玄龄","于志宁","苏世长","薛收","褚亮","姚思廉","陆德时","孔颖达","李玄道","李守素","虞世南","蔡允恭","颜相时","许敬宗","薛元敬","盖文达","苏勖","王勃","杨炯","卢照邻","骆宾王","贺知章","张旭","包融","张若虚","安录山","史思明","薛元敬","薛收","薛德音","李白","李贺","李商隐","鸠摩罗什","真谛","玄奘","延平","延定","延朗","延辉","延德","延昭","延嗣","延顺","程颢","程颐","北宋善画的李伯时","能文的李亮工","工书的李元中","苏洵","苏轼","苏辙","韩愈","柳宗元","欧阳修","苏洵","苏轼","苏辙","曾巩","王安石","李唐","刘松年","马远","夏圭","韩世忠","岳飞","张浚","刘琦","窝阔台","哲别","者勒蔑","速不台","博尔忽","博尔术","木华黎","赤老温","关汉卿","白朴","马致远","郑光祖","唐伯虎","祝枝山","文征明","张梦晋","杨士奇","杨荣","杨溥","袁宏道","袁中道","袁宗道","杨涟","左光斗","魏大中","周朝瑞","袁化中","顾大章","顾炎武","黄宗曦","王夫之","戚继光","袁崇焕","郑成功","田园诗人陶渊明","江州剌史李渤","江州司马白居易","理学大师周敦颐","王阳明","礼亲王","郑亲王","睿亲王","豫亲王","肃亲王","庄亲王","克勤郡王","顺承郡王","三藩","平西王吴三挂","平南王尚之信","靖南王耿精忠","肃顺","载垣","端华","焦佑瀛","杜翰","景寿","穆荫","匡源","林旭","杨锐","谭嗣同","康广仁","刘光第","杨深秀"];


            var name = names[Math.floor(Math.random()*names.length)];

            return name
        }
        </script>
    </div>
</body>

</html>