package tftp

const TFTP_PORT = 16969

const (
	TFTP_OP_INVALID = 0
	TFTP_OP_RRQ     = 1
	TFTP_OP_WRQ     = 2
	TFTP_OP_DATA    = 3
	TFTP_OP_ACK     = 4
	TFTP_OP_ERROR   = 5
)

const (
	TFTP_MODE_BINARY = "octet"
	TFTP_MODE_ASCII  = "netascii"
)

const (
	TFTP_ERROR_Undefined         = 0
	TFTP_ERROR_FileNotFound      = 1
	TFTP_ERROR_InvalidAccess     = 2
	TFTP_ERROR_DiskFull          = 3
	TFTP_ERROR_IllegalOperation  = 4
	TFTP_ERROR_UnknownTransferID = 5
	TFTP_ERROR_FileAlreadyExists = 6
	TFTP_ERROR_NoSuchUser        = 7
)

type any = interface{}

type FileOp interface {
	open(fileName string) (fd int)
	read(fd int, buf []byte) (len int)
	write(fd int, buf []byte) (len int)
	close(fd int)
}

type tftp_packet struct {
	op        uint16
	blockNum  uint16
	fileName  string
	mode      string
	data      []byte
	errorCode uint16
	errorMsg  string
}

func init() {
	println("init")
}
