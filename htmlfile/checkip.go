package htmlfile

import (
	"github.com/kprc/nbsnetwork/tools"
	"log"
)

func NewCheckIPFile(filename string) {

	checkiphtml := `<html>
<head>
<title>check a ip address is a world wide address</title>
</head>
<body>
<form action="/checkip" method="post">
	<h1>IP Address:</h1><input type="text" name="ipaddr">
<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<input type="submit" value="提交">
</form>
</body>
</html>`

	if err := tools.Save2File([]byte(checkiphtml), filename); err != nil {
		log.Fatal("can't generator a checkip html")
	}

}
