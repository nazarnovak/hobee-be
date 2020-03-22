(window.webpackJsonp=window.webpackJsonp||[]).push([[0],[,,,,,,,,,,,,,function(e,t,a){"use strict";a.d(t,"a",function(){return d});var n=a(2),s=a(3),r=a(4),o=a(6),c=a(5),i=a(7),l=a(0),u=a.n(l),p=a(67),h=a(46),d=function(e){function t(e){var a;return Object(s.a)(this,t),(a=Object(o.a)(this,Object(c.a)(t).call(this,e))).state={location:window.location},a}return Object(i.a)(t,e),Object(r.a)(t,[{key:"render",value:function(){var e=!0;return"/chat"===this.state.location.pathname&&(e=!1),u.a.createElement("div",{class:"nav-height"},u.a.createElement("div",{class:"nav nav-height"},u.a.createElement(p.a,{class:"logo-link logo-height",to:"/"},u.a.createElement("img",{class:"logo logo-height",src:h,alt:"Logo"})),u.a.createElement(m,{visible:e})))}}]),t}(u.a.Component),m=(u.a.Component,function(e){function t(e){var a;return Object(s.a)(this,t),(a=Object(o.a)(this,Object(c.a)(t).call(this,e))).state={hover:!1},a.onMouseEnter=a.onMouseEnter.bind(Object(n.a)(Object(n.a)(a))),a.onMouseLeave=a.onMouseLeave.bind(Object(n.a)(Object(n.a)(a))),a}return Object(i.a)(t,e),Object(r.a)(t,[{key:"onMouseEnter",value:function(){this.setState({hover:!0})}},{key:"onMouseLeave",value:function(){this.setState({hover:!1})}},{key:"render",value:function(){return u.a.createElement("div",{className:"chat-button-link-wrapper"},u.a.createElement(p.a,{to:"/chat"},u.a.createElement("button",{className:"chat-button-link scale".concat(this.props.visible?"":" fade")},"Chat")))}}]),t}(u.a.Component))},,,,,,,,,,,,,,function(e,t,a){"use strict";(function(e){a.d(t,"a",function(){return U});var n=a(8),s=a.n(n),r=a(10),o=a(3),c=a(4),i=a(6),l=a(5),u=a(7),p=a(2),h=a(0),d=a.n(h),m=a(13),b="identifying",v="searching",f="matched",g="a",k="c",y="r",O="s",E="t",w="rl",j="rd",x={rdl:"I didn't like the conversation",rsp:"Spam",rse:"Sexism",rha:"Harassment",rra:"Racism",rot:"Other"},N="s",S="d",C=a(53),M=a(54),T=a(55),R=a(56),L=(a(57),a(58),a(59)),W=a(60),B=a(61),D=a(62),U=function(t){function a(e){var t;Object(o.a)(this,a),(t=Object(i.a)(this,Object(l.a)(a).call(this,e))).onFocus=function(){t.scrollToBottom(),t.setState({tabActive:!0,unread:0}),document.title="hobee: Quality conversations"},t.onBlur=function(){t.setState({tabActive:!1})},t.handleEscKey=function(e){if("Escape"===e.key)return t.handleDisconnect(),!0},t.onMessage=function(){return function(e){var a=JSON.parse(e.data);switch(null!==a&&void 0!==a||console.log("Received unexpected JSON:",a),a.type){case O:t.handleSystemMessage(a);break;case k:t.handleChatMessage(a);break;case g:t.handleActivityMessage(a);break;default:console.log("Unexpected json:",a)}}},t.onResize=function(){t.scrollToBottom()},t.handleDisconnect=function(){if(t.state.status!==f)return!1;if(!window.confirm("Are you sure you want to disconnect?"))return!1;var e={type:O,text:S};t.state.websocket.send(JSON.stringify(e))},t.handleMessageClick=function(e){var t=e.target;return t.parentElement.querySelector(".message-timestamp").classList.contains("visible")?(t.parentElement.querySelector(".message-timestamp").classList.remove("visible"),!0):(t.parentElement.querySelector(".message-timestamp").classList.add("visible"),!0)},t.handleReportModalClose=function(e){e.target===e.currentTarget&&t.setState({reportModalOpen:!1})},t.handleReportOptionClick=function(e){var a=e.target.getAttribute("data-key"),n=!1;if(Object.keys(x).map(function(e){e===a&&(n=!0)}),!n)return!1;t.sendWebsocketMessage(y,a),t.setState({reportModalOpen:!1,reported:a})},t.handleSearch=function(){t.clearChat();var e={type:O,text:N};t.state.websocket.send(JSON.stringify(e)),t.setState({status:v,statusShow:!1,reported:"",liked:!1})},t.handleLike=function(){t.state.liked&&t.sendWebsocketMessage(y,j),t.state.liked||t.sendWebsocketMessage(y,w),t.setState({liked:!t.state.liked})},t.handleSave=function(){console.log("Save clicked")},t.handleReport=function(){t.setState({reportModalOpen:!0})};var n=Object(p.a)(Object(p.a)(t));return t.state={messages:[],status:b,websocket:null,liked:!1,reportModalOpen:!1,reported:"",tabActive:!0,unread:0,statusShow:!1,statusText:"",typingTimeout:setTimeout(function(){if("Buddy is typing..."!==n.state.statusText)return!1;n.setState({statusShow:!1})},2e3)},t.connectToChat=t.connectToChat.bind(Object(p.a)(Object(p.a)(t))),t.connectToWs=t.connectToWs.bind(Object(p.a)(Object(p.a)(t))),t.onOpen=t.onOpen.bind(Object(p.a)(Object(p.a)(t))),t.onClose=t.onClose.bind(Object(p.a)(Object(p.a)(t))),t.sendWebsocketMessage=t.sendWebsocketMessage.bind(Object(p.a)(Object(p.a)(t))),t.handleSystemMessage=t.handleSystemMessage.bind(Object(p.a)(Object(p.a)(t))),t.handleChatMessage=t.handleChatMessage.bind(Object(p.a)(Object(p.a)(t))),t.clearChat=t.clearChat.bind(Object(p.a)(Object(p.a)(t))),t}return Object(u.a)(a,t),Object(c.a)(a,[{key:"componentDidUpdate",value:function(e){e.messages!==this.props.messages&&this.setState({messages:this.props.messages})}},{key:"componentWillUnmount",value:function(){window.removeEventListener("keydown",this.handleEscKey,!1),window.removeEventListener("focus",this.onFocus),window.removeEventListener("blur",this.onBlur),window.removeEventListener("resize",this.onResize)}},{key:"clearChat",value:function(){this.setState({messages:[]})}},{key:"connectToChat",value:function(){var e=Object(r.a)(s.a.mark(function e(){return s.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,this.identify();case 2:if(e.sent){e.next=5;break}return e.abrupt("return",!1);case 5:if(this.setState({status:"connecting"}),this.connectToWs()){e.next=9;break}return e.abrupt("return");case 9:return e.abrupt("return",!0);case 10:case"end":return e.stop()}},e,this)}));return function(){return e.apply(this,arguments)}}()},{key:"connectToWs",value:function(){var t=Object(r.a)(s.a.mark(function t(){var a,n,r;return s.a.wrap(function(t){for(;;)switch(t.prev=t.next){case 0:return a=null,8080,e&&Object({NODE_ENV:"production",PUBLIC_URL:""})&&Object({NODE_ENV:"production",PUBLIC_URL:""}).PORT&&Object({NODE_ENV:"production",PUBLIC_URL:""}).PORT,t.prev=3,n=window.location,r="https:"===n.protocol?"wss:":"ws:",r+="//"+n.host,r+="/api/chat",t.next=10,new WebSocket(r);case 10:if(void 0!==(a=t.sent)){t.next=13;break}throw new Error("Could not connect to ws");case 13:t.next=19;break;case 15:return t.prev=15,t.t0=t.catch(3),console.log(t.t0),t.abrupt("return",!1);case 19:return a.binaryType="arraybuffer",a.onopen=this.onOpen,a.onclose=this.onClose,a.onmessage=this.onMessage(),this.setState({websocket:a}),t.abrupt("return",!0);case 25:case"end":return t.stop()}},t,this,[[3,15]])}));return function(){return t.apply(this,arguments)}}()},{key:"componentDidMount",value:function(){var e=Object(r.a)(s.a.mark(function e(){return s.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:window.addEventListener("keydown",this.handleEscKey,!1),window.addEventListener("focus",this.onFocus),window.addEventListener("blur",this.onBlur),window.addEventListener("resize",this.onResize),this.connectToChat();case 5:case"end":return e.stop()}},e,this)}));return function(){return e.apply(this,arguments)}}()},{key:"pullRoomMessages",value:function(){var e=Object(r.a)(s.a.mark(function e(){var t,a;return s.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return"/api/messages",e.prev=1,e.next=4,fetch("/api/messages",{credentials:"include"});case 4:return a=e.sent,e.next=7,a.json();case 7:t=e.sent,e.next=14;break;case 10:return e.prev=10,e.t0=e.catch(1),console.log(e.t0),e.abrupt("return",!1);case 14:if(void 0!==t.messages){e.next=17;break}return console.log("Unknown response:",t),e.abrupt("return",!1);case 17:return e.abrupt("return",t.messages);case 18:case"end":return e.stop()}},e,this,[[1,10]])}));return function(){return e.apply(this,arguments)}}()},{key:"pullResult",value:function(){var e=Object(r.a)(s.a.mark(function e(){var t,a;return s.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return"/api/result",e.prev=1,e.next=4,fetch("/api/result",{credentials:"include"});case 4:return a=e.sent,e.next=7,a.json();case 7:t=e.sent,e.next=14;break;case 10:return e.prev=10,e.t0=e.catch(1),console.log(e.t0),e.abrupt("return",!1);case 14:if(void 0!==t.liked){e.next=17;break}throw new Error("Unknown pull result response",t);case 17:return e.abrupt("return",t);case 18:case"end":return e.stop()}},e,this,[[1,10]])}));return function(){return e.apply(this,arguments)}}()},{key:"chatStatusFromMessages",value:function(e){var t=this.state.statusShow,a=this.state.statusText;e.map(function(e){e.type===g&&"b"===e.authoruuid&&"ui"===e.text&&(t=!0,a="Buddy is inactive"),e.type===g&&"b"===e.authoruuid&&"ua"===e.text&&(t=!1),e.type===O&&e.text===S&&(t=!0,a="You disconnected","b"===e.authoruuid&&(a="Buddy disconnected"))}),this.setState({statusShow:t,statusText:a})}},{key:"handleSystemMessage",value:function(){var e=Object(r.a)(s.a.mark(function e(t){var a,n,r;return s.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:n=this.state.messages.slice(),e.t0=t.text,e.next="c"===e.t0?4:e.t0===S?7:e.t0===N?13:16;break;case 4:return console.log("Matched"),a=f,e.abrupt("break",17);case 7:return console.log("Disconnected"),r="You disconnected","b"===t.authoruuid&&(r="Buddy disconnected"),this.setState({statusShow:!0,statusText:r}),a="disconnected",e.abrupt("break",17);case 13:return console.log("Available for search"),a=v,e.abrupt("break",17);case 16:throw new Error("Unknown system message",t.text);case 17:a!=v&&n.push(t),a===v&&(this.clearChat(),this.sendWebsocketMessage(O,N)),this.setState({status:a,messages:n});case 20:case"end":return e.stop()}},e,this)}));return function(t){return e.apply(this,arguments)}}()},{key:"handleChatMessage",value:function(e){var t=this.state.messages.slice();t.push(e),this.setState({messages:t}),this.scrollToBottom(),this.pushUnreadMessage()}},{key:"handleActivityMessage",value:function(){var e=Object(r.a)(s.a.mark(function e(t){var a,n,r,o;return s.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:a=this.state.status,n=this.state.messages.slice(),e.t0=t.text,e.next="ra"===e.t0?4:"ri"===e.t0?10:"ua"===e.t0?20:"ui"===e.t0?22:e.t0===E?24:29;break;case 4:return a=f,e.next=7,this.pullRoomMessages();case 7:return n=e.sent,this.chatStatusFromMessages(n),e.abrupt("break",30);case 10:return a="disconnected",e.next=13,this.pullRoomMessages();case 13:return n=e.sent,this.chatStatusFromMessages(n),e.next=17,this.pullResult();case 17:return r=e.sent,this.setState({liked:r.liked,reported:r.reported}),e.abrupt("break",30);case 20:return"b"===t.authoruuid&&this.setState({statusShow:!1}),e.abrupt("break",30);case 22:return"b"===t.authoruuid&&this.setState({statusShow:!0,statusText:"Buddy is inactive"}),e.abrupt("break",30);case 24:return this.setState({statusShow:!0,statusText:"Buddy is typing..."}),clearTimeout(this.state.typingTimeout),o=this,this.setState({typingTimeout:setTimeout(function(){if("Buddy is typing..."!==o.state.statusText)return!1;o.setState({statusShow:!1})},2e3)}),e.abrupt("return");case 29:throw new Error("Unknown activity received",t);case 30:n.push(t),this.setState({status:a,messages:n}),this.scrollToBottom();case 33:case"end":return e.stop()}},e,this)}));return function(t){return e.apply(this,arguments)}}()},{key:"pushUnreadMessage",value:function(){if(this.state.tabActive)return!1;this.setState({unread:this.state.unread+1}),document.title="("+this.state.unread+") hobee: Quality conversations"}},{key:"scrollToBottom",value:function(){document.querySelector(".chat-messages").scrollTop=document.querySelector(".chat-messages").scrollHeight}},{key:"identify",value:function(){var e=Object(r.a)(s.a.mark(function e(){var t,a,n;return s.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t="/api/identify"+window.location.search,e.prev=1,e.next=4,fetch(t,{credentials:"include"});case 4:return n=e.sent,e.next=7,n.json();case 7:a=e.sent,e.next=14;break;case 10:return e.prev=10,e.t0=e.catch(1),console.log(e.t0),e.abrupt("return",!1);case 14:if(void 0!==a.error&&!a.error){e.next=17;break}return console.log("Unknown response:",a),e.abrupt("return",!1);case 17:return e.abrupt("return",!0);case 18:case"end":return e.stop()}},e,this,[[1,10]])}));return function(){return e.apply(this,arguments)}}()},{key:"onClose",value:function(){this.setState({status:"disconnected"})}},{key:"onOpen",value:function(){}},{key:"sendWebsocketMessage",value:function(e,t){var a={type:e,text:t};if(null===this.state.websocket||void 0===this.state.websocket)return console.log("Not connected to WS"),!1;try{this.state.websocket.send(JSON.stringify(a))}catch(n){return console.log(n),!1}return!0}},{key:"render",value:function(){return d.a.createElement("div",null,d.a.createElement(m.a,{location:this.props.location}),d.a.createElement("div",{className:"main-content"},d.a.createElement(F,{messages:this.state.messages,searching:"connecting"===this.state.status||this.state.status===v,status:this.state.status,handleMessageClick:this.handleMessageClick}),d.a.createElement(I,{websocket:this.state.websocket,handleDisconnect:this.handleDisconnect,handleSearch:this.handleSearch,disconnected:"disconnected"===this.state.status,matched:this.state.status===f,handleLike:this.handleLike,liked:this.state.liked,reported:this.state.reported,handleSave:this.handleSave,handleReport:this.handleReport,searching:"connecting"===this.state.status||this.state.status===v,statusShow:this.state.statusShow,statusText:this.state.statusText,sendWebsocketMessage:this.sendWebsocketMessage,reportModalOpen:this.state.reportModalOpen,handleReportModalClose:this.handleReportModalClose,handleReportOptionClick:this.handleReportOptionClick})))}}]),a}(d.a.Component),F=function(e){function t(){return Object(o.a)(this,t),Object(i.a)(this,Object(l.a)(t).apply(this,arguments))}return Object(u.a)(t,e),Object(c.a)(t,[{key:"render",value:function(){var e=this;if(this.props.searching)return d.a.createElement("div",{className:"chat-messages"},d.a.createElement("div",{className:"status-wrapper"},d.a.createElement("div",{className:"loader"}),d.a.createElement("div",{className:"status"},this.props.status)));var t=this.props.messages.map(function(t){if(t.type!==O&&t.type!==g){if(t.type!==k)throw new Error("Unknown message type in chat messages",t.type);if("o"===t.authoruuid)return d.a.createElement("div",{className:"my-message-container"},d.a.createElement(J,{direction:"left",timestamp:t.timestamp}),d.a.createElement("div",{className:"chat-message my-message",onClick:e.props.handleMessageClick},t.text));if("b"===t.authoruuid)return d.a.createElement("div",{className:"buddy-message-container"},d.a.createElement("div",{className:"chat-message buddy-message",onClick:e.props.handleMessageClick},t.text),d.a.createElement(J,{direction:"right",timestamp:t.timestamp}));throw new Error("Unknown where to show the message",t)}});return d.a.createElement("div",{className:"chat-messages fade-in"},t)}}]),t}(d.a.Component),z=function(e){function t(){return Object(o.a)(this,t),Object(i.a)(this,Object(l.a)(t).apply(this,arguments))}return Object(u.a)(t,e),Object(c.a)(t,[{key:"render",value:function(){return d.a.createElement("div",{className:"chat-status-container"+(this.props.show?" visible":"")},d.a.createElement("div",{className:"chat-status"},this.props.text))}}]),t}(d.a.Component),I=function(e){function t(e){var a;return Object(o.a)(this,t),(a=Object(i.a)(this,Object(l.a)(t).call(this,e))).sendInputMessage=function(){var e=document.getElementsByTagName("input")[0],t=e.value;if(""===t)return!0;var n={type:"o",text:t};return a.props.websocket.send(JSON.stringify(n)),console.log("Sent message:",t),e.value="",a.setState({inputText:e.value}),!0},a.handleKeyDown=function(e){if("Enter"===e.key)return a.sendInputMessage(),!0},a.handleOnChange=function(e){var t=document.getElementsByTagName("input")[0];if(!t)return!1;if(a.setState({inputText:t.value}),!a.state.typing){a.setState({typing:!0});var n=Object(p.a)(Object(p.a)(a));a.props.sendWebsocketMessage(g,E),setTimeout(function(){n.setState({typing:!1})},1e3)}},a.handleSendClick=function(){a.sendInputMessage()},a.state={inputText:"",typing:!1},a}return Object(u.a)(t,e),Object(c.a)(t,[{key:"render",value:function(){return d.a.createElement("div",{className:"chat-controls connected fade-in"},d.a.createElement(z,{show:this.props.statusShow,text:this.props.statusText}),d.a.createElement(P,{handleDisconnect:this.props.handleDisconnect,handleSearch:this.props.handleSearch,disconnected:this.props.disconnected,matched:this.props.matched}),d.a.createElement(A,{matched:this.props.matched,onKeyDown:this.handleKeyDown,onChange:this.handleOnChange,handleLike:this.props.handleLike,liked:this.props.liked,reported:this.props.reported,handleSave:this.props.handleSave,handleReport:this.props.handleReport,searching:this.props.searching,reportModalOpen:this.props.reportModalOpen,handleReportModalClose:this.props.handleReportModalClose,handleReportOptionClick:this.props.handleReportOptionClick}),d.a.createElement("div",{className:"circle-wrapper send"},d.a.createElement("button",{className:"chat-send-button circle"+(this.props.disconnected||""===this.state.inputText?" disabled":""),onClick:this.handleSendClick,disabled:this.props.disconnected||""===this.state.inputText?" disabled":""},d.a.createElement("img",{className:"button-icon",src:R,alt:"Send"}))))}}]),t}(d.a.Component),P=function(e){function t(){return Object(o.a)(this,t),Object(i.a)(this,Object(l.a)(t).apply(this,arguments))}return Object(u.a)(t,e),Object(c.a)(t,[{key:"render",value:function(){return this.props.disconnected?d.a.createElement("div",{className:"circle-wrapper"},d.a.createElement("button",{className:"chat-next-button circle",onClick:this.props.handleSearch},d.a.createElement("img",{className:"button-icon",src:T,alt:"Next"}))):d.a.createElement("div",{className:"circle-wrapper"},d.a.createElement("button",{className:"chat-disconnect-button circle"+(this.props.matched?"":" disabled"),onClick:this.props.handleDisconnect},d.a.createElement("img",{className:"button-icon x",src:C,alt:"Disconnect"})))}}]),t}(d.a.Component),A=function(e){function t(){return Object(o.a)(this,t),Object(i.a)(this,Object(l.a)(t).apply(this,arguments))}return Object(u.a)(t,e),Object(c.a)(t,[{key:"render",value:function(){return this.props.searching?d.a.createElement("div",{className:"middle-buttons"}):this.props.matched?d.a.createElement("input",{type:"text",placeholder:"Message",className:"chat-input"+(this.props.disconnected?" chat-controls-disabled":""),onKeyDown:this.props.onKeyDown,onChange:this.props.onChange,maxLength:984,disabled:this.props.disconnected?" disabled":""}):d.a.createElement("div",{className:"middle-buttons"},d.a.createElement("div",{className:"circle-wrapper"},d.a.createElement("button",{className:"middle-button circle like-button"+(this.props.liked?"":" active"),onClick:this.props.handleLike},d.a.createElement("img",{className:"button-icon like",src:this.props.liked?W:L,alt:"Like"}))),d.a.createElement("div",{className:"circle-wrapper"},d.a.createElement(K,{open:this.props.reportModalOpen,onClose:this.props.handleReportModalClose,handleReportOptionClick:this.props.handleReportOptionClick,reported:this.props.reported}),d.a.createElement("button",{className:"middle-button circle report-button"+(this.props.reported?" active":""),onClick:this.props.handleReport},d.a.createElement("img",{className:"button-icon",src:this.props.reported?D:B,alt:"Report"}))))}}]),t}(d.a.Component),J=function(e){function t(){return Object(o.a)(this,t),Object(i.a)(this,Object(l.a)(t).apply(this,arguments))}return Object(u.a)(t,e),Object(c.a)(t,[{key:"render",value:function(){var e=new Date(this.props.timestamp),t=e.getHours(),a=e.getMinutes(),n=t>=12?"PM":"AM",s=(t=(t%=12)||12)+":"+(a=a<10?"0"+a:a)+" "+n;return d.a.createElement("span",{className:"message-timestamp"+("bottom"===this.props.direction?" timestamp-system":"")},s)}}]),t}(d.a.Component),K=function(e){function t(){return Object(o.a)(this,t),Object(i.a)(this,Object(l.a)(t).apply(this,arguments))}return Object(u.a)(t,e),Object(c.a)(t,[{key:"render",value:function(){var e=this;if(!this.props.open)return null;var t=""!==this.props.reported,a=Object.keys(x).map(function(a){return t?d.a.createElement("div",{className:"report-option"+(e.props.reported===a?" selected":" disabled"),"data-key":a,key:a},x[a]):d.a.createElement("div",{className:"report-option enabled","data-key":a,key:a,onClick:e.props.handleReportOptionClick},x[a])});return d.a.createElement("div",{className:"backdrop",onClick:this.props.onClose},d.a.createElement("div",{className:"report-modal"},d.a.createElement("div",{className:"report-header"},d.a.createElement("div",{className:"report-header-side"}),d.a.createElement("div",{className:"report-header-title"},"Report"),d.a.createElement("div",{className:"report-header-side"},d.a.createElement("img",{className:"report-x",src:M,alt:"Close",onClick:this.props.onClose}))),d.a.createElement("div",{className:"report-options"},a)))}}]),t}(d.a.Component)}).call(this,a(50))},function(e,t,a){e.exports=a(64)},,,,,function(e,t,a){},,,,,,,,,,,,,function(e,t,a){e.exports=a.p+"static/media/s.b117b408.svg"},function(e,t,a){e.exports=a.p+"static/media/crown2.b2b97c41.svg"},function(e,t,a){e.exports=a.p+"static/media/shield2.c405312f.svg"},function(e,t,a){e.exports=a.p+"static/media/people4.ce737683.svg"},,,,function(e,t,a){e.exports=a.p+"static/media/xWhite.d527b858.svg"},function(e,t,a){e.exports=a.p+"static/media/xGrey.2539a979.svg"},function(e,t,a){e.exports=a.p+"static/media/nextWhite2.ce906727.svg"},function(e,t,a){e.exports=a.p+"static/media/sendWhite.b7ed5a83.svg"},function(e,t,a){e.exports=a.p+"static/media/heartWhite.cbaf32fe.svg"},function(e,t,a){e.exports=a.p+"static/media/exclamationWhite.f1eeb750.svg"},function(e,t,a){e.exports=a.p+"static/media/heartBlueEmpty.a91080c9.svg"},function(e,t,a){e.exports=a.p+"static/media/heartBlueFilled.29ea6b8a.svg"},function(e,t,a){e.exports=a.p+"static/media/reportEmpty.cc62d953.svg"},function(e,t,a){e.exports=a.p+"static/media/reportFilled.a3d45f32.svg"},function(e,t,a){e.exports=a.p+"static/media/circle_checkmark.001302f6.svg"},function(e,t,a){"use strict";a.r(t);var n=a(0),s=a.n(n),r=a(15),o=a.n(r),c=(a(33),a(3)),i=a(4),l=a(6),u=a(5),p=a(7),h=a(2),d=a(69),m=a(19),b=a(66),v=function(e){function t(){return Object(c.a)(this,t),Object(l.a)(this,Object(u.a)(t).apply(this,arguments))}return Object(p.a)(t,e),Object(i.a)(t,[{key:"componentWillMount",value:function(){var e=this;this.unlisten=this.props.history.listen(function(t,a){e.props.onRouteChange(t)})}},{key:"componentWillUnmount",value:function(){this.unlisten()}},{key:"componentDidMount",value:function(){this.props.onRouteChange(this.props.location)}},{key:"render",value:function(){return this.props.children}}]),t}(s.a.Component),f=Object(b.a)(v),g=a(70),k=a(68),y=a(67),O=(a(45),a(13)),E=(a(47),a(48),a(49),function(e){function t(e){var a;return Object(c.a)(this,t),(a=Object(l.a)(this,Object(u.a)(t).call(this,e))).state={didMount:!1},a}return Object(p.a)(t,e),Object(i.a)(t,[{key:"componentDidMount",value:function(){var e=this;setTimeout(function(){e.setState({didMount:!0})},0)}},{key:"render",value:function(){this.props.height,this.props.width,this.state.didMount;return s.a.createElement("div",{class:"home scale"},s.a.createElement(O.a,{location:this.props.location}),s.a.createElement("div",{className:"home-main"},s.a.createElement("div",{className:"motto scale"},s.a.createElement("div",{className:"motto-main scale"},"Quality conversations"),s.a.createElement("div",{className:"motto-extra scale"},"Anonymous one-on-one chat, with the focus on the quality of the conversation")),s.a.createElement("div",{className:"jobs scale"},s.a.createElement("div",{className:"jobs-items scale"},s.a.createElement("div",{className:"jobs-circles"}),s.a.createElement("p",{className:"jobs-titles"},"Discover"),s.a.createElement("p",{className:"jobs-texts"},"New people and share your stories")),s.a.createElement("div",{className:"jobs-items scale"},s.a.createElement("div",{className:"jobs-circles"}),s.a.createElement("p",{className:"jobs-titles"},"Engage"),s.a.createElement("p",{className:"jobs-texts"},"In an interesting discussion or a simple conversation")),s.a.createElement("div",{className:"jobs-items scale"},s.a.createElement("div",{className:"jobs-circles"}),s.a.createElement("p",{className:"jobs-titles"},"Find"),s.a.createElement("p",{className:"jobs-texts"},"Wonderful experiences and relationships")))),s.a.createElement("div",{className:"footer"},s.a.createElement("span",{className:"footer-item"},s.a.createElement(y.a,{to:"/contact"},"Contact"))))}}]),t}(s.a.Component)),w=(s.a.Component,s.a.Component,a(27)),j=(s.a.Component,a(8)),x=a.n(j),N=a(10),S=a(63),C=function(e){function t(e){var a;return Object(c.a)(this,t),(a=Object(l.a)(this,Object(u.a)(t).call(this,e))).state={err:"",sent:!1},a.handleFormSubmit=a.handleFormSubmit.bind(Object(h.a)(Object(h.a)(a))),a}return Object(p.a)(t,e),Object(i.a)(t,[{key:"handleFormSubmit",value:function(){var e=Object(N.a)(x.a.mark(function e(t){var a,n,s,r,o,c;return x.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:if(t.preventDefault(),a=document.getElementById("name").value,n=document.getElementById("email").value,s=document.getElementById("message").value,""!=a){e.next=7;break}return this.setState({err:"Please provide your name"}),e.abrupt("return",!1);case 7:if(""!=n){e.next=10;break}return this.setState({err:"Please provide your email"}),e.abrupt("return",!1);case 10:if(""!=s){e.next=13;break}return this.setState({err:"Please provide your message"}),e.abrupt("return",!1);case 13:return"/api/contact",o={name:a,email:n,message:s},e.prev=15,e.next=18,fetch("/api/contact",{method:"post",headers:{"Content-Type":"application/json"},body:JSON.stringify(o)});case 18:return c=e.sent,e.next=21,c.json();case 21:r=e.sent,e.next=28;break;case 24:return e.prev=24,e.t0=e.catch(15),console.log(e.t0),e.abrupt("return",!1);case 28:if(!r.error){e.next=31;break}return this.setState({err:r.msg}),e.abrupt("return",!1);case 31:this.setState({err:"",sent:!0});case 32:case"end":return e.stop()}},e,this,[[15,24]])}));return function(t){return e.apply(this,arguments)}}()},{key:"render",value:function(){var e=this.state.err;return s.a.createElement("div",null,s.a.createElement(O.a,{location:this.props.location}),s.a.createElement("div",{className:"main-content"},s.a.createElement("div",{className:"auth-page"},s.a.createElement("h1",{className:"heading"},"Contact"),s.a.createElement("form",{className:"contact-form"+(this.state.sent?"":" visible")},s.a.createElement("input",{id:"name",type:"text",placeholder:"name*",className:"auth-input"}),s.a.createElement("input",{id:"email",type:"text",placeholder:"email*",className:"auth-input"}),s.a.createElement("textarea",{id:"message",placeholder:"message*",className:"auth-textarea",rows:"7"}),s.a.createElement("input",{type:"text",className:"error-auth ".concat(e&&"visible"),value:e,readOnly:!0}),s.a.createElement("button",{className:"submit-button",onClick:this.handleFormSubmit},"Contact")),s.a.createElement("div",{className:"contact-success"+(this.state.sent?" visible":"")},s.a.createElement("img",{className:"contact-success-checkmark",src:S,alt:"Success"}),s.a.createElement("h1",{className:"contact-success-heading"},"Thank you for your feedback")))))}}]),t}(s.a.Component),M=(s.a.Component,function(e){function t(){return Object(c.a)(this,t),Object(l.a)(this,Object(u.a)(t).apply(this,arguments))}return Object(p.a)(t,e),Object(i.a)(t,[{key:"render",value:function(){return 0===this.props.chats.length?s.a.createElement("div",{className:"no-history"},"You don't have any chats yet"):s.a.createElement("div",{className:"chats-rows"},this.props.chats.map(function(e){return s.a.createElement("div",{className:"chat-row"},e.messages[0].timestamp," | ",e.duration," | ",e.result.liked?"Liked":"Not liked"," | ",e.result.reported?e.result.reported:"Not reported")}))}}]),t}(s.a.Component)),T=function(e){function t(){return Object(c.a)(this,t),Object(l.a)(this,Object(u.a)(t).apply(this,arguments))}return Object(p.a)(t,e),Object(i.a)(t,[{key:"render",value:function(){return s.a.createElement("div",null,s.a.createElement(O.a,{location:this.props.location}),s.a.createElement("div",{className:"not-found"},s.a.createElement("h1",null,"Page not found")))}}]),t}(s.a.Component),R=function(e){function t(){return Object(c.a)(this,t),Object(l.a)(this,Object(u.a)(t).apply(this,arguments))}return Object(p.a)(t,e),Object(i.a)(t,[{key:"render",value:function(){var e=this;return s.a.createElement(g.a,{location:this.props.location},s.a.createElement(k.a,{path:"/",exact:!0,render:function(){return s.a.createElement(E,{height:e.props.height,width:e.props.width})},location:this.props.location}),s.a.createElement(k.a,{path:"/chat",exact:!0,render:function(){return s.a.createElement(w.a,{height:e.props.height,width:e.props.width,location:e.props.location})}}),s.a.createElement(k.a,{path:"/contact",exact:!0,render:function(){return s.a.createElement(C,{location:e.props.location})}}),s.a.createElement(k.a,{component:T}))}}]),t}(s.a.Component),L=function(e){function t(e){var a;return Object(c.a)(this,t),(a=Object(l.a)(this,Object(u.a)(t).call(this,e))).state={location:window.location,height:window.innerHeight,width:window.innerWidth},a.onResize=a.onResize.bind(Object(h.a)(Object(h.a)(a))),a}return Object(p.a)(t,e),Object(i.a)(t,[{key:"componentDidMount",value:function(){var e=.01*window.innerHeight;document.documentElement.style.setProperty("--vh","".concat(e,"px")),window.addEventListener("resize",this.onResize)}},{key:"componentWillUnmount",value:function(){window.removeEventListener("resize",this.onResize)}},{key:"onResize",value:function(){var e=.01*window.innerHeight;document.documentElement.style.setProperty("--vh","".concat(e,"px")),this.setState({width:window.innerWidth,height:window.innerHeight})}},{key:"handleRouteChange",value:function(e){this.setState({location:e})}},{key:"render",value:function(){var e=this;return s.a.createElement(d.a,null,s.a.createElement(f,{onRouteChange:function(t){return e.handleRouteChange(t)}},s.a.createElement(m.TransitionGroup,null,s.a.createElement(m.CSSTransition,{key:this.state.location.pathname,classNames:"fade",in:!1,appear:!1,timeout:{appear:100,enter:100,exit:100}},s.a.createElement(R,{height:this.state.height,width:this.state.width,location:this.state.location})))))}}]),t}(s.a.Component);Boolean("localhost"===window.location.hostname||"[::1]"===window.location.hostname||window.location.hostname.match(/^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/));o.a.render(s.a.createElement(L,null),document.getElementById("root")),"serviceWorker"in navigator&&navigator.serviceWorker.ready.then(function(e){e.unregister()})}],[[28,2,1]]]);
//# sourceMappingURL=main.07bff618.chunk.js.map