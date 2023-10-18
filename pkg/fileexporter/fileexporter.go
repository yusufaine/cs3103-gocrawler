package fileexporter

import "os"

func WriteToFile(d []byte, filename string) error {
	var f *os.File
	// if file exists, overwrite
	if _, err := os.Stat(filename); err == nil {
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(d)
	f.Write([]byte("\n"))
	return nil
}
