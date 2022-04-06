buffer := make([]byte, 1024)

	for {
		conn, _ := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
		count, _ := conn.Read(buffer)
		conn.Write(bytes.ToUpper(buffer[:count]))
		conn.Close()
	}