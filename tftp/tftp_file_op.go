package tftp

import "os"

type FileMGR struct {
	file *os.File
}

var __file_root_path__ string = "."

func (ins *FileMGR) open(fileName string) (fd int) {
	file, _ := os.OpenFile(__file_root_path__+"/"+fileName, os.O_RDWR|os.O_CREATE, 0664)
	ins.file = file

	return 1
}

func (ins *FileMGR) read(fd int, buf []byte) (len int) {
	n, _ := ins.file.Read(buf)
	return n
}
func (ins *FileMGR) write(fd int, buf []byte) (len int) {
	n, _ := ins.file.Write(buf)
	return n
}
func (ins *FileMGR) close(fd int) {
	ins.file.Sync()
	ins.file.Close()
}
