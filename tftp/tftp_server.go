package tftp

import (
	"fmt"
	"net"
)

func _new_client_handler(client *net.UDPAddr, req *tftp_packet) {
	/**declare new socket used to communicate with client*/
	socket, err := net.DialUDP("udp", nil, client)
	if err != nil {
		fmt.Println("connect to client failed，err:", err)
		return
	}
	defer socket.Close()

	var buffer [4096]byte
	var dataBuffer = make([]byte, 512)
	var fmgr FileOp = &FileMGR{}

	// open file
	fd := fmgr.open(req.fileName)
	defer fmgr.close(fd)

	if req.op == TFTP_OP_RRQ {
		var blockNum uint16 = 0

		fmt.Println("Read request")

		for {
			size := fmgr.read(fd, dataBuffer)
			blockNum++
			n := _pack_msg_data(buffer[:], blockNum, dataBuffer[:size])
			socket.Write(buffer[:n])

			// wait ack
			length, err := socket.Read(buffer[:]) // 接收数据
			if err != nil {
				fmt.Println("read udp failed, err:", err)
				break
			}
			_, msgType, msg := parse_raw_data(buffer[:length])
			if msgType == TFTP_OP_ACK && msg.blockNum == blockNum {

			} else {
				//error
				fmt.Println("Wrong response")
				n := _pack_msg_err(buffer[:], TFTP_ERROR_Undefined, "Wrong response")
				socket.Write(buffer[:n])
				break
			}

			// exit entry
			if size != len(dataBuffer) {
				break
			}
		}

	} else if req.op == TFTP_OP_WRQ {
		var blockNum uint16 = 0

		fmt.Println("Write request")

		for {
			n := _pack_msg_ack(buffer[:], blockNum)
			socket.Write(buffer[:n])
			blockNum++

			// wait data
			length, err := socket.Read(buffer[:]) // 接收数据
			if err != nil {
				fmt.Println("read udp failed, err:", err)
				break
			}
			_, msgType, msg := parse_raw_data(buffer[:length])
			if msgType == TFTP_OP_DATA {
				fmt.Printf("PACKET#%v size: %d/%d\n", msg.blockNum, len(msg.data), len(dataBuffer))
				fmgr.write(fd, msg.data)
			} else {
				//error
				fmt.Println("Wrong response")
				n := _pack_msg_err(buffer[:], TFTP_ERROR_Undefined, "Wrong response")
				socket.Write(buffer[:n])
				break
			}

			if len(msg.data) != len(dataBuffer) {
				n := _pack_msg_ack(buffer[:], blockNum)
				socket.Write(buffer[:n])
				break
			}
		}

	} else {
		return
	}

}

func Server(root string) {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: TFTP_PORT,
	})
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	} else {
		fmt.Printf("listening on port:%v\n", TFTP_PORT)
	}
	defer listen.Close()

	/* set root path */
	__file_root_path__ = root

	for {
		var data [1024]byte
		n, addr, err := listen.ReadFromUDP(data[:]) // 接收数据
		if err != nil {
			fmt.Println("read udp failed, err:", err)
			continue
		}

		fmt.Printf("received: addr:%v count:%v\n", addr, n)
		fmt.Printf("data: %v\n", data[:n])

		_, msgType, msg := parse_raw_data(data[:])

		if msgType == TFTP_OP_WRQ || msgType == TFTP_OP_RRQ {
			fmt.Printf("A new tftp request has been received (%d\n", msgType)
			go _new_client_handler(addr, msg)
		} else {
			fmt.Println("Unrecongnized message")
		}
	}
}
