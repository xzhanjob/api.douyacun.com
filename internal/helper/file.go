package helper

import "os"

func FileExists(filename string) bool {
	info, err := os.Stat(filename)

	return err == nil && !info.IsDir()
}

// Overwrite overwrites the file with data. Creates file if not present.
func FileOverwrite(fileName string, data []byte) bool {
	f, err := os.Create(fileName)
	if err != nil {
		return false
	}

	_, err = f.Write(data)
	return err == nil
}