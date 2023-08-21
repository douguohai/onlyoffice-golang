<!DOCTYPE html>
<html style="height: 100%;">
	<head>
	    <title>Office</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
        <meta name="renderer" content="webkit">
        <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
	</head>

	<body style="height: 100%; margin: 0;">
		<div id="placeholder" style="height: 100%"></div>
    <script type="text/javascript" src="{{.documentServer}}/web-apps/apps/api/documents/api.js"></script>
    <script src="https://cdn.bootcss.com/sockjs-client/1.1.4/sockjs.min.js"></script>
    <script src="https://cdn.bootcdn.net/ajax/libs/jquery/3.6.4/jquery.js"></script>
    <script src="https://cdn.bootcss.com/stomp.js/2.3.3/stomp.min.js"></script>

    <script type="text/javascript">

       window.docId="{{.Doc.Id}}";

    	window.docEditor = new DocsAPI.DocEditor("placeholder",{
            "events": {

            },

            "document": {
              "fileType": "{{.fileType}}",
              "key": "{{.Key}}",
              "title": "{{.Doc.OriginName}}",
              "url": "{{.serverUrl}}/static/{{.Doc.FileName}}",
              "info": {
              },
              "permissions": {
                "comment": false,
                "download": true,
                "edit": {{.Edit}},
                "print": true,
                "review": {{.Review}},//true
                "chat":false,
                },
            },
        "documentType": "{{.documentType}}",
        "editorConfig": {
          "callbackUrl": "{{.serverUrl}}/url-to-callback?id={{.Doc.Id}}",
          "createUrl": "https://example.com/url-to-create-document/",
          "user": {
             "id": {{.Uid}},
             "name": "{{.Username}}"
            },
			"customization": {
                "commentAuthorOnly": false,
                "compactToolbar": false,
                "feedback": {
                  "visible": false
                },
                "forcesave": true,
                "goback": {
                  "text": "Go to Documents",
                  "url": "http://www.baidu.com"
                },
                "zoom": 100,
        	},
        	"embedded": {
                "embedUrl": "https://example.com/embedded?doc=exampledocument1.docx",
                "fullscreenUrl": "https://example.com/embedded?doc=exampledocument1.docx#fullscreen",
                "saveUrl": "https://example.com/download?doc=exampledocument1.docx",
                "shareUrl": "https://example.com/view?doc=exampledocument1.docx",
                "toolbarDocked": "top"
        	},

           "lang": "zh-CN",//"en-US",
           "mode": {{.Mode}},//"view",//edit
            "recent": [
        	]
        },

        "height": "100%",
    	"type": {{.Type}},//"desktop",//desktop//embedded//mobile访问文档的平台类型 网页嵌入
        "width": "100%"
      });



      function connectWebSocket() {
          var socket = new WebSocket("{{.wsServer}}?id="+docId);

          socket.onopen = function () {
              console.log("WebSocket连接已建立");
          }

          socket.onmessage = function (event) {
              console.log("收到消息:", event.data);
              msg = JSON.parse(event.data)
              console.log(window.docId)
              if (msg.type === 0 && msg.data === parseInt(window.docId)) {
                  window.location.reload()
              }
          }

          socket.onclose = function (event) {
              console.log("WebSocket连接已关闭，将在3秒后重新连接");
              setTimeout(connectWebSocket, 3000);
          }

          socket.onerror = function (error) {
              console.error("WebSocket连接出错:", error);
              socket.close();
          }
      }

      $(function () {
          connectWebSocket();
      });






   	</script>
	</body>
</html>
