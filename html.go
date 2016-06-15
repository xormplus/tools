package main

const indexhtml string = `<html>
  <head>
    <title>file name input</title>
	<link type="image/x-icon" href="go.ico" rel="shortcut icon">
	<style>

    body { flow:vertical; border-spacing:8px; }

    div.list
    {
	  background:white white #D5D5D5 #D5D5D5;
      border:1px solid;
      width:*;
      height:*;
      overflow:scroll-indicator;
    }

    </style>
    <script type="text/tiscript">
		view.minSize = (1000,400);
		var msg = {
			A0001: "&#x8BF7;&#x9009;&#x62E9;&#x914D;&#x7F6E;&#x6587;&#x4EF6;&#x76EE;&#x5F55;&#x548C;&#x52A0;&#x5BC6;&#x6587;&#x4EF6;&#x8F93;&#x51FA;&#x76EE;&#x5F55;",
			A0002: "&#x8BF7;&#x9009;&#x62E9;&#x4E0D;&#x540C;&#x76EE;&#x5F55;",
			A0003: "&#x8BF7;&#x8BBE;&#x7F6E;&#x79D8;&#x94A5;&#x957F;&#x5EA6;",
			A0004: "&#x5904;&#x7406;&#x5B8C;&#x6BD5;&#xFF0C;&#x8BF7;&#x67E5;&#x770B;&#x65E5;&#x5FD7;",
			A0005: "&#x79D8;&#x94A5;&#x957F;&#x5EA6;&#x8D85;&#x8FC7;&#x6700;&#x5927;&#x9608;&#x503C;2048&#xFF0C;&#x8BF7;&#x91CD;&#x65B0;&#x8F93;&#x5165;",
			A0006: "&#x79D8;&#x94A5;&#x957F;&#x5EA6;&#x6700;&#x5C0F;&#x4E3A;256&#xFF0C;&#x8BF7;&#x91CD;&#x65B0;&#x8F93;&#x5165;",
			N0001: "&#x9519;&#x8BEF;&#x63D0;&#x793A;"
		};

        $(button.selector1).onClick = function () {
        	var fn = view.selectFolder();
        	if( fn ) $(.path1).value = fn;
        };

		$(button.selector2).onClick = function () {
        	var fn = view.selectFolder();
        	if( fn ) $(.path2).value = fn;
        };

		$(#btn).onClick = function () {
			this.attributes["disabled"]="disabled"
			var path1=$(.path1).value
			var path2=$(.path2).value
			if ( path1=="" || path2=="" ) {
				view.msgbox(#alert,msg.A0001,msg.N0001);
				this.attributes.remove("disabled");
				return;
			}
			if ( path1==path2 ) {
				view.msgbox(#alert,msg.A0002,msg.N0001);
				this.attributes.remove("disabled");
				return;
			}
			if ($(form).value.radioGroup==4 && $(form).value.bitwd==""){
				view.msgbox(#alert,msg.A0003,msg.N0001);
				this.attributes.remove("disabled");
				return;
			}
			if ($(form).value.radioGroup==4){
				if ($(form).value.bitwd>2048){
					view.msgbox(#alert,msg.A0005,msg.N0001);
					this.attributes.remove("disabled");
					return;
				}
				if ($(form).value.bitwd<256){
					view.msgbox(#alert,msg.A0006,msg.N0001);
					this.attributes.remove("disabled");
					return;
				}
			}

        	var res=view.getDir($(form).value);
			if (res["cmd"]=="done" ){
				view.msgbox(#alert,msg.A0004,msg.N0001);
			}
        };

		$(form).on("change",function() {
			var radioGroupValue =this.value.radioGroup;
			if (radioGroupValue==1 || radioGroupValue==2 || radioGroupValue==3){
				$(#pass).style["display"]="block";
				$(#bit).style["display"]="none";
			}else{
				$(#pass).style["display"]="none";
				$(#bit).style["display"]="block";
			}
		});

		function enable(){
			$(#btn).attributes.remove("disabled");
		}

    </script>
  </head>
<body>
	<div>
 <form .table>
	<div>
	&#x914D;&#x7F6E;&#x6587;&#x4EF6;&#x5B58;&#x653E;&#x76EE;&#x5F55;
	<input type="text" name="path1" class="path1" size="80" readonly/><button class="selector1">&#x2026;</button><br/>
	&#x52A0;&#x5BC6;&#x6587;&#x4EF6;&#x8F93;&#x51FA;&#x76EE;&#x5F55;
	<input type="text" name="path2" class="path2" size="80" readonly/><button class="selector2">&#x2026;</button><br/>
	&#x52A0;&#x5BC6;&#x7B97;&#x6CD5;<br/>
	<input type="radio" name="radioGroup" value="1" checked="checked">AES</input>
	<input type="radio" name="radioGroup" value="2">DES</input>
	<input type="radio" name="radioGroup" value="3">TripleDes</input>
	<input type="radio" name="radioGroup" value="4">RSA</input>

	<input type="button" id="btn" class="selector3" value="&#x52A0;&#x5BC6;"/><br/>
	</div>
	<div id="pass" style="display:true">&#x5BC6;&#x7801;<input type="text" name="passwd" class="passwd" size="30"/></div>
	<div id="bit" style="display:none">
		&#x79D8;&#x94A5;&#x957F;&#x5EA6;<input type="text" name="bitwd" class="bitwd" size="8" filter="0~9" maxlength="4"/>
		<input type="radio" name="rsaMode" value="1" checked="checked">&#x516C;&#x94A5;&#x52A0;&#x5BC6;&#x79C1;&#x94A5;&#x89E3;&#x5BC6;&#x6A21;&#x5F0F;</input>
		<input type="radio" name="rsaMode" value="2">&#x79C1;&#x94A5;&#x52A0;&#x5BC6;&#x516C;&#x94A5;&#x89E3;&#x5BC6;&#x6A21;&#x5F0F;</input>
	</div>
 </form>
	</div>

	<div id="result" class="list"></div>

</body>
</html>
`
