package helper

import "testing"

func Test_File_TopDir(t *testing.T) {
	var res string

	var paths = [][]string{
		{"./explain.md", ""},
		{"./mysql/explain", "mysql"},
		{"/mysql/explain.md", "mysql"},
		{"mysql/explain.md", "mysql"},
		{"mysql/explain.php", "mysql"},
		{"/linux/sync_fsync_fdatasync.md", "linux"},
		{"./aof.md#1", ""},
	}
	for _, v := range paths {
		if res = File.TopDir(v[0]); res != v[1] {
			t.Fatalf("路径：%s 预期：\"\"; 实际：%s", v[0], v[1])
		}
	}
}
