package tftp

import "encoding/binary"

func _pack_msg_rrq(buffer []byte, filename string, mode string) (bytes uint16) {
	return 0
}
func _pack_msg_wrq(buffer []byte, filename string, mode string) (bytes uint16) {
	return 0
}
func _pack_msg_data(buffer []byte, blk uint16, data []byte) (bytes int) {
	binary.BigEndian.PutUint16(buffer[:2], TFTP_OP_DATA)
	binary.BigEndian.PutUint16(buffer[2:4], blk)
	copy(buffer[4:], data[:])
	return 4 + len(data)
}
func _pack_msg_ack(buffer []byte, blk uint16) (bytes uint16) {
	binary.BigEndian.PutUint16(buffer[:2], TFTP_OP_ACK)
	binary.BigEndian.PutUint16(buffer[2:4], blk)
	return 4
}
func _pack_msg_err(buffer []byte, errorcode int, errormsg string) (bytes int) {
	binary.BigEndian.PutUint16(buffer[:2], TFTP_OP_ERROR)
	binary.BigEndian.PutUint16(buffer[2:4], uint16(errorcode))
	copy(buffer[4:], errormsg)
	buffer[4+len(errormsg)] = 0

	return 4 + len(errormsg) + 1
}

func parse_raw_data(data []byte) (success bool, msgType int, msg *tftp_packet) {
	var op uint16
	var frame tftp_packet

	op = binary.BigEndian.Uint16(data[0:2])
	frame.op = op

	switch op {
	case TFTP_OP_RRQ, TFTP_OP_WRQ:
		index := 2
		for {
			if data[index] == 0 {
				break
			}
			index++
		}
		frame.fileName = string(data[2:index])

		_index := index
		for {
			index++
			if data[index] == 0 {
				break
			}
		}
		frame.mode = string(data[_index:index])

	case TFTP_OP_DATA:
		frame.blockNum = binary.BigEndian.Uint16(data[2:4])
		frame.data = data[4:]
	case TFTP_OP_ACK:
		frame.blockNum = binary.BigEndian.Uint16(data[2:4])
	case TFTP_OP_ERROR:
		frame.errorCode = binary.BigEndian.Uint16(data[2:4])
		frame.errorMsg = string(data[4:])
	default:
		println("Unkonw op")
		return false, TFTP_OP_INVALID, nil
	}

	return true, int(op), &frame
}
