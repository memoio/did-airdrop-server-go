package main

import "github.com/did-server/cmd"

// @title			DID-Server API
// @version		1.0
// @description	did
// @description	1. get createsigmsg
// @description 2. create/createadmin (createadmin is free and not need createsigmsg)
// @description 3. exist (confirm did exist)
// @description 4. info (get did info)
// @description	mfile (file did)
// @description 1. create (create file did)
// @description 2. confirm (confirm file did)
// @description 3. download (download )
// @Host			localhost:8080
// @BasePath		/
func main() {
	cmd.Exceute()
}
