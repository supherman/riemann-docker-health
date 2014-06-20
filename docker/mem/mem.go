package mem

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	LXC_MEM_DIR        = "/sys/fs/cgroup/memory/docker"
	LXC_MEM_USAGE_FILE = "memory.usage_in_bytes"
	LXC_MEM_LIMIT_FILE = "memory.limit_in_bytes"
)

func GetMetric(id string, metric string) (int64, error) {
	path := fmt.Sprintf("%s/%s/%s", LXC_MEM_DIR, id, metric)
	f, err := os.Open(path)
	if err != nil {
		log.Println("Error while opening : ", err)
		return 0, err
	}
	defer f.Close()

	buffer := make([]byte, 16)
	n, err := f.Read(buffer)
	if err != nil {
		log.Println("Error while reading ", path, " : ", err)
		return 0, err
	}

	buffer = buffer[:n-1]
	val, err := strconv.ParseInt(string(buffer), 10, strconv.IntSize)
	if err != nil {
		log.Println("Error while parsing ", string(buffer), " : ", err)
		return 0, err
	}

	return val, nil
}

func GetUsage(id string)(int64, error) {
  return GetMetric(id, LXC_MEM_USAGE_FILE)
}

func GetLimit(id string)(int64, error) {
  return GetMetric(id, LXC_MEM_LIMIT_FILE)
}

func GetPercentage(id string)(int64, error) {
  usage, error := GetUsage(id)
  limit, error := GetLimit(id)
  percentage := float64(usage) * 100.0 / float64(limit)
  return int64(percentage), error
}
