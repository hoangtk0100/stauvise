package util

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func ToHLS(ctx context.Context, src, dst, dir string) error {
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", src, "-c", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", "-hls_segment_filename", dir+"_%d.ts", dst)
	return cmd.Run()
}

func GetMaxSegmentNumber(dir string) (max int, err error) {
	file, err := os.Open(dir)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	num := strings.ReplaceAll(strings.Split(lines[len(lines)-2], "_")[1], ".ts", "")
	return strconv.Atoi(num)
}
