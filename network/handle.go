package network

func Handle(command int, conn Conn, pack *Package, handle func(p *Package) string) bool {
	if command != pack.Command {
		return false
	}
	_, err := conn.Write([]byte(serializePackage(&Package{
		Command: command,
		Data:    handle(pack),
	}) + EndBytes))
	if err != nil {
		return false
	}

	return true
}
