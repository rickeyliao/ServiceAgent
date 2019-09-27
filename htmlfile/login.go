package htmlfile

import (
	"github.com/kprc/nbsnetwork/tools"

	"log"
)

func NewLoginFile(filename string) {

	loginhtml := `<html>
<head>
<title></title>
</head>
<body align="left">
<form action="/login" method="post">
	用户名:<input type="text" name="username">
	密码:<input type="password" name="password">
	
<br>
地&nbsp;&nbsp;&nbsp;&nbsp;址:<input type="text" name="pubaddr" width="300">

<br><br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<input type="submit" value="提交">
</form>
</body>
</html>`

	if err := tools.Save2File([]byte(loginhtml), filename); err != nil {
		log.Fatal("can't generator a login html")
	}

}
